package module

import (
	"strconv"

	awsclientvpnv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsclientvpn/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsClientVpn *awsclientvpnv1alpha1.AwsClientVpn

	// EndpointName is metadata.name -- Client VPN endpoints have no AWS
	// name argument (identity is the generated cvpn-endpoint-... id), so
	// the name lives on the Name tag, identically on both engines.
	EndpointName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsclientvpnv1alpha1.AwsClientVpnStackInput) *Locals {
	locals := &Locals{}
	locals.AwsClientVpn = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.EndpointName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsClientVpn.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
