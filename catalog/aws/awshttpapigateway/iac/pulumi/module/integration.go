package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// integrations creates the deduplicated API Gateway integrations (see
// dedupIntegrations in locals.go for the whole-object dedup rule) and returns
// them indexed in first-appearance order, matching the route-to-integration
// index that dedupIntegrations computes.
func integrations(
	ctx *pulumi.Context,
	locals *Locals,
	createdApi *apigatewayv2.Api,
	provider *aws.Provider,
) ([]*apigatewayv2.Integration, error) {
	distinct, _ := dedupIntegrations(locals.Spec.Routes)
	result := make([]*apigatewayv2.Integration, 0, len(distinct))

	for i, integration := range distinct {
		resourceName := fmt.Sprintf("%s-integration-%d", locals.ApiName, i+1)

		args := &apigatewayv2.IntegrationArgs{
			ApiId:           createdApi.ID(),
			IntegrationType: pulumi.String(integration.IntegrationType),
		}

		// Proxy integrations carry their target here; AWS service
		// integrations (integration_subtype) express the target in
		// request_parameters and AWS rejects a URI alongside a subtype --
		// the spec CEL enforces the split.
		if integration.IntegrationUri.GetValue() != "" {
			args.IntegrationUri = pulumi.StringPtr(integration.IntegrationUri.GetValue())
		}
		if integration.IntegrationSubtype != "" {
			args.IntegrationSubtype = pulumi.StringPtr(integration.IntegrationSubtype)
		}

		// Only AWS_PROXY (Lambda) and AWS service subtypes carry a payload
		// format. HTTP_PROXY rejects PayloadFormatVersion entirely.
		switch {
		case integration.IntegrationSubtype != "":
			args.PayloadFormatVersion = pulumi.StringPtr("1.0")
		case integration.IntegrationType == "AWS_PROXY":
			payloadVersion := "2.0"
			if integration.PayloadFormatVersion != "" {
				payloadVersion = integration.PayloadFormatVersion
			}
			args.PayloadFormatVersion = pulumi.StringPtr(payloadVersion)
		}

		// Lambda integrations are always invoked with POST regardless of
		// this value.
		if integration.IntegrationMethod != "" {
			args.IntegrationMethod = pulumi.StringPtr(integration.IntegrationMethod)
		}

		if integration.TimeoutMilliseconds > 0 {
			args.TimeoutMilliseconds = pulumi.IntPtr(int(integration.TimeoutMilliseconds))
		}

		// Private integrations reach through a VPC link; the spec CELs
		// guarantee VPC_LINK <=> connection_id and HTTP_PROXY-only.
		if integration.ConnectionType != "" {
			args.ConnectionType = pulumi.StringPtr(integration.ConnectionType)
		}
		if integration.ConnectionId.GetValue() != "" {
			args.ConnectionId = pulumi.StringPtr(integration.ConnectionId.GetValue())
		}

		// The role API Gateway assumes to call an AWS service action
		// (required for subtype integrations by spec CEL); Lambda proxies
		// normally rely on the function's resource policy instead.
		if integration.CredentialsArn.GetValue() != "" {
			args.CredentialsArn = pulumi.StringPtr(integration.CredentialsArn.GetValue())
		}

		// Parameter mappings (proxy) or service-action parameters (subtype).
		if len(integration.RequestParameters) > 0 {
			args.RequestParameters = pulumi.ToStringMap(integration.RequestParameters)
		}

		// Response transforms keyed by backend status code.
		if len(integration.ResponseParameters) > 0 {
			respParams := make(apigatewayv2.IntegrationResponseParameterArray, 0, len(integration.ResponseParameters))
			for _, rp := range integration.ResponseParameters {
				respParams = append(respParams, &apigatewayv2.IntegrationResponseParameterArgs{
					StatusCode: pulumi.String(rp.StatusCode),
					Mappings:   pulumi.ToStringMap(rp.Mappings),
				})
			}
			args.ResponseParameters = respParams
		}

		// SNI override for private integrations whose internal ALB serves a
		// public-domain certificate.
		if integration.TlsServerNameToVerify != "" {
			args.TlsConfig = &apigatewayv2.IntegrationTlsConfigArgs{
				ServerNameToVerify: pulumi.StringPtr(integration.TlsServerNameToVerify),
			}
		}

		if integration.Description != "" {
			args.Description = pulumi.StringPtr(integration.Description)
		}

		created, err := apigatewayv2.NewIntegration(ctx, resourceName, args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create integration %d", i+1)
		}

		result = append(result, created)
	}

	return result, nil
}
