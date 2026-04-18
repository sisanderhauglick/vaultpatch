package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultpatch/internal/compare"
	"github.com/vaultpatch/internal/vault"
)

var compareCmd = &cobra.Command{
	Use:   "compare <src-path> <dst-path>",
	Short: "Compare secrets between two Vault paths",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		res, err := compare.Compare(client, args[0], args[1])
		if err != nil {
			return err
		}

		w := os.Stdout
		fmt.Fprintf(w, "Comparing %s → %s\n\n", res.SourcePath, res.DestPath)

		if len(res.OnlyInSrc) > 0 {
			fmt.Fprintln(w, "Only in source:")
			for k := range res.OnlyInSrc {
				fmt.Fprintf(w, "  + %s\n", k)
			}
		}
		if len(res.OnlyInDst) > 0 {
			fmt.Fprintln(w, "Only in destination:")
			for k := range res.OnlyInDst {
				fmt.Fprintf(w, "  - %s\n", k)
			}
		}
		if len(res.Differ) > 0 {
			fmt.Fprintln(w, "Different values:")
			for k, pair := range res.Differ {
				fmt.Fprintf(w, "  ~ %s  [%s] → [%s]\n", k, pair[0], pair[1])
			}
		}
		if len(res.Match) > 0 {
			fmt.Fprintf(w, "\n%d key(s) match.\n", len(res.Match))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}
