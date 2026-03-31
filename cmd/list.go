package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/notwillk/nara/internal/config"
	"github.com/notwillk/nara/internal/schema"
	"github.com/notwillk/nara/internal/workspace"
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
				cfg, err := config.Load(configPath)
				if err != nil {
					return err
				}
				baseDir := filepath.Dir(configPath)
				discovered, err := schema.Discover(cfg, baseDir)
				if err != nil {
					return err
				}
				for _, item := range discovered {
					rel := relativePath(baseDir, item.Path)
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", item.Name, rel); err != nil {
						return err
					}
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "entries [entryGlob...]",
			Short: "List entry files",
			Args:  cobra.ArbitraryArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := config.Load(configPath)
				if err != nil {
					return err
				}
				baseDir := filepath.Dir(configPath)
				discovered, err := schema.Discover(cfg, baseDir)
				if err != nil {
					return err
				}
				allowed := make(map[string]bool, len(discovered))
				for _, item := range discovered {
					allowed[item.Name] = true
				}

				var entries []workspace.EntryFile
				if len(args) == 0 {
					entries, err = workspace.DiscoverEntries(cfg, baseDir, allowed)
				} else {
					entries, err = workspace.EntriesForGlobs(args, cfg, allowed)
				}
				if err != nil {
					return err
				}
				for _, entry := range entries {
					rel := relativePath(baseDir, entry.Path)
					if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", entry.ID, entry.Schema, rel); err != nil {
						return err
					}
				}
				return nil
			},
		},
	)

	return listCmd
}

func relativePath(baseDir string, path string) string {
	rel, err := filepath.Rel(baseDir, path)
	if err != nil {
		return path
	}
	return rel
}
