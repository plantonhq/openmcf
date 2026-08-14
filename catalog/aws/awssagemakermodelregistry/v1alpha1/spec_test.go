package awssagemakermodelregistryv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsSagemakerModelRegistrySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerModelRegistrySpec Validation Suite")
}

// minimalRegistry is the smallest valid manifest: just the region (the
// group's name derives from metadata.name).
func minimalRegistry() *AwsSagemakerModelRegistrySpec {
	return &AwsSagemakerModelRegistrySpec{
		Region: "us-west-2",
	}
}

var _ = ginkgo.Describe("AwsSagemakerModelRegistrySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalRegistry())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a description and a cross-account policy", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalRegistry()
				spec.Description = "team model registry"
				policy, err := structpb.NewStruct(map[string]interface{}{
					"Version": "2012-10-17",
					"Statement": []interface{}{
						map[string]interface{}{
							"Sid":       "AllowCrossAccountDescribe",
							"Effect":    "Allow",
							"Principal": map[string]interface{}{"AWS": "arn:aws:iam::210987654321:root"},
							"Action":    "sagemaker:DescribeModelPackage",
							"Resource":  "*",
						},
					},
				})
				gomega.Expect(err).To(gomega.BeNil())
				spec.ResourcePolicy = policy
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with an empty region", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRegistry()
				spec.Region = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an oversized description", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRegistry()
				long := make([]byte, 1025)
				for i := range long {
					long[i] = 'x'
				}
				spec.Description = string(long)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
