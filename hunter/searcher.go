package hunter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/concourse/atc"
	"github.com/concourse/go-concourse/concourse"
)

type SearchSpec struct {
	Pattern string
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
	builds, _, err := s.client.Builds(concourse.Page{Limit: 100})
	if err != nil {
		return nil, err
	}

	bs := []Build{}
	for _, build := range builds {
		events, err := s.client.BuildEvents(strconv.Itoa(build.ID))
		if err != nil {
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
	}
	return bs, nil
}
