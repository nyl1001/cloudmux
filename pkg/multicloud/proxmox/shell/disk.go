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

package shell

import (
	"github.com/nyl1001/pkg/util/shellutils"

	"github.com/nyl1001/cloudmux/pkg/multicloud/proxmox"
)

func init() {
	type DiskListOptions struct {
		NODE    string
		STORAGE string
	}
	shellutils.R(&DiskListOptions{}, "disk-list", "list disks", func(cli *proxmox.SRegion, args *DiskListOptions) error {
		disks, err := cli.GetDisks(args.NODE, args.STORAGE)
		if err != nil {
			return err
		}
		printList(disks, 0, 0, 0, []string{})
		return nil
	})

	type DiskResizeOptions struct {
		NODE    string
		VM_ID   string
		DRIVER  string
		SIZE_GB int
	}

	shellutils.R(&DiskResizeOptions{}, "disk-resize", "resize disk size", func(cli *proxmox.SRegion, args *DiskResizeOptions) error {
		return cli.ResizeDisk(args.NODE, args.VM_ID, args.DRIVER, args.SIZE_GB)
	})

}
