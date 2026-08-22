package module

import (
	digitaloceandatabasedbv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceandatabasedb/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module. The logical
// database resource has no tag surface, so no Planton label set applies.
type Locals struct {
	DigitalOceanDatabaseDb *digitaloceandatabasedbv1alpha1.DigitalOceanDatabaseDb
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceandatabasedbv1alpha1.DigitalOceanDatabaseDbStackInput) *Locals {
	return &Locals{
		DigitalOceanDatabaseDb: stackInput.Target,
	}
}
