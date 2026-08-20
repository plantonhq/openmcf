package awseventbridgeschedulerv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsEventBridgeSchedulerSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEventBridgeSchedulerSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func minimalSchedule() *AwsEventBridgeSchedulerSpec {
	return &AwsEventBridgeSchedulerSpec{
		Region:             "us-east-1",
		ScheduleExpression: "rate(5 minutes)",
		FlexibleTimeWindow: &AwsEventBridgeScheduleFlexibleTimeWindow{Mode: "OFF"},
		Target: &AwsEventBridgeScheduleTarget{
			Arn:     svr("arn:aws:sqs:us-east-1:123456789012:jobs"),
			RoleArn: svr("arn:aws:iam::123456789012:role/scheduler-exec"),
		},
	}
}

var _ = ginkgo.Describe("AwsEventBridgeSchedulerSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal schedule", func() {
			gomega.Expect(protovalidate.Validate(minimalSchedule())).To(gomega.BeNil())
		})

		ginkgo.It("accepts an owned group", func() {
			spec := minimalSchedule()
			spec.Group = &AwsEventBridgeScheduleGroup{Name: "batch-jobs"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts joining an existing group by name", func() {
			spec := minimalSchedule()
			spec.GroupName = "shared-jobs"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a flexible window with a size", func() {
			spec := minimalSchedule()
			spec.FlexibleTimeWindow = &AwsEventBridgeScheduleFlexibleTimeWindow{
				Mode:                   "FLEXIBLE",
				MaximumWindowInMinutes: proto.Int32(15),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a zero-retry policy (boundary)", func() {
			spec := minimalSchedule()
			spec.Target.RetryPolicy = &AwsEventBridgeScheduleRetryPolicy{
				MaximumRetryAttempts: proto.Int32(0),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts full ECS parameters", func() {
			spec := minimalSchedule()
			spec.Target.Arn = svr("arn:aws:ecs:us-east-1:123456789012:cluster/apps")
			spec.Target.EcsParameters = &AwsEventBridgeScheduleEcsParameters{
				TaskDefinitionArn: svr("arn:aws:ecs:us-east-1:123456789012:task-definition/nightly:3"),
				LaunchType:        "FARGATE",
				NetworkConfiguration: &AwsEventBridgeScheduleNetworkConfiguration{
					Subnets: []*foreignkeyv1.StringValueOrRef{svr("subnet-0123456789abcdef0")},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects both an owned group and a group_name", func() {
			spec := minimalSchedule()
			spec.Group = &AwsEventBridgeScheduleGroup{Name: "batch-jobs"}
			spec.GroupName = "shared-jobs"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing flexible_time_window", func() {
			spec := minimalSchedule()
			spec.FlexibleTimeWindow = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects mode OFF with a window size", func() {
			spec := minimalSchedule()
			spec.FlexibleTimeWindow.MaximumWindowInMinutes = proto.Int32(15)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects mode FLEXIBLE without a window size", func() {
			spec := minimalSchedule()
			spec.FlexibleTimeWindow.Mode = "FLEXIBLE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target without a role", func() {
			spec := minimalSchedule()
			spec.Target.RoleArn = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects two service parameter blocks", func() {
			spec := minimalSchedule()
			spec.Target.SqsParameters = &AwsEventBridgeScheduleSqsParameters{MessageGroupId: "g1"}
			spec.Target.KinesisParameters = &AwsEventBridgeScheduleKinesisParameters{PartitionKey: "pk"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range retry attempt count", func() {
			spec := minimalSchedule()
			spec.Target.RetryPolicy = &AwsEventBridgeScheduleRetryPolicy{
				MaximumRetryAttempts: proto.Int32(186),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown state", func() {
			spec := minimalSchedule()
			spec.State = "PAUSED"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects ECS parameters without a task definition", func() {
			spec := minimalSchedule()
			spec.Target.EcsParameters = &AwsEventBridgeScheduleEcsParameters{}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects ECS network configuration without subnets", func() {
			spec := minimalSchedule()
			spec.Target.EcsParameters = &AwsEventBridgeScheduleEcsParameters{
				TaskDefinitionArn:    svr("arn:aws:ecs:us-east-1:123456789012:task-definition/nightly:3"),
				NetworkConfiguration: &AwsEventBridgeScheduleNetworkConfiguration{},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid propagate_tags value", func() {
			spec := minimalSchedule()
			spec.Target.EcsParameters = &AwsEventBridgeScheduleEcsParameters{
				TaskDefinitionArn: svr("arn:aws:ecs:us-east-1:123456789012:task-definition/nightly:3"),
				PropagateTags:     "SERVICE",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
