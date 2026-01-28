package cloudflarednsrecordv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestCloudflareDnsRecordSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "CloudflareDnsRecordSpec Custom Validation Tests")
}

var _ = ginkgo.Describe("CloudflareDnsRecordSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("cloudflare_dns_record", func() {

			ginkgo.It("should not return a validation error for minimal valid A record", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-a-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "www",
						Type:   CloudflareDnsRecordType_A,
						Value:  "192.0.2.1",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for AAAA record", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-aaaa-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "www",
						Type:   CloudflareDnsRecordType_AAAA,
						Value:  "2001:db8::1",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for CNAME record", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-cname-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "app",
						Type:   CloudflareDnsRecordType_CNAME,
						Value:  "www.example.com",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for MX record with priority", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-mx-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId:   "abc123def456",
						Name:     "@",
						Type:     CloudflareDnsRecordType_MX,
						Value:    "mail.example.com",
						Priority: 10,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for TXT record", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-txt-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "@",
						Type:   CloudflareDnsRecordType_TXT,
						Value:  "v=spf1 include:_spf.google.com ~all",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for proxied A record", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-proxied-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId:  "abc123def456",
						Name:    "www",
						Type:    CloudflareDnsRecordType_A,
						Value:   "192.0.2.1",
						Proxied: true,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for record with TTL of 1 (auto)", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-auto-ttl-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "www",
						Type:   CloudflareDnsRecordType_A,
						Value:  "192.0.2.1",
						Ttl:    1,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for record with valid TTL", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-ttl-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "www",
						Type:   CloudflareDnsRecordType_A,
						Value:  "192.0.2.1",
						Ttl:    3600,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for record with comment", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-comment-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId:  "abc123def456",
						Name:    "www",
						Type:    CloudflareDnsRecordType_A,
						Value:   "192.0.2.1",
						Comment: "Main web server",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for CAA record", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-caa-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "@",
						Type:   CloudflareDnsRecordType_CAA,
						Value:  "0 issue \"letsencrypt.org\"",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("cloudflare_dns_record", func() {

			ginkgo.It("should return a validation error when zone_id is missing", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						Name:  "www",
						Type:  CloudflareDnsRecordType_A,
						Value: "192.0.2.1",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when name is missing", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Type:   CloudflareDnsRecordType_A,
						Value:  "192.0.2.1",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when type is unspecified", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "www",
						Type:   CloudflareDnsRecordType_dns_record_type_unspecified,
						Value:  "192.0.2.1",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when value is missing", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "www",
						Type:   CloudflareDnsRecordType_A,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for invalid TTL", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "www",
						Type:   CloudflareDnsRecordType_A,
						Value:  "192.0.2.1",
						Ttl:    30,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for TTL exceeding max", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "www",
						Type:   CloudflareDnsRecordType_A,
						Value:  "192.0.2.1",
						Ttl:    100000,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for negative priority", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId:   "abc123def456",
						Name:     "@",
						Type:     CloudflareDnsRecordType_MX,
						Value:    "mail.example.com",
						Priority: -1,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for priority exceeding max", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId:   "abc123def456",
						Name:     "@",
						Type:     CloudflareDnsRecordType_MX,
						Value:    "mail.example.com",
						Priority: 70000,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for comment exceeding max length", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId:  "abc123def456",
						Name:    "www",
						Type:    CloudflareDnsRecordType_A,
						Value:   "192.0.2.1",
						Comment: "This is a very long comment that exceeds the 100 character limit for DNS record comments in Cloudflare's system.",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for proxied TXT record", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId:  "abc123def456",
						Name:    "@",
						Type:    CloudflareDnsRecordType_TXT,
						Value:   "test",
						Proxied: true,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for proxied MX record", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId:   "abc123def456",
						Name:     "@",
						Type:     CloudflareDnsRecordType_MX,
						Value:    "mail.example.com",
						Priority: 10,
						Proxied:  true,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for MX record without priority", func() {
				input := &CloudflareDnsRecord{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareDnsRecord",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-record",
					},
					Spec: &CloudflareDnsRecordSpec{
						ZoneId: "abc123def456",
						Name:   "@",
						Type:   CloudflareDnsRecordType_MX,
						Value:  "mail.example.com",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
