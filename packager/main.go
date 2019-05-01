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
	"flag"
	"fmt"
	"github.com/cloudfoundry/libcfbuildpack/packager/cnbpackager"
	"os"
)

func main() {
	pflags := flag.NewFlagSet("Packager Flags", flag.ExitOnError)
	uncached := pflags.Bool("uncached", false, "cache dependencies")
	archive := pflags.Bool("archive", false, "tar resulting buildpack")

	if err := pflags.Parse(os.Args[1:]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to parse flags: %s\n", err)
		os.Exit(99)
	}

	pflags.Usage = func() {
		pflags.PrintDefaults()
		fmt.Println("-----")
		fmt.Println("Note that the destination should be the last argument")
		fmt.Println("Example: packager --archive --uncached </path/to/destination>")
	}

	if len(pflags.Args()) != 1 {
		pflags.Usage()
		os.Exit(100)
	}

	destination := pflags.Args()[0]
	defaultPackager, err := cnbpackager.DefaultPackager(destination)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Packager: %s\n", err)
		os.Exit(101)
	}

	if err := defaultPackager.Create(!*uncached); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create: %s\n", err)
		os.Exit(102)
	}

	if *archive {
		if err := defaultPackager.Archive(!*uncached); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to archive: %s\n", err)
			os.Exit(103)
		}
	}
}
