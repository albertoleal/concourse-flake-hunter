package commands

import (
	"fmt"
	"regexp"
	"time"

	"github.com/masters-of-cats/concourse-flake-hunter/fly"
	"github.com/masters-of-cats/concourse-flake-hunter/hunter"
	"github.com/urfave/cli"
)

type Failure struct {
	Description  string
	JobName      string
	Date         int64
	ConcourseURL string
}

type FailuresInfo struct {
	Count         int
	LastOccurance int64
	Failures      []Failure
}

var AggregateCommand = cli.Command{
	Name:        "aggregate",
	Usage:       "aggregate",
	Description: "Aggregates flakes",

	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "max-age, m",
			Usage: "Lists builds that failed in the last n hours",
			Value: -1,
		},
		cli.IntFlag{
			Name:  "threshold, t",
			Usage: "Defines a test as a flake it if fails at least n times",
			Value: 3,
		},
	},

	Action: func(ctx *cli.Context) error {
		client := ctx.App.Metadata["client"].(fly.Client)

		searcher := hunter.NewSearcher(client)
		spec := hunter.SearchSpec{
			Pattern: regexp.MustCompile("\\[Fail\\].*"),
		}

		builds := searcher.Search(spec)

		aggregator := NewAggregator()
		maxAge := ctx.Int("max-age")
		for build := range builds {
			if maxAge > 0 && age(build) > maxAge {
				break
			}
			for _, match := range build.Matches {
				aggregator.addFailure(&Failure{
					Description:  match,
					JobName:      build.JobName,
					ConcourseURL: build.ConcourseURL,
					Date:         build.StartTime,
				})
			}
		}

		aggregator.printEntries(ctx.Int("threshold"))
		return nil
	},
}

type Aggregator struct {
	failuresInfo map[string]*FailuresInfo
}

func NewAggregator() *Aggregator {
	return &Aggregator{failuresInfo: make(map[string]*FailuresInfo)}
}

func (a *Aggregator) addFailure(failure *Failure) {
	if info, ok := a.failuresInfo[failure.Description]; ok {
		info.Count++
		info.Failures = append(info.Failures, *failure)
		if info.LastOccurance < failure.Date {
			info.LastOccurance = failure.Date
		}
	} else {
		a.failuresInfo[failure.Description] = &FailuresInfo{
			Count:         1,
			LastOccurance: failure.Date,
			Failures:      []Failure{*failure},
		}
	}
}

func (a *Aggregator) printEntries(threshold int) {
	for description, info := range a.failuresInfo {
		if info.Count < threshold {
			continue
		}
		fmt.Println(description)
		fmt.Printf("\tCount: %d\n", info.Count)
		fmt.Printf("\tLastOccurance: %s\n", time.Unix(info.LastOccurance, 0).String())
		for _, failure := range info.Failures {
			fmt.Printf("\t\tJobName: %s\n", failure.JobName)
			fmt.Printf("\t\tDate: %s\n", time.Unix(failure.Date, 0).String())
			fmt.Printf("\t\tURL: %s\n\n", failure.ConcourseURL)
		}
		fmt.Printf("---------------\n")
	}
}
