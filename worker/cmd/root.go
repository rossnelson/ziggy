package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("temporal-address", "localhost:7233", "Temporal server address")
	rootCmd.PersistentFlags().String("temporal-namespace", "default", "Temporal namespace")
	rootCmd.PersistentFlags().String("task-queue", "ziggy", "Temporal task queue")
	rootCmd.PersistentFlags().String("owner", "", "Owner name for this Ziggy instance")

	viper.BindPFlag("temporal-address", rootCmd.PersistentFlags().Lookup("temporal-address"))
	viper.BindPFlag("temporal-namespace", rootCmd.PersistentFlags().Lookup("temporal-namespace"))
	viper.BindPFlag("task-queue", rootCmd.PersistentFlags().Lookup("task-queue"))
	viper.BindPFlag("owner", rootCmd.PersistentFlags().Lookup("owner"))
}

func initConfig() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}
