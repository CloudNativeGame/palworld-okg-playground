package cluster

import (
	"github.com/CloudNativeGame/palworld-okg-playground/cloudprovider"
	"github.com/CloudNativeGame/palworld-okg-playground/cloudprovider/alibabacloud"
	restclient "k8s.io/client-go/rest"
)

func NewClusterManager() (*ClusterManager, error) {
	defaultProvider, err := alibabacloud.CreateAlibabaCloudManager(&alibabacloud.CloudConfig{})
	return &ClusterManager{
		provider: defaultProvider,
	}, err
}

type ClusterManager struct {
	provider cloudprovider.CloudProvider
}

func (cm *ClusterManager) CreateCluster(options cloudprovider.ClusterOptions) (cloudprovider.KubernetesCluster, error) {
	return cm.provider.CreateCluster(options)
}

func (cm *ClusterManager) DeleteCluster() error {
	return nil
}

func (cm *ClusterManager) GetKubernetesConfig(clusterId string) (*restclient.Config, error) {
	return cm.provider.GetKubernetesConfig(clusterId)
}

func (cm *ClusterManager) GetClusterState(clusterId string) (string, error) {
	return cm.provider.GetClusterState(clusterId)
}

func (cm *ClusterManager) CreateGameServerLoadBalancer() (string, error) {
	return cm.provider.CreateGameServerLoadBalancer()
}
