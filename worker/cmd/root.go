package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ziggy",
	Short: "Ziggy virtual pet Temporal worker",
	Long: `A Temporal worker for the Ziggy virtual pet game.

Ziggy is a tardigrade virtual pet whose entire lifecycle runs as a
durable Temporal workflow. This CLI manages the worker and provides
commands for interacting with Ziggy.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("temporal-address", "localhost:7233", "Temporal server address")
	rootCmd.PersistentFlags().String("temporal-namespace", "default", "Temporal namespace")
	rootCmd.PersistentFlags().String("task-queue", "ziggy", "Temporal task queue")
	rootCmd.PersistentFlags().String("owner", "", "Owner name for this Ziggy instance")
}
