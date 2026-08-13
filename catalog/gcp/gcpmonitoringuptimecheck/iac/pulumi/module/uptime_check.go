package module

import (
	"github.com/pkg/errors"
	gcpmonitoringuptimecheckv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpmonitoringuptimecheck/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/monitoring"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// uptimeCheck provisions the Cloud Monitoring uptime check — the probe
// Google runs against the configured target from multiple regions.
//
// Exactly one target arm (monitored_resource | resource_group |
// synthetic_monitor) and — except for synthetic monitors — exactly one
// check arm (http_check | tcp_check) are present; the proto CELs enforce
// the shape before the module runs, so the expand code here maps arms
// unconditionally when present.
//
// Enum-like strings and durations with API defaults (period, checker_type,
// request_method, content_type) fall through as unset when empty so GCP
// applies its own defaults rather than receiving empty strings it would
// reject.
func uptimeCheck(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpMonitoringUptimeCheck.Spec

	// Enable the Cloud Monitoring API so a fresh project can host the
	// check. disable_on_destroy stays false (the provider default):
	// tearing down one check must never disable monitoring for everything
	// else in the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("monitoring.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"uptimecheck-monitoring.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable monitoring.googleapis.com api")
	}

	args := &monitoring.UptimeCheckConfigArgs{
		DisplayName: pulumi.String(locals.DisplayName),
		Timeout:     pulumi.String(spec.Timeout),
		UserLabels:  pulumi.ToStringMap(locals.GcpLabels),
	}

	if spec.Period != "" {
		args.Period = pulumi.String(spec.Period)
	}
	if spec.CheckerType != "" {
		args.CheckerType = pulumi.String(spec.CheckerType)
	}
	if len(spec.SelectedRegions) > 0 {
		args.SelectedRegions = pulumi.ToStringArray(spec.SelectedRegions)
	}
	if spec.LogCheckFailures {
		args.LogCheckFailures = pulumi.Bool(true)
	}

	// ── Target arms (exactly one present, per the proto CEL) ──

	if spec.MonitoredResource != nil {
		args.MonitoredResource = &monitoring.UptimeCheckConfigMonitoredResourceArgs{
			Type:   pulumi.String(spec.MonitoredResource.Type),
			Labels: pulumi.ToStringMap(spec.MonitoredResource.Labels),
		}
	}

	if spec.ResourceGroup != nil {
		resourceGroupArgs := &monitoring.UptimeCheckConfigResourceGroupArgs{}
		if spec.ResourceGroup.GroupId != "" {
			resourceGroupArgs.GroupId = pulumi.StringPtr(spec.ResourceGroup.GroupId)
		}
		if spec.ResourceGroup.ResourceType != "" {
			resourceGroupArgs.ResourceType = pulumi.StringPtr(spec.ResourceGroup.ResourceType)
		}
		args.ResourceGroup = resourceGroupArgs
	}

	if spec.SyntheticMonitor != nil {
		args.SyntheticMonitor = &monitoring.UptimeCheckConfigSyntheticMonitorArgs{
			CloudFunctionV2: &monitoring.UptimeCheckConfigSyntheticMonitorCloudFunctionV2Args{
				Name: pulumi.String(spec.SyntheticMonitor.CloudFunction.GetValue()),
			},
		}
	}

	// ── Check arms ──

	if spec.HttpCheck != nil {
		args.HttpCheck = expandHttpCheck(spec.HttpCheck)
	}

	if spec.TcpCheck != nil {
		tcpArgs := &monitoring.UptimeCheckConfigTcpCheckArgs{
			Port: pulumi.Int(int(spec.TcpCheck.Port)),
		}
		if spec.TcpCheck.PingConfig != nil {
			tcpArgs.PingConfig = &monitoring.UptimeCheckConfigTcpCheckPingConfigArgs{
				PingsCount: pulumi.Int(int(spec.TcpCheck.PingConfig.PingsCount)),
			}
		}
		args.TcpCheck = tcpArgs
	}

	// ── Content matchers ──

	if len(spec.ContentMatchers) > 0 {
		matchers := monitoring.UptimeCheckConfigContentMatcherArray{}
		for _, matcher := range spec.ContentMatchers {
			matcherArgs := &monitoring.UptimeCheckConfigContentMatcherArgs{
				Content: pulumi.String(matcher.Content),
			}
			if matcher.Matcher != "" {
				matcherArgs.Matcher = pulumi.StringPtr(matcher.Matcher)
			}
			if matcher.JsonPathMatcher != nil {
				jsonPathArgs := &monitoring.UptimeCheckConfigContentMatcherJsonPathMatcherArgs{
					JsonPath: pulumi.String(matcher.JsonPathMatcher.JsonPath),
				}
				if matcher.JsonPathMatcher.JsonMatcher != "" {
					jsonPathArgs.JsonMatcher = pulumi.StringPtr(matcher.JsonPathMatcher.JsonMatcher)
				}
				matcherArgs.JsonPathMatcher = jsonPathArgs
			}
			matchers = append(matchers, matcherArgs)
		}
		args.ContentMatchers = matchers
	}

	// Unset defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (empty string would be sent verbatim and
	// rejected).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdUptimeCheck, err := monitoring.NewUptimeCheckConfig(ctx, "uptime-check", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create uptime check config")
	}

	ctx.Export(OpUptimeCheckName, createdUptimeCheck.Name)
	ctx.Export(OpUptimeCheckId, createdUptimeCheck.UptimeCheckId)

	return nil
}

// expandHttpCheck maps the spec's HTTP probe onto the provider block.
// Booleans with API-default false (use_ssl, validate_ssl, mask_headers) are
// sent only when true — an omitted false and an explicit false are
// indistinguishable to this API, and omitting keeps plans minimal (matching
// the Terraform module).
func expandHttpCheck(httpCheck *gcpmonitoringuptimecheckv1alpha1.GcpMonitoringUptimeCheckHttpCheck) *monitoring.UptimeCheckConfigHttpCheckArgs {
	httpArgs := &monitoring.UptimeCheckConfigHttpCheckArgs{}

	if httpCheck.Path != "" {
		httpArgs.Path = pulumi.StringPtr(httpCheck.Path)
	}
	if httpCheck.Port != 0 {
		httpArgs.Port = pulumi.IntPtr(int(httpCheck.Port))
	}
	if httpCheck.RequestMethod != "" {
		httpArgs.RequestMethod = pulumi.StringPtr(httpCheck.RequestMethod)
	}
	if httpCheck.UseSsl {
		httpArgs.UseSsl = pulumi.BoolPtr(true)
	}
	if httpCheck.ValidateSsl {
		httpArgs.ValidateSsl = pulumi.BoolPtr(true)
	}
	if httpCheck.Body != "" {
		httpArgs.Body = pulumi.StringPtr(httpCheck.Body)
	}
	if httpCheck.ContentType != "" {
		httpArgs.ContentType = pulumi.StringPtr(httpCheck.ContentType)
	}
	if httpCheck.CustomContentType != "" {
		httpArgs.CustomContentType = pulumi.StringPtr(httpCheck.CustomContentType)
	}
	if len(httpCheck.Headers) > 0 {
		httpArgs.Headers = pulumi.ToStringMap(httpCheck.Headers)
	}
	if httpCheck.MaskHeaders {
		httpArgs.MaskHeaders = pulumi.BoolPtr(true)
	}

	if httpCheck.AuthInfo != nil {
		httpArgs.AuthInfo = &monitoring.UptimeCheckConfigHttpCheckAuthInfoArgs{
			Username: pulumi.String(httpCheck.AuthInfo.Username),
			Password: pulumi.String(httpCheck.AuthInfo.Password),
		}
	}

	if httpCheck.ServiceAgentAuthentication != nil && httpCheck.ServiceAgentAuthentication.Type != "" {
		httpArgs.ServiceAgentAuthentication = &monitoring.UptimeCheckConfigHttpCheckServiceAgentAuthenticationArgs{
			Type: pulumi.StringPtr(httpCheck.ServiceAgentAuthentication.Type),
		}
	}

	if len(httpCheck.AcceptedResponseStatusCodes) > 0 {
		statusCodes := monitoring.UptimeCheckConfigHttpCheckAcceptedResponseStatusCodeArray{}
		for _, code := range httpCheck.AcceptedResponseStatusCodes {
			codeArgs := &monitoring.UptimeCheckConfigHttpCheckAcceptedResponseStatusCodeArgs{}
			if code.StatusClass != "" {
				codeArgs.StatusClass = pulumi.StringPtr(code.StatusClass)
			}
			if code.StatusValue != 0 {
				codeArgs.StatusValue = pulumi.IntPtr(int(code.StatusValue))
			}
			statusCodes = append(statusCodes, codeArgs)
		}
		httpArgs.AcceptedResponseStatusCodes = statusCodes
	}

	if httpCheck.PingConfig != nil {
		httpArgs.PingConfig = &monitoring.UptimeCheckConfigHttpCheckPingConfigArgs{
			PingsCount: pulumi.Int(int(httpCheck.PingConfig.PingsCount)),
		}
	}

	return httpArgs
}
