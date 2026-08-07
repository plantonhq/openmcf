package awsbatchjobqueuev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBatchJobQueueSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBatchJobQueueSpec Validation Suite")
}

func stringPtr(s string) *string {
	return &s
}

func svRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func ceOrder(order int32, arn string) *AwsBatchJobQueueComputeEnvironmentOrder {
	return &AwsBatchJobQueueComputeEnvironmentOrder{
		Order:              order,
		ComputeEnvironment: svRef(arn),
	}
}

func minimalSpec() *AwsBatchJobQueueSpec {
	return &AwsBatchJobQueueSpec{
		Region:   "us-west-2",
		Priority: 1,
		ComputeEnvironmentOrder: []*AwsBatchJobQueueComputeEnvironmentOrder{
			ceOrder(1, "arn:aws:batch:us-west-2:123456789012:compute-environment/primary"),
		},
	}
}

var _ = ginkgo.Describe("AwsBatchJobQueueSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalSpec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the Spot-first overflow pattern (two environments)", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalSpec()
				spec.ComputeEnvironmentOrder = []*AwsBatchJobQueueComputeEnvironmentOrder{
					ceOrder(1, "arn:aws:batch:us-west-2:123456789012:compute-environment/spot"),
					ceOrder(2, "arn:aws:batch:us-west-2:123456789012:compute-environment/on-demand"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a scheduling policy reference", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalSpec()
				spec.SchedulingPolicy = svRef("arn:aws:batch:us-west-2:123456789012:scheduling-policy/fair")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with job state time limit actions", func() {
			ginkgo.It("should accept every AWS stuck-job reason matcher", func() {
				spec := minimalSpec()
				spec.JobStateTimeLimitActions = []*AwsBatchJobStateTimeLimitAction{
					{Action: "CANCEL", MaxTimeSeconds: 3600, Reason: "CAPACITY:INSUFFICIENT_INSTANCE_CAPACITY", State: "RUNNABLE"},
					{Action: "CANCEL", MaxTimeSeconds: 3600, Reason: "MISCONFIGURATION:COMPUTE_ENVIRONMENT_MAX_RESOURCE", State: "RUNNABLE"},
					{Action: "CANCEL", MaxTimeSeconds: 3600, Reason: "MISCONFIGURATION:JOB_RESOURCE_REQUIREMENT", State: "RUNNABLE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with state set to DISABLED", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalSpec()
				spec.State = stringPtr("DISABLED")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with priority zero", func() {
			ginkgo.It("should not return a validation error (0 is a valid priority)", func() {
				spec := minimalSpec()
				spec.Priority = 0
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with no region", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalSpec()
				spec.Region = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid state", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalSpec()
				spec.State = stringPtr("PAUSED")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a negative priority", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalSpec()
				spec.Priority = -1
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with no compute environments", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalSpec()
				spec.ComputeEnvironmentOrder = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with more than three compute environments", func() {
			ginkgo.It("should return a validation error (AWS caps a queue at three)", func() {
				spec := minimalSpec()
				spec.ComputeEnvironmentOrder = []*AwsBatchJobQueueComputeEnvironmentOrder{
					ceOrder(1, "arn:aws:batch:us-west-2:123456789012:compute-environment/a"),
					ceOrder(2, "arn:aws:batch:us-west-2:123456789012:compute-environment/b"),
					ceOrder(3, "arn:aws:batch:us-west-2:123456789012:compute-environment/c"),
					ceOrder(4, "arn:aws:batch:us-west-2:123456789012:compute-environment/d"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a compute environment entry missing its reference", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalSpec()
				spec.ComputeEnvironmentOrder = []*AwsBatchJobQueueComputeEnvironmentOrder{
					{Order: 1},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a compute environment order below 1", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalSpec()
				spec.ComputeEnvironmentOrder[0].Order = 0
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid time limit action", func() {
			ginkgo.It("should reject an unknown action", func() {
				spec := minimalSpec()
				spec.JobStateTimeLimitActions = []*AwsBatchJobStateTimeLimitAction{
					{Action: "PAUSE", MaxTimeSeconds: 3600, Reason: "r", State: "RUNNABLE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject TERMINATE (SageMaker-service-environment-only, not modeled)", func() {
				spec := minimalSpec()
				spec.JobStateTimeLimitActions = []*AwsBatchJobStateTimeLimitAction{
					{Action: "TERMINATE", MaxTimeSeconds: 3600, Reason: "r", State: "RUNNABLE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject a free-text reason (AWS accepts only its stuck-job cause matchers)", func() {
				spec := minimalSpec()
				spec.JobStateTimeLimitActions = []*AwsBatchJobStateTimeLimitAction{
					{Action: "CANCEL", MaxTimeSeconds: 3600, Reason: "stuck in RUNNABLE for 1h", State: "RUNNABLE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject max_time_seconds below 600", func() {
				spec := minimalSpec()
				spec.JobStateTimeLimitActions = []*AwsBatchJobStateTimeLimitAction{
					{Action: "CANCEL", MaxTimeSeconds: 100, Reason: "r", State: "RUNNABLE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject max_time_seconds above 86400", func() {
				spec := minimalSpec()
				spec.JobStateTimeLimitActions = []*AwsBatchJobStateTimeLimitAction{
					{Action: "CANCEL", MaxTimeSeconds: 90000, Reason: "r", State: "RUNNABLE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject a state other than RUNNABLE", func() {
				spec := minimalSpec()
				spec.JobStateTimeLimitActions = []*AwsBatchJobStateTimeLimitAction{
					{Action: "CANCEL", MaxTimeSeconds: 3600, Reason: "r", State: "RUNNING"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject a missing reason", func() {
				spec := minimalSpec()
				spec.JobStateTimeLimitActions = []*AwsBatchJobStateTimeLimitAction{
					{Action: "CANCEL", MaxTimeSeconds: 3600, State: "RUNNABLE"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
