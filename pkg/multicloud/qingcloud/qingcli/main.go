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

package qingcli

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/http/httpproxy"

	"github.com/nyl1001/pkg/util/shellutils"
	"github.com/nyl1001/structarg"

	"github.com/nyl1001/cloudmux/pkg/cloudprovider"
	"github.com/nyl1001/cloudmux/pkg/multicloud/qingcloud"
	_ "github.com/nyl1001/cloudmux/pkg/multicloud/qingcloud/shell"
)

type BaseOptions struct {
	Debug           bool   `help:"debug mode"`
	AccessKeyId     string `help:"AccessKeyId" default:"$QINGCLOUD_ACCESS_KEY_ID" metavar:"QINGCLOUD_ACCESS_KEY_ID"`
	AccessKeySecret string `help:"AccessKeySecret" default:"$QINGCLOUD_ACCESS_KEY_SECRET" metavar:"QINGCLOUD_ACCESS_KEY_SECRET"`
	RegionId        string `help:"RegionId" default:"$QINGCLOUD_REGION_ID|pek3" metavar:"QINGCLOUD_REGION_ID"`
	SUBCOMMAND      string `help:"qingcli subcommand" subcommand:"true"`
}

func getSubcommandParser() (*structarg.ArgumentParser, error) {
	parse, e := structarg.NewArgumentParserWithHelp(&BaseOptions{},
		"qingcli",
		"Command-line interface to qingcloudc API.",
		`See "qingcli COMMAND --help" for help on a specific command.`)

	if e != nil {
		return nil, e
	}

	subcmd := parse.GetSubcommand()
	if subcmd == nil {
		return nil, fmt.Errorf("No subcommand argument.")
	}
	for _, v := range shellutils.CommandTable {
		_, e := subcmd.AddSubParserWithHelp(v.Options, v.Command, v.Desc, v.Callback)
		if e != nil {
			return nil, e
		}
	}
	return parse, nil
}

func showErrorAndExit(e error) {
	fmt.Fprintf(os.Stderr, "%s", e)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func newClient(options *BaseOptions) (*qingcloud.SRegion, error) {
	if len(options.AccessKeyId) == 0 {
		return nil, fmt.Errorf("Missing access key id")
	}

	if len(options.AccessKeySecret) == 0 {
		return nil, fmt.Errorf("Missing access key secret")
	}

	cfg := &httpproxy.Config{
		HTTPProxy:  os.Getenv("HTTP_PROXY"),
		HTTPSProxy: os.Getenv("HTTPS_PROXY"),
		NoProxy:    os.Getenv("NO_PROXY"),
	}
	cfgProxyFunc := cfg.ProxyFunc()
	proxyFunc := func(req *http.Request) (*url.URL, error) {
		return cfgProxyFunc(req.URL)
	}

	cli, err := qingcloud.NewQingCloudClient(
		qingcloud.NewQingCloudClientConfig(
			options.AccessKeyId,
			options.AccessKeySecret,
		).Debug(options.Debug).
			CloudproviderConfig(
				cloudprovider.ProviderConfig{
					ProxyFunc: proxyFunc,
				},
			),
	)
	if err != nil {
		return nil, err
	}

	return cli.GetRegion(options.RegionId)
}

func Main() {
	parser, e := getSubcommandParser()
	if e != nil {
		showErrorAndExit(e)
	}

	e = parser.ParseArgs(os.Args[1:], false)
	options := parser.Options().(*BaseOptions)

	if parser.IsHelpSet() {
		fmt.Print(parser.HelpString())
		return
	}
	subcmd := parser.GetSubcommand()
	subparser := subcmd.GetSubParser()
	if e != nil || subparser == nil {
		if subparser != nil {
			fmt.Print(subparser.Usage())
		} else {
			fmt.Print(parser.Usage())
		}
		showErrorAndExit(e)
		return
	}
	suboptions := subparser.Options()
	if subparser.IsHelpSet() {
		fmt.Print(subparser.HelpString())
		return
	}
	var region *qingcloud.SRegion
	region, e = newClient(options)
	if e != nil {
		showErrorAndExit(e)
	}
	e = subcmd.Invoke(region, suboptions)
	if e != nil {
		showErrorAndExit(e)
	}
}
