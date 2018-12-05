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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/fatih/color"
)

// TouchedLayers contains information about the layers that have been touched as part of this execution.
type TouchedLayers struct {
	// Logger logger is used to write debug and info to the console.
	Logger logger.Logger

	// Root is the root location of all layers to inspect for unused layers.
	Root string

	// Touched is the set of layers that have been touched
	Touched map[string]struct{}
}

// Add registers that a given layer has been touched
func (t TouchedLayers) Add(metadata string) {
	t.Logger.Debug("Layer %s touched", metadata)
	t.Touched[metadata] = struct{}{}
}

// Cleanup removes all layers that have not been touched as part of this execution.
func (t TouchedLayers) Cleanup() error {
	candidates, err := filepath.Glob(filepath.Join(t.Root, "*.toml"))
	if err != nil {
		return err
	}

	if t.Logger.IsDebugEnabled() {
		t.Logger.Debug("Existing Layers: %s", candidates)
		t.Logger.Debug("Touched Layers: %s", t.keys(t.Touched))
	}

	remove := t.union(candidates, t.Touched)
	if len(remove) == 0 {
		return nil
	}

	t.Logger.FirstLine("%s unused layers", color.YellowString("Removing"))
	for _, r := range remove {
		t.Logger.SubsequentLine(strings.TrimSuffix(filepath.Base(r), ".toml"))

		if err := os.RemoveAll(r); err != nil {
			return err
		}
	}

	return nil
}

// String makes TouchedLayers satisfy the Stringer interface.
func (t TouchedLayers) String() string {
	return fmt.Sprintf("TouchedLayers{ Logger: %s, Root: %s, Touched: %s }",
		t.Logger, t.Root, t.Touched)
}

func (t TouchedLayers) keys(m map[string]struct{}) []string {
	keys := make([]string, len(t.Touched))
	for k := range t.Touched {
		keys = append(keys, k)
	}

	return keys
}

func (t TouchedLayers) union(a []string, b map[string]struct{}) []string {
	union := make([]string, 0)
	for _, c := range a {
		if _, ok := b[c]; !ok {
			union = append(union, c)
		}
	}

	return union
}

func NewTouchedLayers(root string, logger logger.Logger) TouchedLayers {
	return TouchedLayers{
		logger,
		root,
		make(map[string]struct{}),
	}
}
