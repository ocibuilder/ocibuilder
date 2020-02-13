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

package validate

import (
	"os"

	"github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"
	"github.com/ocibuilder/ocibuilder/pkg/common"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// Validate validates a ocibuilder spec.
func Validate(spec *v1alpha1.OCIBuilderSpec) error {
	if spec == nil {
		return errors.New("builder spec can't be nil")
	}
	if spec.Login == nil && spec.Push == nil && spec.Build == nil {
		return errors.New("all builder specs can't be nil")
	}
	if spec.Login == nil && spec.Push != nil {
		return errors.New("at least one login must be provided")
	}
	return nil
}

// ValidateBuildTemplateStep validates build template step
func ValidateBuildTemplateStep(step v1alpha1.BuildTemplateStep) error {
	if step.Ansible == nil && step.Docker == nil {
		return errors.New("at least one step type should be defined")
	}
	if step.Docker != nil && step.Docker.Inline == nil && step.Docker.Path == "" && step.Docker.Url == "" {
		return errors.New("at least one docker cmd location should be defined")
	}
	return nil
}

// ValidateLoginUsername validates the login spec for a username, and returns the first username found
func ValidateLoginUsername(spec v1alpha1.LoginSpec) (string, error) {
	if spec.Creds.Plain.Username != "" {
		return spec.Creds.Plain.Username, nil
	}
	if spec.Creds.Env.Username != "" {
		return os.Getenv(spec.Creds.Env.Username), nil
	}
	if spec.Creds.K8s.Username != nil {
		return spec.Creds.K8s.Username.Key, nil
	}
	return "", errors.New("at least one login username must be specified")
}

// ValidateLoginPassword validates the login spec for a password, and returns the first password found
func ValidateLoginPassword(spec v1alpha1.LoginSpec) (string, error) {
	if spec.Token != "" {
		return spec.Token, nil
	}
	if spec.Creds.Plain.Password != "" {
		return spec.Creds.Plain.Password, nil
	}
	if spec.Creds.Env.Password != "" {
		return os.Getenv(spec.Creds.Env.Password), nil
	}
	if spec.Creds.K8s.Password != nil {
		return spec.Creds.K8s.Password.Key, nil
	}
	return "", errors.New("at least one login password must be specified")
}

// ValidateLogin validates the top level login specification
func ValidateLogin(spec v1alpha1.OCIBuilderSpec) error {
	if spec.Login == nil {
		return errors.New("at least one login must be provided")
	}
	return nil
}

// ValidatePush validates the top level push specification
func ValidatePush(spec v1alpha1.OCIBuilderSpec) error {
	if spec.Push == nil {
		return errors.New("at least one push spec must be provided")
	}
	return nil
}

// ValidatePushSpec validates the lower level push specification
func ValidatePushSpec(spec *v1alpha1.PushSpec) error {
	if spec.Registry == "" {
		spec.Registry = common.DefaultImageRegistry
	}
	if spec.Image == "" {
		return errors.New("image name must be specified for push")
	}
	if spec.Tag == "" {
		return errors.New("tag must be specified for push")
	}
	return nil
}

// ValidateParams validates path to destination in param section of specs
func ValidateParams(specJSON []byte, src string) error {
	if res := gjson.GetBytes(specJSON, src); res.Str == "" {
		return errors.New("path to dest is invalid in a set param")
	}
	return nil
}

// ValidateContext validates the build context
func ValidateContext(spec *v1alpha1.BuildContext) error {

	if spec == nil {
		return errors.New("no build context has been specified in the build configuration")
	}

	if spec.AliyunOSSContext == nil && spec.AzureBlobContext == nil && spec.GCSContext == nil && spec.GitContext == nil && spec.LocalContext == nil && spec.S3Context == nil {
		return errors.New("no build context has been specified in the build configuration")
	}

	return nil
}

// SetAnsibleDefaultIfNotPresent updates default values if not present
func SetAnsibleDefaultIfNotPresent(spec *v1alpha1.AnsibleStep) error {
	if spec.Playbook == "" {
		spec.Playbook = "playbook.yaml"
	}

	if spec.Requirements == "" {
		spec.Requirements = "requirements.yaml"
	}

	if spec.Workspace == "" {
		return errors.New("ansible workspace name has not been set in ansible step")
	}

	return nil
}
