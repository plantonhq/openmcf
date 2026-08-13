package module

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// satellites provisions the table-scoped settings AWS models as
// standalone resources but that are honestly part of the table's own
// configuration -- each is keyed by the table (replace-on-change), owned
// by exactly one table, and referenced by nothing else: the resource
// policy, the Kinesis change-data destination, and CloudWatch
// contributor insights.
func satellites(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdTable *dynamodb.Table) error {
	spec := locals.AwsDynamodb.Spec

	// A resource-based IAM policy on the table -- cross-account access
	// grants without role assumption. The confirm flag is AWS's guard
	// against locking yourself out: a policy that removes the applying
	// caller's own access is refused unless it is set.
	if spec.ResourcePolicy != nil {
		policyJSON, err := json.Marshal(spec.ResourcePolicy.Policy.AsMap())
		if err != nil {
			return errors.Wrap(err, "failed to serialize resource policy to JSON")
		}
		if _, err := dynamodb.NewResourcePolicy(ctx, "resource-policy", &dynamodb.ResourcePolicyArgs{
			ResourceArn:                     createdTable.Arn,
			Policy:                          pulumi.String(string(policyJSON)),
			ConfirmRemoveSelfResourceAccess: pulumi.BoolPtr(spec.ResourcePolicy.ConfirmRemoveSelfResourceAccess),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "failed to attach DynamoDB resource policy")
		}
	}

	// Item-level change data into a Kinesis Data Stream (independent of
	// DynamoDB Streams). AWS allows exactly one destination per table.
	if spec.KinesisStreamingDestination != nil {
		destinationArgs := &dynamodb.KinesisStreamingDestinationArgs{
			TableName: createdTable.Name,
			StreamArn: pulumi.String(spec.KinesisStreamingDestination.StreamArn.GetValue()),
		}
		if spec.KinesisStreamingDestination.ApproximateCreationDateTimePrecision != "" {
			destinationArgs.ApproximateCreationDateTimePrecision = pulumi.String(spec.KinesisStreamingDestination.ApproximateCreationDateTimePrecision)
		}
		if _, err := dynamodb.NewKinesisStreamingDestination(ctx, "kinesis-streaming-destination",
			destinationArgs, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "failed to create DynamoDB Kinesis streaming destination")
		}
	}

	// CloudWatch contributor insights: one provider resource for the
	// table, plus one per opted-in GSI -- materialized per-name so an
	// index list edit updates in place.
	if spec.ContributorInsights != nil && spec.ContributorInsights.Enabled {
		insightsArgs := &dynamodb.ContributorInsightsArgs{TableName: createdTable.Name}
		if spec.ContributorInsights.Mode != "" {
			insightsArgs.Mode = pulumi.String(spec.ContributorInsights.Mode)
		}
		if _, err := dynamodb.NewContributorInsights(ctx, "contributor-insights",
			insightsArgs, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "failed to enable DynamoDB contributor insights")
		}

		for _, indexName := range spec.ContributorInsights.GsiIndexNames {
			indexInsightsArgs := &dynamodb.ContributorInsightsArgs{
				TableName: createdTable.Name,
				IndexName: pulumi.String(indexName),
			}
			if spec.ContributorInsights.Mode != "" {
				indexInsightsArgs.Mode = pulumi.String(spec.ContributorInsights.Mode)
			}
			if _, err := dynamodb.NewContributorInsights(ctx,
				fmt.Sprintf("contributor-insights-%s", indexName),
				indexInsightsArgs, pulumi.Provider(provider)); err != nil {
				return errors.Wrapf(err, "failed to enable contributor insights on index %s", indexName)
			}
		}
	}

	return nil
}
