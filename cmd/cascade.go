package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpatch/internal/cascade"
	"github.com/your-org/vaultpatch/internal/vault"
)

var cascadeCmd = &cobra.Command{
	Use:   "cascade",
	Short: "Propagate secrets from a source path to one or more destination paths",
	Example: `  vaultpatch cascade --source secret/base --dest secret/staging --dest secret/prod
  vaultpatch cascade --source secret/base --dest secret/staging --key DB_PASS --key API_KEY --dry-run`,
	RunE: runCascade,
}

var (
	cascadeSource string
	cascadeDests  []string
	cascadeKeys   []string
	cascadeDryRun bool
)

func init() {
	cascadeCmd.Flags().StringVar(&cascadeSource, "source", "", "Source secret path (required)")
	cascadeCmd.Flags().StringArrayVar(&cascadeDests, "dest", nil, "Destination secret path (repeatable, required)")
	cascadeCmd.Flags().StringArrayVar(&cascadeKeys, "key", nil, "Keys to cascade (default: all)")
	cascadeCmd.Flags().BoolVar(&cascadeDryRun, "dry-run", false, "Preview changes without writing")
	_ = cascadeCmd.MarkFlagRequired("source")
	_ = cascadeCmd.MarkFlagRequired("dest")
	rootCmd.AddCommand(cascadeCmd)
}

func runCascade(cmd *cobra.Command, _ []string) error {
	client, err := vault.NewClient(vault.Options{})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	results, err := cascade.Cascade(client, cascade.Options{
		Source:       cascadeSource,
		Destinations: cascadeDests,
		Keys:         cascadeKeys,
		DryRun:       cascadeDryRun,
	})
	if err != nil {
		return err
	}

	for _, r := range results {
		prefix := "cascaded"
		if r.DryRun {
			prefix = "[dry-run]"
		}
		fmt.Fprintf(os.Stdout, "%s %s -> %s (%d keys)\n",
			prefix, r.Source, r.Destination, r.KeysCopied)
	}
	return nil
}
