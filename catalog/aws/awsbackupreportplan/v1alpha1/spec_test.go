package awsbackupreportplanv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBackupReportPlanSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBackupReportPlanSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalReportPlan is the smallest valid instance: a backup-job
// report into a bucket.
func minimalReportPlan() *AwsBackupReportPlanSpec {
	return &AwsBackupReportPlanSpec{
		Region:         "us-west-2",
		ReportPlanName: "daily_backup_jobs",
		DeliveryChannel: &AwsBackupReportPlanDeliveryChannel{
			S3BucketName: svr("backup-reports-bucket"),
		},
		ReportSetting: &AwsBackupReportPlanReportSetting{
			ReportTemplate: "BACKUP_JOB_REPORT",
		},
	}
}

var _ = ginkgo.Describe("AwsBackupReportPlanSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal report plan", func() {
			gomega.Expect(protovalidate.Validate(minimalReportPlan())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a compliance report over frameworks with both formats", func() {
			spec := minimalReportPlan()
			spec.DeliveryChannel.S3KeyPrefix = "compliance/"
			spec.DeliveryChannel.Formats = []string{"CSV", "JSON"}
			spec.ReportSetting.ReportTemplate = "CONTROL_COMPLIANCE_REPORT"
			spec.ReportSetting.FrameworkArns = []*foreignkeyv1.StringValueOrRef{
				svr("arn:aws:backup:us-west-2:123456789012:framework:backup_posture-abc123"),
			}
			spec.ReportSetting.NumberOfFrameworks = 1
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts organization-wide coverage", func() {
			spec := minimalReportPlan()
			spec.ReportSetting.Accounts = []string{"*"}
			spec.ReportSetting.Regions = []string{"us-west-2", "us-east-1"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalReportPlan()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a hyphenated report plan name", func() {
			spec := minimalReportPlan()
			spec.ReportPlanName = "daily-backup-jobs"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing delivery channel", func() {
			spec := minimalReportPlan()
			spec.DeliveryChannel = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a delivery channel without a bucket", func() {
			spec := minimalReportPlan()
			spec.DeliveryChannel.S3BucketName = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown report format", func() {
			spec := minimalReportPlan()
			spec.DeliveryChannel.Formats = []string{"PARQUET"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing report setting", func() {
			spec := minimalReportPlan()
			spec.ReportSetting = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown report template", func() {
			spec := minimalReportPlan()
			spec.ReportSetting.ReportTemplate = "COST_REPORT"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
