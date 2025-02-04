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
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
	"github.com/nyl1001/pkg/errors"
)

type SCDNDomainBase struct {
	SResourceBase
}

func (self *SCDNDomainBase) GetCacheKeys() (*cloudprovider.SCDNCacheKeys, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetCacheKeys")
}

func (self *SCDNDomainBase) GetRangeOriginPull() (*cloudprovider.SCDNRangeOriginPull, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetRangeOriginPull")
}

func (self *SCDNDomainBase) GetCache() (*cloudprovider.SCDNCache, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetCache")
}

func (self *SCDNDomainBase) GetHTTPS() (*cloudprovider.SCDNHttps, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetHTTPS")
}

func (self *SCDNDomainBase) GetForceRedirect() (*cloudprovider.SCDNForceRedirect, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetForceRedirect")
}

func (self *SCDNDomainBase) GetReferer() (*cloudprovider.SCDNReferer, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetReferer")
}

func (self *SCDNDomainBase) GetMaxAge() (*cloudprovider.SCDNMaxAge, error) {
	return nil, errors.Wrapf(cloudprovider.ErrNotImplemented, "GetMaxAge")
}
