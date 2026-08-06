package module

import (
	"strconv"
	"strings"

	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpcomputeinstancev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcomputeinstance/v1alpha1"
)

// Locals holds handy references and derived values used across this module.
type Locals struct {
	GcpProviderConfig  *gcpprovider.GcpProviderConfig
	GcpComputeInstance *gcpcomputeinstancev1alpha1.GcpComputeInstance
	GcpLabels          map[string]string
	// InstanceName is the cloud-side name: spec.instance_name when set,
	// metadata.name otherwise — the same explicit conditional as the
	// Terraform module, so both engines derive the identical name.
	InstanceName string
}

// initializeLocals fills the Locals struct from the incoming stack input.
func initializeLocals(stackInput *gcpcomputeinstancev1alpha1.GcpComputeInstanceStackInput) *Locals {
	locals := &Locals{}

	locals.GcpComputeInstance = stackInput.Target
	locals.GcpProviderConfig = stackInput.ProviderConfig

	target := stackInput.Target

	locals.InstanceName = target.Spec.InstanceName
	if locals.InstanceName == "" {
		locals.InstanceName = target.Metadata.Name
	}

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range target.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = strconv.FormatBool(true)
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.InstanceName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpComputeInstance.String())

	if target.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = target.Metadata.Env
	}

	return locals
}
