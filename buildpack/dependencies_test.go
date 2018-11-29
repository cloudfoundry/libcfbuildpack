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
	"reflect"
	"strings"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDependencies(t *testing.T) {
	spec.Run(t, "Dependencies", testDependencies, spec.Report(report.Terminal{}))
}

func testDependencies(t *testing.T, when spec.G, it spec.S) {

	it("filters by id", func() {
		d := buildpack.Dependencies{
			buildpack.Dependency{
				ID:      "test-id-1",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			buildpack.Dependency{
				ID:      "test-id-2",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
		}

		expected := buildpack.Dependency{
			ID:      "test-id-2",
			Name:    "test-name",
			Version: newVersion(t, "1.0"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  []string{"test-stack-1", "test-stack-2"}}

		actual, err := d.Best("test-id-2", "1.0", "test-stack-1")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Dependencies.Best = %s, expected %s", actual, expected)
		}
	})

	it("filters by version constraint", func() {
		d := buildpack.Dependencies{
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "2.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
		}

		expected := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "2.0"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  []string{"test-stack-1", "test-stack-2"}}

		actual, err := d.Best("test-id", "2.0", "test-stack-1")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Dependencies.Best = %s, expected %s", actual, expected)
		}
	})

	it("filters by stack", func() {
		d := buildpack.Dependencies{
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-3"}},
		}

		expected := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "1.0"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  []string{"test-stack-1", "test-stack-3"}}

		actual, err := d.Best("test-id", "1.0", "test-stack-3")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Dependencies.Best = %s, expected %s", actual, expected)
		}
	})

	it("returns the best dependency", func() {
		d := buildpack.Dependencies{
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.1"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-3"}},
		}

		expected := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "1.1"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  []string{"test-stack-1", "test-stack-2"}}

		actual, err := d.Best("test-id", "1.*", "test-stack-1")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Dependencies.Best = %s, expected %s", actual, expected)
		}
	})

	it("returns error if there are no matching dependencies", func() {
		d := buildpack.Dependencies{
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-3"}},
		}

		_, err := d.Best("test-id-2", "1.0", "test-stack-1")
		if !strings.HasPrefix(err.Error(), "no valid dependencies") {
			t.Errorf("Dependencies.Best = %s, expected no valid dependencies...", err.Error())
		}
	})

	it("substitutes all wildcard for unspecified version constraint", func() {
		d := buildpack.Dependencies{
			buildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.1"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
		}

		expected := buildpack.Dependency{
			ID:      "test-id",
			Name:    "test-name",
			Version: newVersion(t, "1.1"),
			URI:     "test-uri",
			SHA256:  "test-sha256",
			Stacks:  []string{"test-stack-1", "test-stack-2"}}

		actual, err := d.Best("test-id", "", "test-stack-1")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Dependencies.Best = %s, expected %s", actual, expected)
		}
	})
}
