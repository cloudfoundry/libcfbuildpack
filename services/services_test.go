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

package services_test

import (
	"testing"

	bp "github.com/buildpack/libbuildpack/services"
	"github.com/cloudfoundry/libcfbuildpack/services"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func TestServices(t *testing.T) {
	spec.Run(t, "Services", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		when("FindServiceCredentials", func() {
			it("matches single service by BindingName", func() {
				s := services.Services{Services: bp.Services{services.Service{BindingName: "test-service-1"}}}

				_, ok := s.FindServiceCredentials("test-service")

				g.Expect(ok).To(gomega.BeTrue())
			})

			it("matches single service by InstanceName", func() {
				s := services.Services{Services: bp.Services{services.Service{InstanceName: "test-service-1"}}}

				_, ok := s.FindServiceCredentials("test-service")

				g.Expect(ok).To(gomega.BeTrue())
			})

			it("matches single service by Label", func() {
				s := services.Services{Services: bp.Services{services.Service{Label: "test-service-1"}}}

				_, ok := s.FindServiceCredentials("test-service")

				g.Expect(ok).To(gomega.BeTrue())
			})

			it("matches single service by Tags", func() {
				s := services.Services{Services: bp.Services{services.Service{Tags: []string{"test-service-1"}}}}

				_, ok := s.FindServiceCredentials("test-service")

				g.Expect(ok).To(gomega.BeTrue())
			})

			it("matches single service with Credentials", func() {
				s := services.Services{Services: bp.Services{services.Service{
					BindingName: "test-service-1",
					Credentials: services.Credentials{"test-credential": "test-payload"},
				}}}

				c, ok := s.FindServiceCredentials("test-service", "test-credential")

				g.Expect(c).To(gomega.Equal(services.Credentials{"test-credential": "test-payload"}))
				g.Expect(ok).To(gomega.BeTrue())
			})

			it("does not match no service", func() {
				s := services.Services{Services: bp.Services{}}

				_, ok := s.FindServiceCredentials("test-service")

				g.Expect(ok).To(gomega.BeFalse())
			})

			it("does not match multiple services", func() {
				s := services.Services{Services: bp.Services{
					services.Service{BindingName: "test-service-1"},
					services.Service{BindingName: "test-service-2"},
				}}

				_, ok := s.FindServiceCredentials("test-service")

				g.Expect(ok).To(gomega.BeFalse())
			})

			it("does not match without Credentials", func() {
				s := services.Services{Services: bp.Services{services.Service{BindingName: "test-service-1"}}}

				_, ok := s.FindServiceCredentials("test-service", "test-credential")

				g.Expect(ok).To(gomega.BeFalse())
			})
		})

		when("HasService", func() {

			it("matches single service by BindingName", func() {
				s := services.Services{Services: bp.Services{services.Service{BindingName: "test-service-1"}}}

				g.Expect(s.HasService("test-service")).To(gomega.BeTrue())
			})

			it("matches single service by InstanceName", func() {
				s := services.Services{Services: bp.Services{services.Service{InstanceName: "test-service-1"}}}

				g.Expect(s.HasService("test-service")).To(gomega.BeTrue())
			})

			it("matches single service by Label", func() {
				s := services.Services{Services: bp.Services{services.Service{Label: "test-service-1"}}}

				g.Expect(s.HasService("test-service")).To(gomega.BeTrue())
			})

			it("matches single service by Tags", func() {
				s := services.Services{Services: bp.Services{services.Service{Tags: []string{"test-service-1"}}}}

				g.Expect(s.HasService("test-service")).To(gomega.BeTrue())
			})

			it("matches single service with Credentials", func() {
				s := services.Services{Services: bp.Services{services.Service{
					BindingName: "test-service-1",
					Credentials: services.Credentials{"test-credential": "test-payload"},
				}}}

				g.Expect(s.HasService("test-service", "test-credential")).To(gomega.BeTrue())
			})

			it("does not match no service", func() {
				s := services.Services{Services: bp.Services{}}

				g.Expect(s.HasService("test-service")).To(gomega.BeFalse())
			})

			it("does not match multiple services", func() {
				s := services.Services{Services: bp.Services{
					services.Service{BindingName: "test-service-1"},
					services.Service{BindingName: "test-service-2"},
				}}

				g.Expect(s.HasService("test-service")).To(gomega.BeFalse())
			})

			it("does not match without Credentials", func() {
				s := services.Services{Services: bp.Services{services.Service{BindingName: "test-service-1"}}}

				g.Expect(s.HasService("test-service", "test-credential")).To(gomega.BeFalse())
			})
		})
	})
}
