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
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/validate"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// the context of an operation on a ocibuilder object.
// the controller creates this context each time it picks a ocibuilder object off its queue.
type operationContext struct {
	// builder is the ocibuilder object
	builder *v1alpha1.OCIBuilder
	// updated indicates whether the ocibuilder object was updated and needs to be persisted back to k8
	updated bool
	// logger is the logging context to correlate logs with a ocibuilder
	logger *logrus.Logger
	// reference to the controller
	controller *Controller
}

// jobConfiguration contains the configuration related to a K8s job lifecycle
type jobConfiguration struct {
	backoffLimit            int32
	completions             int32
	parallelisms            int32
	ttlSecondsAfterFinished int32
	activeDeadlineSeconds   int64
}

// newOperationContext returns a new context of controller operation
func newOperationContext(builder *v1alpha1.OCIBuilder, controller *Controller) *operationContext {
	return &operationContext{
		builder:    builder,
		controller: controller,
		logger: controller.logger.WithFields(map[string]interface{}{
			common.LabelOCIBuilderName: builder.Name,
			common.LabelNamespace:      builder.Namespace,
		}).Logger,
		updated: false,
	}
}

// operate operate on an ocibuilder object and manages its lifecycle
func (opctx *operationContext) Operate() error {
	log := opctx.logger.WithFields(map[string]interface{}{
		common.LabelName:      opctx.builder.Name,
		common.LabelNamespace: opctx.builder.Namespace,
	})

	log.Infoln("operating on the resource...")

	if err := validate.Validate(&opctx.builder.Spec); err != nil {
		return errors.Wrap(err, "failed to validate the resource spec")
	}

	switch opctx.builder.Status.Phase {
	case v1alpha1.NodePhaseNew:
		if _, err := opctx.constructBuilderJob(); err != nil {
			return err
		}
	case v1alpha1.NodePhaseRunning:
	case v1alpha1.NodePhaseCompleted:
	case v1alpha1.NodePhaseError:
	default:
		opctx.logger.WithField(common.LabelPhase, opctx.builder.Status.Phase).Warnln("unknown phase of the resource")
	}

	return nil
}
