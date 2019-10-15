package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os/exec"
)

var _ = Describe("ocictl docker", func() {
	var session *gexec.Session

	BeforeEach(func() {
		ocictlPath := buildOcictl()
		session = runOcictl(ocictlPath)
	})

})


func buildOcictl() string {
	ocictlPath, err := gexec.Build("github.com/ocibuilder/ocictl")
	Expect(err).NotTo(HaveOccurred())

	return ocictlPath
}

func runOcictl(path string) *gexec.Session {
	cmd := exec.Command(path)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}
