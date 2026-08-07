package module

import (
	"github.com/pkg/errors"
	awswafwebaclv1alpha1 "github.com/plantonhq/planton/catalog/aws/awswafwebacl/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/wafv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// webAcl creates the WAFv2 Web ACL resource. Rules are passed via RuleJson —
// the typed statement tree is serialized to the AWS API JSON (see rules.go),
// which handles the recursive statement language and the custom_statement
// escape hatch uniformly.
func webAcl(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*wafv2.WebAcl, error) {
	spec := locals.WebAcl.Spec

	args := &wafv2.WebAclArgs{
		// The ACL's AWS name is the Planton resource name — the stable
		// identity operators see. Name and scope are create-time immutable
		// (ForceNew).
		Name:             pulumi.String(locals.WebAcl.Metadata.Name),
		Scope:            pulumi.String(spec.Scope),
		DefaultAction:    buildDefaultAction(spec.DefaultAction),
		VisibilityConfig: buildAclVisibilityConfig(locals),
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	if len(spec.TokenDomains) > 0 {
		args.TokenDomains = pulumi.ToStringArray(spec.TokenDomains)
	}

	if len(spec.CustomResponseBodies) > 0 {
		bodies := wafv2.WebAclCustomResponseBodyArray{}
		for _, body := range spec.CustomResponseBodies {
			bodies = append(bodies, &wafv2.WebAclCustomResponseBodyArgs{
				Key:         pulumi.String(body.Key),
				Content:     pulumi.String(body.Content),
				ContentType: pulumi.String(body.ContentType),
			})
		}
		args.CustomResponseBodies = bodies
	}

	// Web-ACL-wide immunity windows for CAPTCHA solves and silent-challenge
	// responses (rules can override per rule inside the rule JSON).
	if spec.CaptchaConfig != nil {
		args.CaptchaConfig = &wafv2.WebAclCaptchaConfigArgs{
			ImmunityTimeProperty: &wafv2.WebAclCaptchaConfigImmunityTimePropertyArgs{
				ImmunityTime: pulumi.Int(int(spec.CaptchaConfig.ImmunityTimeSec)),
			},
		}
	}
	if spec.ChallengeConfig != nil {
		args.ChallengeConfig = &wafv2.WebAclChallengeConfigArgs{
			ImmunityTimeProperty: &wafv2.WebAclChallengeConfigImmunityTimePropertyArgs{
				ImmunityTime: pulumi.Int(int(spec.ChallengeConfig.ImmunityTimeSec)),
			},
		}
	}

	// Per-resource-type request-body inspection limits (default 16 KB).
	if spec.AssociationConfig != nil {
		args.AssociationConfig = buildAssociationConfig(spec.AssociationConfig)
	}

	// Field-level masking in ALL WAF outputs (stronger than log redaction).
	if spec.DataProtectionConfig != nil {
		args.DataProtectionConfig = buildDataProtectionConfig(spec.DataProtectionConfig)
	}

	// Rules as AWS API JSON — the single rule surface both engines share.
	if len(spec.Rules) > 0 {
		rulesJSON, err := buildRulesJSON(spec)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build rules JSON")
		}
		args.RuleJson = pulumi.StringPtr(rulesJSON)
	}

	createdAcl, err := wafv2.NewWebAcl(ctx, locals.WebAcl.Metadata.Name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create WAFv2 Web ACL")
	}

	ctx.Export(OpWebAclArn, createdAcl.Arn)
	ctx.Export(OpWebAclId, createdAcl.ID())
	ctx.Export(OpWebAclName, createdAcl.Name)
	ctx.Export(OpCapacity, createdAcl.Capacity)
	ctx.Export(OpApplicationIntegrationUrl, createdAcl.ApplicationIntegrationUrl)

	return createdAcl, nil
}

// buildDefaultAction maps the spec's default action — including the custom
// block response and custom allow request headers the action types carry.
func buildDefaultAction(defaultAction *awswafwebaclv1alpha1.AwsWafWebAclDefaultAction) *wafv2.WebAclDefaultActionArgs {
	if defaultAction.Type == "block" {
		blockArgs := &wafv2.WebAclDefaultActionBlockArgs{}
		if cr := defaultAction.CustomResponse; cr != nil {
			customResponse := &wafv2.WebAclDefaultActionBlockCustomResponseArgs{
				ResponseCode: pulumi.Int(int(cr.ResponseCode)),
			}
			if cr.CustomResponseBodyKey != "" {
				customResponse.CustomResponseBodyKey = pulumi.StringPtr(cr.CustomResponseBodyKey)
			}
			if len(cr.ResponseHeaders) > 0 {
				headers := wafv2.WebAclDefaultActionBlockCustomResponseResponseHeaderArray{}
				for _, header := range cr.ResponseHeaders {
					headers = append(headers, &wafv2.WebAclDefaultActionBlockCustomResponseResponseHeaderArgs{
						Name:  pulumi.String(header.Name),
						Value: pulumi.String(header.Value),
					})
				}
				customResponse.ResponseHeaders = headers
			}
			blockArgs.CustomResponse = customResponse
		}
		return &wafv2.WebAclDefaultActionArgs{Block: blockArgs}
	}

	allowArgs := &wafv2.WebAclDefaultActionAllowArgs{}
	if len(defaultAction.CustomRequestHeaders) > 0 {
		headers := wafv2.WebAclDefaultActionAllowCustomRequestHandlingInsertHeaderArray{}
		for _, header := range defaultAction.CustomRequestHeaders {
			headers = append(headers, &wafv2.WebAclDefaultActionAllowCustomRequestHandlingInsertHeaderArgs{
				Name:  pulumi.String(header.Name),
				Value: pulumi.String(header.Value),
			})
		}
		allowArgs.CustomRequestHandling = &wafv2.WebAclDefaultActionAllowCustomRequestHandlingArgs{
			InsertHeaders: headers,
		}
	}
	return &wafv2.WebAclDefaultActionArgs{Allow: allowArgs}
}

// buildAclVisibilityConfig applies the ACL-level visibility defaults
// (metrics on, sampling on, metric name = resource name) when the spec omits
// the block — identical defaults in both engines.
func buildAclVisibilityConfig(locals *Locals) *wafv2.WebAclVisibilityConfigArgs {
	spec := locals.WebAcl.Spec

	metricName := locals.WebAcl.Metadata.Name
	metricsEnabled := true
	sampledEnabled := true
	if spec.VisibilityConfig != nil {
		if spec.VisibilityConfig.MetricName != "" {
			metricName = spec.VisibilityConfig.MetricName
		}
		metricsEnabled = spec.VisibilityConfig.CloudwatchMetricsEnabled
		sampledEnabled = spec.VisibilityConfig.SampledRequestsEnabled
	}

	return &wafv2.WebAclVisibilityConfigArgs{
		CloudwatchMetricsEnabled: pulumi.Bool(metricsEnabled),
		SampledRequestsEnabled:   pulumi.Bool(sampledEnabled),
		MetricName:               pulumi.String(metricName),
	}
}

// buildAssociationConfig maps the per-resource-type body inspection limits.
func buildAssociationConfig(config *awswafwebaclv1alpha1.AwsWafWebAclAssociationConfig) *wafv2.WebAclAssociationConfigArgs {
	requestBody := &wafv2.WebAclAssociationConfigRequestBodyArgs{}
	if config.CloudfrontRequestBodyLimit != "" {
		requestBody.Cloudfront = &wafv2.WebAclAssociationConfigRequestBodyCloudfrontArgs{
			DefaultSizeInspectionLimit: pulumi.String(config.CloudfrontRequestBodyLimit),
		}
	}
	if config.ApiGatewayRequestBodyLimit != "" {
		requestBody.ApiGateway = &wafv2.WebAclAssociationConfigRequestBodyApiGatewayArgs{
			DefaultSizeInspectionLimit: pulumi.String(config.ApiGatewayRequestBodyLimit),
		}
	}
	if config.CognitoUserPoolRequestBodyLimit != "" {
		requestBody.CognitoUserPool = &wafv2.WebAclAssociationConfigRequestBodyCognitoUserPoolArgs{
			DefaultSizeInspectionLimit: pulumi.String(config.CognitoUserPoolRequestBodyLimit),
		}
	}
	if config.AppRunnerServiceRequestBodyLimit != "" {
		requestBody.AppRunnerService = &wafv2.WebAclAssociationConfigRequestBodyAppRunnerServiceArgs{
			DefaultSizeInspectionLimit: pulumi.String(config.AppRunnerServiceRequestBodyLimit),
		}
	}
	if config.VerifiedAccessInstanceRequestBodyLimit != "" {
		requestBody.VerifiedAccessInstance = &wafv2.WebAclAssociationConfigRequestBodyVerifiedAccessInstanceArgs{
			DefaultSizeInspectionLimit: pulumi.String(config.VerifiedAccessInstanceRequestBodyLimit),
		}
	}
	return &wafv2.WebAclAssociationConfigArgs{
		RequestBodies: wafv2.WebAclAssociationConfigRequestBodyArray{requestBody},
	}
}

// buildDataProtectionConfig maps the field-masking configuration.
func buildDataProtectionConfig(config *awswafwebaclv1alpha1.AwsWafWebAclDataProtectionConfig) *wafv2.WebAclDataProtectionConfigArgs {
	protections := wafv2.WebAclDataProtectionConfigDataProtectionArray{}
	for _, protection := range config.DataProtections {
		field := &wafv2.WebAclDataProtectionConfigDataProtectionFieldArgs{
			FieldType: pulumi.String(protection.FieldType),
		}
		if len(protection.FieldKeys) > 0 {
			field.FieldKeys = pulumi.ToStringArray(protection.FieldKeys)
		}
		args := &wafv2.WebAclDataProtectionConfigDataProtectionArgs{
			Action: pulumi.String(protection.Action),
			Field:  field,
		}
		if protection.ExcludeRuleMatchDetails {
			args.ExcludeRuleMatchDetails = pulumi.BoolPtr(true)
		}
		if protection.ExcludeRateBasedDetails {
			args.ExcludeRateBasedDetails = pulumi.BoolPtr(true)
		}
		protections = append(protections, args)
	}
	return &wafv2.WebAclDataProtectionConfigArgs{DataProtections: protections}
}
