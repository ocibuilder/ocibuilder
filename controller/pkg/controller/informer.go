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
	informers "github.com/ocibuilder/ocibuilder/controller/pkg/client/ocibuilder/informers/externalversions"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/tools/cache"
)

func (ctrl *Controller) instanceIDReq() (*labels.Requirement, error) {
	if ctrl.config.InstanceID == "" {
		return nil, errors.New("instance id is required")
	}
	instanceIDReq, err := labels.NewRequirement(common.LabelKeyControllerInstanceID, selection.Equals, []string{ctrl.config.InstanceID})
	if err != nil {
		return nil, err
	}
	return instanceIDReq, nil
}

// newControllerInformer adds new ocibuilders to the controller's queue based on Add, Update, and Delete Event Handlers for the ocibuilder resources
func (ctrl *Controller) newControllerInformer(labelFilterRequirements *labels.Requirement) cache.SharedIndexInformer {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		ctrl.ociClient,
		resyncPeriod,
		informers.WithNamespace(ctrl.config.Namespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.Everything().String()
			labelSelector := labels.NewSelector().Add(*labelFilterRequirements)
			options.LabelSelector = labelSelector.String()
		}),
	)
	informer := informerFactory.Ocibuilder().V1alpha1().OCIBuilders().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(obj)
				if err == nil {
					ctrl.queue.Add(key)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err == nil {
					ctrl.queue.Add(key)
				}
			},
			DeleteFunc: func(obj interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err == nil {
					ctrl.queue.Add(key)
				}
			},
		},
	)
	return informer
}
