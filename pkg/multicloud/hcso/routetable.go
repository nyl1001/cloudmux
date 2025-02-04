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
	"fmt"

	"github.com/nyl1001/pkg/errors"
	"github.com/nyl1001/pkg/jsonutils"

	api "github.com/nyl1001/cloudmux/pkg/apis/compute"
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
	"github.com/nyl1001/cloudmux/pkg/multicloud"
	"github.com/nyl1001/cloudmux/pkg/multicloud/huawei"
)

// date: 2019.07.15
// In Huawei cloud, there are only two routing tables in a vpc, which are
// self-defined routing tables and peer-to-peer routing tables.
// The routing in these two tables is different, one's NextHop is a IP address and
// the other one's NextHop address is a instance ID of peer-to-peer connection.
// The former has no id and it's Type is ROUTE_TYPR_IP, and the latter's Type is ROUTE_TYPE_PEER.

const (
	ROUTE_TYPR_IP   = "IP"
	ROUTE_TYPE_PEER = "peering"
)

type SRouteEntry struct {
	multicloud.SResourceBase
	huawei.HuaweiTags
	routeTable *SRouteTable

	ID          string // route ID
	Type        string // route type
	Destination string // route destination
	NextHop     string // route next hop (ip or id)
}

func (route *SRouteEntry) GetId() string {
	if len(route.ID) == 0 {
		return route.Destination + ":" + route.NextHop
	}
	return route.ID
}

func (route *SRouteEntry) GetName() string {
	return ""
}

func (route *SRouteEntry) GetGlobalId() string {
	return route.GetId()
}

func (route *SRouteEntry) GetStatus() string {
	return api.ROUTE_ENTRY_STATUS_AVAILIABLE
}

func (route *SRouteEntry) Refresh() error {
	return nil
}

func (route *SRouteEntry) IsEmulated() bool {
	return false
}

func (route *SRouteEntry) GetType() string {
	if route.Type == ROUTE_TYPE_PEER {
		return api.ROUTE_ENTRY_TYPE_CUSTOM
	}
	return api.ROUTE_ENTRY_TYPE_SYSTEM
}

func (route *SRouteEntry) GetCidr() string {
	return route.Destination
}

func (route *SRouteEntry) GetNextHopType() string {
	// In Huawei Cloud, NextHopType is same with itself
	switch route.Type {
	case ROUTE_TYPE_PEER:
		return api.NEXT_HOP_TYPE_VPCPEERING
	default:
		return ""
	}
}

func (route *SRouteEntry) GetNextHop() string {
	return route.NextHop
}

// SRouteTable has no ID and Name because there is no id or name of route table in huawei cloud.
// And some method such as GetId and GetName of ICloudRouteTable has no practical meaning
type SRouteTable struct {
	multicloud.SResourceBase
	huawei.HuaweiTags
	region *SRegion
	vpc    *SVpc

	VpcId       string
	Description string
	Type        string
	Routes      []*SRouteEntry
}

func NewSRouteTable(vpc *SVpc, Type string) SRouteTable {
	return SRouteTable{
		region: vpc.region,
		vpc:    vpc,
		Type:   Type,
		VpcId:  vpc.GetId(),
	}

}

func (self *SRouteTable) GetId() string {
	return self.GetGlobalId()
}

func (self *SRouteTable) GetName() string {
	return ""
}

func (self *SRouteTable) GetGlobalId() string {
	return fmt.Sprintf("%s-%s", self.GetVpcId(), self.GetType())
}

func (self *SRouteTable) GetStatus() string {
	return api.ROUTE_TABLE_AVAILABLE
}

func (self *SRouteTable) Refresh() error {
	return nil
}

func (self *SRouteTable) IsEmulated() bool {
	return false
}

func (self *SRouteTable) GetDescription() string {
	return self.Description
}

func (self *SRouteTable) GetRegionId() string {
	return self.region.GetId()
}

func (self *SRouteTable) GetVpcId() string {
	return self.VpcId
}

func (self *SRouteTable) GetType() cloudprovider.RouteTableType {
	return cloudprovider.RouteTableTypeSystem
}

func (self *SRouteTable) GetIRoutes() ([]cloudprovider.ICloudRoute, error) {
	if self.Routes == nil {
		err := self.fetchRoutes()
		if err != nil {
			return nil, err
		}
	}
	ret := []cloudprovider.ICloudRoute{}
	for i := range self.Routes {
		ret = append(ret, self.Routes[i])
	}
	return ret, nil
}

// fetchRoutes fetch Routes
func (self *SRouteTable) fetchRoutes() error {
	if self.Type == ROUTE_TYPR_IP {
		return self.fetchRoutesForIP()
	}
	return self.fetchRoutesForPeer()
}

// fetchRoutesForIP fetch the Routes which Type is ROUTE_TYPR_IP through vpc's get api
func (self *SRouteTable) fetchRoutesForIP() error {
	ret, err := self.region.ecsClient.Vpcs.Get(self.GetVpcId(), map[string]string{})
	if err != nil {
		return errors.Wrap(err, "get vpc info error")
	}
	routeArray, err := ret.GetArray("routes")
	routes := make([]*SRouteEntry, 0, len(routeArray))
	for i := range routeArray {
		destination, err := routeArray[i].GetString("destination")
		if err != nil {
			return errors.Wrap(err, "get destination of route error")
		}
		nextHop, err := routeArray[i].GetString("nexthop")
		if err != nil {
			return errors.Wrap(err, "get nexthop of route error")
		}
		routes = append(routes, &SRouteEntry{
			routeTable:  self,
			ID:          "",
			Type:        ROUTE_TYPR_IP,
			Destination: destination,
			NextHop:     nextHop,
		})
	}
	self.Routes = routes
	return nil
}

// fetchRoutesForPeer fetch the routes which Type is ROUTE_TYPE_PEER through vpcRoute's list api
func (self *SRouteTable) fetchRoutesForPeer() error {
	retPeer, err := self.region.ecsClient.VpcRoutes.List(map[string]string{"vpc_id": self.GetVpcId()})
	if err != nil {
		return errors.Wrap(err, "get peer route error")
	}
	routesPeer := make([]*SRouteEntry, 0, retPeer.Total)
	for i := range retPeer.Data {
		route := retPeer.Data[i]
		id, err := route.GetString("id")
		if err != nil {
			return errors.Wrap(err, "get id of peer route error")
		}
		destination, err := route.GetString("destination")
		if err != nil {
			return errors.Wrap(err, "get destination of peer route error")
		}
		nextHop, err := route.GetString("nexthop")
		if err != nil {
			return errors.Wrap(err, "get nexthop of peer route error")
		}
		routesPeer = append(routesPeer, &SRouteEntry{
			routeTable:  self,
			ID:          id,
			Type:        ROUTE_TYPE_PEER,
			Destination: destination,
			NextHop:     nextHop,
		})
	}
	self.Routes = routesPeer
	return nil
}

func (self *SRouteTable) GetAssociations() []cloudprovider.RouteTableAssociation {
	result := []cloudprovider.RouteTableAssociation{}
	return result
}

func (self *SRouteTable) CreateRoute(route cloudprovider.RouteSet) error {
	if route.NextHopType != api.NEXT_HOP_TYPE_VPCPEERING {
		return cloudprovider.ErrNotSupported
	}
	err := self.region.CreatePeeringRoute(self.vpc.GetId(), route.Destination, route.NextHop)
	if err != nil {
		return errors.Wrapf(err, " self.region.CreatePeeringRoute(%s,%s,%s)", self.vpc.GetId(), route.Destination, route.NextHop)
	}
	return nil
}

func (self *SRouteTable) UpdateRoute(route cloudprovider.RouteSet) error {
	err := self.RemoveRoute(route)
	if err != nil {
		return errors.Wrap(err, "self.RemoveRoute(route)")
	}
	err = self.CreateRoute(route)
	if err != nil {
		return errors.Wrap(err, "self.CreateRoute(route)")
	}
	return nil
}

func (self *SRouteTable) RemoveRoute(route cloudprovider.RouteSet) error {
	err := self.region.DeletePeeringRoute(route.RouteId)
	if err != nil {
		return errors.Wrapf(err, "self.region.DeletePeeringRoute(%s)", route.RouteId)
	}
	return nil
}

// GetRouteTables return []SRouteTable of self
func (self *SVpc) getRouteTables() ([]SRouteTable, error) {
	// every Vpc has two route table in Huawei Cloud
	routeTableIp := NewSRouteTable(self, ROUTE_TYPR_IP)
	routeTablePeer := NewSRouteTable(self, ROUTE_TYPE_PEER)
	if err := routeTableIp.fetchRoutesForIP(); err != nil {
		return nil, errors.Wrap(err, `get route table whilc type is "ip" error`)
	}
	if err := routeTablePeer.fetchRoutesForPeer(); err != nil {
		return nil, errors.Wrap(err, `get route table whilc type is "peering" error`)
	}
	ret := make([]SRouteTable, 0, 2)
	if len(routeTableIp.Routes) != 0 {
		ret = append(ret, routeTableIp)
	}
	if len(routeTablePeer.Routes) != 0 {
		ret = append(ret, routeTablePeer)
	}
	return ret, nil
}

// GetRouteTables return []SRouteTable of vpc which id is vpcId if vpcId is no-nil,
// otherwise return []SRouteTable of all vpc in this SRegion
func (self *SRegion) GetRouteTables(vpcId string) ([]SRouteTable, error) {
	vpcs, err := self.GetVpcs()
	if err != nil {
		return nil, errors.Wrap(err, "Get Vpcs error")
	}
	if vpcId != "" {
		for i := range vpcs {
			if vpcs[i].GetId() == vpcId {
				vpcs = vpcs[i : i+1]
				break
			}
		}
	}
	ret := make([]SRouteTable, 0, 2*len(vpcs))
	for _, vpc := range vpcs {
		routetables, err := vpc.getRouteTables()
		if err != nil {
			return nil, errors.Wrapf(err, "get vpc's route tables whilch id is %s error", vpc.GetId())
		}
		ret = append(ret, routetables...)

	}
	return ret, nil
}

func (self *SRegion) CreatePeeringRoute(vpcId, destinationCidr, target string) error {
	params := jsonutils.NewDict()
	routeObj := jsonutils.NewDict()
	routeObj.Set("type", jsonutils.NewString("peering"))
	routeObj.Set("nexthop", jsonutils.NewString(target))
	routeObj.Set("destination", jsonutils.NewString(destinationCidr))
	routeObj.Set("vpc_id", jsonutils.NewString(vpcId))
	params.Set("route", routeObj)
	err := DoCreate(self.ecsClient.VpcRoutes.Create, params, nil)
	if err != nil {
		return errors.Wrapf(err, "DoCreate(self.ecsClient.VpcRoutes.Create, %s, &ret)", jsonutils.Marshal(params).String())
	}
	return nil
}

func (self *SRegion) DeletePeeringRoute(routeId string) error {
	err := DoDelete(self.ecsClient.VpcRoutes.Delete, routeId, nil, nil)
	if err != nil {
		return errors.Wrapf(err, "DoDelete(self.ecsClient.VpcRoutes.Delete,%s,nil)", routeId)
	}
	return nil
}
