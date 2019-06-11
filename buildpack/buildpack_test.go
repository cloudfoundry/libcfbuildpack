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

	bp "github.com/buildpack/libbuildpack/buildpack"
	"github.com/buildpack/libbuildpack/stack"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var TestDep1 = map[string]interface{}{
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

var TestDep2 = map[string]interface{}{
	"id":      "test-id-2",
	"name":    "test-name-2",
	"version": "2.0",
	"uri":     "test-uri-2",
	"sha256":  "test-sha256-2",
	"stacks":  []interface{}{"test-stack-2a", "test-stack-2b"},
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

func TestBuildpack(t *testing.T) {
	spec.Run(t, "Buildpack", func(t *testing.T, when spec.G, it spec.S) {
		g := NewWithT(t)

		when("Dependencies", func(){
			it("returns dependencies", func() {
				bp := bp.Buildpack{
					Metadata: bp.Metadata{
						buildpack.DependenciesMetadata: []map[string]interface{}{
							TestDep1, TestDep2,
						},
					},
				}
				expected := buildpack.Dependencies{
					buildpack.Dependency{
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
					},
					buildpack.Dependency{
						ID:      "test-id-2",
						Name:    "test-name-2",
						Version: internal.NewTestVersion(t, "2.0"),
						URI:     "test-uri-2",
						SHA256:  "test-sha256-2",
						Stacks:  buildpack.Stacks{"test-stack-2a", "test-stack-2b"},
						Licenses: buildpack.Licenses{
							buildpack.License{Type: "test-type-1", URI: "test-uri-1"},
							buildpack.License{Type: "test-type-2", URI: "test-uri-2"},
						},
					},
				}

				g.Expect(buildpack.Buildpack{Buildpack: bp}.Dependencies()).To(Equal(expected))
			})
		})

		when("IncludeFiles", func(){
			it("returns include_files if it exists", func() {
				b := bp.Buildpack{
					Metadata: bp.Metadata{
						"include_files": []interface{}{"test-file-1", "test-file-2"},
					},
				}

				g.Expect(buildpack.Buildpack{Buildpack: b}.IncludeFiles()).To(ConsistOf("test-file-1", "test-file-2"))
			})

			it("returns empty []string if include_files does not exist", func() {
				b := bp.Buildpack{}

				g.Expect(buildpack.Buildpack{Buildpack: b}.IncludeFiles()).To(BeEmpty())
			})
		})


		when("PrePackage", func() {
			it("returns pre_package if it exists", func() {
				b := bp.Buildpack{
					Metadata: bp.Metadata{
						"pre_package": "test-package",
					},
				}

				actual, ok := buildpack.Buildpack{Buildpack: b}.PrePackage()
				g.Expect(ok).To(BeTrue())
				g.Expect(actual).To(Equal("test-package"))
			})

			it("returns false if pre_package does not exist", func() {
				b := bp.Buildpack{}

				_, ok := buildpack.Buildpack{Buildpack: b}.PrePackage()
				g.Expect(ok).To(BeFalse())
			})
		})

		when("DefaultVersion", func(){
			it("returns a default dependency if it exists", func() {
				id := "test-id-1"
				version := "1.0"

				b := bp.Buildpack{
					Metadata: bp.Metadata{
						buildpack.DefaultVersions: map[string]interface{}{
							id: version,
						},
					},
				}

				ver, err := buildpack.Buildpack{Buildpack: b}.DefaultVersion(id)
				g.Expect(ver).To(Equal(version))
				g.Expect(err).ToNot(HaveOccurred())

				ver, err = buildpack.Buildpack{}.DefaultVersion("invalid-id")
				g.Expect(ver).To(Equal(""))
				g.Expect(err).ToNot(HaveOccurred())
			})

			it("returns empty string if DefaultVersions has incorrect structure", func() {
				id := "test-id-1"

				b := bp.Buildpack{
					Metadata: bp.Metadata{
						buildpack.DefaultVersions: "foo",
					},
				}

				ver, err := buildpack.Buildpack{Buildpack: b}.DefaultVersion(id)
				g.Expect(ver).To(Equal(""))
				g.Expect(err).ToNot(HaveOccurred())
			})

			it("returns an error if the type of values that DefaultVersions maps to are not strings", func() {
				id := "test-id-1"

				b := bp.Buildpack{
					Metadata: bp.Metadata{
						buildpack.DefaultVersions: map[string]interface{}{
							id: 1,
						},
					},
				}

				ver, err := buildpack.Buildpack{Buildpack: b}.DefaultVersion(id)
				g.Expect(ver).To(Equal(""))
				g.Expect(err).To(HaveOccurred())
			})
		})

		when("RuntimeDependency", func(){
			var (
				expectedDep buildpack.Dependency
				id = "test-id-1"
				stack = stack.Stack("test-stack-1a")
			)

			it.Before(func() {
				expectedDep = buildpack.Dependency{
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
			})

			it("get the Dependency matching the version parameter, if it exists", func() {
				b := bp.Buildpack{
					Metadata: bp.Metadata{
						buildpack.DependenciesMetadata: []map[string]interface{}{
							TestDep1,
						},
					},
				}

				dep, err := buildpack.Buildpack{Buildpack: b}.RuntimeDependency(id, "1.0", stack)
				g.Expect(dep).To(Equal(expectedDep))
				g.Expect(err).NotTo(HaveOccurred())
			})

			it(`get the Dependency matching the DefaultVersion, if "" or "default" version`, func() {
				b := bp.Buildpack{
					Metadata: bp.Metadata{
						buildpack.DefaultVersions: map[string]interface{}{
							id: "1.0",
						},
						buildpack.DependenciesMetadata: []map[string]interface{}{
							TestDep1,
						},
					},
				}

				dep, err := buildpack.Buildpack{Buildpack: b}.RuntimeDependency(id, "", stack)
				g.Expect(dep).To(Equal(expectedDep))
				g.Expect(err).NotTo(HaveOccurred())

				dep, err = buildpack.Buildpack{Buildpack: b}.RuntimeDependency(id, "default", stack)
				g.Expect(dep).To(Equal(expectedDep))
				g.Expect(err).NotTo(HaveOccurred())
			})

			it("returns an error if the Dependency is not present in the Buildpack", func() {
				b := bp.Buildpack{
					Metadata: bp.Metadata{
						buildpack.DefaultVersions: map[string]interface{}{
							id: 1,
						},
						buildpack.DependenciesMetadata: []map[string]interface{}{
							TestDep2,
						},
					},
				}

				dep, err := buildpack.Buildpack{Buildpack: b}.RuntimeDependency(id, "1.0", stack)
				g.Expect(dep).To(Equal(buildpack.Dependency{}))
				g.Expect(err).To(HaveOccurred())
			})
		})

	}, spec.Report(report.Terminal{}))
}
