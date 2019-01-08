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

package layers_test

import (
	"fmt"
	"path/filepath"
	"testing"

	bp "github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayer(t *testing.T) {
	spec.Run(t, "Layer", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var (
			root  string
			layer layers.Layer
		)

		it.Before(func() {
			root = test.ScratchDir(t, "layer")
			layer = layers.NewLayers(bp.Layers{Root: root}, bp.Layers{}, logger.Logger{}).Layer("test-layer")
		})

		it("identifies matching metadata", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
Alpha = "test-value"
Bravo = 1
`)

			g.Expect(layer.MetadataMatches(metadata{"test-value", 1})).To(BeTrue())
		})

		it("identifies non-matching metadata", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
Alpha = "test-value"
Bravo = 2
`)

			g.Expect(layer.MetadataMatches(metadata{"test-value", 1})).To(BeFalse())
		})

		it("identifies invalid metadata", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
Alpha = "test-value"
Bravo = "invalid-value"
`)

			g.Expect(layer.MetadataMatches(metadata{"test-value", 1})).To(BeFalse())
		})

		it("identifies missing metadata", func() {
			g.Expect(layer.MetadataMatches(metadata{"test-value", 1})).To(BeFalse())
		})

		it("does not call contributor for cached layer", func() {
			test.WriteFile(t, layer.Metadata, `[metadata]
Alpha = "test-value"
Bravo = 1
`)

			contributed := false

			g.Expect(layer.Contribute(metadata{"test-value", 1}, func(layer layers.Layer) error {
				contributed = true
				return nil
			})).To(Succeed())

			g.Expect(contributed).To(BeFalse())
		})

		it("calls contributor for uncached layer", func() {
			contributed := false

			g.Expect(layer.Contribute(metadata{"test-value", 1}, func(layer layers.Layer) error {
				contributed = true
				return nil
			})).To(Succeed())

			g.Expect(contributed).To(BeTrue())
		})

		it("does not clean directory for non-matching metadata", func() {
			test.TouchFile(t, layer.Root, "test-file")

			g.Expect(layer.Contribute(metadata{"test-value", 1}, func(layer layers.Layer) error {
				return nil
			})).To(Succeed())

			g.Expect(filepath.Join(layer.Root, "test-file")).To(BeARegularFile())
		})
	}, spec.Report(report.Terminal{}))
}

type metadata struct {
	Alpha string
	Bravo int
}

func (m metadata) Identity() (string, string) {
	return m.Alpha, fmt.Sprintf("%d", m.Bravo)
}
