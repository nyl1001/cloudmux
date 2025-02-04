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

import (
	"fmt"
	"time"

	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
)

type SMachineType struct {
	Id                           string
	CreationTimestamp            time.Time
	Name                         string
	Description                  string
	GuestCpus                    int
	MemoryMb                     int
	ImageSpaceGb                 int
	MaximumPersistentDisks       int
	MaximumPersistentDisksSizeGb int
	Zone                         string
	SelfLink                     string
	IsSharedCpu                  bool
	Kind                         string
}

func (region *SRegion) GetMachineTypes(zone string, maxResults int, pageToken string) ([]SMachineType, error) {
	machines := []SMachineType{}
	params := map[string]string{}
	if len(zone) == 0 {
		return nil, cloudprovider.ErrNotFound
	}
	resource := fmt.Sprintf("zones/%s/machineTypes", zone)
	return machines, region.List(resource, params, maxResults, pageToken, &machines)
}

func (region *SRegion) GetMachineType(id string) (*SMachineType, error) {
	machine := &SMachineType{}
	err := region.client.ecsGet("machineTypes", id, machine)
	if err != nil {
		return nil, err
	}
	return machine, nil
}
