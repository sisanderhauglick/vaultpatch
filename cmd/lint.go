package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpatch/internal/lint"
	"github.com/yourorg/vaultpatch/internal/vault"
)

var lintPaths []string

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Check Vault secrets against lint rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("lint: %w", err)
		}

		results, err := lint.Lint(client, lint.Options{
			Paths: lintPaths,
			Rules: lint.DefaultRules(),
		})
		if err != nil {
			return err
		}

		total := 0
		for _, r := range results {
			for _, v := range r.Violations {
				fmt.Fprintf(os.Stdout, "[%s] %s → %s: %s\n", v.Rule, v.Path, v.Key, v.Message)
				total++
			}
		}
		if total == 0 {
			fmt.Println("no lint violations found")
		} else {
			fmt.Fprintf(os.Stdout, "%d violation(s) found\n", total)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	lintCmd.Flags().StringSliceVar(&lintPaths, "path", nil, "Vault paths to lint (repeatable)")
	_ = lintCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(lintCmd)
}
