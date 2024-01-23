package cmd

import (
	"context"
	"github.com/CloudNativeGame/palworld-okg-playground/pkg"
	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage PalWorld game clusters",
	Long:  `Manage PalWorld game clusters`,
	Run: func(cmd *cobra.Command, args []string) {
		clusterManager := pkg.NewClusterManager()
		ctx := context.WithValue(cmd.Context(), "clusterManager", clusterManager)
		cmd.SetContext(ctx)
		return
	},
}

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PalWorld gameserver cluster",
	Long:  `Create a new PalWorld gameserver cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		//clusterManager := cmd.Context().Value("clusterManager").(*pkg.ClusterManager)
		//cluster, err := clusterManager.CreateCluster()

		return
	},
}

var deleteClusterCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a PalWorld gameserver cluster",
	Long:  "Delete a PalWorld gameserver cluster",
	Run: func(cmd *cobra.Command, args []string) {
		clusterManager := cmd.Context().Value("clusterManager").(*pkg.ClusterManager)
		clusterManager.DeleteCluster()
		return
	},
}

func init() {
	createClusterCmd.Flags().StringP("name", "n", "", "Name of the cluster")
	createClusterCmd.Flags().StringP("region", "r", "", "Region of the cluster")
	clusterCmd.AddCommand(createClusterCmd)
	clusterCmd.AddCommand(deleteClusterCmd)
}
