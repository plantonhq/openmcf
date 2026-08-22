package module

import (
	digitaloceansshkeyv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceansshkey/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module. An SSH key
// has no tag surface, so no Planton label set applies here.
type Locals struct {
	DigitalOceanSshKey *digitaloceansshkeyv1alpha1.DigitalOceanSshKey
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceansshkeyv1alpha1.DigitalOceanSshKeyStackInput) *Locals {
	return &Locals{
		DigitalOceanSshKey: stackInput.Target,
	}
}
