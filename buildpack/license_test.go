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

func TestLicense(t *testing.T) {
	spec.Run(t, "License", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("validates with type set", func() {
			g.Expect(buildpack.License{Type: "test-type"}.Validate()).To(Succeed())
		})

		it("validates with uri set", func() {
			g.Expect(buildpack.License{URI: "test-uri "}.Validate()).To(Succeed())
		})

		it("validates with type and uri set", func() {
			g.Expect(buildpack.License{Type: "test-type", URI: "test-uri "}.Validate()).To(Succeed())
		})

		it("does not validate without type and uri set", func() {
			g.Expect(buildpack.License{}.Validate()).NotTo(Succeed())
		})
	}, spec.Report(report.Terminal{}))
}
