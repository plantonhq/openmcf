package module

import (
	azureredislinkedserverv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureredislinkedserver/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureRedisLinkedServer *azureredislinkedserverv1.AzureRedisLinkedServer
	TargetRedisCacheId     string
	LinkedRedisCacheId     string
}

// serverRoleStrings maps the spec's role enum to ARM's capitalized values.
var serverRoleStrings = map[azureredislinkedserverv1.AzureRedisLinkedServerRole]string{
	azureredislinkedserverv1.AzureRedisLinkedServerRole_PRIMARY:   "Primary",
	azureredislinkedserverv1.AzureRedisLinkedServerRole_SECONDARY: "Secondary",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureredislinkedserverv1.AzureRedisLinkedServerStackInput) *Locals {
	locals := &Locals{}

	locals.AzureRedisLinkedServer = stackInput.Target
	locals.TargetRedisCacheId = stackInput.Target.Spec.TargetRedisCacheId.GetValue()
	locals.LinkedRedisCacheId = stackInput.Target.Spec.LinkedRedisCacheId.GetValue()

	// No Azure tags: ARM does not support tags on linked servers, so the
	// platform's identity tags live on the caches themselves.

	return locals
}
