package awss3objectsetv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func TestAwsS3ObjectSetSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsS3ObjectSetSpec Validation Tests")
}

// minimalValidSpec returns a minimal valid AwsS3ObjectSetSpec.
func minimalValidSpec() *AwsS3ObjectSetSpec {
	return &AwsS3ObjectSetSpec{
		Bucket: &foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
				Value: "my-test-bucket",
			},
		},
		Region: "us-east-1",
		Objects: []*AwsS3Object{
			{
				Key: "config/app.json",
				Source: &AwsS3Object_Content{
					Content: "{\"key\": \"value\"}",
				},
				ContentType: stringPtr("application/json"),
			},
		},
	}
}

var _ = ginkgo.Describe("AwsS3ObjectSetSpec Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("aws_s3_object_set_spec", func() {

			ginkgo.It("should not return a validation error for minimal valid spec", func() {
				spec := minimalValidSpec()
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with multiple objects", func() {
				spec := minimalValidSpec()
				spec.Objects = append(spec.Objects, &AwsS3Object{
					Key: "assets/logo.png",
					Source: &AwsS3Object_ContentBase64{
						ContentBase64: "iVBORw0KGgoAAAANSUhEUg==",
					},
					ContentType: stringPtr("image/png"),
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with set-level tags", func() {
				spec := minimalValidSpec()
				spec.Tags = map[string]string{
					"environment": "production",
					"team":        "platform",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with object-level tags", func() {
				spec := minimalValidSpec()
				spec.Objects[0].Tags = map[string]string{
					"purpose": "config",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with cache_control and content_encoding", func() {
				spec := minimalValidSpec()
				spec.Objects[0].CacheControl = "max-age=86400"
				spec.Objects[0].ContentEncoding = "gzip"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with acl set", func() {
				spec := minimalValidSpec()
				spec.Objects[0].Acl = "private"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with bucket as value_from reference", func() {
				spec := minimalValidSpec()
				spec.Bucket = &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{
							Name: "my-s3-bucket",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with content headers set", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ContentDisposition = "attachment; filename=\"report.pdf\""
				spec.Objects[0].ContentLanguage = "en-US"
				spec.Objects[0].WebsiteRedirect = "/new-page.html"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with lowercase metadata keys", func() {
				spec := minimalValidSpec()
				spec.Objects[0].Metadata = map[string]string{
					"build-commit": "abc123",
					"generated-by": "planton",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for each valid storage class", func() {
				for _, sc := range []string{"STANDARD", "STANDARD_IA", "ONEZONE_IA", "INTELLIGENT_TIERING", "GLACIER", "GLACIER_IR", "DEEP_ARCHIVE", "REDUCED_REDUNDANCY", "EXPRESS_ONEZONE", "OUTPOSTS", "SNOW", "FSX_OPENZFS", "FSX_ONTAP"} {
					spec := minimalValidSpec()
					spec.Objects[0].StorageClass = sc
					err := protovalidate.Validate(spec)
					gomega.Expect(err).To(gomega.BeNil(), "storage class %s should be valid", sc)
				}
			})

			ginkgo.It("should not return a validation error for each valid server_side_encryption", func() {
				for _, sse := range []string{"AES256", "aws:kms", "aws:kms:dsse", "aws:fsx"} {
					spec := minimalValidSpec()
					spec.Objects[0].ServerSideEncryption = sse
					err := protovalidate.Validate(spec)
					gomega.Expect(err).To(gomega.BeNil(), "sse %s should be valid", sse)
				}
			})

			ginkgo.It("should not return a validation error with a kms key and kms encryption", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ServerSideEncryption = "aws:kms"
				spec.Objects[0].KmsKey = &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
						Value: "arn:aws:kms:us-east-1:123456789012:key/abc",
					},
				}
				spec.Objects[0].BucketKeyEnabled = boolPtr(true)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with a kms key and encryption unset (aws:kms implied)", func() {
				spec := minimalValidSpec()
				spec.Objects[0].KmsKey = &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
						Value: "arn:aws:kms:us-east-1:123456789012:key/abc",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error for each valid checksum algorithm", func() {
				for _, ca := range []string{"CRC32", "CRC32C", "CRC64NVME", "SHA1", "SHA256"} {
					spec := minimalValidSpec()
					spec.Objects[0].ChecksumAlgorithm = ca
					err := protovalidate.Validate(spec)
					gomega.Expect(err).To(gomega.BeNil(), "checksum algorithm %s should be valid", ca)
				}
			})

			ginkgo.It("should not return a validation error with an object lock mode + retain-until pair", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ObjectLockMode = "GOVERNANCE"
				spec.Objects[0].ObjectLockRetainUntilDate = "2027-01-01T00:00:00Z"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with a legal hold alone", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ObjectLockLegalHoldStatus = "ON"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should not return a validation error with force_destroy set", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ForceDestroy = true
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("aws_s3_object_set_spec", func() {

			ginkgo.It("should return a validation error when bucket is missing", func() {
				spec := minimalValidSpec()
				spec.Bucket = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when region is empty", func() {
				spec := minimalValidSpec()
				spec.Region = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when objects list is empty", func() {
				spec := minimalValidSpec()
				spec.Objects = []*AwsS3Object{}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when objects is nil", func() {
				spec := minimalValidSpec()
				spec.Objects = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when object key is empty", func() {
				spec := minimalValidSpec()
				spec.Objects[0].Key = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when object has no content source", func() {
				spec := minimalValidSpec()
				spec.Objects[0].Source = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an uppercase metadata key", func() {
				spec := minimalValidSpec()
				spec.Objects[0].Metadata = map[string]string{
					"Build-Commit": "abc123",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid storage class", func() {
				spec := minimalValidSpec()
				spec.Objects[0].StorageClass = "ARCHIVE"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid server_side_encryption", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ServerSideEncryption = "kms"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when a kms key is paired with AES256", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ServerSideEncryption = "AES256"
				spec.Objects[0].KmsKey = &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
						Value: "arn:aws:kms:us-east-1:123456789012:key/abc",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid checksum algorithm", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ChecksumAlgorithm = "MD5"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid object lock mode", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ObjectLockMode = "LEGAL"
				spec.Objects[0].ObjectLockRetainUntilDate = "2027-01-01T00:00:00Z"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when object lock mode lacks retain-until-date", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ObjectLockMode = "GOVERNANCE"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error when retain-until-date lacks object lock mode", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ObjectLockRetainUntilDate = "2027-01-01T00:00:00Z"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for a non-RFC3339 retain-until-date", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ObjectLockMode = "GOVERNANCE"
				spec.Objects[0].ObjectLockRetainUntilDate = "2027-01-01"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid legal hold status", func() {
				spec := minimalValidSpec()
				spec.Objects[0].ObjectLockLegalHoldStatus = "ENABLED"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should return a validation error for an invalid acl", func() {
				spec := minimalValidSpec()
				spec.Objects[0].Acl = "public"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
