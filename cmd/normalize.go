package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/normalize"
	"github.com/your-org/vaultpatch/internal/vault"
)

var normalizeCmd = &cobra.Command{
	Use:   "normalize [paths...]",
	Short: "Normalise secret values and key names in Vault",
	Long: `Reads secrets at each path and applies normalisation rules:
  - trim surrounding whitespace from values
  - lowercase key names
  - uppercase values

Use --dry-run to preview changes without writing to Vault.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		trimSpace, _ := cmd.Flags().GetBool("trim-space")
		lowercaseKeys, _ := cmd.Flags().GetBool("lowercase-keys")
		uppercaseValues, _ := cmd.Flags().GetBool("uppercase-values")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		keysFlag, _ := cmd.Flags().GetString("keys")

		var keys []string
		if keysFlag != "" {
			for _, k := range strings.Split(keysFlag, ",") {
				if k = strings.TrimSpace(k); k != "" {
					keys = append(keys, k)
				}
			}
		}

		client, err := vault.NewClient(vault.Params{
			Address: os.Getenv("VAULT_ADDR"),
			Token:   os.Getenv("VAULT_TOKEN"),
		})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		results, err := normalize.Normalize(client, normalize.Options{
			Paths:           args,
			Keys:            keys,
			TrimSpace:       trimSpace,
			LowercaseKeys:   lowercaseKeys,
			UppercaseValues: uppercaseValues,
			DryRun:          dryRun,
		})
		if err != nil {
			return err
		}

		if len(results) == 0 {
			fmt.Println("no changes detected")
			return nil
		}

		for _, r := range results {
			prefix := ""
			if r.DryRun {
				prefix = "[dry-run] "
			}
			fmt.Printf("%s%s: %d key(s) normalised\n", prefix, r.Path, len(r.Changes))
		}
		return nil
	},
}

func init() {
	normalizeCmd.Flags().Bool("trim-space", true, "trim surrounding whitespace from values")
	normalizeCmd.Flags().Bool("lowercase-keys", false, "lowercase all key names")
	normalizeCmd.Flags().Bool("uppercase-values", false, "uppercase all values")
	normalizeCmd.Flags().Bool("dry-run", false, "preview changes without writing to Vault")
	normalizeCmd.Flags().String("keys", "", "comma-separated list of keys to process (default: all)")
	rootCmd.AddCommand(normalizeCmd)
}
