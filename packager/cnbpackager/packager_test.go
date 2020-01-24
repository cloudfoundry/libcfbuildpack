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

package cnbpackager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/v2/buildpack"
	"github.com/cloudfoundry/libcfbuildpack/v2/test"

	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPackager(t *testing.T) {
	spec.Run(t, "Package", testPackager, spec.Report(report.Terminal{}))
}

func testPackager(t *testing.T, when spec.G, it spec.S) {
	var (
		cnbDir, outputDir, cacheDir, tempDir, depSHA, tarball, newVersion, buildpackTOML string
		pkgr                                                                             Packager
		err                                                                              error
	)

	it.Before(func() {
		gomega.RegisterTestingT(t)

		tempDir, err = ioutil.TempDir("", "")
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		depFile, err := filepath.Abs(filepath.Join("testdata", "hello.tgz"))
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		depSHA = "9299d8c43e7af6797cd7fb5c12d986a90b864daa3c23ee50dab629ef844c1231"

		buildpackTOML = fmt.Sprintf(`
[buildpack]
id = "buildpack-id"
name = "buildpack-name"
version = "buildpack-version"

[metadata]
include_files = ["buildpack.toml"]

[[metadata.dependencies]]
id = "dependency-id"
name = "dependency-name"
sha256 = "%s"
stacks = ["stack-id"]
uri = "file://%s"
version = "1.0.0"

[[stacks]]
id = 'stack-id'`, depSHA, depFile)

		cnbDir = filepath.Join(tempDir, "cnb")
		gomega.Expect(os.MkdirAll(cnbDir, 0777))
		gomega.Expect(ioutil.WriteFile(filepath.Join(cnbDir, "buildpack.toml"), []byte(buildpackTOML), 0666)).To(gomega.Succeed())

		outputDir = filepath.Join(tempDir, "output")
		cacheDir = filepath.Join(tempDir, "cache")
		newVersion = "newVersion"
	})

	it.After(func() {
		os.RemoveAll(tempDir)
		if tarball != "" {
			os.RemoveAll(tarball)
		}
	})

	when("insertTemplateVersion", func() {
		when("buildpack.toml doesn't have a templated version", func() {
			it("doesn't replace anything in the buildpack", func() {
				buildpackTomlPath := filepath.Join(cnbDir, "buildpack.toml")
				gomega.Expect(insertTemplateVersion(cnbDir, newVersion)).To(gomega.Succeed())

				output, err := ioutil.ReadFile(buildpackTomlPath)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(output)).To(gomega.Equal(buildpackTOML))
			})
		})

		when("buildpack.toml has a templated version", func() {
			it("replaces the version in the buildpack", func() {
				buildpackTomlPath := filepath.Join(cnbDir, "buildpack.toml")
				baseBuildpackTOML := `[buildpack]
	id = "buildpack-id"
	name = "buildpack-name"
	version = "%s"`
				oldBPToml := fmt.Sprintf(baseBuildpackTOML, "{{ .Version }}")
				newBPToml := fmt.Sprintf(baseBuildpackTOML, newVersion)
				gomega.Expect(ioutil.WriteFile(buildpackTomlPath, []byte(oldBPToml), 0666)).To(gomega.Succeed())

				gomega.Expect(insertTemplateVersion(cnbDir, newVersion)).To(gomega.Succeed())

				output, err := ioutil.ReadFile(buildpackTomlPath)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(string(output)).To(gomega.Equal(newBPToml))
			})
		})
	})

	when("cached", func() {
		it.Before(func() {
			outputDir += "-cached"
			pkgr, err = New(cnbDir, outputDir, "", cacheDir)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		it("Create makes a cached buildpack", func() {
			gomega.Expect(pkgr.Create(true)).To(gomega.Succeed())
			cacheRoot := filepath.Join(cacheDir, buildpack.CacheRoot)

			gomega.Expect(filepath.Join(cacheRoot, depSHA+".toml")).To(gomega.BeAnExistingFile())
			gomega.Expect(filepath.Join(cacheRoot, depSHA, "hello.tgz")).To(gomega.BeAnExistingFile())
			gomega.Expect(filepath.Join(outputDir, "buildpack.toml")).To(gomega.BeAnExistingFile())
			gomega.Expect(filepath.Join(outputDir, buildpack.CacheRoot, depSHA+".toml")).To(gomega.BeAnExistingFile())
			gomega.Expect(filepath.Join(outputDir, buildpack.CacheRoot, depSHA, "hello.tgz")).To(gomega.BeAnExistingFile())
		})
	})

	when("uncached", func() {
		it.Before(func() {
			pkgr, err = New(cnbDir, outputDir, "", cacheDir)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		it("Create makes an uncached buildpack", func() {
			gomega.Expect(pkgr.Create(false)).To(gomega.Succeed())
			gomega.Expect(filepath.Join(outputDir, "buildpack.toml")).To(gomega.BeAnExistingFile())
			gomega.Expect(filepath.Join(outputDir, buildpack.CacheRoot)).NotTo(gomega.BeAnExistingFile())
		})
	})

	when("archiving", func() {
		it.Before(func() {
			fakeCnbDir := filepath.Join("testdata", "archive-testdata", "fake-cnb")
			pkgr, err = New(fakeCnbDir, outputDir, "", cacheDir)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
		})

		it("Archive can make a tarred up buildpack", func() {
			gomega.Expect(pkgr.Create(false)).To(gomega.Succeed())
			gomega.Expect(pkgr.Archive()).To(gomega.Succeed())

			tarball = filepath.Join(filepath.Dir(outputDir), filepath.Base(outputDir)+".tgz")
			gomega.Expect(tarball).To(gomega.BeAnExistingFile())
			gomega.Expect(outputDir).NotTo(gomega.BeAnExistingFile())
		})

		it("includes the parent directories of included files", func() {
			gomega.Expect(pkgr.Create(false)).To(gomega.Succeed())
			gomega.Expect(pkgr.Archive()).To(gomega.Succeed())

			tarball = filepath.Join(filepath.Dir(outputDir), filepath.Base(outputDir)+".tgz")
			gomega.Expect(tarball).To(gomega.BeAnExistingFile())
			gomega.Expect(outputDir).NotTo(gomega.BeAnExistingFile())

			gomega.Expect(tarball).NotTo(test.HaveArchiveEntry(""))
			gomega.Expect(tarball).To(test.HaveArchiveEntry("bin"))
			gomega.Expect(tarball).To(test.HaveArchiveEntry("bin/detect"))
			gomega.Expect(tarball).To(test.HaveArchiveEntry("bin/build"))
			gomega.Expect(tarball).To(test.HaveArchiveEntry("buildpack.toml"))
		})
	})

	when("summary", func() {
		it("Returns a package Summary of the CNB directory", func() {
			fakeCnbDir := filepath.Join("testdata", "summary-testdata", "fake-cnb")
			pkgr, err = New(fakeCnbDir, "", "", "")
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			solution := `
Packaged binaries:

| name | version | stacks |
|-|-|-|
| dep1 | 7.8.9 | stack1 |
| dep1 | 4.5.6 | stack1, stack2 |
| dep2 | 7.8.9 | stack2 |

Default binary versions:

| name | version |
|-|-|
| dep1 | 4.5.x |

Supported stacks:

| name |
|-|
| stack1 |
| stack2 |
`

			summary, err := pkgr.Summary()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(summary).To(gomega.Equal(solution))
		})

		it("does not have default versions", func() {
			fakeCnbDir := filepath.Join("testdata", "summary-testdata", "fake-cnb-without-defaults")
			pkgr, err = New(fakeCnbDir, "", "", "")
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			solution := `
Packaged binaries:

| name | version | stacks |
|-|-|-|
| dep1 | 4.5.6 | stack1 |
| dep2 | 7.8.9 | stack2 |

Supported stacks:

| name |
|-|
| stack1 |
| stack2 |
`

			summary, err := pkgr.Summary()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(summary).To(gomega.Equal(solution))
		})

		it("does not have any dependencies", func() {
			fakeCnbDir := filepath.Join("testdata", "summary-testdata", "fake-cnb-without-dependencies")
			pkgr, err = New(fakeCnbDir, "", "", "")
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			solution := `
Supported stacks:

| name |
|-|
| stack1 |
| stack2 |
`

			summary, err := pkgr.Summary()
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			gomega.Expect(summary).To(gomega.Equal(solution))
		})
	}, spec.Report(report.Terminal{}))
}
