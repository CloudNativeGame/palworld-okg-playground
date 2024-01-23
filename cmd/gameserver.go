package cmd

import (
	"context"
	"github.com/CloudNativeGame/palworld-okg-playground/pkg"
	"github.com/spf13/cobra"
)

var gameserverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage PalWorld game servers",
	Long:  `Manage PalWorld game servers`,
	Run: func(cmd *cobra.Command, args []string) {
		clusterId := cmd.Context().Value("clusterId").(string)
		gameserverManager := pkg.NewGameServerManager(clusterId)
		ctx := cmd.Context()
		ctx = context.WithValue(ctx, "gameserverManager", gameserverManager)
		cmd.SetContext(ctx)
	},
}

var createGameserverCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PalWorld gameserver",
	Long:  `Create a new PalWorld gameserver`,
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

var listGameserverCmd = &cobra.Command{
	Use:   "list",
	Short: "List PalWorld gameservers",
	Long:  `List PalWorld gameservers`,
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

var describeGameserverCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a PalWorld gameserver",
	Long:  `Describe a PalWorld gameserver`,
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

var deleteGameserverCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a PalWorld gameserver",
	Long:  "Delete a PalWorld gameserver",
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

var upgradeGameserverCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade a PalWorld gameserver",
	Long:  "Upgrade a PalWorld gameserver",
	Run: func(cmd *cobra.Command, args []string) {
		return
	},
}

func init() {
	describeGameserverCmd.Flags().StringP("name", "n", "", "Name of the gameserver")
	upgradeGameserverCmd.Flags().StringP("name", "n", "", "Name of the gameserver")
	deleteGameserverCmd.Flags().StringP("name", "n", "", "Name of the gameserver")
	gameserverCmd.AddCommand(createGameserverCmd)
	gameserverCmd.AddCommand(listGameserverCmd)
	gameserverCmd.AddCommand(describeGameserverCmd)
	gameserverCmd.AddCommand(deleteGameserverCmd)
	gameserverCmd.AddCommand(upgradeGameserverCmd)
}
