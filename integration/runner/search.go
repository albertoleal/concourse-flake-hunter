package runner

import (
	"fmt"

	"github.com/albertoleal/concourse-flake-hunter/hunter"
)

func (r Runner) Search(spec hunter.SearchSpec) (string, error) {
	args := []string{"--pattern", fmt.Sprintf("\"%s\"", spec.Pattern)}

	result, err := r.RunSubcommand("search", args...)
	if err != nil {
		return "", err
	}

	return result, nil
}
