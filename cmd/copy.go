package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/copy"
	"github.com/your-org/vaultpatch/internal/vault"
)

var (
	copySource string
	copyDest   string
	copyKeys   string
	copyDryRun bool
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy secrets from one Vault path to another",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		var keys []string
		if copyKeys != "" {
			for _, k := range strings.Split(copyKeys, ",") {
				if k = strings.TrimSpace(k); k != "" {
					keys = append(keys, k)
				}
			}
		}

		res, err := copy.Copy(client, copy.Options{
			SourcePath: copySource,
			DestPath:   copyDest,
			Keys:       keys,
			DryRun:     copyDryRun,
		})
		if err != nil {
			return err
		}

		if res.DryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would copy %d key(s) from %s → %s\n", res.Keys, res.SourcePath, res.DestPath)
		} else {
			fmt.Fprintf(os.Stdout, "copied %d key(s) from %s → %s\n", res.Keys, res.SourcePath, res.DestPath)
		}
		return nil
	},
}

func init() {
	copyCmd.Flags().StringVar(&copySource, "source", "", "source Vault path (required)")
	copyCmd.Flags().StringVar(&copyDest, "dest", "", "destination Vault path (required)")
	copyCmd.Flags().StringVar(&copyKeys, "keys", "", "comma-separated list of keys to copy (default: all)")
	copyCmd.Flags().BoolVar(&copyDryRun, "dry-run", false, "preview changes without writing")
	_ = copyCmd.MarkFlagRequired("source")
	_ = copyCmd.MarkFlagRequired("dest")
	rootCmd.AddCommand(copyCmd)
}
