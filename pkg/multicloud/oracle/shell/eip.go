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

	"github.com/nyl1001/cloudmux/pkg/multicloud/oracle"
)

func init() {
	type EipListOptions struct {
		Lifetime string `choices:"RESERVED|EPHEMERAL"`
	}
	shellutils.R(&EipListOptions{}, "eip-list", "list eips", func(cli *oracle.SRegion, args *EipListOptions) error {
		eips, err := cli.GetEips(args.Lifetime)
		if err != nil {
			return err
		}
		printList(eips, 0, 0, 0, []string{})
		return nil
	})

	type EipIdOptions struct {
		ID string
	}

	shellutils.R(&EipIdOptions{}, "eip-show", "Show eip", func(cli *oracle.SRegion, args *EipIdOptions) error {
		eip, err := cli.GetEip(args.ID)
		if err != nil {
			return err
		}
		printObject(eip)
		return nil
	})

}
