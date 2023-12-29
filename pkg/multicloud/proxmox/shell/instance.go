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
	"strconv"

	"github.com/nyl1001/pkg/util/shellutils"

	"github.com/nyl1001/cloudmux/pkg/multicloud/proxmox"
)

func init() {
	type InstanceListOptions struct {
		HOST_ID string
	}
	shellutils.R(&InstanceListOptions{}, "instance-list", "list instances", func(cli *proxmox.SRegion, args *InstanceListOptions) error {
		instances, err := cli.GetInstances(args.HOST_ID)
		if err != nil {
			return err
		}
		printList(instances, 0, 0, 0, []string{})
		return nil
	})

	type InstanceIdOptions struct {
		ID string
	}

	shellutils.R(&InstanceIdOptions{}, "instance-show", "show instance", func(cli *proxmox.SRegion, args *InstanceIdOptions) error {
		ret, err := cli.GetInstance(args.ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

	shellutils.R(&InstanceIdOptions{}, "instance-stop", "stop instance", func(cli *proxmox.SRegion, args *InstanceIdOptions) error {
		id, _ := strconv.Atoi(args.ID)
		return cli.StopVm(id)
	})

	shellutils.R(&InstanceIdOptions{}, "instance-delete", "delete instance", func(cli *proxmox.SRegion, args *InstanceIdOptions) error {
		id, _ := strconv.Atoi(args.ID)
		return cli.DeleteVM(id)
	})

	type InstanceStartOptions struct {
		ID        string
		Password  string
		PublicKey string
	}

	shellutils.R(&InstanceStartOptions{}, "instance-start", "Start instance", func(cli *proxmox.SRegion, args *InstanceStartOptions) error {
		id, _ := strconv.Atoi(args.ID)
		return cli.StartVm(id)
	})

	type InstanceAttachDiskOptions struct {
		NODE   string
		VM_ID  int
		DRIVER string
	}

	shellutils.R(&InstanceAttachDiskOptions{}, "instance-detach-disk", "Attach instance disk", func(cli *proxmox.SRegion, args *InstanceAttachDiskOptions) error {
		return cli.DetachDisk(args.NODE, args.VM_ID, args.DRIVER)
	})

	type InstanceChangeConfigOptions struct {
		ID    string
		Cpu   int
		MemMb int
	}

	shellutils.R(&InstanceChangeConfigOptions{}, "instance-change-config", "Change instance config", func(cli *proxmox.SRegion, args *InstanceChangeConfigOptions) error {
		id, _ := strconv.Atoi(args.ID)
		return cli.ChangeConfig(id, args.Cpu, args.MemMb)
	})

	type InstanceCreateOptions struct {
		Name  string
		Node  string
		Cpu   int
		MemMb int
	}

	shellutils.R(&InstanceCreateOptions{}, "instance-create", "create instance ", func(cli *proxmox.SRegion, args *InstanceCreateOptions) error {
		ret, err := cli.GenVM(args.Name, args.Node, args.Cpu, args.MemMb)
		printObject(ret)
		return err
	})

	type InstanceVncOptions struct {
		NODE  string
		VM_ID int
	}

	shellutils.R(&InstanceVncOptions{}, "instance-vnc", "show instance vnc", func(cli *proxmox.SRegion, args *InstanceVncOptions) error {
		ret, err := cli.GetVNCInfo(args.NODE, args.VM_ID)
		if err != nil {
			return err
		}
		printObject(ret)
		return nil
	})

}
