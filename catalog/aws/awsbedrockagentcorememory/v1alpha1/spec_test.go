package awsbedrockagentcorememoryv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockAgentCoreMemorySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockAgentCoreMemorySpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalMemory is the smallest valid manifest: region, name, and the
// short-term retention window.
func minimalMemory() *AwsBedrockAgentCoreMemorySpec {
	return &AwsBedrockAgentCoreMemorySpec{
		Region:          "us-west-2",
		MemoryName:      "support_memory",
		EventExpiryDays: 30,
	}
}

var _ = ginkgo.Describe("AwsBedrockAgentCoreMemorySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalMemory())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalMemory()
				spec.Description = "long-term memory for the support agent"
				spec.EncryptionKeyArn = svr("arn:aws:kms:us-west-2:123456789012:key/abc")
				spec.ExecutionRoleArn = svr("arn:aws:iam::123456789012:role/agentcore-memory")
				spec.IndexedKeys = []*AwsBedrockAgentCoreMemoryIndexedKey{
					{Key: "customer_id", Type: "STRING"},
					{Key: "priority", Type: "NUMBER"},
				}
				spec.KinesisDelivery = &AwsBedrockAgentCoreMemoryKinesisDelivery{
					DataStreamArn: svr("arn:aws:kinesis:us-west-2:123456789012:stream/memory-records"),
					ContentLevel:  "FULL_CONTENT",
				}
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{Name: "facts", Type: "SEMANTIC", NamespaceTemplates: []string{"/facts/{actorId}"}},
					{Name: "summaries", Type: "SUMMARIZATION", NamespaceTemplates: []string{"/summaries/{actorId}/{sessionId}"}},
					{
						Name:               "episodes",
						Type:               "EPISODIC",
						NamespaceTemplates: []string{"/episodes/{actorId}"},
						// A strict prefix of the episodic namespace -- the
						// server-side shape AWS accepts.
						ReflectionNamespaceTemplates: []string{
							"/episodes",
						},
					},
					{
						Name:               "tuned_prefs",
						Type:               "CUSTOM",
						NamespaceTemplates: []string{"/preferences/{actorId}"},
						Custom: &AwsBedrockAgentCoreMemoryCustomStrategy{
							Type: "USER_PREFERENCE_OVERRIDE",
							Extraction: &AwsBedrockAgentCoreMemoryPromptOverride{
								AppendToPrompt: "Focus on stated product preferences.",
								ModelId:        "anthropic.claude-3-5-sonnet-20241022-v2:0",
							},
							Consolidation: &AwsBedrockAgentCoreMemoryPromptOverride{
								AppendToPrompt: "Merge overlapping preferences.",
								ModelId:        "anthropic.claude-3-5-sonnet-20241022-v2:0",
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an EPISODIC_OVERRIDE custom strategy carrying reflection", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{
						Name:               "episodes",
						Type:               "CUSTOM",
						NamespaceTemplates: []string{"/episodes/{actorId}"},
						Custom: &AwsBedrockAgentCoreMemoryCustomStrategy{
							Type: "EPISODIC_OVERRIDE",
							Reflection: &AwsBedrockAgentCoreMemoryReflectionOverride{
								AppendToPrompt:     "Reflect on what worked.",
								ModelId:            "anthropic.claude-3-5-sonnet-20241022-v2:0",
								NamespaceTemplates: []string{"/reflections/{actorId}"},
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with a hyphenated memory name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.MemoryName = "support-memory"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a retention window below AWS's floor", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.EventExpiryDays = 3
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with eleven indexed keys", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				for i := 0; i < 11; i++ {
					spec.IndexedKeys = append(spec.IndexedKeys, &AwsBedrockAgentCoreMemoryIndexedKey{
						Key:  string(rune('a' + i)),
						Type: "STRING",
					})
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate indexed keys", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.IndexedKeys = []*AwsBedrockAgentCoreMemoryIndexedKey{
					{Key: "customer_id", Type: "STRING"},
					{Key: "customer_id", Type: "NUMBER"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an unknown indexed-key type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.IndexedKeys = []*AwsBedrockAgentCoreMemoryIndexedKey{
					{Key: "customer_id", Type: "BOOLEAN"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate strategy names", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{Name: "facts", Type: "SEMANTIC", NamespaceTemplates: []string{"/facts/{actorId}"}},
					{Name: "facts", Type: "SUMMARIZATION", NamespaceTemplates: []string{"/facts/{actorId}"}},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a strategy missing namespace templates", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{Name: "facts", Type: "SEMANTIC"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a CUSTOM strategy missing its overrides", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{Name: "tuned", Type: "CUSTOM", NamespaceTemplates: []string{"/preferences/{actorId}"}},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a built-in strategy carrying overrides", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{
						Name:               "facts",
						Type:               "SEMANTIC",
						NamespaceTemplates: []string{"/facts/{actorId}"},
						Custom: &AwsBedrockAgentCoreMemoryCustomStrategy{
							Type: "SEMANTIC_OVERRIDE",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with reflection namespaces on a non-EPISODIC strategy", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{
						Name:                         "facts",
						Type:                         "SEMANTIC",
						NamespaceTemplates:           []string{"/facts/{actorId}"},
						ReflectionNamespaceTemplates: []string{"/facts/{actorId}"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a reflection namespace unrelated to every episodic namespace", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{
						Name:               "episodes",
						Type:               "EPISODIC",
						NamespaceTemplates: []string{"/episodes/{actorId}"},
						// AWS requires equal-or-hierarchical-prefix; an
						// unrelated root is rejected server-side, so the
						// spec front-loads it.
						ReflectionNamespaceTemplates: []string{"/reflections/{actorId}"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a reflection namespace that is a string prefix but not a hierarchical prefix", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{
						Name:               "episodes",
						Type:               "EPISODIC",
						NamespaceTemplates: []string{"/episodes/{actorId}"},
						// "/epi" starts the episodic namespace string but
						// breaks mid-segment -- AWS's rule is hierarchical
						// (whole path segments), so this is invalid.
						ReflectionNamespaceTemplates: []string{"/epi"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an EPISODIC_OVERRIDE missing reflection", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{
						Name:               "episodes",
						Type:               "CUSTOM",
						NamespaceTemplates: []string{"/episodes/{actorId}"},
						Custom: &AwsBedrockAgentCoreMemoryCustomStrategy{
							Type: "EPISODIC_OVERRIDE",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a SUMMARY_OVERRIDE carrying extraction", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.Strategies = []*AwsBedrockAgentCoreMemoryStrategy{
					{
						Name:               "summaries",
						Type:               "CUSTOM",
						NamespaceTemplates: []string{"/summaries/{actorId}/{sessionId}"},
						Custom: &AwsBedrockAgentCoreMemoryCustomStrategy{
							Type: "SUMMARY_OVERRIDE",
							Extraction: &AwsBedrockAgentCoreMemoryPromptOverride{
								AppendToPrompt: "x",
								ModelId:        "m",
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an unknown kinesis content level", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalMemory()
				spec.KinesisDelivery = &AwsBedrockAgentCoreMemoryKinesisDelivery{
					DataStreamArn: svr("arn:aws:kinesis:us-west-2:123456789012:stream/s"),
					ContentLevel:  "EVERYTHING",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
