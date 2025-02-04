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

package apsara

import (
	"fmt"
	"time"

	"github.com/nyl1001/pkg/errors"
	"github.com/nyl1001/pkg/jsonutils"

	api "github.com/nyl1001/cloudmux/pkg/apis/compute"
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
	"github.com/nyl1001/cloudmux/pkg/multicloud"
)

func (self *SRegion) otsRequest(apiName string, params map[string]string) (jsonutils.JSONObject, error) {
	client, err := self.getSdkClient()
	if err != nil {
		return nil, err
	}
	domain := self.client.getDomain(APSARA_PRODUCT_OTS)
	return self.productRequest(client, APSARA_PRODUCT_OTS, domain, APSARA_OTS_API_VERSION, apiName, params, self.client.debug)
}

type STablestore struct {
	multicloud.SResourceBase
	ApsaraTags
	region *SRegion

	InstanceName string
	Timestamp    time.Time

	DepartmentInfo
}

func (self *STablestore) GetGlobalId() string {
	return self.InstanceName
}

func (self *STablestore) GetName() string {
	return self.InstanceName
}

func (self *STablestore) GetId() string {
	return self.InstanceName
}

func (self *STablestore) GetStatus() string {
	return api.TABLESTORE_STATUS_RUNNING
}

func (self *SRegion) GetTablestoreInstances(pageSize, pageNumber int) ([]STablestore, int, error) {
	params := map[string]string{
		"RegionId":   self.RegionId,
		"PageSize":   fmt.Sprintf("%d", pageSize),
		"PageNumber": fmt.Sprintf("%d", pageNumber),
	}
	resp, err := self.otsRequest("ListInstance", params)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "ListInstance")
	}
	total, _ := resp.Int("TotalCount")
	ret := []STablestore{}
	return ret, int(total), resp.Unmarshal(&ret, "InstanceInfos", "InstanceInfo")
}

func (self *SRegion) GetICloudTablestores() ([]cloudprovider.ICloudTablestore, error) {
	ots, pageNumber := []STablestore{}, 1
	for {
		part, total, err := self.GetTablestoreInstances(50, pageNumber)
		if err != nil {
			return nil, err
		}
		ots = append(ots, part...)
		if len(ots) >= total {
			break
		}
		pageNumber++
	}
	ret := []cloudprovider.ICloudTablestore{}
	for i := range ots {
		ots[i].region = self
		ret = append(ret, &ots[i])
	}
	return ret, nil
}
