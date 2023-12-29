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

	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
	"github.com/nyl1001/cloudmux/pkg/multicloud/aliyun"
)

func init() {
	type KubeClusterListOptions struct {
		PageSize   int
		PageNumber int
	}
	shellutils.R(&KubeClusterListOptions{}, "kube-cluster-list", "List kube clusters", func(cli *aliyun.SRegion, args *KubeClusterListOptions) error {
		clusters, _, err := cli.GetKubeClusters(args.PageSize, args.PageNumber)
		if err != nil {
			return err
		}
		printList(clusters, 0, 0, 0, []string{})
		return nil
	})

	type KubeClusterIdOptions struct {
		ID string
	}

	shellutils.R(&KubeClusterIdOptions{}, "kube-cluster-show", "Show kube cluster", func(cli *aliyun.SRegion, args *KubeClusterIdOptions) error {
		cluster, err := cli.GetKubeCluster(args.ID)
		if err != nil {
			return err
		}
		printObject(cluster)
		return nil
	})

	shellutils.R(&KubeClusterIdOptions{}, "kube-cluster-delete", "Delete kube cluster", func(cli *aliyun.SRegion, args *KubeClusterIdOptions) error {
		return cli.DeleteKubeCluster(args.ID, false)
	})

	shellutils.R(&cloudprovider.KubeClusterCreateOptions{}, "kube-cluster-create", "Create kube cluster", func(cli *aliyun.SRegion, args *cloudprovider.KubeClusterCreateOptions) error {
		cluster, err := cli.CreateKubeCluster(args)
		if err != nil {
			return err
		}
		printObject(cluster)
		return nil
	})

	type KubeClusterKubeconfigOptions struct {
		ID            string
		Private       bool
		ExpireMinutes int
	}

	shellutils.R(&KubeClusterKubeconfigOptions{}, "kube-cluster-kubeconfig", "Get kube cluster kubeconfig", func(cli *aliyun.SRegion, args *KubeClusterKubeconfigOptions) error {
		config, err := cli.GetKubeConfig(args.ID, args.Private, args.ExpireMinutes)
		if err != nil {
			return err
		}
		printObject(config)
		return nil
	})

}
