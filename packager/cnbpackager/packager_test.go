package cnbpackager_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/packager/cnbpackager"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPackager(t *testing.T) {
	spec.Run(t, "Package", testPackager, spec.Report(report.Terminal{}))
}

func testPackager(t *testing.T, when spec.G, it spec.S) {
	var (
		outputDir, cnbDir, tempDir, depSHA, tarball string
		pkgr                                        cnbpackager.Packager
		err                                         error
	)

	it.Before(func() {
		RegisterTestingT(t)

		tempDir, err = ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		depFile, err := filepath.Abs(filepath.Join("testdata", "hello.tgz"))
		Expect(err).NotTo(HaveOccurred())
		depSHA = "9299d8c43e7af6797cd7fb5c12d986a90b864daa3c23ee50dab629ef844c1231"

		buildpackTOML := fmt.Sprintf(`
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
id = 'stack-id'
`, depSHA, depFile)

		outputDir = filepath.Join(tempDir, "output")
		cnbDir = filepath.Join(tempDir, "cnb")
		Expect(os.MkdirAll(cnbDir, 0777))
		Expect(ioutil.WriteFile(filepath.Join(cnbDir, "buildpack.toml"), []byte(buildpackTOML), 0666)).To(Succeed())
	})

	it.After(func() {
		os.RemoveAll(tempDir)
		if tarball != "" {
			os.RemoveAll(tarball)
		}
	})

	when("cached", func() {
		it.Before(func() {
			pkgr, err = cnbpackager.New(cnbDir, outputDir)
			Expect(err).ToNot(HaveOccurred())
		})

		it("Create makes a cached buildpack", func() {
			Expect(pkgr.Create(true)).To(Succeed())
			cacheRoot := filepath.Join(outputDir, "dependency-cache")

			Expect(cacheRoot).To(BeAnExistingFile())
			Expect(filepath.Join(cacheRoot, depSHA)).To(BeAnExistingFile())
			Expect(filepath.Join(cacheRoot, depSHA+".toml")).To(BeAnExistingFile())
			Expect(filepath.Join(cacheRoot, depSHA, "hello.tgz")).To(BeAnExistingFile())
			Expect(filepath.Join(outputDir, "buildpack.toml")).To(BeAnExistingFile())
		})

		it("Archive can make a tarred up cached buildpack", func() {
			Expect(pkgr.Create(true)).To(Succeed())
			Expect(pkgr.Archive(true)).To(Succeed())

			tarball = filepath.Join(filepath.Dir(outputDir), filepath.Base(outputDir)+"-cached.tgz")
			Expect(tarball).To(BeAnExistingFile())
			Expect(outputDir).NotTo(BeAnExistingFile())
		})
	})

	when("uncached", func() {
		it.Before(func() {
			pkgr, err = cnbpackager.New(cnbDir, outputDir)
			Expect(err).ToNot(HaveOccurred())
		})
		it("Create makes an uncached buildpack", func() {
			Expect(pkgr.Create(false)).To(Succeed())
			cacheRoot := filepath.Join(outputDir, "dependency-cache")
			Expect(cacheRoot).NotTo(BeAnExistingFile())
			Expect(filepath.Join(outputDir, "buildpack.toml")).To(BeAnExistingFile())
		})

		it("Archive can make a tarred up buildpack", func() {
			Expect(pkgr.Create(false)).To(Succeed())
			Expect(pkgr.Archive(false)).To(Succeed())

			tarball = filepath.Join(filepath.Dir(outputDir), filepath.Base(outputDir)+".tgz")
			Expect(tarball).To(BeAnExistingFile())
			Expect(outputDir).NotTo(BeAnExistingFile())
		})
	})

	when("summary", func() {

		it("Returns a package Summary of the CNB directory", func() {
			fakeCnbDir := filepath.Join("testdata", "summary-testdata", "fake-cnb")
			pkgr, err = cnbpackager.New(fakeCnbDir, "")
			Expect(err).ToNot(HaveOccurred())
			solution := `
Packaged binaries:

| name | version | cf_stacks |
|-|-|-|
| dep1 | 4.5.6 | stack1, stack2 |
| dep2 | 7.8.9 | stack2 |

Default binary versions:

| name | version |
|-|-|
| dep1 | 4.5.x |

Supported Stacks:

| name |
|-|
| stack1 |
| stack2 |
`

			summary, err := pkgr.Summary()
			Expect(err).ToNot(HaveOccurred())
			Expect(summary).To(Equal(solution))
		})

		it("does not have default versions", func() {
			fakeCnbDir := filepath.Join("testdata", "summary-testdata", "fake-cnb-without-defaults")
			pkgr, err = cnbpackager.New(fakeCnbDir, "")
			Expect(err).ToNot(HaveOccurred())
			solution := `
Packaged binaries:

| name | version | cf_stacks |
|-|-|-|
| dep1 | 4.5.6 | stack1 |
| dep2 | 7.8.9 | stack2 |

Supported Stacks:

| name |
|-|
| stack1 |
| stack2 |
`

			summary, err := pkgr.Summary()
			Expect(err).ToNot(HaveOccurred())
			Expect(summary).To(Equal(solution))
		})

		it("does not have any dependencies", func() {
			fakeCnbDir := filepath.Join("testdata", "summary-testdata", "fake-cnb-without-dependencies")
			pkgr, err = cnbpackager.New(fakeCnbDir, "")
			Expect(err).ToNot(HaveOccurred())
			solution := `
Supported Stacks:

| name |
|-|
| stack1 |
| stack2 |
`

			summary, err := pkgr.Summary()
			Expect(err).ToNot(HaveOccurred())
			Expect(summary).To(Equal(solution))
		})
	})
}
