/*
 * Copyright 2018-2020 the original author or authors.
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

package layers_test

import (
	"fmt"
	"path/filepath"
	"testing"

	buildpackplanBp "github.com/buildpacks/libbuildpack/v2/buildpackplan"
	layersBp "github.com/buildpacks/libbuildpack/v2/layers"
	"github.com/cloudfoundry/libcfbuildpack/v2/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/v2/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/v2/internal"
	"github.com/cloudfoundry/libcfbuildpack/v2/layers"
	"github.com/cloudfoundry/libcfbuildpack/v2/logger"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMultiDependencyLayer(t *testing.T) {
	spec.Run(t, "MultiDependencyLayer", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			root         string
			dependencies []buildpack.Dependency
			ls           layers.Layers
			layer        layers.MultiDependencyLayer
		)

		it.Before(func() {
			root = test.ScratchDir(t, "download-layer")

			dependencies = []buildpack.Dependency{
				{
					ID:      "test-id-1",
					Version: internal.NewTestVersion(t, "1.0"),
					SHA256:  "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
					URI:     "https://test.com/test-path-1",
				},
				{
					ID:      "test-id-2",
					Version: internal.NewTestVersion(t, "2.0"),
					SHA256:  "5fe6dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e262",
					URI:     "https://test.com/test-path-2",
				},
			}

			ls = layers.NewLayers(layersBp.Layers{Root: root}, layersBp.Layers{}, buildpack.Buildpack{}, logger.Logger{})
			layer = ls.MultiDependencyLayer("test-id-0", dependencies...)
		})

		it("creates a dependency layer with the dependency id name", func() {
			g.Expect(layer.Root).To(gomega.Equal(filepath.Join(root, "test-id-0")))
		})

		it("calls contributor to contribute dependency layer", func() {
			for _, d := range dependencies {
				test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", d.SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, d.ID, d.Version.Original(), d.SHA256, d.URI)
			}

			contributed1 := false
			contributed2 := false
			g.Expect(layer.Contribute(map[string]layers.MultiDependencyLayerContributor{
				"test-id-1": func(artifact string, layer layers.MultiDependencyLayer) error {
					contributed1 = true
					return nil
				},
				"test-id-2": func(artifact string, layer layers.MultiDependencyLayer) error {
					contributed2 = true
					return nil
				},
			})).To(gomega.Succeed())

			g.Expect(contributed1).To(gomega.BeTrue())
			g.Expect(contributed2).To(gomega.BeTrue())
		})

		it("does not call contributor for a cached layer", func() {
			test.WriteFile(t, layer.Metadata, `[[metadata]]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"

[[metadata]]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`,
				dependencies[0].ID, dependencies[0].Version.Original(), dependencies[0].SHA256, dependencies[0].URI,
				dependencies[1].ID, dependencies[1].Version.Original(), dependencies[1].SHA256, dependencies[1].URI)

			contributed1 := false
			contributed2 := false
			g.Expect(layer.Contribute(map[string]layers.MultiDependencyLayerContributor{
				"test-id-1": func(artifact string, layer layers.MultiDependencyLayer) error {
					contributed1 = true
					return nil
				},
				"test-id-2": func(artifact string, layer layers.MultiDependencyLayer) error {
					contributed2 = true
					return nil
				},
			})).To(gomega.Succeed())

			g.Expect(contributed1).To(gomega.BeFalse())
			g.Expect(contributed2).To(gomega.BeFalse())
		})

		it("contributes dependency to build plan", func() {
			test.WriteFile(t, layer.Metadata, `[[metadata]]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"

[[metadata]]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`,
				dependencies[0].ID, dependencies[0].Version.Original(), dependencies[0].SHA256, dependencies[0].URI,
				dependencies[1].ID, dependencies[1].Version.Original(), dependencies[1].SHA256, dependencies[1].URI)

			g.Expect(layer.Contribute(map[string]layers.MultiDependencyLayerContributor{
				"test-id-1": func(artifact string, layer layers.MultiDependencyLayer) error {
					return nil
				},
				"test-id-2": func(artifact string, layer layers.MultiDependencyLayer) error {
					return nil
				},
			})).To(gomega.Succeed())

			g.Expect(*ls.Plans).To(gomega.Equal(buildpackplan.Plans{
				Plans: buildpackplanBp.Plans{
					Entries: []buildpackplan.Plan{
						{
							Name:    dependencies[0].ID,
							Version: "1.0",
							Metadata: buildpackplan.Metadata{
								"name":     dependencies[0].Name,
								"uri":      dependencies[0].URI,
								"sha256":   dependencies[0].SHA256,
								"stacks":   dependencies[0].Stacks,
								"licenses": dependencies[0].Licenses,
							},
						},
						{
							Name:    dependencies[1].ID,
							Version: "2.0",
							Metadata: buildpackplan.Metadata{
								"name":     dependencies[1].Name,
								"uri":      dependencies[1].URI,
								"sha256":   dependencies[1].SHA256,
								"stacks":   dependencies[1].Stacks,
								"licenses": dependencies[1].Licenses,
							},
						},
					},
				},
			}))
		})

		it("cleans layer when contributing dependency layer", func() {
			for _, d := range dependencies {
				test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", d.SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, d.ID, d.Version.Original(), d.SHA256, d.URI)
			}
			test.TouchFile(t, layer.Root, "test-file")

			g.Expect(layer.Contribute(map[string]layers.MultiDependencyLayerContributor{
				"test-id-1": func(artifact string, layer layers.MultiDependencyLayer) error {
					return nil
				},
				"test-id-2": func(artifact string, layer layers.MultiDependencyLayer) error {
					return nil
				},
			})).To(gomega.Succeed())

			g.Expect(filepath.Join(layer.Root, "test-file")).NotTo(gomega.BeAnExistingFile())
		})
	}, spec.Report(report.Terminal{}))
}
