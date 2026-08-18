package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/appsync"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// eventsApi creates the Events pivot: per-phase (connect/publish/
// subscribe) authorization against a shared provider list, and
// optional event logging. Channel namespaces render with the other
// satellites (they join data sources by name).
func eventsApi(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*createdApi, error) {
	events := locals.Spec.Events

	authProviders := appsync.ApiEventConfigAuthProviderArray{}
	for _, authProvider := range events.AuthProviders {
		providerArgs := &appsync.ApiEventConfigAuthProviderArgs{
			AuthType: pulumi.String(authProvider.Type),
		}
		if cognito := authProvider.Cognito; cognito != nil {
			cognitoArgs := &appsync.ApiEventConfigAuthProviderCognitoConfigArgs{
				UserPoolId: pulumi.String(cognito.UserPoolId.GetValue()),
				// Required by AWS on Events APIs (the GraphQL arm's
				// asymmetry - there it defaults to the API's region).
				AwsRegion: pulumi.String(cognito.AwsRegion),
			}
			if cognito.AppIdClientRegex != "" {
				cognitoArgs.AppIdClientRegex = pulumi.String(cognito.AppIdClientRegex)
			}
			providerArgs.CognitoConfig = cognitoArgs
		}
		if oidc := authProvider.OpenidConnect; oidc != nil {
			oidcArgs := &appsync.ApiEventConfigAuthProviderOpenidConnectConfigArgs{
				Issuer: pulumi.String(oidc.Issuer),
			}
			if oidc.ClientId != "" {
				oidcArgs.ClientId = pulumi.String(oidc.ClientId)
			}
			if oidc.IatTtl > 0 {
				oidcArgs.IatTtl = pulumi.Int(int(oidc.IatTtl))
			}
			if oidc.AuthTtl > 0 {
				oidcArgs.AuthTtl = pulumi.Int(int(oidc.AuthTtl))
			}
			providerArgs.OpenidConnectConfig = oidcArgs
		}
		if lambda := authProvider.Lambda; lambda != nil {
			lambdaArgs := &appsync.ApiEventConfigAuthProviderLambdaAuthorizerConfigArgs{
				AuthorizerUri: pulumi.String(lambda.AuthorizerUri.GetValue()),
			}
			if lambda.AuthorizerResultTtlInSeconds > 0 {
				lambdaArgs.AuthorizerResultTtlInSeconds = pulumi.Int(int(lambda.AuthorizerResultTtlInSeconds))
			}
			if lambda.IdentityValidationExpression != "" {
				lambdaArgs.IdentityValidationExpression = pulumi.String(lambda.IdentityValidationExpression)
			}
			providerArgs.LambdaAuthorizerConfig = lambdaArgs
		}
		authProviders = append(authProviders, providerArgs)
	}

	connectionModes := appsync.ApiEventConfigConnectionAuthModeArray{}
	for _, mode := range events.ConnectionAuthModes {
		connectionModes = append(connectionModes, &appsync.ApiEventConfigConnectionAuthModeArgs{
			AuthType: pulumi.String(mode),
		})
	}
	publishModes := appsync.ApiEventConfigDefaultPublishAuthModeArray{}
	for _, mode := range events.DefaultPublishAuthModes {
		publishModes = append(publishModes, &appsync.ApiEventConfigDefaultPublishAuthModeArgs{
			AuthType: pulumi.String(mode),
		})
	}
	subscribeModes := appsync.ApiEventConfigDefaultSubscribeAuthModeArray{}
	for _, mode := range events.DefaultSubscribeAuthModes {
		subscribeModes = append(subscribeModes, &appsync.ApiEventConfigDefaultSubscribeAuthModeArgs{
			AuthType: pulumi.String(mode),
		})
	}

	eventConfigArgs := &appsync.ApiEventConfigArgs{
		AuthProviders:             authProviders,
		ConnectionAuthModes:       connectionModes,
		DefaultPublishAuthModes:   publishModes,
		DefaultSubscribeAuthModes: subscribeModes,
	}

	if logConfig := events.LogConfig; logConfig != nil {
		eventConfigArgs.LogConfig = &appsync.ApiEventConfigLogConfigArgs{
			CloudwatchLogsRoleArn: pulumi.String(logConfig.CloudwatchLogsRoleArn.GetValue()),
			LogLevel:              pulumi.String(logConfig.LogLevel),
		}
	}

	apiArgs := &appsync.ApiArgs{
		// Events API names allow hyphens - metadata.name is the naming
		// basis.
		Name:        pulumi.String(locals.Target.Metadata.Name),
		EventConfig: eventConfigArgs,
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}
	if events.OwnerContact != "" {
		apiArgs.OwnerContact = pulumi.String(events.OwnerContact)
	}

	createdEventsApi, err := appsync.NewApi(ctx, "api", apiArgs, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create events api")
	}

	return &createdApi{
		ApiId:          createdEventsApi.ApiId,
		ApiArn:         createdEventsApi.ApiArn,
		GraphqlUrl:     pulumi.String("").ToStringOutput(),
		RealtimeUrl:    pulumi.String("").ToStringOutput(),
		EventsHttp:     createdEventsApi.Dns.MapIndex(pulumi.String("HTTP")),
		EventsRealtime: createdEventsApi.Dns.MapIndex(pulumi.String("REALTIME")),
	}, nil
}
