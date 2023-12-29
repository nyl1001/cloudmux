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

	"github.com/nyl1001/cloudmux/pkg/multicloud/zstack"
)

func init() {
	type QuotaListOptions struct {
	}
	shellutils.R(&QuotaListOptions{}, "quota-list", "List quotas", func(cli *zstack.SRegion, args *QuotaListOptions) error {
		quotas, err := cli.GetQuotas()
		if err != nil {
			return err
		}
		printList(quotas, 0, 0, 0, nil)
		return nil
	})
}
