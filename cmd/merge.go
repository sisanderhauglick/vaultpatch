package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/merge"
	"github.com/your-org/vaultpatch/internal/vault"
)

var (
	mergeSources     []string
	mergeDestination string
	mergeKeys        []string
	mergeDryRun      bool
)

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge secrets from multiple source paths into a destination path",
	Example: `  vaultpatch merge --src secret/app/staging --src secret/app/base \
    --dest secret/app/merged --dry-run`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("merge: %w", err)
		}

		result, err := merge.Merge(client, merge.Options{
			Sources:     mergeSources,
			Destination: mergeDestination,
			Keys:        mergeKeys,
			DryRun:      mergeDryRun,
		})
		if err != nil {
			return err
		}

		if result.DryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would merge %d key(s) into %q\n",
				result.Merged, result.Destination)
		} else {
			fmt.Fprintf(os.Stdout, "merged %d key(s) into %q\n",
				result.Merged, result.Destination)
		}
		return nil
	},
}

func init() {
	mergeCmd.Flags().StringArrayVar(&mergeSources, "src", nil, "source path (repeatable, later sources take precedence)")
	mergeCmd.Flags().StringVar(&mergeDestination, "dest", "", "destination path")
	mergeCmd.Flags().StringArrayVar(&mergeKeys, "key", nil, "restrict merge to specific keys (repeatable)")
	mergeCmd.Flags().BoolVar(&mergeDryRun, "dry-run", false, "preview changes without writing")
	_ = mergeCmd.MarkFlagRequired("src")
	_ = mergeCmd.MarkFlagRequired("dest")
	rootCmd.AddCommand(mergeCmd)
}
