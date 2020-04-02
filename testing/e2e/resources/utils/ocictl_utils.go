package utils

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

// Buildbevalctl builds the bevalctl for e2e testing purposes
func Buildbevalctl() string {
	bevalctlPath, err := gexec.Build("github.com/beval/beval/bevalctl")
	Expect(err).NotTo(HaveOccurred())
	return bevalctlPath
}

// Runbevalctl runs the bevalctl for e2e testing purposes
func Runbevalctl(path string, args []string) *gexec.Session {
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
