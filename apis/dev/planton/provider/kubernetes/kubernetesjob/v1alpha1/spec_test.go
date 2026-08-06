package kubernetesjobv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	kubernetes "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func validContainer() *kubernetes.WorkloadContainer {
	return &kubernetes.WorkloadContainer{
		Image: &kubernetes.ContainerImage{
			Repo: "ghcr.io/acme/migrator",
			Tag:  "v3.0.1",
		},
	}
}

func validSpec() *KubernetesJobSpec {
	return &KubernetesJobSpec{
		Namespace: literal("batch"),
		Container: &KubernetesJobContainer{
			App: validContainer(),
		},
	}
}

func TestKubernetesJobSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesJobSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesJobSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec with namespace and app container", func() {
			err := protovalidate.Validate(validSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a full-featured indexed batch job with policies, sidecars, and scheduling", func() {
			app := validContainer()
			app.Name = "worker"
			spec := &KubernetesJobSpec{
				Namespace:       literal("batch"),
				CreateNamespace: true,
				Container: &KubernetesJobContainer{
					App: app,
					Sidecars: []*kubernetes.WorkloadContainer{
						{
							Name: "log-shipper",
							Image: &kubernetes.ContainerImage{
								Repo: "ghcr.io/acme/fluent-bit",
								Tag:  "2.2.0",
							},
						},
					},
				},
				Pod: &kubernetes.WorkloadPod{
					ServiceAccount: literal("batch-identity"),
					Scheduling: &kubernetes.WorkloadScheduling{
						NodeSelector: map[string]string{"workload-tier": "batch"},
					},
				},
				Parallelism:             proto.Int32(4),
				Completions:             proto.Int32(8),
				CompletionMode:          proto.String("Indexed"),
				BackoffLimit:            proto.Int32(3),
				BackoffLimitPerIndex:    proto.Uint32(2),
				MaxFailedIndexes:        proto.Uint32(1),
				ActiveDeadlineSeconds:   proto.Int64(3600),
				TtlSecondsAfterFinished: proto.Int32(86400),
				RestartPolicy:           proto.String("Never"),
				PodFailurePolicy: &KubernetesJobPodFailurePolicy{
					Rules: []*KubernetesJobPodFailurePolicyRule{
						{
							Action: "FailJob",
							OnExitCodes: &KubernetesJobPodFailurePolicyOnExitCodes{
								ContainerName: "worker",
								Operator:      "In",
								Values:        []int32{42},
							},
						},
						{
							Action: "Ignore",
							OnPodConditions: []*KubernetesJobPodFailurePolicyOnPodCondition{
								{Type: "DisruptionTarget", Status: proto.String("True")},
							},
						},
					},
				},
				SuccessPolicy: &KubernetesJobSuccessPolicy{
					Rules: []*KubernetesJobSuccessPolicyRule{
						{
							SucceededIndexes: "0-2,4",
							SucceededCount:   proto.Int32(3),
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the OnFailure restart policy", func() {
			spec := validSpec()
			spec.RestartPolicy = proto.String("OnFailure")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the NonIndexed completion mode", func() {
			spec := validSpec()
			spec.CompletionMode = proto.String("NonIndexed")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a pod failure policy rule with the NotIn exit-code operator", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{
				Rules: []*KubernetesJobPodFailurePolicyRule{
					{
						Action: "Count",
						OnExitCodes: &KubernetesJobPodFailurePolicyOnExitCodes{
							Operator: "NotIn",
							Values:   []int32{143},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a success policy rule with only a succeeded count", func() {
			spec := validSpec()
			spec.SuccessPolicy = &KubernetesJobSuccessPolicy{
				Rules: []*KubernetesJobSuccessPolicyRule{
					{SucceededCount: proto.Int32(1)},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a success policy rule with only a single succeeded index", func() {
			spec := validSpec()
			spec.SuccessPolicy = &KubernetesJobSuccessPolicy{
				Rules: []*KubernetesJobSuccessPolicyRule{
					{SucceededIndexes: "0"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a suspended job", func() {
			spec := validSpec()
			spec.Suspend = true
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a spec without a namespace", func() {
			spec := validSpec()
			spec.Namespace = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a spec without a container", func() {
			spec := validSpec()
			spec.Container = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a container group without an app container", func() {
			spec := validSpec()
			spec.Container = &KubernetesJobContainer{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an app image without a repo", func() {
			spec := validSpec()
			spec.Container.App.Image = &kubernetes.ContainerImage{Tag: "v3.0.1"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unnamed sidecar", func() {
			spec := validSpec()
			spec.Container.Sidecars = []*kubernetes.WorkloadContainer{
				{
					Image: &kubernetes.ContainerImage{Repo: "busybox", Tag: "1.36"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown completion mode", func() {
			spec := validSpec()
			spec.CompletionMode = proto.String("Partitioned")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects the Always restart policy", func() {
			spec := validSpec()
			spec.RestartPolicy = proto.String("Always")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a negative parallelism", func() {
			spec := validSpec()
			spec.Parallelism = proto.Int32(-1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a negative backoff limit", func() {
			spec := validSpec()
			spec.BackoffLimit = proto.Int32(-1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an active deadline below 1 second", func() {
			spec := validSpec()
			spec.ActiveDeadlineSeconds = proto.Int64(0)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a negative ttl_seconds_after_finished", func() {
			spec := validSpec()
			spec.TtlSecondsAfterFinished = proto.Int32(-1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a pod failure policy with no rules", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a pod failure policy rule with an unknown action", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{
				Rules: []*KubernetesJobPodFailurePolicyRule{
					{
						Action: "Retry",
						OnExitCodes: &KubernetesJobPodFailurePolicyOnExitCodes{
							Operator: "In",
							Values:   []int32{42},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a pod failure policy rule with no trigger", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{
				Rules: []*KubernetesJobPodFailurePolicyRule{
					{Action: "FailJob"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a pod failure policy rule with both triggers set", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{
				Rules: []*KubernetesJobPodFailurePolicyRule{
					{
						Action: "FailJob",
						OnExitCodes: &KubernetesJobPodFailurePolicyOnExitCodes{
							Operator: "In",
							Values:   []int32{42},
						},
						OnPodConditions: []*KubernetesJobPodFailurePolicyOnPodCondition{
							{Type: "DisruptionTarget"},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an on_exit_codes trigger with an unknown operator", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{
				Rules: []*KubernetesJobPodFailurePolicyRule{
					{
						Action: "FailJob",
						OnExitCodes: &KubernetesJobPodFailurePolicyOnExitCodes{
							Operator: "Equals",
							Values:   []int32{42},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an on_exit_codes trigger with no values", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{
				Rules: []*KubernetesJobPodFailurePolicyRule{
					{
						Action: "FailJob",
						OnExitCodes: &KubernetesJobPodFailurePolicyOnExitCodes{
							Operator: "In",
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an on_pod_conditions pattern with an unknown status", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{
				Rules: []*KubernetesJobPodFailurePolicyRule{
					{
						Action: "Ignore",
						OnPodConditions: []*KubernetesJobPodFailurePolicyOnPodCondition{
							{Type: "DisruptionTarget", Status: proto.String("Maybe")},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an on_pod_conditions pattern without a type", func() {
			spec := validSpec()
			spec.PodFailurePolicy = &KubernetesJobPodFailurePolicy{
				Rules: []*KubernetesJobPodFailurePolicyRule{
					{
						Action: "Ignore",
						OnPodConditions: []*KubernetesJobPodFailurePolicyOnPodCondition{
							{Status: proto.String("True")},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a success policy with no rules", func() {
			spec := validSpec()
			spec.SuccessPolicy = &KubernetesJobSuccessPolicy{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a success policy rule with neither criterion", func() {
			spec := validSpec()
			spec.SuccessPolicy = &KubernetesJobSuccessPolicy{
				Rules: []*KubernetesJobSuccessPolicyRule{
					{},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects malformed succeeded_indexes", func() {
			spec := validSpec()
			spec.SuccessPolicy = &KubernetesJobSuccessPolicy{
				Rules: []*KubernetesJobSuccessPolicyRule{
					{SucceededIndexes: "0-2,"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a succeeded count below 1", func() {
			spec := validSpec()
			spec.SuccessPolicy = &KubernetesJobSuccessPolicy{
				Rules: []*KubernetesJobSuccessPolicyRule{
					{SucceededCount: proto.Int32(0)},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
