package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultpatch/internal/lock"
	"github.com/your-org/vaultpatch/internal/vault"
)

var (
	lockOwner  string
	lockTTL    time.Duration
	lockDryRun bool
	unlockPath string
)

var lockCmd = &cobra.Command{
	Use:   "lock <path>",
	Short: "Acquire an advisory lock on a Vault secret path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}
		res, err := lock.Lock(client, args[0], lockOwner, lockTTL, lockDryRun)
		if err != nil {
			return err
		}
		if lockDryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would lock %q for owner %q (ttl: %s)\n", res.Path, lockOwner, lockTTL)
		} else {
			fmt.Fprintf(os.Stdout, "locked %q for owner %q until %s\n", res.Path, lockOwner, res.Entry.ExpiresAt.Format(time.RFC3339))
		}
		return nil
	},
}

var unlockCmd = &cobra.Command{
	Use:   "unlock <path>",
	Short: "Release an advisory lock on a Vault secret path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vault.NewClient(vault.Params{})
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}
		res, err := lock.Unlock(client, args[0], lockOwner, lockDryRun)
		if err != nil {
			return err
		}
		if res.DryRun {
			fmt.Fprintf(os.Stdout, "[dry-run] would unlock %q\n", res.Path)
		} else {
			fmt.Fprintf(os.Stdout, "unlocked %q\n", res.Path)
		}
		return nil
	},
}

func init() {
	lockCmd.Flags().StringVar(&lockOwner, "owner", "", "Owner identifier for the lock (required)")
	lockCmd.Flags().DurationVar(&lockTTL, "ttl", 10*time.Minute, "Lock TTL duration")
	lockCmd.Flags().BoolVar(&lockDryRun, "dry-run", false, "Simulate without writing")
	_ = lockCmd.MarkFlagRequired("owner")

	unlockCmd.Flags().StringVar(&lockOwner, "owner", "", "Owner identifier")
	unlockCmd.Flags().BoolVar(&lockDryRun, "dry-run", false, "Simulate without deleting")

	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
}
