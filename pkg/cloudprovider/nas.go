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

package cloudprovider

import "github.com/nyl1001/pkg/util/billing"

type FileSystemCraeteOptions struct {
	Name           string
	Desc           string
	VpcId          string
	NetworkId      string
	Capacity       int64
	StorageType    string
	Protocol       string
	FileSystemType string
	ZoneId         string

	BillingCycle *billing.SBillingCycle
}
