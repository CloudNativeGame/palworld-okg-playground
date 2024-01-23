package pkg

import (
	"k8s.io/client-go/kubernetes"
)

type GameServerManager struct {
	client kubernetes.Interface // kubernetes client
}

func (gsm *GameServerManager) CreateGameServer() error {
	return nil
}

func (gsm *GameServerManager) ListGameServers() error {
	return nil
}

func (gsm *GameServerManager) DeleteGameServer() error {
	return nil
}

func NewGameServerManager(clusterId string) *GameServerManager {
	return &GameServerManager{}
}
