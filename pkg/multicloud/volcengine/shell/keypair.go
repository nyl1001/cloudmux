// Copyright 2023 Yunion
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
	"github.com/nyl1001/cloudmux/pkg/multicloud/volcengine"
	"github.com/nyl1001/pkg/util/shellutils"
)

func init() {
	type KeypairListOptions struct {
		MaxResult int
		NextToken string
		Finger    string
		Name      string
	}
	shellutils.R(&KeypairListOptions{}, "keypair-list", "List keypairs", func(cli *volcengine.SRegion, args *KeypairListOptions) error {
		keypairs, _, err := cli.GetKeypairs(args.Finger, args.Name, args.MaxResult, args.NextToken)
		if err != nil {
			return err
		}
		printList(keypairs, 0, 0, 0, nil)
		return nil
	})

}
