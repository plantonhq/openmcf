package module

import (
	"fmt"
	"strconv"

	scalewayprovider "github.com/plantonhq/planton/catalog/scaleway"
	scalewayinstancesecuritygroupv1alpha1 "github.com/plantonhq/planton/catalog/scaleway/scalewayinstancesecuritygroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/scaleway/scalewaylabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals bundles handy references for the rest of the module.
type Locals struct {
	ScalewayProviderConfig        *scalewayprovider.ScalewayProviderConfig
	ScalewayInstanceSecurityGroup *scalewayinstancesecuritygroupv1alpha1.ScalewayInstanceSecurityGroup
	ScalewayTags                  []string
}

// initializeLocals copies stack-input fields into the Locals struct and builds
// a reusable tag slice. Tags are formatted as "key=value" strings because
// Scaleway tags are flat strings (not key-value maps).
func initializeLocals(_ *pulumi.Context, stackInput *scalewayinstancesecuritygroupv1alpha1.ScalewayInstanceSecurityGroupStackInput) *Locals {
	locals := &Locals{}

	locals.ScalewayInstanceSecurityGroup = stackInput.Target
	locals.ScalewayProviderConfig = stackInput.ProviderConfig

	// Standard labels applied as Scaleway tags.
	locals.ScalewayTags = []string{
		fmt.Sprintf("%s=%s", scalewaylabelkeys.Resource, strconv.FormatBool(true)),
		fmt.Sprintf("%s=%s", scalewaylabelkeys.ResourceName, locals.ScalewayInstanceSecurityGroup.Metadata.Name),
		fmt.Sprintf("%s=%s", scalewaylabelkeys.ResourceKind, cloudresourcekind.CloudResourceKind_ScalewayInstanceSecurityGroup.String()),
	}

	if locals.ScalewayInstanceSecurityGroup.Metadata.Org != "" {
		locals.ScalewayTags = append(locals.ScalewayTags,
			fmt.Sprintf("%s=%s", scalewaylabelkeys.Organization, locals.ScalewayInstanceSecurityGroup.Metadata.Org))
	}

	if locals.ScalewayInstanceSecurityGroup.Metadata.Env != "" {
		locals.ScalewayTags = append(locals.ScalewayTags,
			fmt.Sprintf("%s=%s", scalewaylabelkeys.Environment, locals.ScalewayInstanceSecurityGroup.Metadata.Env))
	}

	if locals.ScalewayInstanceSecurityGroup.Metadata.Id != "" {
		locals.ScalewayTags = append(locals.ScalewayTags,
			fmt.Sprintf("%s=%s", scalewaylabelkeys.ResourceId, locals.ScalewayInstanceSecurityGroup.Metadata.Id))
	}

	return locals
}
