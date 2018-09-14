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

package libjavabuildpack

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/buildpack/libbuildpack"
)

// Cache is an extension to libbuildpack.Cache that allows additional functionality to be added.
type Cache struct {
	libbuildpack.Cache
}

// DownloadLayer returns a DownloadLayer unique to a dependency.
func (c Cache) DownloadLayer(dependency Dependency) DownloadCacheLayer {
	l := c.Layer(dependency.SHA256)
	return DownloadCacheLayer{l, dependency}
}

// DownloadLayer is an extension to CacheLayer that is unique to a dependency download.
type DownloadCacheLayer struct {
	libbuildpack.CacheLayer

	dependency Dependency
}

// Artifact returns the path to an artifact cached in the layer.
func (d DownloadCacheLayer) Artifact() (string, error) {
	a := filepath.Join(d.Root, filepath.Base(d.dependency.URI))

	err := d.download(a)
	if err != nil {
		return "", err
	}

	err = d.verify(a)
	if err != nil {
		return "", err
	}

	return a, nil
}

func (d DownloadCacheLayer) download(file string) error {
	resp, err := http.Get(d.dependency.URI)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("could not download: %bd", resp.StatusCode)
	}

	return WriteToFile(resp.Body, file, 0644)
}

func (d DownloadCacheLayer) verify(file string) error {
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

	if actualSha256 != d.dependency.SHA256 {
		return fmt.Errorf("dependency sha256 mismatch: expected sha256 %s, actual sha256 %s",
			d.dependency.SHA256, actualSha256)
	}
	return nil
}
