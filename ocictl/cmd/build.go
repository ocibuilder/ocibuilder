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

	"github.com/ocibuilder/ocibuilder/pkg/types"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/ocictl/pkg/utils"
	"github.com/ocibuilder/ocibuilder/pkg/oci"
	"github.com/ocibuilder/ocibuilder/pkg/read"
	"github.com/spf13/cobra"
)

const buildDesc = `
This command runs an image build with the specification defined in your projects ocibuilder.yaml file.
It can run a build in both docker and buildah varieties.
`

type buildCmd struct {
	out           io.Writer
	name          string
	path          string
	builder       string
	overlay       string
	storageDriver string
	debug         bool
}

func newBuildCmd(out io.Writer) *cobra.Command {
	bc := &buildCmd{out: out}
	cmd := &cobra.Command{
		Use:   "build",
		Short: "builds an oci compliant image using either docker or buildah",
		Long:  buildDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return bc.run(args)
		},
	}
	f := cmd.Flags()
	f.StringVarP(&bc.name, "name", "n", "", "Specify the name of your build or defined in ocibuilder.yaml")
	f.StringVarP(&bc.path, "path", "p", "", "Path to your ocibuilder.yaml or build.yaml. By default will look in the current working directory")
	f.StringVarP(&bc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targeted image builder. By default the builder is docker.")
	f.BoolVarP(&bc.debug, "debug", "d", false, "Turn on debug logging")
	f.StringVarP(&bc.overlay, "overlay", "o", "", "Path to your overlay.yaml file")
	f.StringVarP(&bc.storageDriver, "storage-driver", "s", "overlay", "Storage-driver for Buildah. vfs enables the use of buildah within an unprivileged container. By default the storage driver is overlay")

	return cmd
}

func (b *buildCmd) run(args []string) error {
	logger := common.GetLogger(b.debug)
	reader := read.Reader{Logger: logger}

	ociBuilderSpec, err := reader.Read(b.overlay, b.path)
	if err != nil {
		return err
	}

	client, err := utils.GetClient(b.builder, logger)
	if err != nil {
		return err
	}

	// TODO: Instead of having user mention the type of the builder, it makes sense to derive builder value just from Daemon -> true or false
	ociBuilderSpec.Daemon = utils.HasDaemon(b.builder)

	builder := oci.Builder{
		Logger: logger,
		Client: client,
	}

	res := make(chan types.OCIBuildResponse)
	errChan := make(chan error)
	finished := make(chan bool)

	defer func() {
		close(res)
		close(errChan)
		close(finished)
	}()

	go builder.Build(ociBuilderSpec, res, errChan, finished)

	for {
		select {

		case err := <-errChan:
			{
				if err != nil {
					logger.WithError(err).Errorln("error received from error channel whilst building")
					builder.Clean()
					return err
				}
			}

		case buildResponse := <-res:
			{
				logger.Infoln("executing build step")
				if b.builder == "docker" {
					if err := utils.OutputJson(buildResponse.Body); err != nil {
						return err
					}
				} else {
					if err := utils.Output(buildResponse.Body, buildResponse.Stderr); err != nil {
						return err
					}
				}
				logger.Infoln("build step complete")
			}

		case <-finished:
			{
				logger.Infoln("all build steps complete")
				return nil
			}

		}
	}

}
