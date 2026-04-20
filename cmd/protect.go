package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/protect"
	"github.com/your-org/vaultpatch/internal/vault"
)

var (
	protectOwner  string
	protectDryRun bool
)

var protectCmd = &cobra.Command{
	Use:   "protect [flags] <path>...",
	Short: "Mark secret paths as write-protected",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Config{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}
		results, err := protect.Protect(client, protect.Options{
			Paths:  args,
			Owner:  protectOwner,
			DryRun: protectDryRun,
		})
		if err != nil {
			return err
		}
		for _, r := range results {
			status := "protected"
			if protectDryRun {
				status = "would protect"
			}
			fmt.Fprintf(os.Stdout, "%s %s (owner: %s)\n", status, r.Path, r.Owner)
		}
		return nil
	},
}

var unprotectCmd = &cobra.Command{
	Use:   "unprotect [flags] <path>...",
	Short: "Remove write-protection from secret paths",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Config{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}
		results, err := protect.Unprotect(client, protect.Options{
			Paths:  args,
			DryRun: protectDryRun,
		})
		if err != nil {
			return err
		}
		for _, r := range results {
			status := "unprotected"
			if protectDryRun {
				status = "would unprotect"
			}
			fmt.Fprintf(os.Stdout, "%s %s\n", status, r.Path)
		}
		return nil
	},
}

func init() {
	protectCmd.Flags().StringVar(&protectOwner, "owner", "", "owner identifier for the protection record (required)")
	protectCmd.Flags().BoolVar(&protectDryRun, "dry-run", false, "preview changes without writing to Vault")
	_ = protectCmd.MarkFlagRequired("owner")

	unprotectCmd.Flags().BoolVar(&protectDryRun, "dry-run", false, "preview changes without writing to Vault")

	rootCmd.AddCommand(protectCmd)
	rootCmd.AddCommand(unprotectCmd)
}
