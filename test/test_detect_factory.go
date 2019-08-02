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

package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	bp "github.com/buildpack/libbuildpack/services"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/services"
)

// DetectFactory is a factory for creating a test Detect.
type DetectFactory struct {
	// Detect is the configured detect to use.
	Detect detect.Detect

	// Home is the home directory to use.
	Home string

	// Plans is the build plans at termination.
	Plans buildplan.Plans

	// Runner is the used to capture commands executed outside the process.
	Runner *Runner

	t *testing.T
}

// AddService adds an entry to the collection of services.
func (f *DetectFactory) AddService(name string, credentials services.Credentials, tags ...string) {
	f.t.Helper()

	f.Detect.Services.Services = append(f.Detect.Services.Services, bp.Service{
		BindingName: name,
		Credentials: credentials,
		Tags:        tags,
	})
}

// NewDetectFactory creates a new instance of DetectFactory.
func NewDetectFactory(t *testing.T) *DetectFactory {
	t.Helper()

	root := ScratchDir(t, "detect")
	runner := &Runner{}

	f := DetectFactory{Home: filepath.Join(root, "home"), Runner: runner, t: t}

	f.Detect.Application.Root = filepath.Join(root, "application")
	if err := os.MkdirAll(f.Detect.Application.Root, 0755); err != nil {
		t.Fatal(err)
	}
	f.Detect.Buildpack.Info.Version = "1.0"
	f.Detect.Buildpack.Root = filepath.Join(root, "buildpack")
	f.Detect.Platform.Root = filepath.Join(root, "platform")
	f.Detect.Runner = runner
	f.Detect.Services = services.Services{Services: bp.Services{}}
	f.Detect.Stack = "test-stack"
	f.Detect.Writer = func(plans buildplan.Plans) error {
		f.Plans = plans
		return nil
	}

	f.Home = filepath.Join(root, "home")

	return &f
}
