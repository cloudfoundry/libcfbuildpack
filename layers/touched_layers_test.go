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
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestTouchedLayers(t *testing.T) {
	spec.Run(t, "TouchedLayers", testTouchedLayers, spec.Report(report.Terminal{}))
}

func testTouchedLayers(t *testing.T, when spec.G, it spec.S) {

	it("does not remove touched layers", func() {
		root := internal.ScratchDir(t, "touched-layers")
		test.TouchFile(t, filepath.Join(root, "test-layer.toml"))

		tl := layers.TouchedLayers{Root: root, Touched: make(map[string]struct{})}
		tl.Add(filepath.Join(root, "test-layer.toml"))
		if err := tl.Cleanup(); err != nil {
			t.Fatal(err)
		}

		test.FileExists(t, filepath.Join(root, "test-layer.toml"))
	})

	it("removes untouched layers", func() {
		root := internal.ScratchDir(t, "touched-layers")
		test.TouchFile(t, filepath.Join(root, "test-layer.toml"))

		tl := layers.TouchedLayers{Root: root, Touched: make(map[string]struct{})}
		if err := tl.Cleanup(); err != nil {
			t.Fatal(err)
		}

		test.FileDoesNotExist(t, filepath.Join(root, "test-layer.toml"))
	})
}
