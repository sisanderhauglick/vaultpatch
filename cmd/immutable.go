package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpatch/internal/immutable"
	"github.com/yourusername/vaultpatch/internal/vault"
)

var (
	immutableDryRun  bool
	immutableRelease bool
)

var immutableCmd = &cobra.Command{
	Use:   "immutable [paths...]",
	Short: "Mark or release immutable locks on Vault secret paths",
	Long: `Mark one or more Vault secret paths as immutable by writing a sentinel key.
Use --release to remove the immutable lock from the specified paths.
Use --dry-run to preview changes without writing to Vault.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("immutable: %w", err)
		}

		opts := immutable.Options{
			Paths:  args,
			DryRun: immutableDryRun,
		}

		var results []immutable.Result
		if immutableRelease {
			results, err = immutable.Release(client, opts)
		} else {
			results, err = immutable.Immute(client, opts)
		}
		if err != nil {
			return err
		}

		for _, r := range results {
			switch {
			case r.DryRun && r.Released:
				fmt.Fprintf(os.Stdout, "[dry-run] would release immutable lock: %s\n", r.Path)
			case r.DryRun:
				fmt.Fprintf(os.Stdout, "[dry-run] would mark immutable: %s\n", r.Path)
			case r.Released:
				fmt.Fprintf(os.Stdout, "released: %s\n", r.Path)
			default:
				fmt.Fprintf(os.Stdout, "immutable: %s (locked at %s)\n", r.Path, r.LockedAt.Format("2006-01-02T15:04:05Z"))
			}
		}
		return nil
	},
}

func init() {
	immutableCmd.Flags().BoolVar(&immutableDryRun, "dry-run", false, "Preview changes without writing to Vault")
	immutableCmd.Flags().BoolVar(&immutableRelease, "release", false, "Remove the immutable lock from the specified paths")
	rootCmd.AddCommand(immutableCmd)
}
