package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	pkgerrors "github.com/pkg/errors"
)

// kmsKeyVerifier verifies an AwsKmsKey via DescribeKey, keyed on the key_id
// output. KMS keys are never deleted immediately: destroy schedules deletion
// and the key sits in PendingDeletion for the 7-30 day recovery window, still
// fully describable. So existence is "described AND not pending deletion",
// and absence accepts either the (eventual) NotFoundException or the
// PendingDeletion/PendingReplicaDeletion states -- otherwise verify-absent
// could never pass within a test run's lifetime.
type kmsKeyVerifier struct{}

func (*kmsKeyVerifier) IDOutputKey() string { return "key_id" }

func (*kmsKeyVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := kmsKeyExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskmskey verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awskmskey %q not found after deploy", id)
	}
	return nil
}

func (*kmsKeyVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := kmsKeyExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awskmskey verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awskmskey %q still exists after destroy", id)
	}
	return nil
}

func kmsKeyExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := kms.NewFromConfig(cfg, func(o *kms.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeKey(ctx, &kms.DescribeKeyInput{KeyId: &id})
	if err != nil {
		var notFound *types.NotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.KeyMetadata == nil {
		return false, nil
	}
	// A key scheduled for deletion is destroyed from the harness's point of
	// view: the IaC destroy succeeded and only AWS's mandatory recovery window
	// keeps the key describable.
	switch out.KeyMetadata.KeyState {
	case types.KeyStatePendingDeletion, types.KeyStatePendingReplicaDeletion:
		return false, nil
	}
	return true, nil
}
