package main

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"

	"github.com/snowdiceX/metrics-forwarder/collector"
	"github.com/snowdiceX/metrics-forwarder/config"
	"github.com/snowdiceX/metrics-forwarder/log"
	"github.com/spf13/cobra"
)

// MultiFlagVar check for multiple settings for a flag
type MultiFlagVar struct {
	Values []string
}

func (f *MultiFlagVar) String() string {
	return fmt.Sprint(f.Values)
}

// Set a flag value
func (f *MultiFlagVar) Set(value string) error {
	f.Values = append(f.Values, value)
	return nil
}

// Type return type string
func (f *MultiFlagVar) Type() string {
	return "[]string"
}

func addFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&config.GetConfig().LogConfigPath,
		"log", "./log.conf", "log config file path")
	cmd.Flags().StringVar(&config.GetConfig().ConfigPath,
		"config", "./config.conf", "config file path")
	// cmd.Flags().Var(&pullVar, "pull", "pull metrics url")
	// cmd.Flags().StringVar(&config.GetConfig().Push,
	// 	"push", "http://127.0.0.1:9091", "push metrics url")
	// cmd.Flags().StringVar(&config.GetConfig().Job,
	// 	"job", "push", "push job name")
	// // cmd.Flags().StringVar(&config.GetConfig().Instance,
	// // 	"instance", "pushgateway", "instance")
	// cmd.Flags().StringVar(&config.GetConfig().Zone,
	// 	"zone", "east", "zone of server")
	// cmd.Flags().StringVar(&config.GetConfig().Host,
	// 	"host", "127.0.0.1", "host ip")
	// cmd.Flags().StringVar(&config.GetConfig().Group,
	// 	"group", "service", "group label name")
	// cmd.Flags().StringVar(&config.GetConfig().GroupValue,
	// 	"group_value", "blockchain", "group label value")
	// cmd.Flags().Uint32Var(&config.GetConfig().Ticker,
	// 	"ticker", 30000, "time ticker")
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
	err := conf.Load()
	if err != nil {
		log.Warn("Load config error: ", err.Error())
	}
	for i := 0; i < len(conf.Jobs); i++ {
		job := conf.Jobs[i]
		go func() {
			// Since we are dealing with custom Collector implementations, it might
			// be a good idea to try it out with a pedantic registry.
			reg := prometheus.NewPedanticRegistry()

			// Construct forwarder collector.
			collector := collector.NewForwarderCollector(
				job.Zone, job.Host, job.Pull, reg)

			tick := time.NewTicker(time.Millisecond * time.Duration(job.Ticker))
			for range tick.C {
				log.Info("tick...")
				if err := push.New(conf.Push, job.Name).
					Collector(collector).
					Grouping(job.Group, job.GroupValue).
					Push(); err != nil {
					log.Errorf("Could not push metrics to pushgateway: %d", err)
				}
			}
		}()
	}
	return nil, nil
}
