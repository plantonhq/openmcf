package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// clientVpnVerifier verifies an AwsClientVpn via DescribeClientVpnEndpoints
// keyed on the endpoint ID. A destroyed endpoint briefly reports the
// "deleting" state before the API returns the
// InvalidClientVpnEndpointId.NotFound error code (EC2 uses untyped smithy
// codes here) -- both mean "absent" for verification purposes.
type clientVpnVerifier struct{}

func (*clientVpnVerifier) IDOutputKey() string { return "client_vpn_endpoint_id" }

func (*clientVpnVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := clientVpnEndpointExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsclientvpn verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsclientvpn %q not found after deploy", id)
	}
	return nil
}

func (*clientVpnVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := clientVpnEndpointExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsclientvpn verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsclientvpn %q still exists after destroy", id)
	}
	return nil
}

func clientVpnEndpointExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeClientVpnEndpoints(ctx, &ec2.DescribeClientVpnEndpointsInput{
		ClientVpnEndpointIds: []string{id},
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && strings.HasPrefix(apiErr.ErrorCode(), "InvalidClientVpnEndpointId") {
			return false, nil
		}
		return false, err
	}
	for _, ep := range out.ClientVpnEndpoints {
		if aws.ToString(ep.ClientVpnEndpointId) != id {
			continue
		}
		if ep.Status != nil && ep.Status.Code == awstypes.ClientVpnEndpointStatusCodeDeleting {
			return false, nil
		}
		return true, nil
	}
	return false, nil
}
