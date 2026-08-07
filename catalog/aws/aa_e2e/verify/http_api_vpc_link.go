package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	"github.com/pkg/errors"
)

// httpApiVpcLinkVerifier verifies an AwsHttpApiVpcLink via GetVpcLink, keyed
// on the vpc_link_id output. VPC link deletion is asynchronous (AWS reclaims
// the managed ENIs): a link in DELETING status is treated as absent, the
// NAT-gateway lifecycle class.
type httpApiVpcLinkVerifier struct{}

func (*httpApiVpcLinkVerifier) IDOutputKey() string { return "vpc_link_id" }

func (v *httpApiVpcLinkVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	status, found, err := httpApiVpcLinkStatus(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awshttpapivpclink verify-exists failed for %q", id)
	}
	if !found || status == types.VpcLinkStatusDeleting {
		return errors.Errorf("awshttpapivpclink %q not found (or deleting) after deploy", id)
	}
	return nil
}

func (v *httpApiVpcLinkVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	status, found, err := httpApiVpcLinkStatus(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awshttpapivpclink verify-absent failed for %q", id)
	}
	if found && status != types.VpcLinkStatusDeleting {
		return errors.Errorf("awshttpapivpclink %q still exists after destroy", id)
	}
	return nil
}

func httpApiVpcLinkStatus(ctx context.Context, cfg aws.Config, vpcLinkId, region string) (types.VpcLinkStatus, bool, error) {
	client := apigatewayv2Client(cfg, region)
	out, err := client.GetVpcLink(ctx, &apigatewayv2.GetVpcLinkInput{VpcLinkId: aws.String(vpcLinkId)})
	if err != nil {
		var notFound *types.NotFoundException
		if errors.As(err, &notFound) {
			return "", false, nil
		}
		return "", false, err
	}
	return out.VpcLinkStatus, true, nil
}
