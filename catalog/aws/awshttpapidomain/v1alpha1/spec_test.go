package awshttpapidomainv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsHttpApiDomainSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsHttpApiDomainSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

var _ = ginkgo.Describe("AwsHttpApiDomainSpec validations", func() {
	var spec *AwsHttpApiDomainSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: region + domain + certificate.
		spec = &AwsHttpApiDomainSpec{
			Region:         "us-west-2",
			DomainName:     "api.example.com",
			CertificateArn: strRef("arn:aws:acm:us-west-2:123456789012:certificate/abc-123"),
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal spec (region + domain + certificate)", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a wildcard domain", func() {
		spec.DomainName = "*.example.com"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts API mappings with distinct keys", func() {
		spec.ApiMappings = []*AwsHttpApiDomainApiMapping{
			{
				ApiId: strRef("api-abc123"),
				Stage: "$default",
			},
			{
				ApiId:         strRef("api-def456"),
				Stage:         "$default",
				ApiMappingKey: "orders",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts mutual TLS with a pinned truststore version", func() {
		spec.MutualTls = &AwsHttpApiDomainMutualTls{
			TruststoreUri:     "s3://my-bucket/truststore.pem",
			TruststoreVersion: "abc123version",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts ip_address_type ipv4", func() {
		spec.IpAddressType = "ipv4"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Required field validations
	// -------------------------------------------------------------------------

	ginkgo.It("fails when region is empty", func() {
		spec.Region = ""
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when domain_name is empty", func() {
		spec.DomainName = ""
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when domain_name contains uppercase characters", func() {
		spec.DomainName = "API.Example.com"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when certificate_arn is missing", func() {
		spec.CertificateArn = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: ip_address_type_valid
	// -------------------------------------------------------------------------

	ginkgo.It("fails when ip_address_type is invalid", func() {
		spec.IpAddressType = "ipv6"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: api_mapping_keys_unique
	// -------------------------------------------------------------------------

	ginkgo.It("fails when two mappings share a path key", func() {
		spec.ApiMappings = []*AwsHttpApiDomainApiMapping{
			{ApiId: strRef("api-abc123"), Stage: "$default", ApiMappingKey: "orders"},
			{ApiId: strRef("api-def456"), Stage: "$default", ApiMappingKey: "orders"},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when two mappings both claim the domain root", func() {
		spec.ApiMappings = []*AwsHttpApiDomainApiMapping{
			{ApiId: strRef("api-abc123"), Stage: "$default"},
			{ApiId: strRef("api-def456"), Stage: "$default"},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// API mapping field validations
	// -------------------------------------------------------------------------

	ginkgo.It("fails when a mapping omits api_id", func() {
		spec.ApiMappings = []*AwsHttpApiDomainApiMapping{
			{Stage: "$default"},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a mapping omits stage", func() {
		spec.ApiMappings = []*AwsHttpApiDomainApiMapping{
			{ApiId: strRef("api-abc123")},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a mapping key contains a slash", func() {
		spec.ApiMappings = []*AwsHttpApiDomainApiMapping{
			{ApiId: strRef("api-abc123"), Stage: "$default", ApiMappingKey: "v1/orders"},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Mutual TLS validations
	// -------------------------------------------------------------------------

	ginkgo.It("fails when truststore_uri is not an S3 URI", func() {
		spec.MutualTls = &AwsHttpApiDomainMutualTls{
			TruststoreUri: "https://example.com/truststore.pem",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when mutual_tls is present with an empty truststore_uri", func() {
		spec.MutualTls = &AwsHttpApiDomainMutualTls{}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Ownership verification certificate
	// -------------------------------------------------------------------------

	ginkgo.It("accepts an ownership verification certificate ARN", func() {
		spec.OwnershipVerificationCertificateArn = strRef("arn:aws:acm:us-west-2:123456789012:certificate/own-456")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Routing mode and routing rules
	// -------------------------------------------------------------------------

	// helper: one valid base-path rule.
	validRule := func(priority int32) *AwsHttpApiDomainRoutingRule {
		return &AwsHttpApiDomainRoutingRule{
			Priority: priority,
			Conditions: []*AwsHttpApiDomainRoutingRuleCondition{
				{BasePaths: []string{"orders"}},
			},
			ApiId: strRef("api-abc123"),
			Stage: "$default",
		}
	}

	ginkgo.It("accepts routing rules under ROUTING_RULE_THEN_API_MAPPING", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{validRule(10)}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a header-match rule under ROUTING_RULE_ONLY", func() {
		spec.RoutingMode = "ROUTING_RULE_ONLY"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{
			{
				Priority: 5,
				Conditions: []*AwsHttpApiDomainRoutingRuleCondition{
					{Header: &AwsHttpApiDomainRoutingRuleHeaderMatch{Name: "x-tenant-id", ValueGlob: "tenant-a-*"}},
				},
				ApiId:         strRef("api-abc123"),
				Stage:         "$default",
				StripBasePath: true,
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a rule combining a base-path and a header condition", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		rule := validRule(1)
		rule.Conditions = append(rule.Conditions, &AwsHttpApiDomainRoutingRuleCondition{
			Header: &AwsHttpApiDomainRoutingRuleHeaderMatch{Name: "x-env", ValueGlob: "beta"},
		})
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{rule}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts explicit API_MAPPING_ONLY with no rules", func() {
		spec.RoutingMode = "API_MAPPING_ONLY"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("fails when routing_mode is invalid", func() {
		spec.RoutingMode = "RULES_FIRST"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when rules are set without a rule-honoring routing_mode", func() {
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{validRule(10)}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when rules are set under API_MAPPING_ONLY", func() {
		spec.RoutingMode = "API_MAPPING_ONLY"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{validRule(10)}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when api_mappings are combined with a rule-honoring mode", func() {
		spec.ApiMappings = []*AwsHttpApiDomainApiMapping{
			{ApiId: strRef("api-abc123"), Stage: "$default"},
		}
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{validRule(10)}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts api_mappings under explicit API_MAPPING_ONLY", func() {
		spec.ApiMappings = []*AwsHttpApiDomainApiMapping{
			{ApiId: strRef("api-abc123"), Stage: "$default"},
		}
		spec.RoutingMode = "API_MAPPING_ONLY"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("fails when a rule-honoring mode has no rules", func() {
		spec.RoutingMode = "ROUTING_RULE_ONLY"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when two rules share a priority", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{validRule(10), validRule(10)}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when priority is zero", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{validRule(0)}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when priority exceeds 1,000,000", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{validRule(1000001)}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a rule has no conditions", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		rule := validRule(10)
		rule.Conditions = nil
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{rule}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a condition sets both base_paths and header", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		rule := validRule(10)
		rule.Conditions = []*AwsHttpApiDomainRoutingRuleCondition{
			{
				BasePaths: []string{"orders"},
				Header:    &AwsHttpApiDomainRoutingRuleHeaderMatch{Name: "x-env", ValueGlob: "beta"},
			},
		}
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{rule}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a condition sets neither base_paths nor header", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		rule := validRule(10)
		rule.Conditions = []*AwsHttpApiDomainRoutingRuleCondition{{}}
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{rule}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a rule omits api_id", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		rule := validRule(10)
		rule.ApiId = nil
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{rule}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a rule omits stage", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		rule := validRule(10)
		rule.Stage = ""
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{rule}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a header match name exceeds 40 characters", func() {
		spec.RoutingMode = "ROUTING_RULE_ONLY"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{
			{
				Priority: 5,
				Conditions: []*AwsHttpApiDomainRoutingRuleCondition{
					{Header: &AwsHttpApiDomainRoutingRuleHeaderMatch{
						Name:      "x-this-header-name-is-way-over-forty-characters-long",
						ValueGlob: "v",
					}},
				},
				ApiId: strRef("api-abc123"),
				Stage: "$default",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a header match glob exceeds 128 characters", func() {
		longGlob := make([]byte, 129)
		for i := range longGlob {
			longGlob[i] = 'a'
		}
		spec.RoutingMode = "ROUTING_RULE_ONLY"
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{
			{
				Priority: 5,
				Conditions: []*AwsHttpApiDomainRoutingRuleCondition{
					{Header: &AwsHttpApiDomainRoutingRuleHeaderMatch{
						Name:      "x-env",
						ValueGlob: string(longGlob),
					}},
				},
				ApiId: strRef("api-abc123"),
				Stage: "$default",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a base path entry is empty", func() {
		spec.RoutingMode = "ROUTING_RULE_THEN_API_MAPPING"
		rule := validRule(10)
		rule.Conditions = []*AwsHttpApiDomainRoutingRuleCondition{{BasePaths: []string{""}}}
		spec.RoutingRules = []*AwsHttpApiDomainRoutingRule{rule}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})
})
