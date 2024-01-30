package alibabacloud

import (
	//"github.com/kubernetes/autoscaler/cluster-autoscaler/cloudprovider/alicloud/metadata"
	"k8s.io/klog/v2"
	"os"
)

type CloudConfig struct {
	RegionId        string
	AccessKeyID     string
	AccessKeySecret string
}

const (
	accessKeyId     = "ACCESS_KEY_ID"
	accessKeySecret = "ACCESS_KEY_SECRET"
	regionId        = "REGION_ID"
)

func (cc *CloudConfig) ApplyDefaults() {
	if cc.AccessKeyID == "" {
		cc.AccessKeyID = os.Getenv(accessKeyId)
	}

	if cc.AccessKeySecret == "" {
		cc.AccessKeySecret = os.Getenv(accessKeySecret)
	}

	if cc.RegionId == "" {
		cc.RegionId = os.Getenv(regionId)
	}
	//klog.Infof("get config from env %++v", *cc)
}

func (cc *CloudConfig) IsValid() bool {
	if cc.RegionId == "" || cc.AccessKeyID == "" || cc.AccessKeySecret == "" {
		klog.Errorf("Failed to get AccessKeyId:%s,AccessKeySecret:%s,RegionId:%s from CloudConfig and Env\n", cc.AccessKeyID, cc.AccessKeySecret, cc.RegionId)
		return false
	}

	return true
}

func (cc *CloudConfig) GetRegion() string {
	if cc.RegionId != "" {
		return cc.RegionId
	}
	//m := metadata.NewMetaData(nil)
	//r, err := m.Region()
	//if err != nil {
	//	klog.Errorf("Failed to get RegionId from metadata.Because of %s\n", err.Error())
	//}
	return "cn-hangzhou"
}
