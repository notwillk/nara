package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Inspect workspace",
	}

	listCmd.AddCommand(
		&cobra.Command{
			Use:   "schemas",
			Short: "List discovered schemas",
			Args:  cobra.NoArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				_ = configPath
				return fmt.Errorf("list schemas not implemented yet")
			},
		},
		&cobra.Command{
			Use:   "entries [entryGlob...]",
			Short: "List entry files",
			Args:  cobra.ArbitraryArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				_ = configPath
				_ = args
				return fmt.Errorf("list entries not implemented yet")
			},
		},
	)

	return listCmd
}

