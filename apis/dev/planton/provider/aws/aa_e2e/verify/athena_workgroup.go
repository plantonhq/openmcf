package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	athenatypes "github.com/aws/aws-sdk-go-v2/service/athena/types"
	pkgerrors "github.com/pkg/errors"
)

// athenaWorkgroupVerifier verifies an AwsAthenaWorkgroup via GetWorkGroup.
//
// Athena reports a missing workgroup as an InvalidRequestException whose
// message contains "is not found" (no typed NotFound error exists in the
// service API) -- the same signal the Terraform provider's finder keys on.
// Deletion is synchronous, so no transitional deleting state needs handling.
type athenaWorkgroupVerifier struct{}

func (*athenaWorkgroupVerifier) IDOutputKey() string { return "workgroup_name" }

func (*athenaWorkgroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := athenaWorkgroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsathenaworkgroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsathenaworkgroup %q not found after deploy", id)
	}
	return nil
}

func (*athenaWorkgroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := athenaWorkgroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsathenaworkgroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsathenaworkgroup %q still exists after destroy", id)
	}
	return nil
}

func athenaWorkgroupExists(ctx context.Context, cfg aws.Config, name, region string) (bool, error) {
	client := athena.NewFromConfig(cfg, func(o *athena.Options) {
		if region != "" {
			o.Region = region
		}
	})
	_, err := client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{WorkGroup: &name})
	if err != nil {
		var invalidRequest *athenatypes.InvalidRequestException
		if errors.As(err, &invalidRequest) && strings.Contains(invalidRequest.ErrorMessage(), "is not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
