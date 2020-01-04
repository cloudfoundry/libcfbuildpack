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

package buildpack_test

import (
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/buildpack"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestStacks(t *testing.T) {
	spec.Run(t, "Stacks", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("does not validate if there is not at least one stack", func() {
			g.Expect(buildpack.Stacks{}.Validate()).NotTo(gomega.Succeed())
		})
	}, spec.Report(report.Terminal{}))
}
