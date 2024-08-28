package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/opensovereigncloud/dhcp-relay/cmd/commandline"
	"github.com/opensovereigncloud/dhcp-relay/internal/log"
	"github.com/opensovereigncloud/dhcp-relay/internal/service"
)

func main() {
	params := commandline.ParseArgs()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger := log.Setup(params.LogParams)
	svc := service.New(params.KeaEndpoint, params.NicPrefix, params.PidFile, logger)
	if err := svc.Run(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
