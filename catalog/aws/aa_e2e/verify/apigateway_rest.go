package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/pkg/errors"
)

func apigatewayClient(cfg aws.Config, region string) *apigateway.Client {
	return apigateway.NewFromConfig(cfg, func(o *apigateway.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

func isAPIGatewayNotFound(err error) bool {
	var nf *types.NotFoundException
	return errors.As(err, &nf)
}

// restApiGatewayVerifier verifies an AwsRestApiGateway via GetRestApi,
// keyed on the rest_api_id output. REST API deletion is synchronous.
type restApiGatewayVerifier struct{}

func (*restApiGatewayVerifier) IDOutputKey() string { return "rest_api_id" }

func (v *restApiGatewayVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := restApiExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsrestapigateway verify-exists failed for %q", id)
	}
	if !exists {
		return errors.Errorf("awsrestapigateway %q not found after deploy", id)
	}
	return nil
}

func (v *restApiGatewayVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := restApiExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsrestapigateway verify-absent failed for %q", id)
	}
	if exists {
		return errors.Errorf("awsrestapigateway %q still exists after destroy", id)
	}
	return nil
}

func restApiExists(ctx context.Context, cfg aws.Config, restApiId, region string) (bool, error) {
	_, err := apigatewayClient(cfg, region).GetRestApi(ctx, &apigateway.GetRestApiInput{
		RestApiId: aws.String(restApiId),
	})
	if err != nil {
		if isAPIGatewayNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// restApiDomainVerifier verifies an AwsRestApiDomain via GetDomainName,
// keyed on the domain_name output (the hostname is the v1 identifier).
type restApiDomainVerifier struct{}

func (*restApiDomainVerifier) IDOutputKey() string { return "domain_name" }

func (v *restApiDomainVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := restApiDomainExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsrestapidomain verify-exists failed for %q", id)
	}
	if !exists {
		return errors.Errorf("awsrestapidomain %q not found after deploy", id)
	}
	return nil
}

func (v *restApiDomainVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := restApiDomainExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsrestapidomain verify-absent failed for %q", id)
	}
	if exists {
		return errors.Errorf("awsrestapidomain %q still exists after destroy", id)
	}
	return nil
}

func restApiDomainExists(ctx context.Context, cfg aws.Config, domainName, region string) (bool, error) {
	_, err := apigatewayClient(cfg, region).GetDomainName(ctx, &apigateway.GetDomainNameInput{
		DomainName: aws.String(domainName),
	})
	if err != nil {
		if isAPIGatewayNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// restApiUsagePlanVerifier verifies an AwsRestApiUsagePlan via
// GetUsagePlan, keyed on the usage_plan_id output. Attached API keys
// are folded satellites the lane's output map covers.
type restApiUsagePlanVerifier struct{}

func (*restApiUsagePlanVerifier) IDOutputKey() string { return "usage_plan_id" }

func (v *restApiUsagePlanVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := restApiUsagePlanExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsrestapiusageplan verify-exists failed for %q", id)
	}
	if !exists {
		return errors.Errorf("awsrestapiusageplan %q not found after deploy", id)
	}
	return nil
}

func (v *restApiUsagePlanVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := restApiUsagePlanExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsrestapiusageplan verify-absent failed for %q", id)
	}
	if exists {
		return errors.Errorf("awsrestapiusageplan %q still exists after destroy", id)
	}
	return nil
}

func restApiUsagePlanExists(ctx context.Context, cfg aws.Config, usagePlanId, region string) (bool, error) {
	_, err := apigatewayClient(cfg, region).GetUsagePlan(ctx, &apigateway.GetUsagePlanInput{
		UsagePlanId: aws.String(usagePlanId),
	})
	if err != nil {
		if isAPIGatewayNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// restApiVpcLinkVerifier verifies an AwsRestApiVpcLink via GetVpcLink,
// keyed on the vpc_link_id output. VPC link deletion is asynchronous
// (AWS reclaims the NLB attachment): a link in DELETING status is
// treated as absent, the NAT-gateway lifecycle class.
type restApiVpcLinkVerifier struct{}

func (*restApiVpcLinkVerifier) IDOutputKey() string { return "vpc_link_id" }

func (v *restApiVpcLinkVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	status, found, err := restApiVpcLinkStatus(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsrestapivpclink verify-exists failed for %q", id)
	}
	if !found || status == types.VpcLinkStatusDeleting {
		return errors.Errorf("awsrestapivpclink %q not found (or deleting) after deploy", id)
	}
	return nil
}

func (v *restApiVpcLinkVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	status, found, err := restApiVpcLinkStatus(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awsrestapivpclink verify-absent failed for %q", id)
	}
	if found && status != types.VpcLinkStatusDeleting {
		return errors.Errorf("awsrestapivpclink %q still exists after destroy", id)
	}
	return nil
}

func restApiVpcLinkStatus(ctx context.Context, cfg aws.Config, vpcLinkId, region string) (types.VpcLinkStatus, bool, error) {
	out, err := apigatewayClient(cfg, region).GetVpcLink(ctx, &apigateway.GetVpcLinkInput{
		VpcLinkId: aws.String(vpcLinkId),
	})
	if err != nil {
		if isAPIGatewayNotFound(err) {
			return "", false, nil
		}
		return "", false, err
	}
	return out.Status, true, nil
}
