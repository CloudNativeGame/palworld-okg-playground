package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(createClusterCmd)
	rootCmd.AddCommand(gameserverCmd)
	rootCmd.AddCommand(playerCmd)
}

var rootCmd = &cobra.Command{
	Use:   "pal",
	Short: "pal is a fast command to create and manage PalWorld game servers",
	Long:  "pal is a fast command to create and manage PalWorld game servers",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		envFile := NewEnvFile()
		if envFile.Exists() {
			clusterId, _ := envFile.Read()
			if clusterId != "" {
				ctx := context.WithValue(cmd.Context(), "clusterId", clusterId)
				cmd.SetContext(ctx)
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
