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

package logger_test

import (
	"bytes"
	"fmt"
	"testing"

	bp "github.com/buildpacks/libbuildpack/v2/logger"
	"github.com/cloudfoundry/libcfbuildpack/v2/logger"
	"github.com/heroku/color"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

var (
	body             = color.New(color.Faint).SprintfFunc()
	def              = color.New(color.Italic).SprintfFunc()
	description      = color.New(color.FgBlue).SprintfFunc()
	error            = color.New(color.FgRed, color.Bold).SprintfFunc()
	errorName        = color.New(color.FgRed, color.Bold).SprintfFunc()
	errorDescription = color.New(color.FgRed).SprintfFunc()
	name             = color.New(color.FgBlue, color.Bold).SprintfFunc()
	warning          = color.New(color.FgYellow, color.Bold).SprintfFunc()
)

func TestLogger(t *testing.T) {
	spec.Run(t, "Logger", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("writes title with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Title(metadata{"test-name", 1})

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("\n%s %s\n", name("test-name"), description("1"))))
		})

		it("writes header with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Header("test-header")

			g.Expect(info.String()).To(gomega.Equal("  test-header\n"))
		})

		it("writes header error with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.HeaderError("test-header")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("  %s\n", error("test-header"))))
		})

		it("writes header warning with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.HeaderWarning("test-header")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("  %s\n", warning("test-header"))))
		})

		it("writes body with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Body("test-body-1\ntest-body-2")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("%s\n", body("    test-body-1\n    test-body-2"))))
		})

		it("writes body error with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.BodyError("test-body-1\ntest-body-2")

			actual := info.String()
			expected := fmt.Sprintf("%s\n", body("    %s\x1b[%sm", error("test-body-1\n    test-body-2"), color.Faint.String()))
			g.Expect(actual).To(gomega.Equal(expected))
		})

		it("writes body with indent", func() {
			logger := logger.Logger{Logger: bp.NewLogger(nil, nil)}
			s := logger.BodyIndent("test-body-1\ntest-body-2")

			g.Expect(s).To(gomega.Equal("    test-body-1\n    test-body-2"))
		})

		it("writes launch configuration", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.LaunchConfiguration("test-message", "test-default")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("%s\n", body("    test-message. Default %s\x1b[%sm", def("test-default"), color.Faint.String()))))
		})

		it("writes body warning with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.BodyWarning("test-body-1\ntest-body-2")

			actual := info.String()
			expected := fmt.Sprintf("%s\n", body("    %s\x1b[%sm", warning("test-body-1\n    test-body-2"), color.Faint.String()))
			g.Expect(actual).To(gomega.Equal(expected))
		})

		it("writes terminal error with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.TerminalError(metadata{"test-name", 1}, "test-error")

			g.Expect(info.String()).
				To(gomega.Equal(fmt.Sprintf("\n%s %s\n  %s\n", errorName("test-name"), errorDescription("1"),
					error("test-error"))))
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
