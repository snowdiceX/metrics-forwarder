package main

import (
	"context"
	"os"
	"strings"

	"github.com/snowdiceX/metrics-forwarder/config"
	"github.com/snowdiceX/metrics-forwarder/log"
	"github.com/spf13/cobra"
)

const (
	// CommandStart cli command "start"
	CommandStart = "start"

	// CommandVersion cli command "version"
	CommandVersion = "version"
)

// Runner of the command task
type Runner func(conf *config.Config) (context.CancelFunc, error)

// NewRootCommand create root/default command
func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "metrics-forwarder",
		Short: "pull the metrics and push them to pushgateway of prometheus",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if strings.EqualFold(cmd.Use, CommandVersion) {
				return nil
			}
			logger, err := log.LoadLogger(config.GetConfig().LogConfigPath)
			if err != nil {
				log.Warn("Used the default logger because error: ", err)
			} else {
				log.Replace(logger)
			}
			// // init config
			// _, err = config.LoadConfig(config.GetConfig().ConfigFile)
			// if err != nil {
			// 	log.Error("Run root command error: ", err.Error())
			// 	return err
			// }
			// log.Debug("Init config: ", config.GetConfig().ConfigFile)
			return nil
		},
	}
	return root
}

func commandRunner(run Runner, isKeepRunning bool) error {
	cancel, err := run(config.GetConfig())
	if err != nil {
		log.Error("Run command error: ", err.Error())
		return err
	}
	if isKeepRunning {
		KeepRunning(func(sig os.Signal) {
			defer log.Flush()
			if cancel != nil {
				cancel()
			}
			log.Debug("Stopped by signal: ", sig)
		})
	}
	return nil
}
