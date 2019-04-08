/*
 * Copyright 2018-2019 the original author or authors.
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

package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

const (
	pattern      = `(?m)(.*id[\s]+=[\s]+"%s"\n.*\nversion[\s]+=[\s]+")%s("\nuri[\s]+=[\s]+").*("\nsha256[\s]+=[\s]+").*(".*)`
	substitution = "${1}%s${2}%s${3}%s${4}"
)

type dependency struct {
	id             string
	versionPattern string
}

func (d dependency) update(version string, uri string, sha256 string) error {
	if err := d.validate(version, uri, sha256); err != nil {
		return err
	}

	r, err := regexp.Compile(fmt.Sprintf(pattern, d.id, d.versionPattern))
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile("buildpack.toml")
	if err != nil {
		return err
	}

	s := []byte(fmt.Sprintf(substitution, version, uri, sha256))

	b = r.ReplaceAll(b, s)

	return ioutil.WriteFile("buildpack.toml", b, 0644)
}

func (d dependency) validate(version string, uri string, sha256 string) error {
	if d.id == "" {
		return fmt.Errorf("id must be set")
	}

	if d.versionPattern == "" {
		return fmt.Errorf("version_pattern must be set")
	}

	if version == "" {
		return fmt.Errorf("version must be set")
	}

	if uri == "" {
		return fmt.Errorf("uri must be set")
	}

	if sha256 == "" {
		return fmt.Errorf("sha256 must be set")
	}

	return nil
}
