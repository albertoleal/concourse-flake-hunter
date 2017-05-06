package hunter_test

import (
	"errors"
	"fmt"

	"github.com/albertoleal/concourse-flake-hunter/hunter"
	"github.com/albertoleal/concourse-flake-hunter/hunter/hunterfakes"
	"github.com/concourse/atc"
	"github.com/concourse/go-concourse/concourse"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Searcher", func() {
	var (
		fakeClient *hunterfakes.FakeClient
		searcher   *hunter.Searcher
		builds     []atc.Build
	)

	BeforeEach(func() {
		builds = []atc.Build{
			atc.Build{ID: 1, Name: "my-build", URL: "/build-url"},
		}

		fakeClient = new(hunterfakes.FakeClient)
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
			Expect(fakeClient.BuildsCallCount()).To(Equal(1))
			concourseURL := fmt.Sprintf("%s%s", fakeClient.ConcourseURL(), builds[0].URL)
			Expect(bs[0].ConcourseURL).To(Equal(concourseURL))
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
