package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless"
	"github.com/aws/aws-sdk-go-v2/service/redshiftserverless/types"
	pkgerrors "github.com/pkg/errors"
)

// redshiftServerlessWorkgroupVerifier verifies an
// AwsRedshiftServerlessWorkgroup via GetWorkgroup, keyed on the
// workgroup_name output. A workgroup mid-deletion stays describable with
// a DELETING status before the service starts returning the typed
// ResourceNotFoundException -- the same lifecycle class as the RDS-shaped
// kinds -- so existence is "described AND not deleting", and absence
// accepts either signal.
type redshiftServerlessWorkgroupVerifier struct{}

func (*redshiftServerlessWorkgroupVerifier) IDOutputKey() string { return "workgroup_name" }

func (*redshiftServerlessWorkgroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := redshiftServerlessWorkgroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsredshiftserverlessworkgroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsredshiftserverlessworkgroup %q not found after deploy", id)
	}
	return nil
}

func (*redshiftServerlessWorkgroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := redshiftServerlessWorkgroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsredshiftserverlessworkgroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsredshiftserverlessworkgroup %q still exists after destroy", id)
	}
	return nil
}

// redshiftServerlessWorkgroupExists reports whether the workgroup is
// present and not already on its way out. A ResourceNotFoundException is
// treated as absent.
func redshiftServerlessWorkgroupExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := redshiftserverless.NewFromConfig(cfg, func(o *redshiftserverless.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.GetWorkgroup(ctx, &redshiftserverless.GetWorkgroupInput{WorkgroupName: &id})
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	if out.Workgroup == nil {
		return false, nil
	}
	if out.Workgroup.Status == types.WorkgroupStatusDeleting {
		return false, nil
	}
	return true, nil
}
