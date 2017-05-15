package hunter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

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
	client Client
}

type Build struct {
	atc.Build

	ConcourseURL string
}

func NewSearcher(client Client) *Searcher {
	return &Searcher{
		client: client,
	}
}

func (s *Searcher) Search(spec SearchSpec) ([]Build, error) {
	builds, _, err := s.client.Builds(concourse.Page{Limit: spec.Limit})
	if err != nil {
		return nil, err
	}

	bs := []Build{}
	count := 1
	for _, build := range builds {
		if build.Status != string(atc.StatusSucceeded) && build.Status != string(atc.StatusFailed) {
			continue
		}

		events, err := s.client.BuildEvents(strconv.Itoa(build.ID))
		if err != nil {
			// Not sure why, but concourse.Builds returns builds from other teams
			if err.Error() == STATUS_FORBIDDEN {
				continue
			}
			return []Build{}, err
		}

		ok, err := regexp.Match(spec.Pattern, events)
		if err != nil {
			return []Build{}, errors.New("failed trying to match pattern given")
		}

		if ok {
			concourseURL := fmt.Sprintf("%s%s", s.client.ConcourseURL(), build.URL)
			b := Build{build, concourseURL}
			bs = append(bs, b)
		}
		//Debug: Still experimenting with this - I saw concourse failing when count is bigger than 500.
		if count%100 == 0 {
			time.Sleep(10 * time.Second)
		}
		count = count + 1
	}
	return bs, nil
}
