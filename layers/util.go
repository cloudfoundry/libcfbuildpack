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
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// CopyDirectory copies srcDir to destDir
func CopyDirectory(srcDir, destDir string) error {
	destExists, _ := FileExists(destDir)
	if !destExists {
		return errors.New("destination dir must exist")
	}

	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		src := filepath.Join(srcDir, f.Name())
		dest := filepath.Join(destDir, f.Name())

		if m := f.Mode(); m&os.ModeSymlink != 0 {
			target, err := os.Readlink(src)
			if err != nil {
				return fmt.Errorf("Error while reading symlink '%s': %v", src, err)
			}
			if err := os.Symlink(target, dest); err != nil {
				return fmt.Errorf("Error while creating '%s' as symlink to '%s': %v", dest, target, err)
			}
		} else if f.IsDir() {
			err = os.MkdirAll(dest, f.Mode())
			if err != nil {
				return err
			}
			if err := CopyDirectory(src, dest); err != nil {
				return err
			}
		} else {
			rc, err := os.Open(src)
			if err != nil {
				return err
			}

			err = WriteToFile(rc, dest, f.Mode())
			if err != nil {
				rc.Close()
				return err
			}
			rc.Close()
		}
	}

	return nil
}

// CopyFile copies source file to destFile, creating all intermediate directories in destFile
func CopyFile(source, destFile string) error {
	fh, err := os.Open(source)
	if err != nil {
		return err
	}

	fileInfo, err := fh.Stat()
	if err != nil {
		return err
	}

	defer fh.Close()

	return WriteToFile(fh, destFile, fileInfo.Mode())
}

// ExtractTarGz extracts tarfile to destDir
func ExtractTarGz(tarFile, destDir string, stripComponents int) error {
	file, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer file.Close()
	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()
	return extractTar(gz, destDir, stripComponents)
}

// ExtractZip extracts zipfile to destDir
func ExtractZip(zipfile, destDir string, stripComponents int) error {
	r, err := zip.OpenReader(zipfile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		pathComponents := strings.Split(f.Name, string(filepath.Separator))
		if len(pathComponents) <= stripComponents {
			continue
		}

		path := filepath.Join(append([]string{destDir}, pathComponents[stripComponents:]...)...)

		rc, err := f.Open()
		if err != nil {
			return err
		}

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, 0755)
		} else {
			err = WriteToFile(rc, path, f.Mode())
		}

		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// FileExists returns true if a file exists, otherwise false.
func FileExists(file string) (bool, error) {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// FromTomlFile decodes a TOML file into a struct.
func FromTomlFile(file string, v interface{}) error {
	_, err := toml.DecodeFile(file, v)
	return err
}

// WriteToFile writes the contents of an io.Reader to a file.
func WriteToFile(source io.Reader, destFile string, mode os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(destFile), 0755)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(destFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = io.Copy(fh, source)
	if err != nil {
		return err
	}

	return nil
}

func extractTar(src io.Reader, destDir string, stripComponents int) error {
	tr := tar.NewReader(src)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}

		pathComponents := strings.Split(hdr.Name, string(filepath.Separator))
		if len(pathComponents) <= stripComponents {
			continue
		}

		path := filepath.Join(append([]string{destDir}, pathComponents[stripComponents:]...)...)
		fi := hdr.FileInfo()

		if fi.IsDir() {
			err = os.MkdirAll(path, hdr.FileInfo().Mode())
		} else if fi.Mode()&os.ModeSymlink != 0 {
			target := hdr.Linkname
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			if err = os.Symlink(target, path); err != nil {
				return err
			}
		} else {
			err = WriteToFile(tr, path, hdr.FileInfo().Mode())
		}

		if err != nil {
			return err
		}
	}
	return nil
}
