package module

import (
	"strconv"

	hetznercloudprovider "github.com/plantonhq/planton/catalog/hetznercloud"
	hetznerclouddnszonev1alpha1 "github.com/plantonhq/planton/catalog/hetznercloud/hetznerclouddnszone/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/hetznercloud/hcloudlabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	HetznerCloudProviderConfig *hetznercloudprovider.HetznerCloudProviderConfig
	HetznerCloudDnsZone        *hetznerclouddnszonev1alpha1.HetznerCloudDnsZone
	Labels                     map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *hetznerclouddnszonev1alpha1.HetznerCloudDnsZoneStackInput) *Locals {
	locals := &Locals{}

	locals.HetznerCloudDnsZone = stackInput.Target
	locals.HetznerCloudProviderConfig = stackInput.ProviderConfig

	locals.Labels = map[string]string{
		hcloudlabelkeys.Resource:     strconv.FormatBool(true),
		hcloudlabelkeys.ResourceName: locals.HetznerCloudDnsZone.Metadata.Name,
		hcloudlabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_HetznerCloudDnsZone.String(),
	}

	if locals.HetznerCloudDnsZone.Metadata.Org != "" {
		locals.Labels[hcloudlabelkeys.Organization] = locals.HetznerCloudDnsZone.Metadata.Org
	}

	if locals.HetznerCloudDnsZone.Metadata.Env != "" {
		locals.Labels[hcloudlabelkeys.Environment] = locals.HetznerCloudDnsZone.Metadata.Env
	}

	if locals.HetznerCloudDnsZone.Metadata.Id != "" {
		locals.Labels[hcloudlabelkeys.ResourceId] = locals.HetznerCloudDnsZone.Metadata.Id
	}

	for k, v := range locals.HetznerCloudDnsZone.Metadata.Labels {
		if _, exists := locals.Labels[k]; !exists {
			locals.Labels[k] = v
		}
	}

	return locals
}
