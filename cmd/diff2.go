package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/diff2"
	"github.com/your-org/vaultpatch/internal/vault"
)

var diff2Cmd = &cobra.Command{
	Use:   "diff2 <source> <dest>",
	Short: "Structured two-way diff between two Vault secret paths",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, _ := cmd.Flags().GetString("addr")
		token, _ := cmd.Flags().GetString("token")
		keys, _ := cmd.Flags().GetStringSlice("keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		client, err := vault.NewClient(vault.Params{Addr: addr, Token: token})
		if err != nil {
			return fmt.Errorf("diff2: %w", err)
		}

		result, err := diff2.Diff2(client, diff2.Options{
			Source: args[0],
			Dest:   args[1],
			Keys:   keys,
			DryRun: dryRun,
		})
		if err != nil {
			return err
		}

		if !result.HasChanges() {
			fmt.Fprintln(os.Stdout, "No differences found.")
			return nil
		}

		for _, c := range result.Changes {
			switch c.Kind {
			case diff2.Added:
				fmt.Fprintf(os.Stdout, "+ %s = %s\n", c.Key, c.NewValue)
			case diff2.Removed:
				fmt.Fprintf(os.Stdout, "- %s = %s\n", c.Key, c.OldValue)
			case diff2.Modified:
				fmt.Fprintf(os.Stdout, "~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
			}
		}

		if dryRun {
			fmt.Fprintln(os.Stdout, strings.Repeat("-", 40))
			fmt.Fprintln(os.Stdout, "[dry-run] no changes applied")
		}
		return nil
	},
}

func init() {
	diff2Cmd.Flags().String("addr", "", "Vault address (overrides VAULT_ADDR)")
	diff2Cmd.Flags().String("token", "", "Vault token (overrides VAULT_TOKEN)")
	diff2Cmd.Flags().StringSlice("keys", nil, "Comma-separated list of keys to diff (default: all)")
	diff2Cmd.Flags().Bool("dry-run", false, "Show diff without applying changes")
	rootCmd.AddCommand(diff2Cmd)
}
