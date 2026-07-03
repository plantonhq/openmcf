package module

import (
	gcpbackendbucketv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpbackendbucket/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpBackendBucket *gcpbackendbucketv1.GcpBackendBucket

	// The cloud-side name defaults to metadata.name when the spec leaves
	// backend_bucket_name empty — the same naming basis every kind uses.
	BackendBucketName string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpbackendbucketv1.GcpBackendBucketStackInput) *Locals {
	target := stackInput.Target

	backendBucketName := target.Spec.BackendBucketName
	if backendBucketName == "" {
		backendBucketName = target.Metadata.Name
	}

	return &Locals{
		GcpBackendBucket:  target,
		BackendBucketName: backendBucketName,
	}
}
