/*
 * Copyright 2018 the original author or authors.
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
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack/application"
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/buildpack/libbuildpack/platform"
	buildPkg "github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuild(t *testing.T) {
	spec.Run(t, "Build", testBuild, spec.Report(report.Terminal{}))
}

func testBuild(t *testing.T, when spec.G, it spec.S) {

	it("contains default values", func() {
		root := internal.ScratchDir(t, "detect")
		defer internal.ReplaceWorkingDirectory(t, root)()
		defer internal.ReplaceEnv(t, "PACK_STACK_ID", "test-stack")()

		console, d := internal.ReplaceConsole(t)
		defer d()

		console.In(t, `[alpha]
  version = "alpha-version"
  name = "alpha-name"

[bravo]
  name = "bravo-name"
`)

		in := strings.NewReader(`[buildpack]
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

		if err := layers.WriteToFile(in, filepath.Join(root, "buildpack.toml"), 0644); err != nil {
			t.Fatal(err)
		}

		defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), root, root, root, root)()

		build, err := buildPkg.DefaultBuild()
		if err != nil {
			t.Fatal(err)
		}

		if reflect.DeepEqual(build.Application, application.Application{}) {
			t.Errorf("detect.Application should not be empty")
		}

		if reflect.DeepEqual(build.Buildpack, buildpack.Buildpack{}) {
			t.Errorf("detect.Buildpack should not be empty")
		}

		if reflect.DeepEqual(build.BuildPlan, buildplan.BuildPlan{}) {
			t.Errorf("detect.BuildPlan should not be empty")
		}

		if reflect.DeepEqual(build.Layers, layers.Layers{}) {
			t.Errorf("detect.Layers should not be empty")
		}

		if reflect.DeepEqual(build.Logger, logger.Logger{}) {
			t.Errorf("detect.Logger should not be empty")
		}

		if reflect.DeepEqual(build.Platform, platform.Platform{}) {
			t.Errorf("detect.Platform should not be empty")
		}

		if reflect.DeepEqual(build.Stack, "") {
			t.Errorf("detect.Stack should not be empty")
		}
	})

	it("returns 0 when successful", func() {
		root := internal.ScratchDir(t, "build")
		defer internal.ReplaceWorkingDirectory(t, root)()
		defer internal.ReplaceEnv(t, "PACK_STACK_ID", "test-stack")()

		c, d := internal.ReplaceConsole(t)
		defer d()
		c.In(t, "")

		if err := layers.WriteToFile(strings.NewReader(""), filepath.Join(root, "buildpack.toml"), 0644); err != nil {
			t.Fatal(err)
		}

		defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), root, root, root, root)()

		build, err := buildPkg.DefaultBuild()
		if err != nil {
			t.Fatal(err)
		}

		actual, err := build.Success(buildplan.BuildPlan{
			"alpha": buildplan.Dependency{Version: "test-version"},
		})
		if err != nil {
			t.Fatal(err)
		}

		if actual != 0 {
			t.Errorf("Build.Success() = %d, expected 0", actual)
		}

		test.BeFileLike(t, filepath.Join(root, "alpha"), 0644, `version = "test-version"
`)
	})

	it("returns code when failing", func() {
		root := internal.ScratchDir(t, "build")
		defer internal.ReplaceWorkingDirectory(t, root)()
		defer internal.ReplaceEnv(t, "PACK_STACK_ID", "test-stack")()

		c, d := internal.ReplaceConsole(t)
		defer d()
		c.In(t, "")

		if err := layers.WriteToFile(strings.NewReader(""), filepath.Join(root, "buildpack.toml"), 0644); err != nil {
			t.Fatal(err)
		}

		defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), root, root, root, root)()

		build, err := buildPkg.DefaultBuild()
		if err != nil {
			t.Fatal(err)
		}

		actual := build.Failure(42)

		if actual != 42 {
			t.Errorf("Build.Failure() = %d, expected 42", actual)
		}
	})
}
