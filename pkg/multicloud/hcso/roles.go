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

package hcso

import (
	"github.com/nyl1001/pkg/errors"
	"github.com/nyl1001/pkg/jsonutils"

	api "github.com/nyl1001/cloudmux/pkg/apis/cloudid"
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
)

type SRole struct {
	DomainId      string
	Flag          string
	DescriptionCn string
	Catelog       string
	Description   string
	Id            string
	DisplayName   string
	Type          string
	UpdatedTime   string
	CreatedTime   string
	Links         SLink
	Policy        jsonutils.JSONDict
	roleType      string
}

func (role *SRole) GetName() string {
	return role.DisplayName
}

func (role *SRole) GetDescription() string {
	return role.DescriptionCn
}

func (role *SRole) GetPolicyType() string {
	return role.roleType
}

func (role *SRole) GetGlobalId() string {
	return role.DisplayName
}

func (role *SRole) UpdateDocument(document *jsonutils.JSONDict) error {
	return cloudprovider.ErrNotImplemented
}

func (role *SRole) GetDocument() (*jsonutils.JSONDict, error) {
	return &role.Policy, nil
}

func (role *SRole) Delete() error {
	return cloudprovider.ErrNotImplemented
}

func (self *SHuaweiClient) GetISystemCloudpolicies() ([]cloudprovider.ICloudpolicy, error) {
	roles, err := self.GetRoles("", "")
	if err != nil {
		return nil, errors.Wrap(err, "GetRoles")
	}
	ret := []cloudprovider.ICloudpolicy{}
	for i := range roles {
		roles[i].roleType = api.CLOUD_POLICY_TYPE_SYSTEM
		ret = append(ret, &roles[i])
	}
	return ret, nil
}

func (self *SHuaweiClient) GetICustomCloudpolicies() ([]cloudprovider.ICloudpolicy, error) {
	roles, err := self.GetCustomRoles()
	if err != nil {
		return nil, errors.Wrap(err, "GetCustomRoles")
	}
	ret := []cloudprovider.ICloudpolicy{}
	for i := range roles {
		roles[i].roleType = api.CLOUD_POLICY_TYPE_CUSTOM
		ret = append(ret, &roles[i])
	}
	return ret, nil
}

func (self *SHuaweiClient) GetCustomRoles() ([]SRole, error) {
	params := map[string]string{}

	client, err := self.newGeneralAPIClient()
	if err != nil {
		return nil, errors.Wrap(err, "newGeneralAPIClient")
	}

	client.Roles.SetVersion("v3.0/OS-ROLE")
	defer client.Roles.SetVersion("v3.0")

	roles := []SRole{}
	err = doListAllWithNextLink(client.Roles.List, params, &roles)
	if err != nil {
		return nil, errors.Wrap(err, "doListAllWithOffset")
	}
	return roles, nil
}

func (self *SHuaweiClient) CreateICloudpolicy(opts *cloudprovider.SCloudpolicyCreateOptions) (cloudprovider.ICloudpolicy, error) {
	client, err := self.newGeneralAPIClient()
	if err != nil {
		return nil, errors.Wrap(err, "newGeneralAPIClient")
	}

	client.Roles.SetVersion("v3.0/OS-ROLE")
	defer client.Roles.SetVersion("v3.0")

	params := map[string]interface{}{
		"role": map[string]interface{}{
			"display_name": opts.Name,
			"type":         "XA",
			"description":  opts.Desc,
			"policy":       opts.Document,
		},
	}

	resp, err := client.Roles.Create(jsonutils.Marshal(params))
	if err != nil {
		return nil, err
	}
	role := &SRole{roleType: api.CLOUD_POLICY_TYPE_CUSTOM}
	err = resp.Unmarshal(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (self *SHuaweiClient) GetRoles(domainId, name string) ([]SRole, error) {
	params := map[string]string{}
	if len(domainId) > 0 {
		params["domain_id"] = self.ownerId
	}
	if len(name) > 0 {
		params["display_name"] = name
	}

	client, err := self.newGeneralAPIClient()
	if err != nil {
		return nil, errors.Wrap(err, "newGeneralAPIClient")
	}

	roles := []SRole{}
	err = doListAllWithNextLink(client.Roles.List, params, &roles)
	if err != nil {
		return nil, errors.Wrap(err, "doListAllWithOffset")
	}
	return roles, nil
}
