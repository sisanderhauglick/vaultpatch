package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpatch/vaultpatch/internal/redact"
	"github.com/vaultpatch/vaultpatch/internal/vault"
)

var redactCmd = &cobra.Command{
	Use:   "redact",
	Short: "Redact secret values at one or more paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, _ := cmd.Flags().GetString("addr")
		token, _ := cmd.Flags().GetString("token")
		paths, _ := cmd.Flags().GetStringSlice("path")
		keys, _ := cmd.Flags().GetStringSlice("key")
		replacement, _ := cmd.Flags().GetString("replacement")

		client, err := vault.NewClient(vault.Params{Addr: addr, Token: token})
		if err != nil {
			return err
		}

		all := make(map[string]map[string]string)
		for _, p := range paths {
			sm, err := vault.ReadSecrets(client, p)
			if err != nil {
				return fmt.Errorf("read %s: %w", p, err)
			}
			all[p] = sm
		}

		results := redact.Redact(all, redact.Options{
			Keys:        keys,
			Replacement: replacement,
		})

		return json.NewEncoder(os.Stdout).Encode(results)
	},
}

func init() {
	redactCmd.Flags().String("addr", "", "Vault address")
	redactCmd.Flags().String("token", "", "Vault token")
	redactCmd.Flags().StringSlice("path", nil, "Secret paths to redact")
	redactCmd.Flags().StringSlice("key", nil, "Keys to redact (default: all)")
	redactCmd.Flags().String("replacement", "[REDACTED]", "Replacement string")
	rootCmd.AddCommand(redactCmd)
}
