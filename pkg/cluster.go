package pkg

import (
	"github.com/CloudNativeGame/palworld-okg-playground/cloudprovider"
	"github.com/CloudNativeGame/palworld-okg-playground/cloudprovider/alibabacloud"
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

func (cm *ClusterManager) GetKubernetesConfig() error {
	return nil
}
