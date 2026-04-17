package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultpatch/internal/promote"
	"github.com/vaultpatch/internal/vault"
)

func init() {
	var (
		src    string
		dst    string
		dryRun bool
		keys   string
	)

	cmd := &cobra.Command{
		Use:   "promote",
		Short: "Promote secrets from one Vault path to another",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := vault.NewClient(vault.Options{})
			if err != nil {
				return fmt.Errorf("vault client: %w", err)
			}

			var keyList []string
			if keys != "" {
				for _, k := range strings.Split(keys, ",") {
					if k = strings.TrimSpace(k); k != "" {
						keyList = append(keyList, k)
					}
				}
			}

			result, err := promote.Promote(client, promote.Options{
				SrcPath: src,
				DstPath: dst,
				DryRun:  dryRun,
				Keys:    keyList,
			})
			if err != nil {
				return err
			}

			mode := "applied"
			if dryRun {
				mode = "dry-run"
			}
			fmt.Fprintf(os.Stdout, "promote [%s]: %d changes, %d applied, %d skipped\n",
				mode, len(result.Changes), result.Applied, result.Skipped)
			for _, c := range result.Changes {
				fmt.Fprintf(os.Stdout, "  [%s] %s\n", c.Action, c.Key)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&src, "src", "", "source Vault path (required)")
	cmd.Flags().StringVar(&dst, "dst", "", "destination Vault path (required)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview changes without writing")
	cmd.Flags().StringVar(&keys, "keys", "", "comma-separated list of keys to promote")
	_ = cmd.MarkFlagRequired("src")
	_ = cmd.MarkFlagRequired("dst")

	rootCmd.AddCommand(cmd)
}
