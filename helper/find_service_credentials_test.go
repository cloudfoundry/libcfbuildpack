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

package helper_test

import (
	"testing"

	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestFindServiceCredentials(t *testing.T) {
	spec.Run(t, "FindServiceCredentials", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("matches single service by BindingName", func() {
			defer test.ReplaceEnv(t, "CNB_SERVICES", `{
  "": [
    {
      "binding_name": "test-service",
      "credentials": {
        "test-key": "test-value"
      },
      "instance_name": null,
      "label": null,
      "tags": [
      ]
    }
  ]
}`)()

			c, ok, err := helper.FindServiceCredentials("test-service")
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(c).To(gomega.Equal(helper.Credentials{"test-key": "test-value"}))
		})

		it("matches single service by InstanceName", func() {
			defer test.ReplaceEnv(t, "CNB_SERVICES", `{
  "": [
    {
      "binding_name": null,
      "credentials": {
        "test-key": "test-value"
      },
      "instance_name": "test-service",
      "label": null,
      "tags": [
      ]
    }
  ]
}`)()

			c, ok, err := helper.FindServiceCredentials("test-service")
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(c).To(gomega.Equal(helper.Credentials{"test-key": "test-value"}))
		})

		it("matches single service by Label", func() {
			defer test.ReplaceEnv(t, "CNB_SERVICES", `{
  "": [
    {
      "binding_name": null,
      "credentials": {
        "test-key": "test-value"
      },
      "instance_name": null,
      "label": "test-service",
      "tags": [
      ]
    }
  ]
}`)()

			c, ok, err := helper.FindServiceCredentials("test-service")
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(c).To(gomega.Equal(helper.Credentials{"test-key": "test-value"}))
		})

		it("matches single service by Tags", func() {
			defer test.ReplaceEnv(t, "CNB_SERVICES", `{
  "": [
    {
      "binding_name": null,
      "credentials": {
        "test-key": "test-value"
      },
      "instance_name": null,
      "label": null,
      "tags": [
        "test-service"
      ]
    }
  ]
}`)()

			c, ok, err := helper.FindServiceCredentials("test-service")
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(c).To(gomega.Equal(helper.Credentials{"test-key": "test-value"}))
		})

		it("matches single service with Credentials", func() {
			defer test.ReplaceEnv(t, "CNB_SERVICES", `{
  "": [
    {
      "binding_name": null,
      "credentials": {
        "test-key": "test-value"
      },
      "instance_name": null,
      "label": null,
      "tags": [
        "test-service"
      ]
    }
  ]
}`)()

			c, ok, err := helper.FindServiceCredentials("test-service", "test-key")
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeTrue())
			g.Expect(c).To(gomega.Equal(helper.Credentials{"test-key": "test-value"}))
		})

		it("does not match no service", func() {
			_, ok, err := helper.FindServiceCredentials("test-service")
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		})

		it("does not match multiple services", func() {
			defer test.ReplaceEnv(t, "CNB_SERVICES", `{
  "": [
    {
      "binding_name": "test-service-1",
      "credentials": {
        "test-key": "test-value"
      },
      "instance_name": null,
      "label": null,
      "tags": [
      ]
    },
    {
      "binding_name": "test-service-2",
      "credentials": {
        "test-key": "test-value"
      },
      "instance_name": null,
      "label": null,
      "tags": [
      ]
    }
  ]
}`)()

			_, ok, err := helper.FindServiceCredentials("test-service")
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		})

		it("does not match without Credentials", func() {
			defer test.ReplaceEnv(t, "CNB_SERVICES", `{
  "": [
    {
      "binding_name": null,
      "credentials": {
      },
      "instance_name": null,
      "label": null,
      "tags": [
        "test-service"
      ]
    }
  ]
}`)()

			_, ok, err := helper.FindServiceCredentials("test-service", "test-key")
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(ok).To(gomega.BeFalse())
		})
	}, spec.Report(report.Terminal{}))
}
