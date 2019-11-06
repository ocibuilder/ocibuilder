/*
Copyright © 2019 BlackRock Inc.

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

package buildah

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Buildah is  the struct which consists of a logger and context path
type Buildah struct {
	Logger        *logrus.Logger
	StorageDriver string
	Metadata      []v1alpha1.ImageMetadata
	execCmds      []*exec.Cmd
}

var executor = exec.Command

// Build performs a buildah build and returns an array of readclosers
func (b *Buildah) Build(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	log := b.Logger
	buildOpts, err := common.ParseBuildSpec(spec.Build)

	if err != nil {
		log.WithError(err).Errorln("failed to parse build spec...")
		return nil, err
	}

	var buildResponses []io.ReadCloser
	for _, opt := range buildOpts {
		imageName := fmt.Sprintf("%s:%s", opt.Name, opt.Tag)
		opt.Context = common.ValidateContext(opt.Context)

		fullPath := opt.Context.LocalContext.ContextPath + "/" + opt.Dockerfile
		if opt.Context.LocalContext.ContextPath == "" {
			fullPath = "." + fullPath
		}

		buildCommand := createBuildCommand(opt, b.StorageDriver)
		log.WithField("command", buildCommand).Debug("build command to be executed")

		cmd := executor("buildah", buildCommand...)
		out, err := pty.Start(cmd)
		b.execCmds = append(b.execCmds, cmd)

		if err != nil {
			log.WithError(err).Errorln("failed to execute buildah bud...")
			return nil, err
		}
		buildResponses = append(buildResponses, out)

		b.Metadata = append(b.Metadata, v1alpha1.ImageMetadata{
			BuildFile: fullPath,
		})

		if opt.Purge {
			purgeCommand := createPurgeCommand(imageName)
			log.WithFields(logrus.Fields{"command": purgeCommand, "image": imageName}).Debug("purge command to be executed")

			cmd = executor("buildah", purgeCommand...)
			out, err = pty.Start(cmd)
			if err != nil {
				log.WithError(err).Errorln("failed to execute purge")
				return nil, err
			}
			log.WithField("response", out).Infoln("images purged")
		}
	}
	return buildResponses, nil
}

// createBuildCommand is used to generate build command and build args
func createBuildCommand(args v1alpha1.ImageBuildArgs, storageDriver string) []string {
	buildArgs := append([]string{"bud"}, "-f", args.Dockerfile)

	if storageDriver != "" {
		buildArgs = append(buildArgs, "--storage-driver", storageDriver)
	}

	image := ""
	if args.Name != "" {
		image += args.Name
	}
	if args.Tag != "" {
		image += ":" + args.Tag
	}

	if image != "" {
		return append(buildArgs, "-t", image, args.Context.LocalContext.ContextPath)
	}
	return append(buildArgs, args.Context.LocalContext.ContextPath)
}

// Login performs a buildah login on all registries defined in ocibuilder.yaml or login.yaml
func (b *Buildah) Login(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	log := b.Logger

	if err := common.ValidateLogin(spec); err != nil {
		return nil, err
	}

	var loginResponses []io.ReadCloser
	for _, loginSpec := range spec.Login {
		loginCommand, err := createLoginCommand(loginSpec)

		if err != nil {
			log.WithError(err).Errorln("error creating login command")
			return nil, err
		}

		cmd := executor("buildah", loginCommand...)
		b.execCmds = append(b.execCmds, cmd)

		out, err := pty.Start(cmd)
		if err != nil {
			log.WithError(err).Errorln("failed to execute buildah login...")
			return nil, err
		}
		loginResponses = append(loginResponses, out)

		log.Infoln("buildah login has been executed")
	}

	return loginResponses, nil
}

func createLoginCommand(args v1alpha1.LoginSpec) ([]string, error) {
	loginArgs := []string{"login"}

	registry := args.Registry
	if registry == "" {
		return nil, errors.New("no registry has been specified for login")
	}

	username, err := common.ValidateLoginUsername(args)
	if err != nil {
		return nil, err
	}

	password, err := common.ValidateLoginPassword(args)
	if err != nil {
		return nil, err
	}

	return append(loginArgs, "-u", username, "-p", password, registry), nil
}

// Pull performs a buildah pull of a passed in image name. Pull will login to all
// registries specified in the 'login' spec and attempt to pull the image
// uses buildah login to login to directories specified
func (b *Buildah) Pull(spec v1alpha1.OCIBuilderSpec, imageName string) ([]io.ReadCloser, error) {
	log := b.Logger

	var pullResponses []io.ReadCloser
	for _, loginSpec := range spec.Login {
		pullCommand, err := b.createPullCommand(loginSpec.Registry, imageName, spec)

		if err != nil {
			log.WithError(err).Errorln("error attempting to create pull command")
			return nil, err
		}

		cmd := executor("buildah", pullCommand...)
		out, err := pty.Start(cmd)
		b.execCmds = append(b.execCmds, cmd)
		if err != nil {
			log.WithError(err).Errorln("failed to execute buildah pull...")
			return nil, err
		}
		pullResponses = append(pullResponses, out)

		log.Infoln("buildah pull has been executed")
	}

	return pullResponses, nil
}

func (b Buildah) createPullCommand(registry string, imageName string, spec v1alpha1.OCIBuilderSpec) ([]string, error) {
	args := []string{"pull", "--creds"}

	fullImageName := fmt.Sprintf("%s/%s", registry, imageName)
	b.Logger.WithField("command", append(args, fullImageName)).Debugln("push command with AUTH REVOKED")

	authString, err := getPushAuthRegistryString(registry, spec)
	if err != nil {
		return nil, err
	}
	args = append(args, authString)

	return append(args, fullImageName), nil
}

// Push performs a buildah push of a spec image to a chosen registry
// uses buildah login to login to directories specified
func (b *Buildah) Push(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	log := b.Logger

	if err := common.ValidatePush(spec); err != nil {
		return nil, err
	}

	var pushResponses []io.ReadCloser
	for _, pushSpec := range spec.Push {

		if err := common.ValidatePushSpec(pushSpec); err != nil {
			return nil, err
		}

		imageName := fmt.Sprintf("%s:%s", pushSpec.Image, pushSpec.Tag)
		log.WithFields(logrus.Fields{"name": imageName, "registry": pushSpec.Registry}).Infoln("pushing image")

		pushCommand, err := b.createPushCommand(pushSpec.Registry, imageName, spec)

		if err != nil {
			log.WithError(err).Errorln("error attempting to create push command")
			return nil, err
		}

		cmd := executor("buildah", pushCommand...)
		out, err := pty.Start(cmd)
		b.execCmds = append(b.execCmds, cmd)
		if err != nil {
			log.WithError(err).Errorln("failed to execute buildah push...")
			return nil, err
		}
		pushResponses = append(pushResponses, out)

		log.Infoln("buildah push has been executed")

		if pushSpec.Purge {
			fullImageName := fmt.Sprintf("%s:%s", pushSpec.Registry, imageName)
			purgeCommand := createPurgeCommand(fullImageName)
			log.WithFields(logrus.Fields{"command": purgeCommand, "image": fullImageName}).Debug("purge command to be executed")

			cmd = executor("buildah", purgeCommand...)
			out, err = pty.Start(cmd)
			if err != nil {
				log.WithError(err).Errorln("failed to execute purge")
				return nil, err
			}
			log.WithField("response", out).Infoln("images purged")
		}
	}
	return pushResponses, nil
}

func (b Buildah) createPushCommand(registry string, imageName string, spec v1alpha1.OCIBuilderSpec) ([]string, error) {
	args := []string{"push", "--creds"}
	fullImageName := fmt.Sprintf("%s/%s", registry, imageName)
	b.Logger.WithField("command", append(args, fullImageName)).Debugln("push command with AUTH REVOKED")

	authString, err := getPushAuthRegistryString(registry, spec)
	if err != nil {
		return nil, err
	}
	args = append(args, authString)

	return append(args, fullImageName), nil
}

func getPushAuthRegistryString(registry string, spec v1alpha1.OCIBuilderSpec) (string, error) {
	if err := common.ValidateLogin(spec); err != nil {
		return "", err
	}
	for _, spec := range spec.Login {
		if spec.Registry == registry {
			user, err := common.ValidateLoginUsername(spec)
			if err != nil {
				return "", err
			}

			pass, err := common.ValidateLoginPassword(spec)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s:%s", user, pass), nil
		}
	}
	return "", errors.New("no auth credentials matching registry: " + registry + " found")
}

func createPurgeCommand(imageName string) []string {
	return append([]string{"rmi"}, imageName)
}

func (b Buildah) Clean() {
	log := b.Logger
	for _, m := range b.Metadata {
		if m.BuildFile != "" {
			log.WithField("filepath", m.BuildFile).Debugln("attempting to cleanup dockerfile")
			if err := os.Remove(m.BuildFile); err != nil {
				b.Logger.WithError(err).Errorln("error removing generated Dockerfile")
				continue
			}
		}
	}
}

// Wait calls wait for each exec command, handling any output to stderr and exiting the process
func (b Buildah) Wait(idx int) error {
	b.Logger.WithField("execCmds", b.execCmds).Debugln("exec wait called")
	if len(b.execCmds) == 0 {
		return errors.New("error waiting for command to finish executing")
	}
	if err := b.execCmds[idx].Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code is %d\n", exitError.ExitCode())
			error := fmt.Sprintf("error in executing cmd, exited with code %d", exitError.ExitCode())
			return errors.New(error)
		}
	}
	return nil
}
