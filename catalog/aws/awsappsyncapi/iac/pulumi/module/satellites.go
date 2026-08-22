package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/appsync"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// satellites creates the arm satellites (types, functions, resolvers,
// channel namespaces), the shared surfaces (API keys, the custom
// domain), and MERGED source-API associations, then exports outputs.
//
// Lifecycle facts the renders depend on:
//   - resolver mutations on one API serialize behind the provider's
//     per-API mutex (bulk changes apply one at a time);
//   - aws_appsync_type ignores format changes on update (a perpetual
//     diff) - replace the entry to change format;
//   - the api key's SECRET is returned only at creation - the output
//     map carries key IDs, never secrets;
//   - the domain association is 1:1 with the domain (its AWS id IS
//     the domain name); create/delete wait up to 60 minutes upstream.
func satellites(ctx *pulumi.Context, locals *Locals, provider *aws.Provider,
	api *createdApi, createdDatasources map[string]*appsync.DataSource) error {

	// In-spec data source names resolve through the created data
	// source (carrying the dependency edge); externally created data
	// sources pass through as literals.
	datasourceName := func(name string) pulumi.StringInput {
		if created, in := createdDatasources[name]; in {
			return created.Name
		}
		return pulumi.String(name)
	}

	datasourceArns := pulumi.StringMap{}
	for name, created := range createdDatasources {
		datasourceArns[name] = created.Arn
	}

	// --- GraphQL satellites -------------------------------------------------

	functionIds := pulumi.StringMap{}
	createdFunctions := map[string]*appsync.Function{}
	// Import-derivation echo map (the type's composite import ID
	// carries its format, which no output of the resource itself
	// echoes).
	typeFormats := pulumi.StringMap{}

	if graphql := locals.Spec.GetGraphql(); graphql != nil {
		for _, entry := range graphql.Types {
			typeFormats[entry.Name] = pulumi.String(entry.Format).ToStringOutput()
			if _, err := appsync.NewType(ctx, fmt.Sprintf("type-%s", entry.Name),
				&appsync.TypeArgs{
					ApiId:      api.ApiId,
					Definition: pulumi.String(entry.Definition),
					Format:     pulumi.String(entry.Format),
				}, pulumi.Provider(provider)); err != nil {
				return errors.Wrapf(err, "create type %s", entry.Name)
			}
		}

		for _, function := range graphql.Functions {
			functionArgs := &appsync.FunctionArgs{
				ApiId:      api.ApiId,
				Name:       pulumi.String(function.Name),
				DataSource: datasourceName(function.DataSourceName),
			}
			if function.Description != "" {
				functionArgs.Description = pulumi.String(function.Description)
			}
			if function.Code != "" {
				functionArgs.Code = pulumi.String(function.Code)
				functionArgs.Runtime = &appsync.FunctionRuntimeArgs{
					// APPSYNC_JS is the only runtime AWS ships - pinned.
					Name:           pulumi.String("APPSYNC_JS"),
					RuntimeVersion: pulumi.String(runtimeVersionOrDefault(function.RuntimeVersion)),
				}
			}
			if function.RequestMappingTemplate != "" {
				functionArgs.RequestMappingTemplate = pulumi.String(function.RequestMappingTemplate)
			}
			if function.ResponseMappingTemplate != "" {
				functionArgs.ResponseMappingTemplate = pulumi.String(function.ResponseMappingTemplate)
			}
			if function.MaxBatchSize > 0 {
				functionArgs.MaxBatchSize = pulumi.Int(int(function.MaxBatchSize))
			}
			if sync := function.SyncConfig; sync != nil {
				syncArgs := &appsync.FunctionSyncConfigArgs{}
				if sync.ConflictDetection != "" {
					syncArgs.ConflictDetection = pulumi.String(sync.ConflictDetection)
				}
				if sync.ConflictHandler != "" {
					syncArgs.ConflictHandler = pulumi.String(sync.ConflictHandler)
				}
				if sync.LambdaConflictHandlerArn.GetValue() != "" {
					syncArgs.LambdaConflictHandlerConfig = &appsync.FunctionSyncConfigLambdaConflictHandlerConfigArgs{
						LambdaConflictHandlerArn: pulumi.String(sync.LambdaConflictHandlerArn.GetValue()),
					}
				}
				functionArgs.SyncConfig = syncArgs
			}

			createdFunction, err := appsync.NewFunction(ctx,
				fmt.Sprintf("function-%s", function.Name),
				functionArgs, pulumi.Provider(provider))
			if err != nil {
				return errors.Wrapf(err, "create function %s", function.Name)
			}
			createdFunctions[function.Name] = createdFunction
			functionIds[function.Name] = createdFunction.FunctionId
		}

		for _, resolver := range graphql.Resolvers {
			resolverKey := fmt.Sprintf("%s.%s", resolver.Type, resolver.Field)
			resolverArgs := &appsync.ResolverArgs{
				ApiId: api.ApiId,
				Type:  pulumi.String(resolver.Type),
				Field: pulumi.String(resolver.Field),
			}

			// UNIT XOR PIPELINE, derived from which arm the entry
			// carries (the spec's CEL wall keeps it exactly one).
			if len(resolver.PipelineFunctions) > 0 {
				resolverArgs.Kind = pulumi.String("PIPELINE")
				pipeline := pulumi.StringArray{}
				for _, functionName := range resolver.PipelineFunctions {
					// Pipeline entries name spec functions; the module
					// joins names to AWS function ids (externally
					// created functions pass their ids through).
					if createdFunction, in := createdFunctions[functionName]; in {
						pipeline = append(pipeline, createdFunction.FunctionId)
					} else {
						pipeline = append(pipeline, pulumi.String(functionName))
					}
				}
				resolverArgs.PipelineConfig = &appsync.ResolverPipelineConfigArgs{
					Functions: pipeline,
				}
			} else {
				resolverArgs.Kind = pulumi.String("UNIT")
				resolverArgs.DataSource = datasourceName(resolver.DataSourceName)
			}

			if resolver.Code != "" {
				resolverArgs.Code = pulumi.String(resolver.Code)
				resolverArgs.Runtime = &appsync.ResolverRuntimeArgs{
					Name:           pulumi.String("APPSYNC_JS"),
					RuntimeVersion: pulumi.String(runtimeVersionOrDefault(resolver.RuntimeVersion)),
				}
			}
			if resolver.RequestTemplate != "" {
				resolverArgs.RequestTemplate = pulumi.String(resolver.RequestTemplate)
			}
			if resolver.ResponseTemplate != "" {
				resolverArgs.ResponseTemplate = pulumi.String(resolver.ResponseTemplate)
			}
			if resolver.MaxBatchSize > 0 {
				resolverArgs.MaxBatchSize = pulumi.Int(int(resolver.MaxBatchSize))
			}
			if caching := resolver.Caching; caching != nil {
				cachingArgs := &appsync.ResolverCachingConfigArgs{}
				if len(caching.CachingKeys) > 0 {
					cachingArgs.CachingKeys = pulumi.ToStringArray(caching.CachingKeys)
				}
				if caching.Ttl > 0 {
					cachingArgs.Ttl = pulumi.Int(int(caching.Ttl))
				}
				resolverArgs.CachingConfig = cachingArgs
			}
			if sync := resolver.SyncConfig; sync != nil {
				syncArgs := &appsync.ResolverSyncConfigArgs{}
				if sync.ConflictDetection != "" {
					syncArgs.ConflictDetection = pulumi.String(sync.ConflictDetection)
				}
				if sync.ConflictHandler != "" {
					syncArgs.ConflictHandler = pulumi.String(sync.ConflictHandler)
				}
				if sync.LambdaConflictHandlerArn.GetValue() != "" {
					syncArgs.LambdaConflictHandlerConfig = &appsync.ResolverSyncConfigLambdaConflictHandlerConfigArgs{
						LambdaConflictHandlerArn: pulumi.String(sync.LambdaConflictHandlerArn.GetValue()),
					}
				}
				resolverArgs.SyncConfig = syncArgs
			}

			if _, err := appsync.NewResolver(ctx, fmt.Sprintf("resolver-%s", resolverKey),
				resolverArgs, pulumi.Provider(provider)); err != nil {
				return errors.Wrapf(err, "create resolver %s", resolverKey)
			}
		}
	}

	// --- Events satellites --------------------------------------------------

	channelNamespaceArns := pulumi.StringMap{}
	if events := locals.Spec.GetEvents(); events != nil {
		for _, namespace := range events.ChannelNamespaces {
			namespaceArgs := &appsync.ChannelNamespaceArgs{
				ApiId: api.ApiId,
				Name:  pulumi.String(namespace.Name),
				Tags:  pulumi.ToStringMap(locals.AwsTags),
			}
			if namespace.CodeHandlers != "" {
				namespaceArgs.CodeHandlers = pulumi.String(namespace.CodeHandlers)
			}
			if len(namespace.PublishAuthModes) > 0 {
				modes := appsync.ChannelNamespacePublishAuthModeArray{}
				for _, mode := range namespace.PublishAuthModes {
					modes = append(modes, &appsync.ChannelNamespacePublishAuthModeArgs{
						AuthType: pulumi.String(mode),
					})
				}
				namespaceArgs.PublishAuthModes = modes
			}
			if len(namespace.SubscribeAuthModes) > 0 {
				modes := appsync.ChannelNamespaceSubscribeAuthModeArray{}
				for _, mode := range namespace.SubscribeAuthModes {
					modes = append(modes, &appsync.ChannelNamespaceSubscribeAuthModeArgs{
						AuthType: pulumi.String(mode),
					})
				}
				namespaceArgs.SubscribeAuthModes = modes
			}
			if handlers := namespace.HandlerConfigs; handlers != nil {
				handlerArgs := &appsync.ChannelNamespaceHandlerConfigsArgs{}
				if onPublish := handlers.OnPublish; onPublish != nil {
					integrationArgs := &appsync.ChannelNamespaceHandlerConfigsOnPublishIntegrationArgs{
						DataSourceName: datasourceName(onPublish.DataSourceName),
					}
					if onPublish.LambdaInvokeType != "" {
						integrationArgs.LambdaConfig = &appsync.ChannelNamespaceHandlerConfigsOnPublishIntegrationLambdaConfigArgs{
							InvokeType: pulumi.String(onPublish.LambdaInvokeType),
						}
					}
					handlerArgs.OnPublish = &appsync.ChannelNamespaceHandlerConfigsOnPublishArgs{
						Behavior:    pulumi.String(onPublish.Behavior),
						Integration: integrationArgs,
					}
				}
				if onSubscribe := handlers.OnSubscribe; onSubscribe != nil {
					integrationArgs := &appsync.ChannelNamespaceHandlerConfigsOnSubscribeIntegrationArgs{
						DataSourceName: datasourceName(onSubscribe.DataSourceName),
					}
					if onSubscribe.LambdaInvokeType != "" {
						integrationArgs.LambdaConfig = &appsync.ChannelNamespaceHandlerConfigsOnSubscribeIntegrationLambdaConfigArgs{
							InvokeType: pulumi.String(onSubscribe.LambdaInvokeType),
						}
					}
					handlerArgs.OnSubscribe = &appsync.ChannelNamespaceHandlerConfigsOnSubscribeArgs{
						Behavior:    pulumi.String(onSubscribe.Behavior),
						Integration: integrationArgs,
					}
				}
				namespaceArgs.HandlerConfigs = handlerArgs
			}

			createdNamespace, err := appsync.NewChannelNamespace(ctx,
				fmt.Sprintf("channel-namespace-%s", namespace.Name),
				namespaceArgs, pulumi.Provider(provider))
			if err != nil {
				return errors.Wrapf(err, "create channel namespace %s", namespace.Name)
			}
			channelNamespaceArns[namespace.Name] = createdNamespace.ChannelNamespaceArn
		}
	}

	// --- Shared surfaces ----------------------------------------------------

	apiKeyIds := pulumi.StringMap{}
	for _, apiKey := range locals.Spec.ApiKeys {
		apiKeyArgs := &appsync.ApiKeyArgs{
			ApiId: api.ApiId,
		}
		if apiKey.Description != "" {
			apiKeyArgs.Description = pulumi.String(apiKey.Description)
		}
		if apiKey.Expires != "" {
			apiKeyArgs.Expires = pulumi.String(apiKey.Expires)
		}
		createdApiKey, err := appsync.NewApiKey(ctx,
			fmt.Sprintf("api-key-%s", apiKey.Name),
			apiKeyArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create api key %s", apiKey.Name)
		}
		// Key IDs, never secrets: AWS returns a key's secret only at
		// creation, and the provider's key attribute degrades to the
		// ID on refresh.
		apiKeyIds[apiKey.Name] = createdApiKey.ApiKeyId
	}

	appsyncDomainName := pulumi.String("").ToStringOutput()
	domainHostedZoneId := pulumi.String("").ToStringOutput()
	if customDomain := locals.Spec.CustomDomain; customDomain != nil {
		domainArgs := &appsync.DomainNameArgs{
			DomainName: pulumi.String(customDomain.DomainName),
			// Must be in us-east-1 (the CloudFront class; the spec's
			// CEL checks literals).
			CertificateArn: pulumi.String(customDomain.CertificateArn.GetValue()),
		}
		if customDomain.Description != "" {
			domainArgs.Description = pulumi.String(customDomain.Description)
		}
		createdDomain, err := appsync.NewDomainName(ctx, "domain", domainArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create domain name")
		}

		// The association is 1:1 with the domain (its AWS id IS the
		// domain name); create/delete wait up to 60 minutes upstream.
		if _, err := appsync.NewDomainNameApiAssociation(ctx, "domain-association",
			&appsync.DomainNameApiAssociationArgs{
				DomainName: createdDomain.DomainName,
				ApiId:      api.ApiId,
			}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "associate domain")
		}

		appsyncDomainName = createdDomain.AppsyncDomainName
		domainHostedZoneId = createdDomain.HostedZoneId
	}

	// --- MERGED source-API associations --------------------------------------

	sourceApiAssociationIds := pulumi.StringMap{}
	if locals.IsMerged {
		for _, sourceApi := range locals.Spec.Graphql.Merged.SourceApis {
			associationArgs := &appsync.SourceApiAssociationArgs{
				MergedApiId: api.ApiId,
				SourceApiId: pulumi.String(sourceApi.SourceApiId.GetValue()),
			}
			if sourceApi.Description != "" {
				associationArgs.Description = pulumi.String(sourceApi.Description)
			}
			if sourceApi.MergeType != "" {
				associationArgs.SourceApiAssociationConfigs = appsync.SourceApiAssociationSourceApiAssociationConfigArray{
					&appsync.SourceApiAssociationSourceApiAssociationConfigArgs{
						MergeType: pulumi.String(sourceApi.MergeType),
					},
				}
			}
			createdAssociation, err := appsync.NewSourceApiAssociation(ctx,
				fmt.Sprintf("source-api-%s", sourceApi.Name),
				associationArgs, pulumi.Provider(provider))
			if err != nil {
				return errors.Wrapf(err, "associate source api %s", sourceApi.Name)
			}
			sourceApiAssociationIds[sourceApi.Name] = createdAssociation.AssociationId
		}
	}

	ctx.Export(OpApiId, api.ApiId)
	ctx.Export(OpApiArn, api.ApiArn)
	ctx.Export(OpGraphqlUrl, api.GraphqlUrl)
	ctx.Export(OpRealtimeUrl, api.RealtimeUrl)
	ctx.Export(OpEventsHttpEndpoint, api.EventsHttp)
	ctx.Export(OpEventsRealtimeEndpoint, api.EventsRealtime)
	ctx.Export(OpAppsyncDomainName, appsyncDomainName)
	ctx.Export(OpDomainHostedZoneId, domainHostedZoneId)
	ctx.Export(OpDatasourceArns, datasourceArns)
	ctx.Export(OpFunctionIds, functionIds)
	ctx.Export(OpApiKeyIds, apiKeyIds)
	ctx.Export(OpChannelNamespaceArns, channelNamespaceArns)
	ctx.Export(OpSourceApiAssociationIds, sourceApiAssociationIds)
	ctx.Export(OpTypeFormats, typeFormats)
	return nil
}

// runtimeVersionOrDefault resolves the APPSYNC_JS runtime version
// ("1.0.0" - the only version AWS ships - when unset).
func runtimeVersionOrDefault(version string) string {
	if version == "" {
		return "1.0.0"
	}
	return version
}
