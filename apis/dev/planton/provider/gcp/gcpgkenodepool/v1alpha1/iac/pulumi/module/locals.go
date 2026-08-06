package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/apis/dev/planton/provider/gcp"
	gcpgkenodepoolv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpgkenodepool/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpGkeNodePool    *gcpgkenodepoolv1alpha1.GcpGkeNodePool
	// GcpLabels carries the platform attribution labels applied as GCE
	// resource labels on the node VMs, merged over any user resource
	// labels so attribution can never be clobbered.
	GcpLabels map[string]string
	// NodePoolName is the cloud-side pool name: spec.node_pool_name when
	// set, otherwise metadata.name (the spec-level contract).
	NodePoolName string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpgkenodepoolv1alpha1.GcpGkeNodePoolStackInput) *Locals {
	locals := &Locals{}
	locals.GcpGkeNodePool = stackInput.Target

	locals.NodePoolName = locals.GcpGkeNodePool.Spec.NodePoolName
	if locals.NodePoolName == "" {
		locals.NodePoolName = locals.GcpGkeNodePool.Metadata.Name
	}

	// User resource labels merge in first so the platform attribution labels
	// can never be clobbered by a spec label with the same key. The pool
	// name (not metadata.name) keys the name label so the label matches
	// what is visible in the GCP console.
	locals.GcpLabels = map[string]string{}
	if nodeConfig := locals.GcpGkeNodePool.Spec.NodeConfig; nodeConfig != nil {
		for key, value := range nodeConfig.ResourceLabels {
			locals.GcpLabels[key] = value
		}
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.NodePoolName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpGkeNodePool.String())

	if locals.GcpGkeNodePool.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpGkeNodePool.Metadata.Org
	}
	if locals.GcpGkeNodePool.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpGkeNodePool.Metadata.Env
	}
	if locals.GcpGkeNodePool.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpGkeNodePool.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
