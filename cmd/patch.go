package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/youorg/vaultpatch/internal/patch"
	"github.com/youorg/vaultpatch/internal/vault"
)

var patchCmd = &cobra.Command{
	Use:   "patch [paths...]",
	Short: "Apply targeted key-level mutations to Vault secret paths",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return err
		}

		setRaw, _ := cmd.Flags().GetStringArray("set")
		deleteKeys, _ := cmd.Flags().GetStringArray("delete")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		var ops []patch.Op
		for _, s := range setRaw {
			parts := strings.SplitN(s, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("invalid --set value %q: expected key=value", s)
			}
			ops = append(ops, patch.Op{Key: parts[0], Value: parts[1]})
		}
		for _, k := range deleteKeys {
			ops = append(ops, patch.Op{Key: k, Delete: true})
		}

		results, err := patch.Patch(client, args, ops, dryRun)
		if err != nil {
			return err
		}

		for _, r := range results {
			mode := "applied"
			if r.DryRun {
				mode = "dry-run"
			}
			fmt.Fprintf(os.Stdout, "[%s] %s — %d op(s)\n", mode, r.Path, len(r.Applied))
		}
		return nil
	},
}

func init() {
	patchCmd.Flags().StringArray("set", nil, "Key-value pair to set (key=value)")
	patchCmd.Flags().StringArray("delete", nil, "Key to delete")
	patchCmd.Flags().Bool("dry-run", false, "Preview changes without writing")
	rootCmd.AddCommand(patchCmd)
}
