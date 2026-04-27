package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/split"
	"github.com/your-org/vaultpatch/internal/vault"
)

var splitCmd = &cobra.Command{
	Use:   "split <source> <dest=key1,key2>...",
	Short: "Split a secret path into multiple destinations",
	Long: `Read secrets from <source> and write key subsets to each destination.

Each positional argument after <source> must be of the form:
  destination=key1,key2,...

Example:
  vaultpatch split secret/app secret/app/db=DB_HOST,DB_PASS secret/app/api=API_KEY`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		client, err := vault.NewClient(vault.Options{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}

		source := args[0]
		assignments := make(map[string][]string)
		for _, arg := range args[1:] {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return fmt.Errorf("invalid assignment %q: expected dest=key1,key2", arg)
			}
			keys := strings.Split(parts[1], ",")
			assignments[parts[0]] = keys
		}

		results, err := split.Split(client, source, split.Options{
			Assignments: assignments,
			DryRun:      dryRun,
		})
		if err != nil {
			return err
		}

		for _, r := range results {
			status := "applied"
			if r.DryRun {
				status = "dry-run"
			}
			fmt.Fprintf(os.Stdout, "[%s] %s -> %s (%d keys)\n",
				status, r.Source, r.Destination, len(r.Keys))
		}
		return nil
	},
}

func init() {
	splitCmd.Flags().Bool("dry-run", false, "Preview changes without writing to Vault")
	rootCmd.AddCommand(splitCmd)
}
