package e2e

import (
	"os"

	"github.com/ocibuilder/ocibuilder/testing/e2e/resources/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("ocictl docker", func() {
	var session *gexec.Session
	var ocictlPath string

	BeforeEach(func() {
		ocictlPath = utils.BuildOcictl()
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
		if _, err := os.Stat("./spec.yaml"); err == nil {
			os.Remove("spec.yaml")
		}
	})

	It("exits with status code 0", func() {
		session = utils.RunOcictl(ocictlPath, nil)
		Eventually(session).Should(gexec.Exit(0))
	})

	It("completes a build and exits with status code 0", func() {
		args := []string{"build", "-p", "./resources/go-test-service"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 15).Should(gexec.Exit(0))
	}, 15)

	It("completes a push and exits with status code 0", func() {
		args := []string{"push", "-p", "./resources/go-test-service"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 10).Should(gexec.Exit(0))
	}, 10)

	It("completes a pull and exits with status code 0", func() {
		args := []string{"pull", "-i", "ocibuildere2e/go-test-service:v0.1.0", "-p", "./resources/go-test-service"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 10).Should(gexec.Exit(0))
	}, 10)

	It("completes an init and exits with status code 0", func() {
		args := []string{"init"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 2).Should(gexec.Exit(0))
	})

	It("completes a version and exits with status code 0", func() {
		args := []string{"version"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(session).Should(gexec.Exit(0))
	}, 1)

})
