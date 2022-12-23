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
	"yunion.io/x/pkg/util/shellutils"

	"yunion.io/x/cloudmux/pkg/multicloud/qcloud"
)

func init() {
	type ProjectListOptions struct {
		Limit  int `default:"1000"`
		Offset int
	}
	shellutils.R(&ProjectListOptions{}, "project-list", "List project", func(cli *qcloud.SRegion, args *ProjectListOptions) error {
		projects, _, err := cli.GetClient().GetProjects(args.Offset, args.Limit)
		if err != nil {
			return err
		}
		printList(projects, 0, 0, 0, nil)
		return nil
	})

	type ProjectCreateOptions struct {
		NAME string
		Desc string
	}

	shellutils.R(&ProjectCreateOptions{}, "project-create", "Create project", func(cli *qcloud.SRegion, args *ProjectCreateOptions) error {
		project, err := cli.GetClient().CreateProject(args.NAME, args.Desc)
		if err != nil {
			return err
		}
		printObject(project)
		return nil
	})

}
