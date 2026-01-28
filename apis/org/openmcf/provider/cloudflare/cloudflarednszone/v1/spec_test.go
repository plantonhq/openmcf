package cloudflarednszonev1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestCloudflareDnsZoneSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "CloudflareDnsZoneSpec Custom Validation Tests")
}

var _ = ginkgo.Describe("CloudflareDnsZoneSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("cloudflare_dns_zone", func() {

			ginkgo.It("should not return a validation error for minimal valid fields", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "example.com",
						AccountId: "test-account-123",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with all optional fields specified", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:       "example.com",
						AccountId:      "test-account-123",
						Plan:           CloudflareDnsZoneSpec_PRO,
						Paused:         false,
						DefaultProxied: true,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for free plan", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "free-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "free-example.com",
						AccountId: "test-account-123",
						Plan:      CloudflareDnsZoneSpec_FREE,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for business plan", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "business-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "business.example.com",
						AccountId: "test-account-123",
						Plan:      CloudflareDnsZoneSpec_BUSINESS,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for enterprise plan", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "enterprise-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "enterprise.example.com",
						AccountId: "test-account-123",
						Plan:      CloudflareDnsZoneSpec_ENTERPRISE,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for paused zone", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "paused-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "paused.example.com",
						AccountId: "test-account-123",
						Paused:    true,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with default_proxied enabled", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "proxied-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:       "proxied.example.com",
						AccountId:      "test-account-123",
						DefaultProxied: true,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for subdomain zone", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "subdomain-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "sub.domain.example.com",
						AccountId: "test-account-123",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("cloudflare_dns_zone", func() {

			ginkgo.It("should return a validation error when zone_name is missing", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						// ZoneName is missing
						AccountId: "test-account-123",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when account_id is missing", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName: "example.com",
						// AccountId is missing
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when zone_name is empty string", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "",
						AccountId: "test-account-123",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when account_id is empty string", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "example.com",
						AccountId: "",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for invalid zone_name format (no TLD)", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "invalidzone",
						AccountId: "test-account-123",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for invalid zone_name format (contains uppercase)", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "EXAMPLE.COM",
						AccountId: "test-account-123",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for invalid zone_name format (starts with hyphen)", func() {
				input := &CloudflareDnsZone{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsZone",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-dns-zone",
					},
					Spec: &CloudflareDnsZoneSpec{
						ZoneName:  "-example.com",
						AccountId: "test-account-123",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
