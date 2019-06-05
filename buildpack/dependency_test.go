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

package buildpack_test

import (
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var TestDep = map[string]interface{}{
	"id":      "test-id-1",
	"name":    "test-name-1",
	"version": "1.0",
	"uri":     "test-uri-1",
	"sha256":  "test-sha256-1",
	"stacks":  []interface{}{"test-stack-1a", "test-stack-1b"},
	"licenses": []map[string]interface{}{
		{
			"type": "test-type-1",
			"uri":  "test-uri-1",
		},
		{
			"type": "test-type-2",
			"uri":  "test-uri-2",
		},
	},
}

func TestDependency(t *testing.T) {
	spec.Run(t, "Dependency", func(t *testing.T, when spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		when("NewDependency", func(){
			it("constructs a dependency", func(){
				expectedDep := buildpack.Dependency{
					ID:      "test-id-1",
					Name:    "test-name-1",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri-1",
					SHA256:  "test-sha256-1",
					Stacks:  buildpack.Stacks{"test-stack-1a", "test-stack-1b"},
					Licenses: buildpack.Licenses{
						buildpack.License{Type: "test-type-1", URI: "test-uri-1"},
						buildpack.License{Type: "test-type-2", URI: "test-uri-2"},
					},
				}

				g.Expect(buildpack.NewDependency(TestDep)).To(Equal(expectedDep))
			})
		})

		when("Validate", func(){
			it("validates", func() {
				g.Expect(buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack"},
					Licenses: buildpack.Licenses{
						{Type: "test-type"},
					},
				}.Validate()).To(Succeed())
			})

			it("does not validate with invalid id", func() {
				g.Expect(buildpack.Dependency{
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack"},
					Licenses: buildpack.Licenses{
						{Type: "test-type"},
					},
				}.Validate()).NotTo(Succeed())
			})

			it("does not validate with invalid name", func() {
				g.Expect(buildpack.Dependency{
					ID:      "test-id",
					Version: internal.NewTestVersion(t, "1.0.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack"},
					Licenses: buildpack.Licenses{
						{Type: "test-type"},
					},
				}.Validate()).NotTo(Succeed())
			})

			it("does not validate with invalid version", func() {
				g.Expect(buildpack.Dependency{
					ID:     "test-id",
					Name:   "test-name",
					URI:    "test-uri",
					SHA256: "test-sha256",
					Stacks: buildpack.Stacks{"test-stack"},
					Licenses: buildpack.Licenses{
						{Type: "test-type"},
					},
				}.Validate()).NotTo(Succeed())
			})

			it("does not validate with invalid uri", func() {
				g.Expect(buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0.0"),
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack"},
					Licenses: buildpack.Licenses{
						{Type: "test-type"},
					},
				}.Validate()).NotTo(Succeed())
			})

			it("does not validate with invalid sha256", func() {
				g.Expect(buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0.0"),
					URI:     "test-uri",
					Stacks:  buildpack.Stacks{"test-stack"},
					Licenses: buildpack.Licenses{
						{Type: "test-type"},
					},
				}.Validate()).NotTo(Succeed())
			})

			it("does not validate with invalid stacks", func() {
				g.Expect(buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Licenses: buildpack.Licenses{
						{Type: "test-type"},
					},
				}.Validate()).NotTo(Succeed())
			})

			it("does not validate with invalid licenses", func() {
				g.Expect(buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack"},
				}.Validate()).NotTo(Succeed())
			})
		})


	}, spec.Report(report.Terminal{}))
}
