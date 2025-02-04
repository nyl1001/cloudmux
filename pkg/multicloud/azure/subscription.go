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

package azure

import (
	"net/url"

	api "github.com/nyl1001/cloudmux/pkg/apis/compute"
)

type SSubscription struct {
	SubscriptionId string `json:"subscriptionId"`
	State          string
	DisplayName    string `json:"displayName"`
}

func (self *SSubscription) GetHealthStatus() string {
	if self.State == "Enabled" {
		return api.CLOUD_PROVIDER_HEALTH_NORMAL
	}
	return api.CLOUD_PROVIDER_HEALTH_SUSPENDED
}

func (self *SAzureClient) ListSubscriptions() ([]SSubscription, error) {
	result := []SSubscription{}
	err := self.list("subscriptions", url.Values{}, &result)
	return result, err
}
