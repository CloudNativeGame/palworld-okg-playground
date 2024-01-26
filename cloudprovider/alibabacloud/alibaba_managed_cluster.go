package alibabacloud

import restclient "k8s.io/client-go/rest"

type ManagedCluster struct {
	Id     string
	Status string
}

func (m *ManagedCluster) Healthy() (bool, error) {
	return true, nil
}
func (m *ManagedCluster) ClusterId() string {
	return m.Id
}
func (m *ManagedCluster) Description() string {
	return ""
}

func (m *ManagedCluster) GetKubernetesConfig() (*restclient.Config, error) {
	return nil, nil
}
