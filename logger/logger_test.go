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

package logger_test

import (
	"bytes"
	"fmt"
	"testing"

	bp "github.com/buildpack/libbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/fatih/color"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLogger(t *testing.T) {
	spec.Run(t, "Logger", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("writes title with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Title(metadata{"test-name", 1})

			g.Expect(info.String()).
				To(gomega.Equal(fmt.Sprintf("\n%s %s\n", color.New(color.FgBlue, color.Bold).Sprint("test-name"),
					color.New(color.FgBlue).Sprint("1"))))
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

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("  %s\n", color.New(color.FgRed, color.Bold).Sprint("test-header"))))
		})

		it("writes header warning with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.HeaderWarning("test-header")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("  %s\n", color.New(color.FgYellow, color.Bold).Sprint("test-header"))))
		})

		it("writes body with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Body("test-body-1\ntest-body-2")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("%s\n", color.New(color.Faint).Sprint("    test-body-1\n    test-body-2"))))
		})

		it("writes body error with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.BodyError("test-body-1\ntest-body-2")

			actual := info.String()
			expected := fmt.Sprintf("%s\n", color.New(color.Faint).Sprintf("    %s\x1b[%dm", color.New(color.FgRed, color.Bold).Sprint("test-body-1\n    test-body-2"), color.Faint))
			g.Expect(actual).To(gomega.Equal(expected))
		})

		it("writes body with indent", func() {
			logger := logger.Logger{Logger: bp.NewLogger(nil, nil)}
			s := logger.BodyIndent("test-body-1\ntest-body-2")

			g.Expect(s).To(gomega.Equal("    test-body-1\n    test-body-2"))
		})

		it("writes body warning with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.BodyWarning("test-body-1\ntest-body-2")

			actual := info.String()
			expected := fmt.Sprintf("%s\n", color.New(color.Faint).Sprintf("    %s\x1b[%dm", color.New(color.FgYellow, color.Bold).Sprint("test-body-1\n    test-body-2"), color.Faint))
			g.Expect(actual).To(gomega.Equal(expected))
		})

		it("writes terminal error with format", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.TerminalError(metadata{"test-name", 1}, "test-error")

			g.Expect(info.String()).
				To(gomega.Equal(fmt.Sprintf("\n%s %s\n  %s\n", color.New(color.FgRed, color.Bold).Sprint("test-name"),
					color.New(color.FgRed).Sprint("1"), color.New(color.FgRed, color.Bold).Sprint("test-error"))))
		})

		it("writes eye catcher on first line", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.FirstLine("test %s", "message")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("%s test message\n", color.New(color.FgRed, color.Bold).Sprint("----->"))))
		})

		it("writes eye catcher on warning", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Warning("test %s", "message")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("%s test message\n", color.New(color.FgYellow, color.Bold).Sprint("----->"))))
		})

		it("writes eye catcher on error", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Error("test %s", "message")

			g.Expect(info.String()).To(gomega.Equal(fmt.Sprintf("%s test message\n", color.New(color.FgRed, color.Bold).Sprint("----->"))))
		})

		it("writes indent on second line", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.SubsequentLine("test %s", "message")

			g.Expect(info.String()).To(gomega.Equal("       test message\n"))
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
