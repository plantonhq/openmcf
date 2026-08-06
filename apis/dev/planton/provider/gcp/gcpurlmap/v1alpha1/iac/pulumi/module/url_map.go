package module

import (
	"strconv"

	"github.com/pkg/errors"
	gcpurlmapv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpurlmap/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// urlMap provisions the global Compute Engine URL map — the L7 routing brain
// of a global external Application Load Balancer. Host rules map request Host
// headers to named path matchers; path matchers evaluate route_rules (priority-
// ordered, rich matching) then path_rules (longest prefix), then their own
// default; anything unmatched falls through to the URL map's top-level default.
//
// name and project are immutable (ForceNew): changing either destroys and
// recreates the map, briefly breaking every target proxy referencing the old
// self_link. Routing tables, header actions, and tests update in place.
//
// Cross-field exclusivity (exactly one default target, path_rules XOR
// route_rules, redirect vs service vs route_action) is enforced by the spec's
// CEL rules before deploy — no defensive logic lives here. route_action maps
// only weighted_backend_services and url_rewrite; mesh-advanced sub-policies
// (timeout, retry, cors, fault, mirror) are a deliberate coverage boundary.
func urlMap(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpUrlMap.Spec

	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("compute.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"urlmap-compute.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable compute.googleapis.com api")
	}

	args := &compute.URLMapArgs{
		Name: pulumi.String(locals.UrlMapName),
	}

	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	if spec.DefaultService.GetValue() != "" {
		args.DefaultService = pulumi.String(spec.DefaultService.GetValue())
	}
	if spec.DefaultUrlRedirect != nil {
		args.DefaultUrlRedirect = topLevelUrlRedirect(spec.DefaultUrlRedirect)
	}
	if spec.DefaultRouteAction != nil {
		args.DefaultRouteAction = topLevelRouteAction(spec.DefaultRouteAction)
	}
	if spec.DefaultCustomErrorResponsePolicy != nil {
		args.DefaultCustomErrorResponsePolicy = topLevelCustomErrorPolicy(spec.DefaultCustomErrorResponsePolicy)
	}
	if spec.HeaderAction != nil {
		args.HeaderAction = topLevelHeaderAction(spec.HeaderAction)
	}
	if len(spec.HostRules) > 0 {
		args.HostRules = buildHostRules(spec.HostRules)
	}
	if len(spec.PathMatchers) > 0 {
		args.PathMatchers = buildPathMatchers(spec.PathMatchers)
	}
	if len(spec.Tests) > 0 {
		args.Tests = buildTests(spec.Tests)
	}

	createdUrlMap, err := compute.NewURLMap(ctx, "url-map", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create url map")
	}

	ctx.Export(OpSelfLink, createdUrlMap.SelfLink)
	ctx.Export(OpUrlMapName, createdUrlMap.Name)
	ctx.Export(OpMapId, createdUrlMap.MapId.ApplyT(func(id int) string {
		return strconv.Itoa(id)
	}).(pulumi.StringOutput))
	ctx.Export(OpFingerprint, createdUrlMap.Fingerprint)

	return nil
}

func topLevelUrlRedirect(r *gcpurlmapv1alpha1.GcpUrlMapUrlRedirect) *compute.URLMapDefaultUrlRedirectArgs {
	return &compute.URLMapDefaultUrlRedirectArgs{
		HostRedirect:         emptyAsNilString(r.HostRedirect),
		HttpsRedirect:        pulumi.Bool(r.HttpsRedirect),
		PathRedirect:         emptyAsNilString(r.PathRedirect),
		PrefixRedirect:       emptyAsNilString(r.PrefixRedirect),
		RedirectResponseCode: emptyAsNilString(r.RedirectResponseCode),
		StripQuery:           pulumi.Bool(r.StripQuery),
	}
}

func topLevelRouteAction(a *gcpurlmapv1alpha1.GcpUrlMapRouteAction) *compute.URLMapDefaultRouteActionArgs {
	if a == nil {
		return nil
	}
	out := &compute.URLMapDefaultRouteActionArgs{}
	if len(a.WeightedBackendServices) > 0 {
		out.WeightedBackendServices = defaultWeightedBackends(a.WeightedBackendServices)
	}
	if a.UrlRewrite != nil {
		out.UrlRewrite = defaultUrlRewrite(a.UrlRewrite)
	}
	return out
}

func defaultWeightedBackends(backends []*gcpurlmapv1alpha1.GcpUrlMapWeightedBackendService) compute.URLMapDefaultRouteActionWeightedBackendServiceArray {
	result := compute.URLMapDefaultRouteActionWeightedBackendServiceArray{}
	for _, backend := range backends {
		result = append(result, &compute.URLMapDefaultRouteActionWeightedBackendServiceArgs{
			BackendService: pulumi.String(backend.BackendService.GetValue()),
			Weight:         pulumi.Int(int(backend.Weight)),
		})
	}
	return result
}

func defaultUrlRewrite(r *gcpurlmapv1alpha1.GcpUrlMapUrlRewrite) *compute.URLMapDefaultRouteActionUrlRewriteArgs {
	return &compute.URLMapDefaultRouteActionUrlRewriteArgs{
		HostRewrite:       emptyAsNilString(r.HostRewrite),
		PathPrefixRewrite: emptyAsNilString(r.PathPrefixRewrite),
	}
}

func topLevelCustomErrorPolicy(p *gcpurlmapv1alpha1.GcpUrlMapCustomErrorResponsePolicy) *compute.URLMapDefaultCustomErrorResponsePolicyArgs {
	return &compute.URLMapDefaultCustomErrorResponsePolicyArgs{
		ErrorService:       emptyAsNilString(p.ErrorService.GetValue()),
		ErrorResponseRules: topLevelErrorRules(p.ErrorResponseRules),
	}
}

func topLevelErrorRules(rules []*gcpurlmapv1alpha1.GcpUrlMapCustomErrorResponseRule) compute.URLMapDefaultCustomErrorResponsePolicyErrorResponseRuleArray {
	result := compute.URLMapDefaultCustomErrorResponsePolicyErrorResponseRuleArray{}
	for _, rule := range rules {
		codes := pulumi.StringArray{}
		for _, code := range rule.MatchResponseCodes {
			codes = append(codes, pulumi.String(code))
		}
		args := &compute.URLMapDefaultCustomErrorResponsePolicyErrorResponseRuleArgs{
			MatchResponseCodes: codes,
			Path:               pulumi.String(rule.Path),
		}
		if rule.OverrideResponseCode != 0 {
			args.OverrideResponseCode = pulumi.Int(int(rule.OverrideResponseCode))
		}
		result = append(result, args)
	}
	return result
}

func topLevelHeaderAction(h *gcpurlmapv1alpha1.GcpUrlMapHeaderAction) *compute.URLMapHeaderActionArgs {
	return &compute.URLMapHeaderActionArgs{
		RequestHeadersToAdds:     topLevelRequestHeadersToAdd(h.RequestHeadersToAdd),
		RequestHeadersToRemoves:  stringArrayOrNil(h.RequestHeadersToRemove),
		ResponseHeadersToAdds:    topLevelResponseHeadersToAdd(h.ResponseHeadersToAdd),
		ResponseHeadersToRemoves: stringArrayOrNil(h.ResponseHeadersToRemove),
	}
}

func topLevelRequestHeadersToAdd(headers []*gcpurlmapv1alpha1.GcpUrlMapHeaderValue) compute.URLMapHeaderActionRequestHeadersToAddArray {
	result := compute.URLMapHeaderActionRequestHeadersToAddArray{}
	for _, header := range headers {
		result = append(result, &compute.URLMapHeaderActionRequestHeadersToAddArgs{
			HeaderName:  pulumi.String(header.HeaderName),
			HeaderValue: pulumi.String(header.HeaderValue),
			Replace:     pulumi.Bool(header.Replace),
		})
	}
	return result
}

func topLevelResponseHeadersToAdd(headers []*gcpurlmapv1alpha1.GcpUrlMapHeaderValue) compute.URLMapHeaderActionResponseHeadersToAddArray {
	result := compute.URLMapHeaderActionResponseHeadersToAddArray{}
	for _, header := range headers {
		result = append(result, &compute.URLMapHeaderActionResponseHeadersToAddArgs{
			HeaderName:  pulumi.String(header.HeaderName),
			HeaderValue: pulumi.String(header.HeaderValue),
			Replace:     pulumi.Bool(header.Replace),
		})
	}
	return result
}

func buildHostRules(rules []*gcpurlmapv1alpha1.GcpUrlMapHostRule) compute.URLMapHostRuleArray {
	result := compute.URLMapHostRuleArray{}
	for _, rule := range rules {
		hosts := pulumi.StringArray{}
		for _, host := range rule.Hosts {
			hosts = append(hosts, pulumi.String(host))
		}
		args := &compute.URLMapHostRuleArgs{
			Hosts:       hosts,
			PathMatcher: pulumi.String(rule.PathMatcher),
		}
		if rule.Description != "" {
			args.Description = pulumi.String(rule.Description)
		}
		result = append(result, args)
	}
	return result
}

func buildPathMatchers(matchers []*gcpurlmapv1alpha1.GcpUrlMapPathMatcher) compute.URLMapPathMatcherArray {
	result := compute.URLMapPathMatcherArray{}
	for _, matcher := range matchers {
		args := &compute.URLMapPathMatcherArgs{
			Name: pulumi.String(matcher.Name),
		}
		if matcher.Description != "" {
			args.Description = pulumi.String(matcher.Description)
		}
		if matcher.DefaultService.GetValue() != "" {
			args.DefaultService = pulumi.String(matcher.DefaultService.GetValue())
		}
		if matcher.DefaultUrlRedirect != nil {
			args.DefaultUrlRedirect = pathMatcherUrlRedirect(matcher.DefaultUrlRedirect)
		}
		if matcher.DefaultRouteAction != nil {
			args.DefaultRouteAction = pathMatcherRouteAction(matcher.DefaultRouteAction)
		}
		if matcher.DefaultCustomErrorResponsePolicy != nil {
			args.DefaultCustomErrorResponsePolicy = pathMatcherCustomErrorPolicy(matcher.DefaultCustomErrorResponsePolicy)
		}
		if matcher.HeaderAction != nil {
			args.HeaderAction = pathMatcherHeaderAction(matcher.HeaderAction)
		}
		if len(matcher.PathRules) > 0 {
			args.PathRules = buildPathRules(matcher.PathRules)
		}
		if len(matcher.RouteRules) > 0 {
			args.RouteRules = buildRouteRules(matcher.RouteRules)
		}
		result = append(result, args)
	}
	return result
}

func pathMatcherUrlRedirect(r *gcpurlmapv1alpha1.GcpUrlMapUrlRedirect) *compute.URLMapPathMatcherDefaultUrlRedirectArgs {
	return &compute.URLMapPathMatcherDefaultUrlRedirectArgs{
		HostRedirect:         emptyAsNilString(r.HostRedirect),
		HttpsRedirect:        pulumi.Bool(r.HttpsRedirect),
		PathRedirect:         emptyAsNilString(r.PathRedirect),
		PrefixRedirect:       emptyAsNilString(r.PrefixRedirect),
		RedirectResponseCode: emptyAsNilString(r.RedirectResponseCode),
		StripQuery:           pulumi.Bool(r.StripQuery),
	}
}

func pathMatcherRouteAction(a *gcpurlmapv1alpha1.GcpUrlMapRouteAction) *compute.URLMapPathMatcherDefaultRouteActionArgs {
	out := &compute.URLMapPathMatcherDefaultRouteActionArgs{}
	if len(a.WeightedBackendServices) > 0 {
		out.WeightedBackendServices = pathMatcherDefaultWeightedBackends(a.WeightedBackendServices)
	}
	if a.UrlRewrite != nil {
		out.UrlRewrite = pathMatcherDefaultUrlRewrite(a.UrlRewrite)
	}
	return out
}

func pathMatcherDefaultWeightedBackends(backends []*gcpurlmapv1alpha1.GcpUrlMapWeightedBackendService) compute.URLMapPathMatcherDefaultRouteActionWeightedBackendServiceArray {
	result := compute.URLMapPathMatcherDefaultRouteActionWeightedBackendServiceArray{}
	for _, backend := range backends {
		result = append(result, &compute.URLMapPathMatcherDefaultRouteActionWeightedBackendServiceArgs{
			BackendService: pulumi.String(backend.BackendService.GetValue()),
			Weight:         pulumi.Int(int(backend.Weight)),
		})
	}
	return result
}

func pathMatcherDefaultUrlRewrite(r *gcpurlmapv1alpha1.GcpUrlMapUrlRewrite) *compute.URLMapPathMatcherDefaultRouteActionUrlRewriteArgs {
	return &compute.URLMapPathMatcherDefaultRouteActionUrlRewriteArgs{
		HostRewrite:       emptyAsNilString(r.HostRewrite),
		PathPrefixRewrite: emptyAsNilString(r.PathPrefixRewrite),
	}
}

func pathMatcherCustomErrorPolicy(p *gcpurlmapv1alpha1.GcpUrlMapCustomErrorResponsePolicy) *compute.URLMapPathMatcherDefaultCustomErrorResponsePolicyArgs {
	return &compute.URLMapPathMatcherDefaultCustomErrorResponsePolicyArgs{
		ErrorService:       emptyAsNilString(p.ErrorService.GetValue()),
		ErrorResponseRules: pathMatcherErrorRules(p.ErrorResponseRules),
	}
}

func pathMatcherErrorRules(rules []*gcpurlmapv1alpha1.GcpUrlMapCustomErrorResponseRule) compute.URLMapPathMatcherDefaultCustomErrorResponsePolicyErrorResponseRuleArray {
	result := compute.URLMapPathMatcherDefaultCustomErrorResponsePolicyErrorResponseRuleArray{}
	for _, rule := range rules {
		codes := pulumi.StringArray{}
		for _, code := range rule.MatchResponseCodes {
			codes = append(codes, pulumi.String(code))
		}
		args := &compute.URLMapPathMatcherDefaultCustomErrorResponsePolicyErrorResponseRuleArgs{
			MatchResponseCodes: codes,
			Path:               pulumi.String(rule.Path),
		}
		if rule.OverrideResponseCode != 0 {
			args.OverrideResponseCode = pulumi.Int(int(rule.OverrideResponseCode))
		}
		result = append(result, args)
	}
	return result
}

func pathMatcherHeaderAction(h *gcpurlmapv1alpha1.GcpUrlMapHeaderAction) *compute.URLMapPathMatcherHeaderActionArgs {
	return &compute.URLMapPathMatcherHeaderActionArgs{
		RequestHeadersToAdds:     pathMatcherRequestHeadersToAdd(h.RequestHeadersToAdd),
		RequestHeadersToRemoves:  stringArrayOrNil(h.RequestHeadersToRemove),
		ResponseHeadersToAdds:    pathMatcherResponseHeadersToAdd(h.ResponseHeadersToAdd),
		ResponseHeadersToRemoves: stringArrayOrNil(h.ResponseHeadersToRemove),
	}
}

func pathMatcherRequestHeadersToAdd(headers []*gcpurlmapv1alpha1.GcpUrlMapHeaderValue) compute.URLMapPathMatcherHeaderActionRequestHeadersToAddArray {
	result := compute.URLMapPathMatcherHeaderActionRequestHeadersToAddArray{}
	for _, header := range headers {
		result = append(result, &compute.URLMapPathMatcherHeaderActionRequestHeadersToAddArgs{
			HeaderName:  pulumi.String(header.HeaderName),
			HeaderValue: pulumi.String(header.HeaderValue),
			Replace:     pulumi.Bool(header.Replace),
		})
	}
	return result
}

func pathMatcherResponseHeadersToAdd(headers []*gcpurlmapv1alpha1.GcpUrlMapHeaderValue) compute.URLMapPathMatcherHeaderActionResponseHeadersToAddArray {
	result := compute.URLMapPathMatcherHeaderActionResponseHeadersToAddArray{}
	for _, header := range headers {
		result = append(result, &compute.URLMapPathMatcherHeaderActionResponseHeadersToAddArgs{
			HeaderName:  pulumi.String(header.HeaderName),
			HeaderValue: pulumi.String(header.HeaderValue),
			Replace:     pulumi.Bool(header.Replace),
		})
	}
	return result
}

func buildPathRules(rules []*gcpurlmapv1alpha1.GcpUrlMapPathRule) compute.URLMapPathMatcherPathRuleArray {
	result := compute.URLMapPathMatcherPathRuleArray{}
	for _, rule := range rules {
		paths := pulumi.StringArray{}
		for _, path := range rule.Paths {
			paths = append(paths, pulumi.String(path))
		}
		args := &compute.URLMapPathMatcherPathRuleArgs{
			Paths: paths,
		}
		if rule.Service.GetValue() != "" {
			args.Service = pulumi.String(rule.Service.GetValue())
		}
		if rule.UrlRedirect != nil {
			args.UrlRedirect = pathRuleUrlRedirect(rule.UrlRedirect)
		}
		if rule.RouteAction != nil {
			args.RouteAction = pathRuleRouteAction(rule.RouteAction)
		}
		if rule.CustomErrorResponsePolicy != nil {
			args.CustomErrorResponsePolicy = pathRuleCustomErrorPolicy(rule.CustomErrorResponsePolicy)
		}
		result = append(result, args)
	}
	return result
}

func pathRuleUrlRedirect(r *gcpurlmapv1alpha1.GcpUrlMapUrlRedirect) *compute.URLMapPathMatcherPathRuleUrlRedirectArgs {
	return &compute.URLMapPathMatcherPathRuleUrlRedirectArgs{
		HostRedirect:         emptyAsNilString(r.HostRedirect),
		HttpsRedirect:        pulumi.Bool(r.HttpsRedirect),
		PathRedirect:         emptyAsNilString(r.PathRedirect),
		PrefixRedirect:       emptyAsNilString(r.PrefixRedirect),
		RedirectResponseCode: emptyAsNilString(r.RedirectResponseCode),
		StripQuery:           pulumi.Bool(r.StripQuery),
	}
}

func pathRuleRouteAction(a *gcpurlmapv1alpha1.GcpUrlMapRouteAction) *compute.URLMapPathMatcherPathRuleRouteActionArgs {
	out := &compute.URLMapPathMatcherPathRuleRouteActionArgs{}
	if len(a.WeightedBackendServices) > 0 {
		out.WeightedBackendServices = pathRuleWeightedBackends(a.WeightedBackendServices)
	}
	if a.UrlRewrite != nil {
		out.UrlRewrite = pathRuleUrlRewrite(a.UrlRewrite)
	}
	return out
}

func pathRuleWeightedBackends(backends []*gcpurlmapv1alpha1.GcpUrlMapWeightedBackendService) compute.URLMapPathMatcherPathRuleRouteActionWeightedBackendServiceArray {
	result := compute.URLMapPathMatcherPathRuleRouteActionWeightedBackendServiceArray{}
	for _, backend := range backends {
		result = append(result, &compute.URLMapPathMatcherPathRuleRouteActionWeightedBackendServiceArgs{
			BackendService: pulumi.String(backend.BackendService.GetValue()),
			Weight:         pulumi.Int(int(backend.Weight)),
		})
	}
	return result
}

func pathRuleUrlRewrite(r *gcpurlmapv1alpha1.GcpUrlMapUrlRewrite) *compute.URLMapPathMatcherPathRuleRouteActionUrlRewriteArgs {
	return &compute.URLMapPathMatcherPathRuleRouteActionUrlRewriteArgs{
		HostRewrite:       emptyAsNilString(r.HostRewrite),
		PathPrefixRewrite: emptyAsNilString(r.PathPrefixRewrite),
	}
}

func pathRuleCustomErrorPolicy(p *gcpurlmapv1alpha1.GcpUrlMapCustomErrorResponsePolicy) *compute.URLMapPathMatcherPathRuleCustomErrorResponsePolicyArgs {
	return &compute.URLMapPathMatcherPathRuleCustomErrorResponsePolicyArgs{
		ErrorService:       emptyAsNilString(p.ErrorService.GetValue()),
		ErrorResponseRules: pathRuleErrorRules(p.ErrorResponseRules),
	}
}

func pathRuleErrorRules(rules []*gcpurlmapv1alpha1.GcpUrlMapCustomErrorResponseRule) compute.URLMapPathMatcherPathRuleCustomErrorResponsePolicyErrorResponseRuleArray {
	result := compute.URLMapPathMatcherPathRuleCustomErrorResponsePolicyErrorResponseRuleArray{}
	for _, rule := range rules {
		codes := pulumi.StringArray{}
		for _, code := range rule.MatchResponseCodes {
			codes = append(codes, pulumi.String(code))
		}
		args := &compute.URLMapPathMatcherPathRuleCustomErrorResponsePolicyErrorResponseRuleArgs{
			MatchResponseCodes: codes,
			Path:               pulumi.String(rule.Path),
		}
		if rule.OverrideResponseCode != 0 {
			args.OverrideResponseCode = pulumi.Int(int(rule.OverrideResponseCode))
		}
		result = append(result, args)
	}
	return result
}

func buildRouteRules(rules []*gcpurlmapv1alpha1.GcpUrlMapRouteRule) compute.URLMapPathMatcherRouteRuleArray {
	result := compute.URLMapPathMatcherRouteRuleArray{}
	for _, rule := range rules {
		args := &compute.URLMapPathMatcherRouteRuleArgs{
			Priority:   pulumi.Int(int(rule.Priority)),
			MatchRules: buildMatchRules(rule.MatchRules),
		}
		if rule.Service.GetValue() != "" {
			args.Service = pulumi.String(rule.Service.GetValue())
		}
		if rule.UrlRedirect != nil {
			args.UrlRedirect = routeRuleUrlRedirect(rule.UrlRedirect)
		}
		if rule.RouteAction != nil {
			args.RouteAction = routeRuleRouteAction(rule.RouteAction)
		}
		if rule.HeaderAction != nil {
			args.HeaderAction = routeRuleHeaderAction(rule.HeaderAction)
		}
		if rule.CustomErrorResponsePolicy != nil {
			args.CustomErrorResponsePolicy = routeRuleCustomErrorPolicy(rule.CustomErrorResponsePolicy)
		}
		result = append(result, args)
	}
	return result
}

func routeRuleUrlRedirect(r *gcpurlmapv1alpha1.GcpUrlMapUrlRedirect) *compute.URLMapPathMatcherRouteRuleUrlRedirectArgs {
	return &compute.URLMapPathMatcherRouteRuleUrlRedirectArgs{
		HostRedirect:         emptyAsNilString(r.HostRedirect),
		HttpsRedirect:        pulumi.Bool(r.HttpsRedirect),
		PathRedirect:         emptyAsNilString(r.PathRedirect),
		PrefixRedirect:       emptyAsNilString(r.PrefixRedirect),
		RedirectResponseCode: emptyAsNilString(r.RedirectResponseCode),
		StripQuery:           pulumi.Bool(r.StripQuery),
	}
}

func routeRuleRouteAction(a *gcpurlmapv1alpha1.GcpUrlMapRouteAction) *compute.URLMapPathMatcherRouteRuleRouteActionArgs {
	out := &compute.URLMapPathMatcherRouteRuleRouteActionArgs{}
	if len(a.WeightedBackendServices) > 0 {
		out.WeightedBackendServices = routeRuleWeightedBackends(a.WeightedBackendServices)
	}
	if a.UrlRewrite != nil {
		out.UrlRewrite = routeRuleUrlRewrite(a.UrlRewrite)
	}
	return out
}

func routeRuleWeightedBackends(backends []*gcpurlmapv1alpha1.GcpUrlMapWeightedBackendService) compute.URLMapPathMatcherRouteRuleRouteActionWeightedBackendServiceArray {
	result := compute.URLMapPathMatcherRouteRuleRouteActionWeightedBackendServiceArray{}
	for _, backend := range backends {
		result = append(result, &compute.URLMapPathMatcherRouteRuleRouteActionWeightedBackendServiceArgs{
			BackendService: pulumi.String(backend.BackendService.GetValue()),
			Weight:         pulumi.Int(int(backend.Weight)),
		})
	}
	return result
}

func routeRuleUrlRewrite(r *gcpurlmapv1alpha1.GcpUrlMapUrlRewrite) *compute.URLMapPathMatcherRouteRuleRouteActionUrlRewriteArgs {
	return &compute.URLMapPathMatcherRouteRuleRouteActionUrlRewriteArgs{
		HostRewrite:         emptyAsNilString(r.HostRewrite),
		PathPrefixRewrite:   emptyAsNilString(r.PathPrefixRewrite),
		PathTemplateRewrite: emptyAsNilString(r.PathTemplateRewrite),
	}
}

func routeRuleCustomErrorPolicy(p *gcpurlmapv1alpha1.GcpUrlMapCustomErrorResponsePolicy) *compute.URLMapPathMatcherRouteRuleCustomErrorResponsePolicyArgs {
	return &compute.URLMapPathMatcherRouteRuleCustomErrorResponsePolicyArgs{
		ErrorService:       emptyAsNilString(p.ErrorService.GetValue()),
		ErrorResponseRules: routeRuleErrorRules(p.ErrorResponseRules),
	}
}

func routeRuleErrorRules(rules []*gcpurlmapv1alpha1.GcpUrlMapCustomErrorResponseRule) compute.URLMapPathMatcherRouteRuleCustomErrorResponsePolicyErrorResponseRuleArray {
	result := compute.URLMapPathMatcherRouteRuleCustomErrorResponsePolicyErrorResponseRuleArray{}
	for _, rule := range rules {
		codes := pulumi.StringArray{}
		for _, code := range rule.MatchResponseCodes {
			codes = append(codes, pulumi.String(code))
		}
		args := &compute.URLMapPathMatcherRouteRuleCustomErrorResponsePolicyErrorResponseRuleArgs{
			MatchResponseCodes: codes,
			Path:               pulumi.String(rule.Path),
		}
		if rule.OverrideResponseCode != 0 {
			args.OverrideResponseCode = pulumi.Int(int(rule.OverrideResponseCode))
		}
		result = append(result, args)
	}
	return result
}

func routeRuleHeaderAction(h *gcpurlmapv1alpha1.GcpUrlMapHeaderAction) *compute.URLMapPathMatcherRouteRuleHeaderActionArgs {
	return &compute.URLMapPathMatcherRouteRuleHeaderActionArgs{
		RequestHeadersToAdds:     routeRuleRequestHeadersToAdd(h.RequestHeadersToAdd),
		RequestHeadersToRemoves:  stringArrayOrNil(h.RequestHeadersToRemove),
		ResponseHeadersToAdds:    routeRuleResponseHeadersToAdd(h.ResponseHeadersToAdd),
		ResponseHeadersToRemoves: stringArrayOrNil(h.ResponseHeadersToRemove),
	}
}

func routeRuleRequestHeadersToAdd(headers []*gcpurlmapv1alpha1.GcpUrlMapHeaderValue) compute.URLMapPathMatcherRouteRuleHeaderActionRequestHeadersToAddArray {
	result := compute.URLMapPathMatcherRouteRuleHeaderActionRequestHeadersToAddArray{}
	for _, header := range headers {
		result = append(result, &compute.URLMapPathMatcherRouteRuleHeaderActionRequestHeadersToAddArgs{
			HeaderName:  pulumi.String(header.HeaderName),
			HeaderValue: pulumi.String(header.HeaderValue),
			Replace:     pulumi.Bool(header.Replace),
		})
	}
	return result
}

func routeRuleResponseHeadersToAdd(headers []*gcpurlmapv1alpha1.GcpUrlMapHeaderValue) compute.URLMapPathMatcherRouteRuleHeaderActionResponseHeadersToAddArray {
	result := compute.URLMapPathMatcherRouteRuleHeaderActionResponseHeadersToAddArray{}
	for _, header := range headers {
		result = append(result, &compute.URLMapPathMatcherRouteRuleHeaderActionResponseHeadersToAddArgs{
			HeaderName:  pulumi.String(header.HeaderName),
			HeaderValue: pulumi.String(header.HeaderValue),
			Replace:     pulumi.Bool(header.Replace),
		})
	}
	return result
}

func buildMatchRules(rules []*gcpurlmapv1alpha1.GcpUrlMapRouteRuleMatchRule) compute.URLMapPathMatcherRouteRuleMatchRuleArray {
	result := compute.URLMapPathMatcherRouteRuleMatchRuleArray{}
	for _, rule := range rules {
		args := &compute.URLMapPathMatcherRouteRuleMatchRuleArgs{
			PrefixMatch:           emptyAsNilString(rule.PrefixMatch),
			FullPathMatch:         emptyAsNilString(rule.FullPathMatch),
			RegexMatch:            emptyAsNilString(rule.RegexMatch),
			PathTemplateMatch:     emptyAsNilString(rule.PathTemplateMatch),
			IgnoreCase:            pulumi.Bool(rule.IgnoreCase),
			HeaderMatches:         buildHeaderMatches(rule.HeaderMatches),
			QueryParameterMatches: buildQueryParameterMatches(rule.QueryParameterMatches),
			MetadataFilters:       buildMetadataFilters(rule.MetadataFilters),
		}
		result = append(result, args)
	}
	return result
}

func buildHeaderMatches(matches []*gcpurlmapv1alpha1.GcpUrlMapHeaderMatch) compute.URLMapPathMatcherRouteRuleMatchRuleHeaderMatchArray {
	result := compute.URLMapPathMatcherRouteRuleMatchRuleHeaderMatchArray{}
	for _, match := range matches {
		args := &compute.URLMapPathMatcherRouteRuleMatchRuleHeaderMatchArgs{
			HeaderName:   pulumi.String(match.HeaderName),
			ExactMatch:   emptyAsNilString(match.ExactMatch),
			PrefixMatch:  emptyAsNilString(match.PrefixMatch),
			SuffixMatch:  emptyAsNilString(match.SuffixMatch),
			RegexMatch:   emptyAsNilString(match.RegexMatch),
			PresentMatch: pulumi.Bool(match.PresentMatch),
			InvertMatch:  pulumi.Bool(match.InvertMatch),
		}
		if match.RangeMatch != nil {
			args.RangeMatch = &compute.URLMapPathMatcherRouteRuleMatchRuleHeaderMatchRangeMatchArgs{
				RangeStart: pulumi.Int(int(match.RangeMatch.RangeStart)),
				RangeEnd:   pulumi.Int(int(match.RangeMatch.RangeEnd)),
			}
		}
		result = append(result, args)
	}
	return result
}

func buildQueryParameterMatches(matches []*gcpurlmapv1alpha1.GcpUrlMapQueryParameterMatch) compute.URLMapPathMatcherRouteRuleMatchRuleQueryParameterMatchArray {
	result := compute.URLMapPathMatcherRouteRuleMatchRuleQueryParameterMatchArray{}
	for _, match := range matches {
		result = append(result, &compute.URLMapPathMatcherRouteRuleMatchRuleQueryParameterMatchArgs{
			Name:         pulumi.String(match.Name),
			ExactMatch:   emptyAsNilString(match.ExactMatch),
			PresentMatch: pulumi.Bool(match.PresentMatch),
			RegexMatch:   emptyAsNilString(match.RegexMatch),
		})
	}
	return result
}

func buildMetadataFilters(filters []*gcpurlmapv1alpha1.GcpUrlMapMetadataFilter) compute.URLMapPathMatcherRouteRuleMatchRuleMetadataFilterArray {
	result := compute.URLMapPathMatcherRouteRuleMatchRuleMetadataFilterArray{}
	for _, filter := range filters {
		labels := compute.URLMapPathMatcherRouteRuleMatchRuleMetadataFilterFilterLabelArray{}
		for _, label := range filter.FilterLabels {
			labels = append(labels, &compute.URLMapPathMatcherRouteRuleMatchRuleMetadataFilterFilterLabelArgs{
				Name:  pulumi.String(label.Name),
				Value: pulumi.String(label.Value),
			})
		}
		result = append(result, &compute.URLMapPathMatcherRouteRuleMatchRuleMetadataFilterArgs{
			FilterMatchCriteria: pulumi.String(filter.FilterMatchCriteria),
			FilterLabels:        labels,
		})
	}
	return result
}

func buildTests(tests []*gcpurlmapv1alpha1.GcpUrlMapTest) compute.URLMapTestArray {
	result := compute.URLMapTestArray{}
	for _, test := range tests {
		args := &compute.URLMapTestArgs{
			Host: pulumi.String(test.Host),
			Path: pulumi.String(test.Path),
		}
		if test.Service.GetValue() != "" {
			args.Service = pulumi.String(test.Service.GetValue())
		}
		if test.Description != "" {
			args.Description = pulumi.String(test.Description)
		}
		if test.ExpectedOutputUrl != "" {
			args.ExpectedOutputUrl = pulumi.String(test.ExpectedOutputUrl)
		}
		if test.ExpectedRedirectResponseCode != 0 {
			args.ExpectedRedirectResponseCode = pulumi.Int(int(test.ExpectedRedirectResponseCode))
		}
		if len(test.Headers) > 0 {
			headers := compute.URLMapTestHeaderArray{}
			for _, header := range test.Headers {
				headers = append(headers, &compute.URLMapTestHeaderArgs{
					Name:  pulumi.String(header.Name),
					Value: pulumi.String(header.Value),
				})
			}
			args.Headers = headers
		}
		result = append(result, args)
	}
	return result
}

func emptyAsNilString(value string) pulumi.StringPtrInput {
	if value == "" {
		return nil
	}
	return pulumi.String(value)
}

func stringArrayOrNil(values []string) pulumi.StringArrayInput {
	if len(values) == 0 {
		return nil
	}
	result := pulumi.StringArray{}
	for _, value := range values {
		result = append(result, pulumi.String(value))
	}
	return result
}
