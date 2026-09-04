package kubernetesplantonplatformv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestKubernetesPlantonPlatformSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesPlantonPlatformSpec Validation Tests")
}

func literalRef(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

// minimalValidPlatform is the zero-config case: a namespace and a version
// — the whole platform boots from exactly this, reachable over
// port-forward with the first console visitor becoming the admin.
func minimalValidPlatform() *KubernetesPlantonPlatform {
	return &KubernetesPlantonPlatform{
		ApiVersion: "kubernetes.planton.dev/v1alpha1",
		Kind:       "KubernetesPlantonPlatform",
		Metadata: &shared.CloudResourceMetadata{
			Name: "planton",
		},
		Spec: &KubernetesPlantonPlatformSpec{
			Namespace:       literalRef("planton"),
			CreateNamespace: true,
			Version:         "v0.0.45",
		},
	}
}

var _ = ginkgo.Describe("KubernetesPlantonPlatformSpec Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("should not return a validation error for the zero-config platform", func() {
			err := protovalidate.Validate(minimalValidPlatform())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept platform-wide storage settings", func() {
			input := minimalValidPlatform()
			input.Spec.Storage = &KubernetesPlantonPlatformStorage{
				StorageClassName: "gp3",
				Size:             "800Gi",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept an ingress with hostname and cert-manager TLS", func() {
			input := minimalValidPlatform()
			input.Spec.Ingress = &KubernetesPlantonPlatformIngress{
				Enabled:  true,
				Hostname: "planton.example.com",
				Tls: &KubernetesPlantonPlatformIngressTls{
					Issuer: &KubernetesPlantonPlatformCertManagerIssuer{
						Name: "letsencrypt",
					},
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept an ingress with a brought certificate Secret", func() {
			input := minimalValidPlatform()
			input.Spec.Ingress = &KubernetesPlantonPlatformIngress{
				Enabled:  true,
				Hostname: "planton.example.com",
				Tls: &KubernetesPlantonPlatformIngressTls{
					SecretName: "planton-tls",
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept a magic-DNS ingress (enabled, no hostname, no tls)", func() {
			input := minimalValidPlatform()
			input.Spec.Ingress = &KubernetesPlantonPlatformIngress{Enabled: true}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept runner workload identity and database growth", func() {
			input := minimalValidPlatform()
			replicas := int32(2)
			input.Spec.Runner = &KubernetesPlantonPlatformRunner{
				ServiceAccountAnnotations: map[string]string{
					"eks.amazonaws.com/role-arn": "arn:aws:iam::123456789012:role/planton-runner",
				},
			}
			input.Spec.Database = &KubernetesPlantonPlatformDatabase{
				Postgresql: &KubernetesPlantonPlatformPostgresql{
					Replicas:    &replicas,
					StorageSize: "20Gi",
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept the AWS secret backend with its config", func() {
			input := minimalValidPlatform()
			input.Spec.Bootstrap = &KubernetesPlantonPlatformBootstrap{
				SecretBackend: &KubernetesPlantonPlatformSecretBackend{
					Type: "awsSecretsManager",
					AwsSecretsManager: &KubernetesPlantonPlatformAwsSecretsManager{
						Region:    "us-east-1",
						KmsKeyArn: "arn:aws:kms:us-east-1:123456789012:key/abc",
					},
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept a license from a Secret reference", func() {
			input := minimalValidPlatform()
			input.Spec.License = &KubernetesPlantonPlatformLicense{
				SecretKeyRef: &KubernetesPlantonPlatformLicenseSecretKeyRef{
					Name: "planton-license",
					Key:  "key",
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("should accept a pre-seeded admin email", func() {
			input := minimalValidPlatform()
			input.Spec.Identity = &KubernetesPlantonPlatformIdentity{
				AdminEmail: "admin@example.com",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("should fail when version is missing", func() {
			input := minimalValidPlatform()
			input.Spec.Version = ""
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail when namespace is missing", func() {
			input := minimalValidPlatform()
			input.Spec.Namespace = nil
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on TLS without a hostname", func() {
			input := minimalValidPlatform()
			input.Spec.Ingress = &KubernetesPlantonPlatformIngress{
				Enabled: true,
				Tls: &KubernetesPlantonPlatformIngressTls{
					SecretName: "planton-tls",
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on TLS with BOTH a Secret and an issuer", func() {
			input := minimalValidPlatform()
			input.Spec.Ingress = &KubernetesPlantonPlatformIngress{
				Enabled:  true,
				Hostname: "planton.example.com",
				Tls: &KubernetesPlantonPlatformIngressTls{
					SecretName: "planton-tls",
					Issuer: &KubernetesPlantonPlatformCertManagerIssuer{
						Name: "letsencrypt",
					},
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on TLS with NEITHER a Secret nor an issuer", func() {
			input := minimalValidPlatform()
			input.Spec.Ingress = &KubernetesPlantonPlatformIngress{
				Enabled:  true,
				Hostname: "planton.example.com",
				Tls:      &KubernetesPlantonPlatformIngressTls{},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on the AWS secret backend without its config", func() {
			input := minimalValidPlatform()
			input.Spec.Bootstrap = &KubernetesPlantonPlatformBootstrap{
				SecretBackend: &KubernetesPlantonPlatformSecretBackend{
					Type: "awsSecretsManager",
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on a license with BOTH key and secret reference", func() {
			input := minimalValidPlatform()
			input.Spec.License = &KubernetesPlantonPlatformLicense{
				Key: "plk_FAKE_PLACEHOLDER_VALUE",
				SecretKeyRef: &KubernetesPlantonPlatformLicenseSecretKeyRef{
					Name: "planton-license",
					Key:  "key",
				},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on a malformed storage size", func() {
			input := minimalValidPlatform()
			input.Spec.Storage = &KubernetesPlantonPlatformStorage{
				Size: "ten gigs",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on a malformed admin email", func() {
			input := minimalValidPlatform()
			input.Spec.Identity = &KubernetesPlantonPlatformIdentity{
				AdminEmail: "not-an-email",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on an unknown IaC provisioner", func() {
			input := minimalValidPlatform()
			provisioner := "pulumi"
			input.Spec.Bootstrap = &KubernetesPlantonPlatformBootstrap{
				IacProvisioner: &provisioner,
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("should fail on an out-of-range gateway port", func() {
			input := minimalValidPlatform()
			port := int32(70000)
			input.Spec.Gateway = &KubernetesPlantonPlatformGateway{
				LocalPort: &port,
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
