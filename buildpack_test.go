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

package libjavabuildpack_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libjavabuildpack"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpack(t *testing.T) {
	spec.Run(t, "Buildpack", testBuildpack, spec.Report(report.Terminal{}))
}

func testBuildpack(t *testing.T, when spec.G, it spec.S) {

	it("returns error with no defined dependencies", func() {
		b := libbuildpack.Buildpack{}

		_, err := libjavabuildpack.Buildpack{b}.Dependencies()

		if err.Error() != "no dependencies specified" {
			t.Errorf("Buildpack.Dependencies = %s, expected no dependencies specified", err.Error())
		}
	})

	it("returns error with incorrectly defined dependencies", func() {
		b := libbuildpack.Buildpack{
			Metadata: map[string]interface{}{"dependencies": "test-dependency"},
		}

		_, err := libjavabuildpack.Buildpack{b}.Dependencies()

		if err.Error() != "dependencies have invalid structure" {
			t.Errorf("Buildpack.Dependencies = %s, expected dependencies have invalid structure", err.Error())
		}
	})

	it("returns dependencies", func() {
		b := libbuildpack.Buildpack{
			Metadata: map[string]interface{}{
				"dependencies": []map[string]interface{}{
					{
						"id":      "test-id-1",
						"name":    "test-name-1",
						"version": "1.0",
						"uri":     "test-uri-1",
						"sha256":  "test-sha256-1",
						"stacks":  []interface{}{"test-stack-1a", "test-stack-1b"},
					},
					{
						"id":      "test-id-2",
						"name":    "test-name-2",
						"version": "2.0",
						"uri":     "test-uri-2",
						"sha256":  "test-sha256-2",
						"stacks":  []interface{}{"test-stack-2a", "test-stack-2b"},
					},
				},
			},
		}

		expected := libjavabuildpack.Dependencies{
			libjavabuildpack.Dependency{
				ID:      "test-id-1",
				Name:    "test-name-1",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri-1",
				SHA256:  "test-sha256-1",
				Stacks:  []string{"test-stack-1a", "test-stack-1b"}},
			libjavabuildpack.Dependency{
				ID:      "test-id-2",
				Name:    "test-name-2",
				Version: newVersion(t, "2.0"),
				URI:     "test-uri-2",
				SHA256:  "test-sha256-2",
				Stacks:  []string{"test-stack-2a", "test-stack-2b"}},
		}

		actual, err := libjavabuildpack.Buildpack{b}.Dependencies()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Buildpack.Dependencies = %s, expected %s", actual, expected)
		}
	})

	it("filters by id", func() {
		d := libjavabuildpack.Dependencies{
			libjavabuildpack.Dependency{
				ID:      "test-id-1",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			libjavabuildpack.Dependency{
				ID:      "test-id-2",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
		}

		expected := libjavabuildpack.Dependency{
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
		d := libjavabuildpack.Dependencies{
			libjavabuildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			libjavabuildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "2.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
		}

		expected := libjavabuildpack.Dependency{
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
		d := libjavabuildpack.Dependencies{
			libjavabuildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			libjavabuildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-3"}},
		}

		expected := libjavabuildpack.Dependency{
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
		d := libjavabuildpack.Dependencies{
			libjavabuildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.1"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			libjavabuildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-3"}},
		}

		expected := libjavabuildpack.Dependency{
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
		d := libjavabuildpack.Dependencies{
			libjavabuildpack.Dependency{
				ID:      "test-id",
				Name:    "test-name",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri",
				SHA256:  "test-sha256",
				Stacks:  []string{"test-stack-1", "test-stack-2"}},
			libjavabuildpack.Dependency{
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
}

func newVersion(t *testing.T, version string) *semver.Version {
	t.Helper()

	v, err := semver.NewVersion(version)
	if err != nil {
		t.Fatal(err)
	}

	return v
}
