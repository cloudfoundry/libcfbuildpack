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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// BeFileLike tests that a file exists, has a specific mode, and specific content.
func BeFileLike(t *testing.T, file string, mode os.FileMode, content string) {
	t.Helper()

	FileExists(t, file)
	fileModeMatches(t, file, mode)
	fileContentMatches(t, file, content)
}

// BeLayerLike tests that a layer has a specific flag configuration.
func BeLayerLike(t *testing.T, layer layers.Layer, build bool, cache bool, launch bool) {
	t.Helper()

	m := make(map[string]interface{})
	if _, err := toml.DecodeFile(filepath.Join(layer.Metadata), &m); err != nil {
		t.Fatal(err)
	}

	b := m["build"].(bool)
	if b != build {
		t.Errorf("build flag = %t, expected %t", b, build)
	}

	c := m["cache"].(bool)
	if c != cache {
		t.Errorf("cache flag = %t, expected %t", c, cache)
	}

	l := m["launch"].(bool)
	if l != launch {
		t.Errorf("launch flag = %t, expected %t", l, launch)
	}
}

// FileExists tests that a file exists
func FileExists(t *testing.T, file string) {
	t.Helper()

	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("File %s does not exist", file)
		}

		t.Fatal(err)
	}
}

func fileModeMatches(t *testing.T, file string, mode os.FileMode) {
	t.Helper()

	fi, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}

	if fi.Mode() != mode {
		t.Errorf("FileMode = %#o, wanted %#o", fi.Mode(), mode)
	}
}

func fileContentMatches(t *testing.T, file string, content string) {
	t.Helper()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	actual := string(b)
	if actual != content {
		t.Errorf("File content = %s, wanted = %s", actual, content)
	}
}
