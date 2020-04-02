package e2e

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestE2ebevalctl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "bevalctl suite")
}
