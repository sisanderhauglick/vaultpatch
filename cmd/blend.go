package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultpatch/internal/blend"
	"github.com/your-org/vaultpatch/internal/vault"
)

var blendCmd = &cobra.Command{
	Use:   "blend",
	Short: "Blend secrets from multiple source paths into a single destination",
	Long: `Read secrets from one or more source paths and write the merged result
to a destination path. Conflict resolution is controlled by --strategy.`,
	RunE: runBlend,
}

func init() {
	blendCmd.Flags().StringSliceP("sources", "s", nil, "source paths (comma-separated or repeated flag)")
	blendCmd.Flags().StringP("destination", "d", "", "destination path")
	blendCmd.Flags().String("strategy", "last", `conflict resolution strategy: "first" or "last"`)
	blendCmd.Flags().StringSlice("keys", nil, "keys to include (default: all)")
	blendCmd.Flags().Bool("dry-run", false, "preview changes without writing")
	_ = blendCmd.MarkFlagRequired("sources")
	_ = blendCmd.MarkFlagRequired("destination")
	rootCmd.AddCommand(blendCmd)
}

func runBlend(cmd *cobra.Command, _ []string) error {
	sources, _ := cmd.Flags().GetStringSlice("sources")
	destination, _ := cmd.Flags().GetString("destination")
	strategyStr, _ := cmd.Flags().GetString("strategy")
	keys, _ := cmd.Flags().GetStringSlice("keys")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	strategy := blend.Strategy(strings.ToLower(strategyStr))
	if strategy != blend.StrategyFirst && strategy != blend.StrategyLast {
		return fmt.Errorf("blend: unknown strategy %q; use \"first\" or \"last\"", strategyStr)
	}

	client, err := vault.NewClient(vault.Options{})
	if err != nil {
		return fmt.Errorf("blend: %w", err)
	}

	res, err := blend.Blend(client, blend.Options{
		Sources:     sources,
		Destination: destination,
		Strategy:    strategy,
		Keys:        keys,
		DryRun:      dryRun,
	})
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Fprintf(os.Stdout, "[dry-run] would blend %d key(s) into %s\n", len(res.BlendedKeys), res.Destination)
	} else {
		fmt.Fprintf(os.Stdout, "blended %d key(s) into %s at %s\n",
			len(res.BlendedKeys), res.Destination, res.BlendedAt.Format("2006-01-02T15:04:05Z"))
	}
	return nil
}
