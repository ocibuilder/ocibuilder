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
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/buildah"
	"github.com/ocibuilder/ocibuilder/pkg/docker"
	"github.com/ocibuilder/ocibuilder/pkg/oci"
	"github.com/ocibuilder/ocibuilder/pkg/read"
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
	var cli v1alpha1.BuilderClient
	logger := common.GetLogger(l.debug)
	reader := read.Reader{Logger: logger}
	ociBuilderSpec := v1alpha1.OCIBuilderSpec{Daemon: true}

	if err := reader.Read(&ociBuilderSpec, "", l.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	// Prioritise builder passed in as argument, default builder is docker
	builderType := l.builder
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

	res := make(chan v1alpha1.OCILoginResponse)
	errChan := make(chan error)
	finished := make(chan bool)

	defer func() {
		close(res)
		close(errChan)
		close(finished)
	}()

	go builder.Login(ociBuilderSpec, res, errChan, finished)

	for {
		select {

		case err := <-errChan:
			{
				logger.WithError(err).Errorln("error received from error channel whilst logging in")
				return err
			}

		case loginResponse := <-res:
			{
				logger.Infoln("executing login step")
				//TODO: make this output nicer
				fmt.Println(loginResponse)
			}

		case <-finished:
			{
				logger.Infoln("all login steps complete successfully")
				return nil
			}
		}
	}

}
