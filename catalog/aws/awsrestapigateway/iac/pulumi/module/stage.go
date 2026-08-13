package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// stage creates the hash-triggered deployment, the optional generated
// client certificate, the stage, and its per-method settings.
//
// Lifecycle facts the renders below depend on:
//   - REST APIs deploy by explicit snapshot: the deployment's trigger
//     carries the definition hash from locals, so ANY definition
//     change rolls a new deployment (create-before-destroy - the stage
//     repoints, then the old deployment deletes);
//   - deleting a deployment an active canary references fails
//     upstream once and retries after the stage detaches - not a
//     concern here (canary is workflow surface, not modeled);
//   - enabling or resizing the stage cache waits up to 90 minutes
//     upstream while AWS provisions the cluster;
//   - method settings are PATCHes on the stage (a view, not a real
//     object): create and update both PATCH, delete PATCHes the
//     overrides away.
func stage(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, api *apigateway.RestApi, tree *createdTree) error {
	spec := locals.Spec

	deploymentArgs := &apigateway.DeploymentArgs{
		RestApi: api.ID(),
		Triggers: pulumi.StringMap{
			"redeployment": pulumi.String(locals.DeploymentTriggerHash),
		},
	}
	deployment, err := apigateway.NewDeployment(ctx, "deployment", deploymentArgs,
		pulumi.Provider(provider),
		// The snapshot must capture the complete definition.
		pulumi.DependsOn(tree.definitionResources),
		// Replacement order: create the new deployment, repoint the
		// stage, then delete the old one.
		pulumi.DeleteBeforeReplace(false))
	if err != nil {
		return errors.Wrap(err, "create deployment")
	}
	ctx.Export(OpDeploymentId, deployment.ID())

	stageConfig := spec.Stage
	if stageConfig == nil {
		stageConfig = &awsRestApiGatewayStageDefaults
	}

	args := &apigateway.StageArgs{
		RestApi:    api.ID(),
		Deployment: deployment.ID(),
		StageName:  pulumi.String(locals.StageName),
		Tags:       pulumi.ToStringMap(locals.AwsTags),
	}
	if stageConfig.Description != "" {
		args.Description = pulumi.String(stageConfig.Description)
	}
	if len(stageConfig.StageVariables) > 0 {
		args.Variables = pulumi.ToStringMap(stageConfig.StageVariables)
	}
	if stageConfig.XrayTracingEnabled {
		args.XrayTracingEnabled = pulumi.Bool(true)
	}
	if stageConfig.CacheCluster != nil && stageConfig.CacheCluster.Enabled {
		args.CacheClusterEnabled = pulumi.Bool(true)
		args.CacheClusterSize = pulumi.String(stageConfig.CacheCluster.Size)
	}

	// The TLS client certificate presented to HTTP backends: generated
	// with this API, or an existing one by ID.
	if stageConfig.ClientCertificate != nil {
		if stageConfig.ClientCertificate.Generate {
			certificateArgs := &apigateway.ClientCertificateArgs{
				Tags: pulumi.ToStringMap(locals.AwsTags),
			}
			if stageConfig.ClientCertificate.Description != "" {
				certificateArgs.Description = pulumi.String(stageConfig.ClientCertificate.Description)
			}
			certificate, err := apigateway.NewClientCertificate(ctx, "client-certificate", certificateArgs, pulumi.Provider(provider))
			if err != nil {
				return errors.Wrap(err, "create client certificate")
			}
			args.ClientCertificateId = certificate.ID()
			ctx.Export(OpClientCertificateId, certificate.ID())
			ctx.Export(OpClientCertificatePem, certificate.PemEncodedCertificate)
		} else {
			args.ClientCertificateId = pulumi.String(stageConfig.ClientCertificate.ExistingCertificateId)
		}
	}

	if stageConfig.AccessLog != nil {
		args.AccessLogSettings = &apigateway.StageAccessLogSettingsArgs{
			DestinationArn: pulumi.String(stageConfig.AccessLog.DestinationArn.GetValue()),
			Format:         pulumi.String(stageConfig.AccessLog.Format),
		}
	}
	if stageConfig.DocumentationVersion != "" {
		args.DocumentationVersion = pulumi.String(stageConfig.DocumentationVersion)
	}

	created, err := apigateway.NewStage(ctx, "stage", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create stage")
	}

	// Per-method overrides of logging/metrics/throttling/caching.
	for _, m := range stageConfig.MethodSettings {
		settings := &apigateway.MethodSettingsSettingsArgs{}
		if m.MetricsEnabled != nil {
			settings.MetricsEnabled = pulumi.Bool(*m.MetricsEnabled)
		}
		if m.LoggingLevel != "" {
			settings.LoggingLevel = pulumi.String(m.LoggingLevel)
		}
		if m.DataTraceEnabled != nil {
			settings.DataTraceEnabled = pulumi.Bool(*m.DataTraceEnabled)
		}
		if m.ThrottlingBurstLimit != nil {
			settings.ThrottlingBurstLimit = pulumi.Int(int(*m.ThrottlingBurstLimit))
		}
		if m.ThrottlingRateLimit != nil {
			settings.ThrottlingRateLimit = pulumi.Float64(*m.ThrottlingRateLimit)
		}
		if m.CachingEnabled != nil {
			settings.CachingEnabled = pulumi.Bool(*m.CachingEnabled)
		}
		if m.CacheTtlInSeconds != nil {
			settings.CacheTtlInSeconds = pulumi.Int(int(*m.CacheTtlInSeconds))
		}
		if m.CacheDataEncrypted != nil {
			settings.CacheDataEncrypted = pulumi.Bool(*m.CacheDataEncrypted)
		}
		if m.RequireAuthorizationForCacheControl != nil {
			settings.RequireAuthorizationForCacheControl = pulumi.Bool(*m.RequireAuthorizationForCacheControl)
		}
		if m.UnauthorizedCacheControlHeaderStrategy != "" {
			settings.UnauthorizedCacheControlHeaderStrategy = pulumi.String(m.UnauthorizedCacheControlHeaderStrategy)
		}
		if _, err := apigateway.NewMethodSettings(ctx, "method-settings-"+m.MethodPath, &apigateway.MethodSettingsArgs{
			RestApi:    api.ID(),
			StageName:  created.StageName,
			MethodPath: pulumi.String(m.MethodPath),
			Settings:   settings,
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrapf(err, "create method settings %q", m.MethodPath)
		}
	}

	ctx.Export(OpStageName, created.StageName)
	ctx.Export(OpStageArn, created.Arn)
	ctx.Export(OpInvokeUrl, created.InvokeUrl)
	return nil
}
