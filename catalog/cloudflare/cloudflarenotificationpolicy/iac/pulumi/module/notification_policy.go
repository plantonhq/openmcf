package module

import (
	"github.com/pkg/errors"
	cloudflarenotificationpolicyv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarenotificationpolicy/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// notificationPolicy creates the notification policy: one alert type fanned
// out to email, PagerDuty, and webhook destinations, optionally narrowed by
// filters. A plain CRUD resource (real create/update/delete; only the
// account forces replacement).
//
// The spec flattens each mechanism to its identity (an address or a UUID);
// this module rebuilds the API's object rows. `enabled` is sent only when
// set -- Cloudflare's default is true.
func notificationPolicy(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareNotificationPolicy.Spec

	mechanisms := cloudflare.NotificationPolicyMechanismsArgs{}
	if len(spec.Mechanisms.Emails) > 0 {
		emails := cloudflare.NotificationPolicyMechanismsEmailArray{}
		for _, address := range spec.Mechanisms.Emails {
			emails = append(emails, cloudflare.NotificationPolicyMechanismsEmailArgs{
				Id: pulumi.StringPtr(address),
			})
		}
		mechanisms.Emails = emails
	}
	if len(spec.Mechanisms.PagerdutyIds) > 0 {
		pagerduties := cloudflare.NotificationPolicyMechanismsPagerdutyArray{}
		for _, id := range spec.Mechanisms.PagerdutyIds {
			pagerduties = append(pagerduties, cloudflare.NotificationPolicyMechanismsPagerdutyArgs{
				Id: pulumi.StringPtr(id),
			})
		}
		mechanisms.Pagerduties = pagerduties
	}
	if len(spec.Mechanisms.WebhookIds) > 0 {
		webhooks := cloudflare.NotificationPolicyMechanismsWebhookArray{}
		for _, ref := range spec.Mechanisms.WebhookIds {
			webhooks = append(webhooks, cloudflare.NotificationPolicyMechanismsWebhookArgs{
				Id: pulumi.StringPtr(ref.GetValue()),
			})
		}
		mechanisms.Webhooks = webhooks
	}

	args := &cloudflare.NotificationPolicyArgs{
		AccountId:  pulumi.String(spec.AccountId),
		Name:       pulumi.String(spec.Name),
		AlertType:  pulumi.String(spec.AlertType),
		Mechanisms: mechanisms,
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.Enabled != nil {
		args.Enabled = pulumi.BoolPtr(spec.GetEnabled())
	}
	if spec.AlertInterval != "" {
		args.AlertInterval = pulumi.String(spec.AlertInterval)
	}
	if spec.Filters != nil {
		args.Filters = buildFilters(spec.Filters)
	}

	createdPolicy, err := cloudflare.NewNotificationPolicy(
		ctx,
		"notification_policy",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create notification policy")
	}

	ctx.Export(OpPolicyId, createdPolicy.ID())

	return nil
}

// buildFilters maps the declared filter lists; unset filters are never
// sent, so the payload holds exactly the fields the alert type reads.
func buildFilters(filters *cloudflarenotificationpolicyv1alpha1.CloudflareNotificationPolicyFilters) cloudflare.NotificationPolicyFiltersPtrInput {
	filterArgs := cloudflare.NotificationPolicyFiltersArgs{}

	if len(filters.Actions) > 0 {
		filterArgs.Actions = pulumi.ToStringArray(filters.Actions)
	}
	if len(filters.AffectedAsns) > 0 {
		filterArgs.AffectedAsns = pulumi.ToStringArray(filters.AffectedAsns)
	}
	if len(filters.AffectedComponents) > 0 {
		filterArgs.AffectedComponents = pulumi.ToStringArray(filters.AffectedComponents)
	}
	if len(filters.AffectedLocations) > 0 {
		filterArgs.AffectedLocations = pulumi.ToStringArray(filters.AffectedLocations)
	}
	if len(filters.AirportCode) > 0 {
		filterArgs.AirportCodes = pulumi.ToStringArray(filters.AirportCode)
	}
	if len(filters.AlertTriggerPreferences) > 0 {
		filterArgs.AlertTriggerPreferences = pulumi.ToStringArray(filters.AlertTriggerPreferences)
	}
	if len(filters.AlertTriggerPreferencesValue) > 0 {
		filterArgs.AlertTriggerPreferencesValues = pulumi.ToStringArray(filters.AlertTriggerPreferencesValue)
	}
	if len(filters.Enabled) > 0 {
		filterArgs.Enableds = pulumi.ToStringArray(filters.Enabled)
	}
	if len(filters.Environment) > 0 {
		filterArgs.Environments = pulumi.ToStringArray(filters.Environment)
	}
	if len(filters.Event) > 0 {
		filterArgs.Events = pulumi.ToStringArray(filters.Event)
	}
	if len(filters.EventSource) > 0 {
		filterArgs.EventSources = pulumi.ToStringArray(filters.EventSource)
	}
	if len(filters.EventType) > 0 {
		filterArgs.EventTypes = pulumi.ToStringArray(filters.EventType)
	}
	if len(filters.GroupBy) > 0 {
		filterArgs.GroupBies = pulumi.ToStringArray(filters.GroupBy)
	}
	if len(filters.HealthCheckId) > 0 {
		filterArgs.HealthCheckIds = pulumi.ToStringArray(filters.HealthCheckId)
	}
	if len(filters.IncidentImpact) > 0 {
		filterArgs.IncidentImpacts = pulumi.ToStringArray(filters.IncidentImpact)
	}
	if len(filters.InputId) > 0 {
		filterArgs.InputIds = pulumi.ToStringArray(filters.InputId)
	}
	if len(filters.InsightClass) > 0 {
		filterArgs.InsightClasses = pulumi.ToStringArray(filters.InsightClass)
	}
	if len(filters.Limit) > 0 {
		filterArgs.Limits = pulumi.ToStringArray(filters.Limit)
	}
	if len(filters.LogoTag) > 0 {
		filterArgs.LogoTags = pulumi.ToStringArray(filters.LogoTag)
	}
	if len(filters.MegabitsPerSecond) > 0 {
		filterArgs.MegabitsPerSeconds = pulumi.ToStringArray(filters.MegabitsPerSecond)
	}
	if len(filters.NewHealth) > 0 {
		filterArgs.NewHealths = pulumi.ToStringArray(filters.NewHealth)
	}
	if len(filters.NewStatus) > 0 {
		filterArgs.NewStatuses = pulumi.ToStringArray(filters.NewStatus)
	}
	if len(filters.PacketsPerSecond) > 0 {
		filterArgs.PacketsPerSeconds = pulumi.ToStringArray(filters.PacketsPerSecond)
	}
	if len(filters.PoolId) > 0 {
		filterArgs.PoolIds = pulumi.ToStringArray(filters.PoolId)
	}
	if len(filters.PopNames) > 0 {
		filterArgs.PopNames = pulumi.ToStringArray(filters.PopNames)
	}
	if len(filters.Product) > 0 {
		filterArgs.Products = pulumi.ToStringArray(filters.Product)
	}
	if len(filters.ProjectId) > 0 {
		filterArgs.ProjectIds = pulumi.ToStringArray(filters.ProjectId)
	}
	if len(filters.Protocol) > 0 {
		filterArgs.Protocols = pulumi.ToStringArray(filters.Protocol)
	}
	if len(filters.QueryTag) > 0 {
		filterArgs.QueryTags = pulumi.ToStringArray(filters.QueryTag)
	}
	if len(filters.RequestsPerSecond) > 0 {
		filterArgs.RequestsPerSeconds = pulumi.ToStringArray(filters.RequestsPerSecond)
	}
	if len(filters.Selectors) > 0 {
		filterArgs.Selectors = pulumi.ToStringArray(filters.Selectors)
	}
	if len(filters.Services) > 0 {
		filterArgs.Services = pulumi.ToStringArray(filters.Services)
	}
	if len(filters.Slo) > 0 {
		filterArgs.Slos = pulumi.ToStringArray(filters.Slo)
	}
	if len(filters.Status) > 0 {
		filterArgs.Statuses = pulumi.ToStringArray(filters.Status)
	}
	if len(filters.TargetHostname) > 0 {
		filterArgs.TargetHostnames = pulumi.ToStringArray(filters.TargetHostname)
	}
	if len(filters.TargetIp) > 0 {
		filterArgs.TargetIps = pulumi.ToStringArray(filters.TargetIp)
	}
	if len(filters.TargetZoneName) > 0 {
		filterArgs.TargetZoneNames = pulumi.ToStringArray(filters.TargetZoneName)
	}
	if len(filters.TrafficExclusions) > 0 {
		filterArgs.TrafficExclusions = pulumi.ToStringArray(filters.TrafficExclusions)
	}
	if len(filters.TunnelId) > 0 {
		filterArgs.TunnelIds = pulumi.ToStringArray(filters.TunnelId)
	}
	if len(filters.TunnelName) > 0 {
		filterArgs.TunnelNames = pulumi.ToStringArray(filters.TunnelName)
	}
	if len(filters.Type) > 0 {
		filterArgs.Types = pulumi.ToStringArray(filters.Type)
	}
	if len(filters.Where) > 0 {
		filterArgs.Wheres = pulumi.ToStringArray(filters.Where)
	}
	if len(filters.Zones) > 0 {
		filterArgs.Zones = pulumi.ToStringArray(filters.Zones)
	}

	return filterArgs
}
