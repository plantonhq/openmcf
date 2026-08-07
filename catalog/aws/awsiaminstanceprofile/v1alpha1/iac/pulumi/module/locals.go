package module

import (
	"strconv"

	awsiaminstanceprofilev1alpha1 "github.com/plantonhq/planton/catalog/aws/awsiaminstanceprofile/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsIamInstanceProfile *awsiaminstanceprofilev1alpha1.AwsIamInstanceProfile
	AwsTags               map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsiaminstanceprofilev1alpha1.AwsIamInstanceProfileStackInput) *Locals {
	locals := &Locals{}
	locals.AwsIamInstanceProfile = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsIamInstanceProfile.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
