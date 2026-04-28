package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vaultpatch/vaultpatch/internal/enrich"
	"github.com/vaultpatch/vaultpatch/internal/vault"
)

func init() {
	var (
		paths       []string
		annotations []string
		dryRun      bool
	)

	cmd := &cobra.Command{
		Use:   "enrich",
		Short: "Inject metadata annotations into Vault secrets",
		Example: `  vaultpatch enrich --path secret/app --annotation env=prod --annotation owner=team-a`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := vault.NewClient(vault.Options{})
			if err != nil {
				return fmt.Errorf("enrich: %w", err)
			}

			annoMap := make(map[string]string, len(annotations))
			for _, a := range annotations {
				parts := strings.SplitN(a, "=", 2)
				if len(parts) != 2 {
					return fmt.Errorf("enrich: invalid annotation %q (expected key=value)", a)
				}
				annoMap[parts[0]] = parts[1]
			}

			results, err := enrich.Enrich(client, enrich.Options{
				Paths:       paths,
				Annotations: annoMap,
				DryRun:      dryRun,
			})
			if err != nil {
				return err
			}

			for _, r := range results {
				status := "enriched"
				if r.DryRun {
					status = "dry-run"
				}
				fmt.Fprintf(os.Stdout, "[%s] %s (+%d keys)\n", status, r.Path, len(r.Added))
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVar(&paths, "path", nil, "Vault secret paths to enrich (repeatable)")
	cmd.Flags().StringArrayVar(&annotations, "annotation", nil, "Annotation in key=value format (repeatable)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")
	_ = cmd.MarkFlagRequired("path")
	_ = cmd.MarkFlagRequired("annotation")

	rootCmd.AddCommand(cmd)
}
