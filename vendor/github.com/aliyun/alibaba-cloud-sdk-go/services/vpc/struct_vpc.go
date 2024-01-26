package vpc

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// Vpc is a nested struct in vpc response
type Vpc struct {
	CreationTime           string                            `json:"CreationTime" xml:"CreationTime"`
	Status                 string                            `json:"Status" xml:"Status"`
	VpcId                  string                            `json:"VpcId" xml:"VpcId"`
	IsDefault              bool                              `json:"IsDefault" xml:"IsDefault"`
	AdvancedResource       bool                              `json:"AdvancedResource" xml:"AdvancedResource"`
	OwnerId                int64                             `json:"OwnerId" xml:"OwnerId"`
	RegionId               string                            `json:"RegionId" xml:"RegionId"`
	VpcName                string                            `json:"VpcName" xml:"VpcName"`
	VRouterId              string                            `json:"VRouterId" xml:"VRouterId"`
	DhcpOptionsSetStatus   string                            `json:"DhcpOptionsSetStatus" xml:"DhcpOptionsSetStatus"`
	CidrBlock              string                            `json:"CidrBlock" xml:"CidrBlock"`
	Description            string                            `json:"Description" xml:"Description"`
	NetworkAclNum          string                            `json:"NetworkAclNum" xml:"NetworkAclNum"`
	SupportAdvancedFeature bool                              `json:"SupportAdvancedFeature" xml:"SupportAdvancedFeature"`
	ResourceGroupId        string                            `json:"ResourceGroupId" xml:"ResourceGroupId"`
	DhcpOptionsSetId       string                            `json:"DhcpOptionsSetId" xml:"DhcpOptionsSetId"`
	Ipv6CidrBlock          string                            `json:"Ipv6CidrBlock" xml:"Ipv6CidrBlock"`
	CenStatus              string                            `json:"CenStatus" xml:"CenStatus"`
	VSwitchIds             VSwitchIdsInDescribeVpcs          `json:"VSwitchIds" xml:"VSwitchIds"`
	SecondaryCidrBlocks    SecondaryCidrBlocksInDescribeVpcs `json:"SecondaryCidrBlocks" xml:"SecondaryCidrBlocks"`
	UserCidrs              UserCidrsInDescribeVpcs           `json:"UserCidrs" xml:"UserCidrs"`
	NatGatewayIds          NatGatewayIds                     `json:"NatGatewayIds" xml:"NatGatewayIds"`
	RouterTableIds         RouterTableIds                    `json:"RouterTableIds" xml:"RouterTableIds"`
	Tags                   TagsInDescribeVpcs                `json:"Tags" xml:"Tags"`
	Ipv6CidrBlocks         Ipv6CidrBlocksInDescribeVpcs      `json:"Ipv6CidrBlocks" xml:"Ipv6CidrBlocks"`
}