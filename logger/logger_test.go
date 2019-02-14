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

package logger_test

import (
	"bytes"
	"fmt"
	"testing"

	bp "github.com/buildpack/libbuildpack/logger"
	"github.com/cloudfoundry/libcfbuildpack/logger"
	"github.com/fatih/color"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLogger(t *testing.T) {
	spec.Run(t, "Logger", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("writes eye catcher on first line", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.FirstLine("test %s", "message")

			g.Expect(info.String()).To(Equal(fmt.Sprintf("%s test message\n", color.New(color.FgRed, color.Bold).Sprint("----->"))))
		})

		it("writes eye catcher on warning", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Warning("test %s", "message")

			g.Expect(info.String()).To(Equal(fmt.Sprintf("%s test message\n", color.New(color.FgYellow, color.Bold).Sprint("----->"))))
		})

		it("writes eye catcher on error", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.Error("test %s", "message")

			g.Expect(info.String()).To(Equal(fmt.Sprintf("%s test message\n", color.New(color.FgRed, color.Bold).Sprint("----->"))))
		})

		it("writes indent on second line", func() {
			var info bytes.Buffer

			logger := logger.Logger{Logger: bp.NewLogger(nil, &info)}
			logger.SubsequentLine("test %s", "message")

			g.Expect(info.String()).To(Equal("       test message\n"))
		})

		it("formats pretty identity", func() {
			logger := logger.Logger{Logger: bp.NewLogger(nil, nil)}

			g.Expect(logger.PrettyIdentity(metadata{"test-name", 1})).
				To(Equal(fmt.Sprintf("%s %s", color.New(color.FgBlue, color.Bold).Sprint("test-name"),
					color.BlueString("1"))))
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
