package module

import (
	"github.com/pkg/errors"
	awsappsyncapiv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsappsyncapi/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/appsync"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/wafv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createdApi carries the arm-agnostic handles every satellite hangs
// off, plus the endpoint outputs each arm fills for its own surface.
type createdApi struct {
	ApiId  pulumi.StringOutput
	ApiArn pulumi.StringOutput

	GraphqlUrl     pulumi.StringOutput
	RealtimeUrl    pulumi.StringOutput
	EventsHttp     pulumi.StringOutput
	EventsRealtime pulumi.StringOutput
}

// graphqlApi creates the GraphQL pivot with its auth providers, the
// WAF association, and the cache singleton.
//
// Lifecycle facts the render depends on:
//   - api_type (derived from the merged block) and visibility are
//     ForceNew; the schema applies via async StartSchemaCreation with
//     NO provider drift detection;
//   - the cache is one-per-API (its AWS id IS the API id); both
//     encryption flags replace it on change; operations wait up to 60
//     minutes upstream;
//   - the web ACL attaches from the protected resource's side (the
//     AwsAlb pattern); AWS's WAF association supports GraphQL APIs
//     only.
func graphqlApi(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*createdApi, error) {
	graphql := locals.Spec.Graphql

	apiArgs := &appsync.GraphQLApiArgs{
		// GraphQL API names forbid hyphens - the spec carries an
		// explicit api_name (the explicit-name-field convention).
		Name:               pulumi.String(graphql.ApiName),
		AuthenticationType: pulumi.String(graphql.Auth.Type),
		Tags:               pulumi.ToStringMap(locals.AwsTags),
	}

	if locals.IsMerged {
		// MERGED is derived from the merged block - never spec surface.
		apiArgs.ApiType = pulumi.String("MERGED")
		apiArgs.MergedApiExecutionRoleArn = pulumi.String(graphql.Merged.ExecutionRoleArn.GetValue())
	}

	if userPool := graphql.Auth.UserPool; userPool != nil {
		userPoolArgs := &appsync.GraphQLApiUserPoolConfigArgs{
			UserPoolId:    pulumi.String(userPool.UserPoolId.GetValue()),
			DefaultAction: pulumi.String(userPool.DefaultAction),
		}
		if userPool.AppIdClientRegex != "" {
			userPoolArgs.AppIdClientRegex = pulumi.String(userPool.AppIdClientRegex)
		}
		if userPool.AwsRegion != "" {
			userPoolArgs.AwsRegion = pulumi.String(userPool.AwsRegion)
		}
		apiArgs.UserPoolConfig = userPoolArgs
	}

	if oidc := graphql.Auth.OpenidConnect; oidc != nil {
		oidcArgs := &appsync.GraphQLApiOpenidConnectConfigArgs{
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
		apiArgs.OpenidConnectConfig = oidcArgs
	}

	if lambda := graphql.Auth.Lambda; lambda != nil {
		lambdaArgs := &appsync.GraphQLApiLambdaAuthorizerConfigArgs{
			AuthorizerUri: pulumi.String(lambda.AuthorizerUri.GetValue()),
		}
		if lambda.AuthorizerResultTtlInSeconds > 0 {
			lambdaArgs.AuthorizerResultTtlInSeconds = pulumi.Int(int(lambda.AuthorizerResultTtlInSeconds))
		}
		if lambda.IdentityValidationExpression != "" {
			lambdaArgs.IdentityValidationExpression = pulumi.String(lambda.IdentityValidationExpression)
		}
		apiArgs.LambdaAuthorizerConfig = lambdaArgs
	}

	additionalProviders := appsync.GraphQLApiAdditionalAuthenticationProviderArray{}
	for _, additional := range graphql.AdditionalAuthProviders {
		additionalArgs := &appsync.GraphQLApiAdditionalAuthenticationProviderArgs{
			AuthenticationType: pulumi.String(additional.Type),
		}
		// The additional provider's user pool carries NO
		// default_action (AWS's asymmetry; the spec's CEL walls
		// enforce it).
		if userPool := additional.UserPool; userPool != nil {
			userPoolArgs := &appsync.GraphQLApiAdditionalAuthenticationProviderUserPoolConfigArgs{
				UserPoolId: pulumi.String(userPool.UserPoolId.GetValue()),
			}
			if userPool.AppIdClientRegex != "" {
				userPoolArgs.AppIdClientRegex = pulumi.String(userPool.AppIdClientRegex)
			}
			if userPool.AwsRegion != "" {
				userPoolArgs.AwsRegion = pulumi.String(userPool.AwsRegion)
			}
			additionalArgs.UserPoolConfig = userPoolArgs
		}
		if oidc := additional.OpenidConnect; oidc != nil {
			oidcArgs := &appsync.GraphQLApiAdditionalAuthenticationProviderOpenidConnectConfigArgs{
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
			additionalArgs.OpenidConnectConfig = oidcArgs
		}
		if lambda := additional.Lambda; lambda != nil {
			lambdaArgs := &appsync.GraphQLApiAdditionalAuthenticationProviderLambdaAuthorizerConfigArgs{
				AuthorizerUri: pulumi.String(lambda.AuthorizerUri.GetValue()),
			}
			if lambda.AuthorizerResultTtlInSeconds > 0 {
				lambdaArgs.AuthorizerResultTtlInSeconds = pulumi.Int(int(lambda.AuthorizerResultTtlInSeconds))
			}
			if lambda.IdentityValidationExpression != "" {
				lambdaArgs.IdentityValidationExpression = pulumi.String(lambda.IdentityValidationExpression)
			}
			additionalArgs.LambdaAuthorizerConfig = lambdaArgs
		}
		additionalProviders = append(additionalProviders, additionalArgs)
	}
	if len(additionalProviders) > 0 {
		apiArgs.AdditionalAuthenticationProviders = additionalProviders
	}

	// MERGED APIs own no schema - it merges in from the sources.
	if graphql.Schema != "" {
		apiArgs.Schema = pulumi.String(graphql.Schema)
	}
	if graphql.Visibility != "" {
		apiArgs.Visibility = pulumi.String(graphql.Visibility)
	}
	if graphql.DisableIntrospection {
		apiArgs.IntrospectionConfig = pulumi.String("DISABLED")
	}
	if graphql.QueryDepthLimit > 0 {
		apiArgs.QueryDepthLimit = pulumi.Int(int(graphql.QueryDepthLimit))
	}
	if graphql.ResolverCountLimit > 0 {
		apiArgs.ResolverCountLimit = pulumi.Int(int(graphql.ResolverCountLimit))
	}
	if graphql.XrayEnabled {
		apiArgs.XrayEnabled = pulumi.Bool(true)
	}

	if logConfig := graphql.LogConfig; logConfig != nil {
		apiArgs.LogConfig = &appsync.GraphQLApiLogConfigArgs{
			CloudwatchLogsRoleArn: pulumi.String(logConfig.CloudwatchLogsRoleArn.GetValue()),
			FieldLogLevel:         pulumi.String(logConfig.FieldLogLevel),
			ExcludeVerboseContent: pulumi.Bool(logConfig.ExcludeVerboseContent),
		}
	}

	if metrics := graphql.EnhancedMetrics; metrics != nil {
		apiArgs.EnhancedMetricsConfig = &appsync.GraphQLApiEnhancedMetricsConfigArgs{
			DataSourceLevelMetricsBehavior: pulumi.String(metrics.DataSourceLevelMetricsBehavior),
			OperationLevelMetricsConfig:    pulumi.String(metrics.OperationLevelMetricsConfig),
			ResolverLevelMetricsBehavior:   pulumi.String(metrics.ResolverLevelMetricsBehavior),
		}
	}

	createdGraphqlApi, err := appsync.NewGraphQLApi(ctx, "api", apiArgs, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create graphql api")
	}

	if graphql.WebAclArn.GetValue() != "" {
		if _, err := wafv2.NewWebAclAssociation(ctx, "web-acl-association",
			&wafv2.WebAclAssociationArgs{
				ResourceArn: createdGraphqlApi.Arn,
				WebAclArn:   pulumi.String(graphql.WebAclArn.GetValue()),
			}, pulumi.Provider(provider)); err != nil {
			return nil, errors.Wrap(err, "associate web acl")
		}
	}

	if cache := graphql.Cache; cache != nil {
		if _, err := apiCache(ctx, locals, provider, createdGraphqlApi, cache); err != nil {
			return nil, errors.Wrap(err, "create api cache")
		}
	}

	return &createdApi{
		ApiId:          createdGraphqlApi.ID().ToStringOutput(),
		ApiArn:         createdGraphqlApi.Arn,
		GraphqlUrl:     createdGraphqlApi.Uris.MapIndex(pulumi.String("GRAPHQL")),
		RealtimeUrl:    createdGraphqlApi.Uris.MapIndex(pulumi.String("REALTIME")),
		EventsHttp:     pulumi.String("").ToStringOutput(),
		EventsRealtime: pulumi.String("").ToStringOutput(),
	}, nil
}

// apiCache creates the one-per-API cache singleton.
func apiCache(ctx *pulumi.Context, _ *Locals, provider *aws.Provider,
	createdGraphqlApi *appsync.GraphQLApi, cache *awsappsyncapiv1alpha1.AwsAppSyncGraphqlCache) (*appsync.ApiCache, error) {
	cacheArgs := &appsync.ApiCacheArgs{
		ApiId:              createdGraphqlApi.ID(),
		ApiCachingBehavior: pulumi.String(cache.ApiCachingBehavior),
		Ttl:                pulumi.Int(int(cache.Ttl)),
		Type:               pulumi.String(cache.Type),
	}
	if cache.AtRestEncryptionEnabled {
		cacheArgs.AtRestEncryptionEnabled = pulumi.Bool(true)
	}
	if cache.TransitEncryptionEnabled {
		cacheArgs.TransitEncryptionEnabled = pulumi.Bool(true)
	}
	return appsync.NewApiCache(ctx, "cache", cacheArgs, pulumi.Provider(provider))
}
