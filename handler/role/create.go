/*
Copyright © 2019 Portworx

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package role

import (
	"encoding/json"
	"fmt"
	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"github.com/portworx/pxc/cmd"
	"github.com/portworx/pxc/pkg/commander"
	"github.com/portworx/pxc/pkg/portworx"
	"github.com/portworx/pxc/pkg/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	_ "os"
)

type createRoleOpts struct {
	req      *api.SdkRoleCreateRequest
	roleconf string
}

var (
	crOpts        *createRoleOpts
	createRoleCmd *cobra.Command
)

func CreateAddCommand(c *cobra.Command) {
	createRoleCmd.AddCommand(c)
}

var _ = commander.RegisterCommandVar(func() {
	// createRoleCmd represents the createRole command
	crOpts = &createRoleOpts{
		req: &api.SdkRoleCreateRequest{},
	}
	createRoleCmd = &cobra.Command{
		Use:   "role create",
		Short: "Create a role in Portworx",
		Example: `
	# Create a role using a json file which specifies the role and its rules.
	# A role consist of a set of rules defining services and api's which are allowable.
	# e.g. Rule file(test.json) which allows inspection of any object and listings of only volumes:
		{
			"name": "test.view",
			"rules": [
				{
					"services": [
						"volumes"
					],
					"apis": [
						"*enumerate*"
					]
				},
				{
					"services": [
						"*"
					],
					"apis": [
						"inspect*"
					]
				}
			]
		}

	pxc create role --role-config test.json`,

		RunE: createRoleExec,
	}
})

var _ = commander.RegisterCommandInit(func() {
	cmd.CreateAddCommand(createRoleCmd)

	createRoleCmd.Flags().StringVar(&crOpts.roleconf, "role-config", "", "Required role json file'")
	cobra.MarkFlagRequired(createRoleCmd.Flags(), "role-config")
})

func loadRoleCfg(roleFile string) (*api.SdkRole, error) {

	var role api.SdkRole

	data, err := ioutil.ReadFile(roleFile)

	if err != nil {
		fmt.Errorf("unable to read role file %v\n", err)
	}

	if err := json.Unmarshal(data, &role); err != nil {
		fmt.Errorf("Failed to process role definition data, %v", err)
	}

	return &role, err
}

func createRoleExec(c *cobra.Command, args []string) error {

	ctx, conn, err := portworx.PxConnectDefault()
	defer conn.Close()

	s := api.NewOpenStorageRoleClient(conn)

	if _, err := os.Stat(crOpts.roleconf); err != nil {
		fmt.Errorf("unable to read role file %s\n", crOpts.roleconf)
	}

	role, err := loadRoleCfg(crOpts.roleconf)
	if err != nil {
		fmt.Errorf("role create error, %v", err)
	}

	_, err = s.Create(
		ctx,
		&api.SdkRoleCreateRequest{Role: role})
	if err != nil {
		fmt.Errorf("Role create failed %v", err)
	}
	util.Printf("Role " + role.Name + " created ...\n")

	return nil
}
