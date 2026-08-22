package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// monitorAlert provisions the alert policy and exports its outputs.
func monitorAlert(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.MonitorAlert, error) {
	spec := locals.DigitalOceanMonitorAlert.Spec

	// The provider takes one untyped entities list; the spec splits it
	// into three typed reference lists. References are resolved to literal
	// ids before the module runs, so the lists merge back here. Spec
	// validation already guarantees only the metric family's own list is
	// populated.
	var entities pulumi.StringArray
	for _, ref := range spec.DropletIds {
		entities = append(entities, pulumi.String(ref.GetValue()))
	}
	for _, ref := range spec.LoadBalancerIds {
		entities = append(entities, pulumi.String(ref.GetValue()))
	}
	for _, ref := range spec.DatabaseClusterIds {
		entities = append(entities, pulumi.String(ref.GetValue()))
	}

	// The SDK flattens the provider's one-element alerts block to a single
	// object and pluralizes the channel lists (Emails/Slacks).
	alertsArgs := digitalocean.MonitorAlertAlertsArgs{}
	if len(spec.Alerts.Emails) > 0 {
		var emails pulumi.StringArray
		for _, email := range spec.Alerts.Emails {
			emails = append(emails, pulumi.String(email))
		}
		alertsArgs.Emails = emails
	}
	if len(spec.Alerts.Slack) > 0 {
		var slacks digitalocean.MonitorAlertAlertsSlackArray
		for _, slack := range spec.Alerts.Slack {
			slacks = append(slacks, digitalocean.MonitorAlertAlertsSlackArgs{
				Channel: pulumi.String(slack.Channel),
				// The SDK does not flag the webhook URL as secret, so wrap
				// it explicitly -- a credential must never ship as
				// plaintext state.
				Url: pulumi.ToSecret(pulumi.String(slack.Url)).(pulumi.StringOutput),
			})
		}
		alertsArgs.Slacks = slacks
	}

	alertArgs := &digitalocean.MonitorAlertArgs{
		Description: pulumi.String(spec.Description),
		Type:        pulumi.String(spec.MetricType),
		Compare:     pulumi.String(spec.Compare),

		// DigitalOcean stores the threshold as a 32-bit float; more than 7
		// significant digits are truncated server-side.
		Value:  pulumi.Float64(spec.Value),
		Window: pulumi.String(spec.Window),
		Alerts: alertsArgs,
	}

	// Unset defers to the provider's default, enabled; the provider sends
	// the value as a pointer, so an explicit false is transmitted.
	if spec.Enabled != nil {
		alertArgs.Enabled = pulumi.BoolPtr(spec.GetEnabled())
	}

	// An omitted entities set is simply not sent; DigitalOcean then
	// resolves targets from tags.
	if len(entities) > 0 {
		alertArgs.Entities = entities
	}

	if len(spec.Tags) > 0 {
		var tags pulumi.StringArray
		for _, tag := range spec.Tags {
			tags = append(tags, pulumi.String(tag))
		}
		alertArgs.Tags = tags
	}

	createdAlert, err := digitalocean.NewMonitorAlert(
		ctx,
		"alert",
		alertArgs,
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean monitor alert")
	}

	// The resource id IS the policy UUID; the provider's own uuid
	// attribute is declared but never populated at the pinned version.
	ctx.Export(OpAlertId, createdAlert.ID())

	return createdAlert, nil
}
