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

package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalContext_Read(t *testing.T) {
	localContext := LocalContext{
		ContextPath: "../../testing",
	}
	_, _, err := localContext.Read()
	assert.Equal(t, nil, err)
}

func TestLocalContext_Read2(t *testing.T) {
	localContext := LocalContext{
		ContextPath: "",
	}
	_, _, err := localContext.Read()
	assert.Error(t, err, "cannot have empty contextPath: specify . for current directory")
}
