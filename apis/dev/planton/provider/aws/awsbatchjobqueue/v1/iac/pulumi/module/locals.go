package module

import (
	"strconv"

	awsbatchjobqueuev1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsbatchjobqueue/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	AwsBatchJobQueue *awsbatchjobqueuev1.AwsBatchJobQueue
	AwsTags          map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awsbatchjobqueuev1.AwsBatchJobQueueStackInput) *Locals {
	locals := &Locals{}
	locals.AwsBatchJobQueue = stackInput.Target

	// Resource-identity tags follow the catalog convention.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         locals.AwsBatchJobQueue.Metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: locals.AwsBatchJobQueue.Metadata.Org,
		awstagkeys.Environment:  locals.AwsBatchJobQueue.Metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsBatchJobQueue.String(),
		awstagkeys.ResourceId:   locals.AwsBatchJobQueue.Metadata.Id,
	}

	return locals
}
