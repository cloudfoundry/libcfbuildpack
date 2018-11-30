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
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/h2non/gock"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDownloadLayer(t *testing.T) {
	spec.Run(t, "DownloadLayer", testDownloadLayer, spec.Report(report.Terminal{}))
}

func testDownloadLayer(t *testing.T, when spec.G, it spec.S) {

	it("creates a download layer with the dependency SHA256 name", func() {
		root := internal.ScratchDir(t, "download-layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}
		dependency := buildpack.Dependency{SHA256: "test-sha256"}

		l := layers.DownloadLayer(dependency)

		expected := filepath.Join(root, "test-sha256")
		if l.Root != expected {
			t.Errorf("DownloadLayer.Root = %s, expected %s", l.Root, expected)
		}
	})

	it("downloads a dependency", func() {
		root := internal.ScratchDir(t, "download-layer")
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}}

		dependency := buildpack.Dependency{
			Version: newVersion(t, "1.0"),
			SHA256:  "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
			URI:     "http://test.com/test-path",
		}

		defer gock.Off()

		gock.New("http://test.com").
			Get("/test-path").
			Reply(200).
			BodyString("test-payload")

		a, err := layers.DownloadLayer(dependency).Artifact()
		if err != nil {
			t.Fatal(err)
		}

		expected := filepath.Join(root, dependency.SHA256, "test-path")
		if a != expected {
			t.Errorf("DownloadLayer.Artifact() = %s, expected %s", a, expected)
		}

		test.BeFileLike(t, expected, 0644, "test-payload")

		expected = filepath.Join(root, fmt.Sprintf("%s.toml", dependency.SHA256))
		test.BeFileLike(t, expected, 0644, `build = false
cache = true
launch = false

[metadata]
  id = ""
  name = ""
  version = "1.0"
  uri = "http://test.com/test-path"
  sha256 = "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273"
`)
	})

	it("does not download a buildpack cached dependency", func() {
		root := internal.ScratchDir(t, "download-layer")
		layers := layersCf.Layers{
			Layers:         layersBp.Layers{Root: root},
			BuildpackCache: layersBp.Layers{Root: filepath.Join(root, "buildpack")},
		}

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
`), filepath.Join(root, "buildpack", fmt.Sprintf("%s.toml", dependency.SHA256)), 0644); err != nil {
			t.Fatal(err)
		}

		a, err := layers.DownloadLayer(dependency).Artifact()
		if err != nil {
			t.Fatal(err)
		}

		expected := filepath.Join(root, "buildpack", dependency.SHA256, "test-path")
		if a != expected {
			t.Errorf("DownloadLayer.Artifact() = %s, expected %s", a, expected)
		}
	})

	it("does not download a previously cached dependency", func() {
		root := internal.ScratchDir(t, "download-layer")
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

		a, err := layers.DownloadLayer(dependency).Artifact()
		if err != nil {
			t.Fatal(err)
		}

		expected := filepath.Join(root, dependency.SHA256, "test-path")
		if a != expected {
			t.Errorf("DownloadLayer.Artifact() = %s, expected %s", a, expected)
		}
	})
}
