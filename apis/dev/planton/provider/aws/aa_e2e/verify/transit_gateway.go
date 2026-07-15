package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// The Transit Gateway family shares the NAT-gateway lifecycle class: AWS keeps
// deleted records describable in a "deleting"/"deleted" state for a while
// before the typed NotFound error appears. Both mean "absent" for
// verification purposes, so every verifier here checks the state, not just
// presence.

// transitGatewayVerifier verifies an AwsTransitGateway via
// DescribeTransitGateways.
type transitGatewayVerifier struct{}

func (*transitGatewayVerifier) IDOutputKey() string { return "transit_gateway_id" }

func (*transitGatewayVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := transitGatewayExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awstransitgateway verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awstransitgateway %q not found after deploy", id)
	}
	return nil
}

func (*transitGatewayVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := transitGatewayExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awstransitgateway verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awstransitgateway %q still exists after destroy", id)
	}
	return nil
}

func transitGatewayExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2ClientFor(cfg, region)
	out, err := client.DescribeTransitGateways(ctx, &ec2.DescribeTransitGatewaysInput{TransitGatewayIds: []string{id}})
	if err != nil {
		if isEc2NotFound(err, "InvalidTransitGatewayID.NotFound") {
			return false, nil
		}
		return false, err
	}
	for _, tgw := range out.TransitGateways {
		switch tgw.State {
		case awstypes.TransitGatewayStateDeleting, awstypes.TransitGatewayStateDeleted:
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// transitGatewayVpcAttachmentVerifier verifies an
// AwsTransitGatewayVpcAttachment via DescribeTransitGatewayVpcAttachments.
type transitGatewayVpcAttachmentVerifier struct{}

func (*transitGatewayVpcAttachmentVerifier) IDOutputKey() string { return "attachment_id" }

func (*transitGatewayVpcAttachmentVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := transitGatewayVpcAttachmentExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awstransitgatewayvpcattachment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awstransitgatewayvpcattachment %q not found after deploy", id)
	}
	return nil
}

func (*transitGatewayVpcAttachmentVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := transitGatewayVpcAttachmentExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awstransitgatewayvpcattachment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awstransitgatewayvpcattachment %q still exists after destroy", id)
	}
	return nil
}

func transitGatewayVpcAttachmentExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2ClientFor(cfg, region)
	out, err := client.DescribeTransitGatewayVpcAttachments(ctx, &ec2.DescribeTransitGatewayVpcAttachmentsInput{TransitGatewayAttachmentIds: []string{id}})
	if err != nil {
		if isEc2NotFound(err, "InvalidTransitGatewayAttachmentID.NotFound") {
			return false, nil
		}
		return false, err
	}
	for _, attachment := range out.TransitGatewayVpcAttachments {
		switch attachment.State {
		case awstypes.TransitGatewayAttachmentStateDeleting, awstypes.TransitGatewayAttachmentStateDeleted:
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// transitGatewayRouteTableVerifier verifies an AwsTransitGatewayRouteTable via
// DescribeTransitGatewayRouteTables.
type transitGatewayRouteTableVerifier struct{}

func (*transitGatewayRouteTableVerifier) IDOutputKey() string { return "route_table_id" }

func (*transitGatewayRouteTableVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := transitGatewayRouteTableExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awstransitgatewayroutetable verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awstransitgatewayroutetable %q not found after deploy", id)
	}
	return nil
}

func (*transitGatewayRouteTableVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := transitGatewayRouteTableExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awstransitgatewayroutetable verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awstransitgatewayroutetable %q still exists after destroy", id)
	}
	return nil
}

func transitGatewayRouteTableExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2ClientFor(cfg, region)
	out, err := client.DescribeTransitGatewayRouteTables(ctx, &ec2.DescribeTransitGatewayRouteTablesInput{TransitGatewayRouteTableIds: []string{id}})
	if err != nil {
		// AWS reports a missing TGW route table with the generic route-table
		// code, mirroring the provider's finder.
		if isEc2NotFound(err, "InvalidRouteTableID.NotFound") {
			return false, nil
		}
		return false, err
	}
	for _, routeTable := range out.TransitGatewayRouteTables {
		switch routeTable.State {
		case awstypes.TransitGatewayRouteTableStateDeleting, awstypes.TransitGatewayRouteTableStateDeleted:
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}

// ec2ClientFor builds a region-pinned EC2 client.
func ec2ClientFor(cfg aws.Config, region string) *ec2.Client {
	return ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

// isEc2NotFound reports whether err is the given typed EC2 NotFound code.
func isEc2NotFound(err error, code string) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && apiErr.ErrorCode() == code
}
