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

	"github.com/gobuffalo/packr"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/initialize"
	"github.com/spf13/cobra"
)

type initCmd struct {
	out   io.Writer
	dry   bool
	debug bool
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
	f.BoolVarP(&ic.debug, "debug", "d", false, "Turn on debug logging")

	cmd.AddCommand(newFromDockerCmd(out))
	return cmd
}

func (i *initCmd) run(args []string) error {
	initializer := initialize.Initializer{
		Box:    packr.NewBox("../../templates/spec"),
		Dry:    i.dry,
		Logger: common.GetLogger(i.debug),
	}

	if err := initializer.Basic(); err != nil {
		return err
	}

	return nil
}

type fromDockerCmd struct {
	out       io.Writer
	debug     bool
	dry       bool
	imageName string
	path      string
}

func newFromDockerCmd(out io.Writer) *cobra.Command {
	fd := &fromDockerCmd{out: out}
	cmd := &cobra.Command{
		Use:   "from-docker",
		Short: "Initialises a template ocibuilder.yaml file from a docker file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fd.run(args)
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&fd.debug, "debug", "d", false, "Turn on debug logging")
	f.BoolVar(&fd.dry, "dry", false, "Run a dry spec generation which is outputted to the terminal")
	f.StringVarP(&fd.path, "path", "p", "", "Path to your Dockerfile")
	f.StringVarP(&fd.imageName, "tag", "t", "", "The name and tag for your image")

	return cmd
}

func (i *fromDockerCmd) run(args []string) error {
	initializer := initialize.Initializer{
		Dry:    i.dry,
		Logger: common.GetLogger(i.debug),
	}

	if err := initializer.FromDocker(i.imageName, i.path); err != nil {
		return err
	}

	return nil
}
