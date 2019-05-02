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

package runner

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
)

type Runner interface {
	Run(dir, program string, args ...string) error
	RunWithOutput(dir string, program string, args ...string) (string, error)
	CustomRun(dir string, addedEnv []string, stdout, stderr io.Writer, program string, args ...string) error
}

type CommandRunner struct{}

func (r *CommandRunner) Run(dir, program string, args ...string) error {
	return r.run("", nil, nil, nil, program, args...)
}

func (r *CommandRunner) RunWithOutput(dir string, program string, args ...string) (string, error) {
	logs := &bytes.Buffer{}

	if err := r.run(dir, nil, io.MultiWriter(os.Stdout, logs), io.MultiWriter(os.Stderr, logs), program, args...); err != nil {
		return "", err
	}

	return strings.TrimSpace(logs.String()), nil
}

func (r *CommandRunner) CustomRun(dir string, addedEnv []string, stdout, stderr io.Writer, program string, args ...string) error {
	return r.run(dir, addedEnv, stdout, stderr, program, args...)
}

func (r *CommandRunner) run(dir string, addedEnv []string, stdout, stderr io.Writer, program string, args ...string) error {
	cmd := exec.Command(program, args...)
	if stdout != nil {
		cmd.Stdout = stdout
	}
	if stderr != nil {
		cmd.Stderr = stderr
	}
	if dir != "" {
		cmd.Dir = dir
	}
	if addedEnv != nil {
		cmd.Env = append(os.Environ(), addedEnv...)
	}

	return cmd.Run()
}
