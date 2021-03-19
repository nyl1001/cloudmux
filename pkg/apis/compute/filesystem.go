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

package compute

import "yunion.io/x/onecloud/pkg/apis"

const (
	// 可用
	NAS_STATUS_AVAILABLE = "available"
	// 不可用
	NAS_STATUS_UNAVAILABLE = "unavailable"
	// 扩容中
	NAS_STATUS_EXTENDING = "extending"
	// 创建中
	NAS_STATUS_CREATING = "creating"
	// 创建失败
	NAS_STATUS_CREATE_FAILED = "create_failed"
	// 未知
	NAS_STATUS_UNKNOWN = "unknown"
	// 删除中
	NAS_STATUS_DELETING      = "deleting"
	NAS_STATUS_DELETE_FAILED = "delete_failed"
)

type FileSystemListInput struct {
	apis.StatusInfrasResourceBaseListInput
	apis.ExternalizedResourceBaseListInput
	ManagedResourceListInput

	RegionalFilterListInput
}

type FileSystemCreateInput struct {
	apis.StatusInfrasResourceBaseCreateInput
	// 协议类型
	// enum: NFS, SMB, CPFS
	Protocol string `json:"protocol"`

	// 文件系统类型
	// enmu: extreme, standard, cpfs
	FileSystemType string `json:"file_system_type"`

	// 容量大小
	Capacity int64 `json:"capacity"`

	// IP子网Id
	NetworkId string `json:"network_id"`

	// 存储类型
	// enmu: performance, capacity, standard, advance, advance_100, advance_200
	StorageType string `json:"storage_type"`

	// 可用区Id, 若不指定IP子网，此参数必填
	ZoneId string `json:"zone_id"`

	//swagger:ignore
	CloudregionId string `json:"cloudregion_id"`

	// 订阅Id, 若传入network_id此参数可忽略
	ManagerId string `json:"manager_id"`
}

type FileSystemSyncstatusInput struct {
}

type FileSystemDetails struct {
	apis.StatusInfrasResourceBaseDetails
	ManagedResourceInfo
	CloudregionResourceInfo

	Vpc     string
	Network string
	Zone    string
}
