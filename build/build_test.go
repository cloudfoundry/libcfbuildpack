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

package build_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuild(t *testing.T) {
	spec.Run(t, "Build", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var root string

		it.Before(func() {
			root = test.ScratchDir(t, "detect")
		})

		it("contains default values", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "layers"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()
			defer internal.ProtectEnv(t, "TEST_KEY")

			console, d := internal.ReplaceConsole(t)
			defer d()

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

			console.In(t, `[alpha]
  version = "alpha-version"
  name = "alpha-name"

[bravo]
  name = "bravo-name"
`)

			test.WriteFile(t, filepath.Join(root, "platform", "env", "TEST_KEY"), "test-value")

			b, err := build.DefaultBuild()
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(b.Application).NotTo(BeZero())
			g.Expect(b.Buildpack).NotTo(BeZero())
			g.Expect(b.BuildPlan).NotTo(BeZero())
			g.Expect(b.BuildPlanWriter).NotTo(BeZero())
			g.Expect(b.Layers).NotTo(BeZero())
			g.Expect(b.Logger).NotTo(BeZero())
			g.Expect(b.Platform).NotTo(BeZero())
			g.Expect(b.Stack).NotTo(BeZero())

			g.Expect(os.Getenv("TEST_KEY")).To(Equal("test-value"))
		})

		it("returns 0 when successful", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "layers"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			console, d := internal.ReplaceConsole(t)
			defer d()

			test.TouchFile(t, root, "buildpack.toml")
			console.In(t, "")

			b, err := build.DefaultBuild()
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(b.Success(buildplan.BuildPlan{
				"alpha": buildplan.Dependency{Version: "test-version"},
			})).To(Equal(build.SuccessStatusCode))

			g.Expect(filepath.Join(root, "plan.toml")).To(test.HaveContent(`[alpha]
  version = "test-version"
`))
		})

		it("returns code when failing", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "layers"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			console, d := internal.ReplaceConsole(t)
			defer d()

			test.TouchFile(t, root, "buildpack.toml")
			console.In(t, "")

			b, err := build.DefaultBuild()
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(b.Failure(42)).To(Equal(42))
		})
	}, spec.Report(report.Terminal{}))
}
