package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcppubsubsubscriptionv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcppubsubsubscription/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig     *gcpprovider.GcpProviderConfig
	GcpPubSubSubscription *gcppubsubsubscriptionv1alpha1.GcpPubSubSubscription
	GcpLabels             map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcppubsubsubscriptionv1alpha1.GcpPubSubSubscriptionStackInput) *Locals {
	locals := &Locals{}
	locals.GcpPubSubSubscription = stackInput.Target

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpPubSubSubscription.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.GcpPubSubSubscription.Spec.SubscriptionName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpPubSubSubscription.String())

	if locals.GcpPubSubSubscription.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpPubSubSubscription.Metadata.Org
	}
	if locals.GcpPubSubSubscription.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpPubSubSubscription.Metadata.Env
	}
	if locals.GcpPubSubSubscription.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpPubSubSubscription.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
