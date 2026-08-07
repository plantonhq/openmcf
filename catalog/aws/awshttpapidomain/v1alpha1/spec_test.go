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
})
