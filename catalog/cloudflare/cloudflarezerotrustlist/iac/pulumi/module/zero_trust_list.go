package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// zeroTrustList creates the Zero Trust list. The list TYPE is immutable at
// Cloudflare -- changing it replaces the list (new ID), breaking any Gateway
// policy or posture rule referencing the old one. URL-type values are
// normalized by the API, a known upstream drift source at provider v5.23.0.
func zeroTrustList(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustList.Spec

	args := &cloudflare.ZeroTrustListArgs{
		AccountId: pulumi.String(spec.AccountId),
		Name:      pulumi.String(spec.Name),
		Type:      pulumi.String(spec.Type),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Cloudflare treats items as a SET: order-insignificant. Empty descriptions
	// are dropped rather than sent as "".
	if len(spec.Items) > 0 {
		items := make(cloudflare.ZeroTrustListItemArray, 0, len(spec.Items))
		for _, item := range spec.Items {
			itemArgs := &cloudflare.ZeroTrustListItemArgs{
				Value: pulumi.StringPtr(item.Value),
			}
			if item.Description != "" {
				itemArgs.Description = pulumi.StringPtr(item.Description)
			}
			items = append(items, itemArgs)
		}
		args.Items = items
	}

	createdList, err := cloudflare.NewZeroTrustList(
		ctx,
		"zero_trust_list",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create zero trust list")
	}

	ctx.Export(OpListId, createdList.ID())

	return nil
}
