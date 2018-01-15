package hunter_test

import (
	"errors"
	"fmt"

	"github.com/masters-of-cats/concourse-flake-hunter/fly/flyfakes"
	"github.com/masters-of-cats/concourse-flake-hunter/hunter"
	"github.com/concourse/atc"
	"github.com/concourse/go-concourse/concourse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Searcher", func() {
	var (
		fakeClient *flyfakes.FakeClient
		searcher   *hunter.Searcher
		builds     []atc.Build
	)

	BeforeEach(func() {
		builds = []atc.Build{
			atc.Build{ID: 1, Name: "my-build", URL: "/build-url", Status: string(atc.StatusSucceeded)},
			atc.Build{ID: 2, Name: "my-build", URL: "/build-url", Status: string(atc.StatusErrored)},
			atc.Build{ID: 3, Name: "my-build", URL: "/build-url", Status: string(atc.StatusFailed)},
		}

		fakeClient = new(flyfakes.FakeClient)
		searcher = hunter.NewSearcher(fakeClient)

		fakeClient.ConcourseURLReturns("https://concourse.io")
		fakeClient.BuildsReturns(builds, concourse.Pagination{}, nil)
		fakeClient.BuildEventsReturns([]byte("connection reset"), nil)
	})

	Describe("Search", func() {
		It("searches for last builds", func() {
			bs, err := searcher.Search(hunter.SearchSpec{
				Pattern: "connection reset",
			})
			Expect(err).NotTo(HaveOccurred())
			concourseURL := fmt.Sprintf("%s%s", fakeClient.ConcourseURL(), builds[0].URL)
			Expect(bs[0].ConcourseURL).To(Equal(concourseURL))
		})

		It("searches only for succeeded and failed builds", func() {
			_, err := searcher.Search(hunter.SearchSpec{
				Pattern: "connection reset",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(fakeClient.BuildEventsCallCount()).To(Equal(2))
		})

		Context("when faild to get last builds", func() {
			BeforeEach(func() {
				fakeClient.BuildsReturns(nil, concourse.Pagination{}, errors.New("failed to get builds"))
			})

			It("returns an error", func() {
				_, err := searcher.Search(hunter.SearchSpec{
					Pattern: "connection reset",
				})
				Expect(err).To(MatchError(ContainSubstring("failed to get builds")))
			})
		})
	})

})
