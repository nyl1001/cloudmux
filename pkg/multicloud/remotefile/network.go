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

package remotefile

import (
	"github.com/nyl1001/pkg/util/rbacscope"

	api "github.com/nyl1001/cloudmux/pkg/apis/compute"
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
)

type SNetwork struct {
	SResourceBase

	wire    *SWire
	WireId  string
	IpStart string
	IpEnd   string
	IpMask  int8
	Gatway  string
}

func (self *SNetwork) GetIWire() cloudprovider.ICloudWire {
	return self.wire
}

func (self *SNetwork) GetIpStart() string {
	return self.IpStart
}

func (self *SNetwork) GetIpEnd() string {
	return self.IpEnd
}

func (self *SNetwork) GetIpMask() int8 {
	return self.IpMask
}

func (self *SNetwork) GetGateway() string {
	return self.Gatway
}

func (self *SNetwork) GetServerType() string {
	return api.NETWORK_TYPE_GUEST
}

func (self *SNetwork) GetPublicScope() rbacscope.TRbacScope {
	return rbacscope.ScopeDomain
}

func (self *SNetwork) Delete() error {
	return cloudprovider.ErrNotSupported
}

func (self *SNetwork) GetAllocTimeoutSeconds() int {
	return 6
}
