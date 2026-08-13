package awsbedrockguardrailv1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestAwsBedrockGuardrailSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockGuardrailSpec Validation Suite")
}

func boolPtr(b bool) *bool {
	return &b
}

// minimalGuardrail is the smallest valid manifest: region, both blocked
// messagings, and one content filter.
func minimalGuardrail() *AwsBedrockGuardrailSpec {
	return &AwsBedrockGuardrailSpec{
		Region:                  "us-west-2",
		BlockedInputMessaging:   "Sorry, I can't help with that request.",
		BlockedOutputsMessaging: "Sorry, I can't provide that response.",
		ContentPolicy: &AwsBedrockGuardrailContentPolicy{
			Filters: []*AwsBedrockGuardrailContentFilter{
				{Type: "HATE", InputStrength: "HIGH", OutputStrength: "HIGH"},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsBedrockGuardrailSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalGuardrail())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with every policy family configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalGuardrail()
				spec.Description = "full-surface guardrail"
				spec.ContentPolicy = &AwsBedrockGuardrailContentPolicy{
					Tier: "STANDARD",
					Filters: []*AwsBedrockGuardrailContentFilter{
						{
							Type:             "SEXUAL",
							InputStrength:    "HIGH",
							OutputStrength:   "MEDIUM",
							InputAction:      "BLOCK",
							OutputAction:     "NONE",
							InputEnabled:     boolPtr(true),
							OutputEnabled:    boolPtr(false),
							InputModalities:  []string{"TEXT", "IMAGE"},
							OutputModalities: []string{"TEXT"},
						},
						{Type: "PROMPT_ATTACK", InputStrength: "HIGH", OutputStrength: "NONE"},
					},
				}
				spec.TopicPolicy = &AwsBedrockGuardrailTopicPolicy{
					Tier: "CLASSIC",
					Topics: []*AwsBedrockGuardrailTopic{
						{
							Name:       "investment-advice",
							Definition: "Providing investment advice or recommending specific financial products.",
							Examples:   []string{"Should I buy this stock?"},
						},
					},
				}
				spec.WordPolicy = &AwsBedrockGuardrailWordPolicy{
					ProfanityFilter: &AwsBedrockGuardrailManagedWordList{
						InputAction:   "BLOCK",
						OutputAction:  "NONE",
						InputEnabled:  boolPtr(true),
						OutputEnabled: boolPtr(true),
					},
					CustomWords: []*AwsBedrockGuardrailCustomWord{
						{Text: "codename-atlas", InputAction: "BLOCK"},
					},
				}
				spec.SensitiveInformationPolicy = &AwsBedrockGuardrailSensitiveInformationPolicy{
					PiiEntities: []*AwsBedrockGuardrailPiiEntity{
						{Type: "EMAIL", Action: "ANONYMIZE"},
						{Type: "US_SOCIAL_SECURITY_NUMBER", Action: "BLOCK", InputAction: "BLOCK", OutputAction: "ANONYMIZE"},
					},
					Regexes: []*AwsBedrockGuardrailRegex{
						{Name: "employee-id", Pattern: "EMP-[0-9]{6}", Action: "ANONYMIZE", Description: "internal employee identifiers"},
					},
				}
				spec.ContextualGroundingPolicy = &AwsBedrockGuardrailContextualGroundingPolicy{
					Filters: []*AwsBedrockGuardrailContextualGroundingFilter{
						{Type: "GROUNDING", Threshold: 0.75},
						{Type: "RELEVANCE", Threshold: 0.5},
					},
				}
				spec.CrossRegionProfileArn = "arn:aws:bedrock:us-west-2:123456789012:guardrail-profile/us.guardrail.v1:0"
				spec.KmsKeyArn = nil
				spec.Versions = []*AwsBedrockGuardrailVersion{
					{Name: "prod", Description: "initial production pin"},
					{Name: "canary", KeepOnDelete: true},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with grounding thresholds at the domain edges", func() {
			ginkgo.It("should accept 0 and 1", func() {
				spec := minimalGuardrail()
				spec.ContextualGroundingPolicy = &AwsBedrockGuardrailContextualGroundingPolicy{
					Filters: []*AwsBedrockGuardrailContextualGroundingFilter{
						{Type: "GROUNDING", Threshold: 0},
						{Type: "RELEVANCE", Threshold: 1},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Required fields and ranges
	// -----------------------------------------------------------------
	ginkgo.Describe("Required fields and ranges", func() {

		ginkgo.It("should reject a missing region", func() {
			spec := minimalGuardrail()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject missing blocked_input_messaging", func() {
			spec := minimalGuardrail()
			spec.BlockedInputMessaging = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject blocked_outputs_messaging over 500 characters", func() {
			spec := minimalGuardrail()
			spec.BlockedOutputsMessaging = strings.Repeat("x", 501)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a description over 200 characters", func() {
			spec := minimalGuardrail()
			spec.Description = strings.Repeat("d", 201)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Content policy
	// -----------------------------------------------------------------
	ginkgo.Describe("Content policy", func() {

		ginkgo.It("should reject an empty filters list", func() {
			spec := minimalGuardrail()
			spec.ContentPolicy = &AwsBedrockGuardrailContentPolicy{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an unknown filter type", func() {
			spec := minimalGuardrail()
			spec.ContentPolicy.Filters[0].Type = "SPAM"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an unknown strength", func() {
			spec := minimalGuardrail()
			spec.ContentPolicy.Filters[0].InputStrength = "MAXIMUM"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a missing output strength", func() {
			spec := minimalGuardrail()
			spec.ContentPolicy.Filters[0].OutputStrength = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a duplicate filter type", func() {
			spec := minimalGuardrail()
			spec.ContentPolicy.Filters = append(spec.ContentPolicy.Filters,
				&AwsBedrockGuardrailContentFilter{Type: "HATE", InputStrength: "LOW", OutputStrength: "LOW"})
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("content filter type"))
		})

		ginkgo.It("should reject an unknown tier", func() {
			spec := minimalGuardrail()
			spec.ContentPolicy.Tier = "PREMIUM"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an unknown modality", func() {
			spec := minimalGuardrail()
			spec.ContentPolicy.Filters[0].InputModalities = []string{"AUDIO"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an unknown action", func() {
			spec := minimalGuardrail()
			spec.ContentPolicy.Filters[0].InputAction = "REDACT"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Topic policy
	// -----------------------------------------------------------------
	ginkgo.Describe("Topic policy", func() {

		ginkgo.It("should reject an empty topics list", func() {
			spec := minimalGuardrail()
			spec.TopicPolicy = &AwsBedrockGuardrailTopicPolicy{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a topic name with illegal characters", func() {
			spec := minimalGuardrail()
			spec.TopicPolicy = &AwsBedrockGuardrailTopicPolicy{
				Topics: []*AwsBedrockGuardrailTopic{
					{Name: "bad/name", Definition: "d"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a definition over 1000 characters", func() {
			spec := minimalGuardrail()
			spec.TopicPolicy = &AwsBedrockGuardrailTopicPolicy{
				Topics: []*AwsBedrockGuardrailTopic{
					{Name: "t", Definition: strings.Repeat("d", 1001)},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject duplicate topic names", func() {
			spec := minimalGuardrail()
			spec.TopicPolicy = &AwsBedrockGuardrailTopicPolicy{
				Topics: []*AwsBedrockGuardrailTopic{
					{Name: "t", Definition: "a"},
					{Name: "t", Definition: "b"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("topic names"))
		})

		ginkgo.It("should reject an example over 100 characters", func() {
			spec := minimalGuardrail()
			spec.TopicPolicy = &AwsBedrockGuardrailTopicPolicy{
				Topics: []*AwsBedrockGuardrailTopic{
					{Name: "t", Definition: "d", Examples: []string{strings.Repeat("e", 101)}},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Word policy
	// -----------------------------------------------------------------
	ginkgo.Describe("Word policy", func() {

		ginkgo.It("should reject an empty word policy", func() {
			spec := minimalGuardrail()
			spec.WordPolicy = &AwsBedrockGuardrailWordPolicy{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("word_policy"))
		})

		ginkgo.It("should accept a profanity-filter-only word policy", func() {
			spec := minimalGuardrail()
			spec.WordPolicy = &AwsBedrockGuardrailWordPolicy{
				ProfanityFilter: &AwsBedrockGuardrailManagedWordList{},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should reject a custom word with empty text", func() {
			spec := minimalGuardrail()
			spec.WordPolicy = &AwsBedrockGuardrailWordPolicy{
				CustomWords: []*AwsBedrockGuardrailCustomWord{{Text: ""}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject duplicate custom word texts", func() {
			spec := minimalGuardrail()
			spec.WordPolicy = &AwsBedrockGuardrailWordPolicy{
				CustomWords: []*AwsBedrockGuardrailCustomWord{
					{Text: "atlas"},
					{Text: "atlas"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("unique text"))
		})
	})

	// -----------------------------------------------------------------
	// Sensitive information policy
	// -----------------------------------------------------------------
	ginkgo.Describe("Sensitive information policy", func() {

		ginkgo.It("should reject an empty policy", func() {
			spec := minimalGuardrail()
			spec.SensitiveInformationPolicy = &AwsBedrockGuardrailSensitiveInformationPolicy{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject an unknown PII entity type", func() {
			spec := minimalGuardrail()
			spec.SensitiveInformationPolicy = &AwsBedrockGuardrailSensitiveInformationPolicy{
				PiiEntities: []*AwsBedrockGuardrailPiiEntity{{Type: "FAVORITE_COLOR", Action: "BLOCK"}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a duplicate PII entity type", func() {
			spec := minimalGuardrail()
			spec.SensitiveInformationPolicy = &AwsBedrockGuardrailSensitiveInformationPolicy{
				PiiEntities: []*AwsBedrockGuardrailPiiEntity{
					{Type: "EMAIL", Action: "BLOCK"},
					{Type: "EMAIL", Action: "ANONYMIZE"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("PII entity type"))
		})

		ginkgo.It("should reject a missing PII action", func() {
			spec := minimalGuardrail()
			spec.SensitiveInformationPolicy = &AwsBedrockGuardrailSensitiveInformationPolicy{
				PiiEntities: []*AwsBedrockGuardrailPiiEntity{{Type: "EMAIL"}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a regex without a pattern", func() {
			spec := minimalGuardrail()
			spec.SensitiveInformationPolicy = &AwsBedrockGuardrailSensitiveInformationPolicy{
				Regexes: []*AwsBedrockGuardrailRegex{{Name: "r", Action: "BLOCK"}},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject duplicate regex names", func() {
			spec := minimalGuardrail()
			spec.SensitiveInformationPolicy = &AwsBedrockGuardrailSensitiveInformationPolicy{
				Regexes: []*AwsBedrockGuardrailRegex{
					{Name: "r", Pattern: "a+", Action: "BLOCK"},
					{Name: "r", Pattern: "b+", Action: "NONE"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("regex names"))
		})
	})

	// -----------------------------------------------------------------
	// Contextual grounding policy
	// -----------------------------------------------------------------
	ginkgo.Describe("Contextual grounding policy", func() {

		ginkgo.It("should reject an empty filters list", func() {
			spec := minimalGuardrail()
			spec.ContextualGroundingPolicy = &AwsBedrockGuardrailContextualGroundingPolicy{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a threshold above 1", func() {
			spec := minimalGuardrail()
			spec.ContextualGroundingPolicy = &AwsBedrockGuardrailContextualGroundingPolicy{
				Filters: []*AwsBedrockGuardrailContextualGroundingFilter{
					{Type: "GROUNDING", Threshold: 1.5},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a duplicate filter type", func() {
			spec := minimalGuardrail()
			spec.ContextualGroundingPolicy = &AwsBedrockGuardrailContextualGroundingPolicy{
				Filters: []*AwsBedrockGuardrailContextualGroundingFilter{
					{Type: "GROUNDING", Threshold: 0.5},
					{Type: "GROUNDING", Threshold: 0.7},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("grounding filter type"))
		})
	})

	// -----------------------------------------------------------------
	// Cross-region profile
	// -----------------------------------------------------------------
	ginkgo.Describe("Cross-region profile", func() {

		ginkgo.It("should reject a non-guardrail-profile ARN", func() {
			spec := minimalGuardrail()
			spec.CrossRegionProfileArn = "arn:aws:bedrock:us-west-2:123456789012:guardrail/gr-123"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -----------------------------------------------------------------
	// Versions
	// -----------------------------------------------------------------
	ginkgo.Describe("Versions", func() {

		ginkgo.It("should reject duplicate version entry names", func() {
			spec := minimalGuardrail()
			spec.Versions = []*AwsBedrockGuardrailVersion{
				{Name: "prod"},
				{Name: "prod"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("unique names"))
		})

		ginkgo.It("should reject a version name with illegal characters", func() {
			spec := minimalGuardrail()
			spec.Versions = []*AwsBedrockGuardrailVersion{{Name: "prod v1"}}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should reject a version description over 200 characters", func() {
			spec := minimalGuardrail()
			spec.Versions = []*AwsBedrockGuardrailVersion{
				{Name: "prod", Description: strings.Repeat("d", 201)},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
