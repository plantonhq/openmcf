package module

import (
	"strconv"

	awsathenaworkgroup "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsathenaworkgroup/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target         *awsathenaworkgroup.AwsAthenaWorkgroup
	Spec           *awsathenaworkgroup.AwsAthenaWorkgroupSpec
	AwsTags        map[string]string
	WorkgroupName  string
	WorkgroupState string
}

func initializeLocals(ctx *pulumi.Context, in *awsathenaworkgroup.AwsAthenaWorkgroupStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	// The workgroup's cloud name is metadata.name; the name is create-time
	// immutable (changing it replaces the workgroup), which is why it is not
	// spec surface. Same basis as the Terraform module.
	locals.WorkgroupName = in.Target.Metadata.Name

	// The spec defaults state to ENABLED; an omitted value must deploy the
	// same workgroup an explicit ENABLED would.
	locals.WorkgroupState = "ENABLED"
	if s := in.Target.Spec.GetState(); s != "" {
		locals.WorkgroupState = s
	}

	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.Target.Metadata.Org,
		awstagkeys.Environment:  locals.Target.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsAthenaWorkgroup.String(),
		awstagkeys.ResourceId:   locals.Target.Metadata.Id,
	}

	return locals
}
