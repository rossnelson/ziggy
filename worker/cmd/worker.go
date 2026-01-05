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
	workerCmd.Flags().String("timezone", "America/Los_Angeles", "Timezone for time-of-day calculations")
	workerCmd.Flags().Bool("start-workflow", true, "Start the Ziggy workflow if not already running")
}

func runWorker(cmd *cobra.Command, args []string) error {
	address := viper.GetString("temporal-address")
	namespace := viper.GetString("temporal-namespace")
	taskQueue := viper.GetString("task-queue")
	owner := viper.GetString("owner")
	timezone, _ := cmd.Flags().GetString("timezone")
	startWorkflow, _ := cmd.Flags().GetBool("start-workflow")

	if owner == "" {
		owner = "dev"
	}

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

	// Start workflows if requested
	if startWorkflow {
		workflowID := fmt.Sprintf("ziggy-%s", owner)
		chatWorkflowID := fmt.Sprintf("ziggy-chat-%s", owner)
		needUpdaterID := workflowID + "-need-updater"
		poolRegeneratorID := workflowID + "-pool-regenerator"

		ensureWorkflow(ctx, registry, workflowID, workflow.ZiggyWorkflow, workflow.ZiggyInput{
			Owner:      owner,
			Timezone:   timezone,
			Generation: 1,
		})

		ensureWorkflow(ctx, registry, chatWorkflowID, workflow.ChatWorkflow, workflow.ChatInput{
			Owner:   owner,
			ZiggyID: workflowID,
		})

		ensureWorkflow(ctx, registry, needUpdaterID, workflow.NeedUpdaterWorkflow, workflow.NeedUpdaterInput{
			ZiggyWorkflowID: workflowID,
			Iteration:       0,
		})

		ensureWorkflow(ctx, registry, poolRegeneratorID, workflow.PoolRegeneratorWorkflow, workflow.PoolRegeneratorInput{
			ZiggyWorkflowID: workflowID,
		})
	}

	fmt.Println("Worker started. Press Ctrl+C to stop.")
	return registry.StartWorker(ctx)
}

func ensureWorkflow(ctx context.Context, registry *temporal.Registry, workflowID string, workflowFunc interface{}, input interface{}) {
	status, err := registry.DescribeWorkflow(ctx, workflowID)
	if err == nil && status.Status == "WORKFLOW_EXECUTION_STATUS_RUNNING" {
		fmt.Printf("Workflow %s already running\n", workflowID)
		return
	}

	if err == nil {
		fmt.Printf("Workflow %s is %s, restarting...\n", workflowID, status.Status)
	}

	_, err = registry.ExecuteWorkflow(ctx, workflowID, workflowFunc, input)
	if err != nil {
		fmt.Printf("Note: %v (workflow %s)\n", err, workflowID)
	} else {
		fmt.Printf("Started workflow: %s\n", workflowID)
	}
}
