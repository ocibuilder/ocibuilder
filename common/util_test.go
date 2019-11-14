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

package common

import (
	"os"
	"testing"

	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/smartystreets/goconvey/convey"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestReadCredentials(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()

	convey.Convey("Given the credentials in plain text, read the secret", t, func() {
		creds := &v1alpha1.Credentials{
			Plain: "secret",
		}
		secret, err := ReadCredentials(fakeClient, creds)
		convey.So(err, convey.ShouldBeNil)
		convey.So(secret, convey.ShouldEqual, "secret")
	})

	convey.Convey("Given the credentials in environment variable, read the secret", t, func() {
		err := os.Setenv("SECRET", "secret")
		convey.So(err, convey.ShouldBeNil)
		creds := &v1alpha1.Credentials{
			Env: "SECRET",
		}
		secret, err := ReadCredentials(fakeClient, creds)
		convey.So(err, convey.ShouldBeNil)
		convey.So(secret, convey.ShouldEqual, "secret")
	})

	convey.Convey("Given the credentials in Kubernetes secret", t, func() {
		secretObj := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "fake",
			},
			Data: map[string][]byte{
				"accessKey": []byte("access"),
			},
		}
		sensor, err := fakeClient.CoreV1().Secrets(secretObj.Namespace).Create(secretObj)
		convey.So(err, convey.ShouldBeNil)
		convey.So(sensor.Data["accessKey"], convey.ShouldEqual, "access")
	})
}
