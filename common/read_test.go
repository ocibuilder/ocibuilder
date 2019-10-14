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
	"io/ioutil"
	"os"
	"testing"

	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestReadLogin(t *testing.T) {
	spec := v1alpha1.OCIBuilderSpec{}
	err := Read(&spec, "", "../testing")

	assert.Equal(t, nil, err)
	assert.Equal(t, spec.Login, loginSpec, "the login spec to match the expected")
}

func TestApplyParams(t *testing.T) {
	spec := v1alpha1.OCIBuilderSpec{}
	file, err := ioutil.ReadFile("../testing/spec_read_test.yaml")

	spec.Login = loginSpec
	spec.Params = params

	expectedLogin := v1alpha1.EnvCreds{
		Username: "testuser",
		Password: "my-real-password",
	}

	assert.Equal(t, nil, err)
	err = applyParams(file, &spec)
	assert.Equal(t, nil, err)

	assert.Equal(t, expectedLogin, spec.Login[0].Creds.Env)
}

func TestApplyParamsEnvVariable(t *testing.T) {
	spec := v1alpha1.OCIBuilderSpec{}
	file, err := ioutil.ReadFile("../testing/spec_read_test.yaml")

	spec.Login = loginSpec
	spec.Params = paramsEnv

	os.Setenv("$TEST_USERNAME", "test_env_user")
	os.Setenv("$TEST_PASSWORD", "test_env_pass")

	expectedLoginEnv := v1alpha1.EnvCreds{
		Username: "test_env_user",
		Password: "test_env_pass",
	}

	assert.Equal(t, nil, err)
	err = applyParams(file, &spec)

	assert.Equal(t, expectedLoginEnv, spec.Login[0].Creds.Env)

	os.Remove("$TEST_USERNAME")
	os.Remove("$TEST_PASSWORD")
}

var loginSpec = []v1alpha1.LoginSpec{{
	Registry: "example-registry",
	Creds: v1alpha1.RegistryCreds{
		Env: v1alpha1.EnvCreds{
			Username: "art",
			Password: "my-real-password",
		},
		K8s: v1alpha1.K8sCreds{
			Username: nil,
			Password: nil,
		},
		Plain: v1alpha1.PlainCreds{
			Username: "user",
			Password: "pass",
		},
	},
}, {
	Registry: "example-registry-2",
	Creds: v1alpha1.RegistryCreds{
		Plain: v1alpha1.PlainCreds{
			Username: "user2",
			Password: "pass2",
		},
	},
}}

var params = []v1alpha1.Param{{
	Dest:  "login.0.creds.env.username",
	Value: "testuser",
}, {
	Dest:  "login.0.creds.env.password",
	Value: "my-real-password",
}}

var paramsEnv = []v1alpha1.Param{{
	Dest:                 "login.0.creds.env.username",
	ValueFromEnvVariable: "$TEST_USERNAME",
}, {
	Dest:                 "login.0.creds.env.password",
	ValueFromEnvVariable: "$TEST_PASSWORD",
}}
