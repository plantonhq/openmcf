package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcppubsubtopicv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcppubsubtopic/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpPubSubTopic    *gcppubsubtopicv1alpha1.GcpPubSubTopic
	GcpLabels         map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcppubsubtopicv1alpha1.GcpPubSubTopicStackInput) *Locals {
	locals := &Locals{}
	locals.GcpPubSubTopic = stackInput.Target

	// User labels first so platform attribution labels win on key
	// conflicts — identical merge order to the Terraform module.
	locals.GcpLabels = map[string]string{}
	for key, value := range locals.GcpPubSubTopic.Spec.Labels {
		locals.GcpLabels[key] = value
	}
	locals.GcpLabels[gcplabelkeys.Resource] = "true"
	locals.GcpLabels[gcplabelkeys.ResourceName] = locals.GcpPubSubTopic.Spec.TopicName
	locals.GcpLabels[gcplabelkeys.ResourceKind] = strings.ToLower(cloudresourcekind.CloudResourceKind_GcpPubSubTopic.String())

	if locals.GcpPubSubTopic.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpPubSubTopic.Metadata.Org
	}
	if locals.GcpPubSubTopic.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpPubSubTopic.Metadata.Env
	}
	if locals.GcpPubSubTopic.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpPubSubTopic.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
