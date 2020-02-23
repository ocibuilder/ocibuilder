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
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// watches configuration for controller
func (ctrl *Controller) watchControllerConfigMap(ctx context.Context) (cache.Controller, error) {
	ctrl.logger.Info("watching controller config map updates")
	source := ctrl.newControllerConfigMapWatch()
	_, controller := cache.NewInformer(
		source,
		&corev1.ConfigMap{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if cm, ok := obj.(*corev1.ConfigMap); ok {
					ctrl.logger.Info("detected configmap update. updating the controller config")
					if err := ctrl.updateConfig(cm); err != nil {
						ctrl.logger.WithError(err).Errorln("update of config failed")
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				if newCm, ok := new.(*corev1.ConfigMap); ok {
					ctrl.logger.Info("detected configmap update. updating the controller config.")
					if err := ctrl.updateConfig(newCm); err != nil {
						ctrl.logger.WithError(err).Error("update of config failed")
					}
				}
			},
		})

	go controller.Run(ctx.Done())
	return controller, nil
}

// creates a new config map watcher
func (ctrl *Controller) newControllerConfigMapWatch() *cache.ListWatch {
	client := ctrl.kubeClient.CoreV1().RESTClient()
	resource := "configmaps"
	name := ctrl.configmap
	fieldSelector := fields.ParseSelectorOrDie(fmt.Sprintf("metadata.name=%s", name))

	listFunc := func(options metav1.ListOptions) (runtime.Object, error) {
		options.FieldSelector = fieldSelector.String()
		req := client.Get().
			Namespace(ctrl.namespace).
			Resource(resource).
			VersionedParams(&options, metav1.ParameterCodec)
		return req.Do().Get()
	}
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		options.Watch = true
		options.FieldSelector = fieldSelector.String()
		req := client.Get().
			Namespace(ctrl.namespace).
			Resource(resource).
			VersionedParams(&options, metav1.ParameterCodec)
		return req.Watch()
	}
	return &cache.ListWatch{ListFunc: listFunc, WatchFunc: watchFunc}
}

// ResyncConfig reloads the controller config from the configmap
func (ctrl *Controller) ResyncConfig(namespace string) error {
	cmClient := ctrl.kubeClient.CoreV1().ConfigMaps(namespace)
	cm, err := cmClient.Get(ctrl.configmap, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return ctrl.updateConfig(cm)
}

// updates the controller configuration
func (ctrl *Controller) updateConfig(cm *corev1.ConfigMap) error {
	configStr, ok := cm.Data[common.ControllerConfigMapKey]
	if !ok {
		return errors.Errorf("configMap '%s' does not have key '%s'", ctrl.configmap, common.ControllerConfigMapKey)
	}
	var config *ControllerConfig
	err := yaml.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return err
	}
	ctrl.config = config
	return nil
}
