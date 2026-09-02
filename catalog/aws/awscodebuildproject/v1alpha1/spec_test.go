package awscodebuildprojectv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsCodeBuildProjectSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCodeBuildProjectSpec Validation Suite")
}

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

func svRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func minimalGitHubSpec() *AwsCodeBuildProjectSpec {
	return &AwsCodeBuildProjectSpec{
		Region: "us-west-2",
		Source: &AwsCodeBuildSource{
			Type:     "GITHUB",
			Location: "https://github.com/example/repo.git",
		},
		Environment: &AwsCodeBuildEnvironment{
			Type:        "LINUX_CONTAINER",
			ComputeType: "BUILD_GENERAL1_SMALL",
			Image:       "aws/codebuild/amazonlinux2-x86_64-standard:5.0",
		},
		Artifacts: &AwsCodeBuildArtifacts{
			Type: "NO_ARTIFACTS",
		},
		ServiceRole: svRef("arn:aws:iam::123456789012:role/codebuild-role"),
	}
}

func minimalCodePipelineSpec() *AwsCodeBuildProjectSpec {
	return &AwsCodeBuildProjectSpec{
		Region: "us-west-2",
		Source: &AwsCodeBuildSource{
			Type: "CODEPIPELINE",
		},
		Environment: &AwsCodeBuildEnvironment{
			Type:        "LINUX_CONTAINER",
			ComputeType: "BUILD_GENERAL1_MEDIUM",
			Image:       "aws/codebuild/amazonlinux2-x86_64-standard:5.0",
		},
		Artifacts: &AwsCodeBuildArtifacts{
			Type: "CODEPIPELINE",
		},
		ServiceRole: svRef("arn:aws:iam::123456789012:role/codebuild-role"),
	}
}

var _ = ginkgo.Describe("AwsCodeBuildProjectSpec validations", func() {

	// =========================================================================
	// Valid configurations
	// =========================================================================

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal GitHub configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalGitHubSpec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with minimal CodePipeline configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalCodePipelineSpec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with NO_SOURCE and inline buildspec", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = &AwsCodeBuildSource{
					Type:      "NO_SOURCE",
					Buildspec: "version: 0.2\nphases:\n  build:\n    commands:\n      - echo hello",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with S3 source", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = &AwsCodeBuildSource{
					Type:     "S3",
					Location: "my-bucket/source.zip",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with S3 artifacts", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Artifacts = &AwsCodeBuildArtifacts{
					Type:     "S3",
					Location: svRef("my-artifacts-bucket"),
					Name:     "build-output",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with all optional fields populated", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsCodeBuildProjectSpec{
					Region: "us-west-2",
					Source: &AwsCodeBuildSource{
						Type:                "GITHUB",
						Location:            "https://github.com/example/repo.git",
						Buildspec:           "buildspec.yml",
						GitCloneDepth:       1,
						ReportBuildStatus:   true,
						GitSubmodulesConfig: &AwsCodeBuildGitSubmodulesConfig{FetchSubmodules: true},
						BuildStatusConfig:   &AwsCodeBuildBuildStatusConfig{Context: "ci/codebuild"},
						Auth: &AwsCodeBuildSourceAuth{
							Type:     "CODECONNECTIONS",
							Resource: "arn:aws:codeconnections:us-west-2:123456789012:connection/abc",
						},
					},
					Environment: &AwsCodeBuildEnvironment{
						Type:                     "LINUX_CONTAINER",
						ComputeType:              "BUILD_GENERAL1_LARGE",
						Image:                    "aws/codebuild/amazonlinux2-x86_64-standard:5.0",
						PrivilegedMode:           true,
						ImagePullCredentialsType: stringPtr("CODEBUILD"),
						EnvironmentVariables: []*AwsCodeBuildEnvironmentVariable{
							{Name: "ENV", Value: "production", Type: stringPtr("PLAINTEXT")},
							{Name: "DB_PASSWORD", Value: "my-secret", Type: stringPtr("SECRETS_MANAGER")},
						},
					},
					Artifacts: &AwsCodeBuildArtifacts{
						Type:          "S3",
						Location:      svRef("my-bucket"),
						Name:          "output",
						Path:          "builds",
						Packaging:     "ZIP",
						NamespaceType: "BUILD_ID",
					},
					ServiceRole:          svRef("arn:aws:iam::123456789012:role/codebuild-role"),
					Description:          "Production build project",
					EncryptionKey:        svRef("arn:aws:kms:us-east-1:123456789012:key/example"),
					BuildTimeout:         int32Ptr(120),
					QueuedTimeout:        int32Ptr(240),
					ConcurrentBuildLimit: 5,
					SourceVersion:        "main",
					Cache: &AwsCodeBuildCache{
						Type:     stringPtr("S3"),
						Location: svRef("my-cache-bucket/prefix"),
					},
					LogsConfig: &AwsCodeBuildLogsConfig{
						CloudwatchLogs: &AwsCodeBuildCloudWatchLogs{
							Status:     stringPtr("ENABLED"),
							GroupName:  svRef("/aws/codebuild/my-project"),
							StreamName: "build",
						},
						S3Logs: &AwsCodeBuildS3Logs{
							Status: stringPtr("ENABLED"),
							Bucket: svRef("my-log-bucket"),
							Prefix: "build-logs",
						},
					},
					VpcConfig: &AwsCodeBuildVpcConfig{
						VpcId:            svRef("vpc-abc123"),
						SubnetIds:        []*foreignkeyv1.StringValueOrRef{svRef("subnet-aaa"), svRef("subnet-bbb")},
						SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svRef("sg-111")},
					},
					Webhook: &AwsCodeBuildWebhook{
						BuildType: "BUILD",
						FilterGroups: []*AwsCodeBuildWebhookFilterGroup{
							{
								Filters: []*AwsCodeBuildWebhookFilter{
									{Type: "EVENT", Pattern: "PUSH, PULL_REQUEST_CREATED"},
									{Type: "HEAD_REF", Pattern: "^refs/heads/main$"},
								},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with Lambda environment type", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:        "LINUX_LAMBDA_CONTAINER",
					ComputeType: "BUILD_LAMBDA_4GB",
					Image:       "aws/codebuild/amazonlinux-aarch64-lambda-standard:go1.21",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with ARM container", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:        "ARM_CONTAINER",
					ComputeType: "BUILD_GENERAL1_LARGE",
					Image:       "aws/codebuild/amazonlinux2-aarch64-standard:3.0",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with LOCAL cache using Docker layer mode", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Cache = &AwsCodeBuildCache{
					Type:  stringPtr("LOCAL"),
					Modes: []string{"LOCAL_DOCKER_LAYER_CACHE", "LOCAL_SOURCE_CACHE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with webhook excluding a branch", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Webhook = &AwsCodeBuildWebhook{
					FilterGroups: []*AwsCodeBuildWebhookFilterGroup{
						{
							Filters: []*AwsCodeBuildWebhookFilter{
								{Type: "EVENT", Pattern: "PUSH"},
								{Type: "HEAD_REF", Pattern: "^refs/heads/release/.*$", ExcludeMatchedPattern: true},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Required field validations
	// =========================================================================

	ginkgo.Describe("When required fields are missing", func() {

		ginkgo.Context("with missing source", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with missing environment", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with missing artifacts", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Artifacts = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with missing service_role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.ServiceRole = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with missing source.type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source.Type = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with missing environment.image", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.Image = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with missing webhook filter in filter group", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Webhook = &AwsCodeBuildWebhook{
					FilterGroups: []*AwsCodeBuildWebhookFilterGroup{
						{Filters: []*AwsCodeBuildWebhookFilter{}},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Enum / string-in validations
	// =========================================================================

	ginkgo.Describe("When invalid enum values are passed", func() {

		ginkgo.Context("with invalid source.type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source.Type = "INVALID_SOURCE"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid environment.type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.Type = "INVALID_ENV"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid environment.compute_type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.ComputeType = "BUILD_MEGA_XLARGE"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid artifacts.type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Artifacts.Type = "INVALID_ARTIFACTS"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid artifacts.packaging", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Artifacts.Packaging = "TAR"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid cache.type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Cache = &AwsCodeBuildCache{
					Type: stringPtr("REDIS"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid webhook filter type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Webhook = &AwsCodeBuildWebhook{
					FilterGroups: []*AwsCodeBuildWebhookFilterGroup{
						{
							Filters: []*AwsCodeBuildWebhookFilter{
								{Type: "INVALID_FILTER", Pattern: "test"},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Range validations
	// =========================================================================

	ginkgo.Describe("When values are out of range", func() {

		ginkgo.Context("with build_timeout below minimum", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.BuildTimeout = int32Ptr(3)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with build_timeout above maximum", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.BuildTimeout = int32Ptr(3000)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with queued_timeout below minimum", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.QueuedTimeout = int32Ptr(2)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with queued_timeout above maximum", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.QueuedTimeout = int32Ptr(600)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with description exceeding max length", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				longDesc := ""
				for i := 0; i < 260; i++ {
					longDesc += "x"
				}
				spec.Description = longDesc
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Cross-field (CEL) validations
	// =========================================================================

	ginkgo.Describe("When cross-field validations fail", func() {

		ginkgo.Context("with CODEPIPELINE source but non-CODEPIPELINE artifacts", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCodePipelineSpec()
				spec.Artifacts = &AwsCodeBuildArtifacts{Type: "NO_ARTIFACTS"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with non-CODEPIPELINE source but CODEPIPELINE artifacts", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Artifacts = &AwsCodeBuildArtifacts{Type: "CODEPIPELINE"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with GITHUB source but missing location", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source.Location = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with BITBUCKET source but missing location", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source.Type = "BITBUCKET"
				spec.Source.Location = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with NO_SOURCE but missing buildspec", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = &AwsCodeBuildSource{
					Type: "NO_SOURCE",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with S3 artifacts but missing location", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Artifacts = &AwsCodeBuildArtifacts{Type: "S3"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with webhook on CODEPIPELINE source", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCodePipelineSpec()
				spec.Webhook = &AwsCodeBuildWebhook{
					FilterGroups: []*AwsCodeBuildWebhookFilterGroup{
						{
							Filters: []*AwsCodeBuildWebhookFilter{
								{Type: "EVENT", Pattern: "PUSH"},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with webhook on NO_SOURCE", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = &AwsCodeBuildSource{
					Type:      "NO_SOURCE",
					Buildspec: "version: 0.2",
				}
				spec.Webhook = &AwsCodeBuildWebhook{
					FilterGroups: []*AwsCodeBuildWebhookFilterGroup{
						{
							Filters: []*AwsCodeBuildWebhookFilter{
								{Type: "EVENT", Pattern: "PUSH"},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with webhook on S3 source", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = &AwsCodeBuildSource{
					Type:     "S3",
					Location: "my-bucket/source.zip",
				}
				spec.Webhook = &AwsCodeBuildWebhook{
					FilterGroups: []*AwsCodeBuildWebhookFilterGroup{
						{
							Filters: []*AwsCodeBuildWebhookFilter{
								{Type: "EVENT", Pattern: "PUSH"},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// VPC config validations
	// =========================================================================

	ginkgo.Describe("When VPC config is incomplete", func() {

		ginkgo.Context("with vpc_config missing subnet_ids", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.VpcConfig = &AwsCodeBuildVpcConfig{
					VpcId:            svRef("vpc-abc123"),
					SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svRef("sg-111")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with vpc_config missing security_group_ids", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.VpcConfig = &AwsCodeBuildVpcConfig{
					VpcId:     svRef("vpc-abc123"),
					SubnetIds: []*foreignkeyv1.StringValueOrRef{svRef("subnet-aaa")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with vpc_config exceeding max subnets", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				var subnets []*foreignkeyv1.StringValueOrRef
				for i := 0; i < 17; i++ {
					subnets = append(subnets, svRef("subnet-xxx"))
				}
				spec.VpcConfig = &AwsCodeBuildVpcConfig{
					VpcId:            svRef("vpc-abc123"),
					SubnetIds:        subnets,
					SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svRef("sg-111")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with vpc_config exceeding max security groups", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				var sgs []*foreignkeyv1.StringValueOrRef
				for i := 0; i < 6; i++ {
					sgs = append(sgs, svRef("sg-xxx"))
				}
				spec.VpcConfig = &AwsCodeBuildVpcConfig{
					VpcId:            svRef("vpc-abc123"),
					SubnetIds:        []*foreignkeyv1.StringValueOrRef{svRef("subnet-aaa")},
					SecurityGroupIds: sgs,
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Secondary sources / artifacts identifier contract
	// =========================================================================

	ginkgo.Describe("Secondary source and artifact identifiers", func() {

		ginkgo.Context("with an identifier on the primary source", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source.SourceIdentifier = "primary"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a secondary source missing its identifier", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.SecondarySources = []*AwsCodeBuildSource{
					{Type: "GITHUB", Location: "https://github.com/example/tooling.git"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a named secondary source and version pin", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.SecondarySources = []*AwsCodeBuildSource{
					{Type: "GITHUB", Location: "https://github.com/example/tooling.git", SourceIdentifier: "tooling"},
				}
				spec.SecondarySourceVersions = []*AwsCodeBuildSecondarySourceVersion{
					{SourceIdentifier: "tooling", SourceVersion: "main"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a secondary artifact missing its identifier", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.SecondaryArtifacts = []*AwsCodeBuildArtifacts{
					{Type: "S3", Location: svRef("reports-bucket")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a named secondary artifact", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.SecondaryArtifacts = []*AwsCodeBuildArtifacts{
					{Type: "S3", Location: svRef("reports-bucket"), ArtifactIdentifier: "reports"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Lambda environment constraints
	// =========================================================================

	ginkgo.Describe("Lambda environment constraints", func() {

		ginkgo.Context("with an explicit build_timeout on a Lambda environment", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:        "LINUX_LAMBDA_CONTAINER",
					ComputeType: "BUILD_LAMBDA_4GB",
					Image:       "aws/codebuild/amazonlinux-x86_64-lambda-standard:go1.21",
				}
				spec.BuildTimeout = int32Ptr(30)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an explicit queued_timeout on a Lambda environment", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:        "ARM_LAMBDA_CONTAINER",
					ComputeType: "BUILD_LAMBDA_2GB",
					Image:       "aws/codebuild/amazonlinux-aarch64-lambda-standard:go1.21",
				}
				spec.QueuedTimeout = int32Ptr(60)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with privileged_mode on a Lambda environment", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:           "LINUX_LAMBDA_CONTAINER",
					ComputeType:    "BUILD_LAMBDA_1GB",
					Image:          "aws/codebuild/amazonlinux-x86_64-lambda-standard:go1.21",
					PrivilegedMode: true,
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a Lambda environment and no timeouts", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:        "LINUX_LAMBDA_CONTAINER",
					ComputeType: "BUILD_LAMBDA_4GB",
					Image:       "aws/codebuild/amazonlinux-x86_64-lambda-standard:go1.21",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Environment depth (docker server, fleet, registry credential)
	// =========================================================================

	ginkgo.Describe("Environment depth", func() {

		ginkgo.Context("with a docker server and security groups", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.DockerServer = &AwsCodeBuildDockerServer{
					ComputeType:      "BUILD_GENERAL1_MEDIUM",
					SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svRef("sg-111")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid docker server compute type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.DockerServer = &AwsCodeBuildDockerServer{
					ComputeType: "BUILD_LAMBDA_1GB",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a fleet ARN and fleet compute type", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:        "LINUX_EC2",
					ComputeType: "ATTRIBUTE_BASED_COMPUTE",
					Image:       "aws/codebuild/amazonlinux2-x86_64-standard:5.0",
					FleetArn:    "arn:aws:codebuild:us-west-2:123456789012:fleet/my-fleet",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a fleet_arn that is not a CodeBuild fleet ARN", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.FleetArn = "my-fleet"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a host kernel on a Linux container environment", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.HostKernel = "LINUX_KERNEL_6"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a host kernel on an EC2 environment", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:        "ARM_EC2",
					ComputeType: "ATTRIBUTE_BASED_COMPUTE",
					Image:       "aws/codebuild/amazonlinux2-aarch64-standard:3.0",
					FleetArn:    "arn:aws:codebuild:us-west-2:123456789012:fleet/my-fleet",
					HostKernel:  "LINUX_KERNEL_LATEST",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a host kernel on a Windows environment", func() {
			ginkgo.It("should return a validation error (kernel selection is Linux container/EC2 only)", func() {
				spec := minimalGitHubSpec()
				spec.Environment = &AwsCodeBuildEnvironment{
					Type:        "WINDOWS_SERVER_2022_CONTAINER",
					ComputeType: "BUILD_GENERAL1_MEDIUM",
					Image:       "aws/codebuild/windows-base:2022-1.0",
					HostKernel:  "LINUX_KERNEL_6",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid host kernel value", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.HostKernel = "LINUX_KERNEL_5"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a registry credential without SERVICE_ROLE pulls", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.RegistryCredential = &AwsCodeBuildRegistryCredential{
					Credential:         "arn:aws:secretsmanager:us-west-2:123456789012:secret:dockerhub",
					CredentialProvider: "SECRETS_MANAGER",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a registry credential and SERVICE_ROLE pulls", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.ImagePullCredentialsType = stringPtr("SERVICE_ROLE")
				spec.Environment.RegistryCredential = &AwsCodeBuildRegistryCredential{
					Credential:         "arn:aws:secretsmanager:us-west-2:123456789012:secret:dockerhub",
					CredentialProvider: "SECRETS_MANAGER",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an environment certificate not ending in .pem or .zip", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Environment.Certificate = "certs-bucket/private-ca.crt"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Visibility, badge, and source-type couplings
	// =========================================================================

	ginkgo.Describe("Visibility, badge, and source couplings", func() {

		ginkgo.Context("with PUBLIC_READ visibility but no resource access role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.ProjectVisibility = stringPtr("PUBLIC_READ")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with PUBLIC_READ visibility and a resource access role", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.ProjectVisibility = stringPtr("PUBLIC_READ")
				spec.ResourceAccessRole = svRef("arn:aws:iam::123456789012:role/codebuild-public-access")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a badge on a NO_SOURCE project", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = &AwsCodeBuildSource{Type: "NO_SOURCE", Buildspec: "version: 0.2"}
				spec.BadgeEnabled = true
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with git submodules on a GITLAB source", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = &AwsCodeBuildSource{
					Type:                "GITLAB",
					Location:            "https://gitlab.com/example/repo.git",
					GitSubmodulesConfig: &AwsCodeBuildGitSubmodulesConfig{FetchSubmodules: true},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with build status reporting on a CODECOMMIT source", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Source = &AwsCodeBuildSource{
					Type:              "CODECOMMIT",
					Location:          "https://git-codecommit.us-west-2.amazonaws.com/v1/repos/example",
					ReportBuildStatus: true,
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid auto_retry_limit", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.AutoRetryLimit = 11
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Batch config and webhook depth
	// =========================================================================

	ginkgo.Describe("Batch config and webhook depth", func() {

		ginkgo.Context("with a full batch config", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.BuildBatchConfig = &AwsCodeBuildBatchConfig{
					ServiceRole:      svRef("arn:aws:iam::123456789012:role/codebuild-batch"),
					CombineArtifacts: true,
					TimeoutInMins:    120,
					Restrictions: &AwsCodeBuildBatchRestrictions{
						ComputeTypesAllowed:  []string{"BUILD_GENERAL1_SMALL", "BUILD_GENERAL1_MEDIUM"},
						MaximumBuildsAllowed: 10,
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a batch config missing its service role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.BuildBatchConfig = &AwsCodeBuildBatchConfig{}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid batch restriction compute type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.BuildBatchConfig = &AwsCodeBuildBatchConfig{
					ServiceRole: svRef("arn:aws:iam::123456789012:role/codebuild-batch"),
					Restrictions: &AwsCodeBuildBatchRestrictions{
						ComputeTypesAllowed: []string{"BUILD_MEGA"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a scoped webhook and PR build policy", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Webhook = &AwsCodeBuildWebhook{
					BuildType:      "RUNNER_BUILDKITE_BUILD",
					ManualCreation: true,
					ScopeConfiguration: &AwsCodeBuildWebhookScopeConfiguration{
						Name:  "my-org",
						Scope: "GITHUB_ORGANIZATION",
					},
					PullRequestBuildPolicy: &AwsCodeBuildWebhookPullRequestBuildPolicy{
						RequiresCommentApproval: "FORK_PULL_REQUESTS",
						ApproverRoles:           []string{"GITHUB_WRITE", "GITHUB_ADMIN"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid webhook scope", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Webhook = &AwsCodeBuildWebhook{
					ScopeConfiguration: &AwsCodeBuildWebhookScopeConfiguration{
						Name:  "my-org",
						Scope: "BITBUCKET_WORKSPACE",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid PR approver role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Webhook = &AwsCodeBuildWebhook{
					PullRequestBuildPolicy: &AwsCodeBuildWebhookPullRequestBuildPolicy{
						RequiresCommentApproval: "ALL_PULL_REQUESTS",
						ApproverRoles:           []string{"GITHUB_SUPERUSER"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with the new org-scoped webhook filter types", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Webhook = &AwsCodeBuildWebhook{
					FilterGroups: []*AwsCodeBuildWebhookFilterGroup{
						{
							Filters: []*AwsCodeBuildWebhookFilter{
								{Type: "EVENT", Pattern: "WORKFLOW_JOB_QUEUED"},
								{Type: "REPOSITORY_NAME", Pattern: "^service-.*$"},
								{Type: "WORKFLOW_NAME", Pattern: "^ci$"},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid LOCAL cache mode", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.Cache = &AwsCodeBuildCache{
					Type:  stringPtr("LOCAL"),
					Modes: []string{"LOCAL_EVERYTHING_CACHE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a file system location", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.VpcConfig = &AwsCodeBuildVpcConfig{
					VpcId:            svRef("vpc-abc123"),
					SubnetIds:        []*foreignkeyv1.StringValueOrRef{svRef("subnet-aaa")},
					SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svRef("sg-111")},
				}
				spec.FileSystemLocations = []*AwsCodeBuildFileSystemLocation{
					{
						Identifier: "build_cache",
						Location:   "fs-0abc.efs.us-west-2.amazonaws.com:/cache",
						MountPoint: "/mnt/cache",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid file system type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.FileSystemLocations = []*AwsCodeBuildFileSystemLocation{
					{
						Identifier: "cache",
						Location:   "fs-0abc.efs.us-west-2.amazonaws.com:/cache",
						MountPoint: "/mnt/cache",
						Type:       stringPtr("FSX"),
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// S3 build logs: AWS stores them only under "bucket/prefix", so an
	// ENABLED destination needs both halves (the modules compose the
	// provider's single location argument from them).
	// =========================================================================

	ginkgo.Describe("S3 logs bucket and prefix", func() {

		ginkgo.Context("when ENABLED with bucket and prefix", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.LogsConfig = &AwsCodeBuildLogsConfig{
					S3Logs: &AwsCodeBuildS3Logs{
						Status: stringPtr("ENABLED"),
						Bucket: svRef("my-log-bucket"),
						Prefix: "build-logs",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when ENABLED without a prefix", func() {
			ginkgo.It("should return a validation error -- AWS rejects a bare bucket", func() {
				spec := minimalGitHubSpec()
				spec.LogsConfig = &AwsCodeBuildLogsConfig{
					S3Logs: &AwsCodeBuildS3Logs{
						Status: stringPtr("ENABLED"),
						Bucket: svRef("my-log-bucket"),
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when ENABLED without a bucket", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.LogsConfig = &AwsCodeBuildLogsConfig{
					S3Logs: &AwsCodeBuildS3Logs{
						Status: stringPtr("ENABLED"),
						Prefix: "build-logs",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when DISABLED with neither half", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGitHubSpec()
				spec.LogsConfig = &AwsCodeBuildLogsConfig{
					S3Logs: &AwsCodeBuildS3Logs{
						Status: stringPtr("DISABLED"),
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
