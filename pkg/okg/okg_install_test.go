package okg

import (
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

func TestInstallOpenKruiseGame(t *testing.T) {
	cfg, err := clientcmd.BuildConfigFromFlags("", "/Users/liuqiuyang/.kube/goatscaler.conf")
	if err != nil {
		panic(err.Error())
	}

	//cfg := config.GetConfigOrDie()
	err = InstallOpenKruiseGame(cfg)
	if err != nil {
		t.Errorf("Install OKG failed, because of %s", err.Error())
	}
}
