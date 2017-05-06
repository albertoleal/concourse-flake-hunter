package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/albertoleal/concourse-flake-hunter/hunter"
	"github.com/concourse/atc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Search", func() {
	var (
		server *ghttp.Server
		spec   hunter.SearchSpec
		builds []atc.Build
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		concourseURL = fmt.Sprintf("http://%s/", server.Addr())
		builds = []atc.Build{
			atc.Build{
				ID:           1,
				TeamName:     "awesome-team",
				Name:         "my-build",
				PipelineName: "pipeline",
				JobName:      "job",
				URL:          "url"},
		}

		mockConcourseAPICalls(server, builds)

		spec = hunter.SearchSpec{
			Pattern: "connection reset",
		}
	})

	AfterEach(func() {
		server.Close()
	})

	It("searches for flakes on the last builds", func() {
		output, err := Runner.Search(spec)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(ContainSubstring("pipeline/job"))
		Expect(output).To(ContainSubstring(fmt.Sprintf("%s%s", concourseURL, builds[0].URL)))
	})
})

func mockConcourseAPICalls(server *ghttp.Server, builds []atc.Build) {
	bs, err := json.Marshal(builds)
	Expect(err).NotTo(HaveOccurred())

	server.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, "/api/v1/teams/awesome-team/auth/token"),
			ghttp.RespondWith(200, `{"type":"Bearer","value":"token"}`),
		),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, "/api/v1/builds"),
			ghttp.RespondWith(200, bs),
		),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, "/api/v1/teams/awesome-team/auth/token"),
			ghttp.RespondWith(200, `{"type":"Bearer","value":"token"}`),
		),
		ghttp.CombineHandlers(
			ghttp.VerifyRequest(http.MethodGet, "/api/v1/builds/1/events"),
			ghttp.RespondWith(200, `{"origin":{"id":"58f5f81a", "source":"stdout"}, "payload":"using version of resource found in cache"}`),
		),
	)
}
