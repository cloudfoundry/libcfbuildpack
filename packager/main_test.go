/*
 * Copyright 2019-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestIntegrationPackager(t *testing.T) {
	spec.Run(t, "Package Integration", testPackagerIntegration, spec.Report(report.Terminal{}))
}

func testPackagerIntegration(t *testing.T, when spec.G, it spec.S) {
	var (
		Expect                                                              func(interface{}, ...interface{}) Assertion
		cnbDir, outputDir, tempDir, buildpackTomlPath, oldBPToml, newBPToml string
		err                                                                 error
		newVersion                                                          string
	)

	it.Before(func() {
		Expect = NewWithT(t).Expect

		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		cnbDir = filepath.Join(tempDir, "cnb")
		Expect(os.MkdirAll(cnbDir, 0777))
		newVersion = "newVersion"
		buildpackTomlPath = filepath.Join(cnbDir, "buildpack.toml")
		baseBuildpackTOML := `[buildpack]
id = "buildpack-id"
name = "buildpack-name"
version = "%s"

[metadata]
include_files = ["buildpack.toml"]`
		oldBPToml = fmt.Sprintf(baseBuildpackTOML, "{{ .Version }}")
		newBPToml = fmt.Sprintf(baseBuildpackTOML, newVersion)

		outputDir = filepath.Join(tempDir, "output")
	})

	it.After(func() {
		os.RemoveAll(tempDir)
	})

	when("packaging", func() {
		it("does it in a temp directory", func() {
			Expect(ioutil.WriteFile(buildpackTomlPath, []byte(oldBPToml), 0666)).To(Succeed())

			// Set first arg equal to command name, per: https://stackoverflow.com/questions/33723300/how-to-test-the-passing-of-arguments-in-golang
			os.Args = []string{"cmd", "-version=" + newVersion, outputDir}
			Expect(os.Chdir(cnbDir)).To(Succeed())
			main()
			originalOutput, err := ioutil.ReadFile(buildpackTomlPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(originalOutput)).To(Equal(oldBPToml))

			newOutput, err := ioutil.ReadFile(filepath.Join(outputDir, "buildpack.toml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(string(newOutput)).To(Equal(newBPToml))
		})
	})
}
