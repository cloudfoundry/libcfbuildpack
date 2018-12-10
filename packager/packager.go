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
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	buildpackBp "github.com/buildpack/libbuildpack/buildpack"
	layersBp "github.com/buildpack/libbuildpack/layers"
	loggerBp "github.com/buildpack/libbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
)

type packager struct {
	buildpack       buildpack.Buildpack
	layers          layers.Layers
	logger          logger.Logger
	outputDirectory string
}

func (p packager) Create() error {
	p.logger.FirstLine("Packaging %s", p.logger.PrettyIdentity(p.buildpack))

	if err := p.prePackage(); err != nil {
		return err
	}

	includedFiles, err := p.buildpack.IncludeFiles()
	if err != nil {
		return err
	}

	dependencyFiles, err := p.cacheDependencies()
	if err != nil {
		return err
	}
	//
	return p.createPackage(append(includedFiles, dependencyFiles...))
}

func (p packager) cacheDependencies() ([]string, error) {
	var files []string

	deps, err := p.buildpack.Dependencies()
	if err != nil {
		return nil, err
	}

	for _, dep := range deps {
		p.logger.FirstLine("Caching %s", p.logger.PrettyIdentity(dep))

		layer := p.layers.DownloadLayer(dep)

		a, err := layer.Artifact()
		if err != nil {
			return nil, err
		}

		artifact, err := filepath.Rel(p.buildpack.Root, a)
		if err != nil {
			return nil, err
		}

		metadata, err := filepath.Rel(p.buildpack.Root, layer.Metadata)
		if err != nil {
			return nil, err
		}

		files = append(files, artifact, metadata)
	}

	return files, nil
}

func (p packager) createPackage(files []string) error {
	p.logger.FirstLine("Creating package in %s", p.outputDirectory)

	for _, file := range files {
		p.logger.SubsequentLine("Adding %s", file)
		if err := helper.CopyFile(filepath.Join(p.buildpack.Root, file), filepath.Join(p.outputDirectory, file)); err != nil {
			return err
		}
	}

	return nil
}

func (p packager) prePackage() error {
	pp, ok := p.buildpack.PrePackage()
	if !ok {
		return nil
	}

	cmd := exec.Command(pp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = p.buildpack.Root

	p.logger.FirstLine("Pre-Package with %s", strings.Join(cmd.Args, " "))

	return cmd.Run()
}

func defaultPackager(outputDirectory string) (packager, error) {
	l := loggerBp.DefaultLogger()
	logger := logger.Logger{Logger: l}

	b, err := buildpackBp.DefaultBuildpack(l)
	if err != nil {
		return packager{}, err
	}
	buildpack := buildpack.NewBuildpack(b, logger)

	layers := layers.NewLayers(layersBp.NewLayers(buildpack.CacheRoot, l), layersBp.NewLayers(buildpack.CacheRoot, l), logger)

	return packager{
		buildpack,
		layers,
		logger,
		outputDirectory,
	}, nil
}
