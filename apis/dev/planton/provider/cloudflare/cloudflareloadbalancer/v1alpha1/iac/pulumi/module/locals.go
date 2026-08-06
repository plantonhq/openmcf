package module

import (
	cloudflareprovider "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare"
	cloudflareloadbalancerv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare/cloudflareloadbalancer/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals stores quick references for metadata, spec & credentials.
type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareLoadBalancer   *cloudflareloadbalancerv1alpha1.CloudflareLoadBalancer
}

// initializeLocals copies relevant stack‑input fields into Locals.
func initializeLocals(_ *pulumi.Context, stackInput *cloudflareloadbalancerv1alpha1.CloudflareLoadBalancerStackInput) *Locals {
	return &Locals{
		CloudflareProviderConfig: stackInput.ProviderConfig,
		CloudflareLoadBalancer:   stackInput.Target,
	}
}
