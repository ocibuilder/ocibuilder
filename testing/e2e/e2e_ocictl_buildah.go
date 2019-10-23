package e2e

import (
	"github.com/ocibuilder/ocibuilder/testing/e2e/resources/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("ocictl buildah", func() {
	var session *gexec.Session
	var ocictlPath string

	BeforeEach(func() {
		ocictlPath = utils.BuildOcictl()
	})

	It("exits with status code 0", func() {
		session = utils.RunOcictl(ocictlPath, nil)
		Eventually(session).Should(gexec.Exit(0))
	})

	It("completes a build and exits with status code 0", func() {
		args := []string{"build", "-b", "buildah", "-s", "vfs", "-p", "./resources/go-test-service"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 25).Should(gexec.Exit(0))
	}, 25)

	It("completes a push and exits with status code 0", func() {
		args := []string{"push", "-b", "buildah", "-p", "./resources/go-test-service"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 25).Should(gexec.Exit(0))
	}, 25)

	It("completes a pull and exits with status code 0", func() {
		args := []string{"pull", "-b", "buildah", "-i", "ocibuildere2e/go-test-service:v0.1.0", "-p", "./resources/go-test-service"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 25).Should(gexec.Exit(0))
	}, 25)

})
