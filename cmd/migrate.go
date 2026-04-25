package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpatch/internal/migrate"
	"github.com/yourusername/vaultpatch/internal/vault"
)

var (
	migrateSources     []string
	migrateDestination string
	migrateKeyMap      []string // "old=new" pairs
	migrateDryRun      bool
	migrateOverwrite   bool
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Move secrets from one or more paths to a destination, with optional key remapping",
	Example: `  vaultpatch migrate --src secret/app/v1 --dest secret/app/v2 --map DB_PASS=DATABASE_PASSWORD --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := vault.NewClient(vault.Params{})
		if err != nil {
			return fmt.Errorf("migrate: %w", err)
		}

		keyMap := parseKeyMap(migrateKeyMap)

		results, err := migrate.Migrate(c, migrate.Options{
			Sources:     migrateSources,
			Destination: migrateDestination,
			KeyMap:      keyMap,
			DryRun:      migrateDryRun,
			Overwrite:   migrateOverwrite,
		})
		if err != nil {
			return err
		}

		for _, r := range results {
			mode := "applied"
			if r.DryRun {
				mode = "dry-run"
			}
			fmt.Fprintf(os.Stdout, "[%s] %s -> %s (%d keys)\n", mode, r.Source, r.Destination, r.KeysMapped)
		}
		return nil
	},
}

// parseKeyMap converts ["old=new", ...] slice into a map.
func parseKeyMap(pairs []string) map[string]string {
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		for i := 0; i < len(p); i++ {
			if p[i] == '=' {
				m[p[:i]] = p[i+1:]
				break
			}
		}
	}
	return m
}

func init() {
	migrateCmd.Flags().StringArrayVar(&migrateSources, "src", nil, "Source path(s) to migrate from (repeatable)")
	migrateCmd.Flags().StringVar(&migrateDestination, "dest", "", "Destination path to migrate into")
	migrateCmd.Flags().StringArrayVar(&migrateKeyMap, "map", nil, "Key rename pairs in old=new format (repeatable)")
	migrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "Preview changes without writing")
	migrateCmd.Flags().BoolVar(&migrateOverwrite, "overwrite", false, "Overwrite existing keys at destination")
	_ = migrateCmd.MarkFlagRequired("src")
	_ = migrateCmd.MarkFlagRequired("dest")
	rootCmd.AddCommand(migrateCmd)
}
