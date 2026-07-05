package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// notification creates the event-notification satellite when configured. One
// notification resource carries all four arms. SQS/SNS/Lambda targets must
// already permit S3 delivery (queue/topic policy or Lambda invoke permission)
// or AWS rejects the PUT — the EventBridge arm needs no grant.
func notification(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec
	if spec.Notification == nil {
		return nil
	}

	args := &s3.BucketNotificationArgs{
		Bucket:      createdBucket.ID(),
		Eventbridge: pulumi.BoolPtr(spec.Notification.Eventbridge),
	}

	if len(spec.Notification.LambdaFunctions) > 0 {
		lambdas := s3.BucketNotificationLambdaFunctionArray{}
		for _, l := range spec.Notification.LambdaFunctions {
			lambda := &s3.BucketNotificationLambdaFunctionArgs{
				LambdaFunctionArn: pulumi.StringPtr(l.LambdaFunctionArn.GetValue()),
				Events:            pulumi.ToStringArray(l.Events),
			}
			if l.FilterPrefix != "" {
				lambda.FilterPrefix = pulumi.StringPtr(l.FilterPrefix)
			}
			if l.FilterSuffix != "" {
				lambda.FilterSuffix = pulumi.StringPtr(l.FilterSuffix)
			}
			lambdas = append(lambdas, lambda)
		}
		args.LambdaFunctions = lambdas
	}

	if len(spec.Notification.Queues) > 0 {
		queues := s3.BucketNotificationQueueArray{}
		for _, q := range spec.Notification.Queues {
			queue := &s3.BucketNotificationQueueArgs{
				QueueArn: pulumi.String(q.QueueArn.GetValue()),
				Events:   pulumi.ToStringArray(q.Events),
			}
			if q.FilterPrefix != "" {
				queue.FilterPrefix = pulumi.StringPtr(q.FilterPrefix)
			}
			if q.FilterSuffix != "" {
				queue.FilterSuffix = pulumi.StringPtr(q.FilterSuffix)
			}
			queues = append(queues, queue)
		}
		args.Queues = queues
	}

	if len(spec.Notification.Topics) > 0 {
		topics := s3.BucketNotificationTopicArray{}
		for _, t := range spec.Notification.Topics {
			topic := &s3.BucketNotificationTopicArgs{
				TopicArn: pulumi.String(t.TopicArn.GetValue()),
				Events:   pulumi.ToStringArray(t.Events),
			}
			if t.FilterPrefix != "" {
				topic.FilterPrefix = pulumi.StringPtr(t.FilterPrefix)
			}
			if t.FilterSuffix != "" {
				topic.FilterSuffix = pulumi.StringPtr(t.FilterSuffix)
			}
			topics = append(topics, topic)
		}
		args.Topics = topics
	}

	if _, err := s3.NewBucketNotification(ctx, "notification", args, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to configure notifications")
	}
	return nil
}

// objectLock creates the Object Lock default-retention satellite. The root
// resource's object_lock_enabled makes the bucket lock-capable; this adds the
// default retention window applied to every new object.
func objectLock(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec
	if spec.ObjectLockDefaultRetention == nil {
		return nil
	}

	retention := &s3.BucketObjectLockConfigurationV2RuleDefaultRetentionArgs{
		Mode: pulumi.StringPtr(spec.ObjectLockDefaultRetention.Mode),
	}
	if spec.ObjectLockDefaultRetention.Days > 0 {
		retention.Days = pulumi.IntPtr(int(spec.ObjectLockDefaultRetention.Days))
	}
	if spec.ObjectLockDefaultRetention.Years > 0 {
		retention.Years = pulumi.IntPtr(int(spec.ObjectLockDefaultRetention.Years))
	}

	if _, err := s3.NewBucketObjectLockConfigurationV2(ctx, "object-lock", &s3.BucketObjectLockConfigurationV2Args{
		Bucket: createdBucket.ID(),
		Rule: &s3.BucketObjectLockConfigurationV2RuleArgs{
			DefaultRetention: retention,
		},
	}, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to configure object lock retention")
	}
	return nil
}

// accelerate creates the Transfer Acceleration satellite when a state is set.
func accelerate(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec
	if spec.AccelerationStatus == "" {
		return nil
	}

	if _, err := s3.NewBucketAccelerateConfigurationV2(ctx, "accelerate", &s3.BucketAccelerateConfigurationV2Args{
		Bucket: createdBucket.ID(),
		Status: pulumi.String(spec.AccelerationStatus),
	}, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to configure transfer acceleration")
	}
	return nil
}

// requestPayment creates the Requester Pays satellite when a payer is set.
func requestPayment(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec
	if spec.RequestPayer == "" {
		return nil
	}

	if _, err := s3.NewBucketRequestPaymentConfigurationV2(ctx, "request-payment", &s3.BucketRequestPaymentConfigurationV2Args{
		Bucket: createdBucket.ID(),
		Payer:  pulumi.String(spec.RequestPayer),
	}, pulumi.Provider(provider)); err != nil {
		return errors.Wrap(err, "failed to configure requester pays")
	}
	return nil
}

// intelligentTiering creates one archive-configuration resource per named
// spec entry — a many-per-bucket satellite keyed by name, so adding or
// removing one never disturbs the others.
func intelligentTiering(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2) error {
	spec := locals.Spec

	for _, c := range spec.IntelligentTieringConfigurations {
		status := c.Status
		if status == "" {
			status = "Enabled"
		}

		tierings := s3.BucketIntelligentTieringConfigurationTieringArray{}
		for _, t := range c.Tiers {
			tierings = append(tierings, &s3.BucketIntelligentTieringConfigurationTieringArgs{
				AccessTier: pulumi.String(t.AccessTier),
				Days:       pulumi.Int(int(t.Days)),
			})
		}

		args := &s3.BucketIntelligentTieringConfigurationArgs{
			Bucket:   createdBucket.ID(),
			Name:     pulumi.String(c.Name),
			Status:   pulumi.StringPtr(status),
			Tierings: tierings,
		}
		if c.FilterPrefix != "" || len(c.FilterTags) > 0 {
			filter := &s3.BucketIntelligentTieringConfigurationFilterArgs{}
			if c.FilterPrefix != "" {
				filter.Prefix = pulumi.StringPtr(c.FilterPrefix)
			}
			if len(c.FilterTags) > 0 {
				filter.Tags = pulumi.ToStringMap(c.FilterTags)
			}
			args.Filter = filter
		}

		// The configuration name is the per-instance identity, mirroring the
		// Terraform module's for_each key.
		if _, err := s3.NewBucketIntelligentTieringConfiguration(ctx, "intelligent-tiering-"+c.Name, args, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "failed to configure intelligent tiering %q", c.Name)
		}
	}
	return nil
}
