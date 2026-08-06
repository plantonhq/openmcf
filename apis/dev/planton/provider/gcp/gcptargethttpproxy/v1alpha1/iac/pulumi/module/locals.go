package module

import (
	gcptargethttpproxyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcptargethttpproxy/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs.
type Locals struct {
	GcpTargetHttpProxy *gcptargethttpproxyv1alpha1.GcpTargetHttpProxy

	// The cloud-side name defaults to metadata.name when the spec leaves
	// proxy_name empty — the same naming basis every kind uses.
	ProxyName string
}

func initializeLocals(ctx *pulumi.Context, stackInput *gcptargethttpproxyv1alpha1.GcpTargetHttpProxyStackInput) *Locals {
	target := stackInput.Target

	proxyName := target.Spec.ProxyName
	if proxyName == "" {
		proxyName = target.Metadata.Name
	}

	return &Locals{
		GcpTargetHttpProxy: target,
		ProxyName:          proxyName,
	}
}
