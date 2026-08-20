package awsbackupplanv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBackupPlanSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBackupPlanSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func int32Ptr(i int32) *int32 { return &i }

// minimalPlan is the smallest valid instance: one daily rule into a
// vault.
func minimalPlan() *AwsBackupPlanSpec {
	return &AwsBackupPlanSpec{
		Region: "us-west-2",
		Rules: []*AwsBackupPlanRule{{
			Name:            "daily",
			TargetVaultName: svr("app-vault"),
			Schedule:        "cron(0 5 ? * * *)",
		}},
	}
}

var _ = ginkgo.Describe("AwsBackupPlanSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal plan", func() {
			gomega.Expect(protovalidate.Validate(minimalPlan())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a lifecycle honoring the 90-day cold-storage minimum", func() {
			spec := minimalPlan()
			spec.Rules[0].Lifecycle = &AwsBackupPlanLifecycle{
				ColdStorageAfterDays: int32Ptr(30),
				DeleteAfterDays:      int32Ptr(120),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a delete-only lifecycle without cold storage", func() {
			spec := minimalPlan()
			spec.Rules[0].Lifecycle = &AwsBackupPlanLifecycle{DeleteAfterDays: int32Ptr(35)}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts copy actions, scan actions, and air-gapped targeting", func() {
			spec := minimalPlan()
			spec.Rules[0].TargetLogicallyAirGappedBackupVaultArn = svr("arn:aws:backup:us-west-2:123456789012:backup-vault:lag-vault")
			spec.Rules[0].CopyActions = []*AwsBackupPlanCopyAction{{
				DestinationVaultArn: svr("arn:aws:backup:us-east-1:123456789012:backup-vault:dr-vault"),
				Lifecycle:           &AwsBackupPlanLifecycle{DeleteAfterDays: int32Ptr(90)},
			}}
			spec.Rules[0].ScanActions = []*AwsBackupPlanScanAction{{
				MalwareScanner: "GUARDDUTY",
				ScanMode:       "FULL_SCAN",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts Windows VSS advanced settings and a plan-wide scan setting", func() {
			spec := minimalPlan()
			spec.AdvancedBackupSettings = []*AwsBackupPlanAdvancedBackupSetting{{
				ResourceType:  "EC2",
				BackupOptions: map[string]string{"WindowsVSS": "enabled"},
			}}
			spec.ScanSetting = &AwsBackupPlanScanSetting{
				MalwareScanner: "GUARDDUTY",
				ResourceTypes:  []string{"EBS", "EC2"},
				ScannerRoleArn: svr("arn:aws:iam::123456789012:role/backup-scanner"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a selection by ARNs and a selection by conditions", func() {
			spec := minimalPlan()
			spec.Selections = []*AwsBackupPlanSelection{
				{
					Name:       "by-arn",
					IamRoleArn: svr("arn:aws:iam::123456789012:role/backup-service"),
					Resources:  []string{"arn:aws:ec2:us-west-2:123456789012:volume/vol-0abc"},
				},
				{
					Name:          "by-tag",
					IamRoleArn:    svr("arn:aws:iam::123456789012:role/backup-service"),
					SelectionTags: []*AwsBackupPlanSelectionTag{{Type: "STRINGEQUALS", Key: "backup", Value: "true"}},
					Conditions: &AwsBackupPlanSelectionConditions{
						StringLike: []*AwsBackupPlanSelectionConditionPair{{Key: "aws:ResourceTag/environment", Value: "prod*"}},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalPlan()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a plan without rules", func() {
			spec := minimalPlan()
			spec.Rules = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate rule names", func() {
			spec := minimalPlan()
			spec.Rules = append(spec.Rules, &AwsBackupPlanRule{
				Name:            "daily",
				TargetVaultName: svr("other-vault"),
			})
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule without a target vault", func() {
			spec := minimalPlan()
			spec.Rules[0].TargetVaultName = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a rule name with illegal characters", func() {
			spec := minimalPlan()
			spec.Rules[0].Name = "daily rule"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a start window below 60 minutes", func() {
			spec := minimalPlan()
			spec.Rules[0].StartWindowMinutes = 30
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a lifecycle violating the 90-day cold-storage minimum", func() {
			spec := minimalPlan()
			spec.Rules[0].Lifecycle = &AwsBackupPlanLifecycle{
				ColdStorageAfterDays: int32Ptr(30),
				DeleteAfterDays:      int32Ptr(119),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown malware scanner", func() {
			spec := minimalPlan()
			spec.Rules[0].ScanActions = []*AwsBackupPlanScanAction{{
				MalwareScanner: "CLAMAV",
				ScanMode:       "FULL_SCAN",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an advanced setting for a non-EC2 resource type", func() {
			spec := minimalPlan()
			spec.AdvancedBackupSettings = []*AwsBackupPlanAdvancedBackupSetting{{
				ResourceType:  "RDS",
				BackupOptions: map[string]string{"WindowsVSS": "enabled"},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a scan setting without a scanner role", func() {
			spec := minimalPlan()
			spec.ScanSetting = &AwsBackupPlanScanSetting{
				MalwareScanner: "GUARDDUTY",
				ResourceTypes:  []string{"EBS"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate selection names", func() {
			spec := minimalPlan()
			spec.Selections = []*AwsBackupPlanSelection{
				{Name: "sel", IamRoleArn: svr("arn:aws:iam::123456789012:role/backup-service")},
				{Name: "sel", IamRoleArn: svr("arn:aws:iam::123456789012:role/backup-service")},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a selection without an IAM role", func() {
			spec := minimalPlan()
			spec.Selections = []*AwsBackupPlanSelection{{Name: "sel"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a selection tag with an unknown operator", func() {
			spec := minimalPlan()
			spec.Selections = []*AwsBackupPlanSelection{{
				Name:          "sel",
				IamRoleArn:    svr("arn:aws:iam::123456789012:role/backup-service"),
				SelectionTags: []*AwsBackupPlanSelectionTag{{Type: "STRINGLIKE", Key: "backup", Value: "true"}},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
