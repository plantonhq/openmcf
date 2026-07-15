package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	pkgerrors "github.com/pkg/errors"
)

// The EFS family verifies through the elasticfilesystem control plane. Both
// kinds are lifecycle-state-aware (the NAT-gateway class): a deleted resource
// can briefly stay describable in a "deleting"/"deleted" state before the
// typed NotFound appears, so those states count as absent rather than racing
// the API's cleanup.

// efsFileSystemVerifier verifies an AwsElasticFileSystem via
// DescribeFileSystems, keyed on the file_system_id output.
type efsFileSystemVerifier struct{}

func (*efsFileSystemVerifier) IDOutputKey() string { return "file_system_id" }

func (*efsFileSystemVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := efsFileSystemExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awselasticfilesystem verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awselasticfilesystem %q not found after deploy", id)
	}
	return nil
}

func (*efsFileSystemVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := efsFileSystemExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awselasticfilesystem verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awselasticfilesystem %q still exists after destroy", id)
	}
	return nil
}

func efsFileSystemExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := efsClient(cfg, region)
	out, err := client.DescribeFileSystems(ctx, &efs.DescribeFileSystemsInput{
		FileSystemId: aws.String(id),
	})
	if err != nil {
		var notFound *efstypes.FileSystemNotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, fs := range out.FileSystems {
		if !efsLifecycleGone(string(fs.LifeCycleState)) {
			return true, nil
		}
	}
	return false, nil
}

// efsAccessPointVerifier verifies an AwsEfsAccessPoint via
// DescribeAccessPoints, keyed on the access_point_id output.
type efsAccessPointVerifier struct{}

func (*efsAccessPointVerifier) IDOutputKey() string { return "access_point_id" }

func (*efsAccessPointVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := efsAccessPointExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsefsaccesspoint verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsefsaccesspoint %q not found after deploy", id)
	}
	return nil
}

func (*efsAccessPointVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := efsAccessPointExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsefsaccesspoint verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsefsaccesspoint %q still exists after destroy", id)
	}
	return nil
}

func efsAccessPointExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := efsClient(cfg, region)
	out, err := client.DescribeAccessPoints(ctx, &efs.DescribeAccessPointsInput{
		AccessPointId: aws.String(id),
	})
	if err != nil {
		var notFound *efstypes.AccessPointNotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, ap := range out.AccessPoints {
		if !efsLifecycleGone(string(ap.LifeCycleState)) {
			return true, nil
		}
	}
	return false, nil
}

func efsClient(cfg aws.Config, region string) *efs.Client {
	return efs.NewFromConfig(cfg, func(o *efs.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// efsLifecycleGone treats the teardown-side lifecycle states as absent
// (case-insensitive; the API's enum values are lowercase).
func efsLifecycleGone(state string) bool {
	switch strings.ToLower(state) {
	case "deleting", "deleted":
		return true
	default:
		return false
	}
}
