package awsbedrockflowv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockFlowSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockFlowSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalFlow is the smallest valid manifest: region and the execution
// role (an empty shell without a definition).
func minimalFlow() *AwsBedrockFlowSpec {
	return &AwsBedrockFlowSpec{
		Region:           "us-west-2",
		ExecutionRoleArn: svr("arn:aws:iam::123456789012:role/bedrock-flow"),
	}
}

// promptPipelineFlow is an Input -> Prompt(inline) -> Output graph.
func promptPipelineFlow() *AwsBedrockFlowSpec {
	spec := minimalFlow()
	spec.Definition = &AwsBedrockFlowDefinition{
		Nodes: []*AwsBedrockFlowNode{
			{
				Name:    "FlowInput",
				Type:    "Input",
				Outputs: []*AwsBedrockFlowNodeOutput{{Name: "document", Type: "String"}},
			},
			{
				Name: "Summarize",
				Type: "Prompt",
				Inputs: []*AwsBedrockFlowNodeInput{
					{Name: "input", Expression: "$.data", Type: "String"},
				},
				Outputs: []*AwsBedrockFlowNodeOutput{{Name: "modelCompletion", Type: "String"}},
				Prompt: &AwsBedrockFlowPromptNode{
					Inline: &AwsBedrockFlowInlinePrompt{
						ModelId: "amazon.nova-micro-v1:0",
						Text: &AwsBedrockFlowPromptTextTemplate{
							Text:           "Summarize: {{input}}",
							InputVariables: []string{"input"},
						},
					},
				},
			},
			{
				Name: "FlowOutput",
				Type: "Output",
				Inputs: []*AwsBedrockFlowNodeInput{
					{Name: "document", Expression: "$.data", Type: "String"},
				},
			},
		},
		Connections: []*AwsBedrockFlowConnection{
			{
				Name:   "InToPrompt",
				Source: "FlowInput",
				Target: "Summarize",
				Data:   &AwsBedrockFlowDataConnection{SourceOutput: "document", TargetInput: "input"},
			},
			{
				Name:   "PromptToOut",
				Source: "Summarize",
				Target: "FlowOutput",
				Data:   &AwsBedrockFlowDataConnection{SourceOutput: "modelCompletion", TargetInput: "document"},
			},
		},
	}
	return spec
}

var _ = ginkgo.Describe("AwsBedrockFlowSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields (no definition)", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalFlow())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an Input -> Prompt -> Output pipeline", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(promptPipelineFlow())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with condition, knowledge base, and agent nodes", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := promptPipelineFlow()
				spec.Description = "routing pipeline"
				spec.CustomerEncryptionKeyArn = svr("arn:aws:kms:us-west-2:123456789012:key/abc")
				spec.Definition.Nodes = append(spec.Definition.Nodes,
					&AwsBedrockFlowNode{
						Name: "Route",
						Type: "Condition",
						Inputs: []*AwsBedrockFlowNodeInput{
							{Name: "category", Expression: "$.data", Type: "String"},
						},
						Condition: &AwsBedrockFlowConditionNode{
							Conditions: []*AwsBedrockFlowCondition{
								{Name: "billing", Expression: "category == \"billing\""},
								{Name: "default"},
							},
						},
					},
					&AwsBedrockFlowNode{
						Name: "AskKb",
						Type: "KnowledgeBase",
						Inputs: []*AwsBedrockFlowNodeInput{
							{Name: "retrievalQuery", Expression: "$.data", Type: "String"},
						},
						Outputs: []*AwsBedrockFlowNodeOutput{{Name: "outputText", Type: "String"}},
						KnowledgeBase: &AwsBedrockFlowKnowledgeBaseNode{
							KnowledgeBaseId: svr("EMDPPAYPZI"),
							ModelId:         "amazon.nova-lite-v1:0",
							NumberOfResults: 5,
							Guardrail: &AwsBedrockFlowGuardrail{
								GuardrailId: svr("gr-abc"),
								Version:     "1",
							},
							InferenceConfiguration: &AwsBedrockFlowInferenceConfiguration{
								Temperature: float64Ptr(0.2),
							},
						},
					},
					&AwsBedrockFlowNode{
						Name: "Delegate",
						Type: "Agent",
						Inputs: []*AwsBedrockFlowNodeInput{
							{Name: "agentInputText", Expression: "$.data", Type: "String"},
						},
						Outputs: []*AwsBedrockFlowNodeOutput{{Name: "agentResponse", Type: "String"}},
						Agent: &AwsBedrockFlowAgentNode{
							AgentAliasArn: svr("arn:aws:bedrock:us-west-2:123456789012:agent-alias/A/B"),
						},
					},
				)
				spec.Definition.Connections = append(spec.Definition.Connections,
					&AwsBedrockFlowConnection{
						Name:        "RouteBilling",
						Source:      "Route",
						Target:      "Delegate",
						Conditional: &AwsBedrockFlowConditionalConnection{Condition: "billing"},
					},
				)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Node/connection identity
	// -----------------------------------------------------------------
	ginkgo.Describe("Graph identity", func() {

		ginkgo.It("should reject duplicate node names", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes = append(spec.Definition.Nodes, spec.Definition.Nodes[0])
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject duplicate connection names", func() {
			spec := promptPipelineFlow()
			spec.Definition.Connections = append(spec.Definition.Connections, spec.Definition.Connections[0])
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an empty definition", func() {
			spec := minimalFlow()
			spec.Definition = &AwsBedrockFlowDefinition{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Node configuration matching
	// -----------------------------------------------------------------
	ginkgo.Describe("Node configuration matching", func() {

		ginkgo.It("should reject an Agent node without its arm", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes = append(spec.Definition.Nodes, &AwsBedrockFlowNode{
				Name: "Delegate",
				Type: "Agent",
			})
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an Input node carrying a prompt arm", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes[0].Prompt = &AwsBedrockFlowPromptNode{
				PromptArn: svr("arn:aws:bedrock:us-west-2:123456789012:prompt/P"),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should accept a structural Iterator node with no arm", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes = append(spec.Definition.Nodes, &AwsBedrockFlowNode{
				Name: "EachItem",
				Type: "Iterator",
				Inputs: []*AwsBedrockFlowNodeInput{
					{Name: "array", Expression: "$.data", Type: "Array"},
				},
				Outputs: []*AwsBedrockFlowNodeOutput{
					{Name: "arrayItem", Type: "Object"},
				},
			})
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Connections
	// -----------------------------------------------------------------
	ginkgo.Describe("Connections", func() {

		ginkgo.It("should reject a connection with both arms", func() {
			spec := promptPipelineFlow()
			spec.Definition.Connections[0].Conditional = &AwsBedrockFlowConditionalConnection{Condition: "x"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a connection with no arm", func() {
			spec := promptPipelineFlow()
			spec.Definition.Connections[0].Data = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Prompt node shapes
	// -----------------------------------------------------------------
	ginkgo.Describe("Prompt node shapes", func() {

		ginkgo.It("should reject a prompt node with both sources", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes[1].Prompt.PromptArn = svr("arn:aws:bedrock:us-west-2:123456789012:prompt/P")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an inline prompt with both templates", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes[1].Prompt.Inline.Chat = &AwsBedrockFlowPromptChatTemplate{
				Messages: []*AwsBedrockFlowPromptMessage{{Role: "user", Text: "hi"}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a chat message with both text and cache_point", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes[1].Prompt.Inline.Text = nil
			spec.Definition.Nodes[1].Prompt.Inline.Chat = &AwsBedrockFlowPromptChatTemplate{
				Messages: []*AwsBedrockFlowPromptMessage{{Role: "user", Text: "hi", CachePoint: true}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Node fields
	// -----------------------------------------------------------------
	ginkgo.Describe("Node fields", func() {

		ginkgo.It("should reject a node name that starts with a digit", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes[0].Name = "1Input"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject more than five outputs", func() {
			spec := promptPipelineFlow()
			outs := make([]*AwsBedrockFlowNodeOutput, 6)
			for i := range outs {
				outs[i] = &AwsBedrockFlowNodeOutput{Name: "o" + string(rune('a'+i)), Type: "String"}
			}
			spec.Definition.Nodes[0].Outputs = outs
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an unknown IO data type", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes[0].Outputs[0].Type = "Text"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject condition counts above five", func() {
			spec := promptPipelineFlow()
			conds := make([]*AwsBedrockFlowCondition, 6)
			for i := range conds {
				conds[i] = &AwsBedrockFlowCondition{Name: "c" + string(rune('a'+i))}
			}
			spec.Definition.Nodes = append(spec.Definition.Nodes, &AwsBedrockFlowNode{
				Name:      "Route",
				Type:      "Condition",
				Condition: &AwsBedrockFlowConditionNode{Conditions: conds},
			})
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an out-of-range number_of_results", func() {
			spec := promptPipelineFlow()
			spec.Definition.Nodes = append(spec.Definition.Nodes, &AwsBedrockFlowNode{
				Name: "AskKb",
				Type: "KnowledgeBase",
				KnowledgeBase: &AwsBedrockFlowKnowledgeBaseNode{
					KnowledgeBaseId: svr("EMDPPAYPZI"),
					NumberOfResults: 500,
				},
			})
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})

func float64Ptr(v float64) *float64 {
	return &v
}
