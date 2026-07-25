package kuberneteskubeprometheusstackv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesKubePrometheusStack(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesKubePrometheusStack Suite")
}

func int32Ptr(i int32) *int32    { return &i }
func stringPtr(s string) *string { return &s }
func boolPtr(b bool) *bool       { return &b }

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

var _ = ginkgo.Describe("KubernetesKubePrometheusStack Validation Tests", func() {
	var input *KubernetesKubePrometheusStack

	ginkgo.BeforeEach(func() {
		input = &KubernetesKubePrometheusStack{
			ApiVersion: "kubernetes.planton.dev/v1",
			Kind:       "KubernetesKubePrometheusStack",
			Metadata: &shared.CloudResourceMetadata{
				Name: "monitoring",
			},
			Spec: &KubernetesKubePrometheusStackSpec{
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

		ginkgo.It("a tuned prometheus block should be valid", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				Replicas:      int32Ptr(2),
				Retention:     stringPtr("30d"),
				RetentionSize: "45GiB",
				DiskSize:      stringPtr("100Gi"),
				StorageClass:  literal("gp3"),
				ExternalLabels: map[string]string{
					"cluster": "prod-us-east-1",
				},
				ScrapeInterval:     "30s",
				EvaluationInterval: "1m",
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("prometheus ephemeral with the MIDDLEWARE-DEFAULTED disk size should be valid", func() {
			// The platform's defaulting middleware stamps disk_size "50Gi" onto
			// every manifest — an ephemeral manifest must stay expressible.
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				Ephemeral: true,
				DiskSize:  stringPtr("50Gi"),
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("remote write with basic auth should be valid", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{
					Url:  "https://prometheus-us-central1.grafana.net/api/prom/push",
					Name: "grafana-cloud",
					Auth: &KubernetesKubePrometheusStackRemoteWrite_BasicAuth{
						BasicAuth: &KubernetesKubePrometheusStackRemoteWriteBasicAuth{
							Username:       "123456",
							PasswordSecret: &KubernetesKubePrometheusStackSecretKeyRef{Name: "grafana-cloud-token", Key: "token"},
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("keyless SigV4 remote write (region only) should be valid", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{
					Url: "https://aps-workspaces.us-east-1.amazonaws.com/workspaces/ws-abc/api/v1/remote_write",
					Auth: &KubernetesKubePrometheusStackRemoteWrite_Sigv4{
						Sigv4: &KubernetesKubePrometheusStackRemoteWriteSigv4{Region: "us-east-1"},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("SigV4 with BOTH static keys should be valid", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{
					Url: "https://aps-workspaces.us-east-1.amazonaws.com/workspaces/ws-abc/api/v1/remote_write",
					Auth: &KubernetesKubePrometheusStackRemoteWrite_Sigv4{
						Sigv4: &KubernetesKubePrometheusStackRemoteWriteSigv4{
							Region:          "us-east-1",
							AccessKeySecret: &KubernetesKubePrometheusStackSecretKeyRef{Name: "amp-keys", Key: "access-key-id"},
							SecretKeySecret: &KubernetesKubePrometheusStackSecretKeyRef{Name: "amp-keys", Key: "secret-access-key"},
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("azure managed-identity remote write should be valid", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{
					Url: "https://my-workspace.eastus-1.metrics.ingest.monitor.azure.com/dataCollectionRules/dcr/streams/Microsoft-PrometheusMetrics/api/v1/write",
					Auth: &KubernetesKubePrometheusStackRemoteWrite_AzureAd{
						AzureAd: &KubernetesKubePrometheusStackRemoteWriteAzureAd{
							ManagedIdentityClientId: "11111111-2222-3333-4444-555555555555",
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("alertmanager ephemeral with the middleware-defaulted disk size should be valid", func() {
			input.Spec.Alertmanager = &KubernetesKubePrometheusStackAlertmanager{
				Ephemeral: true,
				DiskSize:  stringPtr("2Gi"),
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("bundled grafana with an existing admin secret should be valid", func() {
			input.Spec.Grafana = &KubernetesKubePrometheusStackGrafana{
				AdminSecret: &KubernetesKubePrometheusStackGrafanaAdminSecret{
					Name: "grafana-admin",
				},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("cert-manager-issued admission webhook certificates should be valid", func() {
			input.Spec.Operator = &KubernetesKubePrometheusStackOperator{
				AdmissionWebhooks: &KubernetesKubePrometheusStackAdmissionWebhooks{CertManager: true},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("release-fenced monitor discovery should be valid", func() {
			discovery := KubernetesKubePrometheusStackMonitorDiscovery_release_managed_only
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{Discovery: &discovery}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("the managed-cloud scraper posture should be valid", func() {
			input.Spec.ControlPlaneScrapers = &KubernetesKubePrometheusStackControlPlaneScrapers{
				KubeControllerManager: boolPtr(false),
				KubeEtcd:              boolPtr(false),
				KubeScheduler:         boolPtr(false),
				KubeProxy:             boolPtr(false),
			}
			input.Spec.DefaultRules = &KubernetesKubePrometheusStackDefaultRules{
				DisabledGroups: []string{"etcd", "kubeControllerManager", "kubeSchedulerAlerting", "kubeSchedulerRecording", "kubeProxy"},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("full surface should be valid", func() {
			input.Spec.CreateNamespace = true
			input.Spec.ChartVersion = stringPtr("87.19.1")
			input.Spec.CrdUpgradeJob = true
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				Replicas:                  int32Ptr(2),
				Retention:                 stringPtr("15d"),
				DiskSize:                  stringPtr("200Gi"),
				EnableRemoteWriteReceiver: true,
				Scheduling: &KubernetesKubePrometheusStackScheduling{
					NodeSelector:      map[string]string{"workload": "monitoring"},
					PriorityClassName: "high",
				},
			}
			input.Spec.Alertmanager = &KubernetesKubePrometheusStackAlertmanager{
				Replicas:   int32Ptr(3),
				ConfigYaml: "route:\n  receiver: slack\nreceivers:\n  - name: slack\n",
			}
			input.Spec.Grafana = &KubernetesKubePrometheusStackGrafana{
				Storage: &KubernetesKubePrometheusStackGrafanaStorage{Size: stringPtr("20Gi")},
			}
			input.Spec.Exporters = &KubernetesKubePrometheusStackExporters{
				KubeStateMetricsEnabled: boolPtr(true),
				NodeExporterEnabled:     boolPtr(true),
			}
			input.Spec.ImageRegistry = "mirror.example.com"
			input.Spec.ImagePullSecrets = []string{"mirror-pull"}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.It("a namespace-less spec should fail", func() {
			input.Spec.Namespace = nil
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("zero prometheus replicas should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{Replicas: int32Ptr(0)}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a malformed retention should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{Retention: stringPtr("10days")}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a malformed retention_size should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{RetentionSize: "45G"}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a malformed scrape_interval should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{ScrapeInterval: "half a minute"}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a malformed evaluation_interval should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{EvaluationInterval: "1minute"}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("prometheus ephemeral with a NON-default disk size should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				Ephemeral: true,
				DiskSize:  stringPtr("100Gi"),
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("prometheus ephemeral with a storage class should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				Ephemeral:    true,
				StorageClass: literal("gp3"),
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("alertmanager ephemeral with a non-default disk size should fail", func() {
			input.Spec.Alertmanager = &KubernetesKubePrometheusStackAlertmanager{
				Ephemeral: true,
				DiskSize:  stringPtr("10Gi"),
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("alertmanager ephemeral with a storage class should fail", func() {
			input.Spec.Alertmanager = &KubernetesKubePrometheusStackAlertmanager{
				Ephemeral:    true,
				StorageClass: literal("gp3"),
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a remote write without a url should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{Name: "no-url"}},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("basic auth without a password secret should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{
					Url: "https://example.com/api/v1/write",
					Auth: &KubernetesKubePrometheusStackRemoteWrite_BasicAuth{
						BasicAuth: &KubernetesKubePrometheusStackRemoteWriteBasicAuth{Username: "u"},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("SigV4 with only one static key should fail (keys are a pair)", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{
					Url: "https://aps-workspaces.us-east-1.amazonaws.com/workspaces/ws-abc/api/v1/remote_write",
					Auth: &KubernetesKubePrometheusStackRemoteWrite_Sigv4{
						Sigv4: &KubernetesKubePrometheusStackRemoteWriteSigv4{
							Region:          "us-east-1",
							AccessKeySecret: &KubernetesKubePrometheusStackSecretKeyRef{Name: "amp-keys", Key: "access-key-id"},
						},
					},
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("SigV4 without a region should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{
					Url:  "https://example.com/api/v1/remote_write",
					Auth: &KubernetesKubePrometheusStackRemoteWrite_Sigv4{Sigv4: &KubernetesKubePrometheusStackRemoteWriteSigv4{}},
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("azure remote write without a managed-identity client id should fail", func() {
			input.Spec.Prometheus = &KubernetesKubePrometheusStackPrometheus{
				RemoteWrite: []*KubernetesKubePrometheusStackRemoteWrite{{
					Url:  "https://example.com/api/v1/write",
					Auth: &KubernetesKubePrometheusStackRemoteWrite_AzureAd{AzureAd: &KubernetesKubePrometheusStackRemoteWriteAzureAd{}},
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("disabled admission webhooks combined with cert-manager should fail", func() {
			input.Spec.Operator = &KubernetesKubePrometheusStackOperator{
				AdmissionWebhooks: &KubernetesKubePrometheusStackAdmissionWebhooks{
					Disabled:    true,
					CertManager: true,
				},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a grafana admin secret without a name should fail", func() {
			input.Spec.Grafana = &KubernetesKubePrometheusStackGrafana{
				AdminSecret: &KubernetesKubePrometheusStackGrafanaAdminSecret{},
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})
	})
})
