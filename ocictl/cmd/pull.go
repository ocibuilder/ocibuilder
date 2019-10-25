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

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/ocictl/pkg/utils"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/buildah"
	"github.com/ocibuilder/ocibuilder/pkg/docker"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

const pullDesc = `
This command pulls an image that you have passed in by name. The name should include the path to the image but not
the image registry itself.

e.g. myimage/cool-image:0.0.1

The pull command looks to pull from any registries that have been specified in the login specification. Once the image has
been found in any of the specified registries, a pull is executed.
`

type pullCmd struct {
	out     io.Writer
	name    string
	path    string
	builder string
	debug   bool
}

func newPullCmd(out io.Writer) *cobra.Command {
	pc := &pullCmd{out: out}
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "pulls an image passed in with the name flag",
		Long:  pullDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return pc.run(args)
		},
	}
	f := cmd.Flags()
	f.StringVarP(&pc.name, "name", "i", "", "Specify the name of the image you want to pull")
	f.StringVarP(&pc.path, "path", "p", "", "Path to your spec.yaml. By default will look in the current working directory")
	f.StringVarP(&pc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targetted image puller. By default the builder is docker.")
	f.BoolVarP(&pc.debug, "debug", "d", false, "Turn on debug logging")
	return cmd
}

func (p *pullCmd) run(args []string) error {
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
			log := d.Logger

			res, err := d.Pull(ociBuilderSpec, p.name)
			if err != nil {
				return err
			}

			log.WithField("responses", len(res)).Debugln("received responses and running pull")
			for idx, imageResponse := range res {
				log.WithField("step: ", idx).Infoln("running pull step")

				if err := utils.OutputJson(imageResponse); err != nil {
					return err
				}
				log.WithField("response", idx).Debugln("response has finished executing")
			}
			log.Infoln("docker pull completed")
		}

	case v1alpha1.BuildahFramework:
		{
			b := buildah.Buildah{
				Logger: common.GetLogger(p.debug),
			}
			log := b.Logger

			res, err := b.Pull(ociBuilderSpec, p.name)
			if err != nil {
				return err
			}

			log.WithField("responses", len(res)).Debugln("received responses and running pull")
			for idx, imageResponse := range res {
				log.WithField("step: ", idx).Infoln("running pull step")
				if err := utils.Output(imageResponse); err != nil {
					return err
				}
				if err := b.Wait(idx); err != nil {
					return err
				}
				log.WithField("response", idx).Debugln("response has finished executing")
			}
			log.Infoln("buildah pull complete")
		}

	default:
		{
			return errors.New("invalid builder specified, try --builder=docker or --builder=buildah")
		}
	}

	return nil
}
