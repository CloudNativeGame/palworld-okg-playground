package cmd

import (
	"context"
	"fmt"
	"github.com/CloudNativeGame/palworld-okg-playground/cloudprovider/alibabacloud"
	"github.com/CloudNativeGame/palworld-okg-playground/pkg/cluster"
	"github.com/CloudNativeGame/palworld-okg-playground/pkg/env"
	"github.com/liushuochen/gotable"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"strconv"
	"strings"
)

func init() {
	//clusterCmd.Flags().String("provider", "", "cloud provider")
	//clusterCmd.Flags().String("access_key_id", "", "access_key_id")
	//clusterCmd.Flags().String("access_key_secret", "", "access_key_secret")
	//clusterCmd.Flags().String("region_id", "cn-hangzhou", "cluster_type for which to be created (options: ManagedKubernetes, serverless)")

	createClusterCmd.Flags().String("cluster_type", "ManagedKubernetes", "cluster type for which to be created (options: ManagedKubernetes, serverless)")
	createClusterCmd.Flags().String("cluster_name", "", "cluster name for which to be created")
	createClusterCmd.Flags().String("region", "cn-hangzhou", "cluster region for which to be created")
	createClusterCmd.Flags().String("vpcid", "", "vpc id")
	createClusterCmd.Flags().String("vswitch_ids", "", "example:vsw1 OR vsw1,vsw2 if multi vswitch ids needed")

	clusterCmd.AddCommand(createClusterCmd)
	clusterCmd.AddCommand(deleteClusterCmd)
	clusterCmd.AddCommand(listClusterCmd)
	clusterCmd.AddCommand(chooseClusterCmd)
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage PalWorld game clusters",
	Long:  `Manage PalWorld game clusters`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		clusterManager, err := cluster.NewClusterManager()
		if err != nil {
			klog.Errorf("failed to create cluster manager, err: %s", err.Error())
			return
		}
		ctx = context.WithValue(ctx, "clusterManager", clusterManager)
		envFile := env.NewEnvFile()
		if !envFile.Exists() {
			err := envFile.Create()
			if err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "create envFile failed, because %s \n", err.Error())
			}
		}
		cmd.SetContext(ctx)
		return
	},
}

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PalWorld gameserver cluster",
	Long:  `Create a new PalWorld gameserver cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		clusterManager := cmd.Context().Value("clusterManager").(*cluster.ClusterManager)

		clusterType, _ := cmd.Flags().GetString("cluster_type")
		clusterName, _ := cmd.Flags().GetString("cluster_name")
		region, _ := cmd.Flags().GetString("region")
		vpcid, _ := cmd.Flags().GetString("vpcid")
		vswids, _ := cmd.Flags().GetString("vswitch_ids")
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "got input params: clusterType %s clusterName %s region %s vpcid %s vswitch_ids %s \n", clusterType, clusterName, region, vpcid, vswids)
		gsLoadBalancerId, err := clusterManager.CreateGameServerLoadBalancer()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "failed to create loadbalancer, err %s\n", err.Error())
			return
		}

		conf := &alibabacloud.ClusterConfig{
			ClusterType: clusterType,
			Name:        clusterName,
			RegionId:    region,
			VpcId:       vpcid,
			VswitchIds:  strings.Trim(vswids, " "),
		}
		cluster, err := clusterManager.CreateCluster(conf)
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "failed to create %s cluster, err %s\n", clusterType, err.Error())
			return
		}
		clusterId := cluster.ClusterId()
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "cluster %s (type %s) is creating\n", clusterId, clusterType)

		err = env.NewEnvFile().AddNewCluster(&env.ClusterConfiguration{
			ID:                       clusterId,
			Name:                     clusterName,
			Default:                  false,
			GameServerLoadBalancerId: gsLoadBalancerId,
		})
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Add a new cluster %s to EnvFile failed, because %s. \n", clusterId, err.Error())
		}

		return
	},
}

var deleteClusterCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a PalWorld gameserver cluster",
	Long:  "Delete a PalWorld gameserver cluster",
	Run: func(cmd *cobra.Command, args []string) {
		clusterManager := cmd.Context().Value("clusterManager").(*cluster.ClusterManager)
		clusterManager.DeleteCluster()
		return
	},
}

var listClusterCmd = &cobra.Command{
	Use:   "list",
	Short: "List all PalWorld gameserver cluster",
	Long:  `List all PalWorld gameserver cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		clusterManager := cmd.Context().Value("clusterManager").(*cluster.ClusterManager)

		envfile := env.NewEnvFile()
		clusterMap, err := envfile.ListAll()
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "list all cluster failed, because %s \n", err.Error())
		}

		table, err := gotable.Create("Name", "Id", "State", "GameServerLBID", "IsDefault")
		if err != nil {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Create table failed: %s\n", err.Error())
			return
		}
		for _, cluster := range clusterMap.Clusters {
			state, err := clusterManager.GetClusterState(cluster.ID)
			if err != nil {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "failed to list all cluster, because %s \n", err.Error())
				return
			}

			table.AddRow([]string{cluster.Name, cluster.ID, state, cluster.GameServerLoadBalancerId, strconv.FormatBool(cluster.Default)})
		}
		_, _ = fmt.Fprint(cmd.OutOrStdout(), table)

		return
	},
}

var chooseClusterCmd = &cobra.Command{
	Use:   "choose",
	Short: "Choose a PalWorld gameserver cluster",
	Long:  `Choose a PalWorld gameserver cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			klog.Errorf("cluster id is nil, you should input a cluster id")
		}
		clusterId := args[0]
		envFile := env.NewEnvFile()
		envFile.UpdateDefaultCluster(clusterId)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "choose cluster %s successfully. \n", args[0])
	},
}

//func init() {
//	createClusterCmd.Flags().StringP("name", "n", "", "Name of the cluster")
//	createClusterCmd.Flags().StringP("region", "r", "", "Region of the cluster")
//	clusterCmd.AddCommand(createClusterCmd)
//	clusterCmd.AddCommand(deleteClusterCmd)
//}
