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

package logger

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/buildpacks/libbuildpack/v2/logger"
	"github.com/heroku/color"
)

const (
	BodyIndent   = "    "
	HeaderIndent = "  "
)

var (
	body             = color.New(color.Faint).SprintfFunc()
	def              = color.New(color.Italic).SprintfFunc()
	description      = color.New(color.FgBlue).SprintfFunc()
	error            = color.New(color.FgRed, color.Bold).SprintfFunc()
	errorName        = color.New(color.FgRed, color.Bold).SprintfFunc()
	errorDescription = color.New(color.FgRed).SprintfFunc()
	lines            = regexp.MustCompile(`(?m)^`)
	name             = color.New(color.FgBlue, color.Bold).SprintfFunc()
	warning          = color.New(color.FgYellow, color.Bold).SprintfFunc()
)

func init() {
	color.Enabled()
}

// Logger is an extension to libbuildpack.Logger to add additional functionality.
type Logger struct {
	logger.Logger
}

// Title prints the buildpack description flush left, with an empty line above it.
func (l Logger) Title(v Identifiable) {
	if !l.IsInfoEnabled() {
		return
	}

	l.Info("\n%s", l.prettyIdentity(v, name, description))
}

// Terminal error prints the build description colored red and bold, flush left, with an empty line above it, followed
// by the log message message red and bold, and indented two spaces.
func (l Logger) TerminalError(v Identifiable, format string, args ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}

	l.Info("\n%s", l.prettyIdentity(v, errorName, errorDescription))
	l.HeaderError(format, args...)
}

// Header prints the log message indented two spaces, with an empty line above it.
func (l Logger) Header(format string, args ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}

	l.Info("%s%s", HeaderIndent, fmt.Sprintf(format, args...))
}

// HeaderError prints the log message colored red and bold, indented two spaces, with an empty line above it.
func (l Logger) HeaderError(format string, args ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}

	l.Header(error(format, args...))
}

// HeaderWarning prints the log message colored yellow and bold, indented two spaces, with an empty line above it.
func (l Logger) HeaderWarning(format string, args ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}

	l.Header(warning(format, args...))
}

// Body prints the log message with each line indented four spaces.
func (l Logger) Body(format string, args ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}

	l.Info(body(strings.ReplaceAll(
		l.BodyIndent(format, args...),
		fmt.Sprintf("\x1b[%sm", color.Reset.String()),
		fmt.Sprintf("\x1b[%sm\x1b[%sm", color.Reset.String(), color.Faint.String()))))
}

// BodyError prints the log message colored red and bold with each line indented four spaces.
func (l Logger) BodyError(format string, args ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}

	l.Body(error(format, args...))
}

// BodyIndent indents each line of a log message to the BodyIndent offset.
func (l Logger) BodyIndent(format string, args ...interface{}) string {
	return lines.ReplaceAllString(fmt.Sprintf(format, args...), BodyIndent)
}

// BodyWarning prints the log message colored yellow and bold with each line indented four spaces.
func (l Logger) BodyWarning(format string, args ...interface{}) {
	if !l.IsInfoEnabled() {
		return
	}

	l.Body(warning(format, args...))
}

func (l Logger) LaunchConfiguration(format string, defaultValue string) {
	l.Body("%s. Default %s", format, def(defaultValue))
}

func (l Logger) prettyIdentity(v Identifiable, nameFormatter func(format string, a ...interface{}) string, descriptionFormatter func(format string, a ...interface{}) string) string {
	if v == nil {
		return ""
	}

	name, description := v.Identity()

	if description == "" {
		return nameFormatter(name)
	}

	return fmt.Sprintf("%s %s", nameFormatter(name), descriptionFormatter(description))
}
