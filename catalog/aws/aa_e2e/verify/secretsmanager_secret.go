package verify

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	pkgerrors "github.com/pkg/errors"
)

// secretsManagerSecretVerifier verifies an AwsSecretsManagerSecret via
// DescribeSecret, keyed on secret_arn. A secret with DeletedDate set is
// soft-deleted (scheduled for deletion within its recovery window) and
// counts as ABSENT -- destroy with a non-zero recovery window leaves the
// record describable until the window elapses.
//
// Replication is part of the contract, both directions (live-caught
// 2026-08-13). Exists-time: replication into a region already holding a
// same-name ex-replica FAILS SILENTLY at apply on both engines -- the
// provider has no replication waiter and force_overwrite_replica_secret
// does not clear the class ("Replication failed: Secret name ... is
// currently replicated to <region> with a different arn") -- so a green
// deploy can carry a dead replica claim unless the verifier polls every
// declared replica to InSync. Absent-time: AWS deletes replica secrets
// ASYNCHRONOUSLY after RemoveRegionsFromReplication and the provider does
// not wait for it before deleting the primary, so a replica can strand as
// a live standalone secret the primary-region check never sees. The
// verifier records each secret's replica regions at exists-time and
// asserts every one is gone (or soft-deleted) at absent-time.
type secretsManagerSecretVerifier struct{}

// replicaRegionsBySecretARN carries each verified secret's replica regions
// from VerifyExists to VerifyAbsent (both run in one test process). Keyed
// by the primary's ARN; the value is the region list snapshotted from
// DescribeSecret.ReplicationStatus.
var replicaRegionsBySecretARN sync.Map

const (
	// Replication of a fresh secret reaches InSync in seconds; the async
	// replica deletion after RemoveRegionsFromReplication lands in
	// ~40-90s (probed live 2026-08-13).
	replicationSyncTimeout  = 90 * time.Second
	replicaDeletionTimeout  = 3 * time.Minute
	replicationPollInterval = 5 * time.Second
)

func (*secretsManagerSecretVerifier) IDOutputKey() string { return "secret_arn" }

func (*secretsManagerSecretVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := secretsManagerClient(cfg, region)

	out, err := describeSecret(ctx, client, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssecretsmanagersecret verify-exists failed for %q", id)
	}
	if out == nil || out.DeletedDate != nil {
		return pkgerrors.Errorf("awssecretsmanagersecret %q not found after deploy", id)
	}

	// Poll every declared replica to InSync; a Failed entry is a real
	// defect the engines never surface (no replication waiter upstream).
	deadline := time.Now().Add(replicationSyncTimeout)
	for {
		regions := make([]string, 0, len(out.ReplicationStatus))
		var pending, failed []string
		for _, rs := range out.ReplicationStatus {
			regions = append(regions, aws.ToString(rs.Region))
			switch rs.Status {
			case types.StatusTypeInSync:
			case types.StatusTypeFailed:
				failed = append(failed, aws.ToString(rs.Region)+": "+aws.ToString(rs.StatusMessage))
			default:
				pending = append(pending, aws.ToString(rs.Region))
			}
		}
		if len(failed) > 0 {
			return pkgerrors.Errorf("awssecretsmanagersecret %q replication FAILED (%s)", id, strings.Join(failed, "; "))
		}
		if len(pending) == 0 {
			replicaRegionsBySecretARN.Store(id, regions)
			return nil
		}
		if time.Now().After(deadline) {
			return pkgerrors.Errorf("awssecretsmanagersecret %q replication not InSync in %s (pending: %s)", id, replicationSyncTimeout, strings.Join(pending, ", "))
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(replicationPollInterval):
		}
		if out, err = describeSecret(ctx, client, id); err != nil {
			return pkgerrors.Wrapf(err, "awssecretsmanagersecret verify-exists failed for %q", id)
		}
		if out == nil || out.DeletedDate != nil {
			return pkgerrors.Errorf("awssecretsmanagersecret %q disappeared while awaiting replication", id)
		}
	}
}

func (*secretsManagerSecretVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := secretsManagerSecretExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssecretsmanagersecret verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssecretsmanagersecret %q still exists after destroy", id)
	}

	// Replica regions recorded at exists-time: each replica must be gone
	// (or soft-deleted) too. AWS deletes them asynchronously after
	// RemoveRegionsFromReplication, and a strand outlives the lane as a
	// live standalone secret -- the primary-region check cannot see it.
	regionsVal, ok := replicaRegionsBySecretARN.LoadAndDelete(id)
	if !ok {
		return nil
	}
	name := secretNameFromARN(id)
	deadline := time.Now().Add(replicaDeletionTimeout)
	for _, replicaRegion := range regionsVal.([]string) {
		for {
			exists, err := secretsManagerSecretExists(ctx, cfg, name, replicaRegion)
			if err != nil {
				return pkgerrors.Wrapf(err, "awssecretsmanagersecret %q replica verify-absent failed in %s", name, replicaRegion)
			}
			if !exists {
				break
			}
			if time.Now().After(deadline) {
				return pkgerrors.Errorf("awssecretsmanagersecret %q replica STRANDED in %s after destroy (recover: stop-replication-to-replica in %s, then delete-secret)", name, replicaRegion, replicaRegion)
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(replicationPollInterval):
			}
		}
	}
	return nil
}

func secretsManagerClient(cfg aws.Config, region string) *secretsmanager.Client {
	return secretsmanager.NewFromConfig(cfg, func(o *secretsmanager.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// describeSecret returns nil (no error) when the secret does not exist.
func describeSecret(ctx context.Context, client *secretsmanager.Client, secretID string) (*secretsmanager.DescribeSecretOutput, error) {
	out, err := client.DescribeSecret(ctx, &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretID),
	})
	if err == nil {
		return out, nil
	}
	var notFound *types.ResourceNotFoundException
	if errors.As(err, &notFound) {
		return nil, nil
	}
	return nil, err
}

// secretNameFromARN extracts the friendly name from a secret ARN
// (arn:aws:secretsmanager:<region>:<acct>:secret:<name>-<6-char-suffix>).
// Replica secrets share the primary's suffix but live in another region,
// so the NAME is the cross-region lookup key.
func secretNameFromARN(arn string) string {
	parts := strings.Split(arn, ":")
	last := parts[len(parts)-1]
	if i := strings.LastIndex(last, "-"); i > 0 {
		return last[:i]
	}
	return last
}

func secretsManagerSecretExists(ctx context.Context, cfg aws.Config, secretARN, region string) (bool, error) {
	out, err := describeSecret(ctx, secretsManagerClient(cfg, region), secretARN)
	if err != nil {
		return false, err
	}
	if out == nil {
		return false, nil
	}
	// Soft-deleted secrets stay describable through the recovery
	// window; DeletedDate set means the destroy already happened.
	return out.DeletedDate == nil, nil
}
