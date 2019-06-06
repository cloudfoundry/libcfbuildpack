/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helper_test

import (
	"fmt"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"os"
	"path/filepath"
	"testing"
)

func TestReadBuildpackYaml(t *testing.T) {
	spec.Run(t, "ReadBuildpackYaml", func(t *testing.T, when spec.G, it spec.S) {
		var(
		pythonVersion = "1.2.3"
		goVersion = "4.5.6"
		buildpackYamlPath 	string
		)

		Expect := NewWithT(t).Expect


		when("read buildpack yaml version", func() {
			 it.Before(func() {
				 tmpDir := os.TempDir()
				 buildpackYamlPath = filepath.Join(tmpDir, "buildpack.yml")
				 buildpackYAMLString := fmt.Sprintf("python:\n  version: %s\ngo:\n version: %s", pythonVersion, goVersion)
				 Expect(helper.WriteFile(buildpackYamlPath, 0777, buildpackYAMLString)).To(Succeed())
			 })

			 it.After(func() {
				 Expect(os.RemoveAll(buildpackYamlPath)).To(Succeed())
			 })

			 it("returns the version specified for the given language", func() {
				 buildpackYamlVersion, err := helper.ReadBuildpackYamlVersion(buildpackYamlPath, "python")
				 Expect(err).NotTo(HaveOccurred())
				 Expect(buildpackYamlVersion).To(Equal(pythonVersion))

				 buildpackYamlVersion, err = helper.ReadBuildpackYamlVersion(buildpackYamlPath, "go")
				 Expect(err).NotTo(HaveOccurred())
				 Expect(buildpackYamlVersion).To(Equal(goVersion))
			 })

			 it("returns empty string if no version specified for given language", func() {
				 buildpackYamlVersion, err := helper.ReadBuildpackYamlVersion(buildpackYamlPath, "some-language")
				 Expect(err).NotTo(HaveOccurred())
				 Expect(buildpackYamlVersion).To(Equal(""))
			 })
		})
	}, spec.Report(report.Terminal{}))
}
