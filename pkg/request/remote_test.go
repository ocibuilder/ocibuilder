package request

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"
)

func TestRequestRemoteNoAuth(t *testing.T) {
	url := "https://raw.githubusercontent.com/ocibuilder/ocibuilder/master/testing/dummy/overlay_overlay_test.yaml"
	filepath := "../../testing/dummy/downloaded_overlay.yaml"

	defer os.Remove(filepath)

	err := RequestRemote(url, filepath, types.AuthConfig{})
	assert.Equal(t, nil, err)

	actualFile, err := ioutil.ReadFile(filepath)
	assert.Equal(t, nil, err)

	expectedFile, err := ioutil.ReadFile("../../testing/dummy/overlay_overlay_test.yaml")
	assert.Equal(t, nil, err)

	assert.Equal(t, string(expectedFile), string(actualFile))
}
