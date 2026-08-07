package awsecrrepov1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsEcrRepoSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEcrRepoSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// helper to build a Struct from a map, failing loudly on error.
func mustStruct(m map[string]interface{}) *structpb.Struct {
	s, err := structpb.NewStruct(m)
	if err != nil {
		panic(err)
	}
	return s
}

var _ = ginkgo.Describe("AwsEcrRepoSpec validations", func() {
	var spec *AwsEcrRepoSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec.
		spec = &AwsEcrRepoSpec{
			Region:         "us-west-2",
			RepositoryName: "team-blue/checkout-service",
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal repository", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts the full production surface", func() {
		spec.ImageTagMutability = proto.String("IMMUTABLE_WITH_EXCLUSION")
		spec.ImageTagMutabilityExclusionFilters = []string{"latest", "dev-*"}
		spec.EncryptionType = proto.String("KMS")
		spec.KmsKeyId = strRef("arn:aws:kms:us-west-2:123456789012:key/abc")
		spec.ScanOnPush = proto.Bool(true)
		spec.ForceDelete = true
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{
			{
				RulePriority: 1,
				Description:  "Expire untagged images after 14 days",
				TagStatus:    "untagged",
				CountType:    "sinceImagePushed",
				CountNumber:  14,
			},
			{
				RulePriority: 2,
				TagStatus:    "tagged",
				TagPrefixes:  []string{"pr-"},
				CountType:    "imageCountMoreThan",
				CountNumber:  10,
			},
			{
				RulePriority: 3,
				TagStatus:    "any",
				CountType:    "imageCountMoreThan",
				CountNumber:  100,
			},
		}
		spec.RepositoryPolicy = mustStruct(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []interface{}{map[string]interface{}{
				"Sid":       "AllowLambdaPull",
				"Effect":    "Allow",
				"Principal": map[string]interface{}{"Service": "lambda.amazonaws.com"},
				"Action":    []interface{}{"ecr:BatchGetImage", "ecr:GetDownloadUrlForLayer"},
			}},
		})
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts dual-layer KMS_DSSE encryption without an explicit key", func() {
		spec.EncryptionType = proto.String("KMS_DSSE")
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	ginkgo.It("accepts tag patterns as the tagged-rule selector", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{{
			RulePriority: 1,
			TagStatus:    "tagged",
			TagPatterns:  []string{"*-snapshot"},
			CountType:    "sinceImagePushed",
			CountNumber:  30,
		}}
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// repository_name
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a missing repository name", func() {
		spec.RepositoryName = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects uppercase repository names", func() {
		spec.RepositoryName = "Team/Service"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a trailing slash", func() {
		spec.RepositoryName = "team/"
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Tag mutability + exclusion filters
	// -------------------------------------------------------------------------

	ginkgo.It("rejects an invalid mutability value", func() {
		spec.ImageTagMutability = proto.String("FROZEN")
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects exclusion filters without an exclusion mode", func() {
		spec.ImageTagMutability = proto.String("IMMUTABLE")
		spec.ImageTagMutabilityExclusionFilters = []string{"latest"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an exclusion mode without filters", func() {
		spec.ImageTagMutability = proto.String("MUTABLE_WITH_EXCLUSION")
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a filter with more than two wildcards", func() {
		spec.ImageTagMutability = proto.String("IMMUTABLE_WITH_EXCLUSION")
		spec.ImageTagMutabilityExclusionFilters = []string{"*a*b*"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects more than five exclusion filters", func() {
		spec.ImageTagMutability = proto.String("IMMUTABLE_WITH_EXCLUSION")
		spec.ImageTagMutabilityExclusionFilters = []string{"a", "b", "c", "d", "e", "f"}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Encryption
	// -------------------------------------------------------------------------

	ginkgo.It("rejects an invalid encryption type", func() {
		spec.EncryptionType = proto.String("SSE-S3")
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a KMS key with AES256 encryption", func() {
		spec.EncryptionType = proto.String("AES256")
		spec.KmsKeyId = strRef("arn:aws:kms:us-west-2:123456789012:key/abc")
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a KMS key when encryption type is omitted", func() {
		spec.KmsKeyId = strRef("arn:aws:kms:us-west-2:123456789012:key/abc")
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// Lifecycle rules
	// -------------------------------------------------------------------------

	ginkgo.It("rejects duplicate rule priorities", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{
			{RulePriority: 1, TagStatus: "untagged", CountType: "sinceImagePushed", CountNumber: 7},
			{RulePriority: 1, TagStatus: "any", CountType: "imageCountMoreThan", CountNumber: 50},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects two 'any' rules", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{
			{RulePriority: 1, TagStatus: "any", CountType: "imageCountMoreThan", CountNumber: 50},
			{RulePriority: 2, TagStatus: "any", CountType: "imageCountMoreThan", CountNumber: 100},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects two 'untagged' rules", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{
			{RulePriority: 1, TagStatus: "untagged", CountType: "sinceImagePushed", CountNumber: 7},
			{RulePriority: 2, TagStatus: "untagged", CountType: "sinceImagePushed", CountNumber: 30},
		}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a tagged rule without a selector list", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{{
			RulePriority: 1,
			TagStatus:    "tagged",
			CountType:    "imageCountMoreThan",
			CountNumber:  10,
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a tagged rule with both selector lists", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{{
			RulePriority: 1,
			TagStatus:    "tagged",
			TagPrefixes:  []string{"release-"},
			TagPatterns:  []string{"v1.*"},
			CountType:    "imageCountMoreThan",
			CountNumber:  10,
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an untagged rule carrying tag prefixes", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{{
			RulePriority: 1,
			TagStatus:    "untagged",
			TagPrefixes:  []string{"release-"},
			CountType:    "sinceImagePushed",
			CountNumber:  7,
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects an invalid count type", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{{
			RulePriority: 1,
			TagStatus:    "any",
			CountType:    "olderThan",
			CountNumber:  10,
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	ginkgo.It("rejects a zero count number", func() {
		spec.LifecycleRules = []*AwsEcrRepoLifecycleRule{{
			RulePriority: 1,
			TagStatus:    "any",
			CountType:    "imageCountMoreThan",
			CountNumber:  0,
		}}
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})

	// -------------------------------------------------------------------------
	// region
	// -------------------------------------------------------------------------

	ginkgo.It("rejects a missing region", func() {
		spec.Region = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.Succeed())
	})
})
