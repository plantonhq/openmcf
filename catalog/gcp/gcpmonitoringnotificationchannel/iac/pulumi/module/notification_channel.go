package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/monitoring"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// notificationChannel provisions the Cloud Monitoring notification channel —
// the delivery endpoint alert policies notify when incidents open or close.
//
// Two label surfaces exist on this resource and they must never be
// conflated: the provider's `labels` argument is the TYPE-SPECIFIC channel
// configuration (an email address, a Slack channel name) fed from
// spec.channel_labels, while `user_labels` is freeform user metadata fed
// from spec.labels merged with the platform attribution labels. Credentials
// ride the separate sensitive_labels block, which the provider stores
// API-side (the API redacts them on read; the provider suppresses the
// perpetual diff that would otherwise cause).
//
// `enabled` is sent EXPLICITLY on every apply: it is Optional in the
// provider with a server default of true, and a spec transition
// true -> false must reach the API rather than being omitted (the
// send-true-or-omit class silently no-ops such transitions).
func notificationChannel(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpMonitoringNotificationChannel.Spec

	// Enable the Cloud Monitoring API so a fresh project can host the
	// channel. disable_on_destroy stays false (the provider default):
	// tearing down one channel must never disable monitoring for
	// everything else in the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("monitoring.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"notifchan-monitoring.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable monitoring.googleapis.com api")
	}

	args := &monitoring.NotificationChannelArgs{
		Type:        pulumi.String(spec.Type),
		DisplayName: pulumi.String(locals.DisplayName),
		// Explicit send — see the function comment. Unset optional bool
		// reads false only when the platform default middleware did not
		// apply; the proto default is true.
		Enabled:     pulumi.Bool(spec.Enabled == nil || spec.GetEnabled()),
		ForceDelete: pulumi.Bool(spec.ForceDelete),
		UserLabels:  pulumi.ToStringMap(locals.GcpLabels),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// The per-type channel configuration — only sent when the spec carries
	// keys, so channel types with no config (e.g. some webhook setups) get
	// the API default rather than an empty map.
	if len(spec.ChannelLabels) > 0 {
		args.Labels = pulumi.ToStringMap(spec.ChannelLabels)
	}

	// Credentials ride the dedicated sensitive block; each field is only
	// sent when set so the API never receives empty-string credentials.
	if spec.SensitiveLabels != nil {
		sensitiveArgs := &monitoring.NotificationChannelSensitiveLabelsArgs{}
		if spec.SensitiveLabels.AuthToken != "" {
			sensitiveArgs.AuthToken = pulumi.StringPtr(spec.SensitiveLabels.AuthToken)
		}
		if spec.SensitiveLabels.Password != "" {
			sensitiveArgs.Password = pulumi.StringPtr(spec.SensitiveLabels.Password)
		}
		if spec.SensitiveLabels.ServiceKey != "" {
			sensitiveArgs.ServiceKey = pulumi.StringPtr(spec.SensitiveLabels.ServiceKey)
		}
		args.SensitiveLabels = sensitiveArgs
	}

	// Unset defers to the provider default (DELETE).
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project. Leaving Project unset lets the gcp
	// provider resolve its own project; an empty string would be sent
	// verbatim and rejected.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdChannel, err := monitoring.NewNotificationChannel(ctx, "notification-channel", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create notification channel")
	}

	ctx.Export(OpChannelName, createdChannel.Name)
	ctx.Export(OpVerificationStatus, createdChannel.VerificationStatus)

	return nil
}
