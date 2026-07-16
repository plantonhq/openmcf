package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	"github.com/pkg/errors"
)

// httpApiDomainVerifier verifies an AwsHttpApiDomain via GetDomainName, keyed
// on the domain_name output (the domain name IS the resource identifier in
// the API Gateway v2 API). Domain deletion is synchronous.
type httpApiDomainVerifier struct{}

func (*httpApiDomainVerifier) IDOutputKey() string { return "domain_name" }

func (v *httpApiDomainVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := httpApiDomainExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awshttpapidomain verify-exists failed for %q", id)
	}
	if !exists {
		return errors.Errorf("awshttpapidomain %q not found after deploy", id)
	}
	return nil
}

func (v *httpApiDomainVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := httpApiDomainExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awshttpapidomain verify-absent failed for %q", id)
	}
	if exists {
		return errors.Errorf("awshttpapidomain %q still exists after destroy", id)
	}
	return nil
}

func httpApiDomainExists(ctx context.Context, cfg aws.Config, domainName, region string) (bool, error) {
	client := apigatewayv2Client(cfg, region)
	_, err := client.GetDomainName(ctx, &apigatewayv2.GetDomainNameInput{DomainName: aws.String(domainName)})
	if err != nil {
		var notFound *types.NotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
