package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpatch/vaultpatch/internal/rollback"
	"github.com/vaultpatch/vaultpatch/internal/vault"
)

var (
	rollbackPath   string
	rollbackDryRun bool
)

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Capture or restore a snapshot of Vault secrets at a given path",
	Long: `rollback captures the current state of secrets at a path (snapshot)
or restores a previously captured snapshot back to Vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return fmt.Errorf("rollback: %w", err)
		}

		ctx := context.Background()

		snap, err := rollback.Capture(ctx, client, rollbackPath)
		if err != nil {
			return err
		}

		fmt.Printf("Captured %d keys from %s\n", len(snap.Secrets), snap.Path)

		if err := rollback.Restore(ctx, client, snap, rollbackDryRun); err != nil {
			return err
		}

		if rollbackDryRun {
			fmt.Println("Dry-run complete. No changes written.")
		} else {
			fmt.Println("Restore complete.")
		}
		return nil
	},
}

func init() {
	rollbackCmd.Flags().StringVarP(&rollbackPath, "path", "p", "", "Vault secret path to snapshot/restore (required)")
	rollbackCmd.Flags().BoolVar(&rollbackDryRun, "dry-run", false, "Preview restore without writing to Vault")
	_ = rollbackCmd.MarkFlagRequired("path")

	if root := rootCmd(); root != nil {
		root.AddCommand(rollbackCmd)
	} else {
		fmt.Fprintln(os.Stderr, "warn: rootCmd not available, rollback command not registered")
	}
}
