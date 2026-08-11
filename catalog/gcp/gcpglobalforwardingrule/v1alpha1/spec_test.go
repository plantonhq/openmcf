package gcpglobalforwardingrulev1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestSuite(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "GcpGlobalForwardingRuleSpec Suite")
}

func literalRef(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func proxySelfLink() *foreignkeyv1.StringValueOrRef {
	return literalRef("https://www.googleapis.com/compute/v1/projects/p/global/targetHttpsProxies/web-frontend")
}

func schemePtr(scheme string) *string {
	return &scheme
}

var _ = ginkgo.Describe("GcpGlobalForwardingRuleSpec", func() {
	var validator protovalidate.Validator

	ginkgo.BeforeEach(func() {
		var err error
		validator, err = protovalidate.New()
		gomega.Expect(err).ToNot(gomega.HaveOccurred())
	})

	minimal := func() *GcpGlobalForwardingRule {
		return &GcpGlobalForwardingRule{
			ApiVersion: "gcp.planton.dev/v1alpha1",
			Kind:       "GcpGlobalForwardingRule",
			Metadata: &shared.CloudResourceMetadata{
				Name: "test-forwarding-rule",
			},
			Spec: &GcpGlobalForwardingRuleSpec{
				Target:    proxySelfLink(),
				PortRange: "443",
			},
		}
	}

	// ──────────────── Positive Cases ────────────────

	ginkgo.It("should accept a minimal external HTTPS frontend", func() {
		gomega.Expect(validator.Validate(minimal())).To(gomega.Succeed())
	})

	ginkgo.It("should accept a static ip_address reference", func() {
		target := minimal()
		target.Spec.IpAddress = literalRef("34.120.1.2")
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept a port range", func() {
		target := minimal()
		target.Spec.PortRange = "8080-8090"
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept an EXTERNAL_MANAGED frontend", func() {
		target := minimal()
		target.Spec.LoadBalancingScheme = schemePtr("EXTERNAL_MANAGED")
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept an INTERNAL_SELF_MANAGED frontend with network and metadata filters", func() {
		target := minimal()
		target.Spec.LoadBalancingScheme = schemePtr("INTERNAL_SELF_MANAGED")
		target.Spec.Network = literalRef("https://www.googleapis.com/compute/v1/projects/p/global/networks/mesh-vpc")
		target.Spec.MetadataFilters = []*GcpGlobalForwardingRuleMetadataFilter{
			{
				FilterMatchCriteria: "MATCH_ANY",
				FilterLabels: []*GcpGlobalForwardingRuleMetadataFilterLabel{
					{Name: "env", Value: "prod"},
				},
			},
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept a Private Service Connect Google-APIs frontend", func() {
		target := minimal()
		target.Spec.Target = literalRef("all-apis")
		target.Spec.LoadBalancingScheme = schemePtr("NONE")
		target.Spec.PortRange = ""
		target.Spec.IpAddress = literalRef("10.10.0.5")
		target.Spec.Network = literalRef("https://www.googleapis.com/compute/v1/projects/p/global/networks/main-vpc")
		target.Spec.NoAutomateDnsZone = true
		target.Spec.ServiceDirectoryRegistration = &GcpGlobalForwardingRuleServiceDirectoryRegistration{
			Namespace:              "psc-endpoints",
			ServiceDirectoryRegion: "us-central1",
		}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept the PREMIUM network tier", func() {
		target := minimal()
		target.Spec.NetworkTier = "PREMIUM"
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept labels", func() {
		target := minimal()
		target.Spec.Labels = map[string]string{"env": "prod", "team": "platform"}
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept the backend-bucket migration canary states", func() {
		for _, state := range []string{"PREPARE", "TEST_BY_PERCENTAGE", "TEST_ALL_TRAFFIC"} {
			target := minimal()
			target.Spec.ExternalManagedBackendBucketMigrationState = state
			gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
		}
	})

	ginkgo.It("should accept a migration testing percentage with TEST_BY_PERCENTAGE", func() {
		target := minimal()
		target.Spec.ExternalManagedBackendBucketMigrationState = "TEST_BY_PERCENTAGE"
		target.Spec.ExternalManagedBackendBucketMigrationTestingPercentage = 25.5
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept each ip_protocol", func() {
		for _, protocol := range []string{"TCP", "UDP", "ESP", "AH", "SCTP", "ICMP"} {
			target := minimal()
			target.Spec.IpProtocol = &protocol
			gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
		}
	})

	ginkgo.It("should accept IPV6 for an auto-assigned address", func() {
		target := minimal()
		target.Spec.IpVersion = "IPV6"
		gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
	})

	ginkgo.It("should accept each deletion_policy value", func() {
		for _, v := range []string{"DELETE", "PREVENT", "ABANDON"} {
			target := minimal()
			target.Spec.DeletionPolicy = v
			gomega.Expect(validator.Validate(target)).To(gomega.Succeed())
		}
	})

	// ──────────────── Negative Cases ────────────────

	ginkgo.It("should reject an invalid deletion_policy", func() {
		target := minimal()
		target.Spec.DeletionPolicy = "KEEP"
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject a missing target", func() {
		target := minimal()
		target.Spec.Target = nil
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject an invalid forwarding_rule_name", func() {
		target := minimal()
		target.Spec.ForwardingRuleName = "Frontend-1"
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "RFC1035")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject an invalid ip_protocol", func() {
		target := minimal()
		invalid := "GRE"
		target.Spec.IpProtocol = &invalid
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject an invalid ip_version", func() {
		target := minimal()
		target.Spec.IpVersion = "IPV5"
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject an invalid load_balancing_scheme", func() {
		target := minimal()
		target.Spec.LoadBalancingScheme = schemePtr("INTERNAL")
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject a malformed port_range", func() {
		target := minimal()
		target.Spec.PortRange = "443,80"
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "contiguous range")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject the STANDARD network tier on a global rule", func() {
		target := minimal()
		target.Spec.NetworkTier = "STANDARD"
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "PREMIUM")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject network on an external frontend", func() {
		target := minimal()
		target.Spec.Network = literalRef("https://www.googleapis.com/compute/v1/projects/p/global/networks/main-vpc")
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "internal or Private Service Connect")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject metadata_filters without Traffic Director", func() {
		target := minimal()
		target.Spec.MetadataFilters = []*GcpGlobalForwardingRuleMetadataFilter{
			{
				FilterMatchCriteria: "MATCH_ALL",
				FilterLabels: []*GcpGlobalForwardingRuleMetadataFilterLabel{
					{Name: "env", Value: "prod"},
				},
			},
		}
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "Traffic Director")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject service_directory_registration without the PSC scheme", func() {
		target := minimal()
		target.Spec.ServiceDirectoryRegistration = &GcpGlobalForwardingRuleServiceDirectoryRegistration{
			Namespace: "psc-endpoints",
		}
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject no_automate_dns_zone without the PSC scheme", func() {
		target := minimal()
		target.Spec.NoAutomateDnsZone = true
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject a migration percentage without TEST_BY_PERCENTAGE", func() {
		target := minimal()
		target.Spec.ExternalManagedBackendBucketMigrationTestingPercentage = 50
		err := validator.Validate(target)
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(strings.Contains(err.Error(), "TEST_BY_PERCENTAGE")).To(gomega.BeTrue())
	})

	ginkgo.It("should reject a migration percentage above 100", func() {
		target := minimal()
		target.Spec.ExternalManagedBackendBucketMigrationState = "TEST_BY_PERCENTAGE"
		target.Spec.ExternalManagedBackendBucketMigrationTestingPercentage = 101
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject an invalid migration state", func() {
		target := minimal()
		target.Spec.ExternalManagedBackendBucketMigrationState = "FINISH"
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject a metadata filter with an invalid criteria", func() {
		target := minimal()
		target.Spec.LoadBalancingScheme = schemePtr("INTERNAL_SELF_MANAGED")
		target.Spec.MetadataFilters = []*GcpGlobalForwardingRuleMetadataFilter{
			{
				FilterMatchCriteria: "MATCH_SOME",
				FilterLabels: []*GcpGlobalForwardingRuleMetadataFilterLabel{
					{Name: "env", Value: "prod"},
				},
			},
		}
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject a metadata filter with no labels", func() {
		target := minimal()
		target.Spec.LoadBalancingScheme = schemePtr("INTERNAL_SELF_MANAGED")
		target.Spec.MetadataFilters = []*GcpGlobalForwardingRuleMetadataFilter{
			{FilterMatchCriteria: "MATCH_ALL"},
		}
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject a wrong kind literal", func() {
		target := minimal()
		target.Kind = "GcpForwardingRule"
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})

	ginkgo.It("should reject missing spec", func() {
		target := minimal()
		target.Spec = nil
		gomega.Expect(validator.Validate(target)).ToNot(gomega.Succeed())
	})
})
