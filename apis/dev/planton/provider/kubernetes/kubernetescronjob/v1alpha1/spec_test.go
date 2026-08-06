package kubernetescronjobv1alpha1

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
			Repo: "ghcr.io/acme/backup",
			Tag:  "v1.2.0",
		},
	}
}

func validJobTemplate() *KubernetesCronJobJobTemplate {
	return &KubernetesCronJobJobTemplate{
		Container: &KubernetesCronJobContainer{
			App: validContainer(),
		},
	}
}

func validSpec() *KubernetesCronJobSpec {
	return &KubernetesCronJobSpec{
		Namespace:   literal("ops"),
		Schedule:    "0 3 * * *",
		JobTemplate: validJobTemplate(),
	}
}

func TestKubernetesCronJobSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesCronJobSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesCronJobSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec with namespace, schedule, and job template", func() {
			err := protovalidate.Validate(validSpec())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a full-featured scheduled backup with policies, sidecars, and scheduling", func() {
			app := validContainer()
			app.Name = "backup"
			spec := &KubernetesCronJobSpec{
				Namespace:                  literal("ops"),
				CreateNamespace:            true,
				Schedule:                   "0 9 * * 1-5",
				TimeZone:                   proto.String("America/New_York"),
				StartingDeadlineSeconds:    proto.Int64(300),
				ConcurrencyPolicy:          proto.String("Forbid"),
				SuccessfulJobsHistoryLimit: proto.Int32(3),
				FailedJobsHistoryLimit:     proto.Int32(1),
				JobTemplate: &KubernetesCronJobJobTemplate{
					Container: &KubernetesCronJobContainer{
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
						ServiceAccount: literal("backup-identity"),
						Scheduling: &kubernetes.WorkloadScheduling{
							NodeSelector: map[string]string{"workload-tier": "batch"},
							TopologySpreadConstraints: []*kubernetes.WorkloadTopologySpreadConstraint{
								{MaxSkew: 1, TopologyKey: "topology.kubernetes.io/zone", WhenUnsatisfiable: "DoNotSchedule"},
							},
						},
					},
					Parallelism:             proto.Int32(2),
					Completions:             proto.Int32(4),
					CompletionMode:          proto.String("Indexed"),
					BackoffLimit:            proto.Int32(3),
					ActiveDeadlineSeconds:   proto.Int64(1800),
					TtlSecondsAfterFinished: proto.Int32(3600),
					RestartPolicy:           proto.String("Never"),
					PodFailurePolicy: &KubernetesCronJobPodFailurePolicy{
						Rules: []*KubernetesCronJobPodFailurePolicyRule{
							{
								Action: "FailJob",
								OnExitCodes: &KubernetesCronJobPodFailurePolicyOnExitCodes{
									Operator: "In",
									Values:   []int32{42},
								},
							},
						},
					},
					SuccessPolicy: &KubernetesCronJobSuccessPolicy{
						Rules: []*KubernetesCronJobSuccessPolicyRule{
							{SucceededIndexes: "0"},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a step-interval schedule", func() {
			spec := validSpec()
			spec.Schedule = "*/15 * * * *"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the Allow and Replace concurrency policies", func() {
			spec := validSpec()
			spec.ConcurrencyPolicy = proto.String("Allow")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())

			spec.ConcurrencyPolicy = proto.String("Replace")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts zero history limits", func() {
			spec := validSpec()
			spec.SuccessfulJobsHistoryLimit = proto.Int32(0)
			spec.FailedJobsHistoryLimit = proto.Int32(0)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a suspended cron job", func() {
			spec := validSpec()
			spec.Suspend = true
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts the OnFailure restart policy in the job template", func() {
			spec := validSpec()
			spec.JobTemplate.RestartPolicy = proto.String("OnFailure")
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

		ginkgo.It("rejects a spec without a schedule", func() {
			spec := validSpec()
			spec.Schedule = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a schedule that is not a cron expression", func() {
			spec := validSpec()
			spec.Schedule = "not-a-cron"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a 4-field cron schedule", func() {
			spec := validSpec()
			spec.Schedule = "0 3 * *"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a spec without a job template", func() {
			spec := validSpec()
			spec.JobTemplate = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a job template without a container", func() {
			spec := validSpec()
			spec.JobTemplate = &KubernetesCronJobJobTemplate{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a container group without an app container", func() {
			spec := validSpec()
			spec.JobTemplate.Container = &KubernetesCronJobContainer{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an app image without a tag", func() {
			spec := validSpec()
			spec.JobTemplate.Container.App.Image = &kubernetes.ContainerImage{Repo: "ghcr.io/acme/backup"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unnamed sidecar in the job template", func() {
			spec := validSpec()
			spec.JobTemplate.Container.Sidecars = []*kubernetes.WorkloadContainer{
				{
					Image: &kubernetes.ContainerImage{Repo: "busybox", Tag: "1.36"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown concurrency policy", func() {
			spec := validSpec()
			spec.ConcurrencyPolicy = proto.String("Queue")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a negative successful jobs history limit", func() {
			spec := validSpec()
			spec.SuccessfulJobsHistoryLimit = proto.Int32(-1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a negative failed jobs history limit", func() {
			spec := validSpec()
			spec.FailedJobsHistoryLimit = proto.Int32(-1)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a starting deadline below 1 second", func() {
			spec := validSpec()
			spec.StartingDeadlineSeconds = proto.Int64(0)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown completion mode in the job template", func() {
			spec := validSpec()
			spec.JobTemplate.CompletionMode = proto.String("Partitioned")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects the Always restart policy in the job template", func() {
			spec := validSpec()
			spec.JobTemplate.RestartPolicy = proto.String("Always")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a job template active deadline below 1 second", func() {
			spec := validSpec()
			spec.JobTemplate.ActiveDeadlineSeconds = proto.Int64(0)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a pod failure policy rule with both triggers set", func() {
			spec := validSpec()
			spec.JobTemplate.PodFailurePolicy = &KubernetesCronJobPodFailurePolicy{
				Rules: []*KubernetesCronJobPodFailurePolicyRule{
					{
						Action: "FailJob",
						OnExitCodes: &KubernetesCronJobPodFailurePolicyOnExitCodes{
							Operator: "In",
							Values:   []int32{42},
						},
						OnPodConditions: []*KubernetesCronJobPodFailurePolicyOnPodCondition{
							{Type: "DisruptionTarget"},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a pod failure policy rule with an unknown action", func() {
			spec := validSpec()
			spec.JobTemplate.PodFailurePolicy = &KubernetesCronJobPodFailurePolicy{
				Rules: []*KubernetesCronJobPodFailurePolicyRule{
					{
						Action: "Retry",
						OnExitCodes: &KubernetesCronJobPodFailurePolicyOnExitCodes{
							Operator: "In",
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
			spec.JobTemplate.PodFailurePolicy = &KubernetesCronJobPodFailurePolicy{
				Rules: []*KubernetesCronJobPodFailurePolicyRule{
					{
						Action: "FailJob",
						OnExitCodes: &KubernetesCronJobPodFailurePolicyOnExitCodes{
							Operator: "In",
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a success policy rule with neither criterion", func() {
			spec := validSpec()
			spec.JobTemplate.SuccessPolicy = &KubernetesCronJobSuccessPolicy{
				Rules: []*KubernetesCronJobSuccessPolicyRule{
					{},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects malformed succeeded_indexes in the job template", func() {
			spec := validSpec()
			spec.JobTemplate.SuccessPolicy = &KubernetesCronJobSuccessPolicy{
				Rules: []*KubernetesCronJobSuccessPolicyRule{
					{SucceededIndexes: "a-b"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a topology spread constraint with an unknown whenUnsatisfiable", func() {
			spec := validSpec()
			spec.JobTemplate.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					TopologySpreadConstraints: []*kubernetes.WorkloadTopologySpreadConstraint{
						{MaxSkew: 1, TopologyKey: "topology.kubernetes.io/zone", WhenUnsatisfiable: "Ignore"},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a topology spread constraint with max_skew below 1", func() {
			spec := validSpec()
			spec.JobTemplate.Pod = &kubernetes.WorkloadPod{
				Scheduling: &kubernetes.WorkloadScheduling{
					TopologySpreadConstraints: []*kubernetes.WorkloadTopologySpreadConstraint{
						{MaxSkew: 0, TopologyKey: "topology.kubernetes.io/zone", WhenUnsatisfiable: "DoNotSchedule"},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
