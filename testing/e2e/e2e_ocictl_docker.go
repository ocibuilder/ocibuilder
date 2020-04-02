package e2e

import (
	"os"

	"github.com/beval/beval/testing/e2e/resources/utils"
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
		if _, err := os.Stat("./beval.yaml"); err == nil {
			os.Remove("beval.yaml")
		}
	})

	It("exits with status code 0", func() {
		session = utils.RunOcictl(ocictlPath, nil)
		Eventually(session).Should(gexec.Exit(0))
	})

	It("completes a build and exits with status code 0", func() {
		args := []string{"build", "-p", utils.BuildPath}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 30).Should(gexec.Exit(0))
	}, 30)

	It("completes a build with an overlay and exits with status code 0", func() {
		args := []string{"build", "-p", utils.BuildPath, "--overlay", utils.OverlayPath}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 30).Should(gexec.Exit(0))

		inspectSession := utils.RunDockerInspect(utils.ImageNameOverlayed)
		Eventually(inspectSession).Should(gexec.Exit(0))

	}, 30)

	It("completes a push and exits with status code 0", func() {
		args := []string{"push", "-p", utils.BuildPath}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 20).Should(gexec.Exit(0))
	}, 20)

	It("completes a pull and exits with status code 0", func() {
		args := []string{"pull", "-i", utils.ImageName, "-p", utils.BuildPath}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 15).Should(gexec.Exit(0))
	}, 15)

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
