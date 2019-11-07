package oci

import (
	"context"
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

func (b *Builder) Build(spec v1alpha1.OCIBuilderSpec, res chan<- v1alpha1.OCIBuildResponse, errChan chan<- error) {
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
		res <- v1alpha1.OCIBuildResponse{
			Body: buildResponse.Body,
			Metadata: v1alpha1.ImageMetadata{
				BuildFile: fmt.Sprintf("%s/%s", path, opt.Dockerfile),
				Daemon:    spec.Daemon,
			},
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
	log.Infoln("build complete")
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
