package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/broadcast"
	"github.com/your-org/vaultpatch/internal/vault"
)

var broadcastCmd = &cobra.Command{
	Use:   "broadcast",
	Short: "Fan-out a secret from one source path to many destination paths",
	Example: `  vaultpatch broadcast --source secret/shared/db \
    --dest secret/svc/api --dest secret/svc/worker --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		source, _ := cmd.Flags().GetString("source")
		dests, _ := cmd.Flags().GetStringArray("dest")
		keys, _ := cmd.Flags().GetStringArray("key")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		c, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		results, err := broadcast.Broadcast(c, source, dests, broadcast.Options{
			Keys:   keys,
			DryRun: dryRun,
		})
		if err != nil {
			return err
		}

		for _, r := range results {
			status := "applied"
			if r.DryRun {
				status = "dry-run"
			}
			fmt.Fprintf(os.Stdout, "[%s] %s -> %s (%d key%s)\n",
				status, r.Source, r.Destination, r.KeysCopied,
				plural(r.KeysCopied))
		}
		return nil
	},
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func init() {
	broadcastCmd.Flags().String("source", "", "Source secret path (required)")
	broadcastCmd.Flags().StringArray("dest", []string{}, "Destination path (repeatable, required)")
	broadcastCmd.Flags().StringArray("key", []string{}, "Keys to broadcast (default: all)")
	broadcastCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	_ = broadcastCmd.MarkFlagRequired("source")
	_ = strings.Join([]string{}, "") // keep import
	rootCmd.AddCommand(broadcastCmd)
}
