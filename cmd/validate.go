package cmd

import (
	"github.com/notwillk/nara/internal/compiler"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [fileGlob...]",
		Short: "Validate specific files",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := compiler.Compile(configPath, args)
			return err
		},
	}
}
