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

	"github.com/blackrock/ocibuilder/common"
	"github.com/blackrock/ocibuilder/ocictl/pkg/utils"
	"github.com/blackrock/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/blackrock/ocibuilder/pkg/buildah"
	"github.com/blackrock/ocibuilder/pkg/docker"
	"github.com/docker/docker/client"
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
	f.StringVarP(&pc.path, "path", "p", "", "Path to your spec.yaml or push.yaml. By default will look in the current working directory")
	f.StringVarP(&pc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targetted image builder. By default the builder is docker.")
	f.BoolVarP(&pc.debug, "debug", "d", false, "Turn on debug logging")
	return cmd
}

func (p *pushCmd) run(args []string) error {
	ociBuilderSpec := v1alpha1.OCIBuilderSpec{}
	if err := common.Read(&ociBuilderSpec, "", p.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	switch v1alpha1.Framework(p.builder) {

	case v1alpha1.DockerFramework:
		{
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.WithError(err).Errorln("failed to fetch docker client")
				return err
			}
			d := docker.Docker{
				Client: cli,
				Logger: common.GetLogger(p.debug),
			}
			res, err := d.Push(ociBuilderSpec)
			if err != nil {
				return err
			}

			for idx, imageResponse := range res {
				log.WithField("step: ", idx).Infoln("running push step")
				err := utils.OutputJson(imageResponse)
				if err != nil {
					return err
				}
			}
			log.Infoln("docker push complete")
		}

	case v1alpha1.BuildahFramework:
		{
			b := buildah.Buildah{
				Logger: common.GetLogger(p.debug),
			}
			res, err := b.Push(ociBuilderSpec)
			if err != nil {
				return err
			}

			for idx, imageResponse := range res {
				log.WithField("step: ", idx).Infoln("running push step")
				if err := utils.Output(imageResponse); err != nil {
					return err
				}
			}
			log.Infoln("buildah push complete")
		}

	default:
		{
			return errors.New("invalid builder specified, try --builder=docker or --builder=buildah")
		}

	}
	return nil
}
