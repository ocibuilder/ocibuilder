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

package main

import (
	"context"
	"os"

	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/controller/pkg/controller"
)

func main() {

	logger := common.Logger
	kubeConfig, ok := os.LookupEnv(common.EnvVarKubeConfig)
	if !ok {
		logger.Infoln("kubeconfig not found")
		return
	}

	restConfig, err := common.GetClientConfig(kubeConfig)
	if err != nil {
		return
	}

	namespace, ok := os.LookupEnv(common.EnvVarNamespace)
	if !ok {
		logger.Infoln("namespace not found")
		return
	}

	instanceId, ok := os.LookupEnv(common.EnvVarControllerInstanceID)
	if !ok {
		logger.Infoln("instanceID not found")
		return
	}

	ctrlConfig := &controller.ControllerConfig{
		InstanceID: instanceId,
		Namespace:  namespace,
	}

	configMap, ok := os.LookupEnv(common.EnvVarControllerConfigMap)
	if !ok {
		panic("controller configmap is not provided")
	}

	ctrl := controller.NewController(restConfig, logger, configMap, namespace)

	ctrl.Run(context.Background(), 1, 1)
}
