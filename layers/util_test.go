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

	"github.com/Masterminds/semver"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUtil(t *testing.T) {
	spec.Run(t, "Util", testUtil, spec.Random(), spec.Report(report.Terminal{}))
}

func testUtil(t *testing.T, when spec.G, it spec.S) {

	when("ExtractTarGz", func() {

		it("extracts the archive", func() {
			root := internal.ScratchDir(t, "util")

			err := layers.ExtractTarGz(internal.FixturePath(t, "test-archive.tar.gz"), root, 0)
			if err != nil {
				t.Fatal(err)
			}

			test.BeFileLike(t, filepath.Join(root, "fileA.txt"), 0644, "")
			test.BeFileLike(t, filepath.Join(root, "dirA", "fileB.txt"), 0644, "")
			test.BeFileLike(t, filepath.Join(root, "dirA", "fileC.txt"), 0644, "")
		})

		it("skips stripped components", func() {
			root := internal.ScratchDir(t, "util")

			err := layers.ExtractTarGz(internal.FixturePath(t, "test-archive.tar.gz"), root, 1)
			if err != nil {
				t.Fatal(err)
			}

			exists, err := layers.FileExists(filepath.Join(root, "fileA.txt"))
			if err != nil {
				t.Fatal(err)
			}

			if exists {
				t.Errorf("fileA.txt exists, expected not to")
			}

			test.BeFileLike(t, filepath.Join(root, "fileB.txt"), 0644, "")
			test.BeFileLike(t, filepath.Join(root, "fileC.txt"), 0644, "")
		})

	},spec.Random())

	when("ExtractZip", func() {

		it("extracts the archive", func() {
			root := internal.ScratchDir(t, "util")

			err := layers.ExtractZip(internal.FixturePath(t, "test-archive.zip"), root, 0)
			if err != nil {
				t.Fatal(err)
			}

			test.BeFileLike(t, filepath.Join(root, "fileA.txt"), 0644, "")
			test.BeFileLike(t, filepath.Join(root, "dirA", "fileB.txt"), 0644, "")
			test.BeFileLike(t, filepath.Join(root, "dirA", "fileC.txt"), 0644, "")
		})

		it("skips stripped components", func() {
			root := internal.ScratchDir(t, "util")

			err := layers.ExtractZip(internal.FixturePath(t, "test-archive.zip"), root, 1)
			if err != nil {
				t.Fatal(err)
			}

			exists, err := layers.FileExists(filepath.Join(root, "fileA.txt"))
			if err != nil {
				t.Fatal(err)
			}

			if exists {
				t.Errorf("fileA.txt exists, expected not to")
			}

			test.BeFileLike(t, filepath.Join(root, "fileB.txt"), 0644, "")
			test.BeFileLike(t, filepath.Join(root, "fileC.txt"), 0644, "")
		})

	}, spec.Random())

}

func newVersion(t *testing.T, version string) buildpack.Version {
	t.Helper()

	v, err := semver.NewVersion(version)
	if err != nil {
		t.Fatal(err)
	}

	return buildpack.Version{Version: v}
}
