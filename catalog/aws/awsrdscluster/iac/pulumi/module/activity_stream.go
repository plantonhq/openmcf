package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// activityStream starts the Database Activity Stream when the spec asks
// for one: AWS creates and owns a Kinesis stream that receives every
// audited database event, encrypted with the given KMS key.
// Create/delete-only lifecycle -- every argument forces replacement,
// and AWS walks the stream through starting/stopping states with a
// waiter on each side, so a slow update or destroy here is the
// service's state machine, not this module. The stream needs an
// available instance to start against, hence the explicit dependency
// on the folded instances.
func activityStream(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdCluster *rds.Cluster, createdInstances []*rds.ClusterInstance) (*rds.ClusterActivityStream, error) {
	spec := locals.AwsRdsCluster.Spec
	if spec.ActivityStream == nil {
		return nil, nil
	}

	dependsOn := make([]pulumi.Resource, 0, len(createdInstances))
	for _, createdInstance := range createdInstances {
		dependsOn = append(dependsOn, createdInstance)
	}

	createdStream, err := rds.NewClusterActivityStream(ctx, "activity-stream",
		&rds.ClusterActivityStreamArgs{
			ResourceArn:                     createdCluster.Arn,
			Mode:                            pulumi.String(spec.ActivityStream.Mode),
			KmsKeyId:                        pulumi.String(spec.ActivityStream.KmsKeyId.GetValue()),
			EngineNativeAuditFieldsIncluded: pulumi.Bool(spec.ActivityStream.EngineNativeAuditFieldsIncluded),
		},
		pulumi.Provider(provider), pulumi.DependsOn(dependsOn))
	if err != nil {
		return nil, errors.Wrap(err, "failed to start database activity stream")
	}
	return createdStream, nil
}
