package module

import (
	digitaloceanmonitoralertv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceanmonitoralert/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module. The policy's
// tags SELECT tagged Droplets as alert targets (they are not resource
// labels), so no Planton label set applies here.
type Locals struct {
	DigitalOceanMonitorAlert *digitaloceanmonitoralertv1alpha1.DigitalOceanMonitorAlert
}

// initializeLocals copies stack-input fields into the Locals struct.
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceanmonitoralertv1alpha1.DigitalOceanMonitorAlertStackInput) *Locals {
	return &Locals{
		DigitalOceanMonitorAlert: stackInput.Target,
	}
}
