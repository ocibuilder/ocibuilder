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
	"fmt"
	"log"
	"os"

	"github.com/ghodss/yaml"
	"github.com/ocibuilder/ocibuilder/common"
	ocibv1alpha1 "github.com/ocibuilder/ocibuilder/controller/pkg/client/ocibuilder/clientset/versioned"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	namespace, ok := os.LookupEnv(common.Namespace)
	if !ok {
		namespace = "default"
	}

	resourceName, ok := os.LookupEnv(common.ResourceName)
	if !ok {
		log.Fatalln("ocibuilder resource name not found")
		return
	}

	kubeConfig, ok := os.LookupEnv(common.EnvVarKubeConfig)
	if !ok {
		log.Println("kubeconfig not found")
	}

	restConfig, err := common.GetClientConfig(kubeConfig)
	if err != nil {
		log.Fatalln(err)
		return
	}

	client := ocibv1alpha1.NewForConfigOrDie(restConfig)
	ociObject, err := client.OcibuilderV1alpha1().OCIBuilders(namespace).Get(resourceName, metav1.GetOptions{
		TypeMeta:        metav1.TypeMeta{},
		ResourceVersion: "",
	})

	if err != nil {
		log.Fatalln(err)
		return
	}

	if err := storeBuilderSpecification(ociObject.Spec); err != nil {
		log.Fatalln(err)
		return
	}

}

// storeBuilderSpecification stores the builder specification in a file
func storeBuilderSpecification(spec v1alpha1.OCIBuilderSpec) error {
	file, err := os.Create(fmt.Sprintf("%s/%s", common.VolumeMountPath, common.SpecFilePath))
	if err != nil {
		return err
	}
	body, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}
	if _, err := file.Write(body); err != nil {
		return err
	}
	return nil
}
