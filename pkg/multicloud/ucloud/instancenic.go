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

package ucloud

import (
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
)

type SInstanceNic struct {
	instance *SInstance
	ipAddr   string
	macAddr  string

	cloudprovider.DummyICloudNic
}

func (self *SInstanceNic) GetId() string {
	return ""
}

func (self *SInstanceNic) GetIP() string {
	return self.ipAddr
}

func (self *SInstanceNic) GetMAC() string {
	return self.macAddr
}

func (self *SInstanceNic) GetDriver() string {
	return "virtio"
}

func (self *SInstanceNic) InClassicNetwork() bool {
	return false
}

func (self *SInstanceNic) GetINetworkId() string {
	for _, ip := range self.instance.IPSet {
		if ip.IP == self.ipAddr {
			return ip.SubnetID
		}
	}
	return ""
}
