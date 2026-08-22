package module

import (
	"strconv"

	digitaloceandatabasereplicav1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceandatabasereplica/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/digitalocean/digitaloceanlabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	DigitalOceanDatabaseReplica *digitaloceandatabasereplicav1alpha1.DigitalOceanDatabaseReplica
	DigitalOceanLabels          map[string]string
}

// initializeLocals copies stack-input fields into the Locals struct and
// builds the standard Planton label map (rendered as "key:value" tags on
// the replica -- the identical set the Terraform module applies).
func initializeLocals(_ *pulumi.Context, stackInput *digitaloceandatabasereplicav1alpha1.DigitalOceanDatabaseReplicaStackInput) *Locals {
	locals := &Locals{}

	locals.DigitalOceanDatabaseReplica = stackInput.Target

	locals.DigitalOceanLabels = map[string]string{
		digitaloceanlabelkeys.Resource:     strconv.FormatBool(true),
		digitaloceanlabelkeys.ResourceName: locals.DigitalOceanDatabaseReplica.Metadata.Name,
		digitaloceanlabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_DigitalOceanDatabaseReplica.String(),
	}

	if locals.DigitalOceanDatabaseReplica.Metadata.Org != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Organization] = locals.DigitalOceanDatabaseReplica.Metadata.Org
	}

	if locals.DigitalOceanDatabaseReplica.Metadata.Env != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.Environment] = locals.DigitalOceanDatabaseReplica.Metadata.Env
	}

	if locals.DigitalOceanDatabaseReplica.Metadata.Id != "" {
		locals.DigitalOceanLabels[digitaloceanlabelkeys.ResourceId] = locals.DigitalOceanDatabaseReplica.Metadata.Id
	}

	return locals
}
