package module

import (
	"strconv"

	hetznercloudprovider "github.com/plantonhq/planton/catalog/hetznercloud"
	hetznercloudcertificatev1alpha1 "github.com/plantonhq/planton/catalog/hetznercloud/hetznercloudcertificate/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/hetznercloud/hcloudlabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	HetznerCloudProviderConfig *hetznercloudprovider.HetznerCloudProviderConfig
	HetznerCloudCertificate    *hetznercloudcertificatev1alpha1.HetznerCloudCertificate
	Labels                     map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *hetznercloudcertificatev1alpha1.HetznerCloudCertificateStackInput) *Locals {
	locals := &Locals{}

	locals.HetznerCloudCertificate = stackInput.Target
	locals.HetznerCloudProviderConfig = stackInput.ProviderConfig

	locals.Labels = map[string]string{
		hcloudlabelkeys.Resource:     strconv.FormatBool(true),
		hcloudlabelkeys.ResourceName: locals.HetznerCloudCertificate.Metadata.Name,
		hcloudlabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_HetznerCloudCertificate.String(),
	}

	if locals.HetznerCloudCertificate.Metadata.Org != "" {
		locals.Labels[hcloudlabelkeys.Organization] = locals.HetznerCloudCertificate.Metadata.Org
	}

	if locals.HetznerCloudCertificate.Metadata.Env != "" {
		locals.Labels[hcloudlabelkeys.Environment] = locals.HetznerCloudCertificate.Metadata.Env
	}

	if locals.HetznerCloudCertificate.Metadata.Id != "" {
		locals.Labels[hcloudlabelkeys.ResourceId] = locals.HetznerCloudCertificate.Metadata.Id
	}

	for k, v := range locals.HetznerCloudCertificate.Metadata.Labels {
		if _, exists := locals.Labels[k]; !exists {
			locals.Labels[k] = v
		}
	}

	return locals
}
