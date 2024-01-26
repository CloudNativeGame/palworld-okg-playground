package cmd

import (
	"context"
	"fmt"
	"github.com/CloudNativeGame/palworld-okg-playground/pkg/cluster"
	"github.com/CloudNativeGame/palworld-okg-playground/pkg/gameserver"
	"github.com/liushuochen/gotable"
	gamekruisev1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sort"
	"time"
)

var gameserverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage PalWorld game servers",
	Long:  `Manage PalWorld game servers`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		config := config.GetConfigOrDie()
		clusterId := ctx.Value("clusterId")
		if clusterId != nil {
			clusterManager := ctx.Value("clusterManager").(*cluster.ClusterManager)
			config = clusterManager.GetKubernetesConfig()
		}

		gameserverManager := gameserver.NewGameServerManager(config, "")
		gameserverManager.EnsureOKGInstalled()
		ctx = context.WithValue(ctx, "gameserverManager", gameserverManager)
		cmd.SetContext(ctx)
	},
}

var createGameserverCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PalWorld gameserver",
	Long:  `Create a new PalWorld gameserver`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		gameserverManager := ctx.Value("gameserverManager").(*gameserver.GameServerManager)
		err := gameserverManager.CreateGameServer()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "create a new game server failed, because %s \n", err.Error())
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "create a new game server successfully. \n")
		}
		return
	},
}

var listGameserverCmd = &cobra.Command{
	Use:   "list",
	Short: "List PalWorld gameservers",
	Long:  `List PalWorld gameservers`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		gameserverManager := ctx.Value("gameserverManager").(*gameserver.GameServerManager)
		gameservers, err := gameserverManager.ListGameServers()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "list game servers list failed, because %s \n", err.Error())
		}

		table, err := gotable.Create("Name", "State", "OpsState", "NetworkState", "ResourceType", "Address", "Age")
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Create table failed: %s\n", err.Error())
			return
		}

		sortedGs := gameserver.SortGs(gameservers)
		sort.Sort(sortedGs)
		for _, gs := range sortedGs {
			addr := ""
			if gs.Status.NetworkStatus.CurrentNetworkState == gamekruisev1alpha1.NetworkReady {
				externalAddress := gs.Status.NetworkStatus.ExternalAddresses[0]
				addr = fmt.Sprintf("%s:%s", externalAddress.IP, externalAddress.Ports[0].Port)
			}
			age := time.Since(gs.CreationTimestamp.Time).String()
			resourceType := gameserver.ToResourceType(gs.Spec.Containers)

			table.AddRow([]string{gs.Name, string(gs.Status.CurrentState), string(gs.Spec.OpsState), string(gs.Status.NetworkStatus.CurrentNetworkState), string(resourceType), addr, age})
		}
		_, _ = fmt.Fprint(cmd.OutOrStdout(), table)

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
		ctx := cmd.Context()
		gameserverManager := ctx.Value("gameserverManager").(*gameserver.GameServerManager)
		gsName := cmd.Flag("name").Value.String()
		if gsName == "" {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "game server name shoul not be empty. \n")
			return
		}
		err := gameserverManager.DeleteGameServer(gsName)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "delete game server %s failed, because %s \n", gsName, err.Error())
		} else {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "delete game server %s successfully. \n", gsName)
		}
		return
	},
}

var upgradeGameserverCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade a PalWorld gameserver",
	Long:  "Upgrade a PalWorld gameserver",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		gameserverManager := ctx.Value("gameserverManager").(*gameserver.GameServerManager)
		gsName := cmd.Flag("name").Value.String()
		if gsName == "" {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "game server name shoul not be empty. \n")
			return
		}

		resourceType := cmd.Flag("resources").Value.String()
		if resourceType != "" {
			err := gameserverManager.UpgradeGameServerResources(gsName, resourceType)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "game server %s upgrade resources failed, because %s. \n", gsName, err.Error())
			} else {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "game server %s upgrade resources successfully. \n", gsName)
			}
		}

		players := cmd.Flag("players").Value.String()
		if players != "" {
			err := gameserverManager.UpgradeGameServerEnvConfig(gsName, &gameserver.EnvConfig{Players: players})
			if err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "game server %s upgrade envConfig failed, because %s. \n", gsName, err.Error())
			} else {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "game server %s upgrade envConfig successfully. \n", gsName)
			}
		}

		return
	},
}

func init() {
	describeGameserverCmd.Flags().StringP("name", "n", "", "Name of the gameserver")
	upgradeGameserverCmd.Flags().StringP("name", "n", "", "Name of the gameserver")
	upgradeGameserverCmd.Flags().StringP("resources", "r", "", "Resources standard of the gameserver. You can choose small(4cpu & 8Gi), medium(4cpu & 16Gi) or large(4cpu & 32Gi)")
	upgradeGameserverCmd.Flags().StringP("players", "p", "", "Max amount of players that are able to join the gameserver. You can input number in range [1-31]")
	deleteGameserverCmd.Flags().StringP("name", "n", "", "Name of the gameserver")
	gameserverCmd.AddCommand(createGameserverCmd)
	gameserverCmd.AddCommand(listGameserverCmd)
	gameserverCmd.AddCommand(describeGameserverCmd)
	gameserverCmd.AddCommand(deleteGameserverCmd)
	gameserverCmd.AddCommand(upgradeGameserverCmd)
}
