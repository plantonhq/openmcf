package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

// vpcEndpointVerifier verifies an AwsVpcEndpoint via DescribeVpcEndpoints,
// keyed on the endpoint ID. AWS does not drop the endpoint record immediately
// on delete: a destroyed endpoint lingers describable in a "deleting"/"deleted"
// state for a while and only later returns the typed
// InvalidVpcEndpointId.NotFound error -- the NAT-gateway lifecycle class. Both
// signals mean "absent" for verification purposes. State comparison is
// case-insensitive because the EC2 API reports endpoint states in lowercase
// ("available", "deleted") while the SDK enum spells them capitalized.
type vpcEndpointVerifier struct{}

func (*vpcEndpointVerifier) IDOutputKey() string { return "vpc_endpoint_id" }

func (*vpcEndpointVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := vpcEndpointExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsvpcendpoint verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("awsvpcendpoint %q not found after deploy", id)
	}
	return nil
}

func (*vpcEndpointVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := vpcEndpointExists(ctx, cfg, id, region)
	if err != nil {
		return pkgerrors.Wrapf(err, "awsvpcendpoint verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("awsvpcendpoint %q still exists after destroy", id)
	}
	return nil
}

// vpcEndpointExists reports whether the endpoint is present and not in a
// deleting/deleted state. An InvalidVpcEndpointId.NotFound error is treated as
// absent.
func vpcEndpointExists(ctx context.Context, cfg aws.Config, id, region string) (bool, error) {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) {
		if region != "" {
			o.Region = region
		}
	})
	out, err := client.DescribeVpcEndpoints(ctx, &ec2.DescribeVpcEndpointsInput{VpcEndpointIds: []string{id}})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "InvalidVpcEndpointId.NotFound" {
			return false, nil
		}
		return false, err
	}
	for _, endpoint := range out.VpcEndpoints {
		switch strings.ToLower(string(endpoint.State)) {
		case "deleting", "deleted":
			continue
		default:
			return true, nil
		}
	}
	return false, nil
}
