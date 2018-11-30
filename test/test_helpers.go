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

package test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// FixturePath returns the absolute path to the desired fixture.
func FixturePath(t *testing.T, fixture string) string {
	t.Helper()
	return filepath.Join(findRoot(t), "fixtures", fixture)
}

// TouchFile touches a file with empty content
func TouchFile(t *testing.T, elem ...string) {
	t.Helper()
	if err := layers.WriteToFile(strings.NewReader(""), filepath.Join(elem...), 0644); err != nil {
		t.Fatal(err)
	}
}

func findRoot(t *testing.T) string {
	t.Helper()

	dir, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	for {
		if dir == "/" {
			t.Fatalf("could not find go.mod in the directory hierarchy")
		}
		if exist, err := layers.FileExists(filepath.Join(dir, "go.mod")); err != nil {
			t.Fatal(err)
		} else if exist {
			return dir
		}
		dir, err = filepath.Abs(filepath.Join(dir, ".."))
		if err != nil {
			t.Fatal(err)
		}
	}
}
