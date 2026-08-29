package awssecretsmanagersecretv1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsSecretsManagerSecretSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSecretsManagerSecretSpec Validation Suite")
}

func int32Ptr(i int32) *int32 {
	return &i
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalSecret is the smallest valid manifest: a region and a string value.
func minimalSecret() *AwsSecretsManagerSecretSpec {
	return &AwsSecretsManagerSecretSpec{
		Region:      "us-west-2",
		StringValue: "s3cr3t-value",
	}
}

// lambdaRotation returns a valid self-managed rotation block.
func lambdaRotation() *AwsSecretsManagerSecretRotation {
	return &AwsSecretsManagerSecretRotation{
		RotationLambdaArn:      svr("arn:aws:lambda:us-west-2:123456789012:function:rotate-db"),
		AutomaticallyAfterDays: int32Ptr(30),
	}
}

var _ = ginkgo.Describe("AwsSecretsManagerSecretSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalSecret())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a shell secret (no value arms)", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsSecretsManagerSecretSpec{Region: "us-west-2"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a binary value", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsSecretsManagerSecretSpec{
					Region:      "us-west-2",
					BinaryValue: "aGVsbG8gd29ybGQ=",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with full production configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsSecretsManagerSecretSpec{
					Region:      "us-west-2",
					Description: "Production database credentials",
					KmsKeyId:    svr("arn:aws:kms:us-west-2:123456789012:key/abc-123"),
					StringValue: `{"username":"app","password":"p"}`,
					VersionStages: []string{
						"bluegreen-active",
					},
					ReplicaRegions: []*AwsSecretsManagerSecretReplica{
						{Region: "us-east-1"},
						{Region: "eu-west-1", KmsKeyId: svr("arn:aws:kms:eu-west-1:123456789012:key/def-456")},
					},
					RecoveryWindowInDays: int32Ptr(7),
					Rotation:             lambdaRotation(),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with schedule-expression rotation", func() {
			ginkgo.It("should accept a rate expression", func() {
				spec := minimalSecret()
				spec.Rotation = &AwsSecretsManagerSecretRotation{
					RotationLambdaArn:  svr("arn:aws:lambda:us-west-2:123456789012:function:rotate-db"),
					ScheduleExpression: "rate(10 days)",
					Duration:           "3h",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept a cron expression", func() {
				spec := minimalSecret()
				spec.Rotation = &AwsSecretsManagerSecretRotation{
					RotationLambdaArn:  svr("arn:aws:lambda:us-west-2:123456789012:function:rotate-db"),
					ScheduleExpression: "cron(0 16 1,15 * ? *)",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with external partner rotation", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalSecret()
				spec.Type = "PARTNER/ExamplePartner"
				spec.Rotation = &AwsSecretsManagerSecretRotation{
					ExternalRotationRoleArn: svr("arn:aws:iam::123456789012:role/partner-rotation"),
					ExternalRotationMetadata: []*AwsSecretsManagerSecretRotationMetadata{
						{Key: "tenant", Value: "prod"},
					},
					AutomaticallyAfterDays: int32Ptr(90),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with recovery window boundary values", func() {
			ginkgo.It("should accept 0 (force delete)", func() {
				spec := minimalSecret()
				spec.RecoveryWindowInDays = int32Ptr(0)
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			})
			ginkgo.It("should accept 7 and 30", func() {
				spec := minimalSecret()
				spec.RecoveryWindowInDays = int32Ptr(7)
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
				spec.RecoveryWindowInDays = int32Ptr(30)
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Invalid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("region", func() {
			ginkgo.It("should reject an empty region", func() {
				spec := minimalSecret()
				spec.Region = ""
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("value arms", func() {
			ginkgo.It("should reject both string_value and binary_value", func() {
				spec := minimalSecret()
				spec.BinaryValue = "aGVsbG8="
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject non-base64 binary_value", func() {
				spec := &AwsSecretsManagerSecretSpec{
					Region:      "us-west-2",
					BinaryValue: "not base64!!",
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject version_stages without a value arm", func() {
				spec := &AwsSecretsManagerSecretSpec{
					Region:        "us-west-2",
					VersionStages: []string{"bluegreen-active"},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an empty version stage label", func() {
				spec := minimalSecret()
				spec.VersionStages = []string{""}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("recovery window", func() {
			ginkgo.It("should reject values between 1 and 6", func() {
				spec := minimalSecret()
				spec.RecoveryWindowInDays = int32Ptr(5)
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
			ginkgo.It("should reject values above 30", func() {
				spec := minimalSecret()
				spec.RecoveryWindowInDays = int32Ptr(31)
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("type", func() {
			ginkgo.It("should reject a partner identifier over 256 characters", func() {
				spec := minimalSecret()
				spec.Type = strings.Repeat("x", 257)
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("replicas", func() {
			ginkgo.It("should reject a replica without a region", func() {
				spec := minimalSecret()
				spec.ReplicaRegions = []*AwsSecretsManagerSecretReplica{{}}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("rotation cadence", func() {
			ginkgo.It("should reject both automatically_after_days and schedule_expression", func() {
				spec := minimalSecret()
				spec.Rotation = lambdaRotation()
				spec.Rotation.ScheduleExpression = "rate(10 days)"
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a rotation block with no cadence", func() {
				spec := minimalSecret()
				spec.Rotation = &AwsSecretsManagerSecretRotation{
					RotationLambdaArn: svr("arn:aws:lambda:us-west-2:123456789012:function:rotate-db"),
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject automatically_after_days out of range", func() {
				spec := minimalSecret()
				spec.Rotation = lambdaRotation()
				spec.Rotation.AutomaticallyAfterDays = int32Ptr(1001)
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a malformed schedule expression", func() {
				spec := minimalSecret()
				spec.Rotation = &AwsSecretsManagerSecretRotation{
					RotationLambdaArn:  svr("arn:aws:lambda:us-west-2:123456789012:function:rotate-db"),
					ScheduleExpression: "every 10 days",
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a malformed duration", func() {
				spec := minimalSecret()
				spec.Rotation = lambdaRotation()
				spec.Rotation.Duration = "3d"
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("rotation mechanism", func() {
			ginkgo.It("should reject a rotation block with no mechanism", func() {
				spec := minimalSecret()
				spec.Rotation = &AwsSecretsManagerSecretRotation{
					AutomaticallyAfterDays: int32Ptr(30),
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject both lambda and external role", func() {
				spec := minimalSecret()
				spec.Rotation = lambdaRotation()
				spec.Rotation.ExternalRotationRoleArn = svr("arn:aws:iam::123456789012:role/partner-rotation")
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject external metadata on lambda rotation", func() {
				spec := minimalSecret()
				spec.Rotation = lambdaRotation()
				spec.Rotation.ExternalRotationMetadata = []*AwsSecretsManagerSecretRotationMetadata{
					{Key: "tenant", Value: "prod"},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject metadata items missing key or value", func() {
				spec := minimalSecret()
				spec.Type = "PARTNER/ExamplePartner"
				spec.Rotation = &AwsSecretsManagerSecretRotation{
					ExternalRotationRoleArn:  svr("arn:aws:iam::123456789012:role/partner-rotation"),
					AutomaticallyAfterDays:   int32Ptr(90),
					ExternalRotationMetadata: []*AwsSecretsManagerSecretRotationMetadata{{Key: "tenant"}},
				}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})
	})

	// The explicit AWS-side name: empty defaults to metadata.name; when
	// set, it carries the shapes metadata.name cannot -- hierarchical
	// paths and service-required prefixes.
	ginkgo.Describe("secret_name", func() {

		ginkgo.It("accepts the empty default (metadata.name naming)", func() {
			spec := minimalSecret()
			spec.SecretName = ""
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a hierarchical path", func() {
			spec := minimalSecret()
			spec.SecretName = "prod/payments/db"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a service-required prefix name", func() {
			spec := minimalSecret()
			spec.SecretName = "ecr-pullthroughcache/dockerhub"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects characters outside AWS's secret-name charset", func() {
			spec := minimalSecret()
			spec.SecretName = "prod secrets/db"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
