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
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	buildpackBp "github.com/buildpack/libbuildpack/buildpack"
	"github.com/buildpack/libbuildpack/buildplan"
	layersBp "github.com/buildpack/libbuildpack/layers"
	"github.com/buildpack/libbuildpack/platform"
	"github.com/cloudfoundry/libcfbuildpack/build"
	buildpackCf "github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	layersCf "github.com/cloudfoundry/libcfbuildpack/layers"
)

// BuildFactory is a factory for creating a test Build.
type BuildFactory struct {
	// Build is the configured build to use.
	Build build.Build

	// Home is the home directory to use.
	Home string

	// Output is the BuildPlan output at termination.
	Output buildplan.BuildPlan
}

// AddBuildPlan adds an entry to a build plan.
func (f *BuildFactory) AddBuildPlan(t *testing.T, name string, dependency buildplan.Dependency) {
	t.Helper()
	f.Build.BuildPlan[name] = dependency
}

// AddDependency adds a dependency to the buildpack metadata and copies a fixture into a cached dependency layer.
func (f *BuildFactory) AddDependency(t *testing.T, id string, fixture string) {
	t.Helper()

	d := f.newDependency(t, id, fixture)
	f.cacheFixture(t, d, fixture)
	f.addDependency(t, d)
}

// AddEnv adds an environment variable to the Platform
func (f *BuildFactory) AddEnv(t *testing.T, name string, value string) {
	t.Helper()

	file := filepath.Join(f.Build.Platform.Root, "env", name)
	if err := layersCf.WriteToFile(strings.NewReader(value), file, 0644); err != nil {
		t.Fatal(err)
	}

	f.Build.Platform.Envs = append(f.Build.Platform.Envs, platform.EnvironmentVariable{File: file, Name: name})
}

func (f *BuildFactory) addDependency(t *testing.T, dependency buildpackCf.Dependency) {
	t.Helper()

	metadata := f.Build.Buildpack.Metadata
	dependencies := metadata["dependencies"].([]map[string]interface{})

	var stacks []interface{}
	for _, stack := range dependency.Stacks {
		stacks = append(stacks, stack)
	}

	var licenses []map[string]interface{}
	for _, license := range dependency.Licenses {
		licenses = append(licenses, map[string]interface{}{
			"type": license.Type,
			"uri":  license.URI,
		})
	}

	metadata["dependencies"] = append(dependencies, map[string]interface{}{
		"id":       dependency.ID,
		"name":     dependency.Name,
		"version":  dependency.Version.Version.Original(),
		"uri":      dependency.URI,
		"sha256":   dependency.SHA256,
		"stacks":   stacks,
		"licenses": licenses,
	})
}

func (f *BuildFactory) cacheFixture(t *testing.T, dependency buildpackCf.Dependency, fixture string) {
	t.Helper()

	l := f.Build.Layers.Layer(dependency.SHA256)
	if err := layersCf.CopyFile(FixturePath(t, fixture), filepath.Join(l.Root, filepath.Base(fixture))); err != nil {
		t.Fatal(err)
	}

	d, err := internal.ToTomlString(map[string]interface{}{"metadata": dependency})
	if err != nil {
		t.Fatal(err)
	}
	if err := layersCf.WriteToFile(strings.NewReader(d), l.Metadata, 0644); err != nil {
		t.Fatal(err)
	}
}

func (f *BuildFactory) newDependency(t *testing.T, id string, fixture string) buildpackCf.Dependency {
	t.Helper()

	version, err := semver.NewVersion("1.0")
	if err != nil {
		t.Fatal(err)
	}

	return buildpackCf.Dependency{
		ID:      id,
		Name:    "test-name",
		Version: buildpackCf.Version{Version: version},
		URI:     fmt.Sprintf("http://localhost/%s", filepath.Base(fixture)),
		SHA256:  "test-hash",
		Stacks:  buildpackCf.Stacks{f.Build.Stack},
		Licenses: buildpackCf.Licenses{
			buildpackCf.License{Type: "test-type"},
		},
	}
}

// NewBuildFactory creates a new instance of BuildFactory.
func NewBuildFactory(t *testing.T) *BuildFactory {
	t.Helper()
	f := BuildFactory{}

	root := internal.ScratchDir(t, "test-build-factory")

	f.Build.Application.Root = filepath.Join(root, "application")

	f.Build.Buildpack.Metadata = make(buildpackBp.Metadata)
	f.Build.Buildpack.Metadata["dependencies"] = make([]map[string]interface{}, 0)

	f.Build.BuildPlan = make(buildplan.BuildPlan)

	f.Build.BuildPlanWriter = func(buildPlan buildplan.BuildPlan) error {
		f.Output = buildPlan
		return nil
	}

	f.Build.Layers.Root = filepath.Join(root, "layers")
	f.Build.Layers.BuildpackCache = layersBp.Layers{Root: filepath.Join(root, "buildpack-cache")}

	f.Build.Platform.Root = filepath.Join(root, "platform")
	f.Build.Platform.Envs = make(platform.EnvironmentVariables, 0)

	f.Build.Stack = "test-stack"

	f.Home = filepath.Join(root, "home")

	return &f
}
