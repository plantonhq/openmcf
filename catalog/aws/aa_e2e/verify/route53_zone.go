package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	pkgerrors "github.com/pkg/errors"
)

// route53ZoneVerifier verifies an AwsRoute53Zone via GetHostedZone, keyed on
// the zone_id output. Route 53 is a global service (no regional client
// override needed) and zone deletion is synchronous — a deleted zone returns
// the typed NoSuchHostedZone immediately.
type route53ZoneVerifier struct{}

func (*route53ZoneVerifier) IDOutputKey() string { return "zone_id" }

func (*route53ZoneVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := route53ZoneExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsroute53zone verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsroute53zone %q not found after deploy", id)
	}
	return nil
}

func (*route53ZoneVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := route53ZoneExists(ctx, cfg, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsroute53zone verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsroute53zone %q still exists after destroy", id)
	}
	return nil
}

func route53ZoneExists(ctx context.Context, cfg aws.Config, id string) (bool, error) {
	client := route53.NewFromConfig(cfg)
	_, err := client.GetHostedZone(ctx, &route53.GetHostedZoneInput{Id: aws.String(id)})
	if err != nil {
		var notFound *route53types.NoSuchHostedZone
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
