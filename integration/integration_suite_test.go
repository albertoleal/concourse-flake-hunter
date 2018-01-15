package integration_test

import (
	"github.com/masters-of-cats/concourse-flake-hunter/integration/runner"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var (
	Runner        runner.Runner
	hunterBinPath string
	concourseURL  string
	username      string
	password      string
	team          string
)

func TestConcourseFlakeHunter(t *testing.T) {
	RegisterFailHandler(Fail)

	SynchronizedBeforeSuite(func() []byte {
		var err error

		binPath, err := gexec.Build("github.com/masters-of-cats/concourse-flake-hunter")
		Expect(err).NotTo(HaveOccurred())

		return []byte(binPath)
	}, func(data []byte) {
		hunterBinPath = string(data)
	})

	BeforeEach(func() {
		concourseURL = "https://my.concourse.io"
		username = "alice"
		password = "my-passwd"
		team = "awesome-team"
	})

	JustBeforeEach(func() {
		Runner = runner.Runner{
			ConcourseFlakeHunter: hunterBinPath,
			ConcourseURL:         concourseURL,
			Username:             username,
			Password:             password,
			Team:                 team,
		}
	})

	SynchronizedAfterSuite(func() {
	}, func() {
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "concourse-flake-hunter integration suite")
}
