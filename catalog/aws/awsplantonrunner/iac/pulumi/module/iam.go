package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ecsTasksTrustPolicy lets ECS tasks assume a role -- the trust shape both
// of the runner's roles share.
const ecsTasksTrustPolicy = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

// ecsTaskExecutionManagedPolicyArn grants the baseline execution-role
// permissions: pull images and write CloudWatch logs.
const ecsTaskExecutionManagedPolicyArn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"

// executionRole creates the SETUP identity: the role the ECS agent (not
// the runner) assumes to pull the image, write logs, and read the token
// secret at task start. Secret access is an inline policy scoped to
// exactly the one secret ARN -- the managed execution policy deliberately
// grants no secret permissions.
func executionRole(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdSecret *secretsmanager.Secret) (*iam.Role, error) {
	runnerName := locals.AwsPlantonRunner.Metadata.Name

	createdRole, err := iam.NewRole(ctx,
		"execution-role",
		&iam.RoleArgs{
			Name:             pulumi.String(fmt.Sprintf("%s-exec", runnerName)),
			AssumeRolePolicy: pulumi.String(ecsTasksTrustPolicy),
			Description:      pulumi.String(fmt.Sprintf("ECS task execution role for Planton runner '%s'", runnerName)),
			Tags:             pulumi.ToStringMap(locals.AwsTags),
		},
		pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create execution role")
	}

	_, err = iam.NewRolePolicyAttachment(ctx,
		"execution-role-managed-policy",
		&iam.RolePolicyAttachmentArgs{
			Role:      createdRole.Name,
			PolicyArn: pulumi.String(ecsTaskExecutionManagedPolicyArn),
		},
		pulumi.Provider(provider),
		pulumi.Parent(createdRole))
	if err != nil {
		return nil, errors.Wrap(err, "failed to attach execution policy")
	}

	secretsReadPolicy := createdSecret.Arn.ApplyT(func(secretArn string) (string, error) {
		return fmt.Sprintf(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["secretsmanager:GetSecretValue"],
      "Resource": "%s"
    }
  ]
}`, secretArn), nil
	}).(pulumi.StringOutput)

	_, err = iam.NewRolePolicy(ctx,
		"execution-role-secret-read",
		&iam.RolePolicyArgs{
			Name:   pulumi.String("token-secret-read"),
			Role:   createdRole.Name,
			Policy: secretsReadPolicy,
		},
		pulumi.Provider(provider),
		pulumi.Parent(createdRole))
	if err != nil {
		return nil, errors.Wrap(err, "failed to attach secret-read policy")
	}

	return createdRole, nil
}

// runtimeRole resolves the RUNTIME identity: the role the runner itself
// holds while executing work -- the seam keyless cloud access rides. When
// the spec references a task_role, that role is the runner's identity and
// this module never touches it (modules never mutate resources they
// merely reference). Otherwise a permissionless role is created so the
// seam always exists: permissions can be granted later without replacing
// the runner.
func runtimeRole(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (pulumi.StringOutput, error) {
	spec := locals.AwsPlantonRunner.Spec
	runnerName := locals.AwsPlantonRunner.Metadata.Name

	if spec.TaskRole.GetValue() != "" {
		return pulumi.String(spec.TaskRole.GetValue()).ToStringOutput(), nil
	}

	createdRole, err := iam.NewRole(ctx,
		"runtime-role",
		&iam.RoleArgs{
			Name:             pulumi.String(fmt.Sprintf("%s-runtime", runnerName)),
			AssumeRolePolicy: pulumi.String(ecsTasksTrustPolicy),
			Description:      pulumi.String(fmt.Sprintf("Runtime identity for Planton runner '%s' -- grant it the permissions keyless operations need", runnerName)),
			Tags:             pulumi.ToStringMap(locals.AwsTags),
		},
		pulumi.Provider(provider))
	if err != nil {
		return pulumi.StringOutput{}, errors.Wrap(err, "failed to create runtime role")
	}

	return createdRole.Arn, nil
}
