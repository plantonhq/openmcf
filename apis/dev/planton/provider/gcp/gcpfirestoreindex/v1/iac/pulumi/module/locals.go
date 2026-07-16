package module

import (
	gcpfirestoreindexv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpfirestoreindex/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds resolved values used across the Pulumi module.
// Firestore indexes do not support GCP labels — skip label merge.
type Locals struct {
	GcpFirestoreIndex *gcpfirestoreindexv1.GcpFirestoreIndex
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpfirestoreindexv1.GcpFirestoreIndexStackInput) *Locals {
	return &Locals{
		GcpFirestoreIndex: stackInput.Target,
	}
}
