package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [fileGlob...]",
		Short: "Validate specific files",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = configPath
			_ = args
			return fmt.Errorf("validate not implemented yet")
		},
	}
}

