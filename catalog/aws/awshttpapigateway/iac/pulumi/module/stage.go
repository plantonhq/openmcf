package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// stage creates the API's deployment stage. It runs AFTER the routes exist
// because per-route settings reference routes by key -- AWS rejects settings
// for unknown route keys.
func stage(
	ctx *pulumi.Context,
	locals *Locals,
	createdApi *apigatewayv2.Api,
	createdRoutes []pulumi.Resource,
	provider *aws.Provider,
) error {
	spec := locals.Spec

	// Stage defaults: name "$default" when unset, and auto-deploy ON when the
	// spec does not say otherwise -- a declarative spec should be
	// self-applying, so only an explicit auto_deploy=false turns it off. The
	// Terraform module implements the identical presence rule.
	stageName := "$default"
	autoDeploy := true
	if spec.Stage != nil && spec.Stage.Name != "" {
		stageName = spec.Stage.Name
	}
	if spec.Stage != nil && spec.Stage.AutoDeploy != nil {
		autoDeploy = *spec.Stage.AutoDeploy
	}

	resourceName := locals.ApiName + "-stage"

	args := &apigatewayv2.StageArgs{
		ApiId:      createdApi.ID(),
		Name:       pulumi.String(stageName),
		AutoDeploy: pulumi.BoolPtr(autoDeploy),
		Tags:       pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Stage != nil && spec.Stage.Description != "" {
		args.Description = pulumi.StringPtr(spec.Stage.Description)
	}

	// Access logging streams request records to CloudWatch Logs in the
	// user-supplied format. The log group is referenced, never created here
	// -- log groups are their own composable resource (AwsCloudwatchLogGroup).
	if spec.Stage != nil && spec.Stage.AccessLog != nil {
		accessLog := spec.Stage.AccessLog
		args.AccessLogSettings = &apigatewayv2.StageAccessLogSettingsArgs{
			DestinationArn: pulumi.String(accessLog.DestinationArn.GetValue()),
			Format:         pulumi.String(accessLog.Format),
		}
	}

	// Stage-wide defaults. Rendered when any default is meaningful: a
	// throttle limit or detailed metrics. Zero-valued limits are never sent
	// -- they would otherwise clamp the stage to zero requests/second.
	// (data_trace_enabled and logging_level are WebSocket-only knobs and are
	// deliberately not modeled for HTTP APIs.)
	if spec.Stage != nil {
		burst := int32(0)
		rate := float64(0)
		if spec.Stage.DefaultThrottle != nil {
			burst = spec.Stage.DefaultThrottle.BurstLimit
			rate = spec.Stage.DefaultThrottle.RateLimit
		}
		if burst > 0 || rate > 0 || spec.Stage.DetailedMetricsEnabled {
			defaults := &apigatewayv2.StageDefaultRouteSettingsArgs{
				DetailedMetricsEnabled: pulumi.BoolPtr(spec.Stage.DetailedMetricsEnabled),
			}
			if burst > 0 {
				defaults.ThrottlingBurstLimit = pulumi.IntPtr(int(burst))
			}
			if rate > 0 {
				defaults.ThrottlingRateLimit = pulumi.Float64Ptr(rate)
			}
			args.DefaultRouteSettings = defaults
		}

		// Per-route overrides. Zero-valued limits inherit the stage default
		// -- only real overrides are sent.
		if len(spec.Stage.RouteSettings) > 0 {
			routeSettings := make(apigatewayv2.StageRouteSettingArray, 0, len(spec.Stage.RouteSettings))
			for _, rs := range spec.Stage.RouteSettings {
				rsArgs := &apigatewayv2.StageRouteSettingArgs{
					RouteKey:               pulumi.String(rs.RouteKey),
					DetailedMetricsEnabled: pulumi.BoolPtr(rs.DetailedMetricsEnabled),
				}
				if rs.ThrottlingBurstLimit > 0 {
					rsArgs.ThrottlingBurstLimit = pulumi.IntPtr(int(rs.ThrottlingBurstLimit))
				}
				if rs.ThrottlingRateLimit > 0 {
					rsArgs.ThrottlingRateLimit = pulumi.Float64Ptr(rs.ThrottlingRateLimit)
				}
				routeSettings = append(routeSettings, rsArgs)
			}
			args.RouteSettings = routeSettings
		}

		if len(spec.Stage.StageVariables) > 0 {
			args.StageVariables = pulumi.ToStringMap(spec.Stage.StageVariables)
		}
	}

	createdStage, err := apigatewayv2.NewStage(ctx, resourceName, args,
		pulumi.Provider(provider), pulumi.DependsOn(createdRoutes))
	if err != nil {
		return errors.Wrap(err, "failed to create API stage")
	}

	// Export stage outputs
	ctx.Export(OpStageInvokeUrl, createdStage.InvokeUrl)
	ctx.Export(OpStageName, pulumi.String(stageName))

	return nil
}
