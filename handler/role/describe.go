// Copyright © 2019 Portworx
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

package role

import (
	"encoding/json"

	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"github.com/portworx/pxc/cmd"
	"github.com/portworx/pxc/pkg/commander"
	"github.com/portworx/pxc/pkg/portworx"
	"github.com/portworx/pxc/pkg/util"
	"github.com/spf13/cobra"
)

var describeRoleCmd *cobra.Command
var name string

var describeRoleOpts struct {
	req  *api.SdkRoleInspectRequest
	name string
	all  bool
}

var _ = commander.RegisterCommandVar(func() {
	// describeRoleCmd represents the describeRole command
	describeRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "Describe Role",
		Long:  "Display permission rules for a specific role or for all the roles.",
		Example: `
	# To describe all roles
	  pxc describe role --all
	# To describe a specific role
	  pxc describe role --name <role name>`,

		RunE: describeRoleExec,
	}
})

var _ = commander.RegisterCommandInit(func() {
	cmd.DescribeAddCommand(describeRoleCmd)
	describeRoleCmd.Flags().StringVar(&describeRoleOpts.name, "name", "", "show permission rules of a specified role name")
	describeRoleCmd.Flags().BoolVar(&describeRoleOpts.all, "all", false, "show permission rules of a specified role name")
})

func DescribeAddCommand(cmd *cobra.Command) {
	describeRoleCmd.AddCommand(cmd)
}

func describeRoleExec(c *cobra.Command, args []string) error {
	ctx, conn, err := portworx.PxConnectDefault()
	var result map[string]interface{}
	if err != nil {
		return err
	}
	defer conn.Close()
	s := api.NewOpenStorageRoleClient(conn)

	if describeRoleOpts.all {
		enumRoles, _ := s.Enumerate(ctx, &api.SdkRoleEnumerateRequest{})
		b, _ := json.Marshal(enumRoles)
		json.Unmarshal([]byte(b), &result)
		roleNames := result["names"].([]interface{})
		for _, name := range roleNames {
			roleData, _ := s.Inspect(ctx, &api.SdkRoleInspectRequest{
				Name: name.(string),
			})
			b, _ = json.MarshalIndent(roleData, "", " ")
			util.Printf("%s\n", string(b))
		}
	} else {
		roleData, _ := s.Inspect(ctx, &api.SdkRoleInspectRequest{
			Name: describeRoleOpts.name,
		})
		b, _ := json.MarshalIndent(roleData, "", " ")
		util.Printf(string(b))
	}
	return nil
}
