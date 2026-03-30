package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newFormatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "format",
		Short: "Update config schema only",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = configPath
			return fmt.Errorf("format not implemented yet")
		},
	}
}

