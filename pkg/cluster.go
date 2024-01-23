package pkg

import "github.com/CloudNativeGame/palworld-okg-playground/cloudprovider"

func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		provider: nil,
	}
}

type ClusterManager struct {
	provider cloudprovider.CloudProvider
}

func (cm *ClusterManager) CreateCluster() (cloudprovider.KubernetesCluster, error) {
	return cm.provider.CreateCluster(nil)
}

func (cm *ClusterManager) DeleteCluster() error {
	return nil
}

func (cm *ClusterManager) GetKubernetesConfig() error {
	return nil
}
