package module

import (
	"strconv"

	awsmemorydbuserv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsmemorydbuser/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsMemorydbUser *awsmemorydbuserv1alpha1.AwsMemorydbUser

	// UserName is metadata.name -- in MemoryDB the user name IS the user's
	// single identity (there is no separate user id), it is create-time
	// immutable, and metadata.name is the naming basis both engines share
	// so a manifest deploys identically on either. AWS caps it at 40
	// characters.
	UserName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsmemorydbuserv1alpha1.AwsMemorydbUserStackInput) *Locals {
	locals := &Locals{}
	locals.AwsMemorydbUser = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.UserName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsMemorydbUser.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
