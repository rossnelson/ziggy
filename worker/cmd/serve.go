package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ziggy/internal/api"
	"ziggy/internal/temporal"
	"ziggy/internal/workflow"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	Long:  `Starts an HTTP API server that proxies requests to a running Ziggy workflow.`,
	RunE:  runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().Int("port", 8080, "HTTP server port")
	serveCmd.Flags().String("timezone", "America/Los_Angeles", "Timezone for time-of-day calculations")
	serveCmd.Flags().Bool("start-workflow", true, "Start the Ziggy workflow if not already running")
}

func runServe(cmd *cobra.Command, args []string) error {
	address, _ := cmd.Flags().GetString("temporal-address")
	namespace, _ := cmd.Flags().GetString("temporal-namespace")
	taskQueue, _ := cmd.Flags().GetString("task-queue")
	owner, _ := cmd.Flags().GetString("owner")
	port, _ := cmd.Flags().GetInt("port")
	timezone, _ := cmd.Flags().GetString("timezone")
	startWorkflow, _ := cmd.Flags().GetBool("start-workflow")

	if owner == "" {
		owner = "dev"
	}
	workflowID := fmt.Sprintf("ziggy-%s", owner)
	chatWorkflowID := fmt.Sprintf("ziggy-chat-%s", owner)

	fmt.Printf("Starting Ziggy API server...\n")
	fmt.Printf("  Address: %s\n", address)
	fmt.Printf("  Namespace: %s\n", namespace)
	fmt.Printf("  Task Queue: %s\n", taskQueue)
	fmt.Printf("  Workflow ID: %s\n", workflowID)
	fmt.Printf("  Port: %d\n", port)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the workflows if requested
	if startWorkflow {
		_, err := registry.ExecuteWorkflow(ctx, workflowID, workflow.ZiggyWorkflow, workflow.ZiggyInput{
			Owner:      owner,
			Timezone:   timezone,
			Generation: 1,
		})
		if err != nil {
			fmt.Printf("Note: %v (workflow may already be running)\n", err)
		} else {
			fmt.Printf("Started new Ziggy workflow: %s\n", workflowID)
		}

		_, err = registry.ExecuteWorkflow(ctx, chatWorkflowID, workflow.ChatWorkflow, workflow.ChatInput{
			Owner:   owner,
			ZiggyID: workflowID,
		})
		if err != nil {
			fmt.Printf("Note: %v (chat workflow may already be running)\n", err)
		} else {
			fmt.Printf("Started new Chat workflow: %s\n", chatWorkflowID)
		}

		needUpdaterID := workflowID + "-need-updater"
		_, err = registry.ExecuteWorkflow(ctx, needUpdaterID, workflow.NeedUpdaterWorkflow, workflow.NeedUpdaterInput{
			ZiggyWorkflowID: workflowID,
			Iteration:       0,
		})
		if err != nil {
			fmt.Printf("Note: %v (need updater may already be running)\n", err)
		} else {
			fmt.Printf("Started new NeedUpdater workflow: %s\n", needUpdaterID)
		}
	}

	// Handle shutdown signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down server...")
		cancel()
	}()

	// Start the API server
	server := api.NewServer(registry, workflowID, chatWorkflowID, port)
	return server.Start(ctx)
}
