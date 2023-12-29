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

package loader

import (
	"yunion.io/x/log" // on-premise virtualization technologies

	_ "github.com/nyl1001/cloudmux/pkg/multicloud/aliyun/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/apsara/provider" // aliyun apsara stack
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/aws/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/azure/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/baidu/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/bingocloud/provider" // private clouds
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/ctyun/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/cucloud/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/ecloud/provider" // public clouds
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/esxi/provider"   // private clouds
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/google/provider" // public clouds
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/hcso/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/huawei/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/jdcloud/provider" // public clouds
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/ksyun/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/nutanix/provider" // private clouds
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/objectstore/ceph/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/objectstore/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/objectstore/xsky/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/openstack/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/oracle/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/proxmox/provider" // private clouds
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/qcloud/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/qingcloud/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/remotefile/provider" // private clouds
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/ucloud/provider"     // object storages
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/volcengine/provider"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/zstack/provider" // private clouds
)

func init() {
	log.Infof("Loading cloud providers ...")
}
