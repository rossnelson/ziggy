package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"ziggy/internal/api"
	"ziggy/internal/temporal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
}

func runServe(cmd *cobra.Command, args []string) error {
	address := viper.GetString("temporal-address")
	namespace := viper.GetString("temporal-namespace")
	taskQueue := viper.GetString("task-queue")
	owner := viper.GetString("owner")
	port, _ := cmd.Flags().GetInt("port")

	if owner == "" {
		owner = "dev"
	}
	workflowID := fmt.Sprintf("ziggy-%s", owner)
	chatWorkflowID := fmt.Sprintf("ziggy-chat-%s", owner)

	fmt.Printf("Starting Ziggy API server...\n")
	fmt.Printf("  Address: %s\n", address)
	fmt.Printf("  Namespace: %s\n", namespace)
	fmt.Printf("  Task Queue: %s\n", taskQueue)
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
