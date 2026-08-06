package module

import (
	"strconv"

	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"

	awsplantonrunnerv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsplantonrunner/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals defines local variables used throughout the Pulumi module: the
// target AwsPlantonRunner resource and the resource-identity tags applied
// to everything the module creates.
type Locals struct {
	AwsPlantonRunner *awsplantonrunnerv1alpha1.AwsPlantonRunner
	AwsTags          map[string]string
}

// initializeLocals pulls values from the stack input and populates the
// Locals struct. Similar to Terraform's "locals" concept.
func initializeLocals(ctx *pulumi.Context, stackInput *awsplantonrunnerv1alpha1.AwsPlantonRunnerStackInput) *Locals {
	locals := &Locals{
		AwsPlantonRunner: stackInput.Target,
	}

	locals.AwsTags = map[string]string{
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsPlantonRunner.Metadata.Org,
		awstagkeys.Environment:  locals.AwsPlantonRunner.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsPlantonRunner.String(),
		awstagkeys.ResourceId:   locals.AwsPlantonRunner.Metadata.Id,
	}

	return locals
}
