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

package layers_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	layersBp "github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDependencyLayer(t *testing.T) {
	spec.Run(t, "DependencyLayer", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

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
				URI:     "http://test.com/test-path",
			}

			ls = layers.NewLayers(layersBp.Layers{Root: root}, layersBp.Layers{}, buildpack.Info{}, logger.Logger{})
			layer = ls.DependencyLayer(dependency)
		})

		it("creates a dependency layer with the dependency id name", func() {
			g.Expect(layer.Root).To(Equal(filepath.Join(root, dependency.ID)))
		})

		it("calls contributor to contribute dependency layer", func() {
			test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", dependency.SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)

			contributed := false
			g.Expect(layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
				contributed = true;
				return nil
			})).To(Succeed())

			g.Expect(contributed).To(BeTrue())
		})

		it("does not call contributor for a cached layer", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)

			contributed := false
			g.Expect(layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
				contributed = true;
				return nil
			})).To(Succeed())

			g.Expect(contributed).To(BeFalse())
		})

		it("returns artifact name", func() {
			g.Expect(layer.ArtifactName()).To(Equal("test-path"))
		})

		it("contributes dependency to build plan", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)

			g.Expect(layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
				return nil
			})).To(Succeed())

			g.Expect(ls.DependencyBuildPlans).To(Equal(buildplan.BuildPlan{
				dependency.ID: buildplan.Dependency{
					Version: "1.0",
					Metadata: buildplan.Metadata{
						"name":     dependency.Name,
						"uri":      dependency.URI,
						"sha256":   dependency.SHA256,
						"stacks":   dependency.Stacks,
						"licenses": dependency.Licenses,
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
			})).To(Succeed())

			g.Expect(filepath.Join(layer.Root, "test-file")).NotTo(BeAnExistingFile())
		})
	}, spec.Report(report.Terminal{}))
}
