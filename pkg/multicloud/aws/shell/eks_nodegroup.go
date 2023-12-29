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

package shell

import (
	"github.com/nyl1001/pkg/util/shellutils"

	"github.com/nyl1001/cloudmux/pkg/multicloud/aws"
)

func init() {
	type NodegroupListOptions struct {
		CLUSTER_NAME string
		NextToken    string
	}
	shellutils.R(&NodegroupListOptions{}, "node-group-list", "List node group", func(cli *aws.SRegion, args *NodegroupListOptions) error {
		ret, _, err := cli.GetNodegroups(args.CLUSTER_NAME, args.NextToken)
		if err != nil {
			return err
		}
		printList(ret, 0, 0, 0, []string{})
		return nil
	})

	type NodegroupNameOptions struct {
		CLUSTER_NAME string
		NAME         string
	}

	shellutils.R(&NodegroupNameOptions{}, "node-group-show", "Show node group", func(cli *aws.SRegion, args *NodegroupNameOptions) error {
		ret, err := cli.GetNodegroup(args.CLUSTER_NAME, args.NAME)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&NodegroupNameOptions{}, "node-group-delete", "Delete node group", func(cli *aws.SRegion, args *NodegroupNameOptions) error {
		return cli.DeleteNodegroup(args.CLUSTER_NAME, args.NAME)
	})

}
