package module

import (
	azurecontainerappenvironmentstoragev1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurecontainerappenvironmentstorage/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureContainerAppEnvironmentStorage *azurecontainerappenvironmentstoragev1.AzureContainerAppEnvironmentStorage
}

// accessModeStrings maps the access-mode enum to ARM's values.
var accessModeStrings = map[azurecontainerappenvironmentstoragev1.AzureContainerAppEnvironmentStorageAccessMode]string{
	azurecontainerappenvironmentstoragev1.AzureContainerAppEnvironmentStorageAccessMode_READ_ONLY:  "ReadOnly",
	azurecontainerappenvironmentstoragev1.AzureContainerAppEnvironmentStorageAccessMode_READ_WRITE: "ReadWrite",
}

// The storage registration carries no tags (ARM does not support them on
// managedEnvironments/storages), so locals stay minimal.
func initializeLocals(ctx *pulumi.Context, stackInput *azurecontainerappenvironmentstoragev1.AzureContainerAppEnvironmentStorageStackInput) *Locals {
	return &Locals{
		AzureContainerAppEnvironmentStorage: stackInput.Target,
	}
}
