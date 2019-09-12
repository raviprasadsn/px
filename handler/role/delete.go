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
	"fmt"

	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"github.com/portworx/pxc/cmd"
	"github.com/portworx/pxc/pkg/commander"
	"github.com/portworx/pxc/pkg/portworx"
	"github.com/portworx/pxc/pkg/util"
	"github.com/spf13/cobra"
)

var deleteRoleCmd *cobra.Command

var deleteRoleOpts struct {
	req      *api.SdkRoleDeleteRequest
	roleName string
}

var _ = commander.RegisterCommandVar(func() {
	// deleteRoleCmd represents the deleteRole command
	deleteRoleCmd = &cobra.Command{
		Use:   "role delete",
		Short: "Delete a role",
		Long:  "Remove a role and its permission rules by name.",
		Example: `
	# pxc delete role --name <role name>`,
		RunE: deleteRoleExec,
	}
})

var _ = commander.RegisterCommandInit(func() {
	cmd.DeleteAddCommand(deleteRoleCmd)
	deleteRoleCmd.Flags().StringVar(&deleteRoleOpts.roleName, "name", "", "name of the role to be deleted")
	cobra.MarkFlagRequired(deleteRoleCmd.Flags(), "name")
})

func DeleteAddCommand(cmd *cobra.Command) {
	deleteRoleCmd.AddCommand(cmd)
}

func deleteRoleExec(c *cobra.Command, args []string) error {
	ctx, conn, err := portworx.PxConnectDefault()
	if err != nil {
		return err
	}
	defer conn.Close()
	s := api.NewOpenStorageRoleClient(conn)
	_, err = s.Delete(ctx, &api.SdkRoleDeleteRequest{
		Name: deleteRoleOpts.roleName,
	})
	if err != nil {
		fmt.Printf("error in delete\n")
	}
	util.Printf("Role " + deleteRoleOpts.roleName + " deleted ...\n")
	return nil
}
