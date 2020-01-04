/*
 * Copyright 2018-2020 the original author or authors.
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
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/cloudfoundry/libcfbuildpack/helper"

	"github.com/cloudfoundry/libcfbuildpack/packager/cnbpackager"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	globalCacheDir := filepath.Join(usr.HomeDir, cnbpackager.DefaultCacheBase)
	localCacheDir := "."

	pflags := flag.NewFlagSet("Packager Flags", flag.ExitOnError)
	uncached := pflags.Bool("uncached", false, "cache dependencies")
	archive := pflags.Bool("archive", false, "tar resulting buildpack")
	summary := pflags.Bool("summary", false, "print buildpack.toml summary to stdout")
	globalCache := pflags.Bool("global_cache", false, fmt.Sprintf("use global cache dir at %s", globalCacheDir))
	version := pflags.String("version", "", "version to insert into buildpack.toml")

	if err := pflags.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to parse flags: %s\n", err)
		os.Exit(99)
	}

	pflags.Usage = func() {
		pflags.PrintDefaults()
		fmt.Println("-----")
		fmt.Println("Note that the destination should be the last argument")
		fmt.Println("The version will only be used in replacing a template version in buildpack.toml, if it is present")
		fmt.Println("If you provide a version, we will package it in a temporary directory")
		fmt.Println("Example: packager -archive -uncached -cachedir </path/to/cache> -version <version-to-insert> </path/to/destination> ")
	}

	destination := cnbpackager.DefaultDstDir

	if len(pflags.Args()) == 1 {
		destination = pflags.Args()[0]
	} else if len(pflags.Args()) > 1 {
		pflags.Usage()
		os.Exit(100)
	}

	current, err := filepath.Abs(".")
	bpDir := current
	if *version != "" {
		bpDir, err = ioutil.TempDir("", "cnb")
		if err != nil {
			os.Exit(101)
		}

		if err := helper.CopyDirectory(current, bpDir); err != nil {
			os.Exit(102)
		}

		if !filepath.IsAbs(destination) {
			destination = filepath.Join(current, destination)
		}

		if err := os.Chdir(bpDir); err != nil {
			os.Exit(103)
		}
	}

	var pkgr cnbpackager.Packager
	if *globalCache {
		pkgr, err = cnbpackager.New(bpDir, destination, *version, globalCacheDir)
	} else {
		pkgr, err = cnbpackager.New(bpDir, destination, *version, localCacheDir)
	}

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Packager: %s\n", err)
		os.Exit(104)
	}

	if *summary {
		summaryOutput, err := pkgr.Summary()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to produce summary: %s\n", err)
			os.Exit(105)
		}
		fmt.Println(summaryOutput)
		os.Exit(0)
	}

	if err := pkgr.Create(!*uncached); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create: %s\n", err)
		os.Exit(106)
	}

	if *archive {
		if err := pkgr.Archive(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to archive: %s\n", err)
			os.Exit(107)
		}
	}
}
