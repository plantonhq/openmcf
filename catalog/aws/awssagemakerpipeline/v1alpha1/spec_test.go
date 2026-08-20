package awssagemakerpipelinev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsSagemakerPipelineSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerPipelineSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func inlineDefinition() *structpb.Struct {
	s, _ := structpb.NewStruct(map[string]interface{}{
		"Version": "2020-12-01",
		"Steps": []interface{}{
			map[string]interface{}{
				"Name": "NoOp",
				"Type": "Fail",
				"Arguments": map[string]interface{}{
					"ErrorMessage": "placeholder step",
				},
			},
		},
	})
	return s
}

// minimalPipeline is the smallest valid manifest: region, role, and an
// inline definition.
func minimalPipeline() *AwsSagemakerPipelineSpec {
	return &AwsSagemakerPipelineSpec{
		Region:     "us-west-2",
		RoleArn:    svr("arn:aws:iam::123456789012:role/sagemaker-execution"),
		Definition: inlineDefinition(),
	}
}

var _ = ginkgo.Describe("AwsSagemakerPipelineSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalPipeline())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the S3-location arm and parallelism", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalPipeline()
				spec.Definition = nil
				spec.DefinitionS3Location = &AwsSagemakerPipelineDefinitionS3Location{
					Bucket:    svr("my-pipeline-defs"),
					ObjectKey: "pipelines/train.json",
					VersionId: "abc123",
				}
				spec.DisplayName = "training-pipeline"
				spec.Description = "nightly training"
				spec.ParallelismMaxSteps = proto.Int32(4)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with both definition arms", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalPipeline()
				spec.DefinitionS3Location = &AwsSagemakerPipelineDefinitionS3Location{
					Bucket:    svr("my-pipeline-defs"),
					ObjectKey: "pipelines/train.json",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with neither definition arm", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalPipeline()
				spec.Definition = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an S3 location missing its object key", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalPipeline()
				spec.Definition = nil
				spec.DefinitionS3Location = &AwsSagemakerPipelineDefinitionS3Location{
					Bucket: svr("my-pipeline-defs"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a display name starting with a hyphen", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalPipeline()
				spec.DisplayName = "-bad"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with zero parallelism", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalPipeline()
				spec.ParallelismMaxSteps = proto.Int32(0)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
