package main

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/internal"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPackager(t *testing.T) {
	spec.Run(t, "Package", testPackager, spec.Report(report.Terminal{}))
}

func testPackager(t *testing.T, _ spec.G, it spec.S) {
	var (
		pkgr        packager
		err         error
		outputDir   string
		cnbDir      string
		sha         [32]byte
		depUri      string
		depDir      string
		packageName string
	)

	it.Before(func() {
		RegisterTestingT(t)

		packageName = "temp_package"
		outputDir, err = ioutil.TempDir("", packageName)
		Expect(err).NotTo(HaveOccurred())

		depDir, err := ioutil.TempDir("", "test_dependency")
		Expect(err).NotTo(HaveOccurred())

		content := []byte("hello world")
		var buf bytes.Buffer

		tw := tar.NewWriter(&buf)
		defer tw.Close()

		hdr := &tar.Header{
			Name: "hello.txt",
			Mode: 0666,
			Size: int64(len(content)),
		}

		Expect(tw.WriteHeader(hdr)).To(Succeed())

		_, err = tw.Write(content)
		Expect(err).NotTo(HaveOccurred())

		depUri = filepath.Join(depDir, "hello.tgz")
		Expect(ioutil.WriteFile(depUri, buf.Bytes(), 0666)).To(Succeed())

		sha = sha256.Sum256(buf.Bytes())

		cnbDir, err = ioutil.TempDir("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.WriteFile(filepath.Join(cnbDir, "buildpack.toml"), []byte(`[buildpack]
id = "buildpack-id"
name = "buildpack-name"
version = "buildpack-version"

[metadata]
include_files = ["buildpack.toml"]

[[metadata.dependencies]]
id = "dependency-id"
name = "dependency-name"
sha256 = "`+fmt.Sprintf("%x", sha)+`"
stacks = ["stack-id"]
uri = "`+fmt.Sprintf("file://%s", depUri)+`"
version = "1.0.0"

[[stacks]]
id = 'stack-id'
`), 0666)).To(Succeed())

		defer internal.ReplaceArgs(t, filepath.Join(cnbDir, "buildpack.toml"))()
		defer internal.ReplaceWorkingDirectory(t, cnbDir)()

		pkgr, err = defaultPackager(outputDir)
		Expect(err).ToNot(HaveOccurred())
	})

	it.After(func() {
		os.RemoveAll(outputDir)
		os.RemoveAll(depDir)
		os.RemoveAll(cnbDir)
	})

	it("Create(true) makes a cached buildpack", func() {
		pkgr.Create(true)
		cacheRoot := filepath.Join(outputDir, "dependency-cache")

		Expect(cacheRoot).To(BeAnExistingFile())
		Expect(filepath.Join(cacheRoot, fmt.Sprintf("%x", sha))).To(BeAnExistingFile())
		Expect(filepath.Join(cacheRoot, fmt.Sprintf("%x.toml", sha))).To(BeAnExistingFile())
		Expect(filepath.Join(cacheRoot, fmt.Sprintf("%x", sha), "hello.tgz")).To(BeAnExistingFile())
		Expect(filepath.Join(outputDir, "buildpack.toml")).To(BeAnExistingFile())
	})

	it("Create(false) makes an uncached buildpack", func() {
		pkgr.Create(false)
		cacheRoot := filepath.Join(outputDir, "dependency-cache")
		Expect(cacheRoot).NotTo(BeAnExistingFile())
		Expect(filepath.Join(outputDir, "buildpack.toml")).To(BeAnExistingFile())
	})

	it("Archive can make a tarred up buildpack", func() {
		pkgr.Create(true)
		pkgr.Archive()

		tarball := filepath.Join(filepath.Dir(outputDir), filepath.Base(outputDir)+".tgz")
		Expect(tarball).To(BeAnExistingFile())
		os.RemoveAll(tarball)
	})
}
