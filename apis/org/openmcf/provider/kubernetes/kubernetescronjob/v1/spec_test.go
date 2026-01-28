package kubernetescronjobv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestKubernetesCronJob(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesCronJob Suite")
}

var _ = ginkgo.Describe("KubernetesCronJob Custom Validation Tests", func() {
	var input *KubernetesCronJob

	ginkgo.BeforeEach(func() {
		input = &KubernetesCronJob{
			ApiVersion: "kubernetes.openmcf.org/v1",
			Kind:       "KubernetesCronJob",
			Metadata: &shared.CloudResourceMetadata{
				Name: "my-cron-job",
			},
			Spec: &KubernetesCronJobSpec{
				TargetCluster: &kubernetes.KubernetesClusterSelector{
					ClusterKind: cloudresourcekind.CloudResourceKind_GcpGkeCluster,
					ClusterName: "test-cluster",
				},
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
						Value: "test-namespace",
					},
				},
				CreateNamespace: true,
				Image: &kubernetes.ContainerImage{
					Repo: "busybox",
					Tag:  "latest",
				},
				Resources: &kubernetes.ContainerResources{
					Limits: &kubernetes.CpuMemory{
						Cpu:    "1000m",
						Memory: "1Gi",
					},
					Requests: &kubernetes.CpuMemory{
						Cpu:    "50m",
						Memory: "100Mi",
					},
				},
				Env: &KubernetesCronJobContainerAppEnv{
					Variables: map[string]*foreignkeyv1.StringValueOrRef{
						"ENV_VAR": {
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
								Value: "example",
							},
						},
					},
					Secrets: map[string]*kubernetes.KubernetesSensitiveValue{
						"SECRET_NAME": {
							SensitiveValue: &kubernetes.KubernetesSensitiveValue_Value{
								Value: "secret_value",
							},
						},
					},
				},
				Schedule:          "0 0 * * *",
				ConcurrencyPolicy: proto.String("Forbid"),
				RestartPolicy:     proto.String("Never"),
			},
		}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("cron_job_kubernetes with create_namespace=true", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("cron_job_kubernetes with create_namespace=false", func() {
			ginkgo.It("should not return a validation error", func() {
				input.Spec.CreateNamespace = false
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("Environment secrets validation", func() {
		ginkgo.Context("When secrets have direct string values", func() {
			ginkgo.It("should pass validation", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Secrets: map[string]*kubernetes.KubernetesSensitiveValue{
						"DATABASE_PASSWORD": {
							SensitiveValue: &kubernetes.KubernetesSensitiveValue_Value{
								Value: "my-password",
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("When secrets have Kubernetes Secret references", func() {
			ginkgo.It("should pass validation with valid secret ref", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Secrets: map[string]*kubernetes.KubernetesSensitiveValue{
						"DATABASE_PASSWORD": {
							SensitiveValue: &kubernetes.KubernetesSensitiveValue_SecretRef{
								SecretRef: &kubernetes.KubernetesSecretKeyRef{
									Name: "my-app-secrets",
									Key:  "db-password",
								},
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("When secrets have mixed types", func() {
			ginkgo.It("should pass validation with both string values and secret refs", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Secrets: map[string]*kubernetes.KubernetesSensitiveValue{
						"DEBUG_TOKEN": {
							SensitiveValue: &kubernetes.KubernetesSensitiveValue_Value{
								Value: "debug-only-token",
							},
						},
						"DATABASE_PASSWORD": {
							SensitiveValue: &kubernetes.KubernetesSensitiveValue_SecretRef{
								SecretRef: &kubernetes.KubernetesSecretKeyRef{
									Name: "postgres-credentials",
									Key:  "password",
								},
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("When secret ref is missing required fields", func() {
			ginkgo.It("should fail validation when name is missing", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Secrets: map[string]*kubernetes.KubernetesSensitiveValue{
						"DATABASE_PASSWORD": {
							SensitiveValue: &kubernetes.KubernetesSensitiveValue_SecretRef{
								SecretRef: &kubernetes.KubernetesSecretKeyRef{
									Name: "",
									Key:  "password",
								},
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})

			ginkgo.It("should fail validation when key is missing", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Secrets: map[string]*kubernetes.KubernetesSensitiveValue{
						"DATABASE_PASSWORD": {
							SensitiveValue: &kubernetes.KubernetesSensitiveValue_SecretRef{
								SecretRef: &kubernetes.KubernetesSecretKeyRef{
									Name: "my-secret",
									Key:  "",
								},
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("Environment variables validation", func() {
		ginkgo.Context("When variables have direct string values", func() {
			ginkgo.It("should pass validation", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Variables: map[string]*foreignkeyv1.StringValueOrRef{
						"BACKUP_RETENTION_DAYS": {
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
								Value: "30",
							},
						},
						"LOG_LEVEL": {
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
								Value: "info",
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("When variables have valueFrom references", func() {
			ginkgo.It("should pass validation with valid valueFrom ref", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Variables: map[string]*foreignkeyv1.StringValueOrRef{
						"DATABASE_HOST": {
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
								ValueFrom: &foreignkeyv1.ValueFromRef{
									Name: "my-postgres",
								},
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("When variables have mixed types", func() {
			ginkgo.It("should pass validation with both direct values and valueFrom refs", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Variables: map[string]*foreignkeyv1.StringValueOrRef{
						"BACKUP_RETENTION_DAYS": {
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
								Value: "30",
							},
						},
						"DATABASE_HOST": {
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
								ValueFrom: &foreignkeyv1.ValueFromRef{
									Name: "my-postgres",
								},
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("When valueFrom ref is missing required name", func() {
			ginkgo.It("should fail validation", func() {
				input.Spec.Env = &KubernetesCronJobContainerAppEnv{
					Variables: map[string]*foreignkeyv1.StringValueOrRef{
						"DATABASE_HOST": {
							LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
								ValueFrom: &foreignkeyv1.ValueFromRef{
									Name: "",
								},
							},
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
