package module

import (
	"strconv"

	awsmemorydbaclv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsmemorydbacl/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsMemorydbAcl *awsmemorydbaclv1.AwsMemorydbAcl

	// AclName is metadata.name -- the AWS ACL name is create-time immutable,
	// and metadata.name is the naming basis both engines share so a manifest
	// deploys identically on either. AWS caps it at 40 characters.
	AclName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsmemorydbaclv1.AwsMemorydbAclStackInput) *Locals {
	locals := &Locals{}
	locals.AwsMemorydbAcl = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AclName = metadata.Name

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsMemorydbAcl.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
