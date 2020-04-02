package e2e

import (
	"github.com/beval/beval/testing/e2e/resources/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("bevalctl buildah", func() {
	var session *gexec.Session
	var bevalctlPath string

	BeforeEach(func() {
		bevalctlPath = utils.Buildbevalctl()
	})

	It("exits with status code 0", func() {
		session = utils.Runbevalctl(bevalctlPath, nil)
		Eventually(session).Should(gexec.Exit(0))
	})

	It("completes a build and exits with status code 0", func() {
		args := []string{"build", "-b", "buildah", "-p", "./resources/go-test-service", "-d"}
		session = utils.Runbevalctl(bevalctlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 25).Should(gexec.Exit(0))
	}, 25)

	It("completes a push and exits with status code 0", func() {
		args := []string{"push", "-b", "buildah", "-p", "./resources/go-test-service", "-d"}
		session = utils.Runbevalctl(bevalctlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 25).Should(gexec.Exit(0))
	}, 25)

	It("completes a pull and exits with status code 0", func() {
		args := []string{"pull", "-b", "buildah", "-i", "bevale2e/go-test-service:v0.1.0", "-p", "./resources/go-test-service", "-d"}
		session = utils.Runbevalctl(bevalctlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 25).Should(gexec.Exit(0))
	}, 25)

})
