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

package layers

import (
	"fmt"
	"sort"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/fatih/color"
)

// Layers is an extension allows additional functionality to be added.
type Layers struct {
	layers.Layers

	// DependencyBuildPlans contains all contributed dependencies.
	DependencyBuildPlans buildplan.BuildPlan

	// TouchedLayers registers the layers that have been touched during this execution.
	TouchedLayers TouchedLayers

	buildpackCache layers.Layers
	info           buildpack.Info
	logger         logger.Logger
}

// DependencyLayer returns a DependencyLayer unique to a dependency.
func (l Layers) DependencyLayer(dependency buildpack.Dependency) DependencyLayer {
	return DependencyLayer{
		l.Layer(dependency.ID),
		dependency,
		l.DependencyBuildPlans,
		l.DownloadLayer(dependency),
		l.logger,
	}
}

// DownloadLayer returns a DownloadLayer unique to a dependency.
func (l Layers) DownloadLayer(dependency buildpack.Dependency) DownloadLayer {
	return DownloadLayer{
		l.Layer(dependency.SHA256),
		Layer{l.buildpackCache.Layer(dependency.SHA256), l.logger, l.TouchedLayers},
		dependency,
		l.info,
		l.logger,
	}
}

// Layer creates a Layer with a specified name.
func (l Layers) Layer(name string) Layer {
	return Layer{l.Layers.Layer(name), l.logger, l.TouchedLayers}
}

// String makes Layers satisfy the Stringer interface.
func (l Layers) String() string {
	return fmt.Sprintf("Layers{ Layers: %s, DependencyBuildPlans: %s, TouchedLayers: %s, buildpackCache: %s, info :%s, logger: %s }",
		l.Layers, l.DependencyBuildPlans, l.TouchedLayers, l.buildpackCache, l.info, l.logger)
}

// WriteApplicationMetadata writes application metadata to the filesystem.
func (l Layers) WriteApplicationMetadata(metadata Metadata) error {
	if len(metadata.Slices) > 0 {
		l.logger.FirstLine("%d application slices", len(metadata.Slices))
	}

	if len(metadata.Processes) > 0 {
		l.logger.FirstLine("Process types:")

		p := metadata.Processes
		sort.Slice(p, func(i int, j int) bool {
			return p[i].Type < p[j].Type
		})

		max := l.maximumTypeLength(p)
		for _, p := range p {
			format := fmt.Sprintf("%%s:%%-%ds %%s", max-len(p.Type))
			l.logger.SubsequentLine(format, color.CyanString(p.Type), "", p.Command)
		}
	}

	return l.Layers.WriteApplicationMetadata(metadata)
}

func (l Layers) maximumTypeLength(processes Processes) int {
	max := 0

	for _, t := range processes {
		if l := len(t.Type); l > max {
			max = l
		}
	}

	return max
}

// NewLayers creates a new instance of Layers.
func NewLayers(layers layers.Layers, buildpackCache layers.Layers, info buildpack.Info, logger logger.Logger) Layers {
	return Layers{
		Layers:               layers,
		DependencyBuildPlans: make(buildplan.BuildPlan),
		TouchedLayers:        NewTouchedLayers(layers.Root, logger),
		buildpackCache:       buildpackCache,
		info:                 info,
		logger:               logger,
	}
}
