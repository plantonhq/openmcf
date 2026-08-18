package module

import (
	"github.com/pkg/errors"
	awseventbridgeapidestinationv1alpha1 "github.com/plantonhq/planton/catalog/aws/awseventbridgeapidestination/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudwatch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// apiDestination renders the connection and/or destination arms (the
// spec's CELs guarantee at least one arm and exactly one connection
// source for the destination) and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - both names are fixed for life (replace-on-change);
//   - AuthorizationType is DERIVED from whichever auth block the spec
//     sets - the two can never disagree;
//   - AWS stores credential values in a Secrets Manager secret it
//     creates and owns (the connection_secret_arn output);
//     DescribeConnection never returns them, so imports cannot
//     recover them (declared write-normalized in the import catalog);
//   - connection creates/updates wait through the auth state machine
//     (CREATING/AUTHORIZING -> AUTHORIZED, up to 20 minutes);
//   - neither resource is taggable at AWS (the deliberate
//     tag-convention absence).
func apiDestination(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	var createdConnection *cloudwatch.EventConnection

	if spec.Connection != nil {
		var err error
		createdConnection, err = connectionArm(ctx, locals, provider)
		if err != nil {
			return errors.Wrap(err, "connection arm")
		}
		ctx.Export(OpConnectionArn, createdConnection.Arn)
		ctx.Export(OpConnectionSecretArn, createdConnection.SecretArn)
	} else {
		ctx.Export(OpConnectionArn, pulumi.String(""))
		ctx.Export(OpConnectionSecretArn, pulumi.String(""))
	}

	if spec.Destination != nil {
		destination := spec.Destination

		// The owned connection when the instance has one, else the
		// external ARN (spec CEL: exactly one source).
		var connectionArn pulumi.StringInput
		if createdConnection != nil {
			connectionArn = createdConnection.Arn
		} else {
			connectionArn = pulumi.String(destination.ConnectionArn.GetValue())
		}

		args := &cloudwatch.EventApiDestinationArgs{
			Name:               pulumi.String(destination.Name),
			ConnectionArn:      connectionArn,
			InvocationEndpoint: pulumi.String(destination.InvocationEndpoint),
			HttpMethod:         pulumi.String(destination.HttpMethod),
		}
		if destination.Description != "" {
			args.Description = pulumi.String(destination.Description)
		}
		if destination.InvocationRateLimitPerSecond != nil {
			args.InvocationRateLimitPerSecond = pulumi.Int(int(*destination.InvocationRateLimitPerSecond))
		}

		createdDestination, err := cloudwatch.NewEventApiDestination(ctx, "api_destination", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create api destination")
		}
		ctx.Export(OpApiDestinationArn, createdDestination.Arn)
	} else {
		ctx.Export(OpApiDestinationArn, pulumi.String(""))
	}

	return nil
}

func connectionArm(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*cloudwatch.EventConnection, error) {
	connection := locals.Spec.Connection

	authParameters := &cloudwatch.EventConnectionAuthParametersArgs{}
	var authorizationType string

	switch {
	case connection.ApiKey != nil:
		authorizationType = "API_KEY"
		authParameters.ApiKey = &cloudwatch.EventConnectionAuthParametersApiKeyArgs{
			Key:   pulumi.String(connection.ApiKey.Key),
			Value: pulumi.String(connection.ApiKey.Value),
		}
	case connection.Basic != nil:
		authorizationType = "BASIC"
		authParameters.Basic = &cloudwatch.EventConnectionAuthParametersBasicArgs{
			Username: pulumi.String(connection.Basic.Username),
			Password: pulumi.String(connection.Basic.Password),
		}
	default:
		authorizationType = "OAUTH_CLIENT_CREDENTIALS"
		oauth := connection.Oauth
		authParameters.Oauth = &cloudwatch.EventConnectionAuthParametersOauthArgs{
			AuthorizationEndpoint: pulumi.String(oauth.AuthorizationEndpoint),
			HttpMethod:            pulumi.String(oauth.HttpMethod),
			ClientParameters: &cloudwatch.EventConnectionAuthParametersOauthClientParametersArgs{
				ClientId:     pulumi.String(oauth.ClientId),
				ClientSecret: pulumi.String(oauth.ClientSecret),
			},
			OauthHttpParameters: buildOauthHttpParameters(oauth.OauthHttpParameters),
		}
	}

	if connection.InvocationHttpParameters != nil {
		authParameters.InvocationHttpParameters = buildInvocationHttpParameters(connection.InvocationHttpParameters)
	}

	if connection.PrivateAuthorizationEndpoint != nil {
		authParameters.ConnectivityParameters = &cloudwatch.EventConnectionAuthParametersConnectivityParametersArgs{
			ResourceParameters: &cloudwatch.EventConnectionAuthParametersConnectivityParametersResourceParametersArgs{
				ResourceConfigurationArn: pulumi.String(connection.PrivateAuthorizationEndpoint.ResourceConfigurationArn),
			},
		}
	}

	args := &cloudwatch.EventConnectionArgs{
		Name:              pulumi.String(connection.Name),
		AuthorizationType: pulumi.String(authorizationType),
		AuthParameters:    authParameters,
	}
	if connection.Description != "" {
		args.Description = pulumi.String(connection.Description)
	}
	if connection.PrivateInvocationEndpoint != nil {
		args.InvocationConnectivityParameters = &cloudwatch.EventConnectionInvocationConnectivityParametersArgs{
			ResourceParameters: &cloudwatch.EventConnectionInvocationConnectivityParametersResourceParametersArgs{
				ResourceConfigurationArn: pulumi.String(connection.PrivateInvocationEndpoint.ResourceConfigurationArn),
			},
		}
	}
	if connection.KmsKeyIdentifier.GetValue() != "" {
		args.KmsKeyIdentifier = pulumi.String(connection.KmsKeyIdentifier.GetValue())
	}

	return cloudwatch.NewEventConnection(ctx, "connection", args, pulumi.Provider(provider))
}

func buildOauthHttpParameters(params *awseventbridgeapidestinationv1alpha1.AwsEventBridgeConnectionHttpParameters) *cloudwatch.EventConnectionAuthParametersOauthOauthHttpParametersArgs {
	args := &cloudwatch.EventConnectionAuthParametersOauthOauthHttpParametersArgs{}

	bodies := cloudwatch.EventConnectionAuthParametersOauthOauthHttpParametersBodyArray{}
	for _, parameter := range params.Body {
		bodies = append(bodies, &cloudwatch.EventConnectionAuthParametersOauthOauthHttpParametersBodyArgs{
			Key:           pulumi.String(parameter.Key),
			Value:         pulumi.String(parameter.Value),
			IsValueSecret: pulumi.Bool(parameter.IsValueSecret),
		})
	}
	if len(bodies) > 0 {
		args.Bodies = bodies
	}

	headers := cloudwatch.EventConnectionAuthParametersOauthOauthHttpParametersHeaderArray{}
	for _, parameter := range params.Header {
		headers = append(headers, &cloudwatch.EventConnectionAuthParametersOauthOauthHttpParametersHeaderArgs{
			Key:           pulumi.String(parameter.Key),
			Value:         pulumi.String(parameter.Value),
			IsValueSecret: pulumi.Bool(parameter.IsValueSecret),
		})
	}
	if len(headers) > 0 {
		args.Headers = headers
	}

	queryStrings := cloudwatch.EventConnectionAuthParametersOauthOauthHttpParametersQueryStringArray{}
	for _, parameter := range params.QueryString {
		queryStrings = append(queryStrings, &cloudwatch.EventConnectionAuthParametersOauthOauthHttpParametersQueryStringArgs{
			Key:           pulumi.String(parameter.Key),
			Value:         pulumi.String(parameter.Value),
			IsValueSecret: pulumi.Bool(parameter.IsValueSecret),
		})
	}
	if len(queryStrings) > 0 {
		args.QueryStrings = queryStrings
	}

	return args
}

func buildInvocationHttpParameters(params *awseventbridgeapidestinationv1alpha1.AwsEventBridgeConnectionHttpParameters) *cloudwatch.EventConnectionAuthParametersInvocationHttpParametersArgs {
	args := &cloudwatch.EventConnectionAuthParametersInvocationHttpParametersArgs{}

	bodies := cloudwatch.EventConnectionAuthParametersInvocationHttpParametersBodyArray{}
	for _, parameter := range params.Body {
		bodies = append(bodies, &cloudwatch.EventConnectionAuthParametersInvocationHttpParametersBodyArgs{
			Key:           pulumi.String(parameter.Key),
			Value:         pulumi.String(parameter.Value),
			IsValueSecret: pulumi.Bool(parameter.IsValueSecret),
		})
	}
	if len(bodies) > 0 {
		args.Bodies = bodies
	}

	headers := cloudwatch.EventConnectionAuthParametersInvocationHttpParametersHeaderArray{}
	for _, parameter := range params.Header {
		headers = append(headers, &cloudwatch.EventConnectionAuthParametersInvocationHttpParametersHeaderArgs{
			Key:           pulumi.String(parameter.Key),
			Value:         pulumi.String(parameter.Value),
			IsValueSecret: pulumi.Bool(parameter.IsValueSecret),
		})
	}
	if len(headers) > 0 {
		args.Headers = headers
	}

	queryStrings := cloudwatch.EventConnectionAuthParametersInvocationHttpParametersQueryStringArray{}
	for _, parameter := range params.QueryString {
		queryStrings = append(queryStrings, &cloudwatch.EventConnectionAuthParametersInvocationHttpParametersQueryStringArgs{
			Key:           pulumi.String(parameter.Key),
			Value:         pulumi.String(parameter.Value),
			IsValueSecret: pulumi.Bool(parameter.IsValueSecret),
		})
	}
	if len(queryStrings) > 0 {
		args.QueryStrings = queryStrings
	}

	return args
}
