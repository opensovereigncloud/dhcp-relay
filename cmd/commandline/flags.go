package commandline

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opensovereigncloud/dhcp-relay/internal/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Params struct {
	KeaEndpoint string
	NicPrefix   string
	PidFile     string
	LogParams   log.Params
}

func ParseArgs() Params {
	pflag.Usage = usage
	pflag.ErrHelp = nil
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	pflag.String("kea-endpoint", "", "Kea DHCP endpoint")
	pflag.String("nic-prefix", "Ethernet", "NIC name prefix for filtering")
	pflag.String("pid-file", "/tmp/dhcrelay.pid", "PID file path")
	pflag.String("log-level", "info", "Log level. Valid values: debug, info, warn, error")
	pflag.String("log-format", "text", "Log format. Valid values: text, json")

	var help bool
	pflag.BoolVarP(&help, "help", "h", false, "Show this help message.")
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		exitUsage(err)
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
	err = pflag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		exitUsage(err)
	}
	if help {
		exitUsage(nil)
	}

	return Params{
		KeaEndpoint: viper.GetString("kea-endpoint"),
		NicPrefix:   viper.GetString("nic-prefix"),
		PidFile:     viper.GetString("pid-file"),
		LogParams: log.Params{
			Level:  viper.GetString("log-level"),
			Format: viper.GetString("log-format"),
		},
	}
}

func usage() {
	name := filepath.Base(os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [--option]...\n", name)
	_, _ = fmt.Fprintf(os.Stderr, "Options:\n")
	pflag.PrintDefaults()
}

func exitUsage(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", filepath.Base(os.Args[0]), err)
	}
	pflag.Usage()
	os.Exit(2)
}
