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
// the API Gateway v2 API). Domain deletion is synchronous. Existence
// additionally asserts the routing posture from the domain's own state:
// when the domain's routing mode uses routing rules, ListRoutingRules must
// return at least one rule -- an apply that created the domain but silently
// dropped its rules would otherwise verify clean.
type httpApiDomainVerifier struct{}

func (*httpApiDomainVerifier) IDOutputKey() string { return "domain_name" }

func (v *httpApiDomainVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := apigatewayv2Client(cfg, region)
	domain, err := client.GetDomainName(ctx, &apigatewayv2.GetDomainNameInput{DomainName: aws.String(id)})
	if err != nil {
		return errors.Wrapf(err, "awshttpapidomain verify-exists failed for %q", id)
	}
	mode := domain.RoutingMode
	if mode != types.RoutingModeRoutingRuleOnly && mode != types.RoutingModeRoutingRuleThenApiMapping {
		return nil
	}
	rules, err := client.ListRoutingRules(ctx, &apigatewayv2.ListRoutingRulesInput{DomainName: aws.String(id)})
	if err != nil {
		return errors.Wrapf(err, "awshttpapidomain %q: listing routing rules failed", id)
	}
	if len(rules.RoutingRules) == 0 {
		return errors.Errorf("awshttpapidomain %q: routing mode %s but ListRoutingRules returned no rules", id, mode)
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
