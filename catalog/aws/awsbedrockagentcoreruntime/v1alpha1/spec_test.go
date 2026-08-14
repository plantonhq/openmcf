package awsbedrockagentcoreruntimev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsBedrockAgentCoreRuntimeSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockAgentCoreRuntimeSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalRuntime is the smallest valid manifest: region, name, role,
// a container artifact, and public networking.
func minimalRuntime() *AwsBedrockAgentCoreRuntimeSpec {
	return &AwsBedrockAgentCoreRuntimeSpec{
		Region:      "us-west-2",
		RuntimeName: "support_agent",
		RoleArn:     svr("arn:aws:iam::123456789012:role/agentcore-runtime"),
		Artifact: &AwsBedrockAgentCoreRuntimeArtifact{
			Container: &AwsBedrockAgentCoreRuntimeContainer{
				ImageUri: "123456789012.dkr.ecr.us-west-2.amazonaws.com/agent:v1",
			},
		},
		Network: &AwsBedrockAgentCoreRuntimeNetwork{Mode: "PUBLIC"},
	}
}

var _ = ginkgo.Describe("AwsBedrockAgentCoreRuntimeSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalRuntime())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a code-bundle artifact", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalRuntime()
				spec.Artifact = &AwsBedrockAgentCoreRuntimeArtifact{
					Code: &AwsBedrockAgentCoreRuntimeCode{
						Runtime:    "PYTHON_3_13",
						EntryPoint: []string{"main.py"},
						S3: &AwsBedrockAgentCoreRuntimeCodeS3{
							Bucket: svr("my-code-bucket"),
							Prefix: "bundles/agent.zip",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				policy, perr := structpb.NewStruct(map[string]any{
					"Version": "2012-10-17",
				})
				gomega.Expect(perr).To(gomega.BeNil())
				spec := minimalRuntime()
				spec.Description = "hosts the support agent"
				spec.ServerProtocol = "MCP"
				spec.EnvironmentVariables = map[string]string{"LOG_LEVEL": "info"}
				spec.Lifecycle = &AwsBedrockAgentCoreRuntimeLifecycle{
					IdleRuntimeSessionTimeoutSeconds: 900,
					MaxLifetimeSeconds:               3600,
				}
				spec.Network = &AwsBedrockAgentCoreRuntimeNetwork{
					Mode: "VPC",
					VpcConfig: &AwsBedrockAgentCoreVpcConfig{
						Subnets:        []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
						SecurityGroups: []*foreignkeyv1.StringValueOrRef{svr("sg-abc")},
					},
				}
				spec.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
					DiscoveryUrl:    "https://issuer.example.com/.well-known/openid-configuration",
					AllowedAudience: []string{"agents"},
					CustomClaims: []*AwsBedrockAgentCoreCustomClaim{
						{
							ClaimName:     "org",
							ValueType:     "STRING",
							MatchOperator: "EQUALS",
							MatchValue:    "acme",
						},
					},
					PrivateEndpoint: &AwsBedrockAgentCorePrivateEndpoint{
						ManagedVpc: &AwsBedrockAgentCoreManagedVpcEndpoint{
							VpcId:                 svr("vpc-abc"),
							SubnetIds:             []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
							EndpointIpAddressType: "IPV4",
						},
					},
				}
				spec.RequestHeaderAllowlist = []string{"X-Trace-Id"}
				spec.Filesystems = []*AwsBedrockAgentCoreRuntimeFilesystem{
					{MountPath: "/mnt/scratch", SessionStorage: true},
					{MountPath: "/mnt/shared", EfsAccessPointArn: svr("arn:aws:elasticfilesystem:us-west-2:123456789012:access-point/fsap-abc")},
				}
				spec.Endpoints = []*AwsBedrockAgentCoreRuntimeEndpoint{
					{Name: "live", Description: "production traffic"},
					{Name: "pinned", AgentRuntimeVersion: "1"},
				}
				spec.ResourcePolicy = policy
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with a hyphenated runtime name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.RuntimeName = "support-agent"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with no artifact arm", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Artifact = &AwsBedrockAgentCoreRuntimeArtifact{}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with both artifact arms", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Artifact.Code = &AwsBedrockAgentCoreRuntimeCode{
					Runtime:    "PYTHON_3_13",
					EntryPoint: []string{"main.py"},
					S3: &AwsBedrockAgentCoreRuntimeCodeS3{
						Bucket: svr("b"),
						Prefix: "k",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a code artifact missing its S3 location", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Artifact = &AwsBedrockAgentCoreRuntimeArtifact{
					Code: &AwsBedrockAgentCoreRuntimeCode{
						Runtime:    "PYTHON_3_13",
						EntryPoint: []string{"main.py"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with three entry-point elements", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Artifact = &AwsBedrockAgentCoreRuntimeArtifact{
					Code: &AwsBedrockAgentCoreRuntimeCode{
						Runtime:    "NODE_22",
						EntryPoint: []string{"a", "b", "c"},
						S3: &AwsBedrockAgentCoreRuntimeCodeS3{
							Bucket: svr("b"),
							Prefix: "k",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with VPC mode and no vpc_config", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Network = &AwsBedrockAgentCoreRuntimeNetwork{Mode: "VPC"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with PUBLIC mode carrying a vpc_config", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Network.VpcConfig = &AwsBedrockAgentCoreVpcConfig{
					Subnets:        []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
					SecurityGroups: []*foreignkeyv1.StringValueOrRef{svr("sg-abc")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an unknown server protocol", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.ServerProtocol = "GRPC"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a custom claim setting both match shapes", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
					DiscoveryUrl: "https://issuer.example.com/.well-known/openid-configuration",
					CustomClaims: []*AwsBedrockAgentCoreCustomClaim{
						{
							ClaimName:     "org",
							ValueType:     "STRING",
							MatchOperator: "EQUALS",
							MatchValue:    "acme",
							MatchValues:   []string{"acme"},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a private endpoint setting both arms", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
					DiscoveryUrl: "https://issuer.example.com/.well-known/openid-configuration",
					PrivateEndpoint: &AwsBedrockAgentCorePrivateEndpoint{
						ManagedVpc: &AwsBedrockAgentCoreManagedVpcEndpoint{
							VpcId:                 svr("vpc-abc"),
							SubnetIds:             []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
							EndpointIpAddressType: "IPV4",
						},
						SelfManagedLattice: &AwsBedrockAgentCoreLatticeEndpoint{
							ResourceConfigurationId: "rcfg-abc",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a filesystem mount outside /mnt", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Filesystems = []*AwsBedrockAgentCoreRuntimeFilesystem{
					{MountPath: "/data/scratch", SessionStorage: true},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a filesystem setting two source arms", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Filesystems = []*AwsBedrockAgentCoreRuntimeFilesystem{
					{
						MountPath:         "/mnt/shared",
						SessionStorage:    true,
						EfsAccessPointArn: svr("arn:aws:elasticfilesystem:us-west-2:123456789012:access-point/fsap-abc"),
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate endpoint names", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Endpoints = []*AwsBedrockAgentCoreRuntimeEndpoint{
					{Name: "live"},
					{Name: "live"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate filesystem mount paths", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalRuntime()
				spec.Filesystems = []*AwsBedrockAgentCoreRuntimeFilesystem{
					{MountPath: "/mnt/x", SessionStorage: true},
					{MountPath: "/mnt/x", SessionStorage: true},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with six private endpoint overrides", func() {
			ginkgo.It("should return a validation error", func() {
				override := &AwsBedrockAgentCorePrivateEndpointOverride{
					Domain: "example.com",
					PrivateEndpoint: &AwsBedrockAgentCorePrivateEndpoint{
						SelfManagedLattice: &AwsBedrockAgentCoreLatticeEndpoint{ResourceConfigurationId: "rcfg-abc"},
					},
				}
				spec := minimalRuntime()
				spec.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
					DiscoveryUrl: "https://issuer.example.com/.well-known/openid-configuration",
					PrivateEndpointOverrides: []*AwsBedrockAgentCorePrivateEndpointOverride{
						override, override, override, override, override, override,
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
