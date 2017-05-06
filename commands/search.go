package commands

import (
	"fmt"
	"os"

	"github.com/albertoleal/concourse-flake-hunter/hunter"
	"github.com/urfave/cli"
)

var SearchCommand = cli.Command{
	Name:        "search",
	Usage:       "search <arguments>",
	Description: "Searches for flakes",

	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "pattern, p",
			Usage: "Flake pattern (i.e: connection reset)",
		},
	},

	Action: func(ctx *cli.Context) error {
		client := ctx.App.Metadata["client"].(hunter.Client)

		searcher := hunter.NewSearcher(client)
		spec := hunter.SearchSpec{
			Pattern: ctx.Args().First(),
		}
		builds, err := searcher.Search(spec)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		table := &Table{
			Content: [][]string{},
			Header:  []string{"pipeline/job", "build url"},
		}

		for _, build := range builds {
			line := []string{}
			line = append(line, fmt.Sprintf("%s/%s", build.PipelineName, build.JobName))
			line = append(line, fmt.Sprintf("%s", build.ConcourseURL))
			table.Content = append(table.Content, line)
		}

		context := &Context{Stdout: os.Stdout}
		table.Render(context)
		return nil
	},
}
