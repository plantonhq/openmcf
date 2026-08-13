package awssagemakerimagev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsSagemakerImageSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerImageSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalImage is the smallest valid manifest: region and the pull
// role; the registry entry alone, no versions yet.
func minimalImage() *AwsSagemakerImageSpec {
	return &AwsSagemakerImageSpec{
		Region:  "us-west-2",
		RoleArn: svr("arn:aws:iam::123456789012:role/sagemaker-execution"),
	}
}

var _ = ginkgo.Describe("AwsSagemakerImageSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalImage())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a fully-annotated version", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalImage()
				spec.DisplayName = "PyTorch Kernels"
				spec.Description = "Team PyTorch kernel images"
				spec.Versions = []*AwsSagemakerImageVersion{
					{
						BaseImage:       "123456789012.dkr.ecr.us-west-2.amazonaws.com/kernels:pytorch-2.4",
						Aliases:         []string{"latest", "stable"},
						Horovod:         true,
						JobType:         "NOTEBOOK_KERNEL",
						MlFramework:     "PyTorch 2.4",
						Processor:       "GPU",
						ProgrammingLang: "python 3.12",
						ReleaseNotes:    "CUDA 12.4 base",
						VendorGuidance:  "STABLE",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with a version missing its base image", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalImage()
				spec.Versions = []*AwsSagemakerImageVersion{{Aliases: []string{"latest"}}}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with whitespace in the base image", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalImage()
				spec.Versions = []*AwsSagemakerImageVersion{{BaseImage: "bad image:1"}}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate aliases", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalImage()
				spec.Versions = []*AwsSagemakerImageVersion{
					{
						BaseImage: "123456789012.dkr.ecr.us-west-2.amazonaws.com/kernels:1",
						Aliases:   []string{"latest", "latest"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad job type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalImage()
				spec.Versions = []*AwsSagemakerImageVersion{
					{
						BaseImage: "123456789012.dkr.ecr.us-west-2.amazonaws.com/kernels:1",
						JobType:   "SERVING",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad ml_framework shape", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalImage()
				spec.Versions = []*AwsSagemakerImageVersion{
					{
						BaseImage:   "123456789012.dkr.ecr.us-west-2.amazonaws.com/kernels:1",
						MlFramework: "PyTorch",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
