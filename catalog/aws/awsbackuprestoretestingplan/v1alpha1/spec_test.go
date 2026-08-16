package awsbackuprestoretestingplanv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBackupRestoreTestingPlanSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBackupRestoreTestingPlanSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalTestingPlan is the smallest valid instance: weekly random
// tests over every vault, no selections yet.
func minimalTestingPlan() *AwsBackupRestoreTestingPlanSpec {
	return &AwsBackupRestoreTestingPlanSpec{
		Region:             "us-west-2",
		PlanName:           "weekly_restore_tests",
		ScheduleExpression: "cron(0 5 ? * MON *)",
		RecoveryPointSelection: &AwsBackupRestoreTestingPlanRecoveryPointSelection{
			Algorithm:          "RANDOM_WITHIN_WINDOW",
			IncludeVaults:      []string{"*"},
			RecoveryPointTypes: []string{"SNAPSHOT"},
		},
	}
}

func selectionByArns() *AwsBackupRestoreTestingPlanSelection {
	return &AwsBackupRestoreTestingPlanSelection{
		Name:                  "ebs_volumes",
		ProtectedResourceType: "EBS",
		IamRoleArn:            svr("arn:aws:iam::123456789012:role/restore-testing"),
		ProtectedResourceArns: []string{"*"},
	}
}

var _ = ginkgo.Describe("AwsBackupRestoreTestingPlanSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal testing plan", func() {
			gomega.Expect(protovalidate.Validate(minimalTestingPlan())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a selection by explicit ARNs", func() {
			spec := minimalTestingPlan()
			spec.Selections = []*AwsBackupRestoreTestingPlanSelection{selectionByArns()}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a selection by tag conditions with overrides", func() {
			spec := minimalTestingPlan()
			spec.Selections = []*AwsBackupRestoreTestingPlanSelection{{
				Name:                  "tagged_volumes",
				ProtectedResourceType: "EBS",
				IamRoleArn:            svr("arn:aws:iam::123456789012:role/restore-testing"),
				ProtectedResourceConditions: &AwsBackupRestoreTestingPlanSelectionConditions{
					StringEquals: []*AwsBackupRestoreTestingPlanSelectionConditionPair{{
						Key:   "aws:ResourceTag/backup",
						Value: "true",
					}},
				},
				RestoreMetadataOverrides: map[string]string{"availabilityzone": "us-west-2a"},
				ValidationWindowHours:    12,
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts vault ARNs and window bounds on the recovery point selection", func() {
			spec := minimalTestingPlan()
			spec.RecoveryPointSelection.IncludeVaults = []string{"arn:aws:backup:us-west-2:123456789012:backup-vault:app-vault"}
			spec.RecoveryPointSelection.ExcludeVaults = []string{"arn:aws:backup:us-west-2:123456789012:backup-vault:scratch"}
			spec.RecoveryPointSelection.SelectionWindowDays = 365
			spec.StartWindowHours = 168
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalTestingPlan()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a hyphenated plan name", func() {
			spec := minimalTestingPlan()
			spec.PlanName = "weekly-restore-tests"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing schedule", func() {
			spec := minimalTestingPlan()
			spec.ScheduleExpression = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing recovery point selection", func() {
			spec := minimalTestingPlan()
			spec.RecoveryPointSelection = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown selection algorithm", func() {
			spec := minimalTestingPlan()
			spec.RecoveryPointSelection.Algorithm = "OLDEST_WITHIN_WINDOW"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an include_vaults entry that is neither an ARN nor the wildcard", func() {
			spec := minimalTestingPlan()
			spec.RecoveryPointSelection.IncludeVaults = []string{"app-vault"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a selection window beyond 365 days", func() {
			spec := minimalTestingPlan()
			spec.RecoveryPointSelection.SelectionWindowDays = 366
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a start window beyond 168 hours", func() {
			spec := minimalTestingPlan()
			spec.StartWindowHours = 169
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a selection with BOTH arns and conditions", func() {
			spec := minimalTestingPlan()
			sel := selectionByArns()
			sel.ProtectedResourceConditions = &AwsBackupRestoreTestingPlanSelectionConditions{
				StringEquals: []*AwsBackupRestoreTestingPlanSelectionConditionPair{{Key: "aws:ResourceTag/x", Value: "y"}},
			}
			spec.Selections = []*AwsBackupRestoreTestingPlanSelection{sel}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a selection with NEITHER arns nor conditions", func() {
			spec := minimalTestingPlan()
			sel := selectionByArns()
			sel.ProtectedResourceArns = nil
			spec.Selections = []*AwsBackupRestoreTestingPlanSelection{sel}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty conditions block", func() {
			spec := minimalTestingPlan()
			sel := selectionByArns()
			sel.ProtectedResourceArns = nil
			sel.ProtectedResourceConditions = &AwsBackupRestoreTestingPlanSelectionConditions{}
			spec.Selections = []*AwsBackupRestoreTestingPlanSelection{sel}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate selection names", func() {
			spec := minimalTestingPlan()
			spec.Selections = []*AwsBackupRestoreTestingPlanSelection{selectionByArns(), selectionByArns()}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
