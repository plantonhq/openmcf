package awsbedrockagentcoreevaluationv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockAgentCoreEvaluationSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockAgentCoreEvaluationSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func f64(v float64) *float64 { return &v }

func boolPtr(v bool) *bool { return &v }

// minimalEvaluator is the smallest valid bundle: one code-based
// evaluator.
func minimalEvaluator() *AwsBedrockAgentCoreEvaluationSpec {
	return &AwsBedrockAgentCoreEvaluationSpec{
		Region: "us-west-2",
		Evaluators: []*AwsBedrockAgentCoreEvaluator{
			{
				Name:  "tone_check",
				Level: "TRACE",
				CodeBased: &AwsBedrockAgentCoreCodeEvaluator{
					LambdaArn: svr("arn:aws:lambda:us-west-2:123456789012:function:tone-scorer"),
				},
			},
		},
	}
}

// llmJudgeEvaluator exercises the LLM-judge arm with a categorical
// scale.
func llmJudgeEvaluator() *AwsBedrockAgentCoreEvaluator {
	return &AwsBedrockAgentCoreEvaluator{
		Name:  "helpfulness",
		Level: "SESSION",
		LlmAsAJudge: &AwsBedrockAgentCoreLlmJudge{
			Instructions: "Assess whether the agent resolved the user's request given the {context}.",
			Model: &AwsBedrockAgentCoreJudgeModel{
				ModelId: "us.amazon.nova-2-lite-v1:0",
				Inference: &AwsBedrockAgentCoreJudgeInference{
					Temperature: f64(0),
				},
			},
			RatingScale: &AwsBedrockAgentCoreRatingScale{
				Categorical: []*AwsBedrockAgentCoreCategoricalRating{
					{Label: "helpful", Definition: "the request was resolved"},
					{Label: "unhelpful", Definition: "the request was not resolved"},
				},
			},
		},
	}
}

// minimalHarness is a valid harness with the required model arm.
func minimalHarness() *AwsBedrockAgentCoreHarness {
	return &AwsBedrockAgentCoreHarness{
		Name:             "support_bench",
		ExecutionRoleArn: svr("arn:aws:iam::123456789012:role/agentcore-eval"),
		Model: &AwsBedrockAgentCoreHarnessModel{
			Bedrock: &AwsBedrockAgentCoreHarnessBedrockModel{
				ModelId: "anthropic.claude-3-5-haiku-20241022-v1:0",
			},
		},
	}
}

// minimalOnlineConfig is a valid online evaluation config.
func minimalOnlineConfig() *AwsBedrockAgentCoreOnlineEvaluationConfig {
	return &AwsBedrockAgentCoreOnlineEvaluationConfig{
		Name:             "prod_sampling",
		ExecutionRoleArn: svr("arn:aws:iam::123456789012:role/agentcore-online-eval"),
		DataSource: &AwsBedrockAgentCoreOnlineEvaluationDataSource{
			LogGroupNames: []*foreignkeyv1.StringValueOrRef{svr("/aws/bedrock-agentcore/runtimes/support")},
			ServiceNames:  []string{"support-agent"},
		},
		EvaluatorIds: []string{"Builtin.Helpfulness"},
		Rule: &AwsBedrockAgentCoreOnlineEvaluationRule{
			SamplingPercentage: 5,
		},
	}
}

var _ = ginkgo.Describe("AwsBedrockAgentCoreEvaluationSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with a minimal code-based evaluator", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalEvaluator())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEvaluator()
				spec.Evaluators = append(spec.Evaluators, llmJudgeEvaluator())
				harness := minimalHarness()
				harness.SystemPrompts = []*AwsBedrockAgentCoreHarnessSystemPrompt{
					{Text: "You are the support agent under test."},
				}
				harness.Tools = []*AwsBedrockAgentCoreHarnessTool{
					{
						Name: "kb-gateway",
						Type: "agentcore_gateway",
						AgentcoreGateway: &AwsBedrockAgentCoreHarnessGatewayTool{
							GatewayArn: svr("arn:aws:bedrock-agentcore:us-west-2:123456789012:gateway/kb-gw-123"),
							OutboundAuth: &AwsBedrockAgentCoreHarnessGatewayOutboundAuth{
								AwsIam: true,
							},
						},
					},
					{
						Name:                     "sandbox",
						Type:                     "agentcore_code_interpreter",
						AgentcoreCodeInterpreter: &AwsBedrockAgentCoreHarnessCodeInterpreterTool{},
					},
					{
						Name: "docs-mcp",
						Type: "remote_mcp",
						RemoteMcp: &AwsBedrockAgentCoreHarnessRemoteMcpTool{
							Url:     "https://mcp.internal.example.com/sse",
							Headers: map[string]string{"Authorization": "Bearer token"},
						},
					},
				}
				harness.AllowedTools = []string{"kb-gateway", "docs-mcp"}
				harness.Memory = &AwsBedrockAgentCoreHarnessMemory{
					MemoryArn: svr("arn:aws:bedrock-agentcore:us-west-2:123456789012:memory/mem-123"),
					Retrieval: &AwsBedrockAgentCoreHarnessMemoryRetrieval{
						Namespace:      "/strategies/summary/actors/bench",
						RelevanceScore: f64(0.4),
						TopK:           5,
					},
				}
				harness.RuntimeEnvironment = &AwsBedrockAgentCoreHarnessRuntimeEnvironment{
					Network: &AwsBedrockAgentCoreHarnessNetwork{
						Mode: "VPC",
						VpcConfig: &AwsBedrockAgentCoreHarnessVpcConfig{
							Subnets:        []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
							SecurityGroups: []*foreignkeyv1.StringValueOrRef{svr("sg-abc")},
						},
					},
					Filesystems: []*AwsBedrockAgentCoreHarnessFilesystem{
						{MountPath: "/mnt/scratch", SessionStorage: true},
					},
					Lifecycle: &AwsBedrockAgentCoreHarnessLifecycle{
						IdleRuntimeSessionTimeoutSeconds: 900,
					},
				}
				harness.Truncation = &AwsBedrockAgentCoreHarnessTruncation{
					Strategy: "sliding_window",
					SlidingWindow: &AwsBedrockAgentCoreHarnessSlidingWindow{
						MessagesCount: 40,
					},
				}
				harness.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
					DiscoveryUrl:    "https://accounts.google.com/.well-known/openid-configuration",
					AllowedAudience: []string{"eval-bench"},
					CustomClaims: []*AwsBedrockAgentCoreCustomClaim{
						{
							ClaimName:     "org",
							ValueType:     "STRING",
							MatchOperator: "EQUALS",
							MatchValue:    "acme",
						},
					},
				}
				spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
				config := minimalOnlineConfig()
				config.Enabled = boolPtr(true)
				config.EvaluatorIds = []string{"Builtin.Helpfulness", "helpfulness", "tone_check"}
				config.Rule.Filters = []*AwsBedrockAgentCoreOnlineEvaluationFilter{
					{Key: "session.channel", Operator: "Equals", StringValue: "web"},
				}
				config.Rule.SessionTimeoutMinutes = 15
				spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an in-bundle evaluator name as an online-config entry", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEvaluator()
				config := minimalOnlineConfig()
				config.EvaluatorIds = []string{"tone_check"}
				spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with gemini and openai model arms", func() {
			ginkgo.It("should accept each vendor arm alone", func() {
				spec := minimalEvaluator()
				harness := minimalHarness()
				harness.Model = &AwsBedrockAgentCoreHarnessModel{
					Gemini: &AwsBedrockAgentCoreHarnessGeminiModel{
						ApiKeyArn: svr("arn:aws:secretsmanager:us-west-2:123456789012:secret:gemini-key"),
						ModelId:   "gemini-2.0-flash",
					},
				}
				spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())

				harness.Model = &AwsBedrockAgentCoreHarnessModel{
					Openai: &AwsBedrockAgentCoreHarnessOpenAiModel{
						ApiKeyArn: svr("arn:aws:secretsmanager:us-west-2:123456789012:secret:openai-key"),
						ModelId:   "gpt-4o",
					},
				}
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an empty bundle", func() {
			spec := &AwsBedrockAgentCoreEvaluationSpec{Region: "us-west-2"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("set at least one of evaluators"))
		})

		ginkgo.It("rejects duplicate evaluator names", func() {
			spec := minimalEvaluator()
			dup := minimalEvaluator().Evaluators[0]
			spec.Evaluators = append(spec.Evaluators, dup)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("evaluators entries must have unique names"))
		})

		ginkgo.It("rejects an evaluator with both scoring arms", func() {
			spec := minimalEvaluator()
			spec.Evaluators[0].LlmAsAJudge = llmJudgeEvaluator().LlmAsAJudge
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of llm_as_a_judge or code_based"))
		})

		ginkgo.It("rejects an evaluator name AWS would refuse", func() {
			spec := minimalEvaluator()
			spec.Evaluators[0].Name = "9starts-with-digit"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid evaluator level", func() {
			spec := minimalEvaluator()
			spec.Evaluators[0].Level = "SPAN"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rating scale with both shapes", func() {
			spec := &AwsBedrockAgentCoreEvaluationSpec{
				Region:     "us-west-2",
				Evaluators: []*AwsBedrockAgentCoreEvaluator{llmJudgeEvaluator()},
			}
			spec.Evaluators[0].LlmAsAJudge.RatingScale.Numerical = []*AwsBedrockAgentCoreNumericalRating{
				{Label: "good", Definition: "good", Value: 1},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of categorical or numerical"))
		})

		ginkgo.It("rejects a judge temperature above 1", func() {
			spec := &AwsBedrockAgentCoreEvaluationSpec{
				Region:     "us-west-2",
				Evaluators: []*AwsBedrockAgentCoreEvaluator{llmJudgeEvaluator()},
			}
			spec.Evaluators[0].LlmAsAJudge.Model.Inference.Temperature = f64(1.5)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a code evaluator timeout above 300", func() {
			spec := minimalEvaluator()
			spec.Evaluators[0].CodeBased.TimeoutSeconds = 301
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a harness with two model vendors", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.Model.Gemini = &AwsBedrockAgentCoreHarnessGeminiModel{
				ApiKeyArn: svr("arn:aws:secretsmanager:us-west-2:123456789012:secret:gemini-key"),
				ModelId:   "gemini-2.0-flash",
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of bedrock, gemini, or openai"))
		})

		ginkgo.It("rejects a harness name above 40 characters", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.Name = "this_harness_name_is_way_too_long_for_aws_limits"
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a hyphenated harness name (CreateHarness names the regex)", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.Name = "support-bench"
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects SESSION-level judge instructions without a placeholder", func() {
			spec := &AwsBedrockAgentCoreEvaluationSpec{
				Region:     "us-west-2",
				Evaluators: []*AwsBedrockAgentCoreEvaluator{llmJudgeEvaluator()},
			}
			spec.Evaluators[0].LlmAsAJudge.Instructions = "Assess whether the agent resolved the user's request."
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must embed at least one placeholder"))
		})

		ginkgo.It("accepts placeholder-less instructions on a non-SESSION judge (contract unverified there)", func() {
			spec := &AwsBedrockAgentCoreEvaluationSpec{
				Region:     "us-west-2",
				Evaluators: []*AwsBedrockAgentCoreEvaluator{llmJudgeEvaluator()},
			}
			spec.Evaluators[0].Level = "TRACE"
			spec.Evaluators[0].LlmAsAJudge.Instructions = "Assess whether the agent resolved the user's request."
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("rejects a tool whose config arm mismatches its type", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.Tools = []*AwsBedrockAgentCoreHarnessTool{
				{
					Name: "bad",
					Type: "remote_mcp",
					AgentcoreGateway: &AwsBedrockAgentCoreHarnessGatewayTool{
						GatewayArn: svr("arn:aws:bedrock-agentcore:us-west-2:123456789012:gateway/gw"),
					},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("config arm matching type"))
		})

		ginkgo.It("rejects allowed_tools naming an unconfigured tool", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.AllowedTools = []string{"ghost-tool"}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("must match the name of a tool configured"))
		})

		ginkgo.It("rejects an outbound auth with two arms", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.Tools = []*AwsBedrockAgentCoreHarnessTool{
				{
					Name: "gw",
					Type: "agentcore_gateway",
					AgentcoreGateway: &AwsBedrockAgentCoreHarnessGatewayTool{
						GatewayArn: svr("arn:aws:bedrock-agentcore:us-west-2:123456789012:gateway/gw"),
						OutboundAuth: &AwsBedrockAgentCoreHarnessGatewayOutboundAuth{
							AwsIam: true,
							NoAuth: true,
						},
					},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of aws_iam, no_auth, or oauth"))
		})

		ginkgo.It("rejects VPC network mode without vpc_config", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.RuntimeEnvironment = &AwsBedrockAgentCoreHarnessRuntimeEnvironment{
				Network: &AwsBedrockAgentCoreHarnessNetwork{Mode: "VPC"},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("vpc_config is required when mode is VPC"))
		})

		ginkgo.It("rejects a filesystem with two source arms", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.RuntimeEnvironment = &AwsBedrockAgentCoreHarnessRuntimeEnvironment{
				Filesystems: []*AwsBedrockAgentCoreHarnessFilesystem{
					{
						MountPath:         "/mnt/x",
						SessionStorage:    true,
						EfsAccessPointArn: svr("arn:aws:elasticfilesystem:us-west-2:123456789012:access-point/fsap-1"),
					},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of efs_access_point_arn"))
		})

		ginkgo.It("rejects truncation config mismatching its strategy", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.Truncation = &AwsBedrockAgentCoreHarnessTruncation{
				Strategy: "none",
				SlidingWindow: &AwsBedrockAgentCoreHarnessSlidingWindow{
					MessagesCount: 10,
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("config matching strategy"))
		})

		ginkgo.It("rejects a custom claim with both match-value shapes", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
				DiscoveryUrl: "https://issuer.example.com/.well-known/openid-configuration",
				CustomClaims: []*AwsBedrockAgentCoreCustomClaim{
					{
						ClaimName:     "org",
						ValueType:     "STRING",
						MatchOperator: "EQUALS",
						MatchValue:    "acme",
						MatchValues:   []string{"acme", "globex"},
					},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of match_value or match_values"))
		})

		ginkgo.It("rejects an unresolvable online-config evaluator entry", func() {
			spec := minimalEvaluator()
			config := minimalOnlineConfig()
			config.EvaluatorIds = []string{"not_defined_here"}
			spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("AWS builtin"))
		})

		ginkgo.It("rejects a sampling percentage below the AWS floor", func() {
			spec := minimalEvaluator()
			config := minimalOnlineConfig()
			config.Rule.SamplingPercentage = 0.001
			spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a filter with two typed values", func() {
			spec := minimalEvaluator()
			config := minimalOnlineConfig()
			config.Rule.Filters = []*AwsBedrockAgentCoreOnlineEvaluationFilter{
				{
					Key:          "session.channel",
					Operator:     "Equals",
					StringValue:  "web",
					BooleanValue: boolPtr(true),
				},
			}
			spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of string_value, boolean_value, or double_value"))
		})

		ginkgo.It("rejects more than five filters", func() {
			spec := minimalEvaluator()
			config := minimalOnlineConfig()
			for i := 0; i < 6; i++ {
				config.Rule.Filters = append(config.Rule.Filters, &AwsBedrockAgentCoreOnlineEvaluationFilter{
					Key:         "k",
					Operator:    "Equals",
					StringValue: "v",
				})
			}
			spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a session timeout above 60 minutes", func() {
			spec := minimalEvaluator()
			config := minimalOnlineConfig()
			config.Rule.SessionTimeoutMinutes = 61
			spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than ten evaluators on one online config", func() {
			spec := minimalEvaluator()
			config := minimalOnlineConfig()
			config.EvaluatorIds = nil
			for i := 0; i < 11; i++ {
				config.EvaluatorIds = append(config.EvaluatorIds, "Builtin.Metric"+string(rune('A'+i)))
			}
			spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate evaluator IDs on one online config", func() {
			spec := minimalEvaluator()
			config := minimalOnlineConfig()
			config.EvaluatorIds = []string{"Builtin.Helpfulness", "Builtin.Helpfulness"}
			spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{config}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate harness names", func() {
			spec := minimalEvaluator()
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{minimalHarness(), minimalHarness()}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("harnesses entries must have unique names"))
		})

		ginkgo.It("rejects duplicate online config names", func() {
			spec := minimalEvaluator()
			spec.OnlineEvaluationConfigs = []*AwsBedrockAgentCoreOnlineEvaluationConfig{
				minimalOnlineConfig(), minimalOnlineConfig(),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("online_evaluation_configs entries must have unique names"))
		})

		ginkgo.It("rejects a private endpoint with both arms", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
				DiscoveryUrl: "https://idp.example.com/.well-known/openid-configuration",
				PrivateEndpoint: &AwsBedrockAgentCorePrivateEndpoint{
					ManagedVpc: &AwsBedrockAgentCoreManagedVpcEndpoint{
						VpcId:     svr("vpc-abc"),
						SubnetIds: []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
					},
					SelfManagedLattice: &AwsBedrockAgentCoreLatticeEndpoint{
						ResourceConfigurationId: "rcfg-abc",
					},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of managed_vpc or self_managed_lattice"))
		})

		ginkgo.It("rejects a private endpoint with no arm", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
				DiscoveryUrl:    "https://idp.example.com/.well-known/openid-configuration",
				PrivateEndpoint: &AwsBedrockAgentCorePrivateEndpoint{},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of managed_vpc or self_managed_lattice"))
		})

		ginkgo.It("rejects an evaluator with no scoring arm", func() {
			spec := minimalEvaluator()
			spec.Evaluators[0].CodeBased = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rating scale with no shape", func() {
			spec := minimalEvaluator()
			judge := llmJudgeEvaluator()
			judge.LlmAsAJudge.RatingScale = &AwsBedrockAgentCoreRatingScale{}
			spec.Evaluators = []*AwsBedrockAgentCoreEvaluator{judge}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a harness model with no vendor arm", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.Model = &AwsBedrockAgentCoreHarnessModel{}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a custom claim with no match-value shape", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
				DiscoveryUrl: "https://idp.example.com/.well-known/openid-configuration",
				CustomClaims: []*AwsBedrockAgentCoreCustomClaim{
					{ClaimName: "org", ValueType: "STRING", MatchOperator: "EQUALS"},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects PUBLIC network mode carrying a vpc_config", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.RuntimeEnvironment = &AwsBedrockAgentCoreHarnessRuntimeEnvironment{
				Network: &AwsBedrockAgentCoreHarnessNetwork{
					Mode: "PUBLIC",
					VpcConfig: &AwsBedrockAgentCoreHarnessVpcConfig{
						Subnets:        []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
						SecurityGroups: []*foreignkeyv1.StringValueOrRef{svr("sg-abc")},
					},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an outbound auth with no arm", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.Tools = []*AwsBedrockAgentCoreHarnessTool{
				{
					Name: "gw",
					Type: "agentcore_gateway",
					AgentcoreGateway: &AwsBedrockAgentCoreHarnessGatewayTool{
						GatewayArn:   svr("arn:aws:bedrock-agentcore:us-west-2:123456789012:gateway/gw-1"),
						OutboundAuth: &AwsBedrockAgentCoreHarnessGatewayOutboundAuth{},
					},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of aws_iam, no_auth, or oauth"))
		})

		ginkgo.It("rejects a routing domain shorter than three characters", func() {
			spec := minimalEvaluator()
			harness := minimalHarness()
			harness.CustomJwtAuthorizer = &AwsBedrockAgentCoreJwtAuthorizer{
				DiscoveryUrl: "https://idp.example.com/.well-known/openid-configuration",
				PrivateEndpoint: &AwsBedrockAgentCorePrivateEndpoint{
					ManagedVpc: &AwsBedrockAgentCoreManagedVpcEndpoint{
						VpcId:         svr("vpc-abc"),
						SubnetIds:     []*foreignkeyv1.StringValueOrRef{svr("subnet-abc")},
						RoutingDomain: "ab",
					},
				},
			}
			spec.Harnesses = []*AwsBedrockAgentCoreHarness{harness}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
