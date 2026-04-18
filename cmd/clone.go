package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/clone"
	"github.com/your-org/vaultpatch/internal/vault"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <source> <destination>",
	Short: "Clone secrets from one Vault path to another",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return err
		}

		keysRaw, _ := cmd.Flags().GetString("keys")
		var keys []string
		if keysRaw != "" {
			keys = strings.Split(keysRaw, ",")
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		res, err := clone.Clone(client, clone.Options{
			Source:      args[0],
			Destination: args[1],
			Keys:        keys,
			DryRun:      dryRun,
		})
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would clone %d key(s) from %s to %s\n",
				res.KeysCopied, res.Source, res.Destination)
		} else {
			fmt.Fprintf(os.Stdout, "cloned %d key(s) from %s to %s\n",
				res.KeysCopied, res.Source, res.Destination)
		}
		return nil
	},
}

func init() {
	cloneCmd.Flags().String("keys", "", "comma-separated list of keys to clone (default: all)")
	cloneCmd.Flags().Bool("dry-run", false, "preview changes without writing to Vault")
	rootCmd.AddCommand(cloneCmd)
}
