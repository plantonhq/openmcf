package module

import (
	gcpsslpolicyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpsslpolicy/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpSslPolicy *gcpsslpolicyv1alpha1.GcpSslPolicy

	// The cloud-side name defaults to metadata.name when the spec leaves
	// ssl_policy_name empty — the same naming basis every kind uses.
	SslPolicyName string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcpsslpolicyv1alpha1.GcpSslPolicyStackInput) *Locals {
	target := stackInput.Target

	sslPolicyName := target.Spec.SslPolicyName
	if sslPolicyName == "" {
		sslPolicyName = target.Metadata.Name
	}

	return &Locals{
		GcpSslPolicy:  target,
		SslPolicyName: sslPolicyName,
	}
}
