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

package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	buildpackBp "github.com/buildpack/libbuildpack/buildpack"
	layersBp "github.com/buildpack/libbuildpack/layers"
	loggerBp "github.com/buildpack/libbuildpack/logger"
	buildpackCf "github.com/cloudfoundry/libcfbuildpack/buildpack"
	layersCf "github.com/cloudfoundry/libcfbuildpack/layers"
	loggerCf "github.com/cloudfoundry/libcfbuildpack/logger"
)

type packager struct {
	Buildpack       buildpackCf.Buildpack
	Layers          layersCf.Layers
	Logger          loggerCf.Logger
	OutputDirectory string
}

func (p packager) Create() error {
	p.Logger.FirstLine("Packaging %s", p.Logger.PrettyIdentity(p.Buildpack))

	if err := p.prePackage(); err != nil {
		return err
	}

	includedFiles, err := p.Buildpack.IncludeFiles()
	if err != nil {
		return err
	}

	dependencyFiles, err := p.cacheDependencies()
	if err != nil {
		return err
	}

	return p.createPackage(append(includedFiles, dependencyFiles...))
}

func (p packager) cacheDependencies() ([]string, error) {
	var files []string

	deps, err := p.Buildpack.Dependencies()
	if err != nil {
		return nil, err
	}

	for _, dep := range deps {
		p.Logger.FirstLine("Caching %s", p.Logger.PrettyIdentity(dep))

		layer := p.Layers.DownloadLayer(dep)

		a, err := layer.Artifact()
		if err != nil {
			return nil, err
		}

		artifact, err := filepath.Rel(p.Buildpack.Root, a)
		if err != nil {
			return nil, err
		}

		metadata, err := filepath.Rel(p.Buildpack.Root, layer.Metadata)
		if err != nil {
			return nil, err
		}

		files = append(files, artifact, metadata)
	}

	return files, nil
}

func (p packager) createPackage(files []string) error {
	p.Logger.FirstLine("Creating package in %s", p.OutputDirectory)

	if err := os.MkdirAll(p.OutputDirectory, 0755); err != nil {
		return err
	}

	for _, file := range files {
		if err := p.addFile(file); err != nil {
			return err
		}
	}

	return nil
}

func (p packager) addFile(path string) error {
	p.Logger.SubsequentLine("Adding %s", path)

	in, err := os.OpenFile(filepath.Join(p.Buildpack.Root, path), os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer in.Close()

	stat, err := in.Stat()
	if err != nil {
		return err
	}

	f := filepath.Join(p.OutputDirectory, path)
	if err := os.MkdirAll(filepath.Dir(f), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, stat.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func (p packager) prePackage() error {
	pp, ok := p.Buildpack.PrePackage()
	if !ok {
		return nil
	}

	cmd := exec.Command(pp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = p.Buildpack.Root

	p.Logger.FirstLine("Pre-Package with %s", strings.Join(cmd.Args, " "))

	return cmd.Run()
}

func defaultPackager(outputDirectory string) (packager, error) {
	l := loggerBp.DefaultLogger()
	logger := loggerCf.Logger{Logger: l}

	b, err := buildpackBp.DefaultBuildpack(l)
	if err != nil {
		return packager{}, err
	}
	buildpack := buildpackCf.NewBuildpack(b)

	layers := layersCf.Layers{Layers: layersBp.Layers{Root: buildpack.CacheRoot}, Logger: logger}

	return packager{
		buildpack,
		layers,
		logger,
		outputDirectory,
	}, nil
}
