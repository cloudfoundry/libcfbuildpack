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
	"path/filepath"
	"strings"
	"testing"

	layersBp "github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	layersCf "github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayer(t *testing.T) {
	spec.Run(t, "Layer", testLayer, spec.Report(report.Terminal{}))
}

func testLayer(t *testing.T, when spec.G, it spec.S) {

	type metadata struct {
		Alpha string
		Bravo int
	}

	it("identifies matching metadata", func() {
		root := internal.ScratchDir(t, "layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}

		if err := layersCf.WriteToFile(strings.NewReader(`[metadata]
Alpha = "test-value"
Bravo = 1
`), filepath.Join(root, "test-layer.toml"), 0644); err != nil {
			t.Fatal(err)
		}
		matches, err := layers.Layer("test-layer").MetadataMatches(metadata{"test-value", 1})
		if err != nil {
			t.Fatal(err)
		}

		if !matches {
			t.Errorf("Layer.MetadataMatches() = %t, expected true", matches)
		}
	})

	it("identifies non-matching metadata", func() {
		root := internal.ScratchDir(t, "layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}

		if err := layersCf.WriteToFile(strings.NewReader(`[metadata]
Alpha = "test-value"
Bravo = 2
`), filepath.Join(root, "test-layer.toml"), 0644); err != nil {
			t.Fatal(err)
		}
		matches, err := layers.Layer("test-layer").MetadataMatches(metadata{"test-value", 1})
		if err != nil {
			t.Fatal(err)
		}

		if matches {
			t.Errorf("Layer.MetadataMatches() = %t, expected false", matches)
		}
	})

	it("identifies missing metadata", func() {
		root := internal.ScratchDir(t, "layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}

		matches, err := layers.Layer("test-layer").MetadataMatches(metadata{"test-value", 1})
		if err != nil {
			t.Fatal(err)
		}

		if matches {
			t.Errorf("Layer.MetadataMatches() = %t, expected false", matches)
		}
	})

	it("does not call contributor for cached layer", func() {
		root := internal.ScratchDir(t, "layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}
		layer := layers.Layer("test-layer")

		if err := layersCf.WriteToFile(strings.NewReader(`[metadata]
Alpha = "test-value"
Bravo = 1
`), filepath.Join(root, "test-layer.toml"), 0644); err != nil {
			t.Fatal(err)
		}

		contributed := false

		if err := layer.Contribute(metadata{"test-value", 1}, func(layer layersCf.Layer) error {
			contributed = true
			return nil
		}); err != nil {
			t.Fatal(err)
		}

		if contributed {
			t.Errorf("Expected non-contribution but did contribute")
		}
	})

	it("calls contributor for uncached layer", func() {
		root := internal.ScratchDir(t, "layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}

		contributed := false

		if err := layers.Layer("test-layer").Contribute(metadata{"test-value", 1}, func(layer layersCf.Layer) error {
			contributed = true
			return nil
		}); err != nil {
			t.Fatal(err)
		}

		if !contributed {
			t.Errorf("Expected contribution but didn't contribute")
		}
	})
}
