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
	f.StringVarP(&bc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targetted image builder. By default the builder is docker.")
	f.BoolVarP(&bc.debug, "debug", "d", false, "Turn on debug logging")
	f.StringVarP(&bc.overlay, "overlay", "o", "", "Path to your overlay.yaml file")
	f.StringVarP(&bc.storageDriver, "storage-driver", "s", "overlay", "Storage-driver for Buildah. vfs enables the use of buildah within an unprivileged container. By default the storage driver is overlay")

	return cmd
}

func (b *buildCmd) run(args []string) error {
	logger := common.GetLogger(b.debug)
	ociBuilderSpec := v1alpha1.OCIBuilderSpec{
		Daemon: true,
	}

	reader := common.Reader{
		Logger: logger,
	}

	if err := reader.Read(&ociBuilderSpec, b.overlay, b.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	// Prioritise builder passed in as argument, default builder is docker
	builder := b.builder
	if !ociBuilderSpec.Daemon {
		builder = "buildah"
	}

	switch v1alpha1.Framework(builder) {

	case v1alpha1.DockerFramework:
		{
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.WithError(err).Errorln("failed to fetch docker client")
				return err
			}

			d := docker.Docker{
				Client: cli,
				Logger: logger,
			}
			log := d.Logger

			res, err := d.Build(ociBuilderSpec)
			if err != nil {
				return err
			}

			log.WithField("responses", len(res)).Debugln("received responses and running build")
			for idx, imageResponse := range res {
				log.WithField("step: ", idx).Infoln("running build step")

				if imageResponse == nil {
					return errors.New("no response received from daemon - check if docker is installed and running")
				}

				if err := utils.OutputJson(imageResponse); err != nil {
					return err
				}
				log.WithField("response", idx).Debugln("response has finished executing")
			}
			log.Debugln("running build file cleanup")
			d.Clean()
			log.Infoln("docker build complete")
		}

	case v1alpha1.BuildahFramework:
		{
			b := buildah.Buildah{
				Logger:        logger,
				StorageDriver: b.storageDriver,
			}
			log := b.Logger

			res, err := b.Build(ociBuilderSpec)
			if err != nil {
				log.WithError(err).Errorln("error executing build on ocibuilder spec")
				return err
			}

			log.WithField("responses", len(res)).Debugln("received responses and running build")
			for idx, imageResponse := range res {
				log.WithField("step: ", idx).Infoln("running build step")
				if err := utils.Output(imageResponse); err != nil {
					return err
				}
				if err := b.Wait(idx); err != nil {
					return err
				}
				log.WithField("response", idx).Debugln("response has finished executing")
			}
			log.Debugln("running build file cleanup")
			b.Clean()
			log.Infoln("buildah build complete")
		}

	default:
		{
			return errors.New("invalid builder specified, try --builder=docker or --builder=buildah")
		}

	}

	return nil
}
