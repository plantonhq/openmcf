package awssagemakernotebookinstancev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsSagemakerNotebookInstanceSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerNotebookInstanceSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalNotebook is the smallest valid manifest: region, instance
// type, and the execution role.
func minimalNotebook() *AwsSagemakerNotebookInstanceSpec {
	return &AwsSagemakerNotebookInstanceSpec{
		Region:       "us-west-2",
		InstanceType: "ml.t3.medium",
		RoleArn:      svr("arn:aws:iam::123456789012:role/sagemaker-execution"),
	}
}

var _ = ginkgo.Describe("AwsSagemakerNotebookInstanceSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalNotebook())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full VPC-confined surface", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalNotebook()
				spec.VolumeSizeGb = proto.Int32(50)
				spec.SubnetId = svr("subnet-0abc")
				spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{svr("sg-0abc")}
				spec.KmsKeyArn = svr("arn:aws:kms:us-west-2:123456789012:key/abc")
				spec.DirectInternetAccess = "Disabled"
				spec.RootAccess = "Disabled"
				spec.PlatformIdentifier = "notebook-al2023-v1"
				spec.DefaultCodeRepository = "https://github.com/example/notebooks.git"
				spec.AdditionalCodeRepositories = []string{"https://github.com/example/lib.git"}
				spec.ImdsMinimumVersion = "2"
				spec.LifecycleConfig = &AwsSagemakerNotebookInstanceLifecycleConfig{
					OnCreate: "#!/bin/bash\npip install pandas",
					OnStart:  "#!/bin/bash\necho started",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an on_start-only lifecycle configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalNotebook()
				spec.LifecycleConfig = &AwsSagemakerNotebookInstanceLifecycleConfig{
					OnStart: "#!/bin/bash\necho started",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with disabled internet but no VPC wiring", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalNotebook()
				spec.DirectInternetAccess = "Disabled"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with security groups but no subnet", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalNotebook()
				spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{svr("sg-0abc")}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an empty lifecycle configuration", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalNotebook()
				spec.LifecycleConfig = &AwsSagemakerNotebookInstanceLifecycleConfig{}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a volume below the AWS floor", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalNotebook()
				spec.VolumeSizeGb = proto.Int32(2)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad platform identifier", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalNotebook()
				spec.PlatformIdentifier = "notebook-al2024-v9"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad IMDS version", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalNotebook()
				spec.ImdsMinimumVersion = "3"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with four additional code repositories", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalNotebook()
				spec.AdditionalCodeRepositories = []string{"a", "b", "c", "d"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a non-ml instance type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalNotebook()
				spec.InstanceType = "t3.medium"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
