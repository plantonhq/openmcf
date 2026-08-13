package awsrestapidomainv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsRestApiDomainSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRestApiDomainSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalDomain is the smallest valid domain: name + ACM certificate.
func minimalDomain() *AwsRestApiDomainSpec {
	return &AwsRestApiDomainSpec{
		Region:         "us-west-2",
		DomainName:     "api.example.com",
		CertificateArn: svr("arn:aws:acm:us-west-2:123456789012:certificate/abc"),
	}
}

var _ = ginkgo.Describe("AwsRestApiDomainSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalDomain())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the full surface configured", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalDomain()
				spec.EndpointConfiguration = &AwsRestApiDomainEndpointConfiguration{
					Type:          "REGIONAL",
					IpAddressType: "dualstack",
				}
				spec.EndpointAccessMode = "STRICT"
				spec.SecurityPolicy = "SecurityPolicy_TLS13_1_2_2021_06"
				spec.MutualTls = &AwsRestApiDomainMutualTls{
					TruststoreUri:     "s3://trust-bucket/truststore.pem",
					TruststoreVersion: "3",
				}
				spec.OwnershipVerificationCertificateArn = svr("arn:aws:acm:us-west-2:123456789012:certificate/own")
				spec.RoutingMode = "BASE_PATH_MAPPING_ONLY"
				spec.BasePathMappings = []*AwsRestApiDomainBasePathMapping{
					{
						BasePath:  "orders",
						RestApiId: svr("abc123"),
						StageName: svr("prod"),
					},
					{
						RestApiId: svr("def456"),
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a PRIVATE domain and access associations", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalDomain()
				spec.EndpointConfiguration = &AwsRestApiDomainEndpointConfiguration{Type: "PRIVATE"}
				spec.AccessAssociations = []*AwsRestApiDomainAccessAssociation{
					{VpcEndpointId: svr("vpce-abc123")},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an uploaded certificate", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalDomain()
				spec.CertificateArn = nil
				spec.UploadedCertificate = &AwsRestApiDomainUploadedCertificate{
					Name:       "legacy-cert",
					Body:       "-----BEGIN CERTIFICATE-----...",
					PrivateKey: "-----BEGIN PRIVATE KEY-----...",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a domain without any certificate", func() {
			spec := minimalDomain()
			spec.CertificateArn = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of certificate_arn"))
		})

		ginkgo.It("rejects a domain with both certificate sources", func() {
			spec := minimalDomain()
			spec.UploadedCertificate = &AwsRestApiDomainUploadedCertificate{
				Name:       "legacy",
				Body:       "cert",
				PrivateKey: "key",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of certificate_arn"))
		})

		ginkgo.It("rejects duplicate base paths", func() {
			spec := minimalDomain()
			spec.BasePathMappings = []*AwsRestApiDomainBasePathMapping{
				{BasePath: "orders", RestApiId: svr("a")},
				{BasePath: "orders", RestApiId: svr("b")},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("unique base_path values"))
		})

		ginkgo.It("rejects a base path containing a slash", func() {
			spec := minimalDomain()
			spec.BasePathMappings = []*AwsRestApiDomainBasePathMapping{
				{BasePath: "orders/v1", RestApiId: svr("a")},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects access associations on a non-PRIVATE domain", func() {
			spec := minimalDomain()
			spec.AccessAssociations = []*AwsRestApiDomainAccessAssociation{
				{VpcEndpointId: svr("vpce-abc")},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("access_associations apply only"))
		})

		ginkgo.It("rejects a PRIVATE endpoint forced to ipv4", func() {
			spec := minimalDomain()
			spec.EndpointConfiguration = &AwsRestApiDomainEndpointConfiguration{
				Type:          "PRIVATE",
				IpAddressType: "ipv4",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("PRIVATE endpoints require ip_address_type"))
		})

		ginkgo.It("rejects a truststore that is not an S3 URI", func() {
			spec := minimalDomain()
			spec.MutualTls = &AwsRestApiDomainMutualTls{
				TruststoreUri: "https://trust.example.com/store.pem",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an uploaded certificate without its private key", func() {
			spec := minimalDomain()
			spec.CertificateArn = nil
			spec.UploadedCertificate = &AwsRestApiDomainUploadedCertificate{
				Name: "legacy",
				Body: "cert",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid routing mode", func() {
			spec := minimalDomain()
			spec.RoutingMode = "ROUND_ROBIN"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
