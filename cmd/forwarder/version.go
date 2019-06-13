package main

import (
	"context"
	"fmt"

	"github.com/snowdiceX/metrics-forwarder/config"

	"github.com/spf13/cobra"
)

var (
	// Version of forwarder
	Version = "0.0.0"

	// GitCommit is the current HEAD set using ldflags.
	GitCommit string
)

// NewVersionCommand create version command
func NewVersionCommand(run Runner, isKeepRunning bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show version info",
		RunE: func(cmd *cobra.Command, args []string) error {
			return commandRunner(run, isKeepRunning)
		},
	}
	return cmd
}

var versioner = func(conf *config.Config) (context.CancelFunc, error) {

	s := `Version:	%s
GitCommitID:	%s
`

	fmt.Printf(s, Version, GitCommit)

	return nil, nil
}
