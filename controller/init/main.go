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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"os"

	ocbv1alpha1 "github.com/blackrock/ocibuilder/pkg/client/ocibuilder/clientset/versioned"
	"github.com/ghodss/yaml"
	"github.com/ocibuilder/ocibuilder/common"
	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"k8s.io/client-go/rest"
)

func main() {
	namespace := os.Getenv(common.Namespace)
	resource := os.Getenv(common.Resource)
	resourceName := os.Getenv(common.ResourceName)

	restConfig, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalln(err)
		return
	}

	client := ocbv1alpha1.NewForConfigOrDie(restConfig)

	ociObj, err := client.BlackrockV1alpha1().OCIBuilders(namespace).Get("my-r", metav1.GetOptions{})

	//
	//client := kubernetes.NewForConfigOrDie(restConfig).CoreV1().RESTClient()
	//req := client.Get().
	//	Namespace(namespace).
	//	Resource(resource).
	//	Name(resourceName)
	//
	//result, err := req.Do().Get()
	//if err != nil {
	//	log.Fatalln(err)
	//	return
	//}
	//obj := result.DeepCopyObject()
	//spec := obj.(*v1alpha1.OCIBuilder)
	//if err := storeBuilderSpecification(spec); err != nil {
	//	log.Fatalln(err)
	//	return
	//}
}

// storeBuilderSpecification stores the builder specification in a file
func storeBuilderSpecification(spec *v1alpha1.OCIBuilder) error {
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
