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

func TestDependencyLayer(t *testing.T) {
	spec.Run(t, "DependencyLayer", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			root       string
			dependency buildpack.Dependency
			ls         layers.Layers
			layer      layers.DependencyLayer
		)

		it.Before(func() {
			root = test.ScratchDir(t, "download-layer")

			dependency = buildpack.Dependency{
				ID:      "test-id",
				Version: internal.NewTestVersion(t, "1.0"),
				SHA256:  "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
				URI:     "https://test.com/test-path",
			}

			ls = layers.NewLayers(layersBp.Layers{Root: root}, layersBp.Layers{}, buildpack.Buildpack{}, logger.Logger{})
			layer = ls.DependencyLayer(dependency)
		})

		it("creates a dependency layer with the dependency id name", func() {
			g.Expect(layer.Root).To(gomega.Equal(filepath.Join(root, dependency.ID)))
		})

		it("calls contributor to contribute dependency layer", func() {
			test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", dependency.SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)

			contributed := false
			g.Expect(layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
				contributed = true
				return nil
			})).To(gomega.Succeed())

			g.Expect(contributed).To(gomega.BeTrue())
		})

		it("does not call contributor for a cached layer", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)

			contributed := false
			g.Expect(layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
				contributed = true
				return nil
			})).To(gomega.Succeed())

			g.Expect(contributed).To(gomega.BeFalse())
		})

		it("returns artifact name", func() {
			g.Expect(layer.ArtifactName()).To(gomega.Equal("test-path"))
		})

		it("contributes dependency to build plan", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)

			g.Expect(layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
				return nil
			})).To(gomega.Succeed())

			g.Expect(*ls.Plans).To(gomega.Equal(buildpackplan.Plans{
				Plans: buildpackplanBp.Plans{
					Entries: []buildpackplan.Plan{
						{
							Name:    dependency.ID,
							Version: "1.0",
							Metadata: buildpackplan.Metadata{
								"name":     dependency.Name,
								"uri":      dependency.URI,
								"sha256":   dependency.SHA256,
								"stacks":   dependency.Stacks,
								"licenses": dependency.Licenses,
							},
						},
					},
				},
			}))
		})

		it("cleans layer when contributing dependency layer", func() {
			test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", dependency.SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)
			test.TouchFile(t, layer.Root, "test-file")

			g.Expect(layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
				return nil
			})).To(gomega.Succeed())

			g.Expect(filepath.Join(layer.Root, "test-file")).NotTo(gomega.BeAnExistingFile())
		})
	}, spec.Report(report.Terminal{}))
}
