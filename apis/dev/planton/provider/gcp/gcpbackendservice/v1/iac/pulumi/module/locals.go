package module

import (
	gcpbackendservicev1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpbackendservice/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpBackendService *gcpbackendservicev1.GcpBackendService

	// The cloud-side name defaults to metadata.name when the spec leaves
	// backend_service_name empty — the same naming basis every kind uses.
	BackendServiceName string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpbackendservicev1.GcpBackendServiceStackInput) *Locals {
	target := stackInput.Target

	backendServiceName := target.Spec.BackendServiceName
	if backendServiceName == "" {
		backendServiceName = target.Metadata.Name
	}

	return &Locals{
		GcpBackendService:  target,
		BackendServiceName: backendServiceName,
	}
}
