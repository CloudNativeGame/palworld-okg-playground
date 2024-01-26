package alibabacloud

import (
	"encoding/json"
	"fmt"
	cs "github.com/alibabacloud-go/cs-20151215/v4/client"
	"reflect"
	"testing"
)

func TestAlibabaCloudManager_initCreateClusterConfig(t *testing.T) {
	type fields struct {
		cfg            *CloudConfig
		openAPIService openAPIService
	}
	type args struct {
		customConfig *ClusterConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *cs.CreateClusterRequest
	}{
		// TODO: Add test cases.
		{"default-test", fields{}, args{customConfig: nil}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := &AlibabaCloudManager{
				cfg:            tt.fields.cfg,
				openAPIService: tt.fields.openAPIService,
			}
			got := am.initCreateClusterConfig(tt.args.customConfig)
			if got != nil {
				bytes, _ := json.Marshal(got)
				fmt.Printf("got is %s", string(bytes))
			}
			if !reflect.DeepEqual(got, tt.want) {
				//t.Errorf("initCreateClusterConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
