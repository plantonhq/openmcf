package module

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// restApi creates the REST API itself (with the OpenAPI body when that
// arm defines the API) and its resource policy.
//
// Lifecycle facts the render below depends on:
//   - with a body, the provider runs CreateRestApi -> PutRestApi -> a
//     reconciliation pass that re-applies configured literals the
//     overwrite-mode import wiped (description, endpoint settings,
//     policy, ...) -- expected apply-log noise, not drift;
//   - minimum_compression_size is the provider's nullable-int-as-string
//     quirk: unset means compression disabled, "0" compresses
//     everything;
//   - the standalone rest_api_policy resource owns the policy (clean
//     PATCH updates; delete resets to empty instead of touching the
//     API), so the API's own policy argument stays unset.
func restApi(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*apigateway.RestApi, error) {
	spec := locals.Spec

	args := &apigateway.RestApiArgs{
		// metadata.name is the naming basis on both engines.
		Name: pulumi.String(locals.Target.Metadata.Name),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.ApiKeySource != "" {
		args.ApiKeySource = pulumi.String(spec.ApiKeySource)
	}
	if len(spec.BinaryMediaTypes) > 0 {
		args.BinaryMediaTypes = pulumi.ToStringArray(spec.BinaryMediaTypes)
	}
	if spec.MinimumCompressionSize != nil {
		args.MinimumCompressionSize = pulumi.String(strconv.Itoa(int(*spec.MinimumCompressionSize)))
	}
	if spec.DisableExecuteApiEndpoint {
		args.DisableExecuteApiEndpoint = pulumi.Bool(true)
	}
	if spec.EndpointConfiguration != nil {
		endpoint := &apigateway.RestApiEndpointConfigurationArgs{
			// The bridge flattens the provider's max-1 types list to a
			// singular string.
			Types: pulumi.String(spec.EndpointConfiguration.Type),
		}
		if spec.EndpointConfiguration.IpAddressType != "" {
			endpoint.IpAddressType = pulumi.String(spec.EndpointConfiguration.IpAddressType)
		}
		if len(spec.EndpointConfiguration.VpcEndpointIds) > 0 {
			endpoint.VpcEndpointIds = svrsToStringArray(spec.EndpointConfiguration.VpcEndpointIds)
		}
		args.EndpointConfiguration = endpoint
	}
	if spec.EndpointAccessMode != "" {
		args.EndpointAccessMode = pulumi.String(spec.EndpointAccessMode)
	}
	if spec.SecurityPolicy != "" {
		args.SecurityPolicy = pulumi.String(spec.SecurityPolicy)
	}

	// The OpenAPI arm: the document IS the API definition.
	if spec.Openapi != nil {
		args.Body = pulumi.String(spec.Openapi.Body)
		if spec.Openapi.FailOnWarnings {
			args.FailOnWarnings = pulumi.Bool(true)
		}
		if len(spec.Openapi.Parameters) > 0 {
			args.Parameters = pulumi.ToStringMap(spec.Openapi.Parameters)
		}
		if spec.Openapi.Mode != "" {
			args.PutRestApiMode = pulumi.String(spec.Openapi.Mode)
		}
	}

	api, err := apigateway.NewRestApi(ctx, "rest-api", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create rest api")
	}

	if spec.Policy != nil {
		policyJson, err := json.Marshal(spec.Policy.AsMap())
		if err != nil {
			return nil, errors.Wrap(err, "marshal rest api policy")
		}
		_, err = apigateway.NewRestApiPolicy(ctx, "rest-api-policy", &apigateway.RestApiPolicyArgs{
			RestApiId: api.ID(),
			Policy:    pulumi.String(string(policyJson)),
		}, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrap(err, "create rest api policy")
		}
	}

	ctx.Export(OpRestApiId, api.ID())
	ctx.Export(OpRestApiArn, api.Arn)
	ctx.Export(OpExecutionArn, api.ExecutionArn)
	ctx.Export(OpRootResourceId, api.RootResourceId)
	return api, nil
}
