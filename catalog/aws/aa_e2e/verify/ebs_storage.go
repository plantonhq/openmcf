package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dlm"
	dlmtypes "github.com/aws/aws-sdk-go-v2/service/dlm/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// --- AwsEbsVolume ---------------------------------------------------------------

// ebsVolumeVerifier verifies an AwsEbsVolume via DescribeVolumes,
// keyed on volume_id. Deleting/deleted states count as absent (a
// deleting volume lingers describable briefly). The in-line
// attachments are children of the volume's lifecycle and cannot
// outlive the instance-volume pair.
type ebsVolumeVerifier struct{}

func (*ebsVolumeVerifier) IDOutputKey() string { return "volume_id" }

func (*ebsVolumeVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ebsVolumeExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsebsvolume verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsebsvolume %q not found after deploy", id)
	}
	return nil
}

func (*ebsVolumeVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ebsVolumeExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsebsvolume verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsebsvolume %q still exists after destroy", id)
	}
	return nil
}

func ebsVolumeExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeVolumes(ctx, &ec2.DescribeVolumesInput{
		VolumeIds: []string{id},
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidVolume.NotFound" {
			return false, nil
		}
		return false, err
	}
	for _, volume := range out.Volumes {
		switch volume.State {
		case ec2types.VolumeStateDeleting, ec2types.VolumeStateDeleted:
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// --- AwsEbsSnapshot -------------------------------------------------------------

// ebsSnapshotVerifier verifies an AwsEbsSnapshot via
// DescribeSnapshots, keyed on snapshot_id. The FSR and share
// satellites are children of the snapshot and die with it.
type ebsSnapshotVerifier struct{}

func (*ebsSnapshotVerifier) IDOutputKey() string { return "snapshot_id" }

func (*ebsSnapshotVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ebsSnapshotExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsebssnapshot verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsebssnapshot %q not found after deploy", id)
	}
	return nil
}

func (*ebsSnapshotVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := ebsSnapshotExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsebssnapshot verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsebssnapshot %q still exists after destroy", id)
	}
	return nil
}

func ebsSnapshotExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeSnapshots(ctx, &ec2.DescribeSnapshotsInput{
		SnapshotIds: []string{id},
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidSnapshot.NotFound" {
			return false, nil
		}
		return false, err
	}
	return len(out.Snapshots) > 0, nil
}

// --- AwsDlmLifecyclePolicy ------------------------------------------------------

// dlmLifecyclePolicyVerifier verifies an AwsDlmLifecyclePolicy via
// GetLifecyclePolicy, keyed on policy_id.
type dlmLifecyclePolicyVerifier struct{}

func (*dlmLifecyclePolicyVerifier) IDOutputKey() string { return "policy_id" }

func (*dlmLifecyclePolicyVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := dlmPolicyExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsdlmlifecyclepolicy verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsdlmlifecyclepolicy %q not found after deploy", id)
	}
	return nil
}

func (*dlmLifecyclePolicyVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := dlmPolicyExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsdlmlifecyclepolicy verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsdlmlifecyclepolicy %q still exists after destroy", id)
	}
	return nil
}

func dlmPolicyExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := dlm.NewFromConfig(cfg, func(o *dlm.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.GetLifecyclePolicy(ctx, &dlm.GetLifecyclePolicyInput{
		PolicyId: aws.String(id),
	})
	if err != nil {
		var notFound *dlmtypes.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return out.Policy != nil, nil
}
