package common

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var builder = Builder("buildah")

func TestExec(t *testing.T) {
	executor = fakeExecCommand
	defer func() { executor = exec.Command }()

	_, err := builder.Command("build").Flags(Flag{
		Name:  "f",
		Value: "./Dockerfile",
		Short: true,
	}).Args(".").Build().Exec()
	assert.Equal(t, nil, err)
}

// enabling the mocking of exec commands as in https://npf.io/2015/06/testing-exec-command/
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
