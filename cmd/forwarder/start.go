package main

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/snowdiceX/metrics-forwarder/collector"
	"github.com/snowdiceX/metrics-forwarder/config"
	"github.com/snowdiceX/metrics-forwarder/log"
	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&config.GetConfig().LogConfigPath, "log", "./config/log.conf", "log config file path")
}

// NewStartCommand create start command
func NewStartCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   CommandStart,
		Short: "start relay service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}

	addFlags(cmd)
	return cmd
}

func starter(conf *config.Config) (context.CancelFunc, error) {
	// Since we are dealing with custom Collector implementations, it might
	// be a good idea to try it out with a pedantic registry.
	reg := prometheus.NewPedanticRegistry()

	// Construct forwarder collector.
	collector := collector.NewForwarderCollector(
		"mainnet", "78", "http://127.0.0.1:26660/metrics", reg)

	tick := time.NewTicker(time.Millisecond * 30000)
	for range tick.C {
		if err := push.New("http://127.0.0.1:9091", "irishub").
			Collector(collector).
			Grouping("service", "blockchain").
			Push(); err != nil {
			log.Errorf("Could not push metrics to pushgateway: %d", err)
		}
	}
	return nil, nil
}
