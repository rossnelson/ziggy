package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ziggy/internal/registry"
	_ "ziggy/internal/workflow"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the Temporal worker",
	Long:  `Starts a Temporal worker that processes Ziggy workflows and activities.`,
	RunE:  runWorker,
}

func init() {
	rootCmd.AddCommand(workerCmd)
	workerCmd.Flags().String("timezone", "America/Los_Angeles", "Timezone for time-of-day calculations")
	workerCmd.Flags().Bool("start-workflow", true, "Start the Ziggy workflow if not already running")
}

func runWorker(cmd *cobra.Command, args []string) error {
	reg := registry.NewRegistry()

	startWorkflow, _ := cmd.Flags().GetBool("start-workflow")
	timezone, _ := cmd.Flags().GetString("timezone")
	owner := viper.GetString("owner")
	if owner == "" {
		owner = "dev"
	}

	err := reg.Initialize(registry.Config{
		HostPort:      viper.GetString("temporal-address"),
		Namespace:     viper.GetString("temporal-namespace"),
		TaskQueue:     viper.GetString("task-queue"),
		Owner:         owner,
		Timezone:      timezone,
		StartWorkflow: startWorkflow,
	})
	if err != nil {
		return fmt.Errorf("initialize registry: %w", err)
	}

	ctx, cancel := handleSignals()
	defer cancel()

	fmt.Println("Worker started. Press Ctrl+C to stop.")
	return reg.StartWorker(ctx)
}

func handleSignals() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down worker...")
		cancel()
	}()

	return ctx, cancel
}
