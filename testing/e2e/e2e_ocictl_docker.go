package e2e

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("ocictl docker", func() {
	var session *gexec.Session
	var ocictlPath string

	BeforeEach(func() {
		ocictlPath = buildOcictl()
	})

	AfterSuite(func() {
		gexec.CleanupBuildArtifacts()
		if _, err := os.Stat("./spec.yaml"); err == nil {
			os.Remove("spec.yaml")
		}
	})

	It("exits with status code 0", func() {
		session = runOcictl(ocictlPath, nil)
		Eventually(session).Should(gexec.Exit(0))
	})

	It("completes a build and exits with status code 0", func() {
		args := []string{"build", "-p", "./resources/go-test-service"}
		session = runOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 10).Should(gexec.Exit(0))
	}, 10)

	It("completes a push and exits with status code 0", func() {
		args := []string{"push", "-p", "./resources/go-test-service"}
		session = runOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 10).Should(gexec.Exit(0))
	}, 10)

	It("completes a pull and exits with status code 0", func() {
		args := []string{"pull", "-i", "ocibuildere2e/go-test-service:v0.1.0", "-p", "./resources/go-test-service"}
		session = runOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 10).Should(gexec.Exit(0))
	}, 10)

	It("completes an init and exits with status code 0", func() {
		args := []string{"init"}
		session = runOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 2).Should(gexec.Exit(0))
	})

	It("completes a version and exits with status code 0", func() {
		args := []string{"version"}
		session = runOcictl(ocictlPath, args)
		Eventually(session).Should(gexec.Exit(0))
	}, 1)

})


