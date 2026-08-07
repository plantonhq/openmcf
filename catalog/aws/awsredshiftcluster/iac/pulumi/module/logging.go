package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshift"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// clusterLogging configures audit logging -- a cluster setting keyed by
// the cluster identifier (AWS EnableLogging/DisableLogging), not a
// resource with its own identity, which is why it is folded into this
// module rather than modeled as a standalone node. Attaching it after
// cluster creation is the supported ordering.
func clusterLogging(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdCluster *redshift.Cluster) error {
	spec := locals.AwsRedshiftCluster.Spec
	if spec.Logging == nil {
		return nil
	}

	args := &redshift.LoggingArgs{
		ClusterIdentifier:  createdCluster.ClusterIdentifier,
		LogDestinationType: pulumi.String(spec.Logging.LogDestinationType),
	}

	// S3 delivery needs the bucket (with a policy granting the Redshift
	// service write access); CloudWatch delivery needs the export list.
	// The spec's CEL rules enforce each destination's requirement.
	if spec.Logging.S3BucketName != "" {
		args.BucketName = pulumi.String(spec.Logging.S3BucketName)
	}
	if spec.Logging.S3KeyPrefix != "" {
		args.S3KeyPrefix = pulumi.String(spec.Logging.S3KeyPrefix)
	}
	if len(spec.Logging.LogExports) > 0 {
		logExports := pulumi.StringArray{}
		for _, logExport := range spec.Logging.LogExports {
			logExports = append(logExports, pulumi.String(logExport))
		}
		args.LogExports = logExports
	}

	if _, err := redshift.NewLogging(ctx, "logging", args,
		pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdCluster})); err != nil {
		return errors.Wrap(err, "failed to enable audit logging")
	}
	return nil
}
