package ocibuilder

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/ocibuilder/ocibuilder/common"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strconv"
)

// parseJobConfiguration parses configuration required to manage K8s job lifecycle
func (opctx *operationContext) parseJobConfiguration() *jobConfiguration {
	cfg := &jobConfiguration{
		backoffLimit:            common.BackoffLimit,
		completions:             common.Completions,
		parallelisms:            common.Parallelism,
		ttlSecondsAfterFinished: int32(common.TTLSecondsAfterFinished),
		activeDeadlineSeconds:   int64(common.ActiveDeadlineSeconds),
	}

	backoffLimitEnvVar, ok := os.LookupEnv(common.EnvVarBackoffLimit)
	if ok {
		backoffLimit, err := strconv.ParseInt(backoffLimitEnvVar, 10, 32)
		if err != nil {
			opctx.logger.WithError(err).WithField("env-var", common.EnvVarBackoffLimit).WithField("default-value", common.BackoffLimit).Warnln("failed to parse environment variable. Using the default value")
		} else {
			cfg.backoffLimit = int32(backoffLimit)
		}
	}

	completionsEnvVar, ok := os.LookupEnv(common.EnvVarCompletions)
	if ok {
		completions, err := strconv.ParseInt(completionsEnvVar, 10, 32)
		if err != nil {
			opctx.logger.WithError(err).WithField("env-var", common.EnvVarCompletions).WithField("default-value", common.Completions).Warnln("failed to parse environment variable. Using the default value")
		} else {
			cfg.completions = int32(completions)
		}
	}

	ttlSecondsAfterFinishedEnvVar, ok := os.LookupEnv(common.EnvVarTTLSecondsAfterFinished)
	if ok {
		ttlSecondsAfterFinished, err := strconv.ParseInt(ttlSecondsAfterFinishedEnvVar, 10, 32)
		if err != nil {
			opctx.logger.WithError(err).WithField("env-var", common.EnvVarTTLSecondsAfterFinished).WithField("default-value", common.TTLSecondsAfterFinished).Warnln("failed to parse environment variable. Using the default value")
		} else {
			cfg.ttlSecondsAfterFinished = int32(ttlSecondsAfterFinished)
		}
	}

	activeDeadlineSecondsEnvVar, ok := os.LookupEnv(common.EnvVarActiveDeadlineSeconds)
	if ok {
		activeDeadlineSeconds, err := strconv.ParseInt(activeDeadlineSecondsEnvVar, 10, 32)
		if err != nil {
			opctx.logger.WithError(err).WithField("env-var", common.EnvVarActiveDeadlineSeconds).WithField("default-value", common.ActiveDeadlineSeconds).Warnln("failed to parse environment variable. Using the default value")
		} else {
			cfg.activeDeadlineSeconds = activeDeadlineSeconds
		}
	}

	return cfg
}

// storeBuilderSpecification stores the builder specification in a file
func (opctx *operationContext) storeBuilderSpecification() error {
	file, err := os.Create(fmt.Sprintf("%s/%s", common.ContextDirectory, common.SpecFilePath))
	if err != nil {
		return err
	}
	body, err := yaml.Marshal(opctx.builder)
	if err != nil {
		return err
	}
	if _, err := file.Write(body); err != nil {
		return err
	}
	return nil
}

// constructBuilderJob constructs a K8s job for ocibuilder build step.
func (opctx *operationContext) constructBuilderJob() *batchv1.Job {
	jobCfg := opctx.parseJobConfiguration()

	labels := map[string]string{
		common.LabelOwner:   opctx.builder.Name,
		common.LabelJobType: "ocibuilder",
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
				Spec: corev1.PodSpec{},
			},
		},
	}
}
