package module

import (
	"strconv"

	awseksaccessentryv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awseksaccessentry/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsEksAccessEntry *awseksaccessentryv1alpha1.AwsEksAccessEntry

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awseksaccessentryv1alpha1.AwsEksAccessEntryStackInput) *Locals {
	locals := &Locals{}
	locals.AwsEksAccessEntry = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEksAccessEntry.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
