package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpatch/vaultpatch/internal/export"
	"github.com/vaultpatch/vaultpatch/internal/vault"
)

var (
	exportFormat string
	exportPath   string
)

var exportCmd = &cobra.Command{
	Use:   "export <secret-path>",
	Short: "Export secrets from Vault in a specified format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt, err := export.ParseFormat(exportFormat)
		if err != nil {
			return err
		}

		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		secrets, err := vault.ReadSecrets(client, args[0])
		if err != nil {
			return fmt.Errorf("read secrets: %w", err)
		}

		w := cmd.OutOrStdout()
		if exportPath != "" {
			f, err := os.Create(exportPath)
			if err != nil {
				return fmt.Errorf("open output file: %w", err)
			}
			defer f.Close()
			w = f
		}

		return export.Export(secrets, fmt, w)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "output format: json, yaml, env")
	exportCmd.Flags().StringVarP(&exportPath, "output", "o", "", "write output to file instead of stdout")
	rootCmd.AddCommand(exportCmd)
}
