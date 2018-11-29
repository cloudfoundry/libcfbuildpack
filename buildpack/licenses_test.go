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
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLicenses(t *testing.T) {
	spec.Run(t, "Licenses", testLicenses, spec.Report(report.Terminal{}))
}

func testLicenses(t *testing.T, when spec.G, it spec.S) {

	it("does not validate if there is not at least one license", func() {
		err := buildpack.Licenses{}.Validate()
		if err == nil {
			t.Errorf("Licenses.Validate() = nil expected error")
		}
	})

	it("validates when all licenses are valid", func() {
		err := buildpack.Licenses{
			{Type: "test-type"},
			{URI: "test-uri"},
		}.Validate()

		if err != nil {
			t.Errorf("Licenses.Validate() = %s expected no error", err)
		}
	})

	it("does not validate when a license is invalid", func() {
		err := buildpack.Licenses{
			{Type: "test-type"},
			{},
		}.Validate()

		if err == nil {
			t.Errorf("Licenses.Validate() = nil expected error")
		}
	})
}
