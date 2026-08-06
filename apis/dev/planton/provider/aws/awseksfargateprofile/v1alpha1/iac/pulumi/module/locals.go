package module

import (
	"strconv"

	awseksfargateprofilev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awseksfargateprofile/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsEksFargateProfile *awseksfargateprofilev1alpha1.AwsEksFargateProfile

	// FargateProfileName is metadata.name truncated to AWS's 63-character
	// profile limit, deterministically, so the same manifest always yields
	// the same name on both engines.
	FargateProfileName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awseksfargateprofilev1alpha1.AwsEksFargateProfileStackInput) *Locals {
	locals := &Locals{}
	locals.AwsEksFargateProfile = stackInput.Target

	locals.FargateProfileName = stackInput.Target.Metadata.Name
	if len(locals.FargateProfileName) > 63 {
		locals.FargateProfileName = locals.FargateProfileName[:63]
	}

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEksFargateProfile.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
