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

package layers_test

import (
	"bytes"
	"fmt"
	"testing"

	layersBp "github.com/buildpack/libbuildpack/layers"
	loggerBp "github.com/buildpack/libbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	layersCf "github.com/cloudfoundry/libcfbuildpack/layers"
	loggerCf "github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/fatih/color"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayers(t *testing.T) {
	spec.Run(t, "Layers", testLayers, spec.Random(), spec.Report(report.Terminal{}))
}

func testLayers(t *testing.T, when spec.G, it spec.S) {

	it("logs process types", func() {
		root := internal.ScratchDir(t, "launch")

		var info bytes.Buffer
		logger := loggerCf.Logger{Logger: loggerBp.NewLogger(nil, &info)}
		layers := layersCf.Layers{Layers: layersBp.Layers{Root: root}, Logger: logger}

		if err := layers.WriteMetadata(layersBp.Metadata{
			Processes: []layersBp.Process{
				{"short", "test-command-1"},
				{"a-very-long-type", "test-command-2"},
			},
		}); err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf(`%s Process types:
       %s:            test-command-1
       %s: test-command-2
`, color.New(color.FgRed, color.Bold).Sprint("----->"), color.CyanString("short"),
			color.CyanString("a-very-long-type"))

		if info.String() != expected {
			t.Errorf("Process types log = %s, expected %s", info.String(), expected)
		}
	})
}
