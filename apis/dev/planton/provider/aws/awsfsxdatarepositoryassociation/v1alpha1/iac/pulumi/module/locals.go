package module

import (
	"strconv"

	awsfsxdrav1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsfsxdatarepositoryassociation/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds the association definition from the stack input and the AWS
// tags applied to the created resource.
type Locals struct {
	AwsFsxDataRepositoryAssociation *awsfsxdrav1.AwsFsxDataRepositoryAssociation
	AwsTags                         map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsfsxdrav1.AwsFsxDataRepositoryAssociationStackInput) *Locals {
	locals := &Locals{}
	locals.AwsFsxDataRepositoryAssociation = stackInput.Target

	// Resource-identity tags follow the catalog convention. Associations have
	// no cloud name argument, so the Name tag carries metadata.name — the
	// same basis the Terraform module pins, keeping the two engines' physical
	// identity converged.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsFsxDataRepositoryAssociation.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsFsxDataRepositoryAssociation.Metadata.Org,
		awstagkeys.Environment:  locals.AwsFsxDataRepositoryAssociation.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsFsxDataRepositoryAssociation.String(),
		awstagkeys.ResourceId:   locals.AwsFsxDataRepositoryAssociation.Metadata.Id,
	}

	return locals
}
