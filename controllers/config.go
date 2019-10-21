package ocibuilder

import (
	"context"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/ocibuilder/ocibuilder/common"
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
					err := ctrl.updateConfig(cm)
					if err != nil {
						ctrl.logger.Error("update of config failed", "err", err)
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				if newCm, ok := new.(*corev1.ConfigMap); ok {
					ctrl.logger.Info("detected configmap update. updating the controller config.")
					err := ctrl.updateConfig(newCm)
					if err != nil {
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
	x := ctrl.kubeClient.CoreV1().RESTClient()
	resource := "configmaps"
	name := ctrl.Configmap
	fieldSelector := fields.ParseSelectorOrDie(fmt.Sprintf("metadata.name=%s", name))

	listFunc := func(options metav1.ListOptions) (runtime.Object, error) {
		options.FieldSelector = fieldSelector.String()
		req := x.Get().
			Namespace(ctrl.Namespace).
			Resource(resource).
			VersionedParams(&options, metav1.ParameterCodec)
		return req.Do().Get()
	}
	watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
		options.Watch = true
		options.FieldSelector = fieldSelector.String()
		req := x.Get().
			Namespace(ctrl.Namespace).
			Resource(resource).
			VersionedParams(&options, metav1.ParameterCodec)
		return req.Watch()
	}
	return &cache.ListWatch{ListFunc: listFunc, WatchFunc: watchFunc}
}

// ResyncConfig reloads the controller config from the configmap
func (ctrl *Controller) ResyncConfig(namespace string) error {
	cmClient := ctrl.kubeClient.CoreV1().ConfigMaps(namespace)
	cm, err := cmClient.Get(ctrl.Configmap, metav1.GetOptions{})
	if err != nil {
		return err
	}
	return ctrl.updateConfig(cm)
}

// updates the controller configuration
func (ctrl *Controller) updateConfig(cm *corev1.ConfigMap) error {
	configStr, ok := cm.Data[common.ControllerConfigMapKey]
	if !ok {
		return fmt.Errorf("configMap '%s' does not have key '%s'", ctrl.Configmap, common.ControllerConfigMapKey)
	}
	var config *ControllerConfig
	err := yaml.Unmarshal([]byte(configStr), &config)
	if err != nil {
		return err
	}
	ctrl.Config = config
	return nil
}
