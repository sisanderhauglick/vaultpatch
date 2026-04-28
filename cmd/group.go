package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpatch/internal/group"
	"github.com/your-org/vaultpatch/internal/vault"
)

var (
	groupSources     []string
	groupDestination string
	groupPrefix      bool
	groupDryRun      bool
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Merge secrets from multiple paths into a single destination path",
	Long: `Reads secrets from each source path and writes them combined into the
destination path. Key collisions are resolved by last-writer-wins (source order).

Use --prefix to namespace each key with its source path's final segment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		res, err := group.Group(client, group.Options{
			Sources:     groupSources,
			Destination: groupDestination,
			Prefix:      groupPrefix,
			DryRun:      groupDryRun,
		})
		if err != nil {
			return err
		}

		if groupDryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would merge %d key(s) into %s\n",
				res.KeysMerged, res.Destination)
		} else {
			fmt.Fprintf(os.Stdout, "grouped %d key(s) into %s at %s\n",
				res.KeysMerged, res.Destination, res.GroupedAt.Format("2006-01-02T15:04:05Z"))
		}
		return nil
	},
}

func init() {
	groupCmd.Flags().StringSliceVarP(&groupSources, "source", "s", nil,
		"source paths to read from (repeatable, required)")
	groupCmd.Flags().StringVarP(&groupDestination, "dest", "d", "",
		"destination path to write merged secrets into (required)")
	groupCmd.Flags().BoolVar(&groupPrefix, "prefix", false,
		"prefix each key with its source path segment")
	groupCmd.Flags().BoolVar(&groupDryRun, "dry-run", false,
		"preview changes without writing to Vault")
	_ = groupCmd.MarkFlagRequired("source")
	_ = groupCmd.MarkFlagRequired("dest")
	rootCmd.AddCommand(groupCmd)
}
