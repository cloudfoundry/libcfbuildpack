/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
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
	dependencyPattern      = `(?m)(.*id[\s]+=[\s]+"%s"\n.*\nversion[\s]+=[\s]+")%s("\nuri[\s]+=[\s]+").*("\nsha256[\s]+=[\s]+").*(".*)`
	dependencySubstitution = "${1}%s${2}%s${3}%s${4}"
	orderPattern           = `([\s]+{[\s]+id[\s]+=[\s]+"%s",[\s]+version[\s]+=[\s]+")%s(".+)`
	orderSubstitution      = "${1}%s${2}"
)

type dependency struct {
	id             string
	versionPattern string
}

func (d dependency) update(version string, uri string, sha256 string) error {
	if err := d.validate(version, uri, sha256); err != nil {
		return err
	}

	b, err := ioutil.ReadFile("buildpack.toml")
	if err != nil {
		return err
	}

	b, err = d.updateDependency(version, uri, sha256, b)
	if err != nil {
		return err
	}

	b, err = d.updateOrder(version, b)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("buildpack.toml", b, 0644)
}

func (d dependency) updateDependency(version string, uri string, sha256 string, b []byte) ([]byte, error) {
	r, err := regexp.Compile(fmt.Sprintf(dependencyPattern, d.id, d.versionPattern))
	if err != nil {
		return nil, err
	}

	s := []byte(fmt.Sprintf(dependencySubstitution, version, uri, sha256))

	return r.ReplaceAll(b, s), nil
}

func (d dependency) updateOrder(version string, b []byte) ([]byte, error) {
	r, err := regexp.Compile(fmt.Sprintf(orderPattern, d.id, d.versionPattern))
	if err != nil {
		return nil, err
	}

	s := []byte(fmt.Sprintf(orderSubstitution, version))

	return r.ReplaceAll(b, s), nil
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
