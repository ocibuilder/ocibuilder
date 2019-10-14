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
	"io/ioutil"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/gobuffalo/packr"
	"github.com/spf13/cobra"
)

type initCmd struct {
	out		io.Writer
	dry		bool
	debug	bool
}

func newInitCmd(out io.Writer) *cobra.Command {
	ic := &initCmd{out: out}
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialises a template spec.yaml file for ocibuilder",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ic.run(args)
		},
	}
	f := cmd.Flags()
	f.BoolVar(&ic.dry, "dry", false, "Run a dry spec generation which is outputted to the terminal")
	f.BoolVarP(&ic.debug, "debug", "d", false, "Turn on debug logging")

	return cmd
}

func (i *initCmd) run(args []string) error {
	log := common.GetLogger(i.debug)
	box := packr.NewBox("../../templates/spec")

	template, err := box.Find("simple_spec_template.yaml")
	if err != nil {
		log.WithError(err).Errorln("error reading in template from docs")
		return err
	}

	if i.dry {
		if _, err := i.out.Write(template); err != nil {
			log.WithError(err).Errorln("error writing template to stdout")
			return err
		}
	}

	if err := ioutil.WriteFile("spec.yaml", template, 0777); err != nil {
		log.WithError(err).Errorln("error generating spec.yaml template file")
		return err
	}

	return nil
}
