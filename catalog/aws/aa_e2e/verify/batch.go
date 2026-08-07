package verify

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
	pkgerrors "github.com/pkg/errors"
)

// AWS Batch resources have lifecycle-state semantics rather than hard
// deletes: compute environments and job queues are disable-then-delete
// (passing through DELETING and briefly describable as DELETED), and job
// definitions are never deleted at all -- deregistration flips every
// revision to INACTIVE while DescribeJobDefinitions keeps returning them.
// Every verifier here therefore treats the terminal lifecycle states as
// "absent", the NAT-gateway/ECS class of state-aware verification.

func batchClient(cfg aws.Config, region string) *batch.Client {
	if region != "" {
		cfg.Region = region
	}
	return batch.NewFromConfig(cfg)
}

// batchComputeEnvironmentVerifier verifies an AwsBatchComputeEnvironment via
// DescribeComputeEnvironments, keyed on the compute_environment_arn output.
type batchComputeEnvironmentVerifier struct{}

func (*batchComputeEnvironmentVerifier) IDOutputKey() string { return "compute_environment_arn" }

func (*batchComputeEnvironmentVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := batchComputeEnvironmentExists(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbatchcomputeenvironment verify-exists failed for %q", arn)
	}
	if !exists {
		return pkgerrors.Errorf("awsbatchcomputeenvironment %q not found after deploy", arn)
	}
	return nil
}

func (*batchComputeEnvironmentVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := batchComputeEnvironmentExists(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbatchcomputeenvironment verify-absent failed for %q", arn)
	}
	if exists {
		return pkgerrors.Errorf("awsbatchcomputeenvironment %q still exists after destroy", arn)
	}
	return nil
}

func batchComputeEnvironmentExists(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	out, err := batchClient(cfg, region).DescribeComputeEnvironments(ctx, &batch.DescribeComputeEnvironmentsInput{
		ComputeEnvironments: []string{arn},
	})
	if err != nil {
		return false, err
	}
	for _, ce := range out.ComputeEnvironments {
		switch ce.Status {
		case batchtypes.CEStatusDeleting, batchtypes.CEStatusDeleted:
			// terminal lifecycle -- absent
		default:
			return true, nil
		}
	}
	return false, nil
}

// batchJobQueueVerifier verifies an AwsBatchJobQueue via DescribeJobQueues,
// keyed on the job_queue_arn output.
type batchJobQueueVerifier struct{}

func (*batchJobQueueVerifier) IDOutputKey() string { return "job_queue_arn" }

func (*batchJobQueueVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := batchJobQueueExists(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbatchjobqueue verify-exists failed for %q", arn)
	}
	if !exists {
		return pkgerrors.Errorf("awsbatchjobqueue %q not found after deploy", arn)
	}
	return nil
}

func (*batchJobQueueVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := batchJobQueueExists(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbatchjobqueue verify-absent failed for %q", arn)
	}
	if exists {
		return pkgerrors.Errorf("awsbatchjobqueue %q still exists after destroy", arn)
	}
	return nil
}

func batchJobQueueExists(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	out, err := batchClient(cfg, region).DescribeJobQueues(ctx, &batch.DescribeJobQueuesInput{
		JobQueues: []string{arn},
	})
	if err != nil {
		return false, err
	}
	for _, queue := range out.JobQueues {
		switch queue.Status {
		case batchtypes.JQStatusDeleting, batchtypes.JQStatusDeleted:
			// terminal lifecycle -- absent
		default:
			return true, nil
		}
	}
	return false, nil
}

// batchSchedulingPolicyVerifier verifies an AwsBatchSchedulingPolicy via
// DescribeSchedulingPolicies, keyed on the scheduling_policy_arn output.
// Scheduling policies delete synchronously: absent means an empty result.
type batchSchedulingPolicyVerifier struct{}

func (*batchSchedulingPolicyVerifier) IDOutputKey() string { return "scheduling_policy_arn" }

func (*batchSchedulingPolicyVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := batchSchedulingPolicyExists(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbatchschedulingpolicy verify-exists failed for %q", arn)
	}
	if !exists {
		return pkgerrors.Errorf("awsbatchschedulingpolicy %q not found after deploy", arn)
	}
	return nil
}

func (*batchSchedulingPolicyVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := batchSchedulingPolicyExists(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbatchschedulingpolicy verify-absent failed for %q", arn)
	}
	if exists {
		return pkgerrors.Errorf("awsbatchschedulingpolicy %q still exists after destroy", arn)
	}
	return nil
}

func batchSchedulingPolicyExists(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	out, err := batchClient(cfg, region).DescribeSchedulingPolicies(ctx, &batch.DescribeSchedulingPoliciesInput{
		Arns: []string{arn},
	})
	if err != nil {
		return false, err
	}
	return len(out.SchedulingPolicies) > 0, nil
}

// batchJobDefinitionVerifier verifies an AwsBatchJobDefinition via
// DescribeJobDefinitions, keyed on the revision-carrying job_definition_arn
// output. Deregistered revisions stay describable as INACTIVE forever, so
// only an ACTIVE revision counts as existing.
type batchJobDefinitionVerifier struct{}

func (*batchJobDefinitionVerifier) IDOutputKey() string { return "job_definition_arn" }

func (*batchJobDefinitionVerifier) VerifyExists(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := batchJobDefinitionActive(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbatchjobdefinition verify-exists failed for %q", arn)
	}
	if !exists {
		return pkgerrors.Errorf("awsbatchjobdefinition %q not ACTIVE after deploy", arn)
	}
	return nil
}

func (*batchJobDefinitionVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, arn, region string) error {
	exists, err := batchJobDefinitionActive(ctx, cfg, arn, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsbatchjobdefinition verify-absent failed for %q", arn)
	}
	if exists {
		return pkgerrors.Errorf("awsbatchjobdefinition %q still ACTIVE after destroy", arn)
	}
	return nil
}

func batchJobDefinitionActive(ctx context.Context, cfg aws.Config, arn, region string) (bool, error) {
	out, err := batchClient(cfg, region).DescribeJobDefinitions(ctx, &batch.DescribeJobDefinitionsInput{
		JobDefinitions: []string{arn},
	})
	if err != nil {
		return false, err
	}
	for _, jd := range out.JobDefinitions {
		if jd.Status != nil && strings.EqualFold(*jd.Status, "ACTIVE") {
			return true, nil
		}
	}
	return false, nil
}
