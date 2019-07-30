/*
 * Copyright 2018-2019 the original author or authors.
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

package manifest

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestManifest(t *testing.T) {
	spec.Run(t, "Manifest", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.DetectFactory

		it.Before(func() {
			f = test.NewDetectFactory(t)
		})

		it("returns empty manifest if file doesn't exist", func() {
			m, err := NewManifest(f.Detect.Application, f.Detect.Logger)

			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(m.Len()).To(gomega.Equal(0))
		})

		it("returns populated manifest if file exists", func() {
			test.WriteFile(t, filepath.Join(f.Detect.Application.Root, "META-INF", "MANIFEST.MF"), "test-key=test-value")

			m, err := NewManifest(f.Detect.Application, f.Detect.Logger)

			g.Expect(err).NotTo(gomega.HaveOccurred())

			k, ok := m.Get("test-key")
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(k).To(gomega.Equal("test-value"))
		})

		it("returns proper values when lines are broken", func() {
			test.WriteFile(t, filepath.Join(f.Detect.Application.Root, "META-INF", "MANIFEST.MF"),
				`Manifest-Version: 1.0
Implementation-Title: petclinic
Implementation-Version: 2.1.0.BUILD-SNAPSHOT
Start-Class: org.springframework.samples.petclinic.PetClinicApplicatio
 n
Spring-Boot-Classes: BOOT-INF/classes/
Spring-Boot-Lib: BOOT-INF/lib/
Build-Jdk-Spec: 1.8
Spring-Boot-Version: 2.1.6.RELEASE
Created-By: Maven Archiver 3.4.0
Main-Class: org.springframework.boot.loader.JarLauncher
`)

			m, err := NewManifest(f.Detect.Application, f.Detect.Logger)

			g.Expect(err).NotTo(gomega.HaveOccurred())

			k, ok := m.Get("Start-Class")
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(k).To(gomega.Equal("org.springframework.samples.petclinic.PetClinicApplication"))

		})
	}, spec.Report(report.Terminal{}))
}
