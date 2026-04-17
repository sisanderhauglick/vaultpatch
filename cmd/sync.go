package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/sync"
	"github.com/your-org/vaultpatch/internal/vault"
)

var (
	syncSrc      string
	syncDst      string
	syncDryRun   bool
	syncIncludes string
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronise secrets from a source path to a destination path",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		var includes []string
		if syncIncludes != "" {
			for _, k := range strings.Split(syncIncludes, ",") {
				if k = strings.TrimSpace(k); k != "" {
					includes = append(includes, k)
				}
			}
		}

		result, err := sync.Sync(client, sync.Options{
			SrcPath:  syncSrc,
			DstPath:  syncDst,
			DryRun:   syncDryRun,
			Includes: includes,
		})
		if err != nil {
			return err
		}

		if len(result.Changes) == 0 {
			fmt.Fprintln(os.Stdout, "no changes detected")
			return nil
		}

		for _, c := range result.Changes {
			fmt.Fprintf(os.Stdout, "[%s] %s\n", c.Action, c.Key)
		}

		if syncDryRun {
			fmt.Fprintf(os.Stdout, "dry-run: %d change(s) would be applied\n", len(result.Changes))
		} else {
			fmt.Fprintf(os.Stdout, "applied %d change(s) to %s\n", result.Applied, syncDst)
		}
		return nil
	},
}

func init() {
	syncCmd.Flags().StringVar(&syncSrc, "src", "", "source secret path (required)")
	syncCmd.Flags().StringVar(&syncDst, "dst", "", "destination secret path (required)")
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "preview changes without writing")
	syncCmd.Flags().StringVar(&syncIncludes, "include", "", "comma-separated list of keys to sync")
	_ = syncCmd.MarkFlagRequired("src")
	_ = syncCmd.MarkFlagRequired("dst")
	rootCmd.AddCommand(syncCmd)
}
