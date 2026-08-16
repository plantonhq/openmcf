package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/backup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// reportPlan creates the Backup Audit Manager report plan and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - AWS report plan names forbid hyphens, so the name is
//     spec.report_plan_name (an explicit field), never metadata.name;
//   - report_setting.report_template is ForceNew from INSIDE the
//     nested block - changing it replaces the whole report plan;
//   - number_of_frameworks is sent only when positive (AWS computes it
//     otherwise) - the spec's zero-sentinel mirrors that contract;
//   - the destination bucket needs a policy allowing the Backup report
//     service to write (taught on the spec field).
func reportPlan(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	deliveryChannel := &backup.ReportPlanReportDeliveryChannelArgs{
		S3BucketName: pulumi.String(spec.DeliveryChannel.S3BucketName.GetValue()),
	}
	if spec.DeliveryChannel.S3KeyPrefix != "" {
		deliveryChannel.S3KeyPrefix = pulumi.String(spec.DeliveryChannel.S3KeyPrefix)
	}
	if len(spec.DeliveryChannel.Formats) > 0 {
		deliveryChannel.Formats = pulumi.ToStringArray(spec.DeliveryChannel.Formats)
	}

	reportSetting := &backup.ReportPlanReportSettingArgs{
		ReportTemplate: pulumi.String(spec.ReportSetting.ReportTemplate),
	}
	if len(spec.ReportSetting.FrameworkArns) > 0 {
		frameworkArns := make([]string, 0, len(spec.ReportSetting.FrameworkArns))
		for _, f := range spec.ReportSetting.FrameworkArns {
			frameworkArns = append(frameworkArns, f.GetValue())
		}
		reportSetting.FrameworkArns = pulumi.ToStringArray(frameworkArns)
	}
	if spec.ReportSetting.NumberOfFrameworks != 0 {
		reportSetting.NumberOfFrameworks = pulumi.Int(int(spec.ReportSetting.NumberOfFrameworks))
	}
	if len(spec.ReportSetting.Accounts) > 0 {
		reportSetting.Accounts = pulumi.ToStringArray(spec.ReportSetting.Accounts)
	}
	if len(spec.ReportSetting.OrganizationUnits) > 0 {
		reportSetting.OrganizationUnits = pulumi.ToStringArray(spec.ReportSetting.OrganizationUnits)
	}
	if len(spec.ReportSetting.Regions) > 0 {
		reportSetting.Regions = pulumi.ToStringArray(spec.ReportSetting.Regions)
	}

	args := &backup.ReportPlanArgs{
		Name:                  pulumi.String(spec.ReportPlanName),
		ReportDeliveryChannel: deliveryChannel,
		ReportSetting:         reportSetting,
		Tags:                  pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	createdReportPlan, err := backup.NewReportPlan(ctx, "report-plan", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create report plan")
	}

	ctx.Export(OpReportPlanArn, createdReportPlan.Arn)
	return nil
}
