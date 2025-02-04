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

package aliyun

import (
	"context"
	"fmt"

	"github.com/nyl1001/pkg/jsonutils"

	api "github.com/nyl1001/cloudmux/pkg/apis/compute"
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
	"github.com/nyl1001/cloudmux/pkg/multicloud"
)

type SLoadbalancerUDPListener struct {
	multicloud.SResourceBase
	multicloud.SLoadbalancerRedirectBase
	AliyunTags
	lb *SLoadbalancer

	ListenerPort      int    //	负载均衡实例前端使用的端口。
	BackendServerPort int    //	负载均衡实例后端使用的端口。
	Bandwidth         int    //	监听的带宽峰值。
	Status            string //	当前监听的状态，取值：starting | running | configuring | stopping | stopped
	Description       string

	Scheduler                string //	调度算法
	VServerGroupId           string //	绑定的服务器组ID。
	MasterSlaveServerGroupId string //	绑定的主备服务器组ID。
	AclStatus                string //	是否开启访问控制功能。取值：on | off（默认值）

	AclType string //	访问控制类型：

	AclId string //	监听绑定的访问策略组ID。当AclStatus参数的值为on时，该参数必选。

	HealthCheck               string //	是否开启健康检查。
	HealthyThreshold          int    //	健康检查阈值。
	UnhealthyThreshold        int    //	不健康检查阈值。
	HealthCheckConnectTimeout int    //	每次健康检查响应的最大超时间，单位为秒。
	HealthCheckInterval       int    //	健康检查的时间间隔，单位为秒。
	HealthCheckConnectPort    int    //	健康检查的端口。

	HealthCheckExp string // UDP监听健康检查的响应串
	HealthCheckReq string // UDP监听健康检查的请求串
}

func (listener *SLoadbalancerUDPListener) GetName() string {
	if len(listener.Description) == 0 {
		listener.Refresh()
	}
	if len(listener.Description) > 0 {
		return listener.Description
	}
	return fmt.Sprintf("UDP:%d", listener.ListenerPort)
}

func (listerner *SLoadbalancerUDPListener) GetId() string {
	return fmt.Sprintf("%s/%d", listerner.lb.LoadBalancerId, listerner.ListenerPort)
}

func (listerner *SLoadbalancerUDPListener) GetGlobalId() string {
	return listerner.GetId()
}

func (listerner *SLoadbalancerUDPListener) GetStatus() string {
	switch listerner.Status {
	case "starting", "running":
		return api.LB_STATUS_ENABLED
	case "configuring", "stopping", "stopped":
		return api.LB_STATUS_DISABLED
	default:
		return api.LB_STATUS_UNKNOWN
	}
}

func (listerner *SLoadbalancerUDPListener) GetEgressMbps() int {
	if listerner.Bandwidth < 1 {
		return 0
	}
	return listerner.Bandwidth
}

func (listerner *SLoadbalancerUDPListener) IsEmulated() bool {
	return false
}

func (listerner *SLoadbalancerUDPListener) Refresh() error {
	lis, err := listerner.lb.region.GetLoadbalancerUDPListener(listerner.lb.LoadBalancerId, listerner.ListenerPort)
	if err != nil {
		return err
	}
	return jsonutils.Update(listerner, lis)
}

func (listerner *SLoadbalancerUDPListener) GetListenerType() string {
	return "udp"
}

func (listerner *SLoadbalancerUDPListener) GetListenerPort() int {
	return listerner.ListenerPort
}

func (listerner *SLoadbalancerUDPListener) GetBackendGroupId() string {
	if len(listerner.VServerGroupId) > 0 {
		return listerner.VServerGroupId
	}
	return listerner.MasterSlaveServerGroupId
}

func (listerner *SLoadbalancerUDPListener) GetBackendServerPort() int {
	return listerner.BackendServerPort
}

func (listerner *SLoadbalancerUDPListener) GetScheduler() string {
	return listerner.Scheduler
}

func (listerner *SLoadbalancerUDPListener) GetAclStatus() string {
	return listerner.AclStatus
}

func (listerner *SLoadbalancerUDPListener) GetAclType() string {
	return listerner.AclType
}

func (listerner *SLoadbalancerUDPListener) GetAclId() string {
	return listerner.AclId
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheck() string {
	return listerner.HealthCheck
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckType() string {
	return api.LB_HEALTH_CHECK_UDP
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckDomain() string {
	return ""
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckURI() string {
	return ""
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckCode() string {
	return ""
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckRise() int {
	return listerner.HealthyThreshold
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckFail() int {
	return listerner.UnhealthyThreshold
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckTimeout() int {
	return listerner.HealthCheckConnectTimeout
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckInterval() int {
	return listerner.HealthCheckInterval
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckReq() string {
	return listerner.HealthCheckReq
}

func (listerner *SLoadbalancerUDPListener) GetHealthCheckExp() string {
	return listerner.HealthCheckExp
}

func (listerner *SLoadbalancerUDPListener) GetStickySession() string {
	return ""
}

func (listerner *SLoadbalancerUDPListener) GetStickySessionType() string {
	return ""
}

func (listerner *SLoadbalancerUDPListener) GetStickySessionCookie() string {
	return ""
}

func (listerner *SLoadbalancerUDPListener) GetStickySessionCookieTimeout() int {
	return 0
}

func (listerner *SLoadbalancerUDPListener) XForwardedForEnabled() bool {
	return false
}

func (listerner *SLoadbalancerUDPListener) GzipEnabled() bool {
	return false
}

func (listerner *SLoadbalancerUDPListener) GetCertificateId() string {
	return ""
}

func (listerner *SLoadbalancerUDPListener) GetTLSCipherPolicy() string {
	return ""
}

func (listerner *SLoadbalancerUDPListener) HTTP2Enabled() bool {
	return false
}

func (listerner *SLoadbalancerUDPListener) ChangeCertificate(ctx context.Context, opts *cloudprovider.ListenerCertificateOptions) error {
	return cloudprovider.ErrNotSupported
}

func (self *SLoadbalancerUDPListener) SetAcl(ctx context.Context, opts *cloudprovider.ListenerAclOptions) error {
	params := map[string]string{
		"LoadBalancerId": self.lb.LoadBalancerId,
		"ListenerPort":   fmt.Sprintf("%d", self.ListenerPort),
		"AclStatus":      "off",
	}
	if opts.AclStatus == api.LB_BOOL_ON {
		params["AclStatus"] = "on"
		params["AclType"] = opts.AclType
		params["AclId"] = opts.AclId
	}
	_, err := self.lb.region.lbRequest("SetLoadBalancerUDPListenerAttribute", params)
	return err
}

func (listerner *SLoadbalancerUDPListener) GetILoadbalancerListenerRules() ([]cloudprovider.ICloudLoadbalancerListenerRule, error) {
	return []cloudprovider.ICloudLoadbalancerListenerRule{}, nil
}

func (region *SRegion) GetLoadbalancerUDPListener(loadbalancerId string, listenerPort int) (*SLoadbalancerUDPListener, error) {
	params := map[string]string{}
	params["RegionId"] = region.RegionId
	params["LoadBalancerId"] = loadbalancerId
	params["ListenerPort"] = fmt.Sprintf("%d", listenerPort)
	body, err := region.lbRequest("DescribeLoadBalancerUDPListenerAttribute", params)
	if err != nil {
		return nil, err
	}
	listener := SLoadbalancerUDPListener{}
	return &listener, body.Unmarshal(&listener)
}

func (region *SRegion) CreateLoadbalancerUDPListener(lb *SLoadbalancer, listener *cloudprovider.SLoadbalancerListenerCreateOptions) (cloudprovider.ICloudLoadbalancerListener, error) {
	params := region.constructBaseCreateListenerParams(lb, listener)
	_, err := region.lbRequest("CreateLoadBalancerUDPListener", params)
	if err != nil {
		return nil, err
	}
	iListener, err := region.GetLoadbalancerUDPListener(lb.LoadBalancerId, listener.ListenerPort)
	if err != nil {
		return nil, err
	}
	iListener.lb = lb
	return iListener, nil
}

func (listerner *SLoadbalancerUDPListener) Delete(ctx context.Context) error {
	return listerner.lb.region.DeleteLoadbalancerListener(listerner.lb.LoadBalancerId, listerner.ListenerPort)
}

func (listerner *SLoadbalancerUDPListener) CreateILoadBalancerListenerRule(rule *cloudprovider.SLoadbalancerListenerRule) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (listerner *SLoadbalancerUDPListener) GetILoadBalancerListenerRuleById(ruleId string) (cloudprovider.ICloudLoadbalancerListenerRule, error) {
	return nil, cloudprovider.ErrNotSupported
}

func (listerner *SLoadbalancerUDPListener) Start() error {
	return listerner.lb.region.startListener(listerner.ListenerPort, listerner.lb.LoadBalancerId)
}

func (listerner *SLoadbalancerUDPListener) Stop() error {
	return listerner.lb.region.stopListener(listerner.ListenerPort, listerner.lb.LoadBalancerId)
}

func (self *SLoadbalancerUDPListener) ChangeScheduler(ctx context.Context, opts *cloudprovider.ChangeListenerSchedulerOptions) error {
	params := map[string]string{
		"LoadBalancerId": self.lb.LoadBalancerId,
		"ListenerPort":   fmt.Sprintf("%d", self.ListenerPort),
		"Scheduler":      opts.Scheduler,
	}
	_, err := self.lb.region.lbRequest("SetLoadBalancerUDPListenerAttribute", params)
	return err
}

func (self *SLoadbalancerUDPListener) SetHealthCheck(ctx context.Context, opts *cloudprovider.ListenerHealthCheckOptions) error {
	params := map[string]string{
		"LoadBalancerId":    self.lb.LoadBalancerId,
		"ListenerPort":      fmt.Sprintf("%d", self.ListenerPort),
		"HealthCheckSwitch": "off",
	}
	if opts.HealthCheck == api.LB_BOOL_ON {
		params["HealthCheckSwitch"] = "on"
		if opts.HealthCheckTimeout >= 1 && opts.HealthCheckTimeout <= 300 {
			params["HealthCheckConnectTimeout"] = fmt.Sprintf("%d", opts.HealthCheckTimeout)
		}
		if len(opts.HealthCheckReq) > 0 {
			params["healthCheckReq"] = opts.HealthCheckReq
		}
		if len(opts.HealthCheckExp) > 0 {
			params["healthCheckExp"] = opts.HealthCheckExp
		}

		if len(opts.HealthCheckDomain) > 0 {
			params["HealthCheckDomain"] = opts.HealthCheckDomain
		}

		if len(opts.HealthCheckHttpCode) > 0 {
			params["HealthCheckHttpCode"] = opts.HealthCheckHttpCode
		}

		if len(opts.HealthCheckURI) > 0 {
			params["HealthCheckURI"] = opts.HealthCheckURI
		}

		if opts.HealthCheckRise >= 2 && opts.HealthCheckRise <= 10 {
			params["HealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckRise)
		}

		if opts.HealthCheckFail >= 2 && opts.HealthCheckFail <= 10 {
			params["UnhealthyThreshold"] = fmt.Sprintf("%d", opts.HealthCheckFail)
		}
		if opts.HealthCheckInterval >= 1 && opts.HealthCheckInterval <= 50 {
			params["healthCheckInterval"] = fmt.Sprintf("%d", opts.HealthCheckInterval)
		}
	}
	_, err := self.lb.region.lbRequest("SetLoadBalancerUDPListenerAttribute", params)
	return err
}

func (listerner *SLoadbalancerUDPListener) GetProjectId() string {
	return listerner.lb.GetProjectId()
}

func (listerner *SLoadbalancerUDPListener) GetClientIdleTimeout() int {
	return 0
}

func (listerner *SLoadbalancerUDPListener) GetBackendConnectTimeout() int {
	return 0
}
