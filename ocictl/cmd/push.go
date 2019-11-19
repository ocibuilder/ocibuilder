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
	"github.com/ocibuilder/ocibuilder/pkg/oci"
	"github.com/spf13/cobra"
)

const pushDesc = `
This command pushes all specified images in the push specification to a corresponding registry. You can list many images to push
to many registries.

In order to complete a push to a repository, both the login and push specifications need to be filled in. A push is run with the
authentication passed in the login spec.

The registry, image and tag are used to create a full qualified image path

e.g. my-image-registry.docker.com:1111/myimage/cool-image:0.0.1
`

type pushCmd struct {
	out     io.Writer
	path    string
	builder string
	debug   bool
}

func newPushCmd(out io.Writer) *cobra.Command {
	pc := &pushCmd{out: out}
	cmd := &cobra.Command{
		Use:   "push",
		Short: "pushes container images to one or multiple image registries.",
		Long:  pushDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pc.run(args)
		},
	}
	f := cmd.Flags()
	f.StringVarP(&pc.path, "path", "p", "", "Path to your ocibuilder.yaml or push.yaml. By default will look in the current working directory")
	f.StringVarP(&pc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targetted image builder. By default the builder is docker.")
	f.BoolVarP(&pc.debug, "debug", "d", false, "Turn on debug logging")
	return cmd
}

func (p *pushCmd) run(args []string) error {
	var cli v1alpha1.BuilderClient
	logger := common.GetLogger(p.debug)
	reader := common.Reader{Logger: logger}
	ociBuilderSpec := v1alpha1.OCIBuilderSpec{Daemon: true}

	if err := reader.Read(&ociBuilderSpec, "", p.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	// Prioritise builder passed in as argument, default builder is docker
	builderType := p.builder
	if !ociBuilderSpec.Daemon {
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
		}

	case v1alpha1.BuildahFramework:
		{
			cli = buildah.Client{
				Logger: logger,
			}
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

	res := make(chan v1alpha1.OCIPushResponse)
	errChan := make(chan error)
	go builder.Push(ociBuilderSpec, res, errChan)

	select {

	case err := <-errChan:
		{
			return err
		}

	case pushResponse := <-res:
		{
			if builderType == "docker" {
				if err := utils.OutputJson(pushResponse.Body); err != nil {
					return err
				}
			} else {
				if err := utils.Output(pushResponse.Body, pushResponse.Stderr); err != nil {
					return err
				}
			}
		}

	}
	return nil
}
