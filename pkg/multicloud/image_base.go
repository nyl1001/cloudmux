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

package multicloud

import (
	"github.com/nyl1001/pkg/util/rbacscope"
	"github.com/pkg/errors"

	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
)

type SImageBase struct {
	SVirtualResourceBase
}

func (self *SImageBase) GetPublicScope() rbacscope.TRbacScope {
	return rbacscope.ScopeSystem
}

func (self *SImageBase) GetSubImages() []cloudprovider.SSubImage {
	return nil
}

func (self *SImageBase) Export(opts *cloudprovider.SImageExportOptions) ([]cloudprovider.SImageExportInfo, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "Export")
}
