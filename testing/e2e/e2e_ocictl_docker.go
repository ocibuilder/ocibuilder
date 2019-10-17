package e2e

import (
	"os"
	"os/exec"

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

	AfterEach(func() {
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
		}, 5).Should(gexec.Exit(0))
	}, 5)

	It("completes an init and exits with status code 0", func() {
		args := []string{"init"}
		session = runOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 2).Should(gexec.Exit(0))
	})

	It("completes a version and exits with statue code 0", func() {
		args := []string{"version"}
		session = runOcictl(ocictlPath, args)
		Eventually(session).Should(gexec.Exit(0))
	}, 1)

})

func buildOcictl() string {
	ocictlPath, err := gexec.Build("github.com/ocibuilder/ocibuilder/ocictl")
	Expect(err).NotTo(HaveOccurred())

	return ocictlPath
}

func runOcictl(path string, args []string) *gexec.Session {
	cmd := exec.Command(path, args...)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}
