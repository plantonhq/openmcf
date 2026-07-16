package module

import (
	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpalloydbuserv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpalloydbuser/v1"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpAlloydbUser    *gcpalloydbuserv1.GcpAlloydbUser
}

func initializeLocals(stackInput *gcpalloydbuserv1.GcpAlloydbUserStackInput) *Locals {
	return &Locals{
		GcpAlloydbUser:    stackInput.Target,
		GcpProviderConfig: stackInput.ProviderConfig,
	}
}
