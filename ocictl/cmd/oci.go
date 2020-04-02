/*
Copyright 2019 BlackRock, Inc.

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
	"github.com/beval/beval/pkg/util"
	"github.com/spf13/cobra"
)

var log = util.GetLogger(false)

const rootDesc = `
              _      __  __
  ____  _____(_)____/ /_/ /
 / __ \/ ___/ / ___/ __/ / 
/ /_/ / /__/ / /__/ /_/ /  
\____/\___/_/\___/\__/_/   
                           
                         
The ocictl is a tool offered by beval for pulling, building and
pushing your images using your specified build framework.
`

// NewRootCmd is the root command for the ocictl
func NewRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ocictl",
		Short: "The ocictl provided by beval",
		Long:  rootDesc,
	}
	flags := cmd.PersistentFlags()
	out := cmd.OutOrStdout()

	cmd.AddCommand(
		newBuildCmd(out),
		newLoginCmd(out),
		newPullCmd(out),
		newPushCmd(out),
		newVersionCmd(out),
		newInitCmd(out),
		newSignCmd(out),
	)

	flags.Parse(args) //nolint
	return cmd
}
