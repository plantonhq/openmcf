package module

import (
	"strconv"
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpvpcnetworkv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvpcnetwork/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals keeps frequently‑used values (metadata, labels, credentials) handy for the module.
type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpVpcNetwork     *gcpvpcnetworkv1.GcpVpcNetwork
	GcpLabels         map[string]string
}

// initializeLocals populates the Locals struct from the stack input.
// It mirrors the pattern used in the gcp_gke_cluster module and applies the same Planton label strategy.
func initializeLocals(_ *pulumi.Context, stackInput *gcpvpcnetworkv1.GcpVpcNetworkStackInput) *Locals {
	locals := &Locals{}

	locals.GcpVpcNetwork = stackInput.Target

	// Standard Planton‑wide labels for GCP resources
	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     strconv.FormatBool(true),
		gcplabelkeys.ResourceName: locals.GcpVpcNetwork.Spec.NetworkName,
		gcplabelkeys.ResourceKind: strings.ToLower(cloudresourcekind.CloudResourceKind_GcpVpcNetwork.String()),
	}

	if locals.GcpVpcNetwork.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpVpcNetwork.Metadata.Org
	}

	if locals.GcpVpcNetwork.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpVpcNetwork.Metadata.Env
	}

	if locals.GcpVpcNetwork.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpVpcNetwork.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig

	return locals
}
