package module

import (
	azuredatafactoryintegrationruntimev1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactoryintegrationruntime/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureDataFactoryIntegrationRuntime *azuredatafactoryintegrationruntimev1alpha1.AzureDataFactoryIntegrationRuntime
}

// An integration runtime carries no tags (ARM sub-resources of a
// factory expose none), so there is no tag map to derive.
func initializeLocals(ctx *pulumi.Context, stackInput *azuredatafactoryintegrationruntimev1alpha1.AzureDataFactoryIntegrationRuntimeStackInput) *Locals {
	locals := &Locals{}

	locals.AzureDataFactoryIntegrationRuntime = stackInput.Target

	return locals
}
