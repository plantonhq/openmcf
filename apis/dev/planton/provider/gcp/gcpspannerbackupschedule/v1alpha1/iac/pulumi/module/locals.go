package module

import (
	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpspannerbackupschedulev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpspannerbackupschedule/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds resolved values used across the Pulumi module.
// Note: Spanner backup schedules do not support GCP labels. Labels are
// managed at the instance level only (see GcpSpannerInstance).
type Locals struct {
	GcpProviderConfig        *gcpprovider.GcpProviderConfig
	GcpSpannerBackupSchedule *gcpspannerbackupschedulev1alpha1.GcpSpannerBackupSchedule
	ScheduleName             string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpspannerbackupschedulev1alpha1.GcpSpannerBackupScheduleStackInput) *Locals {
	locals := &Locals{}
	locals.GcpSpannerBackupSchedule = stackInput.Target

	locals.ScheduleName = locals.GcpSpannerBackupSchedule.Spec.ScheduleName
	if locals.ScheduleName == "" {
		locals.ScheduleName = locals.GcpSpannerBackupSchedule.Metadata.Name
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
