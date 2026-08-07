package module

import (
	"github.com/pkg/errors"
	alicloudrocketmqinstancev1alpha1 "github.com/plantonhq/planton/catalog/alicloud/alicloudrocketmqinstance/v1alpha1"
	"github.com/pulumi/pulumi-alicloud/sdk/v3/go/alicloud"
	"github.com/pulumi/pulumi-alicloud/sdk/v3/go/alicloud/rocketmq"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func topic(
	ctx *pulumi.Context,
	provider *alicloud.Provider,
	instance *rocketmq.RocketMQInstance,
	t *alicloudrocketmqinstancev1alpha1.AliCloudRocketmqTopic,
) (*rocketmq.RocketMQTopic, error) {
	args := &rocketmq.RocketMQTopicArgs{
		InstanceId:  instance.ID(),
		TopicName:   pulumi.String(t.TopicName),
		MessageType: pulumi.StringPtr(messageType(t)),
	}

	if t.Remark != "" {
		args.Remark = pulumi.StringPtr(t.Remark)
	}

	args.MaxSendTps = optionalInt(t.MaxSendTps)

	created, err := rocketmq.NewRocketMQTopic(ctx, t.TopicName, args,
		pulumi.Provider(provider),
		pulumi.Parent(instance),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create topic %s", t.TopicName)
	}

	return created, nil
}
