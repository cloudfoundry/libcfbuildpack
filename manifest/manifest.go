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

package manifest

import (
	"path/filepath"

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

	p, err := properties.LoadFile(f, properties.UTF8)
	if err != nil {
		return Manifest{}, err
	}

	m := Manifest{p}

	logger.Debug("Manifest: %s", m)
	return m, nil
}
