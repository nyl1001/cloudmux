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
	"context"
	"fmt"
	"time"

	"github.com/nyl1001/pkg/errors"
	"github.com/nyl1001/pkg/utils"
	"github.com/nyl1001/pkg/jsonutils"
	"yunion.io/x/log"

	api "github.com/nyl1001/cloudmux/pkg/apis/compute"
	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
	"github.com/nyl1001/cloudmux/pkg/multicloud"
)

type SMountInstances struct {
	MountInstance []string
}

type STags struct {
	Tag []string
}

type SDisk struct {
	multicloud.SResourceBase
	ApsaraTags
	storage *SStorage
	multicloud.SDisk

	AttachedTime                  time.Time
	AutoSnapshotPolicyId          string
	Category                      string
	PerformanceLevel              string
	CreationTime                  time.Time
	DeleteAutoSnapshot            bool
	DeleteWithInstance            bool
	Description                   string
	DetachedTime                  time.Time
	Device                        string
	DiskChargeType                TChargeType
	DiskId                        string
	DiskName                      string
	EnableAutoSnapshot            bool
	EnableAutomatedSnapshotPolicy bool
	Encrypted                     bool
	ExpiredTime                   time.Time
	ImageId                       string
	InstanceId                    string
	MountInstances                SMountInstances
	OperationLocks                SOperationLocks
	Portable                      bool
	ProductCode                   string
	RegionId                      string
	Size                          int
	SourceSnapshotId              string
	Status                        string
	Tags                          STags
	Type                          string
	ZoneId                        string
	DepartmentInfo
}

func (self *SRegion) GetDisks(instanceId string, zoneId string, category string, diskIds []string, snapshotpolicyId string) ([]SDisk, error) {
	params := make(map[string]string)
	params["RegionId"] = self.RegionId
	params["PageSize"] = "50"
	pageNum := 1
	params["PageNumber"] = fmt.Sprintf("%d", pageNum)

	if len(instanceId) > 0 {
		params["InstanceId"] = instanceId
	}
	if len(zoneId) > 0 {
		params["ZoneId"] = zoneId
	}
	if len(category) > 0 {
		params["Category"] = category
	}
	if len(diskIds) > 0 {
		params["DiskIds"] = jsonutils.Marshal(diskIds).String()
	}
	if len(snapshotpolicyId) > 0 {
		params["AutoSnapshotPolicyId"] = snapshotpolicyId
	}
	ret := []SDisk{}
	for {
		body, err := self.ecsRequest("DescribeDisks", params)
		if err != nil {
			return nil, errors.Wrapf(err, "DescribeDisks")
		}

		part := make([]SDisk, 0)
		err = body.Unmarshal(&part, "Disks", "Disk")
		if err != nil {
			return nil, errors.Wrapf(err, "Unmarshal")
		}
		ret = append(ret, part...)
		total, _ := body.Int("TotalCount")
		if len(ret) >= int(total) || len(part) == 0 {
			break
		}
		pageNum++
		params["PageNumber"] = fmt.Sprintf("%d", pageNum)
	}
	return ret, nil
}

func (self *SDisk) GetId() string {
	return self.DiskId
}

func (self *SDisk) Delete(ctx context.Context) error {
	_, err := self.storage.zone.region.getDisk(self.DiskId)
	if err != nil {
		if errors.Cause(err) == cloudprovider.ErrNotFound {
			// 未找到disk, 说明disk已经被删除了. 避免回收站中disk-delete循环删除失败
			return nil
		}
		log.Errorf("Failed to find disk %s when delete: %s", self.DiskId, err)
		return err
	}

	for {
		err := self.storage.zone.region.DeleteDisk(self.DiskId)
		if err != nil {
			if isError(err, "IncorrectDiskStatus") {
				log.Infof("The disk is initializing, try later ...")
				time.Sleep(10 * time.Second)
			} else {
				log.Errorf("DeleteDisk fail: %s", err)
				return err
			}
		} else {
			break
		}
	}
	return cloudprovider.WaitDeleted(self, 10*time.Second, 300*time.Second) // 5minutes
}

func (self *SDisk) Resize(ctx context.Context, sizeMb int64) error {
	return self.storage.zone.region.resizeDisk(self.DiskId, sizeMb)
}

func (self *SDisk) GetName() string {
	if len(self.DiskName) > 0 {
		return self.DiskName
	}
	return self.DiskId
}

func (self *SDisk) GetGlobalId() string {
	return self.DiskId
}

func (self *SDisk) IsEmulated() bool {
	return false
}

func (self *SDisk) GetIStorage() (cloudprovider.ICloudStorage, error) {
	return self.storage, nil
}

func (self *SDisk) GetStatus() string {
	// In_use Available Attaching Detaching Creating ReIniting All
	switch self.Status {
	case "Creating", "ReIniting":
		return api.DISK_ALLOCATING
	default:
		return api.DISK_READY
	}
}

func (self *SDisk) Refresh() error {
	new, err := self.storage.zone.region.getDisk(self.DiskId)
	if err != nil {
		return err
	}
	return jsonutils.Update(self, new)
}

func (self *SDisk) ResizeDisk(newSize int64) error {
	// newSize 单位为 GB. 范围在20 ～2000. 只能往大调。不能调小
	// https://help.apsara.com/document_detail/25522.html?spm=a2c4g.11174283.6.897.aHwqkS
	return self.storage.zone.region.resizeDisk(self.DiskId, newSize)
}

func (self *SDisk) GetDiskFormat() string {
	return "vhd"
}

func (self *SDisk) GetDiskSizeMB() int {
	return self.Size * 1024
}

func (self *SDisk) GetIsAutoDelete() bool {
	return self.DeleteWithInstance
}

func (self *SDisk) GetTemplateId() string {
	return self.ImageId
}

func (self *SDisk) GetDiskType() string {
	switch self.Type {
	case "system":
		return api.DISK_TYPE_SYS
	case "data":
		return api.DISK_TYPE_DATA
	default:
		return api.DISK_TYPE_DATA
	}
}

func (self *SDisk) GetFsFormat() string {
	return ""
}

func (self *SDisk) GetIsNonPersistent() bool {
	return false
}

func (self *SDisk) GetDriver() string {
	return "scsi"
}

func (self *SDisk) GetCacheMode() string {
	return "none"
}

func (self *SDisk) GetMountpoint() string {
	return ""
}

func (self *SRegion) CreateDisk(zoneId string, category string, name string, sizeGb int, desc string, projectId string) (string, error) {
	params := make(map[string]string)
	params["ZoneId"] = zoneId
	params["DiskName"] = name
	if len(desc) > 0 {
		params["Description"] = desc
	}
	params["Encrypted"] = "false"
	params["DiskCategory"] = category
	if category == api.STORAGE_CLOUD_ESSD_PL2 {
		params["DiskCategory"] = api.STORAGE_CLOUD_ESSD
		params["PerformanceLevel"] = "PL2"
	}
	if category == api.STORAGE_CLOUD_ESSD_PL3 {
		params["DiskCategory"] = api.STORAGE_CLOUD_ESSD
		params["PerformanceLevel"] = "PL3"
	}

	if len(projectId) > 0 {
		params["ResourceGroupId"] = projectId
	}
	params["Size"] = fmt.Sprintf("%d", sizeGb)
	params["ClientToken"] = utils.GenRequestId(20)

	body, err := self.ecsRequest("CreateDisk", params)
	if err != nil {
		return "", err
	}
	return body.GetString("DiskId")
}

func (self *SRegion) getDisk(diskId string) (*SDisk, error) {
	disks, err := self.GetDisks("", "", "", []string{diskId}, "")
	if err != nil {
		return nil, err
	}
	for i := range disks {
		if disks[i].DiskId == diskId {
			return &disks[i], nil
		}
	}
	return nil, errors.Wrapf(cloudprovider.ErrNotFound, diskId)
}

func (self *SRegion) DeleteDisk(diskId string) error {
	params := make(map[string]string)
	params["DiskId"] = diskId

	_, err := self.ecsRequest("DeleteDisk", params)
	return err
}

func (self *SRegion) resizeDisk(diskId string, sizeMb int64) error {
	sizeGb := sizeMb / 1024
	params := make(map[string]string)
	params["DiskId"] = diskId
	params["NewSize"] = fmt.Sprintf("%d", sizeGb)

	_, err := self.ecsRequest("ResizeDisk", params)
	if err != nil {
		log.Errorf("resizing disk (%s) to %d GiB failed: %s", diskId, sizeGb, err)
		return err
	}

	return nil
}

func (self *SRegion) resetDisk(diskId, snapshotId string) error {
	params := make(map[string]string)
	params["DiskId"] = diskId
	params["SnapshotId"] = snapshotId
	_, err := self.ecsRequest("ResetDisk", params)
	if err != nil {
		log.Errorf("ResetDisk %s to snapshot %s fail %s", diskId, snapshotId, err)
		return err
	}

	return nil
}

func (self *SDisk) CreateISnapshot(ctx context.Context, name, desc string) (cloudprovider.ICloudSnapshot, error) {
	if snapshotId, err := self.storage.zone.region.CreateSnapshot(self.DiskId, name, desc); err != nil {
		log.Errorf("createSnapshot fail %s", err)
		return nil, err
	} else if snapshot, err := self.getSnapshot(snapshotId); err != nil {
		return nil, err
	} else {
		snapshot.region = self.storage.zone.region
		if err := cloudprovider.WaitStatus(snapshot, api.SNAPSHOT_READY, 15*time.Second, 3600*time.Second); err != nil {
			return nil, err
		}
		return snapshot, nil
	}
}

func (self *SRegion) CreateSnapshot(diskId, name, desc string) (string, error) {
	params := make(map[string]string)
	params["RegionId"] = self.RegionId
	params["DiskId"] = diskId
	params["SnapshotName"] = name
	params["Description"] = desc

	if body, err := self.ecsRequest("CreateSnapshot", params); err != nil {
		log.Errorf("CreateSnapshot fail %s", err)
		return "", err
	} else {
		return body.GetString("SnapshotId")
	}
}

func (self *SDisk) GetISnapshot(snapshotId string) (cloudprovider.ICloudSnapshot, error) {
	if snapshot, err := self.getSnapshot(snapshotId); err != nil {
		return nil, err
	} else {
		snapshot.region = self.storage.zone.region
		return snapshot, nil
	}
}

func (self *SDisk) getSnapshot(snapshotId string) (*SSnapshot, error) {
	if snapshots, total, err := self.storage.zone.region.GetSnapshots("", "", "", []string{snapshotId}, 0, 1); err != nil {
		return nil, err
	} else if total != 1 {
		return nil, cloudprovider.ErrNotFound
	} else {
		return &snapshots[0], nil
	}
}

func (self *SDisk) GetISnapshots() ([]cloudprovider.ICloudSnapshot, error) {
	snapshots := make([]SSnapshot, 0)
	for {
		parts, total, err := self.storage.zone.region.GetSnapshots("", self.DiskId, "", []string{}, len(snapshots), 20)
		if err != nil {
			return nil, errors.Wrapf(err, "GetSnapshots")
		}
		snapshots = append(snapshots, parts...)
		if len(snapshots) >= total {
			break
		}
	}
	isnapshots := make([]cloudprovider.ICloudSnapshot, len(snapshots))
	for i := 0; i < len(snapshots); i++ {
		snapshots[i].region = self.storage.zone.region
		isnapshots[i] = &snapshots[i]
	}
	return isnapshots, nil
}

func (self *SDisk) Reset(ctx context.Context, snapshotId string) (string, error) {
	return "", self.storage.zone.region.resetDisk(self.DiskId, snapshotId)
}

func (self *SDisk) GetBillingType() string {
	return convertChargeType(self.DiskChargeType)
}

func (self *SDisk) GetCreatedAt() time.Time {
	return self.CreationTime
}

func (self *SDisk) GetExpiredAt() time.Time {
	return convertExpiredAt(self.ExpiredTime)
}

func (self *SDisk) GetAccessPath() string {
	return ""
}

func (self *SDisk) Rebuild(ctx context.Context) error {
	err := self.storage.zone.region.rebuildDisk(self.DiskId)
	if err != nil {
		if isError(err, "IncorrectInstanceStatus") {
			return nil
		}
		log.Errorf("rebuild disk fail %s", err)
		return err
	}
	return nil
}

func (self *SRegion) rebuildDisk(diskId string) error {
	params := make(map[string]string)
	params["DiskId"] = diskId
	_, err := self.ecsRequest("ReInitDisk", params)
	if err != nil {
		log.Errorf("ReInitDisk %s fail %s", diskId, err)
		return err
	}
	return nil
}
