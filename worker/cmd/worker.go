package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ziggy/internal/ai"
	"ziggy/internal/temporal"
	"ziggy/internal/workflow"

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
}

func runWorker(cmd *cobra.Command, args []string) error {
	address := viper.GetString("temporal-address")
	namespace := viper.GetString("temporal-namespace")
	taskQueue := viper.GetString("task-queue")

	fmt.Printf("Starting Ziggy worker...\n")
	fmt.Printf("  Address: %s\n", address)
	fmt.Printf("  Namespace: %s\n", namespace)
	fmt.Printf("  Task Queue: %s\n", taskQueue)

	// Initialize the Temporal registry
	registry := temporal.Get()
	err := registry.Initialize(temporal.Config{
		HostPort:  address,
		Namespace: namespace,
		TaskQueue: taskQueue,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize temporal: %w", err)
	}
	defer registry.Shutdown()

	// Register workflows and activities
	workflow.RegisterWorkflows()

	aiClient := ai.NewClient()
	activities := workflow.NewActivities(aiClient)
	workflow.RegisterActivities(activities)

	chatActivities := workflow.NewChatActivities(aiClient, registry)
	workflow.RegisterChatActivities(chatActivities)

	if aiClient != nil {
		fmt.Println("  AI: enabled")
	} else {
		fmt.Println("  AI: disabled (no ANTHROPIC_API_KEY)")
	}

	// Create context that cancels on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down worker...")
		cancel()
	}()

	fmt.Println("Worker started. Press Ctrl+C to stop.")
	return registry.StartWorker(ctx)
}
