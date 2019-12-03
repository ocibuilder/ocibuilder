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

package command

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
	// Short determines whether the flag used is a short variation or not
	Short bool
	// OmitEmpty omits the flag if the value is empty
	OmitEmpty bool
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
	var builderFlags []Flag
	for _, f := range flags {
		if !f.OmitEmpty {
			builderFlags = append(builderFlags, f)
		} else if f.Value != "" {
			builderFlags = append(builderFlags, f)
		}
	}

	builder.flags = builderFlags
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

func (c *Command) Exec() (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	command := c.constructCommand()
	cmd := executor(c.name, command...)
	stdout, _ = cmd.StdoutPipe()
	stderr, _ = cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	c.execCmd = cmd
	return stdout, stderr, nil
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

	return append(commandVector, c.args...)
}
