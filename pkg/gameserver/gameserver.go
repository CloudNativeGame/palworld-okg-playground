package gameserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/CloudNativeGame/palworld-okg-playground/pkg/okg"
	kruiseV1beta1 "github.com/openkruise/kruise-api/apps/v1beta1"
	gamekruisev1alpha1 "github.com/openkruise/kruise-game/apis/v1alpha1"
	"github.com/openkruise/kruise-game/cloudprovider/alibabacloud"
	kruisegameclientset "github.com/openkruise/kruise-game/pkg/client/clientset/versioned"
	kruiseV1alpha1 "github.com/openkruise/kruise/apis/apps/v1alpha1"
	kruiseclientset "github.com/openkruise/kruise/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/utils/pointer"
	"strconv"
	"time"
)

const (
	GSName            = "palworld-server"
	GSImage           = "thijsvanloef/palworld-server-docker:latest"
	PlayersAnnotation = "players"
	DefaultPlayers    = "16"
)

type GameServerManager struct {
	client           kubernetes.Interface // kubernetes client
	kruisegameClient kruisegameclientset.Interface
	kruseClient      kruiseclientset.Interface
	cfg              *restclient.Config
	slbId            string
}

func (gsm *GameServerManager) CreateGameServer() error {
	var gss *gamekruisev1alpha1.GameServerSet
	err := fmt.Errorf("init an error")

	for err != nil {
		gss, err = gsm.kruisegameClient.GameV1alpha1().GameServerSets("default").Get(context.Background(), "palworld-server", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			_, _ = gsm.kruisegameClient.GameV1alpha1().GameServerSets("default").Create(context.Background(), defaultGameServerSet(gsm.slbId), metav1.CreateOptions{})
		}
	}

	newReplicas := *gss.Spec.Replicas + 1
	gss.Spec.Replicas = &newReplicas
	_, err = gsm.kruisegameClient.GameV1alpha1().GameServerSets("default").Update(context.Background(), gss, metav1.UpdateOptions{})
	return err
}

func (gsm *GameServerManager) ListGameServers() ([]gamekruisev1alpha1.GameServer, error) {
	labelSelector := labels.SelectorFromSet(map[string]string{
		gamekruisev1alpha1.GameServerOwnerGssKey: "palworld-server",
	}).String()
	gssList, err := gsm.kruisegameClient.GameV1alpha1().GameServers("default").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}

	return gssList.Items, nil
}

type ResourceStandard struct {
	RequestCpu string
	RequestMem string
	LimitCpu   string
	LimitMem   string
}

type ResourceType string

const (
	SmallType  = "small"
	MediumType = "medium"
	LargeType  = "large"
)

func DefaultResourceStandards() map[ResourceType]ResourceStandard {
	return map[ResourceType]ResourceStandard{
		SmallType: {
			RequestCpu: "3",
			RequestMem: "7Gi",
			LimitCpu:   "4",
			LimitMem:   "8Gi",
		},
		MediumType: {
			RequestCpu: "3",
			RequestMem: "15Gi",
			LimitCpu:   "4",
			LimitMem:   "16Gi",
		},
		LargeType: {
			RequestCpu: "3",
			RequestMem: "30Gi",
			LimitCpu:   "4",
			LimitMem:   "32Gi",
		},
	}
}

func ToResourceType(containers []gamekruisev1alpha1.GameServerContainer) ResourceType {
	if len(containers) == 0 {
		return SmallType
	}
	switch containers[0].Resources.Limits.Memory().String() {
	case "32Gi":
		return LargeType
	case "16Gi":
		return MediumType
	}
	return SmallType
}

func ToResourceRequirements(resourceStandard ResourceStandard) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(resourceStandard.LimitCpu),
			corev1.ResourceMemory: resource.MustParse(resourceStandard.LimitMem),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(resourceStandard.RequestCpu),
			corev1.ResourceMemory: resource.MustParse(resourceStandard.RequestMem),
		},
	}
}

func (gsm *GameServerManager) UpgradeGameServerResources(name string, resourceType string) error {
	// get gs & check
	gs, err := gsm.kruisegameClient.GameV1alpha1().GameServers("default").Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	resourceStandard, exist := DefaultResourceStandards()[ResourceType(resourceType)]
	if !exist {
		return fmt.Errorf("Input wrong resource type. You must pick one of {small, medium, large}. \n")
	}

	if ToResourceType(gs.Spec.Containers) == ResourceType(resourceType) {
		return fmt.Errorf("gameserver %s resourceType not changed. Do not upgrade", name)
	}

	// resources diff, update gs
	if len(gs.Spec.Containers) == 0 {
		gs.Spec.Containers = []gamekruisev1alpha1.GameServerContainer{
			{
				Name: GSName,
			},
		}
	}
	gs.Spec.Containers[0].Resources = ToResourceRequirements(resourceStandard)
	_, err = gsm.kruisegameClient.GameV1alpha1().GameServers("default").Update(context.Background(), gs, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	// delete pod
	return gsm.client.CoreV1().Pods("default").Delete(context.Background(), name, metav1.DeleteOptions{})
}

type EnvConfig struct {
	Players string
}

func (gsm *GameServerManager) UpgradeGameServerEnvConfig(name string, envConfig *EnvConfig) error {
	// patch annotation to pod
	patchData := map[string]interface{}{"metadata": map[string]interface{}{"annotations": map[string]string{PlayersAnnotation: envConfig.Players}}}
	patchPodBytes, err := json.Marshal(patchData)
	if err != nil {
		return err
	}
	_, err = gsm.client.CoreV1().Pods("default").Patch(context.TODO(), name, types.StrategicMergePatchType, patchPodBytes, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	// restart pod container
	crr := convContainerRecreateRequest(name)
	_, err = gsm.kruseClient.AppsV1alpha1().ContainerRecreateRequests("default").Create(context.Background(), crr, metav1.CreateOptions{})
	return err
}

func (gsm *GameServerManager) DeleteGameServer(name string) error {
	osJson := map[string]interface{}{"spec": map[string]string{"opsState": string(gamekruisev1alpha1.Kill)}}
	data, err := json.Marshal(osJson)
	if err != nil {
		return err
	}
	_, err = gsm.kruisegameClient.GameV1alpha1().GameServers("default").Patch(context.Background(), name, types.MergePatchType, data, metav1.PatchOptions{})
	return err
}

func (gsm *GameServerManager) EnsureOKGInstalled() {
	err := fmt.Errorf("init an error")
	for err != nil {
		_, err = gsm.client.CoreV1().Namespaces().Get(context.Background(), "kruise-system", metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("OKG is not installed. Installing...\n")
				err = okg.InstallOpenKruiseGame(gsm.cfg)
			} else {
				continue
			}
		}
		_, err = gsm.client.CoreV1().Namespaces().Get(context.Background(), "kruise-game-system", metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("OKG is not installed. Installing...\n")
				err = okg.InstallOpenKruiseGame(gsm.cfg)
			} else {
				continue
			}
		}
	}
	fmt.Printf("OKG is already installed!\n")
}

func NewGameServerManager(config *restclient.Config, sldId string) *GameServerManager {
	kruisegameClient := kruisegameclientset.NewForConfigOrDie(config)
	kubeClient := kubernetes.NewForConfigOrDie(config)
	kruseClient := kruiseclientset.NewForConfigOrDie(config)
	return &GameServerManager{
		kruisegameClient: kruisegameClient,
		kruseClient:      kruseClient,
		client:           kubeClient,
		cfg:              config,
		slbId:            sldId,
	}
}

func defaultGameServerSet(slbId string) *gamekruisev1alpha1.GameServerSet {
	defaultResourceStandard := DefaultResourceStandards()
	maxUnavailable := intstr.FromString("100%")
	return &gamekruisev1alpha1.GameServerSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      GSName,
			Namespace: "default",
		},
		Spec: gamekruisev1alpha1.GameServerSetSpec{
			UpdateStrategy: gamekruisev1alpha1.UpdateStrategy{
				RollingUpdate: &gamekruisev1alpha1.RollingUpdateStatefulSetStrategy{
					MaxUnavailable:  &maxUnavailable,
					PodUpdatePolicy: kruiseV1beta1.InPlaceIfPossiblePodUpdateStrategyType,
				},
			},
			Replicas: pointer.Int32(0),
			Network: &gamekruisev1alpha1.Network{
				NetworkType: alibabacloud.SlbNetwork,
				NetworkConf: []gamekruisev1alpha1.NetworkConfParams{
					{
						Name:  alibabacloud.SlbIdsConfigName,
						Value: slbId,
					},
					{
						Name:  alibabacloud.PortProtocolsConfigName,
						Value: "8211/UDP",
					},
				},
			},
			GameServerTemplate: gamekruisev1alpha1.GameServerTemplate{
				PodTemplateSpec: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							PlayersAnnotation: DefaultPlayers,
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  GSName,
								Image: GSImage,
								Resources: corev1.ResourceRequirements{
									Requests: map[corev1.ResourceName]resource.Quantity{
										corev1.ResourceMemory: resource.MustParse(defaultResourceStandard[SmallType].RequestMem),
										corev1.ResourceCPU:    resource.MustParse(defaultResourceStandard[SmallType].RequestCpu),
									},
									Limits: map[corev1.ResourceName]resource.Quantity{
										corev1.ResourceMemory: resource.MustParse(defaultResourceStandard[SmallType].LimitMem),
										corev1.ResourceCPU:    resource.MustParse(defaultResourceStandard[SmallType].LimitCpu),
									},
								},
								Env: []corev1.EnvVar{
									{
										Name: "PLAYERS",
										ValueFrom: &corev1.EnvVarSource{
											FieldRef: &corev1.ObjectFieldSelector{
												FieldPath: "metadata.annotations['players']",
											},
										},
									},
								},
							},
						},
					},
				},
				ReclaimPolicy: gamekruisev1alpha1.DeleteGameServerReclaimPolicy,
			},
		},
	}
}

func convContainerRecreateRequest(gsName string) *kruiseV1alpha1.ContainerRecreateRequest {
	hour, min, sec := time.Now().Clock()
	return &kruiseV1alpha1.ContainerRecreateRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: gsName + "-" + strconv.Itoa(hour) + strconv.Itoa(min) + strconv.Itoa(sec),
		},
		Spec: kruiseV1alpha1.ContainerRecreateRequestSpec{
			PodName: gsName,
			Containers: []kruiseV1alpha1.ContainerRecreateRequestContainer{
				{
					Name: GSName,
				},
			},
		},
	}
}

type SortGs []gamekruisev1alpha1.GameServer

func (sg SortGs) Len() int {
	return len(sg)
}

func (sg SortGs) Swap(i, j int) {
	sg[i], sg[j] = sg[j], sg[i]
}

func (sg SortGs) Less(i, j int) bool {
	return sg[i].CreationTimestamp.Time.Before(sg[j].CreationTimestamp.Time)
}
