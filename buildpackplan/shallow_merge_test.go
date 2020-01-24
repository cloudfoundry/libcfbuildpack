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

package buildpackplan_test

import (
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/v2/buildpackplan"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestShallowMerge(t *testing.T) {
	spec.Run(t, "ShallowMerge", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		when("version", func() {

			it("chooses neither", func() {
				a := buildpackplan.Plan{Name: "alpha"}
				b := buildpackplan.Plan{Name: "alpha"}

				expected := buildpackplan.Plan{Name: "alpha"}

				g.Expect(buildpackplan.ShallowMerge(a, b)).To(gomega.Equal(expected))
			})

			it("chooses a", func() {
				a := buildpackplan.Plan{Name: "alpha", Version: "version-1"}
				b := buildpackplan.Plan{Name: "alpha"}

				expected := buildpackplan.Plan{Name: "alpha", Version: "version-1"}

				g.Expect(buildpackplan.ShallowMerge(a, b)).To(gomega.Equal(expected))
			})

			it("chooses b", func() {
				a := buildpackplan.Plan{Name: "alpha"}
				b := buildpackplan.Plan{Name: "alpha", Version: "version-2"}

				expected := buildpackplan.Plan{Name: "alpha", Version: "version-2"}

				g.Expect(buildpackplan.ShallowMerge(a, b)).To(gomega.Equal(expected))
			})

			it("combines constrains", func() {
				a := buildpackplan.Plan{Name: "alpha", Version: "version-1"}
				b := buildpackplan.Plan{Name: "alpha", Version: "version-2"}

				expected := buildpackplan.Plan{Name: "alpha", Version: "version-1,version-2"}

				g.Expect(buildpackplan.ShallowMerge(a, b)).To(gomega.Equal(expected))
			})
		})

		when("metadata", func() {

			it("keeps a keys", func() {
				a := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-1": "value-1"}}
				b := buildpackplan.Plan{Name: "alpha"}

				expected := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-1": "value-1"}}

				g.Expect(buildpackplan.ShallowMerge(a, b)).To(gomega.Equal(expected))
			})

			it("keeps b keys", func() {
				a := buildpackplan.Plan{Name: "alpha"}
				b := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-2": "value-2"}}

				expected := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-2": "value-2"}}

				g.Expect(buildpackplan.ShallowMerge(a, b)).To(gomega.Equal(expected))
			})

			it("combines keys", func() {
				a := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-1": "value-1"}}
				b := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-2": "value-2"}}

				expected := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-1": "value-1", "key-2": "value-2"}}

				g.Expect(buildpackplan.ShallowMerge(a, b)).To(gomega.Equal(expected))
			})

			it("overwrites a keys with b keys", func() {
				a := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-1": "value-1"}}
				b := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-1": "value-2"}}

				expected := buildpackplan.Plan{Name: "alpha", Metadata: buildpackplan.Metadata{"key-1": "value-2"}}

				g.Expect(buildpackplan.ShallowMerge(a, b)).To(gomega.Equal(expected))
			})
		})
	}, spec.Report(report.Terminal{}))
}
