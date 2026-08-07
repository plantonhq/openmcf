package module

import (
	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpcloudsqluserv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcloudsqluser/v1alpha1"
)

// Locals holds handy references used across this module.
//
// No label map here: google_sql_user has no labels surface in the API, so
// there is nothing to stamp — attribution lives on the instance node.
type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpCloudSqlUser   *gcpcloudsqluserv1alpha1.GcpCloudSqlUser
}

// initializeLocals fills the Locals struct from the incoming stack input.
func initializeLocals(stackInput *gcpcloudsqluserv1alpha1.GcpCloudSqlUserStackInput) *Locals {
	return &Locals{
		GcpCloudSqlUser:   stackInput.Target,
		GcpProviderConfig: stackInput.ProviderConfig,
	}
}
