package kubernetesserviceaccountv1alpha1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	kubernetes "github.com/plantonhq/planton/catalog/kubernetes"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func TestKubernetesServiceAccountSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesServiceAccountSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesServiceAccountSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal valid spec", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "app-identity",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a name with dots (DNS subdomain)", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "my.dotted.identity",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace provided as a literal value", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "app-identity",
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "prod"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts image pull secrets as literal names and references", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "puller-identity",
				ImagePullSecrets: []*foreignkeyv1.StringValueOrRef{
					{
						LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "registry-creds"},
					},
					{
						LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
							ValueFrom: &foreignkeyv1.ValueFromRef{Name: "ghcr-credential"},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts automount_service_account_token left unset", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "default-mount-identity",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts automount_service_account_token set to true", func() {
			spec := &KubernetesServiceAccountSpec{
				Name:                         "explicit-mount-identity",
				AutomountServiceAccountToken: boolPtr(true),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts automount_service_account_token set to false", func() {
			spec := &KubernetesServiceAccountSpec{
				Name:                         "hardened-identity",
				AutomountServiceAccountToken: boolPtr(false),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts GKE workload identity with a service account email reference", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "gke-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Gke{
						Gke: &kubernetes.KubernetesWorkloadIdentityGke{
							ServiceAccountEmail: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
									ValueFrom: &foreignkeyv1.ValueFromRef{Name: "dns-manager-gsa"},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts EKS workload identity with a literal role ARN", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "eks-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Eks{
						Eks: &kubernetes.KubernetesWorkloadIdentityEksIrsa{
							RoleArn: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
									Value: "arn:aws:iam::123456789012:role/dns-manager",
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts AKS workload identity with a literal client ID", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "aks-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Aks{
						Aks: &kubernetes.KubernetesWorkloadIdentityAks{
							ClientId: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
									Value: "11111111-2222-3333-4444-555555555555",
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts AKS workload identity with a valid tenant ID UUID", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "aks-cross-tenant-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Aks{
						Aks: &kubernetes.KubernetesWorkloadIdentityAks{
							ClientId: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
									Value: "11111111-2222-3333-4444-555555555555",
								},
							},
							TenantId: stringPtr("99999999-8888-7777-6666-555555555555"),
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects an empty name", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name with uppercase letters", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "AppIdentity",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name starting with a dot", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: ".hidden-identity",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name ending with a hyphen", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "app-identity-",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a name longer than 253 characters", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: strings.Repeat("a", 254),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an empty image pull secret entry", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "puller-identity",
				ImagePullSecrets: []*foreignkeyv1.StringValueOrRef{
					{},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects GKE workload identity without a service account email", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "gke-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Gke{
						Gke: &kubernetes.KubernetesWorkloadIdentityGke{},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects GKE workload identity with an empty service account email message", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "gke-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Gke{
						Gke: &kubernetes.KubernetesWorkloadIdentityGke{
							ServiceAccountEmail: &foreignkeyv1.StringValueOrRef{},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects EKS workload identity without a role ARN", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "eks-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Eks{
						Eks: &kubernetes.KubernetesWorkloadIdentityEksIrsa{},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects AKS workload identity without a client ID", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "aks-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Aks{
						Aks: &kubernetes.KubernetesWorkloadIdentityAks{},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects AKS workload identity with a tenant ID that is not a UUID", func() {
			spec := &KubernetesServiceAccountSpec{
				Name: "aks-identity",
				WorkloadIdentity: &kubernetes.KubernetesWorkloadIdentity{
					Provider: &kubernetes.KubernetesWorkloadIdentity_Aks{
						Aks: &kubernetes.KubernetesWorkloadIdentityAks{
							ClientId: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
									Value: "11111111-2222-3333-4444-555555555555",
								},
							},
							TenantId: stringPtr("not-a-uuid"),
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
