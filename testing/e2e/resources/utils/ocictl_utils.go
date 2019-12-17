package utils

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

// BuildOcictl builds the ocictl for e2e testing purposes
func BuildOcictl() string {
	ocictlPath, err := gexec.Build("github.com/ocibuilder/ocibuilder/ocictl")
	Expect(err).NotTo(HaveOccurred())
	return ocictlPath
}

// RunOcictl runs the ocictl for e2e testing purposes
func RunOcictl(path string, args []string) *gexec.Session {
	cmd := exec.Command(path, args...)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}

func RunDockerInspect(imageName string) *gexec.Session {
	cmd := exec.Command("docker", "image", "inspect", imageName)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())

	return session
}
