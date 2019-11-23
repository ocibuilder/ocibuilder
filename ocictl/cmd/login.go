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
	"github.com/ocibuilder/ocibuilder/pkg"
	"io"

	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/ocictl/pkg/utils"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/buildah"
	"github.com/ocibuilder/ocibuilder/pkg/docker"
	"github.com/spf13/cobra"
)

const loginDesc = `
This command logs into all registries that have been defined in the specification. You can login with a number of different credentials.
These can be plain, taken from environment variables or kubernetes secrets.
`

type loginCmd struct {
	out     io.Writer
	path    string
	builder string
	debug   bool
}

func newLoginCmd(out io.Writer) *cobra.Command {
	lc := &loginCmd{out: out}
	cmd := &cobra.Command{
		Use:   "login",
		Short: "logs into all registries defined in the specifcation.",
		Long:  loginDesc,
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.run(args)
		},
	}
	f := cmd.Flags()
	f.StringVarP(&lc.path, "path", "p", "", "Path to your ocibuilder.yaml or login.yaml. By default will look in the current working directory")
	f.StringVarP(&lc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targetted image puller. By default the builder is docker.")
	f.BoolVarP(&lc.debug, "debug", "d", false, "Turn on debug logging")
	return cmd
}

func (l *loginCmd) run(args []string) error {
	ociBuilderSpec := v1alpha1.OCIBuilderSpec{
		Daemon: true,
	}
	logger := common.GetLogger(l.debug)

	reader := pkg.Reader{
		Logger: logger,
	}
	if err := reader.Read(&ociBuilderSpec, "", l.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	// Prioritise builder passed in as argument, default builder is docker
	builder := l.builder
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

			res, err := d.Login(ociBuilderSpec)
			if err != nil {
				log.WithError(err).Errorln("failed to login to registry")
				return err
			}

			log.WithField("responses", len(res)).Debugln("received responses and running login")
			for idx, loginResponse := range res {
				if err := utils.Output(loginResponse); err != nil {
					return err
				}
				log.WithField("response", idx).Debugln("response has finished executing")
			}
			log.Infoln("docker login completed")
		}

	case v1alpha1.BuildahFramework:
		{
			b := buildah.Buildah{
				Logger: logger,
			}
			log := b.Logger

			res, err := b.Login(ociBuilderSpec)
			if err != nil {
				log.WithError(err).Errorln("failed to login to registry")
				return err
			}

			log.WithField("responses", len(res)).Debugln("received responses and running login")
			for idx, loginResponse := range res {
				if err := utils.Output(loginResponse); err != nil {
					return err
				}
				if err := b.Wait(idx); err != nil {
					return err
				}
				log.WithField("response", idx).Debugln("response has finished executing")
			}
			log.Infoln("buildah login complete")
		}

	default:
		{
			return errors.New("invalid builder specified, try --builder=docker or --builder=buildah")
		}
	}
	return nil
}
