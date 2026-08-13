package module

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/msk"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// satellites creates the cluster-scoped settings AWS models as standalone
// resources but that are honestly part of the cluster's own configuration --
// each is keyed by the cluster ARN, owned by exactly one cluster, and
// referenced by nothing else.
func satellites(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdCluster *msk.Cluster) error {
	spec := locals.AwsMskCluster.Spec

	// SASL/SCRAM credential associations: one association per Secrets Manager
	// secret, materialized per-ARN (keyed by the secret's own name, never a
	// list index) so adding or removing a secret updates in place instead of
	// churning its neighbors.
	for _, secretArn := range spec.ScramSecretArns {
		secretName := secretArn
		if idx := strings.LastIndex(secretArn, ":secret:"); idx >= 0 {
			secretName = secretArn[idx+len(":secret:"):]
		}
		_, err := msk.NewSingleScramSecretAssociation(ctx,
			fmt.Sprintf("scram-secret-%s", secretName),
			&msk.SingleScramSecretAssociationArgs{
				ClusterArn: createdCluster.Arn,
				SecretArn:  pulumi.String(secretArn),
			}, pulumi.Provider(provider), pulumi.Parent(createdCluster))
		if err != nil {
			return errors.Wrapf(err, "associate scram secret %s", secretArn)
		}
	}

	// A resource-based IAM policy on the cluster -- the grant behind
	// cross-account PrivateLink access (kafka:CreateVpcConnection and friends).
	// The spec carries the policy as a structured document; it is serialized
	// to JSON here.
	if spec.ClusterPolicy != nil {
		policyJson, err := json.Marshal(spec.ClusterPolicy.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal cluster policy")
		}
		_, err = msk.NewClusterPolicy(ctx, "cluster-policy", &msk.ClusterPolicyArgs{
			ClusterArn: createdCluster.Arn,
			Policy:     pulumi.String(policyJson),
		}, pulumi.Provider(provider), pulumi.Parent(createdCluster))
		if err != nil {
			return errors.Wrap(err, "attach cluster policy")
		}
	}

	// Declared Kafka topics, managed through the MSK topic API -- no Kafka
	// client or bootstrap connectivity needed. Keyed by topic name so adding
	// or removing one topic never churns its neighbors. Topic deletion
	// requires delete.topic.enable=true on the cluster (MSK's default).
	topicArns := pulumi.StringMap{}
	for _, topic := range spec.Topics {
		args := &msk.TopicArgs{
			ClusterArn:        createdCluster.Arn,
			Name:              pulumi.String(topic.Name),
			PartitionCount:    pulumi.Int(topic.PartitionCount),
			ReplicationFactor: pulumi.Int(topic.ReplicationFactor),
		}
		// The provider takes topic configs as a JSON document; entries with
		// no overrides omit the argument entirely.
		if len(topic.Configs) > 0 {
			configsJson, err := json.Marshal(topic.Configs)
			if err != nil {
				return errors.Wrapf(err, "marshal configs for topic %s", topic.Name)
			}
			args.Configs = pulumi.String(configsJson)
		}
		createdTopic, err := msk.NewTopic(ctx,
			fmt.Sprintf("topic-%s", topic.Name),
			args, pulumi.Provider(provider), pulumi.Parent(createdCluster))
		if err != nil {
			return errors.Wrapf(err, "create topic %s", topic.Name)
		}
		topicArns[topic.Name] = createdTopic.Arn
	}
	ctx.Export(OpTopicArns, topicArns)

	return nil
}
