package cmd

import (
	internallint "github.com/notwillk/nara/internal/lint"
	"github.com/spf13/cobra"
)

func newLintCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lint",
		Short: "Validate workspace correctness (no mutations)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return internallint.Run(configPath)
		},
	}
}
