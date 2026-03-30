package cmd

import (
	"fmt"
	"strings"

	"github.com/notwillk/nara/internal/compiler"
	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/emitter"
	"github.com/spf13/cobra"
)

func newCompileCmd() *cobra.Command {
	var outPath string
	var format string

	cmd := &cobra.Command{
		Use:   "compile [entryGlob...]",
		Short: "Compile entrypoint files into a resolved data graph",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if outPath == "" {
				return fmt.Errorf("--out is required")
			}
			result, err := compiler.Compile(configPath, args)
			if err != nil {
				return err
			}

			switch strings.ToLower(format) {
			case "", "yaml", "yml":
				return emitter.WriteYAML(outPath, result.Entries, mustLoadConfigForEmit())
			case "json":
				return emitter.WriteJSON(outPath, result.Entries, mustLoadConfigForEmit())
			case "sqlite":
				return fmt.Errorf("sqlite emitter not implemented yet")
			default:
				return fmt.Errorf("unsupported --format: %s", format)
			}
		},
	}

	cmd.Flags().StringVar(&outPath, "out", "", "output path")
	cmd.Flags().StringVar(&format, "format", "yaml", "output format: yaml|json|sqlite")
	return cmd
}

func mustLoadConfigForEmit() *config.Config {
	cfg, err := config.Load(configPath)
	if err != nil {
		// compile already loads/validates config; this should be unreachable.
		return &config.Config{}
	}
	return cfg
}

