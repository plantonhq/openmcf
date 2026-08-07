package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cluster provisions the ECS cluster itself: Container Insights, ECS Exec
// auditing, Fargate storage encryption, and the Service Connect default
// namespace. Only the name is create-time; every posture block edits in
// place via UpdateCluster.
func cluster(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*ecs.Cluster, error) {
	spec := locals.AwsEcsCluster.Spec

	args := &ecs.ClusterArgs{
		// The cluster name is metadata.name on both engines, so a
		// manifest deploys the same cluster regardless of engine.
		Name: pulumi.String(locals.ClusterName),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// Unset keeps the account's default setting (AWS computes the
	// effective value); the spec values are the exact AWS API strings.
	if spec.ContainerInsights != "" {
		args.Settings = ecs.ClusterSettingArray{
			ecs.ClusterSettingArgs{
				Name:  pulumi.String("containerInsights"),
				Value: pulumi.String(spec.ContainerInsights),
			},
		}
	}

	// Exec auditing and managed-storage encryption share the provider's
	// configuration block; build it only when either is present so an
	// empty block never overwrites AWS defaults.
	if spec.ExecuteCommandConfiguration != nil || spec.ManagedStorageConfiguration != nil {
		configurationArgs := &ecs.ClusterConfigurationArgs{}

		if execConfiguration := spec.ExecuteCommandConfiguration; execConfiguration != nil {
			execArgs := &ecs.ClusterConfigurationExecuteCommandConfigurationArgs{}
			if execConfiguration.Logging != "" {
				execArgs.Logging = pulumi.String(execConfiguration.Logging)
			}
			if execConfiguration.KmsKeyId.GetValue() != "" {
				execArgs.KmsKeyId = pulumi.String(execConfiguration.KmsKeyId.GetValue())
			}
			// Present only with OVERRIDE logging (CEL enforces the
			// coupling both directions).
			if logConfiguration := execConfiguration.LogConfiguration; logConfiguration != nil {
				logArgs := &ecs.ClusterConfigurationExecuteCommandConfigurationLogConfigurationArgs{}
				if logConfiguration.CloudWatchLogGroupName != "" {
					logArgs.CloudWatchLogGroupName = pulumi.String(logConfiguration.CloudWatchLogGroupName)
				}
				if logConfiguration.CloudWatchEncryptionEnabled {
					logArgs.CloudWatchEncryptionEnabled = pulumi.Bool(true)
				}
				if logConfiguration.S3BucketName != "" {
					logArgs.S3BucketName = pulumi.String(logConfiguration.S3BucketName)
				}
				if logConfiguration.S3KeyPrefix != "" {
					logArgs.S3KeyPrefix = pulumi.String(logConfiguration.S3KeyPrefix)
				}
				if logConfiguration.S3BucketEncryptionEnabled {
					logArgs.S3BucketEncryptionEnabled = pulumi.Bool(true)
				}
				execArgs.LogConfiguration = logArgs
			}
			configurationArgs.ExecuteCommandConfiguration = execArgs
		}

		if storageConfiguration := spec.ManagedStorageConfiguration; storageConfiguration != nil {
			storageArgs := &ecs.ClusterConfigurationManagedStorageConfigurationArgs{}
			if storageConfiguration.FargateEphemeralStorageKmsKeyId.GetValue() != "" {
				storageArgs.FargateEphemeralStorageKmsKeyId = pulumi.String(storageConfiguration.FargateEphemeralStorageKmsKeyId.GetValue())
			}
			if storageConfiguration.KmsKeyId.GetValue() != "" {
				storageArgs.KmsKeyId = pulumi.String(storageConfiguration.KmsKeyId.GetValue())
			}
			configurationArgs.ManagedStorageConfiguration = storageArgs
		}

		args.Configuration = configurationArgs
	}

	// One environment-wide Service Connect namespace, overridable per
	// service -- set here so services join the mesh with zero wiring.
	if spec.ServiceConnectNamespaceArn != "" {
		args.ServiceConnectDefaults = &ecs.ClusterServiceConnectDefaultsArgs{
			Namespace: pulumi.String(spec.ServiceConnectNamespaceArn),
		}
	}

	// Stable Pulumi resource name; the cloud identity travels in the
	// Name argument, never in Pulumi auto-naming.
	createdCluster, err := ecs.NewCluster(ctx, "cluster", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ECS cluster")
	}

	return createdCluster, nil
}
