package kubernetesopenbaov1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	kubernetes "github.com/plantonhq/planton/catalog/kubernetes"
	"github.com/plantonhq/planton/shared"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestKubernetesOpenBao(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesOpenBao Suite")
}

func int32Ptr(i int32) *int32 { return &i }
func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func valueFrom(kind cloudresourcekind.CloudResourceKind, name, fieldPath string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
			ValueFrom: &foreignkeyv1.ValueFromRef{
				Kind:      kind,
				Name:      name,
				FieldPath: fieldPath,
			},
		},
	}
}

var _ = ginkgo.Describe("KubernetesOpenBao Validation Tests", func() {
	var input *KubernetesOpenBao

	ginkgo.BeforeEach(func() {
		input = &KubernetesOpenBao{
			ApiVersion: "kubernetes.planton.dev/v1alpha1",
			Kind:       "KubernetesOpenBao",
			Metadata: &shared.CloudResourceMetadata{
				Name: "openbao",
			},
			Spec: &KubernetesOpenBaoSpec{
				Namespace: literal("openbao"),
			},
		}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.It("a minimal spec (namespace only, chart defaults) should be valid", func() {
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("namespace as a reference should be valid", func() {
			input.Spec.Namespace = valueFrom(cloudresourcekind.CloudResourceKind_KubernetesNamespace, "openbao", "spec.name")
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a maximal spec (every block populated) should be valid", func() {
			input.Spec.CreateNamespace = true
			input.Spec.ChartVersion = strPtr("0.28.6")
			input.Spec.Server = &KubernetesOpenBaoServer{
				Mode: &KubernetesOpenBaoServer_Ha{Ha: &KubernetesOpenBaoHaMode{Replicas: int32Ptr(3)}},
				Resources: &kubernetes.ContainerResources{
					Requests: &kubernetes.CpuMemory{Cpu: "100m", Memory: "256Mi"},
					Limits:   &kubernetes.CpuMemory{Cpu: "1", Memory: "1Gi"},
				},
				DataStorage: &KubernetesOpenBaoStorage{
					Size:         strPtr("20Gi"),
					StorageClass: literal("fast-ssd"),
				},
				AuditStorage: &KubernetesOpenBaoStorage{Size: strPtr("5Gi")},
				LogLevel:     strPtr("debug"),
				LogFormat:    strPtr("json"),
				Scheduling: &KubernetesOpenBaoScheduling{
					NodeSelector: map[string]string{"workload": "secrets"},
				},
			}
			input.Spec.Tls = &KubernetesOpenBaoTls{
				Enabled:        true,
				CertSecretName: literal("openbao-server-tls"),
			}
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_AwsKms{
					AwsKms: &KubernetesOpenBaoAwsKmsSeal{
						Region:   "us-west-2",
						KmsKeyId: "alias/openbao-unseal",
					},
				},
			}
			input.Spec.Injector = &KubernetesOpenBaoInjector{
				Enabled:       true,
				Replicas:      int32Ptr(2),
				FailurePolicy: strPtr("Fail"),
				Resources: &kubernetes.ContainerResources{
					Requests: &kubernetes.CpuMemory{Cpu: "50m", Memory: "64Mi"},
					Limits:   &kubernetes.CpuMemory{Cpu: "250m", Memory: "256Mi"},
				},
			}
			input.Spec.UiEnabled = boolPtr(true)
			input.Spec.NetworkPolicyEnabled = true
			input.Spec.Metrics = &KubernetesOpenBaoMetrics{
				Enabled:               true,
				ServiceMonitorEnabled: true,
			}
			input.Spec.SnapshotAgent = &KubernetesOpenBaoSnapshotAgent{
				Enabled:                 true,
				Schedule:                strPtr("0 */6 * * *"),
				S3Host:                  literal("s3.us-east-1.amazonaws.com"),
				S3Bucket:                "openbao-snapshots",
				S3ExpireDays:            int32Ptr(30),
				S3CredentialsSecretName: "openbao-snapshot-s3-creds",
				BaoRole:                 strPtr("snapshot"),
				BaoAuthPath:             strPtr("kubernetes"),
			}
			input.Spec.ServiceAccount = &KubernetesOpenBaoServiceAccount{
				Annotations:          map[string]string{"eks.amazonaws.com/role-arn": "arn:aws:iam::123456789012:role/openbao"},
				AuthDelegatorEnabled: boolPtr(true),
			}
			input.Spec.HelmValues = "server:\n  extraLabels:\n    team: platform\n"
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("dev mode should be valid", func() {
			input.Spec.Server = &KubernetesOpenBaoServer{
				Mode: &KubernetesOpenBaoServer_Dev{Dev: &KubernetesOpenBaoDevMode{}},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("standalone mode should be valid", func() {
			input.Spec.Server = &KubernetesOpenBaoServer{
				Mode: &KubernetesOpenBaoServer_Standalone{Standalone: &KubernetesOpenBaoStandaloneMode{}},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a single-replica HA cluster (Raft cluster of one) should be valid", func() {
			input.Spec.Server = &KubernetesOpenBaoServer{
				Mode: &KubernetesOpenBaoServer_Ha{Ha: &KubernetesOpenBaoHaMode{Replicas: int32Ptr(1)}},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a TLS cert secret from a Certificate reference should be valid", func() {
			input.Spec.Tls = &KubernetesOpenBaoTls{
				Enabled:        true,
				CertSecretName: valueFrom(cloudresourcekind.CloudResourceKind_KubernetesCertificate, "openbao-cert", "status.outputs.secret_name"),
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a GCP KMS seal should be valid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_GcpKms{
					GcpKms: &KubernetesOpenBaoGcpKmsSeal{
						Project:   literal("my-project"),
						Region:    "global",
						KeyRing:   literal("openbao-ring"),
						CryptoKey: literal("openbao-unseal"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("an Azure Key Vault seal should be valid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_AzureKeyVault{
					AzureKeyVault: &KubernetesOpenBaoAzureKeyVaultSeal{
						VaultName: "openbao-vault",
						KeyName:   "unseal-key",
						TenantId:  "00000000-0000-0000-0000-000000000000",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a transit seal should be valid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_Transit{
					Transit: &KubernetesOpenBaoTransitSeal{
						Address:   "https://bao.example.com:8200",
						KeyName:   "autounseal",
						MountPath: strPtr("transit/"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("metrics without a ServiceMonitor should be valid", func() {
			input.Spec.Metrics = &KubernetesOpenBaoMetrics{Enabled: true}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.It("a missing namespace should be invalid", func() {
			input.Spec.Namespace = nil
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an unknown server log level should be invalid", func() {
			input.Spec.Server = &KubernetesOpenBaoServer{LogLevel: strPtr("verbose")}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an unknown server log format should be invalid", func() {
			input.Spec.Server = &KubernetesOpenBaoServer{LogFormat: strPtr("yaml")}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("zero HA replicas should be invalid", func() {
			input.Spec.Server = &KubernetesOpenBaoServer{
				Mode: &KubernetesOpenBaoServer_Ha{Ha: &KubernetesOpenBaoHaMode{Replicas: int32Ptr(0)}},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("HA replicas above 11 should be invalid", func() {
			input.Spec.Server = &KubernetesOpenBaoServer{
				Mode: &KubernetesOpenBaoServer_Ha{Ha: &KubernetesOpenBaoHaMode{Replicas: int32Ptr(12)}},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a storage size without a unit suffix should be invalid", func() {
			input.Spec.Server = &KubernetesOpenBaoServer{
				DataStorage: &KubernetesOpenBaoStorage{Size: strPtr("10GB")},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("TLS enabled without a cert secret should be invalid", func() {
			input.Spec.Tls = &KubernetesOpenBaoTls{Enabled: true}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("TLS enabled with an empty literal cert secret name should be invalid", func() {
			input.Spec.Tls = &KubernetesOpenBaoTls{Enabled: true, CertSecretName: literal("")}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an AWS KMS seal without a region should be invalid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_AwsKms{
					AwsKms: &KubernetesOpenBaoAwsKmsSeal{KmsKeyId: "alias/openbao-unseal"},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an AWS KMS seal without a key id should be invalid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_AwsKms{
					AwsKms: &KubernetesOpenBaoAwsKmsSeal{Region: "us-west-2"},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a GCP KMS seal without a project should be invalid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_GcpKms{
					GcpKms: &KubernetesOpenBaoGcpKmsSeal{
						Region:    "global",
						KeyRing:   literal("openbao-ring"),
						CryptoKey: literal("openbao-unseal"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a GCP KMS seal without a key ring should be invalid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_GcpKms{
					GcpKms: &KubernetesOpenBaoGcpKmsSeal{
						Project:   literal("my-project"),
						Region:    "global",
						CryptoKey: literal("openbao-unseal"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a GCP KMS seal without a crypto key should be invalid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_GcpKms{
					GcpKms: &KubernetesOpenBaoGcpKmsSeal{
						Project: literal("my-project"),
						Region:  "global",
						KeyRing: literal("openbao-ring"),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an Azure Key Vault seal without a tenant id should be invalid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_AzureKeyVault{
					AzureKeyVault: &KubernetesOpenBaoAzureKeyVaultSeal{
						VaultName: "openbao-vault",
						KeyName:   "unseal-key",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a transit seal without an address should be invalid", func() {
			input.Spec.AutoUnseal = &KubernetesOpenBaoAutoUnseal{
				Seal: &KubernetesOpenBaoAutoUnseal_Transit{
					Transit: &KubernetesOpenBaoTransitSeal{KeyName: "autounseal"},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("zero injector replicas should be invalid", func() {
			input.Spec.Injector = &KubernetesOpenBaoInjector{Replicas: int32Ptr(0)}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("injector replicas above 5 should be invalid", func() {
			input.Spec.Injector = &KubernetesOpenBaoInjector{Replicas: int32Ptr(6)}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a lowercase injector failure policy should be invalid", func() {
			input.Spec.Injector = &KubernetesOpenBaoInjector{FailurePolicy: strPtr("ignore")}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a ServiceMonitor without metrics enabled should be invalid", func() {
			input.Spec.Metrics = &KubernetesOpenBaoMetrics{ServiceMonitorEnabled: true}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a snapshot agent without an s3 host should be invalid", func() {
			input.Spec.SnapshotAgent = &KubernetesOpenBaoSnapshotAgent{
				S3Bucket:                "openbao-snapshots",
				S3CredentialsSecretName: "openbao-snapshot-s3-creds",
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a snapshot agent without a bucket should be invalid", func() {
			input.Spec.SnapshotAgent = &KubernetesOpenBaoSnapshotAgent{
				S3Host:                  literal("s3.us-east-1.amazonaws.com"),
				S3CredentialsSecretName: "openbao-snapshot-s3-creds",
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a snapshot agent without a credentials secret name should be invalid", func() {
			input.Spec.SnapshotAgent = &KubernetesOpenBaoSnapshotAgent{
				S3Host:   literal("s3.us-east-1.amazonaws.com"),
				S3Bucket: "openbao-snapshots",
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("zero snapshot expire days should be invalid", func() {
			input.Spec.SnapshotAgent = &KubernetesOpenBaoSnapshotAgent{
				S3Host:                  literal("s3.us-east-1.amazonaws.com"),
				S3Bucket:                "openbao-snapshots",
				S3CredentialsSecretName: "openbao-snapshot-s3-creds",
				S3ExpireDays:            int32Ptr(0),
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})
	})
})
