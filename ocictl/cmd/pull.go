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

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/ocictl/pkg/utils"
	"github.com/ocibuilder/ocibuilder/pkg/oci"
	"github.com/ocibuilder/ocibuilder/pkg/read"
	"github.com/ocibuilder/ocibuilder/pkg/types"
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
	f.StringVarP(&pc.path, "path", "p", "", "Path to your ocibuilder.yaml. By default will look in the current working directory")
	f.StringVarP(&pc.builder, "builder", "b", "docker", "Choose either docker and buildah as the targeted image builder. By default the builder is docker.")
	f.BoolVarP(&pc.debug, "debug", "d", false, "Turn on debug logging")
	return cmd
}

func (p *pullCmd) run(args []string) error {
	logger := common.GetLogger(p.debug)
	reader := read.Reader{Logger: logger}

	ociBuilderSpec, err := reader.Read("", p.path)
	if err != nil {
		return err
	}

	client, err := utils.GetClient(p.builder, logger)
	if err != nil {
		return err
	}

	ociBuilderSpec.Daemon = utils.HasDaemon(p.builder)

	builder := oci.Builder{
		Logger: logger,
		Client: client,
	}

	res := make(chan types.OCIPullResponse)
	errChan := make(chan error)
	finished := make(chan bool)

	defer func() {
		close(res)
		close(errChan)
		close(finished)
	}()

	go builder.Pull(ociBuilderSpec, p.name, res, errChan, finished)

	for {
		select {

		case err := <-errChan:
			{
				if err != nil {
					logger.WithError(err).Errorln("error received from error channel whilst pulling")
					return err
				}
			}

		case pullResponse := <-res:
			{
				logger.Infoln("executing pull step")
				if p.builder == "docker" {
					if err := utils.OutputJson(pullResponse.Body); err != nil {
						return err
					}
				} else {
					if err := utils.Output(pullResponse.Body, pullResponse.Stderr); err != nil {
						return err
					}
				}
				logger.Infoln("pull step complete")
			}

		case <-finished:
			{
				logger.Infoln("all pull steps complete successfully")
				return nil
			}

		}
	}
}
