package awsbedrockpromptv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockPromptSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockPromptSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalPrompt is the smallest valid manifest: region and one text
// variant targeting a model.
func minimalPrompt() *AwsBedrockPromptSpec {
	return &AwsBedrockPromptSpec{
		Region: "us-west-2",
		Variants: []*AwsBedrockPromptVariant{{
			Name:    "main",
			ModelId: "amazon.nova-micro-v1:0",
			Text: &AwsBedrockPromptTextTemplate{
				Text:           "Summarize the following text: {{input}}",
				InputVariables: []string{"input"},
			},
		}},
	}
}

var _ = ginkgo.Describe("AwsBedrockPromptSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalPrompt())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a chat variant, tools, and a default variant", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalPrompt()
				spec.Description = "support answer prompt"
				spec.DefaultVariant = "chat"
				spec.CustomerEncryptionKeyArn = svr("arn:aws:kms:us-west-2:123456789012:key/abc")
				spec.Variants = append(spec.Variants, &AwsBedrockPromptVariant{
					Name:    "chat",
					ModelId: "amazon.nova-lite-v1:0",
					Chat: &AwsBedrockPromptChatTemplate{
						Messages: []*AwsBedrockPromptMessage{
							{Role: "user", Text: "Answer the question: {{question}}"},
							{Role: "assistant", Text: "Certainly:"},
							{Role: "user", CachePoint: true},
						},
						System: []*AwsBedrockPromptSystemBlock{
							{Text: "You are a concise support assistant."},
							{CachePoint: true},
						},
						InputVariables: []string{"question"},
						ToolConfiguration: &AwsBedrockPromptToolConfiguration{
							Tools: []*AwsBedrockPromptTool{
								{Spec: &AwsBedrockPromptToolSpec{Name: "lookup_order", Description: "Look up an order."}},
								{CachePoint: true},
							},
							ToolChoice: &AwsBedrockPromptToolChoice{Auto: true},
						},
					},
					InferenceConfiguration: &AwsBedrockPromptInferenceConfiguration{
						MaxTokens:   int32Ptr(512),
						Temperature: float64Ptr(0),
						TopP:        float64Ptr(0.9),
					},
					Metadata: map[string]string{"team": "support"},
				})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a variant executed through an agent alias", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalPrompt()
				spec.Variants[0].ModelId = ""
				spec.Variants[0].AgentAliasArn = svr("arn:aws:bedrock:us-west-2:123456789012:agent-alias/A/B")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Spec-level rules
	// -----------------------------------------------------------------
	ginkgo.Describe("Spec-level rules", func() {

		ginkgo.It("should reject a prompt without variants", func() {
			spec := minimalPrompt()
			spec.Variants = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject duplicate variant names", func() {
			spec := minimalPrompt()
			spec.Variants = append(spec.Variants, spec.Variants[0])
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a default_variant that names no variant", func() {
			spec := minimalPrompt()
			spec.DefaultVariant = "missing"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("default_variant"))
		})
	})

	// -----------------------------------------------------------------
	// Variant shape
	// -----------------------------------------------------------------
	ginkgo.Describe("Variant shape", func() {

		ginkgo.It("should reject a variant with both model_id and agent_alias_arn", func() {
			spec := minimalPrompt()
			spec.Variants[0].AgentAliasArn = svr("arn:aws:bedrock:us-west-2:123456789012:agent-alias/A/B")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a variant with neither execution target", func() {
			spec := minimalPrompt()
			spec.Variants[0].ModelId = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a variant with both templates", func() {
			spec := minimalPrompt()
			spec.Variants[0].Chat = &AwsBedrockPromptChatTemplate{
				Messages: []*AwsBedrockPromptMessage{{Role: "user", Text: "hi"}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a variant with no template", func() {
			spec := minimalPrompt()
			spec.Variants[0].Text = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a chat template without messages", func() {
			spec := minimalPrompt()
			spec.Variants[0].Text = nil
			spec.Variants[0].Chat = &AwsBedrockPromptChatTemplate{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a message with both text and cache_point", func() {
			spec := minimalPrompt()
			spec.Variants[0].Text = nil
			spec.Variants[0].Chat = &AwsBedrockPromptChatTemplate{
				Messages: []*AwsBedrockPromptMessage{{Role: "user", Text: "hi", CachePoint: true}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a tool entry with both spec and cache_point", func() {
			spec := minimalPrompt()
			spec.Variants[0].Text = nil
			spec.Variants[0].Chat = &AwsBedrockPromptChatTemplate{
				Messages: []*AwsBedrockPromptMessage{{Role: "user", Text: "hi"}},
				ToolConfiguration: &AwsBedrockPromptToolConfiguration{
					Tools: []*AwsBedrockPromptTool{{
						Spec:       &AwsBedrockPromptToolSpec{Name: "t"},
						CachePoint: true,
					}},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a tool choice with two modes", func() {
			spec := minimalPrompt()
			spec.Variants[0].Text = nil
			spec.Variants[0].Chat = &AwsBedrockPromptChatTemplate{
				Messages: []*AwsBedrockPromptMessage{{Role: "user", Text: "hi"}},
				ToolConfiguration: &AwsBedrockPromptToolConfiguration{
					Tools:      []*AwsBedrockPromptTool{{Spec: &AwsBedrockPromptToolSpec{Name: "t"}}},
					ToolChoice: &AwsBedrockPromptToolChoice{Auto: true, ToolName: "t"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an out-of-range temperature", func() {
			spec := minimalPrompt()
			spec.Variants[0].InferenceConfiguration = &AwsBedrockPromptInferenceConfiguration{
				Temperature: float64Ptr(1.5),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})

func int32Ptr(v int32) *int32 {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}
