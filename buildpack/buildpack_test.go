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
	"testing"

	buildpackBp "github.com/buildpack/libbuildpack/buildpack"
	buildpackCf "github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpack(t *testing.T) {
	spec.Run(t, "Buildpack", testBuildpack, spec.Report(report.Terminal{}))
}

func testBuildpack(t *testing.T, when spec.G, it spec.S) {

	it("returns error with incorrectly defined dependencies", func() {
		b := buildpackBp.Buildpack{
			Metadata: buildpackBp.Metadata{"dependencies": "test-dependency"},
		}

		_, err := buildpackCf.Buildpack{Buildpack: b}.Dependencies()

		if err.Error() != "dependencies have invalid structure" {
			t.Errorf("Buildpack.Dependencies = %s, expected dependencies have invalid structure", err.Error())
		}
	})

	it("returns dependencies", func() {
		b := buildpackBp.Buildpack{
			Metadata: buildpackBp.Metadata{
				"dependencies": []map[string]interface{}{
					{
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
					},
					{
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
					},
				},
			},
		}

		expected := buildpackCf.Dependencies{
			buildpackCf.Dependency{
				ID:      "test-id-1",
				Name:    "test-name-1",
				Version: newVersion(t, "1.0"),
				URI:     "test-uri-1",
				SHA256:  "test-sha256-1",
				Stacks:  buildpackCf.Stacks{"test-stack-1a", "test-stack-1b"},
				Licenses: buildpackCf.Licenses{
					buildpackCf.License{Type: "test-type-1", URI: "test-uri-1"},
					buildpackCf.License{Type: "test-type-2", URI: "test-uri-2"},
				},
			},
			buildpackCf.Dependency{
				ID:      "test-id-2",
				Name:    "test-name-2",
				Version: newVersion(t, "2.0"),
				URI:     "test-uri-2",
				SHA256:  "test-sha256-2",
				Stacks:  buildpackCf.Stacks{"test-stack-2a", "test-stack-2b"},
				Licenses: buildpackCf.Licenses{
					buildpackCf.License{Type: "test-type-1", URI: "test-uri-1"},
					buildpackCf.License{Type: "test-type-2", URI: "test-uri-2"},
				},
			},
		}

		actual, err := buildpackCf.Buildpack{Buildpack: b}.Dependencies()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Buildpack.Dependencies = %s, expected %s", actual, expected)
		}
	})

	it("returns include_files if it exists", func() {
		b := buildpackBp.Buildpack{
			Metadata: buildpackBp.Metadata{
				"include_files": []interface{}{"test-file-1", "test-file-2"},
			},
		}

		actual, err := buildpackCf.Buildpack{Buildpack: b}.IncludeFiles()
		if err != nil {
			t.Fatal(err)
		}

		expected := []string{"test-file-1", "test-file-2"}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Buildpack.IncludeFiles = %s, expected empty []string", actual)
		}
	})

	it("returns empty []string if include_files does not exist", func() {
		b := buildpackBp.Buildpack{}

		actual, err := buildpackCf.Buildpack{Buildpack: b}.IncludeFiles()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(actual, []string{}) {
			t.Errorf("Buildpack.IncludeFiles = %s, expected empty []string", actual)
		}
	})

	it("returns false if include_files is not []string", func() {
		b := buildpackBp.Buildpack{
			Metadata: buildpackBp.Metadata{
				"include_files": 1,
			},
		}

		_, err := buildpackCf.Buildpack{Buildpack: b}.IncludeFiles()
		if err.Error() != "include_files is not an array of strings" {
			t.Errorf("Buildpack.IncludeFiles = %s, expected include_files is not an array of strings", err.Error())
		}
	})

	it("returns pre_package if it exists", func() {
		b := buildpackBp.Buildpack{
			Metadata: buildpackBp.Metadata{
				"pre_package": "test-package",
			},
		}

		actual, ok := buildpackCf.Buildpack{Buildpack: b}.PrePackage()
		if !ok {
			t.Errorf("Buildpack.PrePackage() = %t, expected true", ok)
		}

		if actual != "test-package" {
			t.Errorf("Buildpack.PrePackage() %s, expected test-package", actual)
		}
	})

	it("returns false if pre_package does not exist", func() {
		b := buildpackBp.Buildpack{}

		_, ok := buildpackCf.Buildpack{Buildpack: b}.PrePackage()
		if ok {
			t.Errorf("Buildpack.PrePackage() = %t, expected false", ok)
		}
	})

	it("returns false if pre_package is not string", func() {
		b := buildpackBp.Buildpack{
			Metadata: buildpackBp.Metadata{
				"pre_package": 1,
			},
		}

		_, ok := buildpackCf.Buildpack{Buildpack: b}.PrePackage()
		if ok {
			t.Errorf("Buildpack.PrePackage() = %t, expected false", ok)
		}
	})

}
