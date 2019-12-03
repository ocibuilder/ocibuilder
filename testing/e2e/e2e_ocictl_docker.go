/*
Copyright 2019 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"os"

	"github.com/ocibuilder/ocibuilder/testing/e2e/resources/utils"
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
		if _, err := os.Stat("./ocibuilder.yaml"); err == nil {
			os.Remove("ocibuilder.yaml")
		}
	})

	It("exits with status code 0", func() {
		session = utils.RunOcictl(ocictlPath, nil)
		Eventually(session).Should(gexec.Exit(0))
	})

	It("completes a build and exits with status code 0", func() {
		args := []string{"build", "-p", "./resources/go-test-service"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 30).Should(gexec.Exit(0))
	}, 30)

	It("completes a push and exits with status code 0", func() {
		args := []string{"push", "-p", "./resources/go-test-service"}
		session = utils.RunOcictl(ocictlPath, args)
		Eventually(func() *gexec.Session {
			return session
		}, 20).Should(gexec.Exit(0))
	}, 20)

	It("completes a pull and exits with status code 0", func() {
		args := []string{"pull", "-i", "ocibuildere2e/go-test-service:v0.1.0", "-p", "./resources/go-test-service"}
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
