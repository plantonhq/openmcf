package verify

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
	"github.com/pkg/errors"
)

// httpApiGatewayVerifier verifies an AwsHttpApiGateway via GetApi, keyed on
// the api_id output. API deletion is synchronous, so existence is a plain
// found/not-found check on the typed NotFoundException.
type httpApiGatewayVerifier struct{}

func (*httpApiGatewayVerifier) IDOutputKey() string { return "api_id" }

func (v *httpApiGatewayVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := httpApiExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awshttpapigateway verify-exists failed for %q", id)
	}
	if !exists {
		return errors.Errorf("awshttpapigateway %q not found after deploy", id)
	}
	return nil
}

func (v *httpApiGatewayVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	exists, err := httpApiExists(ctx, cfg, id, region)
	if err != nil {
		return errors.Wrapf(err, "awshttpapigateway verify-absent failed for %q", id)
	}
	if exists {
		return errors.Errorf("awshttpapigateway %q still exists after destroy", id)
	}
	return nil
}

func httpApiExists(ctx context.Context, cfg aws.Config, apiId, region string) (bool, error) {
	client := apigatewayv2Client(cfg, region)
	_, err := client.GetApi(ctx, &apigatewayv2.GetApiInput{ApiId: aws.String(apiId)})
	if err != nil {
		var notFound *types.NotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// apigatewayv2Client builds the shared API Gateway v2 client the HTTP API,
// VPC link, and domain verifiers all use.
func apigatewayv2Client(cfg aws.Config, region string) *apigatewayv2.Client {
	return apigatewayv2.NewFromConfig(cfg, func(o *apigatewayv2.Options) {
		if region != "" {
			o.Region = region
		}
	})
}
