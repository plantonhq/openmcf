package kubernetesotelcollectorv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	kubernetes "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesOtelCollector(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesOtelCollector Suite")
}

func int32Ptr(i int32) *int32 { return &i }

func modePtr(m KubernetesOtelCollectorMode) *KubernetesOtelCollectorMode { return &m }

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

// otlpFanInConfig is a minimal, realistic collector pipeline: OTLP in,
// debug exporter out.
const otlpFanInConfig = `receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318
exporters:
  debug: {}
service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [debug]
`

var _ = ginkgo.Describe("KubernetesOtelCollector Validation Tests", func() {
	var input *KubernetesOtelCollector

	ginkgo.BeforeEach(func() {
		input = &KubernetesOtelCollector{
			ApiVersion: "kubernetes.planton.dev/v1",
			Kind:       "KubernetesOtelCollector",
			Metadata: &shared.CloudResourceMetadata{
				Name: "telemetry",
			},
			Spec: &KubernetesOtelCollectorSpec{
				Namespace:  literal("observability"),
				ConfigYaml: otlpFanInConfig,
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

		ginkgo.It("an explicit deployment mode with replicas should be valid", func() {
			input.Spec.Mode = modePtr(KubernetesOtelCollectorMode_deployment)
			input.Spec.Replicas = int32Ptr(3)
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("daemonset mode with the middleware-defaulted replicas should be valid", func() {
			// The platform's defaulting middleware stamps replicas 1 onto
			// every manifest — daemonset manifests must stay expressible
			// with it present.
			input.Spec.Mode = modePtr(KubernetesOtelCollectorMode_daemonset)
			input.Spec.Replicas = int32Ptr(1)
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("a daemonset log collector with hostPath volumes and tolerations should be valid", func() {
			input.Spec.Mode = modePtr(KubernetesOtelCollectorMode_daemonset)
			input.Spec.Volumes = []*kubernetes.VolumeMount{{
				Name:      "varlogpods",
				MountPath: "/var/log/pods",
				ReadOnly:  true,
				HostPath:  &kubernetes.HostPathVolumeSource{Path: "/var/log/pods"},
			}}
			input.Spec.Scheduling = &KubernetesOtelCollectorScheduling{
				Tolerations: []*kubernetes.WorkloadToleration{{
					Key:      "node-role.kubernetes.io/control-plane",
					Operator: "Exists",
					Effect:   "NoSchedule",
				}},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("statefulset mode with an autoscaler should be valid", func() {
			input.Spec.Mode = modePtr(KubernetesOtelCollectorMode_statefulset)
			input.Spec.Autoscaler = &KubernetesOtelCollectorAutoscaler{
				MinReplicas: int32Ptr(2),
				MaxReplicas: 10,
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("sidecar mode should be valid", func() {
			input.Spec.Mode = modePtr(KubernetesOtelCollectorMode_sidecar)
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("env and secret-fed env should be valid", func() {
			input.Spec.Env = map[string]string{"GOMEMLIMIT": "400MiB"}
			input.Spec.EnvFromSecrets = []string{"tempo-basic-auth"}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("additional ports should be valid", func() {
			input.Spec.AdditionalPorts = []*KubernetesOtelCollectorPort{
				{Name: "syslog", Port: 5514, Protocol: func() *string { s := "UDP"; return &s }()},
			}
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})

		ginkgo.It("full surface should be valid", func() {
			input.Spec.CreateNamespace = true
			input.Spec.Mode = modePtr(KubernetesOtelCollectorMode_deployment)
			input.Spec.Autoscaler = &KubernetesOtelCollectorAutoscaler{
				MinReplicas:             int32Ptr(2),
				MaxReplicas:             8,
				TargetCpuUtilization:    int32Ptr(70),
				TargetMemoryUtilization: int32Ptr(80),
			}
			input.Spec.Image = "mirror.example.com/otel/custom-collector:1.0.0"
			input.Spec.ServiceAccount = "otel-collector"
			gomega.Expect(protovalidate.Validate(input)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.It("a namespace-less spec should fail", func() {
			input.Spec.Namespace = nil
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("a spec without config_yaml should fail", func() {
			input.Spec.ConfigYaml = ""
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("daemonset mode with declared replicas should fail", func() {
			input.Spec.Mode = modePtr(KubernetesOtelCollectorMode_daemonset)
			input.Spec.Replicas = int32Ptr(3)
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("sidecar mode with an autoscaler should fail", func() {
			input.Spec.Mode = modePtr(KubernetesOtelCollectorMode_sidecar)
			input.Spec.Autoscaler = &KubernetesOtelCollectorAutoscaler{MaxReplicas: 4}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an autoscaler with declared replicas should fail", func() {
			input.Spec.Replicas = int32Ptr(3)
			input.Spec.Autoscaler = &KubernetesOtelCollectorAutoscaler{MaxReplicas: 4}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an autoscaler with max below min should fail", func() {
			input.Spec.Autoscaler = &KubernetesOtelCollectorAutoscaler{
				MinReplicas: int32Ptr(5),
				MaxReplicas: 2,
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an autoscaler without max_replicas should fail", func() {
			input.Spec.Autoscaler = &KubernetesOtelCollectorAutoscaler{}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an out-of-range CPU target should fail", func() {
			input.Spec.Autoscaler = &KubernetesOtelCollectorAutoscaler{
				MaxReplicas:          4,
				TargetCpuUtilization: int32Ptr(150),
			}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an uppercase port name should fail", func() {
			input.Spec.AdditionalPorts = []*KubernetesOtelCollectorPort{{Name: "Syslog", Port: 5514}}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an out-of-range port should fail", func() {
			input.Spec.AdditionalPorts = []*KubernetesOtelCollectorPort{{Name: "syslog", Port: 70000}}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})

		ginkgo.It("an unsupported port protocol should fail", func() {
			input.Spec.AdditionalPorts = []*KubernetesOtelCollectorPort{{
				Name:     "sctp-port",
				Port:     5000,
				Protocol: func() *string { s := "SCTP"; return &s }(),
			}}
			gomega.Expect(protovalidate.Validate(input)).NotTo(gomega.BeNil())
		})
	})
})
