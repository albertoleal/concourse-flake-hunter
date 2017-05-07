package runner

import (
	"fmt"
	"strconv"

	"github.com/albertoleal/concourse-flake-hunter/hunter"
)

func (r Runner) Search(spec hunter.SearchSpec) (string, error) {
	args := []string{"--limit", strconv.Itoa(spec.Limit)}
	if spec.Pattern != "" {
		args = append(args, fmt.Sprintf("\"%s\"", spec.Pattern))
	}

	result, err := r.RunSubcommand("search", args...)
	if err != nil {
		return err.Error(), err
	}

	return result, nil
}
