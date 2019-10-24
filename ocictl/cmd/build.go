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
	"errors"
	"io"

	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/ocictl/pkg/utils"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/buildah"
	"github.com/ocibuilder/ocibuilder/pkg/docker"
	"github.com/spf13/cobra"
)

const buildDesc = `
This command runs an image build with the specification defined in your projects spec.yaml file.
It can run a build in both docker and buildah varieties.
`

type buildCmd struct {
	out io.Writer
	name string
	path string
	builder string
	overlay string
	storageDriver string
	debug bool
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
	f.StringVarP(&bc.name, "name", "n", "", "Specify the name of your build or defined in spec.yaml")
	f.StringVarP(&bc.path, "path", "p", "", "Path to your spec.yaml or build.yaml. By default will look in the current working directory")
	f.StringVarP(&bc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targetted image builder. By default the builder is docker.")
	f.BoolVarP(&bc.debug, "debug", "d", false, "Turn on debug logging")
	f.StringVarP(&bc.overlay, "overlay", "o", "", "Path to your overlay.yaml file")
	f.StringVarP(&bc.storageDriver, "storage-driver", "s", "", "Storage-driver for Buildah. vfs enables the use of buildah within an unprivileged container. By default the storage driver is overlay")

	return cmd
}

func (b *buildCmd) run(args []string) error {
	ociBuilderSpec := v1alpha1.OCIBuilderSpec{}
	if err := common.Read(&ociBuilderSpec, b.overlay, b.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	switch v1alpha1.Framework(b.builder) {

	case v1alpha1.DockerFramework:
		{
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.WithError(err).Errorln("failed to fetch docker client")
				return err
			}

			d := docker.Docker{
				Client:      cli,
				Logger:      common.GetLogger(b.debug),
			}
			res, err := d.Build(ociBuilderSpec)
			if err != nil {
				return err
			}

			for idx, imageResponse := range res {
				log.WithField("step: ", idx).Infoln("running build step")

				if imageResponse == nil {
					return errors.New("no response received from daemon - check if docker is installed and running")
				}

				err := utils.OutputJson(imageResponse)
				if err != nil {
					return err
				}
			}
			d.Clean()
		}

	case v1alpha1.BuildahFramework:
		{
			b := buildah.Buildah{
				Logger: common.GetLogger(b.debug),
				StorageDriver: b.storageDriver,
			}

			res, err := b.Build(ociBuilderSpec)
			if err != nil {
				log.WithError(err).Errorln("error executing build on ocibuilder spec")
				return err
			}

			for idx, imageResponse := range res {
				log.WithField("step: ", idx).Infoln("running build step")
				if err := utils.Output(imageResponse); err != nil {
					return err
				}
				if err := b.Wait(); err != nil {
					return err
				}
			}
			b.Clean()
		}

	default:
		{
			return errors.New("invalid builder specified, try --builder=docker or --builder=buildah")
		}

	}

	return nil
}
