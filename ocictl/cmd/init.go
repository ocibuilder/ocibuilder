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
	"io"

	"github.com/ocibuilder/ocibuilder/pkg/init"

	"github.com/gobuffalo/packr"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/spf13/cobra"
)

type initCmd struct {
	out        io.Writer
	dry        bool
	debug      bool
	fromDocker string
}

func newInitCmd(out io.Writer) *cobra.Command {
	ic := &initCmd{out: out}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialises a template ocibuilder.yaml file for ocibuilder",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ic.run(args)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&ic.dry, "dry", false, "Run a dry spec generation which is outputted to the terminal")
	f.StringVarP(&ic.fromDocker, "from-docker", "-f", "", "Generate an ocibuilder build spec from a Dockerfile. Expects a path to a Dockerfile")
	f.BoolVarP(&ic.debug, "debug", "d", false, "Turn on debug logging")

	return cmd
}

func (i *initCmd) run(args []string) error {
	initializer := init.Initializer{
		Box:    packr.NewBox("../../templates/spec"),
		Dry:    i.dry,
		Logger: common.GetLogger(i.debug),
	}

	if i.fromDocker != "" {
		if err := initializer.FromDocker(i.fromDocker); err != nil {
			return err
		}
	}

	if err := initializer.Basic(); err != nil {
		return err
	}

	return nil
}
