package alibabacloud

import (
	"fmt"
	cs "github.com/alibabacloud-go/cs-20151215/v4/client"
	apiconf "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"k8s.io/klog/v2"
	"time"
)

const (
	refreshClientInterval = 5 * time.Minute
	clientScheme          = "HTTPS"
)

type openAPIService interface {
	vpcService
	ackService
	slbService
}

type openAPIServiceImpl struct {
	vpcService
	ackService
	slbService
	cfg *CloudConfig
}

type vpcService interface {
	CreateVpc(request *vpc.CreateVpcRequest) (response *vpc.CreateVpcResponse, err error)
	DescribeZones(request *vpc.DescribeZonesRequest) (response *vpc.DescribeZonesResponse, err error)
	CreateVSwitch(request *vpc.CreateVSwitchRequest) (response *vpc.CreateVSwitchResponse, err error)
	DescribeVSwitches(req *vpc.DescribeVSwitchesRequest) (*vpc.DescribeVSwitchesResponse, error)
}

type ackService interface {
	CreateCluster(request *cs.CreateClusterRequest) (*cs.CreateClusterResponse, error)
	DescribeClusterNodes(ClusterId *string, request *cs.DescribeClusterNodesRequest) (*cs.DescribeClusterNodesResponse, error)
	RemoveNodePoolNodes(ClusterId *string, NodepoolId *string, request *cs.RemoveNodePoolNodesRequest) (*cs.RemoveNodePoolNodesResponse, error)
	DescribeClusterNodePools(ClusterId *string) (response *cs.DescribeClusterNodePoolsResponse, err error)
	DescribeClusterNodePoolDetail(ClusterId *string, NodepoolId *string) (*cs.DescribeClusterNodePoolDetailResponse, error)
	ScaleClusterNodePool(ClusterId *string, NodepoolId *string, request *cs.ScaleClusterNodePoolRequest) (*cs.ScaleClusterNodePoolResponse, error)
	ModifyClusterNodePool(ClusterId *string, NodepoolId *string, request *cs.ModifyClusterNodePoolRequest) (*cs.ModifyClusterNodePoolResponse, error)
	DescribeTaskInfo(taskId *string) (*cs.DescribeTaskInfoResponse, error)
	DescribeClusterUserKubeconfig(ClusterId *string, request *cs.DescribeClusterUserKubeconfigRequest) (*cs.DescribeClusterUserKubeconfigResponse, error)
	DescribeClusterDetail(ClusterId *string) (*cs.DescribeClusterDetailResponse, error)
}

type slbService interface {
	CreateLoadBalancer(*slb.CreateLoadBalancerRequest) (*slb.CreateLoadBalancerResponse, error)
}

type slbServiceImpl struct {
	*slb.Client
}

type vpcServiceImpl struct {
	*vpc.Client
}

type ackServiceImpl struct {
	*cs.Client
}

func (openAPI *openAPIServiceImpl) Run() {
	timer := time.NewTicker(refreshClientInterval)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			err := openAPI.refreshClient()
			if err != nil {
				// Must Fatal because the client is not valid.
				klog.Fatalf("refresh client error: %v", err)
			}
		}
	}
}

func (openAPI *openAPIServiceImpl) refreshClient() (err error) {
	cfg := openAPI.cfg

	// add endpoint mapping to vpc endpoint
	region := cfg.GetRegion()
	//internal := cfg.Internal
	//enableSts := cfg.STSEnabled

	// add endpoint mapping to vpc endpoint
	//if region != "" && internal {
	//	endpoints.AddEndpointMapping(region, "Ess", fmt.Sprintf("ess-vpc.%s.aliyuncs.com", region))
	//	endpoints.AddEndpointMapping(region, "Ecs", fmt.Sprintf("ecs-vpc.%s.aliyuncs.com", region))
	//	endpoints.AddEndpointMapping(region, "Vpc", fmt.Sprintf("vpc-vpc.%s.aliyuncs.com", region))
	//	endpoints.AddEndpointMapping(region, "Ack", fmt.Sprintf("cs.%s.aliyuncs.com", region))
	//}

	accessKeyId := cfg.AccessKeyID
	accessKeySecret := cfg.AccessKeySecret
	config := &apiconf.Config{
		AccessKeyId:     &accessKeyId,
		AccessKeySecret: &accessKeySecret,
		RegionId:        &region,
	}
	//essClient, err := ess.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	//if err != nil {
	//	klog.Errorf("failed to create ess client,Because of %s", err.Error())
	//	return err
	//}
	//openAPI.essService = &essServiceImpl{essClient}

	//ecsClient, err := ecs.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	//if err != nil {
	//	klog.Errorf("failed to create ecs client,Because of %s", err.Error())
	//	return err
	//}
	//openAPI.ecsService = &ecsServiceImpl{ecsClient}

	slbClient, err := slb.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	if err != nil {
		klog.Errorf("failed to create slb client,Because of %s", err.Error())
		return err
	}
	openAPI.slbService = &slbServiceImpl{slbClient}

	vpcClient, err := vpc.NewClientWithAccessKey(region, accessKeyId, accessKeySecret)
	if err != nil {
		klog.Errorf("failed to create vpc client,Because of %s", err.Error())
		return err
	}
	openAPI.vpcService = &vpcServiceImpl{vpcClient}

	csClient, err := cs.NewClient(config)
	if err != nil {
		klog.Errorf("failed to create ack client,Because of %s", err.Error())
		return err
	}
	openAPI.ackService = &ackServiceImpl{csClient}

	return nil
}

func NewOpenAPIService(cfg *CloudConfig) (openAPIService, error) {
	if cfg.IsValid() == false {
		//Never reach here.
		return nil, fmt.Errorf("your cloud config is not valid")
	}

	openAPI := &openAPIServiceImpl{
		cfg: cfg,
	}

	err := openAPI.refreshClient()
	if err != nil {
		klog.Errorf("failed to refresh client,Because of %s", err.Error())
	}

	// TODO not a good way to run a goroutine here.
	go openAPI.Run()

	return openAPI, err
}
