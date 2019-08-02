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

package detect_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var root string

		it.Before(func() {
			root = test.ScratchDir(t, "detect")
		})

		it("contains default values", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.WriteFile(t, filepath.Join(root, "buildpack.toml"), `[buildpack]
id = "buildpack-id"
name = "buildpack-name"
version = "buildpack-version"

[[stacks]]
id = 'stack-id'
build-images = ["build-image-tag"]
run-images = ["run-image-tag"]

[metadata]
test-key = "test-value"
`)

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(d.Application).NotTo(gomega.BeZero())
			g.Expect(d.Buildpack).NotTo(gomega.BeZero())
			g.Expect(d.Logger).NotTo(gomega.BeZero())
			g.Expect(d.Platform).NotTo(gomega.BeZero())
			g.Expect(d.Services).NotTo(gomega.BeZero())
			g.Expect(d.Stack).NotTo(gomega.BeZero())
			g.Expect(d.Writer).NotTo(gomega.BeZero())
		})

		it("returns code when erroring", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.TouchFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(d.Error(42)).To(gomega.Equal(42))
		})

		it("returns 100 when failing", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.TouchFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(d.Fail()).To(gomega.Equal(detect.FailStatusCode))
		})

		it("returns 0 and Plan when passing", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.TouchFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(d.Pass(buildplan.Plan{
				Provides: []buildplan.Provided{
					{"test-provided-1a"},
					{"test-provided-1b"},
				},
				Requires: []buildplan.Required{
					{"test-required-1a", "test-version-1a", buildplan.Metadata{"test-key-1a": "test-value-1a"}},
					{"test-required-1b", "test-version-1b", buildplan.Metadata{"test-key-1b": "test-value-1b"}},
				},
			},
				buildplan.Plan{
					Provides: []buildplan.Provided{
						{"test-provided-2a"},
						{"test-provided-2b"},
					},
					Requires: []buildplan.Required{
						{"test-required-2a", "test-version-2a", buildplan.Metadata{"test-key-2a": "test-value-2a"}},
						{"test-required-2b", "test-version-2b", buildplan.Metadata{"test-key-2b": "test-value-2b"}},
					},
				},
				buildplan.Plan{
					Provides: []buildplan.Provided{
						{"test-provided-3a"},
						{"test-provided-3b"},
					},
					Requires: []buildplan.Required{
						{"test-required-3a", "test-version-3a", buildplan.Metadata{"test-key-3a": "test-value-3a"}},
						{"test-required-3b", "test-version-3b", buildplan.Metadata{"test-key-3b": "test-value-3b"}},
					},
				})).To(gomega.Equal(detect.PassStatusCode))

			g.Expect(filepath.Join(root, "plan.toml")).To(test.HaveContent(`[[provides]]
  name = "test-provided-1a"

[[provides]]
  name = "test-provided-1b"

[[requires]]
  name = "test-required-1a"
  version = "test-version-1a"
  [requires.metadata]
    test-key-1a = "test-value-1a"

[[requires]]
  name = "test-required-1b"
  version = "test-version-1b"
  [requires.metadata]
    test-key-1b = "test-value-1b"

[[or]]

  [[or.provides]]
    name = "test-provided-2a"

  [[or.provides]]
    name = "test-provided-2b"

  [[or.requires]]
    name = "test-required-2a"
    version = "test-version-2a"
    [or.requires.metadata]
      test-key-2a = "test-value-2a"

  [[or.requires]]
    name = "test-required-2b"
    version = "test-version-2b"
    [or.requires.metadata]
      test-key-2b = "test-value-2b"

[[or]]

  [[or.provides]]
    name = "test-provided-3a"

  [[or.provides]]
    name = "test-provided-3b"

  [[or.requires]]
    name = "test-required-3a"
    version = "test-version-3a"
    [or.requires.metadata]
      test-key-3a = "test-value-3a"

  [[or.requires]]
    name = "test-required-3b"
    version = "test-version-3b"
    [or.requires.metadata]
      test-key-3b = "test-value-3b"
`))
		})
	}, spec.Report(report.Terminal{}))
}

