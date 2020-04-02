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

	"github.com/beval/beval/bevalctl/pkg/utils"
	"github.com/beval/beval/pkg/apis/beval/v1alpha1"
	"github.com/beval/beval/pkg/buildah"
	"github.com/beval/beval/pkg/docker"
	"github.com/beval/beval/pkg/oci"
	"github.com/beval/beval/pkg/read"
	"github.com/beval/beval/pkg/util"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const buildDesc = `
This command runs an image build with the specification defined in your projects beval.yaml file.
It can run a build in both docker and buildah varieties.
`

type buildCmd struct {
	out     io.Writer
	name    string
	path    string
	builder string
	overlay string
	debug   bool
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
	f.StringVarP(&bc.name, "name", "n", "", "Specify the name of your build or defined in beval.yaml")
	f.StringVarP(&bc.path, "path", "p", "", "Path to your beval.yaml or build.yaml. By default will look in the current working directory")
	f.StringVarP(&bc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targeted image builder. By default the builder is docker.")
	f.BoolVarP(&bc.debug, "debug", "d", false, "Turn on debug logging")
	f.StringVarP(&bc.overlay, "overlay", "o", "", "Path to your overlay.yaml file")

	return cmd
}

func (b *buildCmd) run(args []string) error {
	var cli v1alpha1.BuilderClient
	logger := util.GetLogger(b.debug)
	reader := read.Reader{Logger: logger}
	bevalSpec := v1alpha1.bevalSpec{Daemon: true}

	if err := reader.Read(&bevalSpec, b.overlay, b.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	// Prioritise builder passed in as argument, default builder is docker
	builderType := b.builder
	if !bevalSpec.Daemon {
		builderType = "buildah"
	}

	switch v1alpha1.Framework(builderType) {

	case v1alpha1.DockerFramework:
		{
			apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.WithError(err).Errorln("failed to fetch docker api client")
				return err
			}

			cli = docker.Client{
				APIClient: apiClient,
				Logger:    logger,
			}

			bevalSpec.Daemon = true
		}

	case v1alpha1.BuildahFramework:
		{
			cli = buildah.Client{
				Logger: logger,
			}

			bevalSpec.Daemon = false
		}

	default:
		{
			return errors.New("invalid builder specified, try --builder=docker or --builder=buildah")
		}

	}

	builder := oci.Builder{
		Logger: logger,
		Client: cli,
	}

	res := make(chan v1alpha1.OCIBuildResponse)
	errChan := make(chan error)
	finished := make(chan bool)

	defer func() {
		close(res)
		close(errChan)
	}()

	go builder.Build(bevalSpec, res, errChan, finished)

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
				if builderType == "docker" {
					if err := utils.OutputJson(buildResponse.Body); err != nil {
						return err
					}
				} else {
					if err := utils.Output(buildResponse.Body, buildResponse.Stderr); err != nil {
						return err
					}
				}
				buildResponse.Finished = true
				res <- buildResponse
				logger.Infoln("build step complete")
			}

		case <-finished:
			{
				logger.Infoln("all build steps complete")
				close(finished)
				return nil
			}

		}
	}

}
