package module

import (
	"strconv"

	awstgwv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awstransitgateway/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals mirrors Terraform-style locals: the target resource and the identity
// tag set applied to the Transit Gateway.
type Locals struct {
	TransitGateway *awstgwv1.AwsTransitGateway
	AwsTags        map[string]string
}

// initializeLocals reads the stack input and builds the Locals instance.
func initializeLocals(ctx *pulumi.Context, stackInput *awstgwv1.AwsTransitGatewayStackInput) *Locals {
	locals := &Locals{}

	locals.TransitGateway = stackInput.Target

	// Identity tags match the Terraform module key-for-key. The Name tag IS
	// the gateway's console identity -- Transit Gateways have no name
	// attribute of their own.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.TransitGateway.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.TransitGateway.Metadata.Org,
		awstagkeys.Environment:  locals.TransitGateway.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsTransitGateway.String(),
		awstagkeys.ResourceId:   locals.TransitGateway.Metadata.Id,
	}

	return locals
}

// enableDisable converts a plain boolean to the AWS-style "enable"/"disable"
// string the Transit Gateway API speaks.
func enableDisable(b bool) string {
	if b {
		return "enable"
	}
	return "disable"
}

// enableDisableTriState converts a proto `optional bool` to the AWS-style
// string, preserving the tri-state: nil (unset in the manifest) returns nil
// so the argument is omitted and the provider/AWS default applies, exactly
// like the Terraform module's null fall-through.
func enableDisableTriState(b *bool) pulumi.StringPtrInput {
	if b == nil {
		return nil
	}
	return pulumi.StringPtr(enableDisable(*b))
}
