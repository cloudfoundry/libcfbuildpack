/*
 * Copyright 2018-2019 the original author or authors.
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

package helper_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/helper"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestReadBuildpackYaml(t *testing.T) {
	type BuildpackYaml struct {
		Python struct {
			Version string `yaml:"version"`
		} `yaml:"python"`

		Go struct {
			Version string `yaml:"version"`
		} `yaml:"go"`
	}

	spec.Run(t, "ReadBuildpackYaml", func(t *testing.T, when spec.G, it spec.S) {
		var (
			pythonVersion     = "1.2.3"
			goVersion         = "4.5.6"
			buildpackYamlPath string
			config            *BuildpackYaml
		)

		Expect := NewWithT(t).Expect

		when("read buildpack yaml version", func() {
			it.Before(func() {
				tmpDir := os.TempDir()
				buildpackYamlPath = filepath.Join(tmpDir, "buildpack.yml")
				buildpackYAMLString := fmt.Sprintf("python:\n  version: %s\ngo:\n version: %s", pythonVersion, goVersion)
				Expect(helper.WriteFile(buildpackYamlPath, 0777, buildpackYAMLString)).To(Succeed())
				config = &BuildpackYaml{}
			})

			it.After(func() {
				Expect(os.RemoveAll(buildpackYamlPath)).To(Succeed())
			})

			it("unmarshals a user defined config when given a buildpackyml path", func() {
				err := helper.ReadBuildpackYaml(buildpackYamlPath, config)
				Expect(err).NotTo(HaveOccurred())
				Expect(config.Python.Version).To(Equal(pythonVersion))
				Expect(config.Go.Version).To(Equal(goVersion))
			})

		})
	}, spec.Report(report.Terminal{}))
}
