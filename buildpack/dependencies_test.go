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

package buildpack_test

import (
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/v2/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/v2/internal"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDependencies(t *testing.T) {
	spec.Run(t, "Dependencies", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("filters by id", func() {
			d := buildpack.Dependencies{
				buildpack.Dependency{
					ID:      "test-id-1",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
				buildpack.Dependency{
					ID:      "test-id-2",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
			}

			expected := buildpack.Dependency{
				ID:      "test-id-2",
				Name:    "test-name",
				Version: internal.NewTestVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}}

			g.Expect(d.Best("test-id-2", "1.0", "test-stack-1")).To(gomega.Equal(expected))
		})

		it("filters by version constraint", func() {
			d := buildpack.Dependencies{
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "2.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
			}

			expected := buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: internal.NewTestVersion(t, "2.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}}

			g.Expect(d.Best("test-id", "2.0", "test-stack-1")).To(gomega.Equal(expected))
		})

		it("filters by stack", func() {
			d := buildpack.Dependencies{
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-3"}},
			}

			expected := buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: internal.NewTestVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-3"}}

			g.Expect(d.Best("test-id", "1.0", "test-stack-3")).To(gomega.Equal(expected))
		})

		it("returns the best dependency", func() {
			d := buildpack.Dependencies{
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.1"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-3"}},
			}

			expected := buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: internal.NewTestVersion(t, "1.1"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}}

			g.Expect(d.Best("test-id", "1.*", "test-stack-1")).To(gomega.Equal(expected))
		})

		it("returns the best dependency after filtering", func() {
			d := buildpack.Dependencies{
				buildpack.Dependency{
					ID:      "test-id-1",
					Name:    "test-name-1",
					Version: internal.NewTestVersion(t, "1.9.1"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1"},
				},
				buildpack.Dependency{
					ID:      "test-id-1",
					Name:    "test-name-1",
					Version: internal.NewTestVersion(t, "1.9.1"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-2"},
				},
				buildpack.Dependency{
					ID:      "test-id-2",
					Name:    "test-name-2",
					Version: internal.NewTestVersion(t, "1.8.5"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-2"},
				},
				buildpack.Dependency{
					ID:      "test-id-2",
					Name:    "test-name-2",
					Version: internal.NewTestVersion(t, "1.8.6"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1"},
				},
				buildpack.Dependency{
					ID:      "test-id-2",
					Name:    "test-name-2",
					Version: internal.NewTestVersion(t, "1.8.6"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-2"},
				},
				buildpack.Dependency{
					ID:      "test-id-2",
					Name:    "test-name-2",
					Version: internal.NewTestVersion(t, "1.9.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1"},
				},
				buildpack.Dependency{
					ID:      "test-id-2",
					Name:    "test-name-2",
					Version: internal.NewTestVersion(t, "1.9.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-2"},
				},
			}

			expected := buildpack.Dependency{
				ID:      "test-id-2",
				Name:    "test-name-2",
				Version: internal.NewTestVersion(t, "1.9.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  buildpack.Stacks{"test-stack-2"},
			}

			g.Expect(d.Best("test-id-2", "", "test-stack-2")).To(gomega.Equal(expected))
		})

		it("returns error if there are no matching dependencies", func() {
			d := buildpack.Dependencies{
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.0"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-3"}},
				buildpack.Dependency{
					ID:      "test-id-2",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.1"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-3"}},
			}

			_, err := d.Best("test-id-2", "1.0", "test-stack-1")
			g.Expect(err).To(gomega.HaveOccurred())
			expectedError := "no valid dependencies for test-id-2, 1.0, and test-stack-1 in [(test-id, 1.0, [test-stack-1 test-stack-2]), (test-id, 1.0, [test-stack-1 test-stack-3]), (test-id-2, 1.1, [test-stack-1 test-stack-3])]"
			g.Expect(err).To(gomega.MatchError(expectedError))
		})

		it("substitutes all wildcard for unspecified version constraint", func() {
			d := buildpack.Dependencies{
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.1"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
			}

			expected := buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: internal.NewTestVersion(t, "1.1"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}}

			g.Expect(d.Best("test-id", "", "test-stack-1")).To(gomega.Equal(expected))
		})

		it("indicates that dependency exists", func() {
			d := buildpack.Dependencies{
				buildpack.Dependency{
					ID:      "test-id",
					Name:    "test-name",
					Version: internal.NewTestVersion(t, "1.1"),
					URI:     "test-uri",
					SHA256:  "test-sha256",
					Stacks:  buildpack.Stacks{"test-stack-1", "test-stack-2"}},
			}

			g.Expect(d.Has("test-id")).To(gomega.BeTrue())
		})

		it("indicates that dependency does not exist", func() {
			d := buildpack.Dependencies{}

			g.Expect(d.Has("test-id")).To(gomega.BeFalse())
		})
	}, spec.Report(report.Terminal{}))
}
