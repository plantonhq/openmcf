package module

import (
	"fmt"
	"strings"

	azurenetworkwatcherflowlogv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurenetworkwatcherflowlog/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureNetworkWatcherFlowLog *azurenetworkwatcherflowlogv1alpha1.AzureNetworkWatcherFlowLog

	// NetworkWatcherName and NetworkWatcherResourceGroup address the
	// regional Network Watcher the flow log attaches to. Unset spec
	// fields resolve to the AUTO-CREATED singleton -- Azure names it
	// "NetworkWatcher_{region}" and homes it in "NetworkWatcherRG" the
	// moment the region hosts a virtual network (one watcher per region
	// per subscription). Both fields override together (spec-validated)
	// for subscriptions running a self-managed watcher.
	NetworkWatcherName          string
	NetworkWatcherResourceGroup string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azurenetworkwatcherflowlogv1alpha1.AzureNetworkWatcherFlowLogStackInput) *Locals {
	locals := &Locals{}

	locals.AzureNetworkWatcherFlowLog = stackInput.Target
	target := stackInput.Target

	locals.NetworkWatcherName = target.Spec.NetworkWatcherName
	if locals.NetworkWatcherName == "" {
		locals.NetworkWatcherName = fmt.Sprintf("NetworkWatcher_%s", target.Spec.Region)
	}

	locals.NetworkWatcherResourceGroup = target.Spec.NetworkWatcherResourceGroup.GetValue()
	if locals.NetworkWatcherResourceGroup == "" {
		locals.NetworkWatcherResourceGroup = "NetworkWatcherRG"
	}

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		// PARITY-EXCEPTION: resource_kind here is the lowered
		// CloudResourceKind enum string and resource_id is omitted when
		// metadata.id is empty, while the Terraform module emits the
		// family-wide snake-case literal and falls back to metadata.name.
		// Output-neutral (tags never feed stack outputs); aligning the two
		// shapes is a family-wide convention change, not a per-kind fix.
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureNetworkWatcherFlowLog.String()),
	}

	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}

	for k, v := range target.Spec.Tags {
		locals.AzureTags[k] = v
	}

	return locals
}
