package awsbedrockagentcoreidentityv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBedrockAgentCoreIdentitySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockAgentCoreIdentitySpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalIdentity is the smallest valid manifest: region plus one arm.
func minimalIdentity() *AwsBedrockAgentCoreIdentitySpec {
	return &AwsBedrockAgentCoreIdentitySpec{
		Region: "us-west-2",
		WorkloadIdentities: []*AwsBedrockAgentCoreWorkloadIdentity{
			{Name: "support-agent"},
		},
	}
}

var _ = ginkgo.Describe("AwsBedrockAgentCoreIdentitySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalIdentity())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalIdentity()
				spec.WorkloadIdentities[0].AllowedResourceOauth2ReturnUrls = []string{"https://app.example.com/callback"}
				spec.ApiKeyCredentialProviders = []*AwsBedrockAgentCoreApiKeyProvider{
					{Name: "docs-api", ApiKey: "sk-test-abc123"},
				}
				spec.Oauth2CredentialProviders = []*AwsBedrockAgentCoreOauth2Provider{
					{
						Name:         "github",
						Vendor:       "GITHUB",
						ClientId:     "Iv1.abc",
						ClientSecret: "ghp_secret",
					},
					{
						Name:         "internal-idp",
						Vendor:       "CUSTOM",
						ClientId:     "planton-agents",
						ClientSecret: "s3cret",
						OauthDiscovery: &AwsBedrockAgentCoreOauth2Discovery{
							DiscoveryUrl: "https://idp.example.com/.well-known/openid-configuration",
						},
					},
					{
						Name:         "legacy-idp",
						Vendor:       "CUSTOM",
						ClientId:     "legacy",
						ClientSecret: "s3cret",
						OauthDiscovery: &AwsBedrockAgentCoreOauth2Discovery{
							AuthorizationServerMetadata: &AwsBedrockAgentCoreOauth2ServerMetadata{
								Issuer:                "https://legacy.example.com",
								AuthorizationEndpoint: "https://legacy.example.com/authorize",
								TokenEndpoint:         "https://legacy.example.com/token",
								ResponseTypes:         []string{"code"},
							},
						},
					},
				}
				spec.PolicyEngine = &AwsBedrockAgentCorePolicyEngine{
					EngineName:       "agent_authz",
					Description:      "tool-call authorization",
					EncryptionKeyArn: svr("arn:aws:kms:us-west-2:123456789012:key/abc"),
					Policies: []*AwsBedrockAgentCoreCedarPolicy{
						{
							Name:           "allow_reads",
							CedarStatement: `permit(principal, action == Action::"get_order", resource);`,
							ValidationMode: "FAIL_ON_ANY_FINDINGS",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with only a policy engine", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsBedrockAgentCoreIdentitySpec{
					Region: "us-west-2",
					PolicyEngine: &AwsBedrockAgentCorePolicyEngine{
						EngineName: "authz",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with an empty bundle", func() {
			ginkgo.It("should return a validation error", func() {
				spec := &AwsBedrockAgentCoreIdentitySpec{Region: "us-west-2"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a two-character workload identity name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.WorkloadIdentities[0].Name = "ab"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate workload identity names", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.WorkloadIdentities = append(spec.WorkloadIdentities, &AwsBedrockAgentCoreWorkloadIdentity{Name: "support-agent"})
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an api-key provider missing its key", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.ApiKeyCredentialProviders = []*AwsBedrockAgentCoreApiKeyProvider{
					{Name: "docs-api"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate api-key provider names", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.ApiKeyCredentialProviders = []*AwsBedrockAgentCoreApiKeyProvider{
					{Name: "docs-api", ApiKey: "k1"},
					{Name: "docs-api", ApiKey: "k2"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an unknown oauth vendor", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.Oauth2CredentialProviders = []*AwsBedrockAgentCoreOauth2Provider{
					{Name: "x", Vendor: "OKTA", ClientId: "c", ClientSecret: "s"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a CUSTOM vendor missing discovery", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.Oauth2CredentialProviders = []*AwsBedrockAgentCoreOauth2Provider{
					{Name: "x", Vendor: "CUSTOM", ClientId: "c", ClientSecret: "s"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a well-known vendor carrying discovery", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.Oauth2CredentialProviders = []*AwsBedrockAgentCoreOauth2Provider{
					{
						Name:         "gh",
						Vendor:       "GITHUB",
						ClientId:     "c",
						ClientSecret: "s",
						OauthDiscovery: &AwsBedrockAgentCoreOauth2Discovery{
							DiscoveryUrl: "https://x.example.com/.well-known/openid-configuration",
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a discovery block setting both arms", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.Oauth2CredentialProviders = []*AwsBedrockAgentCoreOauth2Provider{
					{
						Name:         "x",
						Vendor:       "CUSTOM",
						ClientId:     "c",
						ClientSecret: "s",
						OauthDiscovery: &AwsBedrockAgentCoreOauth2Discovery{
							DiscoveryUrl: "https://x.example.com/.well-known/openid-configuration",
							AuthorizationServerMetadata: &AwsBedrockAgentCoreOauth2ServerMetadata{
								Issuer:                "https://x.example.com",
								AuthorizationEndpoint: "https://x.example.com/authorize",
								TokenEndpoint:         "https://x.example.com/token",
							},
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with a hyphenated policy engine name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.PolicyEngine = &AwsBedrockAgentCorePolicyEngine{EngineName: "agent-authz"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with duplicate cedar policy names", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.PolicyEngine = &AwsBedrockAgentCorePolicyEngine{
					EngineName: "authz",
					Policies: []*AwsBedrockAgentCoreCedarPolicy{
						{Name: "p1", CedarStatement: "permit(principal, action, resource);"},
						{Name: "p1", CedarStatement: "forbid(principal, action, resource);"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("with an unknown validation mode", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalIdentity()
				spec.PolicyEngine = &AwsBedrockAgentCorePolicyEngine{
					EngineName: "authz",
					Policies: []*AwsBedrockAgentCoreCedarPolicy{
						{Name: "p1", CedarStatement: "permit(principal, action, resource);", ValidationMode: "WARN"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
