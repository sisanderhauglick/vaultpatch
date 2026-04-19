package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/vaultpatch/internal/revert"
	"github.com/user/vaultpatch/internal/vault"
)

var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "Revert secret keys to previous values",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, _ := cmd.Flags().GetStringSlice("path")
		keys, _ := cmd.Flags().GetStringSlice("key")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// before values supplied as KEY=VALUE pairs via --before flag
		beforeRaw, _ := cmd.Flags().GetStringSlice("before")
		before := make(map[string]string, len(beforeRaw))
		for _, pair := range beforeRaw {
			for i, ch := range pair {
				if ch == '=' {
					before[pair[:i]] = pair[i+1:]
					break
				}
			}
		}

		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		results, err := revert.Revert(client, revert.Options{
			Paths:  paths,
			Before: before,
			Keys:   keys,
			DryRun: dryRun,
		})
		if err != nil {
			return err
		}

		for _, r := range results {
			mode := "applied"
			if r.DryRun {
				mode = "dry-run"
			}
			fmt.Fprintf(os.Stdout, "[%s] %s: reverted=%v skipped=%v\n",
				mode, r.Path, r.Reverted, r.Skipped)
		}
		return nil
	},
}

func init() {
	revertCmd.Flags().StringSlice("path", nil, "Vault paths to revert")
	revertCmd.Flags().StringSlice("key", nil, "Keys to revert (default: all in --before)")
	revertCmd.Flags().StringSlice("before", nil, "Previous values as KEY=VALUE pairs")
	revertCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	rootCmd.AddCommand(revertCmd)
}
