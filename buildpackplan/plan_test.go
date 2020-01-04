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
	"fmt"
	"testing"

	bp "github.com/buildpack/libbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpackPlan(t *testing.T) {
	spec.Run(t, "Plan", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		ps := buildpackplan.Plans{
			Plans: bp.Plans{
				Entries: []bp.Plan{
					{Name: "test-entry-1", Metadata: bp.Metadata{"alpha": 1}},
					{Name: "test-entry-1", Metadata: bp.Metadata{"alpha": 2}},
					{Name: "test-entry-2", Metadata: bp.Metadata{"alpha": 1}},
				},
			},
		}

		it("gets filtered elements", func() {
			g.Expect(ps.Get("test-entry-1")).To(gomega.HaveLen(2))
		})

		when("GetMerged", func() {

			it("gets merged element", func() {
				p, ok, err := ps.GetMerged("test-entry-1", func(a, b buildpackplan.Plan) (buildpackplan.Plan, error) {
					a.Metadata["alpha"] = a.Metadata["alpha"].(int) + b.Metadata["alpha"].(int)
					return a, nil
				})

				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(ok).To(gomega.BeTrue())
				g.Expect(p).To(gomega.Equal(buildpackplan.Plan{
					Name:     "test-entry-1",
					Metadata: bp.Metadata{"alpha": 3},
				}))
			})

			it("returns false if no matches", func() {
				_, ok, err := ps.GetMerged("test-entry-3", func(a, b buildpackplan.Plan) (buildpackplan.Plan, error) {
					return a, nil
				})

				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(ok).To(gomega.BeFalse())
			})

			it("returns error if created", func() {
				_, _, err := ps.GetMerged("test-entry-1", func(a, b buildpackplan.Plan) (buildpackplan.Plan, error) {
					return buildpackplan.Plan{}, fmt.Errorf("")
				})

				g.Expect(err).To(gomega.HaveOccurred())
			})
		})

		when("has", func() {

			it("has plan", func() {
				g.Expect(ps.Has("test-entry-1")).To(gomega.BeTrue())
			})

			it("does not have plan", func() {
				g.Expect(ps.Has("test-entry-3")).To(gomega.BeFalse())
			})
		})
	}, spec.Report(report.Terminal{}))
}
