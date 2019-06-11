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

package buildpack_test

import (
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLicenses(t *testing.T) {
	spec.Run(t, "Licenses", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("does not validate if there is not at least one license", func() {
			g.Expect(buildpack.Licenses{}.Validate()).NotTo(Succeed())
		})

		it("validates when all licenses are valid", func() {
			g.Expect(buildpack.Licenses{
				{Type: "test-type"},
				{URI: "test-uri"},
			}.Validate()).To(Succeed())
		})

		it("does not validate when a license is invalid", func() {
			g.Expect(buildpack.Licenses{
				{Type: "test-type"},
				{},
			}.Validate()).NotTo(Succeed())
		})
	}, spec.Report(report.Terminal{}))
}
