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

package controller

import (
	"fmt"
	"path"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/command"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// parseJobConfiguration parses configuration required to manage K8s job lifecycle
func (opctx *operationContext) readJobConfiguration() *jobConfiguration {
	cfg := &jobConfiguration{
		backoffLimit:            common.BackoffLimit,
		completions:             common.Completions,
		parallelisms:            common.Parallelism,
		ttlSecondsAfterFinished: int32(common.TTLSecondsAfterFinished),
		activeDeadlineSeconds:   int64(common.ActiveDeadlineSeconds),
	}
	if opctx.builder.Spec.JobTemplate != nil {
		if opctx.builder.Spec.JobTemplate.BackoffLimit != nil {
			cfg.backoffLimit = *opctx.builder.Spec.JobTemplate.BackoffLimit
		}
		if opctx.builder.Spec.JobTemplate.Completions != nil {
			cfg.completions = *opctx.builder.Spec.JobTemplate.Completions
		}
		if opctx.builder.Spec.JobTemplate.TTLSecondsAfterFinished != nil {
			cfg.ttlSecondsAfterFinished = *opctx.builder.Spec.JobTemplate.TTLSecondsAfterFinished
		}
		if opctx.builder.Spec.JobTemplate.ActiveDeadlineSeconds != nil {
			cfg.activeDeadlineSeconds = *opctx.builder.Spec.JobTemplate.ActiveDeadlineSeconds
		}
	}
	return cfg
}

// generateCommands generates commands to be executed for the job
func (opctx *operationContext) generateCommands() []string {
	specificationFilePath := path.Clean(fmt.Sprintf("%s/%s", common.ContextDirectory, common.SpecFilePath))

	flags := []command.Flag{
		{
			Name:  "--path",
			Value: specificationFilePath,
		},
	}
	if opctx.builder.Spec.OverlayPath == "" {
		flags = append(flags, command.Flag{
			Name:  "--overlay",
			Value: path.Clean(fmt.Sprintf("%s/%s", common.ContextDirectoryUncompressed, opctx.builder.Spec.OverlayPath)),
		})
	}

	buildCmd := command.Builder(common.CmdOcictl).Command(common.CmdBuild).Flags(flags...).Build()
	pushCmd := command.Builder(common.CmdOcictl).Command(common.CmdPush).Flags(flags...).Build()

	cmd := append(buildCmd.String(), "&&")
	cmd = append(cmd, pushCmd.String()...)
	return cmd
}

// constructBuilderJob constructs a K8s job for ocibuilder build step.
func (opctx *operationContext) constructBuilderJob() (*batchv1.Job, error) {
	jobCfg := opctx.readJobConfiguration()

	labels := map[string]string{
		common.LabelOwner:   opctx.builder.Name,
		common.LabelJobType: "ocibuilder",
	}

	container := corev1.Container{
		Name:      common.Name,
		Image:     fmt.Sprintf("%s:%s", common.Image, common.Tag),
		Command:   opctx.generateCommands(),
		Resources: corev1.ResourceRequirements{},
	}

	initContainer := corev1.Container{
		Name:  common.InitName,
		Image: fmt.Sprintf("%s:%s", common.InitImage, common.InitTag),
		Env: []corev1.EnvVar{
			{
				Name:  common.Namespace,
				Value: opctx.controller.namespace,
			},
			{
				Name:  common.ResourceName,
				Value: opctx.builder.Name,
			},
		},
		Resources: corev1.ResourceRequirements{},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      common.VolumeName,
				MountPath: common.VolumeMountPath,
			},
		},
	}

	volume := corev1.Volume{
		Name: common.VolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium: corev1.StorageMediumDefault,
			},
		},
	}

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: opctx.builder.Name,
			Namespace:    opctx.builder.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(opctx.builder, opctx.builder.GroupVersionKind()),
			},
		},
		Spec: batchv1.JobSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			BackoffLimit:            &jobCfg.backoffLimit,
			Completions:             &jobCfg.completions,
			Parallelism:             &jobCfg.parallelisms,
			TTLSecondsAfterFinished: &jobCfg.ttlSecondsAfterFinished,
			ActiveDeadlineSeconds:   &jobCfg.activeDeadlineSeconds,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: opctx.builder.Name,
					Namespace:    opctx.builder.Namespace,
					Labels:       labels,
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						initContainer,
					},
					Containers: []corev1.Container{
						container,
					},
					Volumes: []corev1.Volume{
						volume,
					},
				},
			},
		},
	}, nil
}
