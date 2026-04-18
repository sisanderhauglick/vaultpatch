package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/tag"
	"github.com/your-org/vaultpatch/internal/vault"
)

var (
	tagAdd    []string
	tagRemove []string
	tagDryRun bool
)

var tagCmd = &cobra.Command{
	Use:   "tag <path>",
	Short: "Add or remove metadata tags on a Vault secret path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		result, err := tag.Tag(client, args[0], tagAdd, tagRemove, tagDryRun)
		if err != nil {
			return err
		}

		if tagDryRun {
			fmt.Printf("[dry-run] %s → tags: [%s]\n", result.Path, strings.Join(result.Tags, ", "))
		} else {
			fmt.Printf("tagged %s → [%s]\n", result.Path, strings.Join(result.Tags, ", "))
		}
		return nil
	},
}

func init() {
	tagCmd.Flags().StringSliceVar(&tagAdd, "add", nil, "Tags to add (comma-separated)")
	tagCmd.Flags().StringSliceVar(&tagRemove, "remove", nil, "Tags to remove (comma-separated)")
	tagCmd.Flags().BoolVar(&tagDryRun, "dry-run", false, "Preview changes without writing")
	rootCmd.AddCommand(tagCmd)
}
