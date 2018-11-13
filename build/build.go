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

package build

import (
	"fmt"

	"github.com/buildpack/libbuildpack/build"
	layersBp "github.com/buildpack/libbuildpack/layers"
	buildpackPkg "github.com/cloudfoundry/libcfbuildpack/buildpack"
	layersCf "github.com/cloudfoundry/libcfbuildpack/layers"
	loggerPkg "github.com/cloudfoundry/libcfbuildpack/logger"
)

// Build is an extension to libbuildpack.Build that allows additional functionality to be added.
type Build struct {
	build.Build

	// Buildpack represents the metadata associated with a buildpack.
	Buildpack buildpackPkg.Buildpack

	// Layers represents the launch layers contributed by a buildpack.
	Layers layersCf.Layers

	// Logger is used to write debug and info to the console.
	Logger loggerPkg.Logger
}

// String makes Build satisfy the Stringer interface.
func (b Build) String() string {
	return fmt.Sprintf("Build{ Build: %s, Buildpack: %s, Layers: %s, Logger: %s }",
		b.Build, b.Buildpack, b.Layers, b.Logger)
}

// DefaultBuild creates a new instance of Build using default values.
func DefaultBuild() (Build, error) {
	b, err := build.DefaultBuild()
	if err != nil {
		return Build{}, err
	}

	logger := loggerPkg.Logger{Logger: b.Logger}
	buildpack := buildpackPkg.NewBuildpack(b.Buildpack)
	layers := layersCf.Layers{
		Layers: b.Layers,
		BuildpackCache: layersBp.Layers{
			Root:   buildpack.CacheRoot,
			Logger: b.Logger,
		},
		Logger: logger,
	}

	return Build{
		b,
		buildpack,
		layers,
		logger,
	}, nil
}
