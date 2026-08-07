package module

import (
	azurecontainerappenvironmentdaprcomponentv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurecontainerappenvironmentdaprcomponent/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureContainerAppEnvironmentDaprComponent *azurecontainerappenvironmentdaprcomponentv1alpha1.AzureContainerAppEnvironmentDaprComponent
}

// The Dapr component carries no tags (ARM does not support them on
// managedEnvironments/daprComponents), so locals stay minimal.
func initializeLocals(ctx *pulumi.Context, stackInput *azurecontainerappenvironmentdaprcomponentv1alpha1.AzureContainerAppEnvironmentDaprComponentStackInput) *Locals {
	return &Locals{
		AzureContainerAppEnvironmentDaprComponent: stackInput.Target,
	}
}
