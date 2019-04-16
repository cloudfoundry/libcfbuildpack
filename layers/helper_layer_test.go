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
	"path/filepath"
	"testing"

	buildpackBp "github.com/buildpack/libbuildpack/buildpack"
	"github.com/buildpack/libbuildpack/buildplan"
	layersBp "github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestHelperLayer(t *testing.T) {
	spec.Run(t, "HelperLayer", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var (
			bp    buildpack.Buildpack
			root  string
			id    string
			ls    layers.Layers
			layer layers.HelperLayer
		)

		it.Before(func() {
			root = test.ScratchDir(t, "helper-layer")

			bp = buildpack.Buildpack{
				Buildpack: buildpackBp.Buildpack{
					Info: buildpackBp.Info{
						ID:      "test-id",
						Name:    "test-name",
						Version: "test-version",
					},
				},
			}

			id = "test-id"

			ls = layers.NewLayers(layersBp.Layers{Root: root}, layersBp.Layers{}, bp, logger.Logger{})
			layer = ls.HelperLayer(id, "Test Name")
		})

		it("creates a helper layer with the helper id name", func() {
			g.Expect(layer.Root).To(Equal(filepath.Join(root, id)))
		})

		it("calls contributor to contribute dependency layer", func() {
			contributed := false
			g.Expect(layer.Contribute(func(artifact string, layer layers.HelperLayer) error {
				contributed = true
				return nil
			})).To(Succeed())

			g.Expect(contributed).To(BeTrue())
		})

		it("does not call contributor for a cached layer", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
  id = "%s"
  name = "%s"
  version = "%s"
  display_name = "Test Name"`, bp.Info.ID, bp.Info.Name, bp.Info.Version)

			contributed := false
			g.Expect(layer.Contribute(func(artifact string, layer layers.HelperLayer) error {
				contributed = true
				return nil
			})).To(Succeed())

			g.Expect(contributed).To(BeFalse())
		})

		it("contributes dependency to build plan", func() {
			g.Expect(layer.Contribute(func(artifact string, layer layers.HelperLayer) error {
				return nil
			})).To(Succeed())

			g.Expect(ls.DependencyBuildPlans).To(Equal(buildplan.BuildPlan{
				id: buildplan.Dependency{
					Version: bp.Info.Version,
					Metadata: buildplan.Metadata{
						"id":   bp.Info.ID,
						"name": bp.Info.Name,
					},
				},
			}))
		})

		it("cleans layer when contributing dependency layer", func() {
			test.TouchFile(t, layer.Root, "test-file")

			g.Expect(layer.Contribute(func(artifact string, layer layers.HelperLayer) error {
				return nil
			})).To(Succeed())

			g.Expect(filepath.Join(layer.Root, "test-file")).NotTo(BeAnExistingFile())
		})
	}, spec.Report(report.Terminal{}))
}
