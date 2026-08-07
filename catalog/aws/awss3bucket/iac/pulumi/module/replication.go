package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// replication creates the replication configuration satellite when the spec
// defines one.
func replication(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	createdBucket *s3.BucketV2, versioning *s3.BucketVersioningV2) error {
	spec := locals.Spec
	if spec.Replication == nil {
		return nil
	}

	rules := s3.BucketReplicationConfigRuleArray{}
	for _, r := range spec.Replication.Rules {
		status := r.Status
		if status == "" {
			status = "Enabled"
		}

		rule := &s3.BucketReplicationConfigRuleArgs{
			Id:       pulumi.StringPtr(r.Id),
			Priority: pulumi.IntPtr(int(r.Priority)),
			Status:   pulumi.String(status),
			// AWS requires the delete-marker block on the modern (V2)
			// replication schema; the spec's bool maps to Enabled/Disabled.
			DeleteMarkerReplication: &s3.BucketReplicationConfigRuleDeleteMarkerReplicationArgs{
				Status: pulumi.String(map[bool]string{true: "Enabled", false: "Disabled"}[r.DeleteMarkerReplication]),
			},
		}

		// Same single-vs-and document shaping as lifecycle filters: one
		// predicate goes directly on the filter, several go inside `and`.
		if r.Filter != nil {
			predicates := 0
			if r.Filter.Prefix != "" {
				predicates++
			}
			predicates += len(r.Filter.Tags)

			filter := &s3.BucketReplicationConfigRuleFilterArgs{}
			switch {
			case predicates > 1:
				and := &s3.BucketReplicationConfigRuleFilterAndArgs{}
				if r.Filter.Prefix != "" {
					and.Prefix = pulumi.StringPtr(r.Filter.Prefix)
				}
				if len(r.Filter.Tags) > 0 {
					and.Tags = pulumi.ToStringMap(r.Filter.Tags)
				}
				filter.And = and
			case r.Filter.Prefix != "":
				filter.Prefix = pulumi.StringPtr(r.Filter.Prefix)
			case len(r.Filter.Tags) == 1:
				for k, v := range r.Filter.Tags {
					filter.Tag = &s3.BucketReplicationConfigRuleFilterTagArgs{
						Key:   pulumi.String(k),
						Value: pulumi.String(v),
					}
				}
			}
			rule.Filter = filter
		}

		if r.ExistingObjectReplication {
			rule.ExistingObjectReplication = &s3.BucketReplicationConfigRuleExistingObjectReplicationArgs{
				Status: pulumi.String("Enabled"),
			}
		}

		if r.ReplicateReplicaModifications || r.ReplicateSseKmsEncryptedObjects {
			criteria := &s3.BucketReplicationConfigRuleSourceSelectionCriteriaArgs{}
			if r.ReplicateReplicaModifications {
				criteria.ReplicaModifications = &s3.BucketReplicationConfigRuleSourceSelectionCriteriaReplicaModificationsArgs{
					Status: pulumi.String("Enabled"),
				}
			}
			if r.ReplicateSseKmsEncryptedObjects {
				criteria.SseKmsEncryptedObjects = &s3.BucketReplicationConfigRuleSourceSelectionCriteriaSseKmsEncryptedObjectsArgs{
					Status: pulumi.String("Enabled"),
				}
			}
			rule.SourceSelectionCriteria = criteria
		}

		destination := &s3.BucketReplicationConfigRuleDestinationArgs{
			Bucket: pulumi.String(r.Destination.BucketArn.GetValue()),
		}
		if r.Destination.Account != "" {
			destination.Account = pulumi.StringPtr(r.Destination.Account)
		}
		if r.Destination.StorageClass != "" {
			destination.StorageClass = pulumi.StringPtr(r.Destination.StorageClass)
		}
		if r.Destination.ChangeReplicaOwnershipToDestination {
			destination.AccessControlTranslation = &s3.BucketReplicationConfigRuleDestinationAccessControlTranslationArgs{
				Owner: pulumi.String("Destination"),
			}
		}
		if r.Destination.ReplicaKmsKeyId.GetValue() != "" {
			destination.EncryptionConfiguration = &s3.BucketReplicationConfigRuleDestinationEncryptionConfigurationArgs{
				ReplicaKmsKeyId: pulumi.String(r.Destination.ReplicaKmsKeyId.GetValue()),
			}
		}
		// AWS requires metrics whenever RTC is enabled (CEL mirrors this);
		// both use the fixed 15-minute threshold AWS accepts.
		if r.Destination.MetricsEnabled {
			destination.Metrics = &s3.BucketReplicationConfigRuleDestinationMetricsArgs{
				Status: pulumi.String("Enabled"),
				EventThreshold: &s3.BucketReplicationConfigRuleDestinationMetricsEventThresholdArgs{
					Minutes: pulumi.Int(15),
				},
			}
		}
		if r.Destination.ReplicationTimeControlEnabled {
			destination.ReplicationTime = &s3.BucketReplicationConfigRuleDestinationReplicationTimeArgs{
				Status: pulumi.String("Enabled"),
				Time: &s3.BucketReplicationConfigRuleDestinationReplicationTimeTimeArgs{
					Minutes: pulumi.Int(15),
				},
			}
		}
		rule.Destination = destination

		rules = append(rules, rule)
	}

	// AWS rejects replication configuration until versioning is Enabled on
	// the source bucket, and the two PUTs race without an explicit edge.
	// (CEL guarantees versioning_status is Enabled whenever replication is
	// set, so the satellite always exists here.)
	opts := []pulumi.ResourceOption{pulumi.Provider(provider)}
	if versioning != nil {
		opts = append(opts, pulumi.DependsOn([]pulumi.Resource{versioning}))
	}
	if _, err := s3.NewBucketReplicationConfig(ctx, "replication", &s3.BucketReplicationConfigArgs{
		Bucket: createdBucket.ID(),
		Role:   pulumi.String(spec.Replication.RoleArn.GetValue()),
		Rules:  rules,
	}, opts...); err != nil {
		return errors.Wrap(err, "failed to configure replication")
	}
	return nil
}
