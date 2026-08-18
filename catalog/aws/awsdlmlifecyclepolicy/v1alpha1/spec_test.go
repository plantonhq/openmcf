package awsdlmlifecyclepolicyv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsDlmLifecyclePolicySpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsDlmLifecyclePolicySpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func dailySchedule() *AwsDlmSchedule {
	return &AwsDlmSchedule{
		Name:       "daily",
		CreateRule: &AwsDlmCreateRule{IntervalHours: 24, Times: []string{"03:00"}},
		RetainRule: &AwsDlmRetainRule{Count: 7},
	}
}

func customPolicy() *AwsDlmLifecyclePolicySpec {
	return &AwsDlmLifecyclePolicySpec{
		Region:           "us-west-2",
		ExecutionRoleArn: literal("arn:aws:iam::111122223333:role/dlm"),
		CustomPolicy: &AwsDlmCustomPolicy{
			ResourceTypes: []string{"VOLUME"},
			TargetTags:    map[string]string{"backup": "true"},
			Schedules:     []*AwsDlmSchedule{dailySchedule()},
		},
	}
}

func defaultPolicy() *AwsDlmLifecyclePolicySpec {
	return &AwsDlmLifecyclePolicySpec{
		Region:           "us-west-2",
		ExecutionRoleArn: literal("arn:aws:iam::111122223333:role/dlm"),
		DefaultPolicy:    &AwsDlmDefaultPolicy{ResourceType: "VOLUME"},
	}
}

var _ = ginkgo.Describe("AwsDlmLifecyclePolicySpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal default policy", func() {
			gomega.Expect(protovalidate.Validate(defaultPolicy())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a default policy with volume exclusions", func() {
			spec := defaultPolicy()
			spec.DefaultPolicy.Exclusions = &AwsDlmDefaultPolicyExclusions{
				ExcludeBootVolumes: true,
				ExcludeVolumeTypes: []string{"standard"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a custom snapshot policy with one schedule", func() {
			gomega.Expect(protovalidate.Validate(customPolicy())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cron schedule with a full rule set", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules = []*AwsDlmSchedule{{
				Name:       "weekly",
				CreateRule: &AwsDlmCreateRule{CronExpression: "cron(0 3 ? * SUN *)"},
				RetainRule: &AwsDlmRetainRule{Interval: 30, IntervalUnit: "DAYS"},
				ArchiveRule: &AwsDlmArchiveRule{
					Interval:     90,
					IntervalUnit: "DAYS",
				},
				CrossRegionCopyRules: []*AwsDlmCrossRegionCopyRule{{
					TargetRegion: "us-east-1",
					Encrypted:    true,
					RetainRule:   &AwsDlmCopyRetainRule{Interval: 14, IntervalUnit: "DAYS"},
				}},
				FastRestoreRule: &AwsDlmFastRestoreRule{
					AvailabilityZones: []string{"us-west-2a"},
					Count:             2,
				},
				ShareRule: &AwsDlmShareRule{TargetAccounts: []string{"111122223333"}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an event-based policy", func() {
			spec := &AwsDlmLifecyclePolicySpec{
				Region:           "us-west-2",
				ExecutionRoleArn: literal("arn:aws:iam::111122223333:role/dlm"),
				CustomPolicy: &AwsDlmCustomPolicy{
					PolicyType: "EVENT_BASED_POLICY",
					EventSource: &AwsDlmEventSource{
						EventType:        "shareSnapshot",
						DescriptionRegex: "^.*$",
						SnapshotOwners:   []string{"444455556666"},
					},
					Action: &AwsDlmAction{
						Name: "copy-shared",
						CrossRegionCopies: []*AwsDlmActionCrossRegionCopy{{
							Target: "us-east-1",
						}},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an instance policy with parameters", func() {
			spec := customPolicy()
			spec.CustomPolicy.ResourceTypes = []string{"INSTANCE"}
			spec.CustomPolicy.Parameters = &AwsDlmParameters{ExcludeBootVolume: true}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects zero modes", func() {
			spec := &AwsDlmLifecyclePolicySpec{
				Region:           "us-west-2",
				ExecutionRoleArn: literal("arn:aws:iam::111122223333:role/dlm"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects both modes at once", func() {
			spec := defaultPolicy()
			spec.CustomPolicy = customPolicy().CustomPolicy
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects boot-volume exclusion on an instance default policy", func() {
			spec := defaultPolicy()
			spec.DefaultPolicy.ResourceType = "INSTANCE"
			spec.DefaultPolicy.Exclusions = &AwsDlmDefaultPolicyExclusions{ExcludeBootVolumes: true}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a create interval outside 1-7", func() {
			spec := defaultPolicy()
			spec.DefaultPolicy.CreateIntervalDays = 9
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a snapshot policy without target tags", func() {
			spec := customPolicy()
			spec.CustomPolicy.TargetTags = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a snapshot policy without schedules", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects five schedules", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules = []*AwsDlmSchedule{}
			for _, name := range []string{"a", "b", "c", "d", "e"} {
				schedule := dailySchedule()
				schedule.Name = name
				spec.CustomPolicy.Schedules = append(spec.CustomPolicy.Schedules, schedule)
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate schedule names", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules = []*AwsDlmSchedule{dailySchedule(), dailySchedule()}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an event-based policy with target tags", func() {
			spec := customPolicy()
			spec.CustomPolicy.PolicyType = "EVENT_BASED_POLICY"
			spec.CustomPolicy.EventSource = &AwsDlmEventSource{
				EventType:        "shareSnapshot",
				DescriptionRegex: "^.*$",
				SnapshotOwners:   []string{"444455556666"},
			}
			spec.CustomPolicy.Action = &AwsDlmAction{
				Name:              "copy",
				CrossRegionCopies: []*AwsDlmActionCrossRegionCopy{{Target: "us-east-1"}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an event source on a schedule policy", func() {
			spec := customPolicy()
			spec.CustomPolicy.EventSource = &AwsDlmEventSource{
				EventType:        "shareSnapshot",
				DescriptionRegex: "^.*$",
				SnapshotOwners:   []string{"444455556666"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects parameters on a volume policy", func() {
			spec := customPolicy()
			spec.CustomPolicy.Parameters = &AwsDlmParameters{ExcludeBootVolume: true}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a create rule with both interval and cron", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules[0].CreateRule = &AwsDlmCreateRule{
				IntervalHours:  24,
				CronExpression: "cron(0 3 * * ? *)",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-canonical create interval", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules[0].CreateRule = &AwsDlmCreateRule{IntervalHours: 5}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects times alongside a cron expression", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules[0].CreateRule = &AwsDlmCreateRule{
				CronExpression: "cron(0 3 * * ? *)",
				Times:          []string{"03:00"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a retain rule with both count and interval", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules[0].RetainRule = &AwsDlmRetainRule{
				Count:        7,
				Interval:     30,
				IntervalUnit: "DAYS",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a retain interval without a unit", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules[0].RetainRule = &AwsDlmRetainRule{Interval: 30}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a cmk on an unencrypted cross-region copy", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules[0].CrossRegionCopyRules = []*AwsDlmCrossRegionCopyRule{{
				TargetRegion: "us-east-1",
				CmkArn:       literal("arn:aws:kms:us-east-1:111122223333:key/abc"),
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a fast restore rule without zones", func() {
			spec := customPolicy()
			spec.CustomPolicy.Schedules[0].FastRestoreRule = &AwsDlmFastRestoreRule{Count: 2}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a description outside the DLM character set", func() {
			spec := customPolicy()
			spec.Description = "daily backups: volumes"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
