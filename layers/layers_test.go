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
	"bytes"
	"fmt"
	"path/filepath"
	"testing"

	layersBp "github.com/buildpacks/libbuildpack/v2/layers"
	loggerBp "github.com/buildpacks/libbuildpack/v2/logger"
	"github.com/cloudfoundry/libcfbuildpack/v2/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/v2/layers"
	"github.com/cloudfoundry/libcfbuildpack/v2/logger"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"
	"github.com/heroku/color"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayers(t *testing.T) {
	spec.Run(t, "Layers", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			root string
			info bytes.Buffer
			l    layers.Layers
		)

		it.Before(func() {
			root = test.ScratchDir(t, "layers")
			logger := logger.Logger{Logger: loggerBp.NewLogger(nil, &info)}
			l = layers.NewLayers(layersBp.Layers{Root: root}, layersBp.Layers{}, buildpack.Buildpack{}, logger)
		})

		it("logs process types", func() {
			g.Expect(l.WriteApplicationMetadata(layers.Metadata{
				Processes: []layers.Process{
					{Type: "short", Command: "test-command-1"},
					{Type: "a-very-long-type", Command: "test-command-2"},
				},
			})).To(gomega.Succeed())

			actual := info.String()
			expected := fmt.Sprintf(`  Process types:
    %s: test-command-2
    %s:            test-command-1
`, color.CyanString("a-very-long-type"), color.CyanString("short"))
			g.Expect(actual).To(gomega.Equal(expected))
		})

		it("logs number of slices", func() {
			g.Expect(l.WriteApplicationMetadata(layers.Metadata{
				Slices: layers.Slices{
					layers.Slice{},
					layers.Slice{},
				},
			})).To(gomega.Succeed())

			g.Expect(info.String()).To(gomega.Equal("  2 application slices\n"))
		})

		it("registers touched layers", func() {
			test.TouchFile(t, l.Root, "test-layer-1.toml")
			test.TouchFile(t, l.Root, "test-layer-2.toml")

			g.Expect(l.Layer("test-layer-1").Contribute(nil, func(layer layers.Layer) error {
				return nil
			})).To(gomega.Succeed())

			g.Expect(l.TouchedLayers.Cleanup()).To(gomega.Succeed())
			g.Expect(filepath.Join(l.Root, "test-layer-1.toml")).To(gomega.BeAnExistingFile())
			g.Expect(filepath.Join(l.Root, "test-layer-2.toml")).NotTo(gomega.BeAnExistingFile())
		})
	}, spec.Report(report.Terminal{}))
}
