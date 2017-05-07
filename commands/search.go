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
		cli.IntFlag{
			Name:  "limit, l",
			Usage: "Limit number of builds to check",
			Value: 50,
		},
	},

	Action: func(ctx *cli.Context) error {
		if ctx.Args().First() == "" {
			return cli.NewExitError("need to provide a pattern", 1)
		}

		client := ctx.App.Metadata["client"].(hunter.Client)

		searcher := hunter.NewSearcher(client)
		spec := hunter.SearchSpec{
			Pattern: ctx.Args().First(),
			Limit:   ctx.Int("limit"),
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
