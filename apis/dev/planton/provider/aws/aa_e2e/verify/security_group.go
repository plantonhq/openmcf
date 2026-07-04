package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// securityGroupVerifier verifies an AwsSecurityGroup via DescribeSecurityGroups.
// It exists so an AwsSecurityGroup can be used as a deployed E2E prerequisite
// (e.g. for AwsMskCluster, whose brokers must attach at least one group at
// creation) and confirmed live. A deleted group returns the typed
// InvalidGroup.NotFound error (the "absent" signal).
type securityGroupVerifier struct{}

func (*securityGroupVerifier) IDOutputKey() string { return "security_group_id" }

func (*securityGroupVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := securityGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssecuritygroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awssecuritygroup %q not found after deploy", id)
	}
	return nil
}

func (*securityGroupVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := securityGroupExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awssecuritygroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awssecuritygroup %q still exists after destroy", id)
	}
	return nil
}

func securityGroupExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{GroupIds: []string{id}})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidGroup.NotFound" {
			return false, nil
		}
		return false, err
	}
	return len(out.SecurityGroups) > 0, nil
}
