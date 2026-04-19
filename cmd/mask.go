package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/mask"
	"github.com/your-org/vaultpatch/internal/vault"
)

var maskCmd = &cobra.Command{
	Use:   "mask [path]",
	Short: "Display secrets with sensitive values redacted",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return err
		}

		secrets, err := vault.ReadSecrets(client, args[0])
		if err != nil {
			return err
		}

		keysFlag, _ := cmd.Flags().GetString("keys")
		placeholder, _ := cmd.Flags().GetString("placeholder")

		var keys []string
		if keysFlag != "" {
			for _, k := range strings.Split(keysFlag, ",") {
				if t := strings.TrimSpace(k); t != "" {
					keys = append(keys, t)
				}
			}
		}

		r := mask.Mask(args[0], secrets, mask.Options{
			Keys:        keys,
			Placeholder: placeholder,
		})

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(r.Masked); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "redacted %d key(s) at %s\n", r.Redacted, r.Path)
		return nil
	},
}

func init() {
	maskCmd.Flags().String("keys", "", "Comma-separated keys to redact (default: all)")
	maskCmd.Flags().String("placeholder", "***", "Replacement value for redacted keys")
	rootCmd.AddCommand(maskCmd)
}
