package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command for the vaultpatch CLI.
var rootCmd = &cobra.Command{
	Use:   "vaultpatch",
	Short: "Diff and apply HashiCorp Vault secret changes across environments",
	Long: `vaultpatch is a CLI tool for managing HashiCorp Vault secrets across
environments. It supports diffing, applying, promoting, syncing, copying,
merging, renaming, rolling back, validating, and exporting secrets.`,
	SilenceUsage: true,
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().String("token", "", "Vault token (overrides VAULT_TOKEN)")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Preview changes without applying them")
}
