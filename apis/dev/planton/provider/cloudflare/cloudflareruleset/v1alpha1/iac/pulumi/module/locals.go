package module

import (
	cloudflareprovider "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare"
	cloudflarerulesetv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/cloudflare/cloudflareruleset/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	CloudflareProviderConfig *cloudflareprovider.CloudflareProviderConfig
	CloudflareRuleset        *cloudflarerulesetv1alpha1.CloudflareRuleset
}

func initializeLocals(_ *pulumi.Context, stackInput *cloudflarerulesetv1alpha1.CloudflareRulesetStackInput) *Locals {
	return &Locals{
		CloudflareRuleset:        stackInput.Target,
		CloudflareProviderConfig: stackInput.ProviderConfig,
	}
}
