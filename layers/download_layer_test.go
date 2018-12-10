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
	"net/http"
	"path/filepath"
	"testing"

	layersBp "github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDownloadLayer(t *testing.T) {
	spec.Run(t, "DownloadLayer", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var (
			root       string
			dependency buildpack.Dependency
			layer      layers.DownloadLayer
			server     *ghttp.Server
		)

		it.Before(func() {
			root = internal.ScratchDir(t, "download-layer")

			server = ghttp.NewServer()

			dependency = buildpack.Dependency{
				ID:      "test-id",
				Version: internal.NewTestVersion(t, "1.0"),
				SHA256:  "6f06dd0e26608013eff30bb1e951cda7de3fdd9e78e907470e0dd5c0ed25e273",
				URI:     fmt.Sprintf("%s/test-path", server.URL()),
			}

			layers := layers.NewLayers(layersBp.Layers{Root: root}, layersBp.Layers{Root: filepath.Join(root, "buildpack")}, logger.Logger{})
			layer = layers.DownloadLayer(dependency)
		})

		it.After(func() {
			server.Close()
		})

		it("creates a download layer with the dependency SHA256 name", func() {
			g.Expect(layer.Root).To(Equal(filepath.Join(root, dependency.SHA256)))
		})

		it("downloads a dependency", func() {
			server.AppendHandlers(ghttp.RespondWith(http.StatusOK, "test-payload"))

			g.Expect(layer.Artifact()).To(SatisfyAll(
				Equal(filepath.Join(layer.Root, "test-path")),
				test.HaveContent("test-payload")))

			g.Expect(layer).To(test.HaveLayerMetadata(false, true, false))
		})

		it("does not download a buildpack cached dependency", func() {
			test.WriteFile(t, filepath.Join(root, "buildpack", fmt.Sprintf("%s.toml", dependency.SHA256)), `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)

			g.Expect(layer.Artifact()).To(Equal(filepath.Join(root, "buildpack", dependency.SHA256, "test-path")))
		})

		it("does not download a previously cached dependency", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
ID = "%s"
Version = "%s"
SHA256 = "%s"
URI = "%s"`, dependency.ID, dependency.Version.Original(), dependency.SHA256, dependency.URI)

			g.Expect(layer.Artifact()).To(Equal(filepath.Join(layer.Root, "test-path")))
		})
	}, spec.Report(report.Terminal{}))
}
