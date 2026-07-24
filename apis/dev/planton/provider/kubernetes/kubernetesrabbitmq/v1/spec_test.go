package kubernetesrabbitmqv1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	kubernetes "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesRabbitMqSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesRabbitMqSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func int32Ptr(v int32) *int32 { return &v }

func int64Ptr(v int64) *int64 { return &v }

func strPtr(v string) *string { return &v }

var _ = ginkgo.Describe("KubernetesRabbitMqSpec validations", func() {
	var spec *KubernetesRabbitMqSpec

	ginkgo.BeforeEach(func() {
		spec = &KubernetesRabbitMqSpec{
			Namespace: literal("messaging"),
		}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.It("accepts a minimal single-node spec", func() {
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a production 3-node spec with persistence and spreading", func() {
			spec.Replicas = int32Ptr(3)
			spec.DiskSize = strPtr("50Gi")
			spec.StorageClass = literal("gp3")
			spec.SpreadAcrossNodes = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an ephemeral dev spec without storage fields", func() {
			spec.Ephemeral = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts ephemeral with the middleware-defaulted disk_size", func() {
			// The platform's defaulting middleware stamps 10Gi onto every
			// manifest — ephemeral must stay expressible through it (the
			// module ignores the defaulted value).
			spec.Ephemeral = true
			spec.DiskSize = strPtr("10Gi")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts additional plugins with valid names", func() {
			spec.Configuration = &KubernetesRabbitMqConfiguration{
				AdditionalPlugins: []string{"rabbitmq_shovel", "rabbitmq_mqtt", "rabbitmq_stream"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts rabbitmq.conf additions", func() {
			spec.Configuration = &KubernetesRabbitMqConfiguration{
				AdditionalConfig: "vm_memory_high_watermark.relative = 0.8\nconsumer_timeout = 1800000",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts env_config without shell substitution", func() {
			spec.Configuration = &KubernetesRabbitMqConfiguration{
				EnvConfig: "RABBITMQ_LOGS=-",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts TLS with a certificate secret", func() {
			spec.Tls = &KubernetesRabbitMqTls{
				SecretName: literal("orders-mq-tls"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts mutual TLS with a CA secret and closed plain listeners", func() {
			spec.Tls = &KubernetesRabbitMqTls{
				SecretName:             literal("orders-mq-tls"),
				CaSecretName:           literal("orders-mq-ca"),
				DisableNonTlsListeners: true,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a LoadBalancer service with annotations", func() {
			spec.Service = &KubernetesRabbitMqService{
				Type: KubernetesRabbitMqService_load_balancer,
				Annotations: map[string]string{
					"service.beta.kubernetes.io/aws-load-balancer-type": "nlb",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a vault secret backend", func() {
			spec.SecretBackend = &KubernetesRabbitMqSecretBackend{
				Backend: &KubernetesRabbitMqSecretBackend_Vault{
					Vault: &KubernetesRabbitMqVaultBackend{
						Role:            "rabbitmq",
						DefaultUserPath: "secret/data/rabbitmq/config",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an external-secret backend", func() {
			spec.SecretBackend = &KubernetesRabbitMqSecretBackend{
				Backend: &KubernetesRabbitMqSecretBackend_ExternalSecretName{
					ExternalSecretName: "orders-mq-creds",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts placement and tuning knobs", func() {
			spec.NodeSelector = map[string]string{"workload": "messaging"}
			spec.TerminationGracePeriodSeconds = int64Ptr(300)
			spec.DelayStartSeconds = int32Ptr(0)
			spec.SkipPostDeploySteps = true
			spec.AutoEnableAllFeatureFlags = true
			spec.Resources = &kubernetes.ContainerResources{
				Requests: &kubernetes.CpuMemory{Cpu: "1000m", Memory: "2Gi"},
				Limits:   &kubernetes.CpuMemory{Cpu: "2000m", Memory: "2Gi"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.It("rejects a missing namespace", func() {
			spec.Namespace = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects zero replicas", func() {
			spec.Replicas = int32Ptr(0)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed disk size", func() {
			spec.DiskSize = strPtr("fifty-gigs")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects ephemeral combined with a non-default disk_size", func() {
			spec.Ephemeral = true
			spec.DiskSize = strPtr("50Gi")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("ephemeral"))
		})

		ginkgo.It("rejects ephemeral combined with storage_class", func() {
			spec.Ephemeral = true
			spec.StorageClass = literal("gp3")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects plugin names with invalid characters", func() {
			spec.Configuration = &KubernetesRabbitMqConfiguration{
				AdditionalPlugins: []string{"rabbitmq-shovel"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects more than 100 plugins", func() {
			plugins := make([]string, 101)
			for i := range plugins {
				plugins[i] = "p" + strings.Repeat("x", 3)
			}
			spec.Configuration = &KubernetesRabbitMqConfiguration{AdditionalPlugins: plugins}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects env_config with command substitution", func() {
			spec.Configuration = &KubernetesRabbitMqConfiguration{
				EnvConfig: "RABBITMQ_LOGS=$(rm -rf /)",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("shell"))
		})

		ginkgo.It("rejects env_config with backticks", func() {
			spec.Configuration = &KubernetesRabbitMqConfiguration{
				EnvConfig: "RABBITMQ_LOGS=`id`",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects TLS without a certificate secret", func() {
			spec.Tls = &KubernetesRabbitMqTls{
				DisableNonTlsListeners: true,
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a vault backend without a role", func() {
			spec.SecretBackend = &KubernetesRabbitMqSecretBackend{
				Backend: &KubernetesRabbitMqSecretBackend_Vault{
					Vault: &KubernetesRabbitMqVaultBackend{
						DefaultUserPath: "secret/data/rabbitmq/config",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a vault backend without a default user path", func() {
			spec.SecretBackend = &KubernetesRabbitMqSecretBackend{
				Backend: &KubernetesRabbitMqSecretBackend_Vault{
					Vault: &KubernetesRabbitMqVaultBackend{
						Role: "rabbitmq",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects negative termination grace period", func() {
			spec.TerminationGracePeriodSeconds = int64Ptr(-1)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects negative delay start seconds", func() {
			spec.DelayStartSeconds = int32Ptr(-5)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
