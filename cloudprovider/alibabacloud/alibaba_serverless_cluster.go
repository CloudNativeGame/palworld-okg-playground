package alibabacloud

import restclient "k8s.io/client-go/rest"

type ServerlessCluster struct {
	Id     string
	Status string
}

func (m *ServerlessCluster) Healthy() (bool, error) {
	return true, nil
}
func (m *ServerlessCluster) ClusterId() string {
	return m.Id
}
func (m *ServerlessCluster) Description() string {
	return ""
}

func (m *ServerlessCluster) GetKubernetesConfig() (*restclient.Config, error) {
	return nil, nil
}
