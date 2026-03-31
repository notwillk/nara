package cmd

import (
	"github.com/notwillk/nara/internal/scaffold"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize project",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return scaffold.Init(configPath)
		},
	}
}
