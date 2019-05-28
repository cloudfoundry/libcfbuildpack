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

package manifest

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestManifest(t *testing.T) {
	spec.Run(t, "Manifest", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var f *test.DetectFactory

		it.Before(func() {
			f = test.NewDetectFactory(t)
		})

		it("returns empty manifest if file doesn't exist", func() {
			m, err := NewManifest(f.Detect.Application, f.Detect.Logger)

			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(m.Len()).To(Equal(0))
		})

		it("returns populated manifest if file exists", func() {
			test.WriteFile(t, filepath.Join(f.Detect.Application.Root, "META-INF", "MANIFEST.MF"), "test-key=test-value")

			m, err := NewManifest(f.Detect.Application, f.Detect.Logger)

			g.Expect(err).NotTo(HaveOccurred())

			k, ok := m.Get("test-key")
			g.Expect(ok).To(BeTrue())
			g.Expect(k).To(Equal("test-value"))
		})
	}, spec.Report(report.Terminal{}))
}
