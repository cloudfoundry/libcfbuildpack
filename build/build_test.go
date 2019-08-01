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

package build_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuild(t *testing.T) {
	spec.Run(t, "Build", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var root string

		it.Before(func() {
			root = test.ScratchDir(t, "detect")
		})

		it("contains default values", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "layers"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

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

			test.WriteFile(t, filepath.Join(root, "plan.toml"), `[[entries]]
  name = "test-entry-1a"
  version = "test-version-1a"
  [entries.metadata]
    test-key-1a = "test-value-1a"

[[entries]]
  name = "test-entry-1b"
  version = "test-version-1b"
  [entries.metadata]
    test-key-1b = "test-value-1b"
`)

			b, err := build.DefaultBuild()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(b.Application).NotTo(gomega.BeZero())
			g.Expect(b.Buildpack).NotTo(gomega.BeZero())
			g.Expect(b.Plans).NotTo(gomega.BeZero())
			g.Expect(b.Layers).NotTo(gomega.BeZero())
			g.Expect(b.Logger).NotTo(gomega.BeZero())
			g.Expect(b.Platform).NotTo(gomega.BeZero())
			g.Expect(b.Services).NotTo(gomega.BeZero())
			g.Expect(b.Stack).NotTo(gomega.BeZero())
			g.Expect(b.Writer).NotTo(gomega.BeZero())
		})

		it("returns 0 when successful", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "layers"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.TouchFile(t, root, "buildpack.toml")
			test.TouchFile(t, root, "plan.toml")

			b, err := build.DefaultBuild()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(b.Success(
				buildpackplan.Plan{
					Name:     "test-entry-1a",
					Version:  "test-version-1a",
					Metadata: buildpackplan.Metadata{"test-key-1a": "test-value-1a"},
				},
				buildpackplan.Plan{
					Name:     "test-entry-1b",
					Version:  "test-version-1b",
					Metadata: buildpackplan.Metadata{"test-key-1b": "test-value-1b"},
				})).To(gomega.Equal(build.SuccessStatusCode))

			g.Expect(filepath.Join(root, "plan.toml")).To(test.HaveContent(`[[entries]]
  name = "test-entry-1a"
  version = "test-version-1a"
  [entries.metadata]
    test-key-1a = "test-value-1a"

[[entries]]
  name = "test-entry-1b"
  version = "test-version-1b"
  [entries.metadata]
    test-key-1b = "test-value-1b"
`))
		})

		it("returns code when failing", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "layers"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.TouchFile(t, root, "buildpack.toml")
			test.TouchFile(t, root, "plan.toml")

			b, err := build.DefaultBuild()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(b.Failure(42)).To(gomega.Equal(42))
		})
	}, spec.Report(report.Terminal{}))
}
