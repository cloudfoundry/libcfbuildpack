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
	"path/filepath"

	"github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/logger"
)

// DependencyLayer is an extension to Layer that is unique to a dependency.
type DependencyLayer struct {
	Layer

	// Dependency is the dependency provided by this layer
	Dependency buildpack.Dependency

	// Logger is used to write debug and info to the console.
	Logger logger.Logger

	downloadLayer DownloadLayer
}

// ArtifactName returns the name portion of the download path for the dependency.
func (l DependencyLayer) ArtifactName() string {
	return filepath.Base(l.Dependency.URI)
}

// DependencyLayerContributor defines a callback function that is called when a dependency needs to be contributed.
type DependencyLayerContributor func(artifact string, layer DependencyLayer) error

// Contribute facilitates custom contribution of an artifact to a layer.  If the artifact has already been contributed,
// the contribution is validated and the contributor is not called.
func (l DependencyLayer) Contribute(contributor DependencyLayerContributor, flags ...layers.Flag) error {
	return l.Layer.Contribute(l.Dependency, func(layer Layer) error {
		a, err := l.downloadLayer.Artifact()
		if err != nil {
			return err
		}

		return contributor(a, l)
	}, flags...)
}

// String makes DependencyLayer satisfy the Stringer interface.
func (l DependencyLayer) String() string {
	return fmt.Sprintf("DependencyLayer{ Layer: %s, Dependency: %s, Logger: %s, downloadLayer: %s }",
		l.Layer, l.Dependency, l.Logger, l.downloadLayer)
}
