package commands

import (
	"fmt"

	"github.com/albertoleal/concourse-flake-hunter/fly"
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

		client := ctx.App.Metadata["client"].(fly.Client)

		searcher := hunter.NewSearcher(client)
		spec := hunter.SearchSpec{
			Pattern: ctx.Args().First(),
			Limit:   ctx.Int("limit"),
		}
		builds := searcher.Search(spec)

		fmt.Printf("+%-32s+%s\n", "----------------------------------", "-----------------------------------------------------")
		fmt.Printf("| %-32s | %s\n", "JOB", "URL")
		fmt.Printf("+%-32s+%s\n", "----------------------------------", "-----------------------------------------------------")

		for build := range builds {
			fmt.Printf("| %-32s | %s\n", build.PipelineName+"/"+build.JobName, build.ConcourseURL)
		}

		return nil
	},
}
