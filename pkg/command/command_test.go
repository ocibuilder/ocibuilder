package command

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var builder = Builder("buildah")

func TestCommand_constructCommand(t *testing.T) {
	constructedCmd := cmd.constructCommand()

	assert.Equal(t, []string{"build", "-f", "./Dockerfile", "--storage-driver", "vfs", ".", "one", "two"}, constructedCmd)
}

func TestCommand_Exec(t *testing.T) {
	executor = fakeExecCommand
	defer func() { executor = exec.Command }()

	_, _, err := cmd.Exec()
	assert.Equal(t, nil, err)
}

func TestCommandBuilder_Flags(t *testing.T) {

	flags := []Flag{
		{"f", "Dockerfile", true, true},
		{"storage-driver", "", false, true},
		{"t", "image:tag", true, true},
	}

	builder := Builder("test").Flags(flags...)
	assert.Equal(t, Command{
		name:    "test",
		command: "",
		flags:   expectedFlags,
		args:    []string{},
	}, builder.Build())
}

var expectedFlags = []Flag{
	{"f", "Dockerfile", true, true},
	{"t", "image:tag", true, true},
}

var cmd = builder.Command("build").Flags([]Flag{
	{
		Name:  "f",
		Value: "./Dockerfile",
		Short: true,
	},
	{
		Name:  "storage-driver",
		Value: "vfs",
		Short: false,
	},
}...).Args(".", "one", "two").Build()

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

func TestCommand_constructCommand_emptyCommand(t *testing.T) {
	flags := []Flag{
		{"testFlag", "flagValue", false, false},
	}
	command := Builder("test").Flags(flags...).Args("testArg").Build()
	commandVector := command.constructCommand()

	expectedCommandVector := []string{"--testFlag", "flagValue", "testArg"}
	assert.Equal(t, expectedCommandVector, commandVector)
}
