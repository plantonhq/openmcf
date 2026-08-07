package module

import (
	"strings"

	gcpprovider "github.com/plantonhq/planton/catalog/gcp"
	gcpalloydbclusterv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpalloydbcluster/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/gcplabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpProviderConfig *gcpprovider.GcpProviderConfig
	GcpAlloydbCluster *gcpalloydbclusterv1alpha1.GcpAlloydbCluster
	GcpLabels         map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *gcpalloydbclusterv1alpha1.GcpAlloydbClusterStackInput) *Locals {
	locals := &Locals{}
	locals.GcpAlloydbCluster = stackInput.Target
	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     "true",
		gcplabelkeys.ResourceName: locals.GcpAlloydbCluster.Spec.ClusterName,
		gcplabelkeys.ResourceKind: strings.ToLower(cloudresourcekind.CloudResourceKind_GcpAlloydbCluster.String()),
	}

	if locals.GcpAlloydbCluster.Metadata.Org != "" {
		locals.GcpLabels[gcplabelkeys.Organization] = locals.GcpAlloydbCluster.Metadata.Org
	}
	if locals.GcpAlloydbCluster.Metadata.Env != "" {
		locals.GcpLabels[gcplabelkeys.Environment] = locals.GcpAlloydbCluster.Metadata.Env
	}
	if locals.GcpAlloydbCluster.Metadata.Id != "" {
		locals.GcpLabels[gcplabelkeys.ResourceId] = locals.GcpAlloydbCluster.Metadata.Id
	}

	locals.GcpProviderConfig = stackInput.ProviderConfig
	return locals
}
