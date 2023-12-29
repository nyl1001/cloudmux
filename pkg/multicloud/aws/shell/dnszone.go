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
	"github.com/nyl1001/cloudmux/pkg/multicloud/aws"
)

func init() {
	type Route53LocationListOptions struct{}
	shellutils.R(&Route53LocationListOptions{}, "dns-location-list", "List route53location", func(cli *aws.SRegion, args *Route53LocationListOptions) error {
		locations, err := cli.GetClient().ListGeoLocations()
		if err != nil {
			return err
		}
		printList(locations, len(locations), 0, 20, []string{})
		return nil
	})

	type HostedZoneListOptions struct{}
	shellutils.R(&HostedZoneListOptions{}, "dns-zone-list", "List hostedzone", func(cli *aws.SRegion, args *HostedZoneListOptions) error {
		hostzones, err := cli.GetClient().GetDnsZones()
		if err != nil {
			return err
		}
		printList(hostzones, len(hostzones), 0, 20, []string{})
		return nil
	})

	type HostedZoneCreateOptions struct {
		NAME   string `help:"Domain name"`
		Type   string `choices:"PublicZone|PrivateZone"`
		Vpc    string `help:"vpc id"`
		Region string `help:"region id"`
	}
	shellutils.R(&HostedZoneCreateOptions{}, "dns-zone-create", "Create hostedzone", func(cli *aws.SRegion, args *HostedZoneCreateOptions) error {
		opts := cloudprovider.SDnsZoneCreateOptions{}
		opts.Name = args.NAME
		opts.ZoneType = cloudprovider.TDnsZoneType(args.Type)
		if len(args.Vpc) > 0 && len(args.Region) > 0 {
			vpc := cloudprovider.SPrivateZoneVpc{}
			vpc.Id = args.Vpc
			vpc.RegionId = args.Region
			opts.Vpcs = []cloudprovider.SPrivateZoneVpc{vpc}
		}
		hostzones, err := cli.GetClient().CreateDnsZone(&opts)
		if err != nil {
			return err
		}
		printObject(hostzones)
		return nil
	})
	type HostedZoneGetOptions struct {
		HOSTEDZONEID string
	}
	shellutils.R(&HostedZoneGetOptions{}, "dns-zone-show", "get hostedzone by id", func(cli *aws.SRegion, args *HostedZoneGetOptions) error {
		hostedzone, err := cli.GetClient().GetDnsZone(args.HOSTEDZONEID)
		if err != nil {
			return err
		}
		printObject(hostedzone)
		return nil
	})

	type HostedZoneAddVpcOptions struct {
		HOSTEDZONEID string
		VPC          string
		REGION       string
	}
	shellutils.R(&HostedZoneAddVpcOptions{}, "dns-zone-add-vpc", "associate vpc with hostedzone", func(cli *aws.SRegion, args *HostedZoneAddVpcOptions) error {
		err := cli.GetClient().AssociateVPCWithHostedZone(args.VPC, args.REGION, args.HOSTEDZONEID)
		if err != nil {
			return err
		}
		return nil
	})

	type HostedZoneRemoveVpcOptions struct {
		HOSTEDZONEID string
		VPC          string
		REGION       string
	}
	shellutils.R(&HostedZoneRemoveVpcOptions{}, "dns-zone-remove-vpc", "disassociate vpc with hostedzone", func(cli *aws.SRegion, args *HostedZoneRemoveVpcOptions) error {
		err := cli.GetClient().DisassociateVPCFromHostedZone(args.VPC, args.REGION, args.HOSTEDZONEID)
		if err != nil {
			return err
		}
		return nil
	})

	type HostedZoneDeleteOptions struct {
		HOSTEDZONEID string
	}
	shellutils.R(&HostedZoneDeleteOptions{}, "dns-zone-delete", "delete hostedzone", func(cli *aws.SRegion, args *HostedZoneDeleteOptions) error {
		err := cli.GetClient().DeleteDnsZone(args.HOSTEDZONEID)
		if err != nil {
			return err
		}
		return nil
	})

	type DnsRecordSetListOptions struct {
		HOSTEDZONEID string
	}
	shellutils.R(&DnsRecordSetListOptions{}, "dns-record-list", "List dnsrecordset", func(cli *aws.SRegion, args *DnsRecordSetListOptions) error {
		dnsrecordsets, err := cli.GetClient().ListResourceRecordSet(args.HOSTEDZONEID)
		if err != nil {
			return err
		}
		printList(dnsrecordsets, len(dnsrecordsets), 0, 20, []string{})
		return nil
	})

	type DnsRecordSetCreateOptions struct {
		HOSTEDZONEID string `help:"HostedzoneId"`
		NAME         string `help:"Domain name"`
		VALUE        string `help:"dns record value"`
		TTL          int64  `help:"ttl"`
		TYPE         string `help:"dns type"`
		PolicyType   string `help:"PolicyType"`
		Id           string
		Identify     string `help:"Identify"`
	}
	shellutils.R(&DnsRecordSetCreateOptions{}, "dns-record-create", "create dnsrecordset", func(cli *aws.SRegion, args *DnsRecordSetCreateOptions) error {
		opts := cloudprovider.DnsRecord{}
		opts.DnsName = args.NAME
		opts.DnsType = cloudprovider.TDnsType(args.TYPE)
		opts.DnsValue = args.VALUE
		opts.Ttl = args.TTL
		_, err := cli.GetClient().ChangeResourceRecordSets("CREATE", args.HOSTEDZONEID, args.NAME, args.Id, opts)
		return err
	})

	type DnsRecordSetupdateOptions struct {
		HOSTEDZONEID string `help:"HostedzoneId"`
		NAME         string `help:"Domain name"`
		VALUE        string `help:"dns record value"`
		TTL          int64  `help:"ttl"`
		TYPE         string `help:"dns type"`
		Identify     string `help:"Identify"`
	}
	shellutils.R(&DnsRecordSetupdateOptions{}, "dns-record-update", "update dnsrecordset", func(cli *aws.SRegion, args *DnsRecordSetupdateOptions) error {
		opts := cloudprovider.DnsRecord{}
		opts.DnsName = args.NAME
		opts.DnsType = cloudprovider.TDnsType(args.TYPE)
		opts.DnsValue = args.VALUE
		opts.Ttl = args.TTL
		_, err := cli.GetClient().ChangeResourceRecordSets("UPSERT", args.HOSTEDZONEID, args.NAME, args.Identify, opts)
		return err
	})

	type DnsRecordDeleteOptions struct {
		HOSTEDZONEID string `help:"HostedzoneId"`
		ID           string `help:"Identify"`
		cloudprovider.DnsRecord
	}
	shellutils.R(&DnsRecordDeleteOptions{}, "dns-record-delete", "delete dnsrecordset", func(cli *aws.SRegion, args *DnsRecordDeleteOptions) error {
		_, err := cli.GetClient().ChangeResourceRecordSets("DELETE", args.HOSTEDZONEID, args.DnsName, args.ID, args.DnsRecord)
		return err
	})

}
