package request

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/beval/beval/pkg/apis/beval/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestRequestRemoteNoAuth(t *testing.T) {
	url := "https://raw.githubusercontent.com/beval/beval/master/testing/dummy/overlay_overlay_test.yaml"
	filepath := "../../testing/dummy/downloaded_overlay.yaml"

	defer os.Remove(filepath)

	err := RequestRemote(url, filepath, v1alpha1.RemoteCreds{})
	assert.Equal(t, nil, err)

	actualFile, err := ioutil.ReadFile(filepath)
	assert.Equal(t, nil, err)

	expectedFile, err := ioutil.ReadFile("../../testing/dummy/overlay_overlay_test.yaml")
	assert.Equal(t, nil, err)

	assert.Equal(t, string(expectedFile), string(actualFile))
}
