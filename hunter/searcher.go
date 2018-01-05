package hunter

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/albertoleal/concourse-flake-hunter/fly"
	"github.com/concourse/atc"
	"github.com/concourse/go-concourse/concourse"
)

const (
	STATUS_FORBIDDEN = "forbidden"
)

type SearchSpec struct {
	Pattern string
	Limit   int
}

type Searcher struct {
	client fly.Client
}

type Build struct {
	atc.Build

	ConcourseURL string
}

func NewSearcher(client fly.Client) *Searcher {
	return &Searcher{
		client: client,
	}
}

func (s *Searcher) Search(spec SearchSpec) chan Build {
	ch := make(chan Build, 100)
	go s.getBuildsFromPage(ch, concourse.Page{Limit: 300}, spec)
	return ch
}

func (s *Searcher) getBuildsFromPage(flakesChan chan Build, page concourse.Page, spec SearchSpec) {
	var (
		buildsChan = make(chan []atc.Build, 2)
		pages      = concourse.Pagination{Next: &page}
		builds     []atc.Build
		err        error
	)

	go s.processBuilds(flakesChan, buildsChan, spec)
	go s.processBuilds(flakesChan, buildsChan, spec)

	for ; pages.Next != nil; page = *pages.Next {
		builds, pages, err = s.client.Builds(page)
		if err != nil {
			panic(err)
		}

		buildsChan <- builds

		time.Sleep(time.Second)
	}
}

func (s *Searcher) processBuilds(flakesCh chan Build, buildsCh chan []atc.Build, spec SearchSpec) {
	for builds := range buildsCh {
		for _, build := range builds {
			if build.Status != string(atc.StatusFailed) {
				continue
			}

			events, err := s.client.BuildEvents(strconv.Itoa(build.ID))
			if err != nil {
				// Not sure why, but concourse.Builds returns builds from other teams
				if err.Error() == STATUS_FORBIDDEN {
					continue
				}
				println(err.Error())
				continue
			}

			ok, err := regexp.Match(spec.Pattern, events)
			if err != nil {
				println("failed trying to match pattern", err.Error())
				continue
			}

			if ok {
				concourseURL := fmt.Sprintf("%s%s", s.client.ConcourseURL(), build.URL)
				b := Build{build, concourseURL}
				flakesCh <- b
			}
		}
	}
}
