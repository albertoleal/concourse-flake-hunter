package runner

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Runner struct {
	ConcourseFlakeHunter string
	ConcourseURL         string
	Username             string
	Password             string
	Team                 string
}

func (r Runner) RunSubcommand(subcommand string, args ...string) (string, error) {
	stdoutBuffer := bytes.NewBuffer([]byte{})
	cmd := r.makeCmd(subcommand, args)
	cmd.Stdout = stdoutBuffer

	if err := cmd.Run(); err != nil {
		fmt.Printf("stdoutBuffer.String() = %+v\n", stdoutBuffer.String())
		return "", err
	}

	return strings.TrimSpace(stdoutBuffer.String()), nil
}

func (r Runner) makeCmd(subcommand string, args []string) *exec.Cmd {
	allArgs := []string{"--concourse-url", r.ConcourseURL, "--username",
		r.Username, "--password", r.Password, "--team-name", r.Team}

	allArgs = append(allArgs, subcommand)
	allArgs = append(allArgs, args...)

	return exec.Command(r.ConcourseFlakeHunter, allArgs...)
}
