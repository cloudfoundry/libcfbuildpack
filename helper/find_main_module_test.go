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

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestFindMainModule(t *testing.T) {
	spec.Run(t, "FindMainModule", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			app  application.Application
			root string
		)

		it.Before(func() {
			root = test.ScratchDir(t, "find-main-module")
			app = application.Application{Root: root}
		})

		it("returns false if no package.json", func() {
			_, ok, err := helper.FindMainModule(app)

			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		})

		it("returns false if no main", func() {
			test.WriteFile(t, filepath.Join(root, "package.json"), `{ }`)

			_, ok, err := helper.FindMainModule(app)

			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		})

		it("returns false if main file does not exist", func() {
			test.WriteFile(t, filepath.Join(root, "package.json"), `{ "main": "test.js" }`)

			_, ok, err := helper.FindMainModule(app)

			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		})

		it("returns true if main file does exist", func() {
			test.WriteFile(t, filepath.Join(root, "package.json"), `{ "main": "test.js" }`)
			test.TouchFile(t, root, "test.js")

			_, ok, err := helper.FindMainModule(app)

			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())
		})
	}, spec.Report(report.Terminal{}))
}
