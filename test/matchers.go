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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/buildpack/libbuildpack/buildplan"
	layersPkg "github.com/cloudfoundry/libcfbuildpack/layers"
)

// BeAppendBuildEnvLike tests that an append build env has specific content.
func BeAppendBuildEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beBuildEnvFileLike(t, layer, fmt.Sprintf("%s.append", name), format, args...)
}

// BeAppendLaunchEnvLike tests that an append launch env has specific content.
func BeAppendLaunchEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beLaunchEnvFileLike(t, layer, fmt.Sprintf("%s.append", name), format, args...)
}

// BeAppendSharedEnvLike tests that an append shared env has specific content.
func BeAppendSharedEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beSharedEnvFileLike(t, layer, fmt.Sprintf("%s.append", name), format, args...)
}

// BeAppendPathBuildEnvLike tests that an append path build env has specific content.
func BeAppendPathBuildEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beBuildEnvFileLike(t, layer, name, format, args...)
}

// BeAppendPathLaunchEnvLike tests that an append path launch env has specific content.
func BeAppendPathLaunchEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beLaunchEnvFileLike(t, layer, name, format, args...)
}

// BeAppendPathSharedEnvLike tests that an append path shared env has specific content.
func BeAppendPathSharedEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beSharedEnvFileLike(t, layer, name, format, args...)
}

// BeFileLike tests that a file exists, has a specific mode, and specific content.
func BeFileLike(t *testing.T, file string, mode os.FileMode, format string, args ...interface{}) {
	t.Helper()

	FileExists(t, file)
	fileModeMatches(t, file, mode)
	fileContentMatches(t, file, format, args...)
}

// BeLayerLike tests that a layer has a specific flag configuration.
func BeLayerLike(t *testing.T, layer layersPkg.Layer, build bool, cache bool, launch bool) {
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

// BeBuildPlanLike tests that a build plan, like the output of detect or build has specific content.
func BeBuildPlanLike(t *testing.T, actual buildplan.BuildPlan, expected buildplan.BuildPlan) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("BuildPlan = %s, expected %s", actual, expected)
	}
}

// BeLaunchMetadata tests that launch metadata has a specific configuration.
func BeLaunchMetadataLike(t *testing.T, layers layersPkg.Layers, expected layersPkg.Metadata) {
	t.Helper()

	var actual layersPkg.Metadata
	if _, err := toml.DecodeFile(filepath.Join(layers.Root, "launch.toml"), &actual); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("launch.toml = %s, expected %s", actual, expected)
	}
}

// BeOverrideBuildEnvLike tests that an override build env has specific content.
func BeOverrideBuildEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beBuildEnvFileLike(t, layer, fmt.Sprintf("%s.override", name), format, args...)
}

// BeOverrideLaunchEnvLike tests that an override launch env has specific content.
func BeOverrideLaunchEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beLaunchEnvFileLike(t, layer, fmt.Sprintf("%s.override", name), format, args...)
}

// BeOverrideSharedEnvLike tests that an override shared env has specific content.
func BeOverrideSharedEnvLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	beSharedEnvFileLike(t, layer, fmt.Sprintf("%s.override", name), format, args...)
}

// BeProfileLike tests that a profile.d file has specific content.
func BeProfileLike(t *testing.T, layer layersPkg.Layer, name string, format string, args ...interface{}) {
	t.Helper()
	BeFileLike(t, filepath.Join(layer.Root, "profile.d", name), 0644, format, args...)
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

func beBuildEnvFileLike(t *testing.T, layer layersPkg.Layer, file string, format string, args ...interface{}) {
	t.Helper()
	beEnvFileLike(t, layer, filepath.Join("env.build", file), format, args...)
}

func beEnvFileLike(t *testing.T, layer layersPkg.Layer, file string, format string, args ...interface{}) {
	t.Helper()
	BeFileLike(t, filepath.Join(layer.Root, file), 0644, format, args...)
}

func beLaunchEnvFileLike(t *testing.T, layer layersPkg.Layer, file string, format string, args ...interface{}) {
	t.Helper()
	beEnvFileLike(t, layer, filepath.Join("env.launch", file), format, args...)
}

func beSharedEnvFileLike(t *testing.T, layer layersPkg.Layer, file string, format string, args ...interface{}) {
	t.Helper()
	beEnvFileLike(t, layer, filepath.Join("env", file), format, args...)
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

func fileContentMatches(t *testing.T, file string, format string, args ...interface{}) {
	t.Helper()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	actual := string(b)
	expected := fmt.Sprintf(format, args...)

	if actual != expected {
		t.Errorf("File content = %s, wanted = %s", actual, expected)
	}
}
