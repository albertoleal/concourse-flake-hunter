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

func (s *Searcher) getBuildsFromPage(ch chan Build, page concourse.Page, spec SearchSpec) {
	before := time.Now()
	builds, pages, err := s.client.Builds(page)
	duration := time.Since(before)
	fmt.Println("Response time", duration/time.Millisecond)

	if err != nil {
		panic(err)
	}

	go s.processBuilds(ch, builds, spec)

	if pages.Next == nil {
		fmt.Println("No More Pages")
		return
	}

	time.Sleep(time.Second)

	s.getBuildsFromPage(ch, *pages.Next, spec)
	if err != nil {
		panic(err)
	}
}

func (s *Searcher) processBuilds(ch chan Build, builds []atc.Build, spec SearchSpec) {
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
			ch <- b
		}
	}
}
