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

package openstack

import (
	"fmt"

	"github.com/nyl1001/pkg/errors"
	"github.com/nyl1001/pkg/util/regutils"

	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
)

type SInstanceNic struct {
	MacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
	Version int    `json:"version"`
	Addr    string `json:"addr"`
	Type    string `json:"OS-EXT-IPS:type"`
}

type SFixedIp struct {
	IpAddress string
	SubnetId  string
}

type SInstancePort struct {
	region    *SRegion
	FixedIps  []SFixedIp
	MacAddr   string
	NetId     string
	PortId    string
	PortState string

	cloudprovider.DummyICloudNic
}

func (region *SRegion) GetInstancePorts(instanceId string) ([]SInstancePort, error) {
	resource := fmt.Sprintf("/servers/%s/os-interface", instanceId)
	resp, err := region.ecsList(resource, nil)
	if err != nil {
		return nil, errors.Wrap(err, "ecsList")
	}
	ports := []SInstancePort{}
	err = resp.Unmarshal(&ports, "interfaceAttachments")
	if err != nil {
		return nil, errors.Wrap(err, "resp.Unmarshal")
	}
	return ports, nil
}

func (nic *SInstancePort) GetId() string {
	return ""
}

func (nic *SInstancePort) GetIP() string {
	for i := range nic.FixedIps {
		if regutils.MatchIPAddr(nic.FixedIps[i].IpAddress) {
			return nic.FixedIps[i].IpAddress
		}
	}
	return ""
}

func (nic *SInstancePort) GetMAC() string {
	return nic.MacAddr
}

func (nic *SInstancePort) GetDriver() string {
	return "virtio"
}

func (nic *SInstancePort) InClassicNetwork() bool {
	return false
}

func (nic *SInstancePort) GetINetworkId() string {
	for i := range nic.FixedIps {
		if regutils.MatchIPAddr(nic.FixedIps[i].IpAddress) {
			network, err := nic.region.GetNetwork(nic.FixedIps[i].SubnetId)
			if err != nil {
				continue
			}
			pools := network.AllocationPools
			for j := range pools {
				if !pools[j].IsValid() {
					continue
				}
				network.AllocationPools = []AllocationPool{pools[j]}
				if network.Contains(nic.FixedIps[i].IpAddress) {
					return network.GetGlobalId()
				}
			}
		}
	}
	return ""
}
