package module

import (
	azurerediscacheaccesspolicyv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurerediscacheaccesspolicy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureRedisCacheAccessPolicy *azurerediscacheaccesspolicyv1.AzureRedisCacheAccessPolicy
	RedisCacheId                string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurerediscacheaccesspolicyv1.AzureRedisCacheAccessPolicyStackInput) *Locals {
	locals := &Locals{}

	locals.AzureRedisCacheAccessPolicy = stackInput.Target
	locals.RedisCacheId = stackInput.Target.Spec.RedisCacheId.GetValue()

	// No Azure tags: ARM does not support tags on access policies (cache
	// children), so the platform's identity tags live on the cache.

	return locals
}
