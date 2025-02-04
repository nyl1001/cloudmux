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

package google

import "github.com/nyl1001/pkg/errors"

type SCloudbuildBuild struct {
	Id     string
	Status string
	LogUrl string
}

type SCloudbuildMetadata struct {
	Build SCloudbuildBuild
}

type SCloudbuildOperation struct {
	Name     string
	Metadata SCloudbuildMetadata
}

func (region *SRegion) GetCloudbuildOperation(name string) (*SCloudbuildOperation, error) {
	operation := SCloudbuildOperation{}
	err := region.cloudbuildGet(name, &operation)
	if err != nil {
		return nil, errors.Wrap(err, "region.cloudbuildGet")
	}
	return &operation, nil
}
