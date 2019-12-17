package command

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

var executor = exec.Command

// Command is a single executable command
type Command struct {
	name    string
	command string
	flags   []Flag
	args    []string

	execCmd *exec.Cmd
}

// CommandBuilder is a builder for a Command
type CommandBuilder struct {
	name    string
	command string
	flags   []Flag
	args    []string
}

// Flag is a command flag
type Flag struct {
	// Name is the name of the flag
	Name string
	// Value is the value of the flag
	Value string
	// Short determines whether the flag used is a short variation or not
	Short bool
	// OmitEmpty omits the flag if the value is empty
	OmitEmpty bool
}

// Build builds a command from a CommandBuilder
func (builder *CommandBuilder) Build() Command {
	return Command{
		name:    builder.name,
		command: builder.command,
		flags:   builder.flags,
		args:    builder.args,
	}
}

// Command specifies the command to exec for the builder
func (builder *CommandBuilder) Command(command string) *CommandBuilder {
	builder.command = command
	return builder
}

// Flags specifies the flags to exec for the builder
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

// Args specifies the args to exec for the builder
func (builder *CommandBuilder) Args(args ...string) *CommandBuilder {
	builder.args = args
	return builder
}

// Builder initializes the CommandBuilder with default values
func Builder(name string) *CommandBuilder {
	cmdBuilder := new(CommandBuilder)
	cmdBuilder.args = make([]string, 0)
	cmdBuilder.command = ""
	cmdBuilder.flags = make([]Flag, 0)
	cmdBuilder.name = name
	return cmdBuilder
}

// Exec executes a command, returning readers for both stdout and stderr pipes
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

// Wait calls wait on a started exec command
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
