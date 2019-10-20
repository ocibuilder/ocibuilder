package e2e

import (
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ocictl suite")
}

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
