/*
Copyright © 2019 BlackRock Inc.

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

package dummy

import "github.com/ocibuilder/ocibuilder/pkg/apis/ocibuilder/v1alpha1"

var Spec = v1alpha1.OCIBuilderSpec{
	Build: BuildSpec,
	Login: LoginSpec,
	Push:  PushSpec,
}

var BuildSpec = &v1alpha1.BuildSpec{
	Steps: []v1alpha1.BuildStep{
		{
			Metadata: &v1alpha1.Metadata{
				Name: "test-build",
			},
			Stages: []v1alpha1.Stage{
				{
					Base: v1alpha1.Base{
						Image: "alpine",
					},
					Metadata: &v1alpha1.Metadata{
						Name: "stage-one",
					},
					Cmd: []v1alpha1.BuildTemplateStep{
						{
							Docker: &v1alpha1.DockerStep{
								Inline: []string{"echo", "done"},
							},
						},
					},
				},
			},
		},
	},
}

var LoginSpec = []v1alpha1.LoginSpec{
	{
		Registry: "example-registry",
		Token:    "ThiSiSalOgInToK3N",
		Creds: v1alpha1.RegistryCreds{
			Plain: v1alpha1.PlainCreds{
				Username: "username",
				Password: "password",
			},
		},
	},
}

var PushSpec = []v1alpha1.PushSpec{
	{
		Registry: "example-registry",
		Image:    "example-image",
		User:     "namespace",
		Tag:      "1.0.0",
	},
}
