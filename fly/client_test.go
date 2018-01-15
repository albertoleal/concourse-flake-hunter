package fly_test

import (
	"fmt"
	"net/http"

	"github.com/albertoleal/concourse-flake-hunter/fly"
	"github.com/concourse/go-concourse/concourse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		server *ghttp.Server

		client       fly.Client
		concourseURL string
		username     string
		password     string
		team         string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()

		concourseURL = fmt.Sprintf("http://%s/", server.Addr())
		username = "alice"
		password = "my-passwd"
		team = "awesome-team"
	})

	JustBeforeEach(func() {
		client = fly.NewClient(concourseURL, username, password, team)
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Builds", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/api/v1/teams/awesome-team/auth/token"),
					ghttp.RespondWith(200, `{"type":"Bearer","value":"token"}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/api/v1/builds"),
					ghttp.RespondWith(200, `[{"id":1,"team_name":"awesome-team","name":"my-job"}]`),
				),
			)
		})

		It("returns a list of builds", func() {
			builds, _, err := client.Builds(concourse.Page{Limit: 5})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(builds)).To(Equal(1))
			Expect(builds[0].ID).To(Equal(1))
			Expect(builds[0].TeamName).To(Equal("awesome-team"))
			Expect(builds[0].Name).To(Equal("my-job"))
		})
	})

	PDescribe("BuildEvents", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/api/v1/teams/awesome-team/auth/token"),
					ghttp.RespondWith(200, `{"type":"Bearer","value":"token"}`),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/api/v1/builds/123/events"),
					ghttp.RespondWith(200, `{"origin":{"id":"58f5f81a", "source":"stdout"}, "payload":"using version of resource found in cache"}`),
				),
			)
		})

		It("returns the events for a given build", func() {
			events, err := client.BuildEvents("123")
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(Equal("asda"))
		})
	})

	Context("when credentials are invalid", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/api/v1/teams/awesome-team/auth/token"),
					ghttp.RespondWith(401, "not authorized"),
				),
			)
		})

		It("returns an error", func() {
			builds, _, err := client.Builds(concourse.Page{Limit: 5})
			Expect(err).To(HaveOccurred())
			Expect(builds).To(BeEmpty())
		})
	})
})
