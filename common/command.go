package common

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

var executor = exec.Command

type Command struct {
	name    string
	command string
	flags   []Flag
	args    []string

	execCmd *exec.Cmd
}

type CommandBuilder struct {
	name    string
	command string
	flags   []Flag
	args    []string
}

type Flag struct {
	Name  string
	Value string
	Short bool
}

func (builder *CommandBuilder) Build() Command {
	return Command{
		name:    builder.name,
		command: builder.command,
		flags:   builder.flags,
		args:    builder.args,
	}
}

func (builder *CommandBuilder) Command(command string) *CommandBuilder {
	builder.command = command
	return builder
}

func (builder *CommandBuilder) Flags(flags ...Flag) *CommandBuilder {
	builder.flags = flags
	return builder
}

func (builder *CommandBuilder) Args(args ...string) *CommandBuilder {
	builder.args = args
	return builder
}

func Builder(name string) *CommandBuilder {
	cmdBuilder := new(CommandBuilder)
	cmdBuilder.args = make([]string, 0)
	cmdBuilder.command = ""
	cmdBuilder.flags = make([]Flag, 0)
	cmdBuilder.name = name
	return cmdBuilder
}

func (c *Command) Exec() (io.ReadCloser, error) {
	command := c.constructCommand()
	cmd := executor(c.name, command...)
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	c.execCmd = cmd
	return stdout, nil
}

func (c Command) Wait() error {
	if err := c.execCmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code is %d\n", exitError.ExitCode())
			errorString := fmt.Sprintf("error in executing cmd, exited with code %d", exitError.ExitCode())
			return errors.New(errorString)
		}
	}
	return nil
}

func (c Command) constructCommand() []string {
	var commandVector = []string{c.command}

	for _, flag := range c.flags {
		if flag.Short {
			commandVector = append(commandVector, fmt.Sprintf("-%s", flag.Name), flag.Value)
		} else {
			commandVector = append(commandVector, fmt.Sprintf("--%s", flag.Name), flag.Value)
		}
	}

	for _, arg := range c.args {
		commandVector = append(commandVector, arg)
	}

	return commandVector
}
