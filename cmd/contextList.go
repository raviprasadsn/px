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
package cmd

import (
	"github.com/portworx/px/pkg/contextconfig"
	"github.com/portworx/px/pkg/util"
	"github.com/spf13/cobra"
)

// contextGetCmd represents the contextGet command
var contextListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"contexts", "ctx"},
	Short:   "List all context configurations",
	Long: `List all context configurations
px get context`,
	RunE: contextListExec,
}

func init() {
	contextCmd.AddCommand(contextListCmd)
}

func contextListExec(cmd *cobra.Command, args []string) error {
	contextManager, err := contextconfig.NewContextManager(GetConfigFile())
	if err != nil {
		return util.PxErrorMessagef(err, "Failed to get context configuration at location %s",
			GetConfigFile())
	}
	cfg := contextManager.GetAll()

	// add extra information
	cfg = contextconfig.AddClaimsInfo(cfg)
	cfg = contextconfig.MarkInvalidTokens(cfg)

	// Print out config
	util.PrintYaml(cfg)

	return nil
}
