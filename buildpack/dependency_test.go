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

func TestDependency(t *testing.T) {
	spec.Run(t, "Dependency", testDependency, spec.Random(), spec.Report(report.Terminal{}))
}

func testDependency(t *testing.T, when spec.G, it spec.S) {

	it("validates", func() {
		err := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "1.0.0"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  buildpack.Stacks{"test-stack"},
			Licenses: buildpack.Licenses{
				{Type: "test-type"},
			},
		}.Validate()
		if err != nil {
			t.Errorf("Dependency.Validate() = %s expected no error", err)
		}
	})

	it("does not validate with invalid id", func() {
		err := buildpack.Dependency{
			Name:    "test-name",
			Version: newVersion(t, "1.0.0"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  buildpack.Stacks{"test-stack"},
			Licenses: buildpack.Licenses{
				{Type: "test-type"},
			},
		}.Validate()
		if err == nil {
			t.Errorf("Dependency.Validate() = nil expected error")
		}
	})

	it("does not validate with invalid name", func() {
		err := buildpack.Dependency{
			ID:      "test-id",
			Version: newVersion(t, "1.0.0"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  buildpack.Stacks{"test-stack"},
			Licenses: buildpack.Licenses{
				{Type: "test-type"},
			},
		}.Validate()
		if err == nil {
			t.Errorf("Dependency.Validate() = nil expected error")
		}
	})

	it("does not validate with invalid version", func() {
		err := buildpack.Dependency{
			ID:     "test-id",
			Name:   "test-name",
			URI:    "test-uri",
			SHA256: "test-sha256",
			Stacks: buildpack.Stacks{"test-stack"},
			Licenses: buildpack.Licenses{
				{Type: "test-type"},
			},
		}.Validate()
		if err == nil {
			t.Errorf("Dependency.Validate() = nil expected error")
		}
	})

	it("does not validate with invalid uri", func() {
		err := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "1.0.0"),
			SHA256:  "test-sha256",
			Stacks:  buildpack.Stacks{"test-stack"},
			Licenses: buildpack.Licenses{
				{Type: "test-type"},
			},
		}.Validate()
		if err == nil {
			t.Errorf("Dependency.Validate() = nil expected error")
		}
	})

	it("does not validate with invalid sha256", func() {
		err := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "1.0.0"),
			URI:     "test-uri",
			Stacks:  buildpack.Stacks{"test-stack"},
			Licenses: buildpack.Licenses{
				{Type: "test-type"},
			},
		}.Validate()
		if err == nil {
			t.Errorf("Dependency.Validate() = nil expected error")
		}
	})

	it("does not validate with invalid stacks", func() {
		err := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "1.0.0"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Licenses: buildpack.Licenses{
				{Type: "test-type"},
			},
		}.Validate()
		if err == nil {
			t.Errorf("Dependency.Validate() = nil expected error")
		}
	})

	it("does not validate with invalid licenses", func() {
		err := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "1.0.0"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  buildpack.Stacks{"test-stack"},
		}.Validate()
		if err == nil {
			t.Errorf("Dependency.Validate() = nil expected error")
		}
	})

}
