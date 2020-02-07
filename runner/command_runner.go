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

package runner

import (
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/creack/pty"
	"github.com/mattn/go-isatty"
)

// CommandRunner is an empty struct to hang the Run method on.
type CommandRunner struct {
}

// Run makes CommandRunner satisfy the Runner interface.  This implementation delegates to exec.Command.
func (r CommandRunner) Run(bin string, dir string, args ...string) error {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir

	if !isatty.IsTerminal(os.Stdout.Fd()) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		return cmd.Run()
	}

	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(os.Stdout, f); err != nil {
		if !r.isEIO(err) {
			return err
		}
	}

	return cmd.Wait()
}

// RunWithOutput makes CommandRunner satisfy the Runner interface.  This implementation delegates to exec.Command.
func (r CommandRunner) RunWithOutput(bin string, dir string, args ...string) ([]byte, error) {
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

func (CommandRunner) isEIO(err error) bool {
	pe, ok := err.(*os.PathError)
	if !ok {
		return false
	}

	return pe.Err == syscall.EIO
}
