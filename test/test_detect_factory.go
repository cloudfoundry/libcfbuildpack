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

package test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack/buildpack"
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/buildpack/libbuildpack/platform"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// DetectFactory is a factory for creating a test Detect.
type DetectFactory struct {
	// Detect is the configured detect to use.
	Detect detect.Detect

	// Home is the home directory to use.
	Home string

	// Output is the BuildPlan output at termination.
	Output buildplan.BuildPlan
}

// AddBuildPlan adds an entry to a build plan.
func (f *DetectFactory) AddBuildPlan(t *testing.T, name string, dependency buildplan.Dependency) {
	t.Helper()
	f.Detect.BuildPlan[name] = dependency
}

// AddEnv adds an environment variable to the Platform
func (f *DetectFactory) AddEnv(t *testing.T, name string, value string) {
	t.Helper()

	file := filepath.Join(f.Detect.Platform.Root, "env", name)
	if err := layers.WriteToFile(strings.NewReader(value), file, 0644); err != nil {
		t.Fatal(err)
	}

	f.Detect.Platform.Envs = append(f.Detect.Platform.Envs, platform.EnvironmentVariable{File: file, Name: name})
}

// NewDetectFactory creates a new instance of DetectFactory.
func NewDetectFactory(t *testing.T) *DetectFactory {
	t.Helper()
	f := DetectFactory{}

	root := internal.ScratchDir(t, "test-detect-factory")

	f.Detect.Application.Root = filepath.Join(root, "application")

	f.Detect.Buildpack.Metadata = make(buildpack.Metadata)
	f.Detect.Buildpack.Metadata["dependencies"] = make([]map[string]interface{}, 0)

	f.Detect.BuildPlan = make(buildplan.BuildPlan)

	f.Detect.BuildPlanWriter = func(buildPlan buildplan.BuildPlan) error {
		f.Output = buildPlan
		return nil
	}

	f.Detect.Platform.Root = filepath.Join(root, "platform")
	f.Detect.Platform.Envs = make(platform.EnvironmentVariables, 0)

	f.Detect.Stack = "test-stack"

	f.Home = filepath.Join(root, "home")

	return &f
}
