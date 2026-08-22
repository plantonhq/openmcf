package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	pkgerrors "github.com/pkg/errors"
)

// rdsProxyVerifier verifies an AwsRdsProxy via DescribeDBProxies,
// keyed on the proxy name (the provider's import ID - the API has no
// ARN-addressed reads). Exists demands status AVAILABLE, never mere
// presence: a proxy whose secrets or role are broken sits in
// INCOMPATIBLE_* states while still existing. From outputs, the
// additional endpoints and the registered target are asserted too.
type rdsProxyVerifier struct{}

func (*rdsProxyVerifier) IDOutputKey() string { return "proxy_name" }

func (*rdsProxyVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	out, err := rds.NewFromConfig(cfg).DescribeDBProxies(ctx, &rds.DescribeDBProxiesInput{
		DBProxyName: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsrdsproxy verify-exists failed for %q", id)
	}
	if len(out.DBProxies) == 0 {
		return pkgerrors.Errorf("awsrdsproxy %q not found after deploy", id)
	}
	if status := out.DBProxies[0].Status; status != rdstypes.DBProxyStatusAvailable {
		return pkgerrors.Errorf("awsrdsproxy %q status %q, want available", id, status)
	}
	return nil
}

func (*rdsProxyVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	out, err := rds.NewFromConfig(cfg).DescribeDBProxies(ctx, &rds.DescribeDBProxiesInput{
		DBProxyName: aws.String(id),
	})
	if err != nil {
		var notFound *rdstypes.DBProxyNotFoundFault
		if pkgerrors.As(err, &notFound) {
			return nil
		}
		return pkgerrors.Wrapf(err, "awsrdsproxy verify-absent failed for %q", id)
	}
	if len(out.DBProxies) > 0 {
		return pkgerrors.Errorf("awsrdsproxy %q still exists after destroy", id)
	}
	return nil
}

// VerifyExistsFromOutputs additionally walks the declared endpoints
// (each must be AVAILABLE) and the registered target.
func (v *rdsProxyVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	proxyName, _ := outputs["proxy_name"].(string)
	if proxyName == "" {
		return pkgerrors.New("awsrdsproxy outputs carry no proxy_name")
	}
	if err := v.VerifyExists(ctx, cfg, proxyName, region); err != nil {
		return err
	}
	client := rds.NewFromConfig(cfg)

	endpointAddresses, _ := outputs["endpoint_addresses"].(map[string]interface{})
	for endpointName := range endpointAddresses {
		out, err := client.DescribeDBProxyEndpoints(ctx, &rds.DescribeDBProxyEndpointsInput{
			DBProxyName:         aws.String(proxyName),
			DBProxyEndpointName: aws.String(endpointName),
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awsrdsproxy endpoint %q read failed", endpointName)
		}
		if len(out.DBProxyEndpoints) == 0 {
			return pkgerrors.Errorf("awsrdsproxy endpoint %q not found after deploy", endpointName)
		}
		if status := out.DBProxyEndpoints[0].Status; status != rdstypes.DBProxyEndpointStatusAvailable {
			return pkgerrors.Errorf("awsrdsproxy endpoint %q status %q, want available", endpointName, status)
		}
	}

	if targetResourceId, _ := outputs["target_rds_resource_id"].(string); targetResourceId != "" {
		out, err := client.DescribeDBProxyTargets(ctx, &rds.DescribeDBProxyTargetsInput{
			DBProxyName: aws.String(proxyName),
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awsrdsproxy target read failed for %q", proxyName)
		}
		found := false
		for _, target := range out.Targets {
			if target.RdsResourceId != nil && *target.RdsResourceId == targetResourceId {
				found = true
				break
			}
		}
		if !found {
			return pkgerrors.Errorf("awsrdsproxy %q has no registered target %q", proxyName, targetResourceId)
		}
	}
	return nil
}
