package module

import (
	"strconv"

	awsecstaskdefinitionv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsecstaskdefinition/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsEcsTaskDefinition *awsecstaskdefinitionv1.AwsEcsTaskDefinition

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, stackInput *awsecstaskdefinitionv1.AwsEcsTaskDefinitionStackInput) *Locals {
	locals := &Locals{}
	locals.AwsEcsTaskDefinition = stackInput.Target

	metadata := stackInput.Target.Metadata
	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsEcsTaskDefinition.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
