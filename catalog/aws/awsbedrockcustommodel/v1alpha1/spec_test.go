package awsbedrockcustommodelv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockCustomModelSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockCustomModelSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalCustomModel is the smallest valid manifest.
func minimalCustomModel() *AwsBedrockCustomModelSpec {
	return &AwsBedrockCustomModelSpec{
		Region:            "us-east-1",
		BaseModelArn:      "arn:aws:bedrock:us-east-1::foundation-model/amazon.titan-text-lite-v1",
		Hyperparameters:   map[string]string{"epochCount": "1"},
		RoleArn:           svr("arn:aws:iam::123456789012:role/bedrock-customization"),
		TrainingDataS3Uri: "s3://my-training-bucket/data/train.jsonl",
		OutputDataS3Uri:   "s3://my-training-bucket/output/",
	}
}

var _ = ginkgo.Describe("AwsBedrockCustomModelSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalCustomModel())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCustomModel()
				spec.CustomizationType = "FINE_TUNING"
				spec.JobName = "titan-ft-2026-08-13"
				spec.Hyperparameters = map[string]string{
					"epochCount":   "1",
					"batchSize":    "1",
					"learningRate": "0.00001",
				}
				spec.CustomModelKmsKeyArn = svr("arn:aws:kms:us-east-1:123456789012:key/abc-123")
				spec.ValidationDataS3Uris = []string{
					"s3://my-training-bucket/data/validate-a.jsonl",
					"s3://my-training-bucket/data/validate-b.jsonl",
				}
				spec.VpcConfig = &AwsBedrockCustomModelVpcConfig{
					SubnetIds:        []*foreignkeyv1.StringValueOrRef{svr("subnet-0abc")},
					SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svr("sg-0abc")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with continued pre-training", func() {
			ginkgo.It("should accept the customization type", func() {
				spec := minimalCustomModel()
				spec.CustomizationType = "CONTINUED_PRE_TRAINING"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Required fields
	// -----------------------------------------------------------------
	ginkgo.Describe("Required fields", func() {

		ginkgo.It("should reject a missing region", func() {
			spec := minimalCustomModel()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a missing base model ARN", func() {
			spec := minimalCustomModel()
			spec.BaseModelArn = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject empty hyperparameters", func() {
			spec := minimalCustomModel()
			spec.Hyperparameters = map[string]string{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a missing role", func() {
			spec := minimalCustomModel()
			spec.RoleArn = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a missing training data URI", func() {
			spec := minimalCustomModel()
			spec.TrainingDataS3Uri = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a missing output data URI", func() {
			spec := minimalCustomModel()
			spec.OutputDataS3Uri = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Formats and domains
	// -----------------------------------------------------------------
	ginkgo.Describe("Formats and domains", func() {

		ginkgo.It("should reject a base model identifier that is not a foundation-model ARN", func() {
			spec := minimalCustomModel()
			spec.BaseModelArn = "amazon.titan-text-lite-v1"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an unknown customization type", func() {
			spec := minimalCustomModel()
			spec.CustomizationType = "IMPORTED"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a job name starting with a hyphen", func() {
			spec := minimalCustomModel()
			spec.JobName = "-bad-job"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a non-S3 training URI", func() {
			spec := minimalCustomModel()
			spec.TrainingDataS3Uri = "https://bucket/train.jsonl"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject more than 10 validation URIs", func() {
			spec := minimalCustomModel()
			for i := 0; i < 11; i++ {
				spec.ValidationDataS3Uris = append(spec.ValidationDataS3Uris, "s3://b/v.jsonl")
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a VPC config without subnets", func() {
			spec := minimalCustomModel()
			spec.VpcConfig = &AwsBedrockCustomModelVpcConfig{
				SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svr("sg-0abc")},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a VPC config without security groups", func() {
			spec := minimalCustomModel()
			spec.VpcConfig = &AwsBedrockCustomModelVpcConfig{
				SubnetIds: []*foreignkeyv1.StringValueOrRef{svr("subnet-0abc")},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
