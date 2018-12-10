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

package buildpack_test

import (
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDependencies(t *testing.T) {
	spec.Run(t, "Dependencies", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

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

			g.Expect(d.Best("test-id-2", "1.0", "test-stack-1")).To(Equal(expected))
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

			g.Expect(d.Best("test-id", "2.0", "test-stack-1")).To(Equal(expected))
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

			g.Expect(d.Best("test-id", "1.0", "test-stack-3")).To(Equal(expected))
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

			g.Expect(d.Best("test-id", "1.*", "test-stack-1")).To(Equal(expected))
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
			}

			_, err := d.Best("test-id-2", "1.0", "test-stack-1")
			g.Expect(err).To(HaveOccurred())
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

			g.Expect(d.Best("test-id", "", "test-stack-1")).To(Equal(expected))
		})
	}, spec.Report(report.Terminal{}))
}
