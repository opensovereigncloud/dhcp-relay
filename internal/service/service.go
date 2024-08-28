package service

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"syscall"
	"time"
)

const (
	dhcrelayBin   = "/usr/sbin/dhcrelay"
	inet6FilePath = "/proc/net/if_inet6"
)

type logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warn(err error, msg string, keysAndValues ...any)
	Error(err error, msg string, keysAndValues ...any)
}

type RelayService struct {
	keaEndpoint  string
	nicPrefix    string
	pidFile      string
	listenString []string
	p            *os.Process
	log          logger
}

func New(keaEndpoint, nicPrefix, pidFile string, log logger) *RelayService {
	return &RelayService{
		keaEndpoint: keaEndpoint,
		nicPrefix:   nicPrefix,
		pidFile:     pidFile,
		log:         log,
	}
}

func (svc *RelayService) Run(ctx context.Context) error {
	procContext, cancel := context.WithCancel(ctx)
	defer cancel()

	go svc.loop(procContext)

	<-ctx.Done()
	return nil
}

func (svc *RelayService) loop(ctx context.Context) {
	svc.log.Info("starting relay service in 5 seconds")
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			err := svc.cleanupDHCPRelay()
			if err != nil {
				svc.log.Warn(err, "failed to cleanup DHCP relay process")
			}
			svc.log.Info("stopped relay service")
			return
		case <-t.C:
			diff, err := svc.setListenString()
			if err != nil {
				svc.log.Warn(err, "failed to set listen string")
				continue
			}
			if diff {
				_ = svc.cleanupDHCPRelay()
				continue
			}
			if err := svc.ensureDHCPRelay(); err != nil {
				svc.log.Warn(err, "failed to ensure DHCP relay process")
			}
		}
	}
}

func (svc *RelayService) ensureDHCPRelay() error {
	if svc.p == nil {
		return svc.tryRunDHCPRelay()
	}
	if err := svc.p.Signal(syscall.Signal(0)); err == nil {
		return nil
	}
	return svc.cleanupDHCPRelay()
}

func (svc *RelayService) tryRunDHCPRelay() error {
	args, err := svc.dhcpRelayArgs()
	if err != nil {
		return err
	}

	var attrs os.ProcAttr
	attrs.Files = []*os.File{os.Stdin, os.Stdout, os.Stderr}
	svc.log.Debug("starting DHCP relay process", "args", strings.Join(args, " "))
	p, err := os.StartProcess(dhcrelayBin, args, &attrs)
	if err != nil {
		return err
	}
	svc.p = p
	return nil
}

func (svc *RelayService) dhcpRelayArgs() ([]string, error) {
	if len(svc.listenString) == 0 {
		return nil, fmt.Errorf("no listen string provided")
	}
	args := []string{"-6", "-pf", svc.pidFile, "-I", "-u", svc.keaEndpoint}
	args = append(args, svc.listenString...)
	return args, nil
}

func (svc *RelayService) setListenString() (bool, error) {
	if svc.listenString == nil {
		svc.listenString = make([]string, 0)
	}
	listenString := make([]string, 0)
	f, err := os.Open(inet6FilePath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	regexString := fmt.Sprintf("(?P<addr>[a-f,0-9]{32}).*(?P<nic>%s[0-9]{1,3})", svc.nicPrefix)
	re, err := regexp.Compile(regexString)
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !re.MatchString(line) {
			continue
		}
		matches := re.FindStringSubmatch(line)
		addrIndex := re.SubexpIndex("addr")
		addr := matches[addrIndex]
		if !strings.HasPrefix(addr, "fe80") {
			continue
		}
		nicIndex := re.SubexpIndex("nic")
		nic := matches[nicIndex]
		listenString = append(listenString, []string{"-l", nic}...)
	}
	if len(listenString) == 0 {
		svc.listenString = listenString
		return false, fmt.Errorf("no interfaces to listen found")
	}
	if reflect.DeepEqual(listenString, svc.listenString) {
		return false, nil
	}
	svc.listenString = listenString
	return true, nil
}

func (svc *RelayService) cleanupDHCPRelay() error {
	if svc.p == nil {
		return nil
	}
	if err := svc.p.Kill(); err != nil {
		return err
	}
	if err := os.Remove(svc.pidFile); err != nil {
		return err
	}
	svc.p = nil
	return nil
}
