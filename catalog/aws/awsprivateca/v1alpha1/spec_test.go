package awsprivatecav1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsPrivateCaSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsPrivateCaSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

const testCsr = "-----BEGIN CERTIFICATE REQUEST-----\nMIIB...\n-----END CERTIFICATE REQUEST-----\n"

func minimalConfig() *AwsPrivateCaSpec {
	return &AwsPrivateCaSpec{
		Region:           "us-west-2",
		Type:             "ROOT",
		KeyAlgorithm:     "RSA_2048",
		SigningAlgorithm: "SHA256WITHRSA",
		Subject: &AwsPrivateCaSubject{
			CommonName: "corp-root-ca",
		},
	}
}

var _ = ginkgo.Describe("AwsPrivateCaSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal root CA", func() {
			gomega.Expect(protovalidate.Validate(minimalConfig())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an EC key with an ECDSA signing algorithm", func() {
			spec := minimalConfig()
			spec.KeyAlgorithm = "EC_prime256v1"
			spec.SigningAlgorithm = "SHA256WITHECDSA"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a post-quantum key with its exact signing level", func() {
			spec := minimalConfig()
			spec.KeyAlgorithm = "ML_DSA_65"
			spec.SigningAlgorithm = "ML_DSA_65"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts CRL and OCSP revocation with a root validity", func() {
			spec := minimalConfig()
			spec.Revocation = &AwsPrivateCaRevocation{
				Crl: &AwsPrivateCaCrl{
					Enabled:          true,
					ExpirationInDays: 7,
					S3BucketName:     literal("corp-ca-crl"),
					S3ObjectAcl:      "BUCKET_OWNER_FULL_CONTROL",
					CustomCname:      "crl.corp.example.com",
				},
				Ocsp: &AwsPrivateCaOcsp{Enabled: true},
			}
			spec.RootCaValidity = &AwsPrivateCaValidity{Type: "YEARS", Value: "10"}
			spec.PermanentDeletionTimeInDays = 7
			spec.AcmRenewalPermission = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a subordinate activated from a parent CA", func() {
			spec := minimalConfig()
			spec.Type = "SUBORDINATE"
			spec.SubordinateActivation = &AwsPrivateCaSubordinateActivation{
				ParentCaArn: literal("arn:aws:acm-pca:us-west-2:111111111111:certificate-authority/1234abcd-12ab-34cd-56ef-1234567890ab"),
				PathLength:  0,
				Validity:    &AwsPrivateCaValidity{Type: "YEARS", Value: "5"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts issued certificates with an END_DATE validity", func() {
			spec := minimalConfig()
			spec.IssuedCertificates = []*AwsPrivateCaIssuedCertificate{
				{
					Name:             "orders-mtls",
					Csr:              testCsr,
					SigningAlgorithm: "SHA256WITHRSA",
					Validity:         &AwsPrivateCaValidity{Type: "END_DATE", Value: "2036-01-01T00:00:00Z"},
					TemplateArn:      "arn:aws:acm-pca:::template/EndEntityClientAuthCertificate/V1",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts short-lived mode without revocation", func() {
			spec := minimalConfig()
			spec.UsageMode = "SHORT_LIVED_CERTIFICATE"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an RSA key with an ECDSA signing algorithm", func() {
			spec := minimalConfig()
			spec.SigningAlgorithm = "SHA256WITHECDSA"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a post-quantum key with a mismatched signing level", func() {
			spec := minimalConfig()
			spec.KeyAlgorithm = "ML_DSA_44"
			spec.SigningAlgorithm = "ML_DSA_87"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty subject", func() {
			spec := minimalConfig()
			spec.Subject = &AwsPrivateCaSubject{}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a root validity on a subordinate", func() {
			spec := minimalConfig()
			spec.Type = "SUBORDINATE"
			spec.RootCaValidity = &AwsPrivateCaValidity{Type: "YEARS", Value: "10"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a subordinate activation on a root", func() {
			spec := minimalConfig()
			spec.SubordinateActivation = &AwsPrivateCaSubordinateActivation{
				ParentCaArn: literal("arn:aws:acm-pca:us-west-2:111111111111:certificate-authority/1234abcd-12ab-34cd-56ef-1234567890ab"),
				Validity:    &AwsPrivateCaValidity{Type: "YEARS", Value: "5"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects short-lived mode with an enabled CRL", func() {
			spec := minimalConfig()
			spec.UsageMode = "SHORT_LIVED_CERTIFICATE"
			spec.Revocation = &AwsPrivateCaRevocation{
				Crl: &AwsPrivateCaCrl{
					Enabled:          true,
					ExpirationInDays: 7,
					S3BucketName:     literal("corp-ca-crl"),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an enabled CRL without a bucket", func() {
			spec := minimalConfig()
			spec.Revocation = &AwsPrivateCaRevocation{
				Crl: &AwsPrivateCaCrl{Enabled: true, ExpirationInDays: 7},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an enabled CRL without an expiration", func() {
			spec := minimalConfig()
			spec.Revocation = &AwsPrivateCaRevocation{
				Crl: &AwsPrivateCaCrl{Enabled: true, S3BucketName: literal("corp-ca-crl")},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an END_DATE validity carrying a bare number", func() {
			spec := minimalConfig()
			spec.RootCaValidity = &AwsPrivateCaValidity{Type: "END_DATE", Value: "3650"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a YEARS validity carrying a timestamp", func() {
			spec := minimalConfig()
			spec.RootCaValidity = &AwsPrivateCaValidity{Type: "YEARS", Value: "2036-01-01T00:00:00Z"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate issued certificate names", func() {
			spec := minimalConfig()
			cert := &AwsPrivateCaIssuedCertificate{
				Name:             "dup",
				Csr:              testCsr,
				SigningAlgorithm: "SHA256WITHRSA",
				Validity:         &AwsPrivateCaValidity{Type: "DAYS", Value: "398"},
			}
			spec.IssuedCertificates = []*AwsPrivateCaIssuedCertificate{cert, cert}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-PEM CSR", func() {
			spec := minimalConfig()
			spec.IssuedCertificates = []*AwsPrivateCaIssuedCertificate{
				{
					Name:             "bad-csr",
					Csr:              "MIIB-not-pem",
					SigningAlgorithm: "SHA256WITHRSA",
					Validity:         &AwsPrivateCaValidity{Type: "DAYS", Value: "398"},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed template ARN", func() {
			spec := minimalConfig()
			spec.IssuedCertificates = []*AwsPrivateCaIssuedCertificate{
				{
					Name:             "bad-template",
					Csr:              testCsr,
					SigningAlgorithm: "SHA256WITHRSA",
					Validity:         &AwsPrivateCaValidity{Type: "DAYS", Value: "398"},
					TemplateArn:      "arn:aws:acm-pca:us-west-2:111111111111:template/EndEntityCertificate/V1",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range deletion window", func() {
			spec := minimalConfig()
			spec.PermanentDeletionTimeInDays = 3
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range path length", func() {
			spec := minimalConfig()
			spec.Type = "SUBORDINATE"
			spec.SubordinateActivation = &AwsPrivateCaSubordinateActivation{
				ParentCaArn: literal("arn:aws:acm-pca:us-west-2:111111111111:certificate-authority/1234abcd-12ab-34cd-56ef-1234567890ab"),
				PathLength:  4,
				Validity:    &AwsPrivateCaValidity{Type: "YEARS", Value: "5"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
