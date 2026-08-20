package awssagemakerendpointv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsSagemakerEndpointSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSagemakerEndpointSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalEndpoint is the smallest valid manifest: one serverless
// variant serving a model.
func minimalEndpoint() *AwsSagemakerEndpointSpec {
	return &AwsSagemakerEndpointSpec{
		Region: "us-west-2",
		ProductionVariants: []*AwsSagemakerEndpointVariant{
			{
				Model: svr("my-model"),
				Serverless: &AwsSagemakerEndpointServerlessConfig{
					MaxConcurrency: 5,
					MemorySizeMb:   2048,
				},
			},
		},
	}
}

// instanceVariant is a dedicated-instance variant.
func instanceVariant(name string) *AwsSagemakerEndpointVariant {
	return &AwsSagemakerEndpointVariant{
		VariantName:          name,
		Model:                svr("my-model"),
		InstanceType:         "ml.m5.large",
		InitialInstanceCount: proto.Int32(1),
	}
}

var _ = ginkgo.Describe("AwsSagemakerEndpointSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalEndpoint())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full instance-backed surface", func() {
			ginkgo.It("should not return a validation error", func() {
				v := instanceVariant("primary")
				v.InitialVariantWeight = proto.Float32(1)
				v.RoutingStrategy = "LEAST_OUTSTANDING_REQUESTS"
				v.VolumeSizeGb = proto.Int32(64)
				v.ContainerStartupHealthCheckTimeoutSeconds = proto.Int32(300)
				v.ModelDataDownloadTimeoutSeconds = proto.Int32(600)
				v.EnableSsmAccess = true
				v.InferenceAmiVersion = "al2023-ami-sagemaker-inference-gpu-4-1"
				v.ManagedInstanceScaling = &AwsSagemakerEndpointManagedInstanceScaling{
					Status:           "ENABLED",
					MinInstanceCount: proto.Int32(1),
					MaxInstanceCount: proto.Int32(4),
				}
				v.CoreDump = &AwsSagemakerEndpointCoreDump{
					DestinationS3Uri: "s3://my-dumps/endpoint/",
					KmsKeyArn:        svr("arn:aws:kms:us-west-2:123456789012:key/abc"),
				}
				spec := &AwsSagemakerEndpointSpec{
					Region:             "us-west-2",
					ProductionVariants: []*AwsSagemakerEndpointVariant{v},
					KmsKeyArn:          svr("arn:aws:kms:us-west-2:123456789012:key/abc"),
					AsyncInference: &AwsSagemakerEndpointAsyncInference{
						OutputS3Path:                        "s3://my-bucket/async-out/",
						FailureS3Path:                       "s3://my-bucket/async-fail/",
						MaxConcurrentInvocationsPerInstance: proto.Int32(4),
						SuccessTopicArn:                     svr("arn:aws:sns:us-west-2:123456789012:ok"),
						ErrorTopicArn:                       svr("arn:aws:sns:us-west-2:123456789012:err"),
						IncludeInferenceResponseIn:          []string{"SUCCESS_NOTIFICATION_TOPIC"},
					},
					DataCapture: &AwsSagemakerEndpointDataCapture{
						DestinationS3Uri:          "s3://my-bucket/capture/",
						InitialSamplingPercentage: 100,
						CaptureModes:              []string{"Input", "Output"},
						EnableCapture:             true,
						JsonContentTypes:          []string{"application/json"},
					},
					Deployment: &AwsSagemakerEndpointDeployment{
						BlueGreen: &AwsSagemakerEndpointBlueGreenPolicy{
							TrafficRoutingType:  "CANARY",
							WaitIntervalSeconds: 300,
							CanarySize: &AwsSagemakerEndpointCapacitySize{
								Type:  "CAPACITY_PERCENT",
								Value: 20,
							},
							TerminationWaitSeconds:         proto.Int32(120),
							MaximumExecutionTimeoutSeconds: proto.Int32(3600),
						},
						AutoRollbackAlarmNames: []string{"endpoint-5xx"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with shadow testing one-and-one", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsSagemakerEndpointSpec{
					Region:             "us-west-2",
					ProductionVariants: []*AwsSagemakerEndpointVariant{instanceVariant("prod")},
					ShadowVariants:     []*AwsSagemakerEndpointVariant{instanceVariant("shadow")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a rolling deployment policy", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEndpoint()
				spec.Deployment = &AwsSagemakerEndpointDeployment{
					Rolling: &AwsSagemakerEndpointRollingPolicy{
						MaximumBatchSize: &AwsSagemakerEndpointCapacitySize{
							Type:  "CAPACITY_PERCENT",
							Value: 25,
						},
						WaitIntervalSeconds: 120,
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an inference-components variant and configuration role", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsSagemakerEndpointSpec{
					Region: "us-west-2",
					ProductionVariants: []*AwsSagemakerEndpointVariant{
						{
							VariantName:          "components",
							InstanceType:         "ml.m5.large",
							InitialInstanceCount: proto.Int32(1),
						},
					},
					ExecutionRoleArn: svr("arn:aws:iam::123456789012:role/sagemaker-execution"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with no production variants", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEndpoint()
				spec.ProductionVariants = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate variant names across sides", func() {
			ginkgo.It("should return a validation error", func() {
				spec := &AwsSagemakerEndpointSpec{
					Region:             "us-west-2",
					ProductionVariants: []*AwsSagemakerEndpointVariant{instanceVariant("same")},
					ShadowVariants:     []*AwsSagemakerEndpointVariant{instanceVariant("same")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with shadow testing against two production variants", func() {
			ginkgo.It("should return a validation error", func() {
				spec := &AwsSagemakerEndpointSpec{
					Region: "us-west-2",
					ProductionVariants: []*AwsSagemakerEndpointVariant{
						instanceVariant("a"),
						instanceVariant("b"),
					},
					ShadowVariants: []*AwsSagemakerEndpointVariant{instanceVariant("shadow")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a model-less variant and no configuration role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := &AwsSagemakerEndpointSpec{
					Region: "us-west-2",
					ProductionVariants: []*AwsSagemakerEndpointVariant{
						{
							InstanceType:         "ml.m5.large",
							InitialInstanceCount: proto.Int32(1),
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with serverless and instance settings on one variant", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEndpoint()
				spec.ProductionVariants[0].InstanceType = "ml.m5.large"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a variant that declares no compute", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEndpoint()
				spec.ProductionVariants[0].Serverless = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with provisioned concurrency above max", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEndpoint()
				spec.ProductionVariants[0].Serverless.ProvisionedConcurrency = proto.Int32(10)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an off-step serverless memory size", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEndpoint()
				spec.ProductionVariants[0].Serverless.MemorySizeMb = 1536
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a scaling floor above the ceiling", func() {
			ginkgo.It("should return a validation error", func() {
				v := instanceVariant("primary")
				v.ManagedInstanceScaling = &AwsSagemakerEndpointManagedInstanceScaling{
					MinInstanceCount: proto.Int32(5),
					MaxInstanceCount: proto.Int32(2),
				}
				spec := &AwsSagemakerEndpointSpec{
					Region:             "us-west-2",
					ProductionVariants: []*AwsSagemakerEndpointVariant{v},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with both deployment strategies", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEndpoint()
				spec.Deployment = &AwsSagemakerEndpointDeployment{
					BlueGreen: &AwsSagemakerEndpointBlueGreenPolicy{
						TrafficRoutingType:  "ALL_AT_ONCE",
						WaitIntervalSeconds: 60,
					},
					Rolling: &AwsSagemakerEndpointRollingPolicy{
						MaximumBatchSize: &AwsSagemakerEndpointCapacitySize{
							Type:  "INSTANCE_COUNT",
							Value: 1,
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a canary size on a LINEAR rollout", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEndpoint()
				spec.Deployment = &AwsSagemakerEndpointDeployment{
					BlueGreen: &AwsSagemakerEndpointBlueGreenPolicy{
						TrafficRoutingType:  "LINEAR",
						WaitIntervalSeconds: 60,
						CanarySize: &AwsSagemakerEndpointCapacitySize{
							Type:  "CAPACITY_PERCENT",
							Value: 10,
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with data capture missing its modes", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEndpoint()
				spec.DataCapture = &AwsSagemakerEndpointDataCapture{
					DestinationS3Uri:          "s3://my-bucket/capture/",
					InitialSamplingPercentage: 50,
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a bad instance type prefix", func() {
			ginkgo.It("should return a validation error", func() {
				v := instanceVariant("primary")
				v.InstanceType = "m5.large"
				spec := &AwsSagemakerEndpointSpec{
					Region:             "us-west-2",
					ProductionVariants: []*AwsSagemakerEndpointVariant{v},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
