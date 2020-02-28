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

package oci

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/common"
	"github.com/ocibuilder/ocibuilder/pkg/parser"
	"github.com/ocibuilder/ocibuilder/pkg/validate"
	"github.com/sirupsen/logrus"
)

type Builder struct {
	Logger     *logrus.Logger
	Client     v1alpha1.BuilderClient
	Provenance []v1alpha1.BuildProvenance
}

func (b *Builder) Build(spec v1alpha1.OCIBuilderSpec, res chan<- v1alpha1.OCIBuildResponse, errChan chan<- error, finished chan<- bool) {
	log := b.Logger
	cli := b.Client

	defer func() {
		b.Clean()
		finished <- true
	}()

	buildOpts, err := parser.ParseBuildSpec(spec.Build)
	if err != nil {
		log.WithError(err).Errorln("error in parsing build spec")
		errChan <- err
	}

	for idx, opt := range buildOpts {
		buildProvenance := v1alpha1.BuildProvenance{
			BuildFile:        opt.Dockerfile,
			ContextDirectory: opt.BuildContextPath,
			Creator:          opt.Creator,
			Source:           opt.Source,
			Name:             opt.Name,
			Tag:              opt.Tag,
		}

		log.WithField("step: ", idx).Debugln("running build step")
		log.WithField("path", opt.BuildContextPath).Debugln("building with build context at path")
		buildContext, err := os.Open(opt.BuildContextPath + common.ContextDirectory + common.ContextFile)
		if err != nil {
			log.WithError(err).Errorln("error reading image build context")
			errChan <- err
			return
		}

		imageName := fmt.Sprintf("%s:%s", opt.Name, opt.Tag)

		builderOptions := v1alpha1.OCIBuildOptions{
			Ctx:         context.Background(),
			ContextPath: opt.BuildContextPath + common.ContextDirectory,
			Context:     buildContext,
			ImageBuildOptions: types.ImageBuildOptions{
				Dockerfile: opt.Dockerfile,
				Tags:       []string{imageName},
				Context:    buildContext,
				Labels:     opt.Labels,
				NoCache:    !opt.Cache,
			},
		}

		buildProvenance.StartTime = time.Now()
		log.WithField("imageName", imageName).Debugln("building image with name")
		buildResponse, err := cli.ImageBuild(builderOptions)
		if err != nil {
			log.WithError(err).Errorln("error building image")
			errChan <- err
			return
		}

		res <- buildResponse
		if buildResponse.Exec != nil {
			log.Debugln("executing wait on build response")
			if err := buildResponse.Exec.Wait(); err != nil {
				errChan <- err
				return
			}
		}

		buildProvenance.EndTime = time.Now()
		if spec.Metadata.StoreConfig != nil {
			log.Debugln("metadata specification present")
			mw := NewMetadataWriter(log, spec.Metadata)

			if err := mw.ParseMetadata(imageName, b.Client, buildProvenance); err != nil {
				log.Errorln("error parsing image metadata to push to store")
				errChan <- err
			}

			if err := mw.Write(); err != nil {
				errChan <- err
			}

		}

		if opt.Purge {
			if err := b.Purge(imageName); err != nil {
				log.WithError(err).Errorln("unable to complete image purge")
				errChan <- err
				return
			}
		}
		b.Provenance = append(b.Provenance, buildProvenance)
		log.WithField("step", idx).Debugln("build step has finished excuting")
	}
}

func (b *Builder) Push(spec v1alpha1.OCIBuilderSpec, res chan<- v1alpha1.OCIPushResponse, errChan chan<- error, finished chan<- bool) {
	log := b.Logger
	cli := b.Client

	for idx, pushSpec := range spec.Push {
		log.WithField("step: ", idx).Debugln("running push step")
		if err := validate.ValidatePushSpec(&pushSpec); err != nil {
			errChan <- err
		}

		pushFullImageName := fmt.Sprintf("%s/%s:%s", pushSpec.Registry, pushSpec.Image, pushSpec.Tag)
		log.WithField("name", pushFullImageName).Infoln("pushing image with name")

		authString, err := b.generateAuthRegistryString(pushSpec.Registry, spec)
		if err != nil {
			log.WithError(err).Errorln("unable to find login spec")
			errChan <- err
			return
		}

		pushOptions := v1alpha1.OCIPushOptions{
			Ctx: context.Background(),
			Ref: pushFullImageName,
			ImagePushOptions: types.ImagePushOptions{
				RegistryAuth: authString,
			},
		}

		pushResponse, err := cli.ImagePush(pushOptions)
		if err != nil {
			log.WithError(err).Debugln("failed to push image")
			errChan <- err
			return
		}

		res <- pushResponse
		if pushResponse.Exec != nil {
			log.Debugln("executing wait on push response")
			if err := pushResponse.Exec.Wait(); err != nil {
				errChan <- err
				return
			}
		}

		if pushSpec.Purge {
			if err := b.Purge(pushFullImageName); err != nil {
				log.WithError(err).Errorln("unable to complete image purge")
				errChan <- err
				return
			}
		}
		log.WithField("step", idx).Debugln("push step has finished executing")
	}
	finished <- true
}

func (b *Builder) Pull(spec v1alpha1.OCIBuilderSpec, imageName string, res chan<- v1alpha1.OCIPullResponse, errChan chan<- error, finished chan<- bool) {
	log := b.Logger
	cli := b.Client

	for idx, loginSpec := range spec.Login {
		registry := loginSpec.Registry
		log.WithField("registry", registry).Debugln("attempting to pull from logged in registry")
		authString, err := b.generateAuthRegistryString(registry, spec)

		if registry != "" {
			registry = registry + "/"
		}

		if err != nil {
			errChan <- err
			return
		}

		pullOptions := v1alpha1.OCIPullOptions{
			Ctx: context.Background(),
			Ref: registry + imageName,
			ImagePullOptions: types.ImagePullOptions{
				RegistryAuth: authString,
			},
		}

		pullResponse, err := cli.ImagePull(pullOptions)
		if err != nil {
			log.WithError(err).Errorln("failed to pull image")
			errChan <- err
			return
		}

		res <- pullResponse
		if pullResponse.Exec != nil {
			log.Debugln("executing wait on pull response")
			if err := pullResponse.Exec.Wait(); err != nil {
				errChan <- err
				return
			}
		}

		log.WithField("step", idx).Debugln("finished pull attempt from registry")
	}
	finished <- true
}

func (b *Builder) Login(spec v1alpha1.OCIBuilderSpec, res chan<- v1alpha1.OCILoginResponse, errChan chan<- error, finished chan<- bool) {
	log := b.Logger
	cli := b.Client

	if err := validate.ValidateLogin(spec); err != nil {
		errChan <- err
		return
	}

	for idx, loginSpec := range spec.Login {
		log.WithField("registry", loginSpec.Registry).Debugln("attempting to login to registry")
		username, err := validate.ValidateLoginUsername(loginSpec)
		if err != nil {
			errChan <- err
			return
		}

		password, err := validate.ValidateLoginPassword(loginSpec)
		if err != nil {
			errChan <- err
			return
		}
		loginOptions := v1alpha1.OCILoginOptions{
			Ctx: context.Background(),
			AuthConfig: types.AuthConfig{
				Username:      username,
				Password:      password,
				ServerAddress: loginSpec.Registry,
			},
		}

		loginResponse, err := cli.RegistryLogin(loginOptions)
		if err != nil {
			log.WithError(err).Errorln("failed to pull image")
			errChan <- err
			return
		}

		res <- loginResponse
		log.WithField("step", idx).Debugln("login step has finished executing")
	}
	finished <- true
}

func (b *Builder) Purge(imageName string) error {
	log := b.Logger
	cli := b.Client

	log.WithField("image", imageName).Debugln("attempting to purge image")

	removeOptions := v1alpha1.OCIRemoveOptions{
		Image:              imageName,
		Ctx:                context.Background(),
		ImageRemoveOptions: types.ImageRemoveOptions{},
	}

	res, err := cli.ImageRemove(removeOptions)
	if err != nil {
		log.WithError(err).Errorln("unable to complete image purge")
		return err
	}

	log.WithFields(logrus.Fields{"response": res}).Infoln("images purged")
	return nil
}

func (b *Builder) Clean() {
	log := b.Logger
	log.WithField("metadata", b.Provenance).Debugln("attempting to cleanup files listed in metadata")
	for _, m := range b.Provenance {
		if m.ContextDirectory != "" {
			log.WithField("filepath", m.ContextDirectory).Debugln("attempting to cleanup context")
			if err := os.RemoveAll(m.ContextDirectory + "/ocib"); err != nil {
				b.Logger.WithError(err).Errorln("error removing generated context")
				continue
			}
			if _, err := os.Stat(common.DockerStepPath); err == nil {
				if err := os.Remove(common.DockerStepPath); err != nil {
					b.Logger.WithError(err).Errorln("error removing downloaded step file")
					continue
				}
			}
		}
	}
}

func (b Builder) generateAuthRegistryString(registry string, spec v1alpha1.OCIBuilderSpec) (string, error) {
	if err := validate.ValidateLogin(spec); err != nil {
		return "", err
	}
	for _, spec := range spec.Login {
		if spec.Registry == registry {
			user, err := validate.ValidateLoginUsername(spec)
			if err != nil {
				return "", err
			}

			pass, err := validate.ValidateLoginPassword(spec)
			if err != nil {
				return "", err
			}
			return b.Client.GenerateAuthRegistryString(types.AuthConfig{
				Username: user,
				Password: pass,
			}), nil
		}
	}
	return "", errors.New("no auth credentials matching registry: " + registry + " found")
}
