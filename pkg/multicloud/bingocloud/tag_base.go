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

package bingocloud

import (
	"github.com/nyl1001/pkg/errors"

	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
)

type BingoTags struct {
	TagSet []struct {
		Key   string
		Value string
	}
}

func (self *BingoTags) GetTags() (map[string]string, error) {
	tags := map[string]string{}
	for _, tag := range self.TagSet {
		tags[tag.Key] = tag.Value
	}
	return tags, nil
}

func (self *BingoTags) GetSysTags() map[string]string {
	return nil
}

func (self *BingoTags) SetTags(tags map[string]string, replace bool) error {
	return errors.Wrap(cloudprovider.ErrNotImplemented, "SetTags")
}
