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

	"github.com/nyl1001/cloudmux/pkg/multicloud/ctyun"
)

func init() {
	type SImageListOptions struct {
		ImageType string `help:"image type" choices:"gold|private|shared"`
	}
	shellutils.R(&SImageListOptions{}, "image-list", "List images", func(cli *ctyun.SRegion, args *SImageListOptions) error {
		images, e := cli.GetImages(args.ImageType)
		if e != nil {
			return e
		}
		printList(images, 0, 0, 0, nil)
		return nil
	})

	type ImageIdOptions struct {
		ID string
	}

	shellutils.R(&ImageIdOptions{}, "image-show", "Show image", func(cli *ctyun.SRegion, args *ImageIdOptions) error {
		image, e := cli.GetImage(args.ID)
		if e != nil {
			return e
		}
		printObject(image)
		return nil
	})

}
