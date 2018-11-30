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

package buildpack

import (
	"fmt"
	"path/filepath"

	"github.com/buildpack/libbuildpack/buildpack"
	"github.com/mitchellh/mapstructure"
)

// Buildpack is an extension to libbuildpack.Buildpack that adds additional opinionated behaviors.
type Buildpack struct {
	buildpack.Buildpack

	// CacheRoot is the path to the root directory for the buildpack's dependency cache.
	CacheRoot string
}

// Dependencies returns the collection of dependencies extracted from the generic buildpack metadata.
func (b Buildpack) Dependencies() (Dependencies, error) {
	d, ok := b.Metadata["dependencies"]
	if !ok {
		return Dependencies{}, nil
	}

	deps, ok := d.([]map[string]interface{})
	if !ok {
		return Dependencies{}, fmt.Errorf("dependencies have invalid structure")
	}

	var dependencies Dependencies
	for _, dep := range deps {
		d, err := b.dependency(dep)
		if err != nil {
			return Dependencies{}, err
		}

		dependencies = append(dependencies, d)
	}

	b.Logger.Debug("Dependencies: %s", dependencies)
	return dependencies, nil
}

// Identity make Buildpack satisfy the Identifiable interface.
func (b Buildpack) Identity() (string, string) {
	return b.Info.Name, b.Info.Version
}

// IncludeFiles returns the include_files buildpack metadata.
func (b Buildpack) IncludeFiles() ([]string, error) {
	i, ok := b.Metadata["include_files"]
	if !ok {
		return []string{}, nil
	}

	files, ok := i.([]interface{})
	if !ok {
		return []string{}, fmt.Errorf("include_files is not an array of strings")
	}

	var includes []string
	for _, candidate := range files {
		file, ok := candidate.(string)
		if !ok {
			return []string{}, fmt.Errorf("include_files is not an array of strings")
		}

		includes = append(includes, file)
	}

	return includes, nil
}

// PrePackage returns the pre_package buildpack metadata.
func (b Buildpack) PrePackage() (string, bool) {
	p, ok := b.Metadata["pre_package"]
	if !ok {
		return "", false
	}

	s, ok := p.(string)
	return s, ok
}

func (b Buildpack) dependency(dep map[string]interface{}) (Dependency, error) {
	var d Dependency

	config := mapstructure.DecoderConfig{
		DecodeHook: unmarshalText,
		Result:     &d,
	}

	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return Dependency{}, err
	}

	if err := decoder.Decode(dep); err != nil {
		return Dependency{}, err
	}

	return d, nil
}

// NewBuildpack creates a new instance of Buildpack from a specified buildpack.Buildpack.
func NewBuildpack(buildpack buildpack.Buildpack) Buildpack {
	return Buildpack{
		buildpack,
		filepath.Join(buildpack.Root, "cache"),
	}
}
