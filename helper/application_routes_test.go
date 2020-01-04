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

package helper_test

import (
	"os"
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/internal"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestApplicationRoutes(t *testing.T) {
	spec.Run(t, "ApplicationRoutes", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("extracts value from CNB_APP_ROUTES", func() {
			defer test.ReplaceEnv(t, "CNB_APP_ROUTES", `{
  "test-type-1": {
    "port": 1
  },
  "test-type-2": {
    "uri": "test-uri"
  }
}`)()

			g.Expect(helper.DefaultApplicationRoutes()).To(gomega.Equal(helper.ApplicationRoutes{
				"test-type-1": {Port: 1},
				"test-type-2": {URI: "test-uri"},
			}))
		})

		it("returns error when CNB_APP_ROUTES not set", func() {
			defer internal.ProtectEnv(t, "CNB_APP_ROUTES")()
			g.Expect(os.Unsetenv("CNB_APP_ROUTES")).Should(gomega.Succeed())

			_, err := helper.DefaultApplicationRoutes()
			g.Expect(err).To(gomega.MatchError("CNB_APP_ROUTES not set"))
		})
	}, spec.Report(report.Terminal{}))
}
