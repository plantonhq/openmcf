package module

import (
	azurecontainerappenvironmentdaprcomponentv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurecontainerappenvironmentdaprcomponent/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureContainerAppEnvironmentDaprComponent *azurecontainerappenvironmentdaprcomponentv1.AzureContainerAppEnvironmentDaprComponent
}

// The Dapr component carries no tags (ARM does not support them on
// managedEnvironments/daprComponents), so locals stay minimal.
func initializeLocals(ctx *pulumi.Context, stackInput *azurecontainerappenvironmentdaprcomponentv1.AzureContainerAppEnvironmentDaprComponentStackInput) *Locals {
	return &Locals{
		AzureContainerAppEnvironmentDaprComponent: stackInput.Target,
	}
}
