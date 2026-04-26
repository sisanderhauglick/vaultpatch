package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpatch/internal/audit2"
)

var (
	audit2Operation string
	audit2Path      string
	audit2DryRun    bool
	audit2Output    string
)

var audit2Cmd = &cobra.Command{
	Use:   "audit2",
	Short: "Write a structured audit log entry for a Vault operation",
	RunE: func(cmd *cobra.Command, args []string) error {
		w := os.Stdout
		if audit2Output != "" && audit2Output != "-" {
			f, err := os.OpenFile(audit2Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
			if err != nil {
				return fmt.Errorf("open audit log: %w", err)
			}
			defer f.Close()
			w = f
		}

		logger, err := audit2.NewLogger(w)
		if err != nil {
			return err
		}

		entry := audit2.Entry{
			Operation: audit2Operation,
			Path:      audit2Path,
			DryRun:    audit2DryRun,
		}
		if err := logger.Log(entry); err != nil {
			return fmt.Errorf("write audit entry: %w", err)
		}
		return nil
	},
}

func init() {
	audit2Cmd.Flags().StringVar(&audit2Operation, "operation", "", "Operation name to record (required)")
	audit2Cmd.Flags().StringVar(&audit2Path, "path", "", "Vault secret path (required)")
	audit2Cmd.Flags().BoolVar(&audit2DryRun, "dry-run", false, "Mark entry as a dry-run")
	audit2Cmd.Flags().StringVar(&audit2Output, "output", "-", "Output file path (default: stdout)")
	_ = audit2Cmd.MarkFlagRequired("operation")
	_ = audit2Cmd.MarkFlagRequired("path")
	rootCmd.AddCommand(audit2Cmd)
}
