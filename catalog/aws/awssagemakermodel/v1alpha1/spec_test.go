package awssagemakermodelv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsSagemakerModelSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerModelSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalModel is the smallest valid manifest: region, execution role,
// and a single container serving a prebuilt image.
func minimalModel() *AwsSagemakerModelSpec {
	return &AwsSagemakerModelSpec{
		Region:           "us-west-2",
		ExecutionRoleArn: svr("arn:aws:iam::123456789012:role/sagemaker-execution"),
		PrimaryContainer: &AwsSagemakerModelContainer{
			Image: "246618743249.dkr.ecr.us-west-2.amazonaws.com/sagemaker-scikit-learn:1.2-1",
		},
	}
}

var _ = ginkgo.Describe("AwsSagemakerModelSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalModel())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full single-container surface", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalModel()
				spec.EnableNetworkIsolation = true
				spec.PrimaryContainer = &AwsSagemakerModelContainer{
					Image:             "123456789012.dkr.ecr.us-west-2.amazonaws.com/serve:latest",
					ContainerHostname: "primary",
					Environment:       map[string]string{"SAGEMAKER_PROGRAM": "serve.py"},
					Mode:              "MultiModel",
					MultiModelCache:   "Disabled",
					ModelDataSource: &AwsSagemakerModelS3DataSource{
						S3Uri:           "s3://my-models/llm/",
						S3DataType:      "S3Prefix",
						CompressionType: "None",
						AcceptEula:      true,
					},
					AdditionalModelDataSources: []*AwsSagemakerModelAdditionalDataSource{
						{
							ChannelName: "adapters",
							Source: &AwsSagemakerModelS3DataSource{
								S3Uri:           "s3://my-models/adapters/",
								S3DataType:      "S3Prefix",
								CompressionType: "None",
							},
						},
					},
					ImageConfig: &AwsSagemakerModelImageConfig{
						RepositoryAccessMode:             "Vpc",
						RepositoryCredentialsProviderArn: "arn:aws:lambda:us-west-2:123456789012:function:creds",
					},
				}
				spec.VpcConfig = &AwsSagemakerModelVpcConfig{
					SubnetIds:        []*foreignkeyv1.StringValueOrRef{svr("subnet-0abc")},
					SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svr("sg-0abc")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an inference pipeline and Direct invocation", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer = nil
				spec.Containers = []*AwsSagemakerModelContainer{
					{Image: "123456789012.dkr.ecr.us-west-2.amazonaws.com/pre:1", ContainerHostname: "pre"},
					{Image: "123456789012.dkr.ecr.us-west-2.amazonaws.com/infer:1", ContainerHostname: "infer"},
				}
				spec.InferenceExecutionMode = "Direct"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a model package container", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer = &AwsSagemakerModelContainer{
					ModelPackageArn:            "arn:aws:sagemaker:us-west-2:123456789012:model-package/my-group/3",
					InferenceSpecificationName: "default",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with both container forms set", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.Containers = []*AwsSagemakerModelContainer{{Image: "123456789012.dkr.ecr.us-west-2.amazonaws.com/x:1"}}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with neither container form set", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an execution mode on a single-container model", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.InferenceExecutionMode = "Serial"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a container that has neither image nor model package", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer = &AwsSagemakerModelContainer{ContainerHostname: "empty"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with both artifact forms on one container", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer.ModelDataUrl = "s3://my-models/model.tar.gz"
				spec.PrimaryContainer.ModelDataSource = &AwsSagemakerModelS3DataSource{
					S3Uri:           "s3://my-models/llm/",
					S3DataType:      "S3Prefix",
					CompressionType: "None",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an environment key starting with a digit", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer.Environment = map[string]string{"1BAD": "x"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with cache tuning on a SingleModel container", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer.MultiModelCache = "Enabled"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with registry credentials on Platform access mode", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer.ImageConfig = &AwsSagemakerModelImageConfig{
					RepositoryAccessMode:             "Platform",
					RepositoryCredentialsProviderArn: "arn:aws:lambda:us-west-2:123456789012:function:creds",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a whitespace image path", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer.Image = "bad image"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with too many pipeline containers", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.PrimaryContainer = nil
				for i := 0; i < 16; i++ {
					spec.Containers = append(spec.Containers, &AwsSagemakerModelContainer{
						Image: "123456789012.dkr.ecr.us-west-2.amazonaws.com/x:1",
					})
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an empty vpc_config subnet list", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalModel()
				spec.VpcConfig = &AwsSagemakerModelVpcConfig{
					SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svr("sg-0abc")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
