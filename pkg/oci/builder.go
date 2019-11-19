package oci

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/sirupsen/logrus"
)

type Builder struct {
	Logger   *logrus.Logger
	Client   v1alpha1.BuilderClient
	Metadata []v1alpha1.ImageMetadata
}

func (b *Builder) Build(spec v1alpha1.OCIBuilderSpec, res chan<- v1alpha1.OCIBuildResponse, errChan chan<- error, finished chan bool) {
	log := b.Logger
	cli := b.Client

	reader := common.Reader{
		Logger: log,
	}

	buildOpts, err := common.ParseBuildSpec(spec.Build)
	if err != nil {
		log.WithError(err).Errorln("error in parsing build spec")
		errChan <- err
	}

	for idx, opt := range buildOpts {
		log.WithField("step: ", idx).Infoln("running build step")
		ctx, path, err := reader.ReadContext(opt.Context)
		if err != nil {
			log.WithError(err).Errorln("error reading image build context")
			errChan <- err
			return
		}

		imageName := fmt.Sprintf("%s:%s", opt.Name, opt.Tag)

		builderOptions := v1alpha1.OCIBuildOptions{
			Ctx:         context.Background(),
			ContextPath: path,
			Context:     ctx,
			ImageBuildOptions: types.ImageBuildOptions{
				Dockerfile: opt.Dockerfile,
				Tags:       []string{imageName},
				Context:    ctx,
			},
		}

		log.WithField("imageName", imageName).Debugln("building image with name")
		buildResponse, err := cli.ImageBuild(builderOptions)
		if err != nil {
			log.WithError(err).Errorln("error building image")
			errChan <- err
			return
		}
		log.Debugln("sending step build response")
		res <- buildResponse
		log.Debugln("finished sending step build response")
		if buildResponse.Exec != nil {
			log.Debugln("executing wait on build response")
			if err := buildResponse.Exec.Wait(); err != nil {
				errChan <- err
				return
			}
		}

		if opt.Purge {
			if err := b.Purge(imageName); err != nil {
				log.WithError(err).Errorln("unable to complete image purge")
				errChan <- err
				return
			}
		}
		log.WithField("response", idx).Debugln("response has finished executing")
	}
	close(res)
	close(errChan)
	log.Debugln("running build file cleanup")
	b.Clean()
	<-finished
}

func (b *Builder) Push(spec v1alpha1.OCIBuilderSpec, res chan<- v1alpha1.OCIPushResponse, errChan chan<- error) {
	log := b.Logger
	cli := b.Client

	for idx, pushSpec := range spec.Push {
		log.WithField("step: ", idx).Infoln("running push step")
		if err := common.ValidatePushSpec(pushSpec); err != nil {
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
			log.WithError(err).Errorln("failed to push image")
			errChan <- err
			return
		}
		res <- pushResponse

		if pushSpec.Purge {
			if err := b.Purge(pushFullImageName); err != nil {
				log.WithError(err).Errorln("unable to complete image purge")
				errChan <- err
				return
			}
		}
		log.WithField("response", idx).Debugln("response has finished executing")
	}
	close(res)
	close(errChan)
}

func (b *Builder) Pull(spec v1alpha1.OCIBuilderSpec, imageName string, res chan<- v1alpha1.OCIPullResponse, errChan chan<- error) {
	log := b.Logger
	cli := b.Client

	log.Infoln("attempting to pull from logged in registries")
	for idx, loginSpec := range spec.Login {
		registry := loginSpec.Registry
		if registry != "" {
			registry = registry + "/"
		}

		authString, err := b.generateAuthRegistryString(registry, spec)
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

		log.WithField("response", idx).Debugln("response has finished executing")
	}
	close(res)
	close(errChan)
}

func (b *Builder) Login(spec v1alpha1.OCIBuilderSpec, res chan<- v1alpha1.OCILoginResponse, errChan chan<- error) {
	log := b.Logger
	cli := b.Client

	if err := common.ValidateLogin(spec); err != nil {
		errChan <- err
		return
	}

	for _, loginSpec := range spec.Login {
		log.WithFields(logrus.Fields{"registry": loginSpec.Registry}).Infoln("attempting to login to registry")
		username, err := common.ValidateLoginUsername(loginSpec)
		if err != nil {
			errChan <- err
			return
		}

		password, err := common.ValidateLoginPassword(loginSpec)
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
	}
	close(res)
	close(errChan)
	log.Infoln("login complete")
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

func (b Builder) Clean() {
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

func (b Builder) generateAuthRegistryString(registry string, spec v1alpha1.OCIBuilderSpec) (string, error) {
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
			return b.Client.GenerateAuthRegistryString(types.AuthConfig{
				Username: user,
				Password: pass,
			}), nil
		}
	}
	return "", errors.New("no auth credentials matching registry: " + registry + " found")
}
