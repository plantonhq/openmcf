package awscodepipelinev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsCodePipelineSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCodePipelineSpec Validation Suite")
}

func stringPtr(s string) *string {
	return &s
}

func svRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalV2Spec returns a minimal valid V2 pipeline spec with a CodeStar
// source stage and a CodeBuild build stage.
func minimalV2Spec() *AwsCodePipelineSpec {
	return &AwsCodePipelineSpec{
		Region:  "us-west-2",
		RoleArn: svRef("arn:aws:iam::123456789012:role/codepipeline-role"),
		ArtifactStores: []*AwsCodePipelineArtifactStore{
			{Location: svRef("my-artifact-bucket")},
		},
		Stages: []*AwsCodePipelineStage{
			{
				Name: "Source",
				Actions: []*AwsCodePipelineAction{
					{
						Name:            "SourceAction",
						Category:        "Source",
						Owner:           "AWS",
						Provider:        "CodeStarSourceConnection",
						Version:         "1",
						OutputArtifacts: []string{"SourceOutput"},
						Configuration: map[string]string{
							"ConnectionArn":    "arn:aws:codestar-connections:us-east-1:123456789012:connection/abc",
							"FullRepositoryId": "my-org/my-repo",
							"BranchName":       "main",
						},
					},
				},
			},
			{
				Name: "Build",
				Actions: []*AwsCodePipelineAction{
					{
						Name:            "BuildAction",
						Category:        "Build",
						Owner:           "AWS",
						Provider:        "CodeBuild",
						Version:         "1",
						InputArtifacts:  []string{"SourceOutput"},
						OutputArtifacts: []string{"BuildOutput"},
						Configuration: map[string]string{
							"ProjectName": "my-build-project",
						},
					},
				},
			},
		},
	}
}

// minimalV1Spec returns a minimal valid V1 pipeline spec.
func minimalV1Spec() *AwsCodePipelineSpec {
	spec := minimalV2Spec()
	spec.PipelineType = stringPtr("V1")
	return spec
}

var _ = ginkgo.Describe("AwsCodePipelineSpec validations", func() {

	// =========================================================================
	// Valid configurations
	// =========================================================================

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal V2 pipeline", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalV2Spec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with minimal V1 pipeline", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalV1Spec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with explicit pipeline_type V2 and execution_mode SUPERSEDED", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.PipelineType = stringPtr("V2")
				spec.ExecutionMode = stringPtr("SUPERSEDED")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with V2 pipeline and QUEUED execution mode", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.PipelineType = stringPtr("V2")
				spec.ExecutionMode = stringPtr("QUEUED")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with V2 pipeline and PARALLEL execution mode", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.PipelineType = stringPtr("V2")
				spec.ExecutionMode = stringPtr("PARALLEL")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with cross-region artifact stores", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.ArtifactStores = []*AwsCodePipelineArtifactStore{
					{Location: svRef("bucket-us-east-1"), Region: "us-east-1"},
					{Location: svRef("bucket-eu-west-1"), Region: "eu-west-1"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with artifact store encryption key", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.ArtifactStores[0].EncryptionKeyId = svRef("arn:aws:kms:us-east-1:123456789012:key/abc")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with triggers on V2 pipeline", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
							Push: []*AwsCodePipelineGitPush{
								{
									Branches: &AwsCodePipelineGitFilter{
										Includes: []string{"main", "release/*"},
									},
								},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with variables on V2 pipeline", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Variables = []*AwsCodePipelineVariable{
					{Name: "DEPLOY_ENV", DefaultValue: "staging", Description: "Target deployment environment"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with action run_order", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].RunOrder = 1
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a timeout on a Manual Approval action", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions = append(spec.Stages[1].Actions, &AwsCodePipelineAction{
					Name:             "ApproveRelease",
					Category:         "Approval",
					Owner:            "AWS",
					Provider:         "Manual",
					Version:          "1",
					RunOrder:         2,
					TimeoutInMinutes: 2880,
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a timeout on a Build action", func() {
			ginkgo.It("should return a validation error (AWS supports the override only on Approval actions)", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].TimeoutInMinutes = 60
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a Compute action carrying commands, output variables, and compute artifacts", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions = append(spec.Stages[1].Actions, &AwsCodePipelineAction{
					Name:            "InlineChecks",
					Category:        "Compute",
					Owner:           "AWS",
					Provider:        "Commands",
					Version:         "1",
					RunOrder:        2,
					InputArtifacts:  []string{"SourceOutput"},
					Commands:        []string{"npm ci", "npm test"},
					Namespace:       "ComputeVars",
					OutputVariables: []string{"TEST_REPORT_URL"},
					OutputArtifactsForComputeAction: []*AwsCodePipelineComputeOutputArtifact{
						{Name: "TestReports", Files: []string{"reports/junit.xml"}},
					},
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with commands on a Build action (AWS ignores them outside Compute; no invented restriction)", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].Commands = []string{"echo build-side command"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with action role_arn for cross-account", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].RoleArn = svRef("arn:aws:iam::987654321098:role/cross-account-deploy")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with pull request trigger", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
							PullRequest: []*AwsCodePipelineGitPullRequest{
								{
									Branches: &AwsCodePipelineGitFilter{
										Includes: []string{"main"},
									},
									Events: []string{"OPEN", "UPDATED"},
								},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with multi-stage pipeline (source, build, approval, deploy)", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages = append(spec.Stages,
					&AwsCodePipelineStage{
						Name: "Approval",
						Actions: []*AwsCodePipelineAction{
							{
								Name:     "ManualApproval",
								Category: "Approval",
								Owner:    "AWS",
								Provider: "Manual",
								Version:  "1",
								Configuration: map[string]string{
									"CustomData": "Please review the build output",
								},
							},
						},
					},
					&AwsCodePipelineStage{
						Name: "Deploy",
						Actions: []*AwsCodePipelineAction{
							{
								Name:           "DeployToECS",
								Category:       "Deploy",
								Owner:          "AWS",
								Provider:       "ECS",
								Version:        "1",
								InputArtifacts: []string{"BuildOutput"},
								Configuration: map[string]string{
									"ClusterName": "my-cluster",
									"ServiceName": "my-service",
									"FileName":    "imagedefinitions.json",
								},
							},
						},
					},
				)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Required field validations
	// =========================================================================

	ginkgo.Describe("When required fields are missing", func() {

		ginkgo.Context("when role_arn is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.RoleArn = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when artifact_stores is empty", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.ArtifactStores = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when stages is empty", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when stages has only one stage", func() {
			ginkgo.It("should return a validation error for min_items 2", func() {
				spec := minimalV2Spec()
				spec.Stages = spec.Stages[:1]
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when artifact_store location is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.ArtifactStores[0].Location = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when stage name is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Name = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when stage actions is empty", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Actions = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when action name is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Actions[0].Name = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when action category is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Actions[0].Category = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when action owner is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Actions[0].Owner = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when action provider is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Actions[0].Provider = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when action version is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Actions[0].Version = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when trigger provider_type is wrong", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "InvalidProvider",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when trigger git_configuration is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when trigger source_action_name is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("when variable name is missing", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Variables = []*AwsCodePipelineVariable{
					{Name: ""},
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

		ginkgo.Context("with invalid pipeline_type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.PipelineType = stringPtr("V3")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid execution_mode", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.ExecutionMode = stringPtr("INVALID")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid action category", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Actions[0].Category = "InvalidCategory"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid action owner", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[0].Actions[0].Owner = "Nobody"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Range validations
	// =========================================================================

	ginkgo.Describe("When values are out of range", func() {

		ginkgo.Context("with stage name exceeding 100 characters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				longName := ""
				for i := 0; i < 101; i++ {
					longName += "a"
				}
				spec.Stages[0].Name = longName
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with action run_order below 1", func() {
			ginkgo.It("should pass because zero is ignored", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].RunOrder = 0
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with action run_order above 999", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].RunOrder = 1000
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with action timeout_in_minutes below 5", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].TimeoutInMinutes = 3
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with action timeout_in_minutes above 86400", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].TimeoutInMinutes = 86401
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with action namespace exceeding 100 characters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				longNs := ""
				for i := 0; i < 101; i++ {
					longNs += "x"
				}
				spec.Stages[0].Actions[0].Namespace = longNs
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with git filter includes exceeding max 8 patterns", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
							Push: []*AwsCodePipelineGitPush{
								{
									Branches: &AwsCodePipelineGitFilter{
										Includes: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
									},
								},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with more than 3 push filters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
							Push: []*AwsCodePipelineGitPush{
								{Branches: &AwsCodePipelineGitFilter{Includes: []string{"a"}}},
								{Branches: &AwsCodePipelineGitFilter{Includes: []string{"b"}}},
								{Branches: &AwsCodePipelineGitFilter{Includes: []string{"c"}}},
								{Branches: &AwsCodePipelineGitFilter{Includes: []string{"d"}}},
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
	// Cross-field (CEL) validations
	// =========================================================================

	ginkgo.Describe("When cross-field constraints are violated", func() {

		ginkgo.Context("with triggers on a V1 pipeline", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV1Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
							Push: []*AwsCodePipelineGitPush{
								{Branches: &AwsCodePipelineGitFilter{Includes: []string{"main"}}},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with variables on a V1 pipeline", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV1Spec()
				spec.Variables = []*AwsCodePipelineVariable{
					{Name: "MY_VAR"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with QUEUED execution mode on a V1 pipeline", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV1Spec()
				spec.ExecutionMode = stringPtr("QUEUED")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with PARALLEL execution mode on a V1 pipeline", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV1Spec()
				spec.ExecutionMode = stringPtr("PARALLEL")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with SUPERSEDED execution mode on a V1 pipeline", func() {
			ginkgo.It("should not return a validation error (SUPERSEDED is valid for V1)", func() {
				spec := minimalV1Spec()
				spec.ExecutionMode = stringPtr("SUPERSEDED")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with plain output_artifacts on a Compute action", func() {
			ginkgo.It("should return a validation error (Compute actions export via output_artifacts_for_compute_action)", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions = append(spec.Stages[1].Actions, &AwsCodePipelineAction{
					Name:            "InlineChecks",
					Category:        "Compute",
					Owner:           "AWS",
					Provider:        "Commands",
					Version:         "1",
					RunOrder:        2,
					Commands:        []string{"make verify"},
					OutputArtifacts: []string{"ComputeOutput"},
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with output_artifacts_for_compute_action on a Build action", func() {
			ginkgo.It("should return a validation error (file-based artifacts are Compute-only)", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions[0].OutputArtifactsForComputeAction = []*AwsCodePipelineComputeOutputArtifact{
					{Name: "BuildFiles", Files: []string{"dist/app.zip"}},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a compute output artifact name carrying illegal characters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].Actions = append(spec.Stages[1].Actions, &AwsCodePipelineAction{
					Name:     "InlineChecks",
					Category: "Compute",
					Owner:    "AWS",
					Provider: "Commands",
					Version:  "1",
					RunOrder: 2,
					Commands: []string{"make verify"},
					OutputArtifactsForComputeAction: []*AwsCodePipelineComputeOutputArtifact{
						{Name: "bad name!", Files: []string{"out.txt"}},
					},
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with more than 15 output variables on a Compute action", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				vars := make([]string, 16)
				for i := range vars {
					vars[i] = "VAR_" + string(rune('A'+i))
				}
				spec.Stages[1].Actions = append(spec.Stages[1].Actions, &AwsCodePipelineAction{
					Name:            "InlineChecks",
					Category:        "Compute",
					Owner:           "AWS",
					Provider:        "Commands",
					Version:         "1",
					RunOrder:        2,
					Commands:        []string{"make verify"},
					OutputVariables: vars,
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an action configuration key longer than 50 characters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				longKey := ""
				for range 51 {
					longKey += "k"
				}
				spec.Stages[1].Actions[0].Configuration[longKey] = "v"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Stage conditions (V2 entry gates, success checks, failure handling)
	// =========================================================================

	ginkgo.Describe("Stage conditions", func() {

		deploymentWindowRule := func() *AwsCodePipelineRule {
			return &AwsCodePipelineRule{
				Name:       "BusinessHoursOnly",
				RuleTypeId: &AwsCodePipelineRuleTypeId{Provider: "DeploymentWindow"},
				Configuration: map[string]string{
					"Cron":     "* 9-17 ? * MON-FRI *",
					"TimeZone": "America/Los_Angeles",
				},
			}
		}

		ginkgo.Context("with a before_entry deployment window gate", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].BeforeEntry = &AwsCodePipelineStageCondition{
					Rules: []*AwsCodePipelineRule{deploymentWindowRule()},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an on_success CloudWatch alarm check", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].OnSuccess = &AwsCodePipelineStageCondition{
					Rules: []*AwsCodePipelineRule{
						{
							Name:       "VerifyNoAlarm",
							RuleTypeId: &AwsCodePipelineRuleTypeId{Provider: "CloudWatchAlarm"},
							Configuration: map[string]string{
								"AlarmName": "deploy-errors",
								"WaitTime":  "300",
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an on_failure automatic rollback", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].OnFailure = &AwsCodePipelineFailureHandling{
					Result: "ROLLBACK",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an on_failure retry and retry configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].OnFailure = &AwsCodePipelineFailureHandling{
					Result:             "RETRY",
					RetryConfiguration: &AwsCodePipelineRetryConfiguration{RetryMode: "FAILED_ACTIONS"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a retry configuration but a non-RETRY result", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].OnFailure = &AwsCodePipelineFailureHandling{
					Result:             "ROLLBACK",
					RetryConfiguration: &AwsCodePipelineRetryConfiguration{RetryMode: "ALL_ACTIONS"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a stage condition on a V1 pipeline", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV1Spec()
				spec.Stages[1].BeforeEntry = &AwsCodePipelineStageCondition{
					Rules: []*AwsCodePipelineRule{deploymentWindowRule()},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a condition that has no rules", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].BeforeEntry = &AwsCodePipelineStageCondition{}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with more than five rules in a condition", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				var rules []*AwsCodePipelineRule
				for i := 0; i < 6; i++ {
					rules = append(rules, deploymentWindowRule())
				}
				spec.Stages[1].BeforeEntry = &AwsCodePipelineStageCondition{Rules: rules}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid rule_type_id category", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				rule := deploymentWindowRule()
				rule.RuleTypeId.Category = stringPtr("Action")
				spec.Stages[1].BeforeEntry = &AwsCodePipelineStageCondition{
					Rules: []*AwsCodePipelineRule{rule},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a Commands rule carrying shell commands", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.Stages[1].OnSuccess = &AwsCodePipelineStageCondition{
					Rules: []*AwsCodePipelineRule{
						{
							Name:             "SmokeTest",
							RuleTypeId:       &AwsCodePipelineRuleTypeId{Provider: "Commands"},
							Commands:         []string{"curl -sf https://example.com/health"},
							InputArtifacts:   []string{"BuildOutput"},
							TimeoutInMinutes: 10,
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an out-of-range rule timeout", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				rule := deploymentWindowRule()
				rule.TimeoutInMinutes = 3
				spec.Stages[1].OnSuccess = &AwsCodePipelineStageCondition{
					Rules: []*AwsCodePipelineRule{rule},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a rule input artifact name carrying illegal characters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				rule := deploymentWindowRule()
				rule.InputArtifacts = []string{"has spaces"}
				spec.Stages[1].OnSuccess = &AwsCodePipelineStageCondition{
					Rules: []*AwsCodePipelineRule{rule},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a rule configuration value longer than 10000 characters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				rule := deploymentWindowRule()
				longValue := make([]byte, 10001)
				for i := range longValue {
					longValue[i] = 'v'
				}
				rule.Configuration = map[string]string{"Cron": string(longValue)}
				spec.Stages[1].OnSuccess = &AwsCodePipelineStageCondition{
					Rules: []*AwsCodePipelineRule{rule},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	// =========================================================================
	// Trigger contract honesty
	// =========================================================================

	ginkgo.Describe("Trigger contract", func() {

		ginkgo.Context("with an invalid pull-request event", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
							PullRequest: []*AwsCodePipelineGitPullRequest{
								{Events: []string{"UPDATE"}},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a git configuration carrying no filters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an empty git filter block", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.Triggers = []*AwsCodePipelineTrigger{
					{
						ProviderType: "CodeStarSourceConnection",
						GitConfiguration: &AwsCodePipelineGitConfiguration{
							SourceActionName: "SourceAction",
							Push: []*AwsCodePipelineGitPush{
								{Branches: &AwsCodePipelineGitFilter{}},
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
	// Artifact store region shape
	// =========================================================================

	ginkgo.Describe("Artifact store region shape", func() {

		ginkgo.Context("with multiple stores each carrying a region", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalV2Spec()
				spec.ArtifactStores = []*AwsCodePipelineArtifactStore{
					{Location: svRef("artifacts-usw2"), Region: "us-west-2"},
					{Location: svRef("artifacts-use1"), Region: "us-east-1"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with multiple stores where one is missing a region", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalV2Spec()
				spec.ArtifactStores = []*AwsCodePipelineArtifactStore{
					{Location: svRef("artifacts-usw2"), Region: "us-west-2"},
					{Location: svRef("artifacts-use1")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
