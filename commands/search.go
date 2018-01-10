package commands

import (
	"fmt"
	"time"

	"github.com/albertoleal/concourse-flake-hunter/fly"
	"github.com/albertoleal/concourse-flake-hunter/hunter"
	"github.com/urfave/cli"
)

var SearchCommand = cli.Command{
	Name:        "search",
	Usage:       "search <arguments>",
	Description: "Searches for flakes",

	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "show-one-offs",
			Usage: "If set one off failures will be reported as well",
		},
		cli.IntFlag{
			Name:  "max-age, m",
			Usage: "Lists builds that failed in the last n hours",
			Value: -1,
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
		}

		if ctx.Bool("show-one-offs") {
			spec.ShowOneOffs = true
		}
		builds := searcher.Search(spec)

		fmt.Printf("+-------+%-32s+%s\n", "----------------------------------", "-----------------------------------------------------")
		fmt.Printf("| %-5s | %-32s | %s\n", "Ended", "Job", "Url")
		fmt.Printf("+-------+%-32s+%s\n", "----------------------------------", "-----------------------------------------------------")

		maxAge := ctx.Int("max-age")
		for build := range builds {
			if maxAge > 0 && age(build) > maxAge {
				break
			}

			fmt.Printf("| %-5s | %-32s | %s\n", timeSince(build.EndTime), build.PipelineName+"/"+build.JobName, build.ConcourseURL)
		}

		return nil
	},
}

func age(build hunter.Build) int {
	endTime := time.Unix(build.EndTime, 0)
	return int(time.Since(endTime) / time.Hour)
}

func timeSince(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	timeSince := time.Since(t)

	hoursSince := timeSince / time.Hour
	if hoursSince < 1 {
		return fmt.Sprintf("%dm", timeSince/time.Minute)
	}
	if hoursSince < 24 {
		return fmt.Sprintf("%dh", hoursSince)
	}
	return fmt.Sprintf("%dd", hoursSince/24)
}
