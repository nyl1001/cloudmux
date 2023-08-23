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
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/esxi"
)

func init() {
	type NetworkListOptions struct {
		// DATACENTER string `help:"List datastores in datacenter"`
	}
	shellutils.R(&NetworkListOptions{}, "network-list", "List networks in datacenter", func(cli *esxi.SESXiClient, args *NetworkListOptions) error {
		nets, err := cli.GetNetworks()
		if err != nil {
			return err
		}
		printList(nets, nil)
		return nil
	})

	shellutils.R(&NetworkListOptions{}, "wire-list", "List wires in datacenter", func(cli *esxi.SESXiClient, args *NetworkListOptions) error {
		vpcs, err := cli.GetIVpcs()
		if err != nil {
			return err
		}
		wires, err := vpcs[0].GetIWires()
		if err != nil {
			return err
		}
		printList(wires, nil)
		return nil
	})
}
