module github.com/cloudfoundry/libcfbuildpack

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Masterminds/semver v1.4.2
	github.com/buildpack/libbuildpack v1.11.0
	github.com/fatih/color v1.7.0
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.6 // indirect
	github.com/mitchellh/mapstructure v1.1.2
	github.com/onsi/gomega v1.4.3
	github.com/sclevine/spec v1.2.0
)

replace github.com/buildpack/libbuildpack => /users/pivotal/workspace/buildpack/libbuildpack
