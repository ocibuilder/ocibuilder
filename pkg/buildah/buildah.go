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
	"os/exec"

	"github.com/creack/pty"
	"github.com/ocibuilder/ocibuilder/pkg"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Buildah contains configuration required to perform builah related operation
type Buildah struct {
	// Logger to log stuff
	Logger *logrus.Logger
	// StorageDrive to use when building image
	StorageDriver string
}

var executor = exec.Command

// Build performs a buildah build and returns an array of readclosers
func (b *Buildah) Build(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	var buildResponses []io.ReadCloser
	buildOpts, err := pkg.ParseBuildSpec(spec.Build)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse the build specification")
	}
	for _, opt := range buildOpts {
		buildCommand := createBuildCommand(opt, b.StorageDriver)
		b.Logger.WithField("command", buildCommand).Debug("build command to be executed")
		cmd := executor("b", buildCommand...)
		out, err := pty.Start(cmd)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to execute b build")
		}
		if err := cmd.Wait(); err != nil {
			return nil, errors.Wrapf(err, "error waiting for cmd execution")
		}
		if opt.Purge {
			imageName := fmt.Sprintf("%s:%s", opt.Name, opt.Tag)
			purgeCommand := createPurgeCommand(imageName)
			b.Logger.WithFields(logrus.Fields{
				"command": purgeCommand,
				"image":   imageName,
			}).Debug("purge command to be executed")
			cmd = executor("b", purgeCommand...)
			out, err = pty.Start(cmd)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to execute purge")
			}
			b.Logger.WithField("image", imageName).Infoln("image purged")
		}
		buildResponses = append(buildResponses, out)
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
		return append(buildArgs, "-t", image, args.BuildContextPath)
	}
	return append(buildArgs, args.BuildContextPath)
}

// Login performs a buildah login on all registries defined in spec.yaml or login.yaml
func (b *Buildah) Login(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	if err := pkg.ValidateLogin(spec); err != nil {
		return nil, err
	}
	var loginResponses []io.ReadCloser
	for _, loginSpec := range spec.Login {
		loginCommand, err := createLoginCommand(loginSpec)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create login command for %s", loginSpec.Registry)
		}
		b.Logger.WithField("command", loginCommand).Debug("login command to be executed")
		cmd := executor("b", loginCommand...)
		out, err := pty.Start(cmd)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to execute login command for %s", loginSpec.Registry)
		}
		if err := cmd.Wait(); err != nil {
			return nil, errors.Wrapf(err, "error waiting for login command execution")
		}
		loginResponses = append(loginResponses, out)
	}
	return loginResponses, nil
}

func createLoginCommand(args v1alpha1.LoginSpec) ([]string, error) {
	loginArgs := []string{"login"}
	registry := args.Registry
	if registry == "" {
		return nil, errors.New("no registry has been specified for login")
	}
	username, err := pkg.ValidateLoginUsername(args)
	if err != nil {
		return nil, err
	}
	password, err := pkg.ValidateLoginPassword(args)
	if err != nil {
		return nil, err
	}
	return append(loginArgs, "-u", username, "-p", password, registry), nil
}

// Pull performs a buildah pull of a passed in image name. Pull will login to all
// registries specified in the 'login' spec and attempt to pull the image
// uses buildah login to login to directories specified
func (b *Buildah) Pull(spec v1alpha1.OCIBuilderSpec, imageName string) ([]io.ReadCloser, error) {
	var pullResponses []io.ReadCloser
	if _, err := b.Login(spec); err != nil {
		return nil, err
	}
	for _, loginSpec := range spec.Login {
		pullCommand, err := createPullCommand(imageName, loginSpec.Registry)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create pull command for %s", imageName)
		}
		b.Logger.WithField("command", pullCommand).Debug("pull command to be executed")
		cmd := executor("b", pullCommand...)
		out, err := pty.Start(cmd)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to execute pull command for %s", imageName)
		}
		if err := cmd.Wait(); err != nil {
			return nil, errors.Wrapf(err, "error waiting for pull command execution for %s", imageName)
		}
		pullResponses = append(pullResponses, out)
	}
	return pullResponses, nil
}

func createPullCommand(imageName string, registry string) ([]string, error) {
	pullArgs := []string{"pull"}
	if imageName == "" {
		return nil, errors.New("no image name specified to pull")
	}
	if registry != "" {
		registry = registry + "/"
	}
	return append(pullArgs, registry+imageName), nil
}

// Push performs a buildah push of a spec image to a chosen registry
// uses buildah login to login to directories specified
func (b *Buildah) Push(spec v1alpha1.OCIBuilderSpec) ([]io.ReadCloser, error) {
	if _, err := b.Login(spec); err != nil {
		return nil, err
	}
	if err := pkg.ValidatePush(spec); err != nil {
		return nil, err
	}
	var pushResponses []io.ReadCloser
	for _, pushSpec := range spec.Push {
		imageName := fmt.Sprintf("%s/%s:%s", pushSpec.Registry, pushSpec.Image, pushSpec.Tag)
		b.Logger.WithField("name", imageName).Infoln("pushing image...")
		pushCommand, err := createPushCommand(pushSpec, imageName)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create push command for %s", imageName)
		}
		cmd := executor("b", pushCommand...)
		out, err := pty.Start(cmd)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to push image %s", imageName)
		}
		if err := cmd.Wait(); err != nil {
			return nil, errors.Wrapf(err, "failed to wait for push command to execute for %s", imageName)
		}
		if pushSpec.Purge {
			purgeCommand := createPurgeCommand(imageName)
			b.Logger.WithFields(logrus.Fields{"command": purgeCommand, "image": imageName}).Debug("purge command to be executed")
			cmd = executor("b", purgeCommand...)
			out, err = pty.Start(cmd)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to purge the image %s", imageName)
			}
		}
		pushResponses = append(pushResponses, out)
	}
	return pushResponses, nil
}

func createPushCommand(spec v1alpha1.PushSpec, imageName string) ([]string, error) {
	pushArgs := []string{"push"}
	if err := pkg.ValidatePushSpec(spec); err != nil {
		return nil, err
	}
	return append(pushArgs, imageName), nil
}

func createPurgeCommand(imageName string) []string {
	return append([]string{"rmi"}, imageName)
}
