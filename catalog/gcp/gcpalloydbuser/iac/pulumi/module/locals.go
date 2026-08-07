package module

import (
	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpalloydbuserv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpalloydbuser/v1alpha1"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpAlloydbUser    *gcpalloydbuserv1alpha1.GcpAlloydbUser
}

func initializeLocals(stackInput *gcpalloydbuserv1alpha1.GcpAlloydbUserStackInput) *Locals {
	return &Locals{
		GcpAlloydbUser:    stackInput.Target,
		GcpProviderConfig: stackInput.ProviderConfig,
	}
}
