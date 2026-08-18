package module

import (
	digitaloceanprojectv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceanproject/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module. A project has
// no tag surface, so no Planton label set applies here.
type Locals struct {
	DigitalOceanProject *digitaloceanprojectv1alpha1.DigitalOceanProject
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceanprojectv1alpha1.DigitalOceanProjectStackInput) *Locals {
	return &Locals{
		DigitalOceanProject: stackInput.Target,
	}
}
