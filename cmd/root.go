package cmd

import (
	"fmt"
	"os"

	"github.com/notwillk/nara/internal/version"
	"github.com/spf13/cobra"
)

var configPath string

// Execute runs the CLI.
func Execute() {
	if err := rootCmd().Execute(); err != nil {
		// Cobra already prints errors for us in most cases, but we make sure
		// exit status is non-zero for CI/usage.
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "nara",
		Short:        "nara: validate, resolve, and compile entity graphs",
		Version:      version.Version,
		SilenceUsage: true,
	}

	cmd.PersistentFlags().StringVar(&configPath, "config", "nara.yaml", "path to config file")

	cmd.AddCommand(
		newCompileCmd(),
		newLintCmd(),
		newValidateCmd(),
		newFormatCmd(),
		newListCmd(),
		newInitCmd(),
	)

	return cmd
}

