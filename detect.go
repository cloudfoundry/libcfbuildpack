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
	"github.com/buildpack/libbuildpack"
)

type Detect struct {
	*libbuildpack.Detect

	// Buildpack represents the metadata associated with a buildpack.
	Buildpack Buildpack

	// Logger is used to write debug and info to the console.
	Logger Logger
}

// DefaultDetect creates a new instance of Detect using default values.
func DefaultDetect() (Detect, error) {
	d, err := libbuildpack.DefaultDetect()
	if err != nil {
		return Detect{}, err
	}

	return Detect{
		&d,
		Buildpack{d.Buildpack},
		Logger{d.Logger},
	}, nil
}
