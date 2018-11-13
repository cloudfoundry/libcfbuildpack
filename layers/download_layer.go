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

package layers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/fatih/color"
)

// DownloadLayer is an extension to Layer that is unique to a dependency download.
type DownloadLayer struct {
	Layer

	cacheLayer Layer
	dependency buildpack.Dependency
	logger     logger.Logger
}

// Artifact returns the path to an artifact cached in the layer.  If the artifact has already been downloaded, the cache
// will be validated and used directly.
func (l DownloadLayer) Artifact() (string, error) {
	artifact := filepath.Join(l.cacheLayer.Root, filepath.Base(l.dependency.URI))

	if err := l.cacheLayer.Contribute(l.dependency, func(layer Layer) error {
		return fmt.Errorf("buildpack cached dependency does not exist")
	}); err == nil {
		l.logger.SubsequentLine("%s cached download from buildpack", color.GreenString("Reusing"))
		return artifact, nil
	}

	artifact = filepath.Join(l.Layer.Root, filepath.Base(l.dependency.URI))

	if err := l.Layer.Contribute(l.dependency, func(layer Layer) error {
		l.logger.SubsequentLine("%s from %s", color.YellowString("Downloading"), l.dependency.URI)
		if err := l.download(artifact); err != nil {
			return err
		}

		l.logger.SubsequentLine("Verifying checksum")
		return l.verify(artifact)
	}, layers.Build, layers.Cache); err != nil {
		return "", err
	}

	return artifact, nil
}

// String makes DownloadLayer satisfy the Stringer interface.
func (l DownloadLayer) String() string {
	return fmt.Sprintf("DownloadLayer{ Layer: %s, cacheLayer:%s, dependency: %s, logger: %s }",
		l.Layer, l.cacheLayer, l.dependency, l.logger)
}

func (l DownloadLayer) download(file string) error {
	resp, err := http.Get(l.dependency.URI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("could not download: %bd", resp.StatusCode)
	}

	return WriteToFile(resp.Body, file, 0644)
}

func (l DownloadLayer) verify(file string) error {
	s := sha256.New()

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(s, f)
	if err != nil {
		return err
	}

	actualSha256 := hex.EncodeToString(s.Sum(nil))

	if actualSha256 != l.dependency.SHA256 {
		return fmt.Errorf("dependency sha256 mismatch: expected sha256 %s, actual sha256 %s",
			l.dependency.SHA256, actualSha256)
	}
	return nil
}
