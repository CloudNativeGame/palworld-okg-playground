package cloudprovider

import (
	restclient "k8s.io/client-go/rest"
)

type CloudProvider interface {
	CreateCluster(options ClusterOptions) (KubernetesCluster, error)
	ListClusters() ([]KubernetesCluster, error)
	DeleteCluster(clusterId string) error
	GetCluster(clusterId string) (KubernetesCluster, error)
	GetKubernetesConfig(clusterId string) (*restclient.Config, error)
}

type KubernetesCluster interface {
	Healthy() (bool, error)
	ClusterId() string
	Description() string // tostring for cluster
	GetKubernetesConfig() (*restclient.Config, error)
}

type ClusterOptions interface {
	Options() map[string]string
}
