package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerGenerator_Generate(t *testing.T) {
	dockerGen := DockerGenerator{
		Filepath: "../../testing/dummy/Dockerfile_Test",
	}
	_, err := dockerGen.Generate()
	assert.Equal(t, nil, err)

}
