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
	"github.com/ocibuilder/ocibuilder/provenance"
	"github.com/spf13/cobra"
	"io"
)

type versionCmd struct {
	out 	io.Writer
	verbose bool
}

func newVersionCmd(out io.Writer) *cobra.Command {
	vc := &versionCmd{out: out}
	cmd := &cobra.Command{
		Use: "version",
		Short: "prints the version of ocibuilder",
		Run: func(cmd *cobra.Command, args []string) {
			vc.run()
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&vc.verbose, "verbose", "v", false, "Detailed output of ocictl version")

	return cmd
}

func (v *versionCmd) run() {
	prov := provenance.GetProvenance()

	if v.verbose {
		prov.PrintVerbose(v.out)
		return
	}

	prov.Print(v.out)
	return
}
