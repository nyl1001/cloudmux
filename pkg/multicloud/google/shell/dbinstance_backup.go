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

	"yunion.io/x/cloudmux/pkg/multicloud/google"
)

func init() {
	type DBInstanceBackupListOptions struct {
		RDS string
	}
	shellutils.R(&DBInstanceBackupListOptions{}, "dbinstance-backup-list", "List dbinstance backup", func(cli *google.SRegion, args *DBInstanceBackupListOptions) error {
		backups, err := cli.GetDBInstanceBackups(args.RDS)
		if err != nil {
			return err
		}
		printList(backups, 0, 0, 0, nil)
		return nil
	})

	type DBInstanceBackupCreateOptions struct {
		RDS  string
		NAME string
		Desc string
	}

	shellutils.R(&DBInstanceBackupCreateOptions{}, "dbinstance-backup-create", "Create dbinstance backup", func(cli *google.SRegion, args *DBInstanceBackupCreateOptions) error {
		return cli.CreateDBInstanceBackup(args.RDS, args.NAME, args.Desc)
	})
}
