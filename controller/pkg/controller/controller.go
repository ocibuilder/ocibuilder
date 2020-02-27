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
	"context"
	"log"
	"time"

	"github.com/ocibuilder/ocibuilder/common"
	ociv1alpha1 "github.com/ocibuilder/ocibuilder/controller/pkg/client/ocibuilder/clientset/versioned"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/provenance"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	resyncPeriod         = 30 * time.Minute
	rateLimiterBaseDelay = 5 * time.Second
	rateLimiterMaxDelay  = 1000 * time.Second
)

// ControllerConfig contain the configuration settings for the controller
type ControllerConfig struct {
	// InstanceID is a label selector to limit the controller's watch of ocibuilders to a specific instance.
	InstanceID string
	// namespace is a label selector filter to limit controller's watch to specific namespace
	Namespace string
}

// Controller listens for new ocibuilder resources and hands off handling of each resource on the queue to the operator
type Controller struct {
	// configmap is the name of the K8s configmap which contains controller configuration
	configmap string
	// namespace for controller
	namespace string
	// config is the controller's configuration
	config *ControllerConfig
	// kubernetes config and apis
	kubeConfig *rest.Config
	// kubeClient communicates with Kubernetes API server
	kubeClient kubernetes.Interface
	// ociClient is the client to operates on ocibuilder resource
	ociClient ociv1alpha1.Interface
	// informer provides eventually consistent linkage of its clients to the authoritative state of a given collection of objects.
	informer cache.SharedIndexInformer
	// queue is an interface that rate limits items being added to the queue.
	queue workqueue.RateLimitingInterface
	// logger is the logger for a controller
	logger *logrus.Logger
}

// NewController creates a new controller
func NewController(rest *rest.Config, logger *logrus.Logger, configmap, namespace string) *Controller {
	rateLimiter := workqueue.NewItemExponentialFailureRateLimiter(rateLimiterBaseDelay, rateLimiterMaxDelay)
	return &Controller{
		namespace:  namespace,
		configmap:  configmap,
		kubeConfig: rest,
		kubeClient: kubernetes.NewForConfigOrDie(rest),
		ociClient:  ociv1alpha1.NewForConfigOrDie(rest),
		queue:      workqueue.NewRateLimitingQueue(rateLimiter),
		logger:     logger,
	}
}

func (ctrl *Controller) processNextItem() bool {
	// Wait until there is a new item in the queue
	key, quit := ctrl.queue.Get()
	if quit {
		return false
	}
	defer ctrl.queue.Done(key)

	obj, exists, err := ctrl.informer.GetIndexer().GetByKey(key.(string))
	if err != nil {
		ctrl.logger.WithError(err).WithField("key", key).Errorln("failed to get ocibuilder from informer index")
		return true
	}

	if !exists {
		// this happens after ocibuilder was deleted, but work queue still had entry in it
		return true
	}

	builder, ok := obj.(*v1alpha1.OCIBuilder)
	if !ok {
		ctrl.logger.WithError(err).WithField("key", key).Errorln("key in index is not a builder")
		return true
	}

	ctx := newOperationContext(builder, ctrl)

	err = ctx.Operate()
	if err != nil {
		ctrl.logger.WithError(err).WithField(common.LabelOCIBuilderName, builder.Name).Errorln("failed to operate on the ocibuilder obejct")
	}

	err = ctrl.handleErr(err, key)
	if err != nil {
		ctrl.logger.WithError(err).Errorln("controller is unable to handle the error")
	}
	return true
}

// handleErr checks if an error happened and make sure we will retry later
// returns an error if unable to handle the error
func (ctrl *Controller) handleErr(err error, key interface{}) error {
	if err == nil {
		// Forget about the #AddRateLimited history of key on every successful sync
		// Ensure future updates for this key are not delayed because of outdated error history
		ctrl.queue.Forget(key)
		return nil
	}

	// due to the base delay of 5ms of the DefaultControllerRateLimiter
	// requeues will happen very quickly even after a ocibuilder pod goes down
	// we want to give the ocibuilder pod a chance to come back up so we give a generous number of retries
	if ctrl.queue.NumRequeues(key) < 20 {
		// Re-enqueue the key rate limited. This key will be processed later again.
		ctrl.queue.AddRateLimited(key)
		return nil
	}
	return errors.New("exceeded max requeues")
}

// Run executes the controller
func (ctrl *Controller) Run(ctx context.Context, gwThreads, eventThreads int) {
	defer ctrl.queue.ShutDown()
	ctrl.logger.WithFields(
		map[string]interface{}{
			common.LabelKeyControllerInstanceID: ctrl.config.InstanceID,
			common.LabelVersion:                 provenance.GetProvenance().Version,
		}).Info("starting controller")
	if _, err := ctrl.watchControllerConfigMap(ctx); err != nil {
		ctrl.logger.WithError(err).Error("failed to register watch for controller config map")
		return
	}

	labelFilters, err := ctrl.instanceIDReq()
	if err != nil {
		ctrl.logger.WithError(err).Error("failed to get instance id filter")
		return
	}

	ctrl.informer = ctrl.newControllerInformer(labelFilters)
	go ctrl.informer.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), ctrl.informer.HasSynced) {
		log.Panicf("timed out waiting for the caches to sync")
		return
	}

	for i := 0; i < gwThreads; i++ {
		go wait.Until(ctrl.runWorker, time.Second, ctx.Done())
	}

	<-ctx.Done()
}

func (ctrl *Controller) runWorker() {
	for ctrl.processNextItem() {
	}
}
