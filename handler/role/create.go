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
package volume

import (

	api "github.com/libopenstorage/openstorage-sdk-clients/sdk/golang"
	"github.com/portworx/px/cmd"
	"github.com/portworx/px/pkg/commander"
//	"github.com/portworx/px/pkg/portworx"
	"github.com/portworx/px/pkg/util"
	"github.com/spf13/cobra"
)

type createRoleOpts struct {
	req                *api.SdkRoleCreateRequest
	roleconf	   string
}

var (
	crOpts          *createRoleOpts
	createRoleCmd   *cobra.Command
)

func CreateAddCommand(c *cobra.Command) {
	createRoleCmd.AddCommand(c)
}

var _ = commander.RegisterCommandVar(func() {
	// createRoleCmd represents the createRole command
	crOpts = &createRoleOpts{
		req: &api.SdkRoleCreateRequest{
		},
	}
	createRoleCmd = &cobra.Command{
		Use:   "role create",
		Short: "Create a role in Portworx",

		// TODO:
		Example: `1. Create volume called "myvolume" with size as 3GiB:
	$ px create volume myvolume --size=3
2. Create volume called "myvolume" with size as 3GiB and replica set to 3:
	$ px create volume myvolume --size=3 --replicas=3
3. Create shared volume called "myvolume" with size as 3GiB:
	$ px create volume myvolume --size=3 --shared
4. Create shared volume called "myvolume" with size as 2GiB and replicas set to 3:
	$ px create volume myvolume --size=3 --shared --replicas=3
5. Create volume called "myvolume" with label as "access=slow" and size as 3 GiB:
	$ px create volume myvolume --size=3 --labels 'access=slow'`,
/*
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Must supply a something for role")
			}
			return nil
		},
*/
		RunE: createRoleExec,
	}
})

var _ = commander.RegisterCommandInit(func() {
	cmd.CreateAddCommand(createRoleCmd)

	createRoleCmd.Flags().StringVar(&crOpts.roleconf, "role-config", "", "Required role json file'")
	cobra.MarkFlagRequired(createRoleCmd.Flags(), "role-config")
})

func createRoleExec(c *cobra.Command, args []string) error {

util.Printf("In role create: %s \n", crOpts.roleconf)
	jsonFile, err := os.Open(crOpts.roleconf)
	if err != nil {
		fmt.Primntln(err)
	}
	defer jsonFile.Close()

/*
	ctx, conn, err := portworx.PxConnectDefault()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Get name
	cvOpts.req.Name = args[0]

	// Get labels
	if len(cvOpts.labelsAsString) != 0 {
		var err error
		cvOpts.req.Labels, err = util.CommaStringToStringMap(cvOpts.labelsAsString)
		if err != nil {
			return fmt.Errorf("Failed to parse labels: %v\n", err)
		}
	}

	// Convert size to bytes in uint64
	cvOpts.req.Spec.Size = uint64(cvOpts.sizeInGi) * uint64(cmd.Gi)

	// Add fs to request
	switch {
	case cvOpts.filesystemAsString == "ext4":
		cvOpts.req.Spec.Format = api.FSType_FS_TYPE_EXT4
	case cvOpts.filesystemAsString == "none":
		cvOpts.req.Spec.Format = api.FSType_FS_TYPE_NONE
	default:
		return fmt.Errorf("Error: --fs valid values are [none, ext4]\n")
	}

	// Update default EXT4, if it fs is 'none' and shared volume
	if cvOpts.req.Spec.Format == api.FSType_FS_TYPE_NONE && cvOpts.req.Spec.Shared {
		cvOpts.req.Spec.Format = api.FSType_FS_TYPE_EXT4
	}

	// setting replica set nodes if provided
	if len(cvOpts.replicaSet) != 0 {
		cvOpts.req.Spec.ReplicaSet = &api.ReplicaSet{
			Nodes: cvOpts.replicaSet,
		}
	}

	// setting IO profile if provided.
	if len(cvOpts.IoProfile) > 0 {
		switch cvOpts.IoProfile {
		case "db":
			cvOpts.req.Spec.IoProfile = api.IoProfile_IO_PROFILE_DB
		case "cms":
			cvOpts.req.Spec.IoProfile = api.IoProfile_IO_PROFILE_CMS
		case "db_remote":
			cvOpts.req.Spec.IoProfile = api.IoProfile_IO_PROFILE_DB_REMOTE
		case "sync_shared":
			cvOpts.req.Spec.IoProfile = api.IoProfile_IO_PROFILE_SYNC_SHARED
		default:
			flagError := errors.New("Invalid IO profile")
			return flagError
		}
	}

	// Send request
	volumes := api.NewOpenStorageVolumeClient(conn)
	resp, err := volumes.Create(ctx, cvOpts.req)
	if err != nil {
		return util.PxErrorMessage(err, "Failed to create volume")
	}

	// Show user information
	msg := fmt.Sprintf("Volume %s created with id %s\n",
		cvOpts.req.GetName(),
		resp.GetVolumeId())

	formattedOut := &util.DefaultFormatOutput{
		Cmd:  "create volume",
		Desc: msg,
		Id:   []string{resp.GetVolumeId()},
	}
	util.PrintFormatted(formattedOut)
*/
	return nil
}
