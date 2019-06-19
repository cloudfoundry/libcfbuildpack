/*
 * Copyright 2018-2019 the original author or authors.
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

func TestMultiDependencyLayer(t *testing.T) {
	spec.Run(t, "MultiDependencyLayer", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var (
			root         string
			dependencies []buildpack.Dependency
			ls           layers.Layers
			layer        layers.MultiDependencyLayer
		)

		it.Before(func() {
			root = test.ScratchDir(t, "download-layer")

			dependencies = append(dependencies,
				buildpack.Dependency{
					ID:      "test-id-1",
					Version: internal.NewTestVersion(t, "1.0"),
					SHA256:  "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
					URI:     "https://test.com/test-path",
				},
				buildpack.Dependency{
					ID:      "test-id-2",
					Version: internal.NewTestVersion(t, "1.0"),
					SHA256:  "7f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
					URI:     "https://test.com/test-path",
				})

			ls = layers.NewLayers(layersBp.Layers{Root: root}, layersBp.Layers{}, buildpack.Buildpack{}, logger.Logger{})
			layer = ls.MultiDependencyLayer("test-name", dependencies)
		})

		it("creates a multi-dependency layer with the name", func() {
			g.Expect(layer.Root).To(Equal(filepath.Join(root, "test-name")))
		})

		it("calls contributors to contribute multi-dependency layer", func() {
			test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", dependencies[0].SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependencies[0].ID, dependencies[0].Version.Original(), dependencies[0].SHA256, dependencies[0].URI)

			test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", dependencies[1].SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependencies[1].ID, dependencies[1].Version.Original(), dependencies[1].SHA256, dependencies[1].URI)

			contributed := []bool{false, false}
			g.Expect(layer.Contribute(map[string]layers.MultiDependencyLayerContributor{
				"test-id-1": func(artifact string, layer layers.MultiDependencyLayer) error {
					contributed[0] = true
					return nil
				},
				"test-id-2": func(artifact string, layer layers.MultiDependencyLayer) error {
					contributed[1] = true
					return nil
				},
			})).To(Succeed())

			g.Expect(contributed[0]).To(BeTrue())
			g.Expect(contributed[1]).To(BeTrue())
		})

		it("does not call contributors for a cached layer", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
[[metadata.dependencies]]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"

[[metadata.dependencies]]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`,
				dependencies[0].ID, dependencies[0].Version.Original(), dependencies[0].SHA256, dependencies[0].URI,
				dependencies[1].ID, dependencies[1].Version.Original(), dependencies[1].SHA256, dependencies[1].URI)

			contributed := []bool{false, false}
			g.Expect(layer.Contribute(map[string]layers.MultiDependencyLayerContributor{
				"test-id-1": func(artifact string, layer layers.MultiDependencyLayer) error {
					contributed[0] = true
					return nil
				},
				"test-id-2": func(artifact string, layer layers.MultiDependencyLayer) error {
					contributed[1] = true
					return nil
				},
			})).To(Succeed())

			g.Expect(contributed[0]).To(BeFalse())
			g.Expect(contributed[1]).To(BeFalse())
		})

		it("contributes dependencies to build plan", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
[[metadata.dependencies]]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"

[[metadata.dependencies]]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`,
				dependencies[0].ID, dependencies[0].Version.Original(), dependencies[0].SHA256, dependencies[0].URI,
				dependencies[1].ID, dependencies[1].Version.Original(), dependencies[1].SHA256, dependencies[1].URI)

			g.Expect(layer.Contribute(map[string]layers.MultiDependencyLayerContributor{})).To(Succeed())

			g.Expect(ls.DependencyBuildPlans).To(Equal(buildplan.BuildPlan{
				dependencies[0].ID: buildplan.Dependency{
					Version: "1.0",
					Metadata: buildplan.Metadata{
						"name":     dependencies[0].Name,
						"uri":      dependencies[0].URI,
						"sha256":   dependencies[0].SHA256,
						"stacks":   dependencies[0].Stacks,
						"licenses": dependencies[0].Licenses,
					},
				},
				dependencies[1].ID: buildplan.Dependency{
					Version: "1.0",
					Metadata: buildplan.Metadata{
						"name":     dependencies[1].Name,
						"uri":      dependencies[1].URI,
						"sha256":   dependencies[1].SHA256,
						"stacks":   dependencies[1].Stacks,
						"licenses": dependencies[1].Licenses,
					},
				},
			}))
		})

		it("cleans layer when contributing dependency layer", func() {
			test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", dependencies[0].SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependencies[0].ID, dependencies[0].Version.Original(), dependencies[0].SHA256, dependencies[0].URI)

			test.WriteFile(t, filepath.Join(root, fmt.Sprintf("%s.toml", dependencies[1].SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependencies[1].ID, dependencies[1].Version.Original(), dependencies[1].SHA256, dependencies[1].URI)
			test.TouchFile(t, layer.Root, "test-file")

			g.Expect(layer.Contribute(map[string]layers.MultiDependencyLayerContributor{
				"test-id-1": func(artifact string, layer layers.MultiDependencyLayer) error {
					return nil
				},
				"test-id-2": func(artifact string, layer layers.MultiDependencyLayer) error {
					return nil
				},
			})).To(Succeed())

			g.Expect(filepath.Join(layer.Root, "test-file")).NotTo(BeAnExistingFile())
		})
	}, spec.Report(report.Terminal{}))
}
