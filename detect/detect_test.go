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

package detect_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var root string

		it.Before(func() {
			root = test.ScratchDir(t, "detect")
		})

		it("contains default values", func() {
			defer test.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()
			defer internal.ProtectEnv(t, "TEST_KEY")

			console, e := internal.ReplaceConsole(t)
			defer e()

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

			test.WriteFile(t, filepath.Join(root, "platform", "env", "TEST_KEY"), "test-value")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(HaveOccurred())

			console.In(t, `[alpha]
  version = "alpha-version"
  name = "alpha-name"

[bravo]
  name = "bravo-name"
`)

			g.Expect(d.BuildPlan.Init()).To(Succeed())

			g.Expect(d.Application).NotTo(BeZero())
			g.Expect(d.Buildpack).NotTo(BeZero())
			g.Expect(d.BuildPlan).NotTo(BeZero())
			g.Expect(d.BuildPlanWriter).NotTo(BeZero())
			g.Expect(d.Logger).NotTo(BeZero())
			g.Expect(d.Platform).NotTo(BeZero())
			g.Expect(d.Stack).NotTo(BeZero())

			g.Expect(os.Getenv("TEST_KEY")).To(Equal("test-value"))
		})

		it("returns code when erroring", func() {
			defer test.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.TouchFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(d.Error(42)).To(Equal(42))
		})

		it("returns 100 when failing", func() {
			defer test.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.TouchFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(d.Fail()).To(Equal(detect.FailStatusCode))
		})

		it("returns 0 and BuildPlan when passing", func() {
			defer test.ReplaceWorkingDirectory(t, root)()
			defer test.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			test.TouchFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(d.Pass(buildplan.BuildPlan{
				"alpha": buildplan.Dependency{Version: "test-version"},
			})).To(Equal(detect.PassStatusCode))

			g.Expect(filepath.Join(root, "plan.toml")).To(test.HaveContent(`[alpha]
  version = "test-version"
`))
		})
	}, spec.Report(report.Terminal{}))
}
