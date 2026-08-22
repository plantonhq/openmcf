package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// uptimeCheck provisions the uptime check, one alert resource per spec
// alert row, and exports the check's outputs.
func uptimeCheck(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.UptimeCheck, error) {
	spec := locals.DigitalOceanUptimeCheck.Spec

	// Always declared (spec-required): the provider never reconciles a
	// DigitalOcean-defaulted region set, so an omitted value would leave
	// every subsequent plan trying to remove what the API chose.
	var regions pulumi.StringArray
	for _, region := range spec.Regions {
		regions = append(regions, pulumi.String(region))
	}

	checkArgs := &digitalocean.UptimeCheckArgs{
		Name:    pulumi.StringPtr(spec.CheckName),
		Target:  pulumi.String(spec.Target),
		Regions: regions,
	}

	// Unset defers to the provider's default, https.
	if spec.Type != "" {
		checkArgs.Type = pulumi.StringPtr(spec.Type)
	}

	// Unset defers to the provider's default, enabled.
	if spec.Enabled != nil {
		checkArgs.Enabled = pulumi.BoolPtr(spec.GetEnabled())
	}

	createdCheck, err := digitalocean.NewUptimeCheck(
		ctx,
		"check",
		checkArgs,
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean uptime check")
	}

	// One alert resource per spec row, composed on the check (the
	// standalone resource's mutable parent id is a corruption class this
	// composition makes unrepresentable). Names carry the row index so two
	// rows may share a display name without colliding.
	for idx, alert := range spec.Alerts {
		// The SDK keeps the provider's unbounded notifications list; the
		// provider reads only the first element, so exactly one is sent.
		notificationArgs := digitalocean.UptimeAlertNotificationArgs{}
		if len(alert.Notifications.Emails) > 0 {
			var emails pulumi.StringArray
			for _, email := range alert.Notifications.Emails {
				emails = append(emails, pulumi.String(email))
			}
			notificationArgs.Emails = emails
		}
		if len(alert.Notifications.Slack) > 0 {
			var slacks digitalocean.UptimeAlertNotificationSlackArray
			for _, slack := range alert.Notifications.Slack {
				slacks = append(slacks, digitalocean.UptimeAlertNotificationSlackArgs{
					Channel: pulumi.String(slack.Channel),
					// The SDK does not flag the webhook URL as secret, so
					// wrap it explicitly -- a credential must never ship
					// as plaintext state.
					Url: pulumi.ToSecret(pulumi.String(slack.Url)).(pulumi.StringOutput),
				})
			}
			notificationArgs.Slacks = slacks
		}

		alertArgs := &digitalocean.UptimeAlertArgs{
			CheckId:       createdCheck.ID(),
			Name:          pulumi.StringPtr(alert.AlertName),
			Type:          pulumi.String(alert.Type),
			Notifications: digitalocean.UptimeAlertNotificationArray{notificationArgs},
		}

		// Milliseconds for latency, days before expiry for ssl_expiry;
		// down and down_global carry no threshold.
		if alert.Threshold != nil {
			alertArgs.Threshold = pulumi.IntPtr(int(alert.GetThreshold()))
		}
		if alert.Comparison != "" {
			alertArgs.Comparison = pulumi.StringPtr(alert.Comparison)
		}
		if alert.Period != "" {
			alertArgs.Period = pulumi.StringPtr(alert.Period)
		}

		if _, err := digitalocean.NewUptimeAlert(
			ctx,
			fmt.Sprintf("alert-%d-%s", idx, alert.AlertName),
			alertArgs,
			pulumi.Provider(digitalOceanProvider),
			pulumi.Parent(createdCheck),
		); err != nil {
			return nil, errors.Wrapf(err, "failed to create uptime alert %q", alert.AlertName)
		}
	}

	ctx.Export(OpCheckId, createdCheck.ID())

	return createdCheck, nil
}
