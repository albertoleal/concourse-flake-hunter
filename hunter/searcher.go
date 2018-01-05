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
	start  time.Time
}

type Build struct {
	atc.Build

	ConcourseURL string
}

func NewSearcher(client fly.Client) *Searcher {
	return &Searcher{
		client: client,
		start:  time.Now(),
	}
}

func (s *Searcher) Search(spec SearchSpec) chan Build {
	ch := make(chan Build, 100)
	go s.getBuildsFromPage(ch, concourse.Page{Limit: 300}, spec)
	return ch
}

func (s *Searcher) getBuildsFromPage(flakesChan chan Build, page concourse.Page, spec SearchSpec) {
	var (
		buildsChan = make(chan []atc.Build)
		pages      = concourse.Pagination{Next: &page}
		builds     []atc.Build
		err        error
	)

	go s.processBuilds(flakesChan, buildsChan, spec)
	go s.processBuilds(flakesChan, buildsChan, spec)
	go s.processBuilds(flakesChan, buildsChan, spec)
	go s.processBuilds(flakesChan, buildsChan, spec)
	go s.processBuilds(flakesChan, buildsChan, spec)
	go s.processBuilds(flakesChan, buildsChan, spec)

	for i := 0; pages.Next != nil; page, i = *pages.Next, i+1 {
		builds, pages, err = s.client.Builds(page)
		// retrier := retrier.New(retrier.ConstantBackoff(2, 5*time.Second), nil)

		// err = retrier.Run(func() error {
		// 	return err
		// })

		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Backing off")
			time.Sleep(time.Second * 5)
			continue
		}

		buildsretrier := retrier.New(retrier.ConstantBackoff(20, 500*time.Millisecond), nil)
var (
  pid = -1
  err error
)
retrier.Run(func() error {
  pid, err = parsePid(pidFilePath)
  return err
})Chan <- builds
		fmt.Println("batch", i, "builds processed", i*300, "time elapsed ", time.Since(s.start))

		if i > 0 && i%30 == 0 {
			// fmt.Println("Backing off")
			// time.Sleep(5 * time.Second)
		}

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
