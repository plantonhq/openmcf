package module

import (
	"github.com/pkg/errors"
	awsrdsinstancev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsrdsinstance/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources provisions the RDS DB instance: a standalone database server
// with its own EBS-backed storage -- optionally Multi-AZ with a
// synchronous standby, or a read replica of another instance. Subnets,
// security groups, KMS keys, and the monitoring role compose by
// reference -- this module never creates or mutates a resource that
// deserves to be its own node.
func Resources(ctx *pulumi.Context, stackInput *awsrdsinstancev1alpha1.AwsRdsInstanceStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder,
	// which resolves the right credential mechanism (static keys, keyless
	// web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsRdsInstance.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdSubnetGroup, err := subnetGroup(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create subnet group")
	}

	createdInstance, err := rdsInstance(ctx, locals, provider, createdSubnetGroup)
	if err != nil {
		return errors.Wrap(err, "failed to create RDS instance")
	}

	ctx.Export(OpInstanceIdentifier, createdInstance.Identifier)
	ctx.Export(OpArn, createdInstance.Arn)
	ctx.Export(OpResourceId, createdInstance.ResourceId)
	ctx.Export(OpEndpoint, createdInstance.Endpoint)
	ctx.Export(OpAddress, createdInstance.Address)
	ctx.Export(OpPort, createdInstance.Port)
	ctx.Export(OpHostedZoneId, createdInstance.HostedZoneId)
	ctx.Export(OpEngineVersionActual, createdInstance.EngineVersionActual)

	// The AWS-managed master-user secret exists only when
	// manage_master_user_password is on; export empty otherwise so the
	// output shape is stable across both password strategies.
	ctx.Export(OpMasterUserSecretArn, createdInstance.MasterUserSecrets.ApplyT(func(secrets []rds.InstanceMasterUserSecret) string {
		if len(secrets) == 0 || secrets[0].SecretArn == nil {
			return ""
		}
		return *secrets[0].SecretArn
	}).(pulumi.StringOutput))

	ctx.Export(OpDbSubnetGroupName, createdInstance.DbSubnetGroupName)

	return nil
}
