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
	"fmt"
	"io"

	"github.com/docker/docker/client"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/crypto"
	"github.com/ocibuilder/ocibuilder/pkg/docker"
	"github.com/ocibuilder/ocibuilder/pkg/oci"
	"github.com/ocibuilder/ocibuilder/pkg/read"
	"github.com/ocibuilder/ocibuilder/pkg/util"
	"github.com/spf13/cobra"
)

type signCmd struct {
	out   io.Writer
	path  string
	push  bool
	debug bool
	name  string
}

func newSignCmd(out io.Writer) *cobra.Command {
	sc := &signCmd{out: out}
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign signs a container image ID and optionally pushes attestation to a metadata store",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sc.run(args)
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&sc.debug, "debug", "d", false, "Turn on debug logging")
	f.StringVarP(&sc.path, "path", "p", ".", "Path to your ocibuilder.yaml file")
	f.BoolVar(&sc.push, "push", false, "Push to specified metadata store")
	f.StringVarP(&sc.name, "name", "n", "", "The image name to sign")

	return cmd
}

func (sc *signCmd) run(args []string) error {
	logger := util.GetLogger(sc.debug)
	reader := read.Reader{Logger: logger}
	ociBuilderSpec := v1alpha1.OCIBuilderSpec{Daemon: true}

	if err := reader.Read(&ociBuilderSpec, "", sc.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.WithError(err).Errorln("failed to fetch docker api client")
		return err
	}

	cli := docker.Client{
		APIClient: apiClient,
		Logger:    logger,
	}

	if sc.push {
		mw := oci.NewMetadataWriter(log, ociBuilderSpec.Metadata)
		if err := mw.ParseMetadata(sc.name, cli, v1alpha1.BuildProvenance{}); err != nil {
			return err
		}

		logger.Infoln("image attestation has been pushed to metadata store")
	} else {

		inspectResponse, err := cli.ImageInspect(sc.name)
		if err != nil || inspectResponse.ID == "" {
			log.Errorln("error in inspecting image or no response ID returned - cannot push metadata")
			return err
		}

		privKey, pubKey, err := crypto.ValidateKeysPacket(ociBuilderSpec.Metadata.Key)
		if err != nil {
			return err
		}

		e := crypto.CreateEntityFromKeys(privKey, pubKey)
		id, sig, err := crypto.SignDigest(inspectResponse.ID, ociBuilderSpec.Metadata.Key.Passphrase, e)
		if err != nil {
			return err
		}

		logger.WithField("imageId", id).Infoln("successfully signed image")
		logger.Infoln("image signature")
		fmt.Println(sig)
	}

	return nil
}
