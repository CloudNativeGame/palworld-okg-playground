package alibabacloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CloudNativeGame/palworld-okg-playground/cloudprovider"
	cs "github.com/alibabacloud-go/cs-20151215/v4/client"
	restclient "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"strings"
)

// const DefaultCreateClusterConfig = `{"name":"test-create","cluster_type":"ManagedKubernetes","disable_rollback":true,"timeout_mins":60,"kubernetes_version":"1.28.3-aliyun.1","region_id":"cn-hangzhou","snat_entry":true,"cloud_monitor_flags":true,"endpoint_public_access":true,"deletion_protection":true,"proxy_mode":"ipvs","cis_enable_risk_check":true,"tags":[],"timezone":"Asia/Shanghai","addons":[{"name":"security-inspector"},{"name":"terway-eniip","config":"{\"IPVlan\":\"false\",\"NetworkPolicy\":\"false\",\"ENITrunking\":\"false\"}"},{"name":"csi-plugin"},{"name":"csi-provisioner"},{"name":"storage-operator","config":"{\"CnfsOssEnable\":\"false\",\"CnfsNasEnable\":\"true\"}"},{"name":"logtail-ds","config":"{\"IngressDashboardEnabled\":\"true\"}"},{"name":"ack-node-problem-detector","config":"{\"sls_project_name\":\"\"}"},{"name":"nginx-ingress-controller","config":"{\"IngressSlbNetworkType\":\"internet\",\"IngressSlbSpec\":\"slb.s2.small\"}"},{"name":"ack-node-local-dns"},{"name":"arms-prometheus"},{"name":"alicloud-monitor-controller","config":"{\"group_contact_ids\":\"[41399]\"}"}],"cluster_spec":"ack.pro.small","os_type":"Linux","platform":"AliyunLinux","image_type":"AliyunLinux3","pod_vswitch_ids":["vsw-bp1njm4435ly60z7adj0j"],"runtime":{"name":"containerd","version":"1.6.20"},"charge_type":"PostPaid","vpcid":"vpc-bp1qoxk37f4wgf59e2avn","service_cidr":"172.16.0.0/16","vswitch_ids":["vsw-bp1njm4435ly60z7adj0j"],"ip_stack":"ipv4","key_pair":"packer_64899245-d623-9e4b-42ee-33ee5a81f84f","logging_type":"SLS","cpu_policy":"none","service_account_issuer":"https://kubernetes.default.svc","api_audiences":"https://kubernetes.default.svc","is_enterprise_security_group":true,"controlplane_log_ttl":"30","controlplane_log_components":["apiserver","kcm","scheduler","ccm","controlplane-events","alb"],"nodepools":[{"nodepool_info":{"name":"default-nodepool"},"scaling_group":{"vswitch_ids":["vsw-bp1njm4435ly60z7adj0j"],"system_disk_category":"cloud_essd","system_disk_size":120,"system_disk_performance_level":"PL0","system_disk_encrypted":false,"data_disks":[],"instance_types":["ecs.c6.xlarge"],"tags":[],"instance_charge_type":"PostPaid","soc_enabled":false,"cis_enabled":false,"key_pair":"packer_64899245-d623-9e4b-42ee-33ee5a81f84f","login_as_non_root":false,"security_group_ids":[],"platform":"AliyunLinux","image_id":"aliyun_3_x64_20G_alibase_20230727.vhd","image_type":"AliyunLinux3","desired_size":3,"rds_instances":[],"multi_az_policy":"BALANCE"},"kubernetes_config":{"cpu_policy":"none","cms_enabled":true,"unschedulable":false,"runtime":"containerd","runtime_version":"1.6.20"}}],"num_of_nodes":0}`
const DefaultCreateClusterConfig = `{"addons":[{"name":"ack-goatscaler"},{"name":"terway-eniip","config":"{\"IPVlan\":\"false\",\"NetworkPolicy\":\"false\",\"ENITrunking\":\"false\"}"}],"api_audiences":"https://kubernetes.default.svc","charge_type":"PostPaid","cloud_monitor_flags":true,"cluster_spec":"ack.pro.small","cluster_type":"ManagedKubernetes","controlplane_log_components":["apiserver","kcm","scheduler","ccm","controlplane-events","alb"],"controlplane_log_ttl":"30","cpu_policy":"none","deletion_protection":true,"disable_rollback":true,"endpoint_public_access":true,"image_type":"AliyunLinux3","ip_stack":"ipv4","is_enterprise_security_group":true,"kubernetes_version":"1.28.3-aliyun.1","logging_type":"SLS","name":"test-create0126","nodepools":[{"kubernetes_config":{"cms_enabled":true,"cpu_policy":"none","runtime":"containerd","runtime_version":"1.6.20"},"nodepool_info":{"name":"default-nodepool"},"scaling_group":{"image_id":"aliyun_3_x64_20G_alibase_20230727.vhd","image_type":"AliyunLinux3","instance_charge_type":"PostPaid","instance_types":["ecs.c6.xlarge"],"multi_az_policy":"BALANCE","platform":"AliyunLinux","system_disk_category":"cloud_essd","system_disk_performance_level":"PL0","system_disk_size":120,"vswitch_ids":["vsw-bp1njm4435ly60z7adj0j"]},"auto_scaling":{"enable":true,"max_instances":10,"min_instances":3}}],"num_of_nodes":0,"os_type":"Linux","platform":"AliyunLinux","pod_vswitch_ids":["vsw-bp1njm4435ly60z7adj0j"],"proxy_mode":"ipvs","region_id":"cn-hangzhou","runtime":{"name":"containerd","version":"1.6.20"},"service_account_issuer":"https://kubernetes.default.svc","service_cidr":"172.16.0.0/16","snat_entry":true,"timeout_mins":60,"timezone":"Asia/Shanghai","vpcid":"vpc-bp1qoxk37f4wgf59e2avn","service_cidr":"172.16.0.0/21","vswitch_ids":["vsw-bp1njm4435ly60z7adj0j"]}`

// "github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

type AlibabaCloudManager struct {
	cfg *CloudConfig
	openAPIService
}

func CreateAlibabaCloudManager(cfg *CloudConfig) (manager *AlibabaCloudManager, err error) {
	// TODO add cloud provider token support
	cfg.ApplyDefaults()

	if cfg.IsValid() == false {
		return nil, errors.New("please check whether you have provided correct AccessKeyId,AccessKeySecret,RegionId or sts token")
	}

	manager = &AlibabaCloudManager{
		cfg: cfg,
		//nodePools: make(map[string]nodepool),
	}

	openApiSvc, err := NewOpenAPIService(cfg)
	manager.openAPIService = openApiSvc

	//err = manager.updateClusterNodePools()
	//
	//if err != nil {
	//	log.Fatalf("update cluster node pools failed, err: %s and exit", err.Error())
	//}

	return manager, nil
}

type NodePool struct {
	Name string
}
type ClusterConfig struct {
	RegionId          string `json:"region_id"`
	Name              string `json:"name"`
	ClusterType       string `json:"cluster_type"`
	KubernetesVersion string `json:"kubernetes_version"`
	ChargeType        string

	// net
	SnatEntry            bool
	EndpointPublicAccess bool
	ProxyMode            string
	VpcId                string `json:"vpc_id"`
	VswitchIds           string `json:"vswitch_ids"`
	PodVswitchIds        []string
	ServiceCidr          string
	IpStack              string

	// nodepool
	NodePools []NodePool
}

func (c *ClusterConfig) Options() map[string]string {
	optionStr, _ := json.Marshal(*c)
	res := map[string]string{}
	json.Unmarshal(optionStr, &res)
	return res
}

func defaultClusterConfig() ClusterConfig {
	return ClusterConfig{
		Name:              "palworld-cluster",
		ClusterType:       "ManagedKubernetes",
		KubernetesVersion: "1.28.3-aliyun.1",
		ChargeType:        "PostPaid",

		SnatEntry:            true,
		EndpointPublicAccess: true,
		ProxyMode:            "ipvs",
		VpcId:                "",
		VswitchIds:           "",
		ServiceCidr:          "172.19.0.0/20",
	}
}

func (am *AlibabaCloudManager) initCreateClusterConfig(opts cloudprovider.ClusterOptions) *cs.CreateClusterRequest {
	def2 := &cs.CreateClusterRequest{}
	err := json.Unmarshal([]byte(DefaultCreateClusterConfig), def2)
	if err != nil {
		return nil
	}
	if opts != nil {
		options := opts.Options()
		if v, ok := options["vpc_id"]; ok && len(v) > 0 {
			def2.Vpcid = &v
		}
		if v, ok := options["cluster_type"]; ok {
			def2.ClusterType = &v
		}
		if v, ok := options["name"]; ok {
			def2.Name = &v
		}
		if v, ok := options["region_id"]; ok {
			def2.RegionId = &v
		}
		if v, ok := options["vswitch_ids"]; ok && len(v) > 0 {
			ids := []*string{}
			for _, id := range strings.Split(v, ",") {
				id = strings.Trim(id, " ")
				if id != "" {
					ids = append(ids, &id)
				}
			}
			def2.VswitchIds = ids
		}
	}

	return def2
	//return &cs.CreateClusterRequest{
	//	RegionId:          &region,
	//	Name:              &defaultConf.Name,
	//	ClusterType:       &defaultConf.ClusterType,
	//	KubernetesVersion: &defaultConf.KubernetesVersion,
	//	ChargeType:        &defaultConf.ChargeType,
	//
	//	SnatEntry:            &defaultConf.SnatEntry,
	//	EndpointPublicAccess: &defaultConf.EndpointPublicAccess,
	//	ProxyMode:            &defaultConf.ProxyMode,
	//}
}

func (am *AlibabaCloudManager) CreateCluster(options cloudprovider.ClusterOptions) (cloudprovider.KubernetesCluster, error) {
	req := am.initCreateClusterConfig(options)
	resp, err := am.openAPIService.CreateCluster(req)
	if err != nil {
		klog.Errorf("failed to create cluster, err %v", err)
		return nil, err
	}
	if resp == nil || resp.Body == nil || resp.Body.ClusterId == nil {
		return nil, fmt.Errorf("failed to creat cluster due to response is empty")
	}

	switch *(req.ClusterType) {
	case "serverless":
		return &ServerlessCluster{
			Id:     *(resp.Body.ClusterId),
			Status: "creating",
		}, nil
	default:
		return &ManagedCluster{
			Id:     *(resp.Body.ClusterId),
			Status: "creating",
		}, nil

	}

	return nil, nil
}

func (am *AlibabaCloudManager) ListClusters() ([]cloudprovider.KubernetesCluster, error) {
	return []cloudprovider.KubernetesCluster{}, nil
}

func (am *AlibabaCloudManager) DeleteCluster(clusterId string) error {
	return nil
}

func (am *AlibabaCloudManager) GetCluster(clusterId string) (cloudprovider.KubernetesCluster, error) {
	return nil, nil
}

func (am *AlibabaCloudManager) GetKubernetesConfig(clusterId string) (*restclient.Config, error) {
	return nil, nil
}
