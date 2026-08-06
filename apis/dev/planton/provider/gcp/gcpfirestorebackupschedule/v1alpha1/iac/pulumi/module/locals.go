package module

import (
	gcpfirestorebackupschedulev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpfirestorebackupschedule/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds resolved values used across the Pulumi module.
// Firestore backup schedules do not support GCP labels — skip label merge.
type Locals struct {
	GcpFirestoreBackupSchedule *gcpfirestorebackupschedulev1alpha1.GcpFirestoreBackupSchedule
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpfirestorebackupschedulev1alpha1.GcpFirestoreBackupScheduleStackInput) *Locals {
	return &Locals{
		GcpFirestoreBackupSchedule: stackInput.Target,
	}
}
