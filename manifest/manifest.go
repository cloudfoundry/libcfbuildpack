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

package manifest

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/magiconair/properties"
)

type Manifest struct {
	*properties.Properties
}

func NewManifest(application application.Application, logger logger.Logger) (Manifest, error) {
	f := filepath.Join(application.Root, "META-INF", "MANIFEST.MF")

	if exists, err := helper.FileExists(f); err != nil {
		return Manifest{}, err
	} else if !exists {
		return Manifest{properties.NewProperties()}, nil
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		return Manifest{}, err
	}

	p, err := properties.LoadString(normalizeManifest(string(b)))
	if err != nil {
		return Manifest{}, err
	}

	m := Manifest{p}

	logger.Debug("Manifest: %s", m)
	return m, nil
}

func normalizeManifest(manifest string) string {
	// The full grammar for manifests can be found here:
	// https://docs.oracle.com/javase/8/docs/technotes/guides/jar/jar.html#JARManifest

	// Convert Windows style line endings to UNIX
	n := strings.ReplaceAll(manifest, "\r\n", "\n")

	// The spec allows newlines to be single carriage-returns
	// this is a legacy line ending only supported on System 9
	// and before.
	n = strings.ReplaceAll(n, "\r", "\n")

	// The spec only allowed for line lengths of 78 bytes.
	// All lines are blank, start a property name or are
	// a continuation of the previous lines (indicated by a leading space).
	return strings.ReplaceAll(n, "\n ", "")
}
