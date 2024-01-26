package cmd

import (
	"context"
	"github.com/CloudNativeGame/palworld-okg-playground/cloudprovider/alibabacloud"
	"github.com/CloudNativeGame/palworld-okg-playground/pkg"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
	"strings"
)

func init() {
	clusterCmd.Flags().String("provider", "", "cloud provider")
	clusterCmd.Flags().String("access_key_id", "", "access_key_id")
	clusterCmd.Flags().String("access_key_secret", "", "access_key_secret")
	clusterCmd.Flags().String("region_id", "cn-hangzhou", "cluster_type for which to be created (options: ManagedKubernetes, serverless)")

	createClusterCmd.Flags().String("cluster_type", "ManagedKubernetes", "cluster type for which to be created (options: ManagedKubernetes, serverless)")
	createClusterCmd.Flags().String("cluster_name", "", "cluster name for which to be created")
	createClusterCmd.Flags().String("region", "", "cluster region for which to be created")
	createClusterCmd.Flags().String("vpcid", "", "vpc id")
	createClusterCmd.Flags().String("vswitch_ids", "", "example:vsw1 OR vsw1,vsw2 if multi vswitch ids needed")

	clusterCmd.AddCommand(createClusterCmd)
	clusterCmd.AddCommand(deleteClusterCmd)
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage PalWorld game clusters",
	Long:  `Manage PalWorld game clusters`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		clusterManager, err := pkg.NewClusterManager()
		if err != nil {
			klog.Errorf("failed to create cluster manager, err: %s", err.Error())
			return
		}
		klog.Infof("got cluster manager %v", clusterManager)
		ctx = context.WithValue(cmd.Context(), "clusterManager", clusterManager)
		cmd.SetContext(ctx)
		klog.Infof("recheck cluster manager %v", cmd.Context().Value("clusterManager"))
		return
	},
}

var createClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new PalWorld gameserver cluster",
	Long:  `Create a new PalWorld gameserver cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		var clusterManager *pkg.ClusterManager
		var err error
		//clusterManager := cmd.Context().Value("clusterManager").(*pkg.ClusterManager)
		if cmd.Context().Value("clusterManager") != nil {
			clusterManager = cmd.Context().Value("clusterManager").(*pkg.ClusterManager)
			klog.Infof("clusterManager exists %v", clusterManager)
		} else {
			klog.Infof("clusterManager is nil, try to new one")
			clusterManager, err = pkg.NewClusterManager()
			if err != nil {
				klog.Errorf("failed to create cluster manager, err: %s", err.Error())
				return
			}
		}

		clusterType, _ := cmd.Flags().GetString("cluster_type")
		clusterName, _ := cmd.Flags().GetString("cluster_name")
		region, _ := cmd.Flags().GetString("region")
		vpcid, _ := cmd.Flags().GetString("vpcid")
		vswids, _ := cmd.Flags().GetString("vswitch_ids")
		klog.Infof("got input params: clusterType %s clusterName %s region %s vpcid %s vswitch_ids %s ", clusterType, clusterName, region, vpcid, vswids)
		conf := &alibabacloud.ClusterConfig{
			ClusterType: clusterType,
			Name:        clusterName,
			RegionId:    region,
			VpcId:       vpcid,
			VswitchIds:  strings.Trim(vswids, " "),
		}
		cluster, err := clusterManager.CreateCluster(conf)
		if err != nil {
			klog.Errorf("failed to create %s cluster, err %s", clusterType, err.Error())
		} else {
			klog.Infof("cluster %s (type %s) is creating", cluster.ClusterId(), clusterType)
		}

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

//func init() {
//	createClusterCmd.Flags().StringP("name", "n", "", "Name of the cluster")
//	createClusterCmd.Flags().StringP("region", "r", "", "Region of the cluster")
//	clusterCmd.AddCommand(createClusterCmd)
//	clusterCmd.AddCommand(deleteClusterCmd)
//}
