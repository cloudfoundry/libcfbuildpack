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

package layers_test

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	layersBp "github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	layersCf "github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDependencyLayer(t *testing.T) {
	spec.Run(t, "DependencyLayer", testDependencyLayer, spec.Report(report.Terminal{}))
}

func testDependencyLayer(t *testing.T, when spec.G, it spec.S) {

	it("creates a dependency later with the dependency id name", func() {
		root := internal.ScratchDir(t, "dependency-layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}
		dependency := buildpack.Dependency{ID: "test-id"}

		l := layers.DependencyLayer(dependency)

		expected := filepath.Join(root, "test-id")
		if l.Root != expected {
			t.Errorf("DependencyLayer.Root = %s, expected %s", l.Root, expected)
		}
	})

	it("calls contributor to contribute dependency layer", func() {
		root := internal.ScratchDir(t, "dependency-layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}

		dependency := buildpack.Dependency{
			Version: newVersion(t, "1.0"),
			SHA256:  "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
			URI:     "http://test.com/test-path",
		}

		if err := layersCf.WriteToFile(strings.NewReader(`[metadata]
  id = ""
  name = ""
  version = "1.0"
  uri = "http://test.com/test-path"
  sha256 = "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273"
`), filepath.Join(root, fmt.Sprintf("%s.toml", dependency.SHA256)), 0644); err != nil {
			t.Fatal(err)
		}

		contributed := false

		if err := layers.DependencyLayer(dependency).Contribute(func(artifact string, layer layersCf.DependencyLayer) error {
			contributed = true;
			return nil
		}); err != nil {
			t.Fatal(err)
		}

		if !contributed {
			t.Errorf("Expected contribution but didn't contribute")
		}
	})

	it("does not call contributor for a cached launch layer", func() {
		root := internal.ScratchDir(t, "dependency-layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}

		dependency := buildpack.Dependency{
			Version: newVersion(t, "1.0"),
			SHA256:  "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
			URI:     "http://test.com/test-path",
		}

		if err := layersCf.WriteToFile(strings.NewReader(`[metadata]
  id = ""
  name = ""
  version = "1.0"
  uri = "http://test.com/test-path"
  sha256 = "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273"
`), filepath.Join(root, fmt.Sprintf("%s.toml", dependency.ID)), 0644); err != nil {
			t.Fatal(err)
		}

		contributed := false

		if err := layers.DependencyLayer(dependency).Contribute(func(artifact string, layer layersCf.DependencyLayer) error {
			contributed = true;
			return nil
		}); err != nil {
			t.Fatal(err)
		}

		if contributed {
			t.Errorf("Expected non-contribution but did contribute")
		}
	})

	it("returns artifact name", func() {
		root := internal.ScratchDir(t, "dependency-layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}

		dependency := buildpack.Dependency{ID: "test-id", URI: "http://localhost/path/test-artifact-name"}

		l := layers.DependencyLayer(dependency)

		if l.ArtifactName() != "test-artifact-name" {
			t.Errorf("DependencyLaunchLayer.ArtifactName = %s, expected test-artifact-name", l.ArtifactName())
		}
	})
}
