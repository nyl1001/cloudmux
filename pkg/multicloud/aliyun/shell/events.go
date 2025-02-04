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
	"fmt"

	"github.com/nyl1001/pkg/errors"
	"github.com/nyl1001/pkg/util/shellutils"
	"github.com/nyl1001/pkg/util/timeutils"

	"github.com/nyl1001/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type EventListOptions struct {
		START     string
		END       string
		Token     string
		EventRW   string `choices:"Read|Write|All"`
		RequestId string
	}
	shellutils.R(&EventListOptions{}, "event-list", "List event", func(cli *aliyun.SRegion, args *EventListOptions) error {
		start, err := timeutils.ParseTimeStr(args.START)
		if err != nil {
			return errors.Wrap(err, "timeutils.ParseTimeStr.Start")
		}
		end, err := timeutils.ParseTimeStr(args.END)
		if err != nil {
			return errors.Wrap(err, "timeutils.ParseTimeStr.End")
		}
		events, token, err := cli.GetEvents(start, end, args.Token, args.EventRW, args.RequestId)
		if err != nil {
			return err
		}
		fmt.Printf("token: %s", token)
		printList(events, 0, 0, 0, []string{})
		return nil
	})

}
