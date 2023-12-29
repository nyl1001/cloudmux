// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package qcloud

import (
	"fmt"

	api "github.com/nyl1001/cloudmux/pkg/apis/compute"
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
	"github.com/nyl1001/cloudmux/pkg/multicloud"
)

type SDTable struct {
	multicloud.SResourceBase
	QcloudTags
	nat *SNatGateway

	CreatedTime      string `json:"CreatedTime"`
	Description      string `json:"Description"`
	IpProtocol       string `json:"IpProtocol"`
	NatGatewayId     string `json:"NatGatewayId"`
	PrivateIpAddress string `json:"PrivateIpAddress"`
	PrivatePort      int    `json:"PrivatePort"`
	PublicIpAddress  string `json:"PublicIpAddress"`
	PublicPort       int    `json:"PublicPort"`
	VpcId            string `json:"VpcId"`
}

func (table *SDTable) GetName() string {
	if len(table.Description) > 0 {
		return table.Description
	}
	return fmt.Sprintf("%s/%s/%d", table.PublicIpAddress, table.IpProtocol, table.PublicPort)
}

func (table *SDTable) GetId() string {
	return fmt.Sprintf("%s/%s/%d", table.NatGatewayId, table.PublicIpAddress, table.PublicPort)
}

func (table *SDTable) GetGlobalId() string {
	return table.GetId()
}

func (table *SDTable) GetStatus() string {
	return api.NAT_STAUTS_AVAILABLE
}

func (table *SDTable) GetExternalIp() string {
	return table.PublicIpAddress
}

func (table *SDTable) GetExternalPort() int {
	return table.PublicPort
}

func (table *SDTable) GetInternalIp() string {
	return table.PrivateIpAddress
}

func (table *SDTable) GetInternalPort() int {
	return table.PrivatePort
}

func (table *SDTable) GetIpProtocol() string {
	return table.IpProtocol
}

func (table *SDTable) Delete() error {
	return cloudprovider.ErrNotImplemented
}

func (region *SRegion) GetDTables(natId string, offset, limit int) ([]SDTable, int, error) {
	params := map[string]string{}
	if limit < 1 || limit > 50 {
		limit = 50
	}
	params["Limit"] = fmt.Sprintf("%d", limit)
	params["Offset"] = fmt.Sprintf("%d", offset)
	params["NatGatewayIds.0"] = natId
	body, err := region.vpcRequest("DescribeNatGatewayDestinationIpPortTranslationNatRules", params)
	if err != nil {
		return nil, 0, err
	}
	tables := []SDTable{}
	err = body.Unmarshal(&tables, "NatGatewayDestinationIpPortTranslationNatRuleSet")
	if err != nil {
		return nil, 0, err
	}
	total, _ := body.Float("TotalCount")
	return tables, int(total), nil
}
