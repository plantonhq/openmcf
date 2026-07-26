package kubernetesgrafanav1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesGrafana(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesGrafana Suite")
}

func int32Ptr(i int32) *int32    { return &i }
func stringPtr(s string) *string { return &s }

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

func testDatabase() *KubernetesGrafanaDatabase {
	return &KubernetesGrafanaDatabase{
		Engine:         KubernetesGrafanaDatabaseEngine_postgres,
		Host:           literal("grafana-db-rw.data.svc.cluster.local:5432"),
		Name:           "grafana",
		User:           "grafana",
		PasswordSecret: &KubernetesGrafanaSecretKeyRef{Name: "grafana-db-app", Key: "password"},
	}
}

var _ = ginkgo.Describe("KubernetesGrafana Validation Tests", func() {
	var input *KubernetesGrafana

	ginkgo.BeforeEach(func() {
		input = &KubernetesGrafana{
			ApiVersion: "kubernetes.planton.dev/v1",
			Kind:       "KubernetesGrafana",
			Metadata: &shared.CloudResourceMetadata{
				Name: "dashboards",
			},
			Spec: &KubernetesGrafanaSpec{
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

		ginkgo.It("a single replica with persistent storage should be valid", func() {
			input.Spec.Storage = &KubernetesGrafanaStorage{
				Size:         stringPtr("20Gi"),
				StorageClass: literal("gp3"),
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("multiple replicas WITH an external database should be valid", func() {
			input.Spec.Replicas = int32Ptr(3)
			input.Spec.Database = testDatabase()
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a database host referencing a KubernetesPostgres should be valid", func() {
			db := testDatabase()
			db.Host = valueFrom(cloudresourcekind.CloudResourceKind_KubernetesPostgres, "grafana-db", "status.outputs.kube_endpoint")
			input.Spec.Database = db
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a literal-url datasource should be valid", func() {
			input.Spec.Datasources = []*KubernetesGrafanaDatasource{{
				Name: "Prometheus",
				Url:  literal("http://monitoring-prometheus.observability.svc.cluster.local:9090"),
			}}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a datasource referencing a KubernetesKubePrometheusStack should be valid", func() {
			input.Spec.Datasources = []*KubernetesGrafanaDatasource{{
				Name:      "Prometheus",
				Url:       valueFrom(cloudresourcekind.CloudResourceKind_KubernetesKubePrometheusStack, "monitoring", "status.outputs.prometheus_endpoint"),
				IsDefault: true,
			}}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a datasource with basic auth should be valid", func() {
			input.Spec.Datasources = []*KubernetesGrafanaDatasource{{
				Name: "External Mimir",
				Url:  literal("https://mimir.example.com/prometheus"),
				BasicAuth: &KubernetesGrafanaDatasourceBasicAuth{
					Username:       "tenant-1",
					PasswordSecret: &KubernetesGrafanaSecretKeyRef{Name: "mimir-credentials", Key: "password"},
				},
			}}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a community dashboard pinned to a revision should be valid", func() {
			input.Spec.Datasources = []*KubernetesGrafanaDatasource{{
				Name: "Prometheus",
				Url:  literal("http://monitoring-prometheus.observability.svc.cluster.local:9090"),
			}}
			input.Spec.CommunityDashboards = []*KubernetesGrafanaCommunityDashboard{{
				GnetId:     1860,
				Revision:   37,
				Datasource: "Prometheus",
			}}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("an existing admin secret should be valid", func() {
			input.Spec.AdminSecret = &KubernetesGrafanaAdminSecret{Name: "grafana-admin"}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("smtp with a credentials secret should be valid", func() {
			input.Spec.Smtp = &KubernetesGrafanaSmtp{
				Host:                  "smtp.example.com:587",
				FromAddress:           "grafana@example.com",
				CredentialsSecretName: "smtp-credentials",
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("full surface should be valid", func() {
			input.Spec.CreateNamespace = true
			input.Spec.ChartVersion = stringPtr("12.8.0")
			input.Spec.Replicas = int32Ptr(2)
			input.Spec.Database = testDatabase()
			input.Spec.Datasources = []*KubernetesGrafanaDatasource{{
				Name:      "Prometheus",
				Url:       literal("http://monitoring-prometheus.observability.svc.cluster.local:9090"),
				IsDefault: true,
				Uid:       "prom-main",
				JsonData:  "httpMethod: POST\n",
			}}
			input.Spec.Plugins = []string{"grafana-clock-panel", "elasticsearch"}
			input.Spec.Server = &KubernetesGrafanaServer{RootUrl: "https://grafana.example.com"}
			input.Spec.Auth = &KubernetesGrafanaAuth{AnonymousEnabled: true}
			input.Spec.ServiceMonitorEnabled = true
			input.Spec.Image = &KubernetesGrafanaImage{Repository: "mirror.example.com/grafana/grafana"}
			input.Spec.Scheduling = &KubernetesGrafanaScheduling{
				NodeSelector: map[string]string{"workload": "monitoring"},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.It("a namespace-less spec should fail", func() {
			input.Spec.Namespace = nil
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("multiple replicas WITHOUT a database should fail", func() {
			input.Spec.Replicas = int32Ptr(3)
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("multiple replicas with storage should fail even with a database (single-writer volume)", func() {
			input.Spec.Replicas = int32Ptr(2)
			input.Spec.Database = testDatabase()
			input.Spec.Storage = &KubernetesGrafanaStorage{}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("zero replicas should fail", func() {
			input.Spec.Replicas = int32Ptr(0)
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a database without an engine should fail", func() {
			db := testDatabase()
			db.Engine = KubernetesGrafanaDatabaseEngine_kubernetes_grafana_database_engine_unspecified
			input.Spec.Database = db
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a database without a password secret should fail", func() {
			db := testDatabase()
			db.PasswordSecret = nil
			input.Spec.Database = db
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a database password secret without a key should fail", func() {
			db := testDatabase()
			db.PasswordSecret = &KubernetesGrafanaSecretKeyRef{Name: "grafana-db-app"}
			input.Spec.Database = db
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a datasource without a name should fail", func() {
			input.Spec.Datasources = []*KubernetesGrafanaDatasource{{
				Url: literal("http://prometheus:9090"),
			}}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a datasource without a url should fail", func() {
			input.Spec.Datasources = []*KubernetesGrafanaDatasource{{Name: "Prometheus"}}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("datasource basic auth without a password secret should fail", func() {
			input.Spec.Datasources = []*KubernetesGrafanaDatasource{{
				Name:      "External",
				Url:       literal("https://mimir.example.com"),
				BasicAuth: &KubernetesGrafanaDatasourceBasicAuth{Username: "u"},
			}}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a community dashboard without a gnet id should fail", func() {
			input.Spec.CommunityDashboards = []*KubernetesGrafanaCommunityDashboard{{
				Datasource: "Prometheus",
			}}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a community dashboard without a datasource should fail", func() {
			input.Spec.CommunityDashboards = []*KubernetesGrafanaCommunityDashboard{{GnetId: 1860}}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("smtp without a host should fail", func() {
			input.Spec.Smtp = &KubernetesGrafanaSmtp{FromAddress: "grafana@example.com"}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a malformed storage size should fail", func() {
			input.Spec.Storage = &KubernetesGrafanaStorage{Size: stringPtr("twenty gigs")}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an admin secret without a name should fail", func() {
			input.Spec.AdminSecret = &KubernetesGrafanaAdminSecret{}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})
	})
})
