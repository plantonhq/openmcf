package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cognito"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// logDelivery wires the pool's event sources to their log destinations. AWS
// models this as ONE pool-scoped configuration carrying every route, so all
// spec entries materialize into a single resource (two resources would fight
// over the same setting on every apply).
func logDelivery(ctx *pulumi.Context, locals *Locals, createdPool *cognito.UserPool, provider *aws.Provider) error {
	if len(locals.Spec.LogConfigurations) == 0 {
		return nil
	}

	var configurations cognito.LogDeliveryConfigurationLogConfigurationArray
	for _, lc := range locals.Spec.LogConfigurations {
		entry := &cognito.LogDeliveryConfigurationLogConfigurationArgs{
			EventSource: pulumi.String(lc.EventSource),
			LogLevel:    pulumi.String(lc.LogLevel),
		}
		// The spec's CEL guarantees exactly one destination per entry.
		if lc.CloudwatchLogGroupArn.GetValue() != "" {
			entry.CloudWatchLogsConfiguration = &cognito.LogDeliveryConfigurationLogConfigurationCloudWatchLogsConfigurationArgs{
				LogGroupArn: pulumi.StringPtr(lc.CloudwatchLogGroupArn.GetValue()),
			}
		}
		if lc.FirehoseStreamArn.GetValue() != "" {
			entry.FirehoseConfiguration = &cognito.LogDeliveryConfigurationLogConfigurationFirehoseConfigurationArgs{
				StreamArn: pulumi.StringPtr(lc.FirehoseStreamArn.GetValue()),
			}
		}
		if lc.S3BucketArn.GetValue() != "" {
			entry.S3Configuration = &cognito.LogDeliveryConfigurationLogConfigurationS3ConfigurationArgs{
				BucketArn: pulumi.StringPtr(lc.S3BucketArn.GetValue()),
			}
		}
		configurations = append(configurations, entry)
	}

	_, err := cognito.NewLogDeliveryConfiguration(ctx,
		locals.Target.Metadata.Name+"-log-delivery",
		&cognito.LogDeliveryConfigurationArgs{
			UserPoolId:        createdPool.ID(),
			LogConfigurations: configurations,
		}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create Cognito log delivery configuration")
	}

	return nil
}
