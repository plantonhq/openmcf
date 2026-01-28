package cloudflared1databasev1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
)

func TestCloudflareD1DatabaseSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "CloudflareD1DatabaseSpec Custom Validation Tests")
}

var _ = ginkgo.Describe("CloudflareD1DatabaseSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("cloudflare_d1_database", func() {

			ginkgo.It("should not return a validation error for minimal valid fields", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:    "test-account-123",
						DatabaseName: "test-database",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with optional region specified", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:    "test-account-123",
						DatabaseName: "test-database",
						Region:       CloudflareD1Region_weur,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with read replication enabled", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:    "test-account-123",
						DatabaseName: "test-database",
						Region:       CloudflareD1Region_enam,
						ReadReplication: &CloudflareD1ReadReplication{
							Mode: "auto",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for database_name at max length (64 chars)", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:    "test-account-123",
						DatabaseName: "a234567890123456789012345678901234567890123456789012345678901234", // exactly 64 chars
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for all regions", func() {
				regions := []CloudflareD1Region{
					CloudflareD1Region_weur,
					CloudflareD1Region_eeur,
					CloudflareD1Region_apac,
					CloudflareD1Region_oc,
					CloudflareD1Region_wnam,
					CloudflareD1Region_enam,
				}

				for _, region := range regions {
					input := &CloudflareD1Database{
						ApiVersion: "cloudflare.openmcf.org/v1",
						Kind:       "CloudflareD1Database",
						Metadata: &shared.CloudResourceMetadata{
							Name: "test-d1-database",
						},
						Spec: &CloudflareD1DatabaseSpec{
							AccountId:    "test-account-123",
							DatabaseName: "test-database",
							Region:       region,
						},
					}
					err := protovalidate.Validate(input)
					gomega.Expect(err).To(gomega.BeNil())
				}
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("cloudflare_d1_database", func() {

			ginkgo.It("should return a validation error when account_id is missing", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						DatabaseName: "test-database",
						// AccountId is missing
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when database_name is missing", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId: "test-account-123",
						// DatabaseName is missing
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when database_name exceeds max length (65 chars)", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:    "test-account-123",
						DatabaseName: "a2345678901234567890123456789012345678901234567890123456789012345", // 65 chars
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when account_id is empty string", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:    "",
						DatabaseName: "test-database",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when database_name is empty string", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:    "test-account-123",
						DatabaseName: "",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when read_replication mode is missing", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:       "test-account-123",
						DatabaseName:    "test-database",
						ReadReplication: &CloudflareD1ReadReplication{
							// Mode is missing
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when read_replication mode is empty string", func() {
				input := &CloudflareD1Database{
					ApiVersion: "cloudflare.openmcf.org/v1",
					Kind:       "CloudflareD1Database",
					Metadata: &shared.CloudResourceMetadata{
						Name: "test-d1-database",
					},
					Spec: &CloudflareD1DatabaseSpec{
						AccountId:    "test-account-123",
						DatabaseName: "test-database",
						ReadReplication: &CloudflareD1ReadReplication{
							Mode: "",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
