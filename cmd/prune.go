package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/prune"
	"github.com/your-org/vaultpatch/internal/vault"
)

var (
	prunePaths     []string
	pruneOlderThan string
	pruneDryRun    bool
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Remove secrets older than a specified duration",
	Long: `Prune reads each Vault path and deletes secrets whose _created_at
timestamp is older than the given duration. Use --dry-run to preview
what would be removed without making any changes.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		duration, err := time.ParseDuration(pruneOlderThan)
		if err != nil {
			return fmt.Errorf("invalid --older-than value %q: %w", pruneOlderThan, err)
		}

		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		results, err := prune.Prune(client, prune.Options{
			Paths:     prunePaths,
			OlderThan: duration,
			DryRun:    pruneDryRun,
		})
		if err != nil {
			return err
		}

		for _, r := range results {
			switch {
			case r.Pruned && r.DryRun:
				fmt.Fprintf(os.Stdout, "[dry-run] would prune %s (%s)\n", r.Path, r.Reason)
			case r.Pruned:
				fmt.Fprintf(os.Stdout, "pruned %s at %s\n", r.Path, r.PrunedAt.Format(time.RFC3339))
			default:
				fmt.Fprintf(os.Stdout, "skipped %s: %s\n", r.Path, r.Reason)
			}
		}
		return nil
	},
}

func init() {
	pruneCmd.Flags().StringSliceVar(&prunePaths, "path", nil, "Vault paths to evaluate (required)")
	pruneCmd.Flags().StringVar(&pruneOlderThan, "older-than", "720h", "Remove secrets older than this duration (e.g. 720h, 30d)")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false, "Preview changes without applying them")
	_ = pruneCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(pruneCmd)
}
