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

package modules

import (
	"github.com/nyl1001/pkg/jsonutils"

	"github.com/nyl1001/cloudmux/pkg/multicloud/hcso/client/manager"
	"github.com/nyl1001/cloudmux/pkg/multicloud/hcso/client/responses"
)

type SfsTurboManager struct {
	SResourceManager
}

func NewSfsTurboManager(cfg manager.IManagerConfig) *SfsTurboManager {
	return &SfsTurboManager{SResourceManager: SResourceManager{
		SBaseManager:  NewBaseManager(cfg),
		ServiceName:   ServiceNameSFSTurbo,
		Region:        cfg.GetRegionId(),
		ProjectId:     cfg.GetProjectId(),
		version:       "v1",
		Keyword:       "",
		KeywordPlural: "shares",

		ResourceKeyword: "sfs-turbo/shares",
	}}
}

func (self *SfsTurboManager) List(querys map[string]string) (*responses.ListResult, error) {
	return self.ListInContextWithSpec(nil, "detail", querys, self.KeywordPlural)
}

func (self *SfsTurboManager) Create(params jsonutils.JSONObject) (jsonutils.JSONObject, error) {
	return self.CreateInContextWithSpec(self.ctx, "", params, "")
}
