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
	"github.com/nyl1001/cloudmux/pkg/multicloud/hcso/client/manager"
)

type SElbL7policiesManager struct {
	SResourceManager
}

func NewElbL7policiesManager(cfg manager.IManagerConfig) *SElbL7policiesManager {
	var requestHook portProject
	if len(cfg.GetProjectId()) > 0 {
		requestHook = portProject{projectId: cfg.GetProjectId()}
	}

	return &SElbL7policiesManager{SResourceManager: SResourceManager{
		SBaseManager:  NewBaseManager2(cfg, &requestHook),
		ServiceName:   ServiceNameELB,
		Region:        cfg.GetRegionId(),
		ProjectId:     "",
		version:       "v2.0",
		Keyword:       "l7policy",
		KeywordPlural: "l7policies",

		ResourceKeyword: "lbaas/l7policies",
	}}
}
