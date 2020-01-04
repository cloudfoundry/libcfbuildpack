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

package helper_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestExtractTarGz(t *testing.T) {
	spec.Run(t, "ExtractTarGz", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var root string

		it.Before(func() {
			root = test.ScratchDir(t, "extract-tar-gz")
		})

		it("extracts the archive", func() {
			g.Expect(helper.ExtractTarGz(filepath.Join("testdata", "test-archive.tar.gz"), root, 0)).To(gomega.Succeed())
			g.Expect(filepath.Join(root, "fileA.txt")).To(gomega.BeARegularFile())
			g.Expect(filepath.Join(root, "dirA", "fileB.txt")).To(gomega.BeARegularFile())
			g.Expect(filepath.Join(root, "dirA", "fileC.txt")).To(gomega.BeARegularFile())
		})

		it("skips stripped components", func() {
			g.Expect(helper.ExtractTarGz(filepath.Join("testdata", "test-archive.tar.gz"), root, 1)).To(gomega.Succeed())
			g.Expect(filepath.Join(root, "fileB.txt")).To(gomega.BeARegularFile())
			g.Expect(filepath.Join(root, "fileC.txt")).To(gomega.BeARegularFile())
		})
	}, spec.Report(report.Terminal{}))
}
