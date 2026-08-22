package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// notificationWebhook registers the webhook destination notification
// policies deliver alerts to. Cloudflare infers the destination type from
// the URL and reports it as an output. A plain CRUD resource (real
// create/update/delete; only the account forces replacement).
//
// `secret` is WRITE-ONLY at the API: sent on create/update, never returned
// by any read -- it cannot be drift-detected, and an imported webhook has no
// secret in state.
func notificationWebhook(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareNotificationWebhook.Spec

	args := &cloudflare.NotificationPolicyWebhooksArgs{
		AccountId: pulumi.String(spec.AccountId),
		Name:      pulumi.String(spec.Name),
		Url:       pulumi.String(spec.Url),
	}

	if spec.Secret.GetValue() != "" {
		args.Secret = pulumi.StringPtr(spec.Secret.GetValue())
	}

	createdWebhook, err := cloudflare.NewNotificationPolicyWebhooks(
		ctx,
		"notification_webhook",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create notification webhook")
	}

	ctx.Export(OpWebhookId, createdWebhook.ID())
	ctx.Export(OpType, createdWebhook.Type)

	return nil
}
