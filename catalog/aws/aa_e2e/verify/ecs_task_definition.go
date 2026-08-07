package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	pkgerrors "github.com/pkg/errors"
)

// ecsTaskDefinitionVerifier verifies an AwsEcsTaskDefinition via
// DescribeTaskDefinition, keyed on the revision-carrying ARN output.
//
// Deregistration does NOT delete a revision: a destroyed task definition
// stays describable as INACTIVE (AWS keeps revisions for running-task
// forensics and deletes them asynchronously later). So existence is
// "describable AND ACTIVE", and absence is "not describable, or any
// non-ACTIVE status" -- checking describability alone would report a
// destroyed revision as still existing, forever.
type ecsTaskDefinitionVerifier struct{}

func (*ecsTaskDefinitionVerifier) IDOutputKey() string { return "task_definition_arn" }

func (*ecsTaskDefinitionVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsTaskDefinitionActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecstaskdefinition verify-exists failed for %q", id)
	}
	if !active {
		return pkgerrors.Errorf("awsecstaskdefinition %q not ACTIVE after deploy", id)
	}
	return nil
}

func (*ecsTaskDefinitionVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	active, err := ecsTaskDefinitionActive(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsecstaskdefinition verify-absent failed for %q", id)
	}
	if active {
		return pkgerrors.Errorf("awsecstaskdefinition %q still ACTIVE after destroy", id)
	}
	return nil
}

func ecsTaskDefinitionActive(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	client := ecs.NewFromConfig(cfg, func(o *ecs.Options) {
		if region != "" {
			o.Region = region
		}
	})
	described, err := client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &arn,
	})
	if err != nil {
		// ECS has no typed not-found for task definitions: an unknown ARN
		// surfaces as a ClientException ("Unable to describe task
		// definition") -- treat it as absence, and anything else as a real
		// verification error.
		var clientErr *types.ClientException
		if errors.As(err, &clientErr) {
			return false, nil
		}
		return false, err
	}
	if described.TaskDefinition == nil {
		return false, nil
	}
	return described.TaskDefinition.Status == types.TaskDefinitionStatusActive, nil
}
