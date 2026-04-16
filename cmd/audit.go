package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultpatch/vaultpatch/internal/audit"
	"github.com/vaultpatch/vaultpatch/internal/diff"
	"github.com/vaultpatch/vaultpatch/internal/vault"
)

var auditPath string

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Log current secrets at a path as an audit snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, auditPath)
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		changes := make([]diff.Change, 0, len(secrets))
		for k, v := range secrets {
			changes = append(changes, diff.Change{
				Key:      k,
				Action:   diff.ActionAdded,
				NewValue: fmt.Sprintf("%v", v),
			})
		}

		logger := audit.NewLogger(os.Stdout)
		if err := logger.Log(auditPath, changes, false, false); err != nil {
			return fmt.Errorf("audit log: %w", err)
		}
		return nil
	},
}

func init() {
	auditCmd.Flags().StringVarP(&auditPath, "path", "p", "", "Vault secret path to audit (required)")
	_ = auditCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(auditCmd)
}
