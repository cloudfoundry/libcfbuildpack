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

	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPriorityMerge(t *testing.T) {
	spec.Run(t, "PriorityMerge", testPriorityMerge, spec.Report(report.Terminal{}))
}

func testPriorityMerge(t *testing.T, when spec.G, it spec.S) {
	var (
		planA, planB, expected buildpackplan.Plan
		Expect                 func(interface{}, ...interface{}) Assertion
		priorities             = map[interface{}]int{
			"buildpack.yml": 3,
			"package.json":  2,
			".nvmrc":        1,
			"":              -1,
		}
	)
	var testPriorityMerge = func(planA, planB, expected buildpackplan.Plan) {
		output, err := buildpackplan.PriorityMerge(priorities)(planA, planB)
		Expect(err).NotTo(HaveOccurred())
		Expect(output).To(Equal(expected))
	}

	it.Before(func() {
		Expect = NewWithT(t).Expect
	})

	when("version key is empty for both", func() {
		it("leaves the version empty and merges metadata", func() {
			planA := createTestBuildPlan("", buildpackplan.Metadata{"key": "1"})
			planB := createTestBuildPlan("", buildpackplan.Metadata{"key": "2"})
			expected := createTestBuildPlan("", buildpackplan.Metadata{"key": "1,2"})

			testPriorityMerge(planA, planB, expected)
		})
	})

	when("version key is empty for one", func() {
		it("picks the other version and merges metadata", func() {
			planA := createTestBuildPlan("", buildpackplan.Metadata{"key": "1"})
			planB := createTestBuildPlan("1", buildpackplan.Metadata{"key": "2"})
			expected := createTestBuildPlan("1", buildpackplan.Metadata{"key": "1,2"})

			testPriorityMerge(planA, planB, expected)
		})
	})

	when("both have versions and VersionSource key is unset for both", func() {
		it("should pick latest version of the two", func() {
			planA := createTestBuildPlan("1.0", buildpackplan.Metadata{})
			planB := createTestBuildPlan("2.0", buildpackplan.Metadata{})
			expected := createTestBuildPlan("2.0", buildpackplan.Metadata{})
			testPriorityMerge(planA, planB, expected)
			testPriorityMerge(planB, planA, expected)

			planA = createTestBuildPlan("1.0", buildpackplan.Metadata{buildpackplan.VersionSource: ""})
			planB = createTestBuildPlan("2.0", buildpackplan.Metadata{buildpackplan.VersionSource: ""})
			expected = createTestBuildPlan("2.0", buildpackplan.Metadata{buildpackplan.VersionSource: ""})
			testPriorityMerge(planA, planB, expected)
			testPriorityMerge(planB, planA, expected)
		})

	})

	when("both have versions and VersionSource key is empty for one", func() {
		it.Before(func() {
			planA = createTestBuildPlan("2.0", buildpackplan.Metadata{buildpackplan.VersionSource: ""})
			planB = createTestBuildPlan("1.0", buildpackplan.Metadata{buildpackplan.VersionSource: "package.json"})
			expected = createTestBuildPlan("1.0", buildpackplan.Metadata{buildpackplan.VersionSource: "package.json"})
		})
		it("picks the version and source from the second input", func() {
			testPriorityMerge(planA, planB, expected)
		})

		it("picks the version and source from the first input", func() {
			testPriorityMerge(planB, planA, expected)
		})
	})

	when("different priority versions", func() {
		it.Before(func() {
			planA = createTestBuildPlan("2.0", buildpackplan.Metadata{buildpackplan.VersionSource: "package.json"})
			planB = createTestBuildPlan("1.0", buildpackplan.Metadata{buildpackplan.VersionSource: "buildpack.yml"})
			expected = createTestBuildPlan("1.0", buildpackplan.Metadata{buildpackplan.VersionSource: "buildpack.yml"})
		})

		it("picks the version and source from the second input", func() {
			testPriorityMerge(planA, planB, expected)
		})

		it("picks the version and source from the first input", func() {
			testPriorityMerge(planB, planA, expected)
		})
	})

	when("they are the same priority", func() {
		it("picks the higher version", func() {
			planA := createTestBuildPlan("2.0", buildpackplan.Metadata{buildpackplan.VersionSource: "buildpack.yml"})
			planB := createTestBuildPlan("1.0", buildpackplan.Metadata{buildpackplan.VersionSource: "buildpack.yml"})
			expected := createTestBuildPlan("2.0", buildpackplan.Metadata{buildpackplan.VersionSource: "buildpack.yml"})
			testPriorityMerge(planA, planB, expected)
			testPriorityMerge(planB, planA, expected)
		})
	})

	when("different build and launch metadata are set", func() {
		it("does an or-map of each key", func() {
			planA := createTestBuildPlan("", buildpackplan.Metadata{"build": true})
			planB := createTestBuildPlan("", buildpackplan.Metadata{"launch": "true", "build": false})
			expected := createTestBuildPlan("", buildpackplan.Metadata{"build": true, "launch": true})
			testPriorityMerge(planA, planB, expected)
		})
	})

}

func createTestBuildPlan(version string, metadata buildpackplan.Metadata) buildpackplan.Plan {
	return buildpackplan.Plan{
		Name:     "test",
		Version:  version,
		Metadata: metadata,
	}
}
