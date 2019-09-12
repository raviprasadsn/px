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

var enumerateRoleCmd *cobra.Command

var _ = commander.RegisterCommandVar(func() {
	// enumerateRoleCmd represents the enumerateRole command
	enumerateRoleCmd = &cobra.Command{
		Use:   "role",
		Short: "List avaliable roles",
		Long:  "Display the role names available for use by a user.",
		Example: `
	# pxc get role [flags]`,

		RunE: enumerateRoleExec,
	}
})

var _ = commander.RegisterCommandInit(func() {
	cmd.GetAddCommand(enumerateRoleCmd)
})

func EnumerateAddCommand(cmd *cobra.Command) {
	enumerateRoleCmd.AddCommand(cmd)
}

func enumerateRoleExec(c *cobra.Command, args []string) error {
	ctx, conn, err := portworx.PxConnectDefault()
	var result map[string]interface{}
	if err != nil {
		return err
	}
	defer conn.Close()
	s := api.NewOpenStorageRoleClient(conn)

	enumRoles, _ := s.Enumerate(ctx, &api.SdkRoleEnumerateRequest{})
	b, _ := json.Marshal(enumRoles)
	json.Unmarshal([]byte(b), &result)
	roleNames := result["names"].([]interface{})
	for _, name := range roleNames {
		util.Printf("%s\n", name.(string))
	}
	return nil
}
