package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	pkgerrors "github.com/pkg/errors"
)

// rdsInstanceVerifier verifies an AwsRdsInstance via DescribeDBInstances,
// keyed on the instance identifier output. Like the cluster, an instance
// mid-deletion stays describable with a "deleting" status before RDS starts
// returning the typed DBInstanceNotFoundFault, so both signals count as
// absent.
type rdsInstanceVerifier struct{}

func (*rdsInstanceVerifier) IDOutputKey() string { return "instance_identifier" }

func (*rdsInstanceVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := rdsInstanceExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsrdsinstance verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsrdsinstance %q not found after deploy", id)
	}
	return nil
}

func (*rdsInstanceVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := rdsInstanceExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsrdsinstance verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsrdsinstance %q still exists after destroy", id)
	}
	return nil
}

// rdsInstanceExists reports whether the instance is present and not already
// on its way out. A DBInstanceNotFoundFault is treated as absent.
func rdsInstanceExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := rds.NewFromConfig(cfg, func(o *rds.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeDBInstances(ctx, &rds.DescribeDBInstancesInput{DBInstanceIdentifier: &id})
	if err != nil {
		var notFound *types.DBInstanceNotFoundFault
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	for _, instance := range out.DBInstances {
		if instance.DBInstanceStatus == nil {
			continue
		}
		switch strings.ToLower(*instance.DBInstanceStatus) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}
