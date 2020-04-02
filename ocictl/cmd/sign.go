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
	"io/ioutil"

	"github.com/beval/beval/pkg/apis/beval/v1alpha1"
	"github.com/beval/beval/pkg/crypto"
	"github.com/beval/beval/pkg/docker"
	"github.com/beval/beval/pkg/oci"
	"github.com/beval/beval/pkg/read"
	"github.com/beval/beval/pkg/util"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

type signCmd struct {
	out    io.Writer
	path   string
	push   bool
	debug  bool
	name   string
	output string
}

func newSignCmd(out io.Writer) *cobra.Command {
	sc := &signCmd{out: out}
	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Sign signs a container image ID and optionally pushes the attestation to a metadata store",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sc.run(args)
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&sc.debug, "debug", "d", false, "Turn on debug logging")
	f.StringVarP(&sc.path, "path", "p", ".", "Path to your beval.yaml file")
	f.BoolVar(&sc.push, "push", false, "Push to specified metadata store")
	f.StringVarP(&sc.name, "name", "n", "", "The image name to sign")
	f.StringVarP(&sc.output, "output", "o", "", "Filepath to output image signature to")

	return cmd
}

func (sc *signCmd) run(args []string) error {
	logger := util.GetLogger(sc.debug)
	reader := read.Reader{Logger: logger}
	bevalSpec := v1alpha1.bevalSpec{Daemon: true}

	if err := reader.Read(&bevalSpec, "", sc.path); err != nil {
		log.WithError(err).Errorln("failed to read spec")
		return err
	}

	apiClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.WithError(err).Errorln("failed to fetch docker api client")
		return err
	}

	cli := docker.Client{
		APIClient: apiClient,
		Logger:    logger,
	}

	if sc.push {
		bevalSpec.Metadata.Data = []v1alpha1.MetadataType{"attestation"}
		mw := oci.NewMetadataWriter(logger, bevalSpec.Metadata)
		if err := mw.ParseMetadata(sc.name, cli, &v1alpha1.BuildProvenance{
			Name: sc.name,
		}); err != nil {
			return err
		}
		if err := mw.Write(); err != nil {
			return err
		}

		logger.Infoln("image attestation has been pushed to metadata store")
		return nil
	}

	inspectResponse, err := cli.ImageInspect(sc.name)
	if err != nil || inspectResponse.ID == "" {
		logger.Errorln("error in inspecting image or no response ID returned - cannot push metadata")
		return err
	}

	privKey, pubKey, err := crypto.ValidateKeysPacket(bevalSpec.Metadata.Key)
	if err != nil {
		return err
	}

	e := crypto.CreateEntityFromKeys(privKey, pubKey)
	_, sig, err := crypto.SignDigest(inspectResponse.ID, bevalSpec.Metadata.Key.Passphrase, e)
	if err != nil {
		return err
	}

	logger.WithField("imageId", inspectResponse.ID).Infoln("successfully signed image")
	if sc.output != "" {
		logger.WithField("path", sc.output).Infoln("outputting signature to file")
		if err := ioutil.WriteFile(sc.output, []byte(sig), 0644); err != nil {
			return err
		}
		return nil
	}

	logger.Infoln("image signature")
	fmt.Println(sig)
	return nil
}
