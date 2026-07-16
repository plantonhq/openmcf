package module

import (
	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpcloudsqldatabasev1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcloudsqldatabase/v1"
)

// Locals holds handy references used across this module.
//
// No label map here: google_sql_database has no labels surface in the API,
// so there is nothing to stamp — attribution lives on the instance node.
type Locals struct {
	GcpProviderConfig   *gcpprovider.GcpProviderConfig
	GcpCloudSqlDatabase *gcpcloudsqldatabasev1.GcpCloudSqlDatabase
}

// initializeLocals fills the Locals struct from the incoming stack input.
func initializeLocals(stackInput *gcpcloudsqldatabasev1.GcpCloudSqlDatabaseStackInput) *Locals {
	return &Locals{
		GcpCloudSqlDatabase: stackInput.Target,
		GcpProviderConfig:   stackInput.ProviderConfig,
	}
}
