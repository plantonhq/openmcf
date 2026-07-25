package kubernetessignozv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesSignoz(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesSignoz Suite")
}

func int32Ptr(i int32) *int32 { return &i }

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

// testExternalClickHouse returns an external-database arm shaped like a
// KubernetesClickHouse composition with literal values.
func testExternalClickHouse() *KubernetesSignozExternalClickHouse {
	return &KubernetesSignozExternalClickHouse{
		Host:     literal("clickhouse-analytics.data.svc.cluster.local"),
		Username: "signoz",
		PasswordSecret: &KubernetesSignozExternalClickHousePassword{
			SecretName: literal("analytics-clickhouse-auth"),
			SecretKey:  "signoz",
		},
	}
}

var _ = ginkgo.Describe("KubernetesSignoz Validation Tests", func() {
	var input *KubernetesSignoz

	ginkgo.BeforeEach(func() {
		input = &KubernetesSignoz{
			ApiVersion: "kubernetes.planton.dev/v1",
			Kind:       "KubernetesSignoz",
			Metadata: &shared.CloudResourceMetadata{
				Name: "observe",
			},
			Spec: &KubernetesSignozSpec{
				Namespace: literal("observability"),
			},
		}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.It("minimal spec should not return a validation error (every optional block omitted)", func() {
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("namespace as a reference should be valid", func() {
			input.Spec.Namespace = valueFrom(cloudresourcekind.CloudResourceKind_KubernetesNamespace, "observability", "spec.name")
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a fully declared managed-clickhouse arm should be valid", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					Shards:            int32Ptr(2),
					Replicas:          int32Ptr(2),
					DiskSize:          stringPtr("100Gi"),
					AllowedNetworkIps: []string{"10.0.0.0/8"},
					Zookeeper: &KubernetesSignozZookeeper{
						Replicas: int32Ptr(3),
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("managed clickhouse with keyless IRSA cold storage should be valid", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					ColdStorage: &KubernetesSignozColdStorage{
						Backend: &KubernetesSignozColdStorage_S3{
							S3: &KubernetesSignozColdStorageS3{
								Endpoint:    "https://telemetry-cold.s3-us-east-1.amazonaws.com/data/",
								IrsaRoleArn: "arn:aws:iam::123456789012:role/signoz-cold-storage",
							},
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("managed clickhouse with declared-key S3 cold storage should be valid", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					ColdStorage: &KubernetesSignozColdStorage{
						Backend: &KubernetesSignozColdStorage_S3{
							S3: &KubernetesSignozColdStorageS3{
								Endpoint:  "https://telemetry-cold.s3-us-east-1.amazonaws.com/data/",
								AccessKey: "AKIAEXAMPLE",
								SecretKey: "secretmaterial",
							},
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("managed clickhouse with GCS cold storage should be valid", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					ColdStorage: &KubernetesSignozColdStorage{
						Backend: &KubernetesSignozColdStorage_Gcs{
							Gcs: &KubernetesSignozColdStorageGcs{
								Endpoint:  "https://storage.googleapis.com/telemetry-cold/data/",
								AccessKey: "GOOGEXAMPLE",
								SecretKey: "secretmaterial",
							},
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("an external clickhouse with literal values should be valid", func() {
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{
				ExternalClickhouse: testExternalClickHouse(),
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("an external clickhouse composed entirely from KubernetesClickHouse references should be valid", func() {
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{
				ExternalClickhouse: &KubernetesSignozExternalClickHouse{
					Host:        valueFrom(cloudresourcekind.CloudResourceKind_KubernetesClickHouse, "analytics", "status.outputs.service_name"),
					ClusterName: valueFrom(cloudresourcekind.CloudResourceKind_KubernetesClickHouse, "analytics", "status.outputs.cluster_name"),
					Username:    "signoz",
					PasswordSecret: &KubernetesSignozExternalClickHousePassword{
						SecretName: valueFrom(cloudresourcekind.CloudResourceKind_KubernetesClickHouse, "analytics", "status.outputs.auth_secret_name"),
						SecretKey:  "signoz",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("an external clickhouse over verified TLS should be valid", func() {
			external := testExternalClickHouse()
			external.Secure = true
			external.Verify = true
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{ExternalClickhouse: external}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a fully declared server block should be valid", func() {
			input.Spec.Server = &KubernetesSignozServer{
				DiskSize:    stringPtr("2Gi"),
				ExternalUrl: "https://signoz.example.com",
				Smtp: &KubernetesSignozSmtp{
					Address:  "smtp.example.com:587",
					From:     "signoz@example.com",
					Username: "smtp-user",
					PasswordSecret: &KubernetesSignozSecretKeyRef{
						Name: "smtp-auth",
						Key:  "password",
					},
				},
				Env: map[string]string{
					"signoz_prometheus_active__query__tracker_enabled": "true",
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("unauthenticated SMTP (an internal relay) should be valid", func() {
			input.Spec.Server = &KubernetesSignozServer{
				Smtp: &KubernetesSignozSmtp{
					Address: "mail-relay.internal:25",
					From:    "signoz@example.com",
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("collector autoscaling with sane bounds should be valid", func() {
			input.Spec.OtelCollector = &KubernetesSignozOtelCollector{
				Autoscaling: &KubernetesSignozOtelCollectorAutoscaling{
					Enabled:                        true,
					MinReplicas:                    int32Ptr(2),
					MaxReplicas:                    int32Ptr(10),
					TargetCpuUtilizationPercent:    int32Ptr(60),
					TargetMemoryUtilizationPercent: int32Ptr(70),
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.It("an even zookeeper replica count should fail (quorum teaching)", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					Zookeeper: &KubernetesSignozZookeeper{Replicas: int32Ptr(2)},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("zero clickhouse shards should fail", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{Shards: int32Ptr(0)},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a malformed clickhouse disk size should fail", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{DiskSize: stringPtr("twenty-gigs")},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("S3 cold storage with BOTH a role and declared keys should fail", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					ColdStorage: &KubernetesSignozColdStorage{
						Backend: &KubernetesSignozColdStorage_S3{
							S3: &KubernetesSignozColdStorageS3{
								Endpoint:    "https://telemetry-cold.s3.amazonaws.com/data/",
								IrsaRoleArn: "arn:aws:iam::123456789012:role/signoz",
								AccessKey:   "AKIAEXAMPLE",
								SecretKey:   "secretmaterial",
							},
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("S3 cold storage with NO auth posture should fail", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					ColdStorage: &KubernetesSignozColdStorage{
						Backend: &KubernetesSignozColdStorage_S3{
							S3: &KubernetesSignozColdStorageS3{
								Endpoint: "https://telemetry-cold.s3.amazonaws.com/data/",
							},
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("S3 cold storage with only HALF a key pair should fail", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					ColdStorage: &KubernetesSignozColdStorage{
						Backend: &KubernetesSignozColdStorage_S3{
							S3: &KubernetesSignozColdStorageS3{
								Endpoint:  "https://telemetry-cold.s3.amazonaws.com/data/",
								AccessKey: "AKIAEXAMPLE",
							},
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("S3 cold storage without an endpoint should fail", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					ColdStorage: &KubernetesSignozColdStorage{
						Backend: &KubernetesSignozColdStorage_S3{
							S3: &KubernetesSignozColdStorageS3{
								IrsaRoleArn: "arn:aws:iam::123456789012:role/signoz",
							},
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("GCS cold storage without keys should fail", func() {
			input.Spec.Database = &KubernetesSignozSpec_ManagedClickhouse{
				ManagedClickhouse: &KubernetesSignozManagedClickHouse{
					ColdStorage: &KubernetesSignozColdStorage{
						Backend: &KubernetesSignozColdStorage_Gcs{
							Gcs: &KubernetesSignozColdStorageGcs{
								Endpoint: "https://storage.googleapis.com/telemetry-cold/data/",
							},
						},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an external clickhouse without a host should fail", func() {
			external := testExternalClickHouse()
			external.Host = nil
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{ExternalClickhouse: external}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an external clickhouse without a username should fail", func() {
			external := testExternalClickHouse()
			external.Username = ""
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{ExternalClickhouse: external}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an external clickhouse without a password secret should fail", func() {
			external := testExternalClickHouse()
			external.PasswordSecret = nil
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{ExternalClickhouse: external}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an external clickhouse password secret without a key should fail", func() {
			external := testExternalClickHouse()
			external.PasswordSecret.SecretKey = ""
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{ExternalClickhouse: external}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("certificate verification without TLS should fail", func() {
			external := testExternalClickHouse()
			external.Verify = true
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{ExternalClickhouse: external}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an out-of-range external clickhouse port should fail", func() {
			external := testExternalClickHouse()
			external.TcpPort = int32Ptr(70000)
			input.Spec.Database = &KubernetesSignozSpec_ExternalClickhouse{ExternalClickhouse: external}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an SMTP address without a port should fail", func() {
			input.Spec.Server = &KubernetesSignozServer{
				Smtp: &KubernetesSignozSmtp{
					Address: "smtp.example.com",
					From:    "signoz@example.com",
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("SMTP without a From address should fail", func() {
			input.Spec.Server = &KubernetesSignozServer{
				Smtp: &KubernetesSignozSmtp{Address: "smtp.example.com:587"},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an SMTP password without a username should fail", func() {
			input.Spec.Server = &KubernetesSignozServer{
				Smtp: &KubernetesSignozSmtp{
					Address: "smtp.example.com:587",
					From:    "signoz@example.com",
					PasswordSecret: &KubernetesSignozSecretKeyRef{
						Name: "smtp-auth",
						Key:  "password",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a malformed server disk size should fail", func() {
			input.Spec.Server = &KubernetesSignozServer{DiskSize: stringPtr("big")}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("zero collector replicas should fail", func() {
			input.Spec.OtelCollector = &KubernetesSignozOtelCollector{Replicas: int32Ptr(0)}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("collector autoscaling with max below min should fail", func() {
			input.Spec.OtelCollector = &KubernetesSignozOtelCollector{
				Autoscaling: &KubernetesSignozOtelCollectorAutoscaling{
					Enabled:     true,
					MinReplicas: int32Ptr(5),
					MaxReplicas: int32Ptr(2),
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})
	})
})

func stringPtr(s string) *string { return &s }
