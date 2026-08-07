package module

import (
	gcpbigtabletablev1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpbigtabletable/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors the Terraform module's locals {} convention: the resolved
// resource plus any derived values the module needs. Bigtable tables have
// no labels surface, so no label merge exists here.
type Locals struct {
	GcpBigtableTable *gcpbigtabletablev1alpha1.GcpBigtableTable

	// Table name defaults to metadata.name when table_name is omitted.
	TableName string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpbigtabletablev1alpha1.GcpBigtableTableStackInput) *Locals {
	locals := &Locals{}
	locals.GcpBigtableTable = stackInput.Target

	locals.TableName = locals.GcpBigtableTable.Spec.TableName
	if locals.TableName == "" {
		locals.TableName = locals.GcpBigtableTable.Metadata.Name
	}

	return locals
}
