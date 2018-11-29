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

package layers

import (
	"fmt"

	"github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/fatih/color"
)

// Layers is an extension allows additional functionality to be added.
type Layers struct {
	layers.Layers

	// BuildpackCacheRoot is the root of the cache of dependencies in the buildpack.
	BuildpackCache layers.Layers

	// Logger logger is used to write debug and info to the console.
	Logger logger.Logger
}

// DependencyLayer returns a DependencyLayer unique to a dependency.
func (l Layers) DependencyLayer(dependency buildpack.Dependency) DependencyLayer {
	return DependencyLayer{
		l.Layer(dependency.ID),
		dependency,
		l.Logger,
		l.DownloadLayer(dependency),
	}
}

// DownloadLayer returns a DownloadLayer unique to a dependency.
func (l Layers) DownloadLayer(dependency buildpack.Dependency) DownloadLayer {
	return DownloadLayer{
		l.Layer(dependency.SHA256),
		Layer{l.BuildpackCache.Layer(dependency.SHA256), l.Logger},
		dependency,
		l.Logger,
	}
}

// Layer creates a Layer with a specified name.
func (l Layers) Layer(name string) Layer {
	return Layer{l.Layers.Layer(name), l.Logger}
}

// String makes Layers satisfy the Stringer interface.
func (l Layers) String() string {
	return fmt.Sprintf("Layers{ Layers: %s, Logger: %s }",
		l.Layers, l.Logger)
}

// WriteMetadata writes Launch metadata to the filesystem.
func (l Layers) WriteMetadata(metadata Metadata) error {
	l.Logger.FirstLine("Process types:")

	max := l.maximumTypeLength(metadata)

	for _, t := range metadata.Processes {
		s := color.CyanString(t.Type) + ":"

		for i := 0; i < (max - len(t.Type)); i++ {
			s += " "
		}

		l.Logger.SubsequentLine("%s %s", s, t.Command)
	}

	return l.Layers.WriteMetadata(metadata)
}

func (l Layers) maximumTypeLength(metadata Metadata) int {
	max := 0

	for _, t := range metadata.Processes {
		l := len(t.Type)

		if l > max {
			max = l
		}
	}

	return max
}
