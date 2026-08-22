package module

import (
	"github.com/pkg/errors"
	cloudflareaigatewayv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareaigateway/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// aiGateway creates the gateway itself. The spec groups Cloudflare's flat
// retry_* and log_management* arguments into retry{} / log_management{}
// messages for authoring clarity; this function fans them back out to the
// provider's flat shape. The gateway id (spec.gateway_id) is the provider's
// `id` argument -- user-chosen, create-only, and the URL slug every client
// calls.
func aiGateway(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) (*cloudflare.AiGateway, error) {
	spec := locals.CloudflareAiGateway.Spec

	args := &cloudflare.AiGatewayArgs{
		AccountId:   pulumi.String(spec.AccountId),
		AiGatewayId: pulumi.String(spec.GatewayId),
		// The five required scalars: Cloudflare demands an explicit choice
		// for each, and the spec's CEL required rules guarantee presence.
		CacheInvalidateOnUpdate: pulumi.Bool(spec.GetCacheInvalidateOnUpdate()),
		CacheTtl:                pulumi.Int(int(spec.GetCacheTtl())),
		CollectLogs:             pulumi.Bool(spec.GetCollectLogs()),
		RateLimitingInterval:    pulumi.Int(int(spec.GetRateLimitingInterval())),
		RateLimitingLimit:       pulumi.Int(int(spec.GetRateLimitingLimit())),
	}

	if spec.RateLimitingTechnique != "" {
		args.RateLimitingTechnique = pulumi.StringPtr(spec.RateLimitingTechnique)
	}

	// retry{} fans out to the provider's three flat retry_* arguments.
	if spec.Retry != nil {
		if spec.Retry.Backoff != "" {
			args.RetryBackoff = pulumi.StringPtr(spec.Retry.Backoff)
		}
		if spec.Retry.Delay != nil {
			args.RetryDelay = pulumi.IntPtr(int(spec.Retry.GetDelay()))
		}
		if spec.Retry.MaxAttempts != nil {
			args.RetryMaxAttempts = pulumi.IntPtr(int(spec.Retry.GetMaxAttempts()))
		}
	}

	// log_management{} fans out to the flat log_management (record cap) and
	// log_management_strategy arguments.
	if spec.LogManagement != nil {
		if spec.LogManagement.MaxRecords != nil {
			args.LogManagement = pulumi.IntPtr(int(spec.LogManagement.GetMaxRecords()))
		}
		if spec.LogManagement.Strategy != "" {
			args.LogManagementStrategy = pulumi.StringPtr(spec.LogManagement.Strategy)
		}
	}

	if spec.Authentication != nil {
		args.Authentication = pulumi.BoolPtr(spec.GetAuthentication())
	}
	if spec.Logpush != nil {
		args.Logpush = pulumi.BoolPtr(spec.GetLogpush())
	}
	if spec.LogpushPublicKey != "" {
		args.LogpushPublicKey = pulumi.StringPtr(spec.LogpushPublicKey)
	}
	if spec.Zdr != nil {
		args.Zdr = pulumi.BoolPtr(spec.GetZdr())
	}
	if spec.WorkersAiBillingMode != "" {
		args.WorkersAiBillingMode = pulumi.StringPtr(spec.WorkersAiBillingMode)
	}
	if spec.StoreId.GetValue() != "" {
		args.StoreId = pulumi.StringPtr(spec.StoreId.GetValue())
	}

	if spec.Dlp != nil {
		args.Dlp = buildDlp(spec.Dlp)
	}
	if spec.Guardrails != nil {
		args.Guardrails = buildGuardrails(spec.Guardrails)
	}
	if len(spec.Otel) > 0 {
		args.Otels = buildOtels(spec.Otel)
	}
	if spec.Stripe != nil {
		args.Stripe = buildStripe(spec.Stripe)
	}
	if spec.SpendLimits != nil {
		args.SpendLimits = buildSpendLimits(spec.SpendLimits)
	}

	createdGateway, err := cloudflare.NewAiGateway(
		ctx,
		"ai_gateway",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ai gateway")
	}

	ctx.Export(OpGatewayId, createdGateway.AiGatewayId)

	return createdGateway, nil
}

func buildDlp(dlp *cloudflareaigatewayv1alpha1.CloudflareAiGatewayDlp) cloudflare.AiGatewayDlpPtrInput {
	args := cloudflare.AiGatewayDlpArgs{
		Enabled: pulumi.Bool(dlp.GetEnabled()),
	}
	if dlp.Action != "" {
		args.Action = pulumi.StringPtr(dlp.Action)
	}
	if len(dlp.Profiles) > 0 {
		args.Profiles = pulumi.ToStringArray(dlp.Profiles)
	}
	if len(dlp.Policies) > 0 {
		policies := cloudflare.AiGatewayDlpPolicyArray{}
		for _, policy := range dlp.Policies {
			policies = append(policies, cloudflare.AiGatewayDlpPolicyArgs{
				Id:       pulumi.String(policy.Id),
				Enabled:  pulumi.Bool(policy.GetEnabled()),
				Action:   pulumi.String(policy.Action),
				Checks:   pulumi.ToStringArray(policy.Check),
				Profiles: pulumi.ToStringArray(policy.Profiles),
			})
		}
		args.Policies = policies
	}
	return args
}

// buildGuardrails maps the per-hazard-code controls. The prompt and
// response sides carry identical fields but distinct SDK types, so each is
// built explicitly. An unset code stays nil -- that category is simply not
// evaluated.
func buildGuardrails(guardrails *cloudflareaigatewayv1alpha1.CloudflareAiGatewayGuardrails) cloudflare.AiGatewayGuardrailsPtrInput {
	prompt := guardrails.Prompt
	response := guardrails.Response
	return cloudflare.AiGatewayGuardrailsArgs{
		Prompt: cloudflare.AiGatewayGuardrailsPromptArgs{
			P1:  optionalControl(prompt.P1),
			S1:  optionalControl(prompt.S1),
			S2:  optionalControl(prompt.S2),
			S3:  optionalControl(prompt.S3),
			S4:  optionalControl(prompt.S4),
			S5:  optionalControl(prompt.S5),
			S6:  optionalControl(prompt.S6),
			S7:  optionalControl(prompt.S7),
			S8:  optionalControl(prompt.S8),
			S9:  optionalControl(prompt.S9),
			S10: optionalControl(prompt.S10),
			S11: optionalControl(prompt.S11),
			S12: optionalControl(prompt.S12),
			S13: optionalControl(prompt.S13),
		},
		Response: cloudflare.AiGatewayGuardrailsResponseArgs{
			P1:  optionalControl(response.P1),
			S1:  optionalControl(response.S1),
			S2:  optionalControl(response.S2),
			S3:  optionalControl(response.S3),
			S4:  optionalControl(response.S4),
			S5:  optionalControl(response.S5),
			S6:  optionalControl(response.S6),
			S7:  optionalControl(response.S7),
			S8:  optionalControl(response.S8),
			S9:  optionalControl(response.S9),
			S10: optionalControl(response.S10),
			S11: optionalControl(response.S11),
			S12: optionalControl(response.S12),
			S13: optionalControl(response.S13),
		},
	}
}

// optionalControl turns an unset guardrail control into nil so the API
// leaves that hazard category unevaluated.
func optionalControl(value string) pulumi.StringPtrInput {
	if value == "" {
		return nil
	}
	return pulumi.StringPtr(value)
}

func buildOtels(otels []*cloudflareaigatewayv1alpha1.CloudflareAiGatewayOtel) cloudflare.AiGatewayOtelArrayInput {
	built := cloudflare.AiGatewayOtelArray{}
	for _, otel := range otels {
		// headers is Required at the API -- send an empty map when the
		// manifest declares none.
		headers := pulumi.StringMap{}
		for key, value := range otel.Headers {
			headers[key] = pulumi.String(value)
		}
		otelArgs := cloudflare.AiGatewayOtelArgs{
			Url:     pulumi.String(otel.Url),
			Headers: headers,
		}
		// The authorization header value is a credential: kept secret in
		// Pulumi state even though the provider leaves it unmarked.
		if otel.Authorization.GetValue() != "" {
			otelArgs.Authorization = pulumi.ToSecret(pulumi.String(otel.Authorization.GetValue())).(pulumi.StringOutput)
		}
		if otel.ContentType != "" {
			otelArgs.ContentType = pulumi.StringPtr(otel.ContentType)
		}
		built = append(built, otelArgs)
	}
	return built
}

func buildStripe(stripe *cloudflareaigatewayv1alpha1.CloudflareAiGatewayStripe) cloudflare.AiGatewayStripePtrInput {
	usageEvents := cloudflare.AiGatewayStripeUsageEventArray{}
	for _, event := range stripe.UsageEvents {
		usageEvents = append(usageEvents, cloudflare.AiGatewayStripeUsageEventArgs{
			Payload: pulumi.String(event.Payload),
		})
	}
	return cloudflare.AiGatewayStripeArgs{
		// The Stripe credential: kept secret in Pulumi state even though the
		// provider leaves it unmarked.
		Authorization: pulumi.ToSecret(pulumi.String(stripe.Authorization.GetValue())).(pulumi.StringOutput),
		UsageEvents:   usageEvents,
	}
}

func buildSpendLimits(spendLimits *cloudflareaigatewayv1alpha1.CloudflareAiGatewaySpendLimits) cloudflare.AiGatewaySpendLimitsPtrInput {
	args := cloudflare.AiGatewaySpendLimitsArgs{}
	if spendLimits.Enabled != nil {
		args.Enabled = pulumi.BoolPtr(spendLimits.GetEnabled())
	}
	if len(spendLimits.Rules) > 0 {
		rules := cloudflare.AiGatewaySpendLimitsRuleArray{}
		for _, rule := range spendLimits.Rules {
			// The rule id is always sent explicitly: the provider schema's
			// default is a leaked example value shared by every omitted id,
			// which would collapse multiple rules into one -- the spec
			// requires an explicit unique id per rule.
			ruleArgs := cloudflare.AiGatewaySpendLimitsRuleArgs{
				Id:        pulumi.StringPtr(rule.Id),
				Limit:     pulumi.Float64(rule.GetLimit()),
				LimitType: pulumi.String(rule.LimitType),
				Window:    pulumi.Int(int(rule.GetWindow())),
			}
			if rule.Enabled != nil {
				ruleArgs.Enabled = pulumi.BoolPtr(rule.GetEnabled())
			}
			if rule.Technique != "" {
				ruleArgs.Technique = pulumi.StringPtr(rule.Technique)
			}
			if len(rule.Metadata) > 0 {
				metadata := cloudflare.AiGatewaySpendLimitsRuleMetadataMap{}
				for key, filter := range rule.Metadata {
					metadata[key] = cloudflare.AiGatewaySpendLimitsRuleMetadataArgs{
						Mode:   pulumi.String(filter.Mode),
						Values: pulumi.ToStringArray(filter.Values),
					}
				}
				ruleArgs.Metadata = metadata
			}
			if rule.Model != nil {
				ruleArgs.Model = cloudflare.AiGatewaySpendLimitsRuleModelArgs{
					Mode:   pulumi.String(rule.Model.Mode),
					Values: pulumi.ToStringArray(rule.Model.Values),
				}
			}
			// The spec calls this field `provider` (the wire name); the
			// Terraform provider renamed it ai_gateway_provider to dodge the
			// reserved word.
			if rule.Provider != nil {
				ruleArgs.AiGatewayProvider = cloudflare.AiGatewaySpendLimitsRuleAiGatewayProviderArgs{
					Mode:   pulumi.String(rule.Provider.Mode),
					Values: pulumi.ToStringArray(rule.Provider.Values),
				}
			}
			rules = append(rules, ruleArgs)
		}
		args.Rules = rules
	}
	return args
}
