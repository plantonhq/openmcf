package module

import (
	"github.com/pkg/errors"
	cloudflarednszonev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarednszone/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// subscription applies a plan subscription to the zone. Creating a paid rate
// plan is a real billing action (the deploying token needs Billing Write
// scope). The spec's flat rate_plan/scope pair maps into the provider's nested
// rate_plan object; unset levers are left nil so Cloudflare's defaults stay in
// effect.
func subscription(
	ctx *pulumi.Context,
	resourceName string,
	zone *cloudflare.Zone,
	spec *cloudflarednszonev1alpha1.CloudflareDnsZoneSubscription,
	cloudflareProvider *cloudflare.Provider,
) error {
	args := &cloudflare.ZoneSubscriptionArgs{
		ZoneId: zone.ID(),
	}
	if spec.Frequency != "" {
		args.Frequency = pulumi.String(spec.Frequency)
	}
	if spec.RatePlan != "" {
		ratePlan := cloudflare.ZoneSubscriptionRatePlanArgs{
			Id: pulumi.String(spec.RatePlan),
		}
		if spec.Scope != "" {
			ratePlan.Scope = pulumi.String(spec.Scope)
		}
		args.RatePlan = ratePlan
	}

	if _, err := cloudflare.NewZoneSubscription(
		ctx,
		resourceName+"-subscription",
		args,
		pulumi.Provider(cloudflareProvider),
	); err != nil {
		return errors.Wrap(err, "failed to apply zone subscription")
	}
	return nil
}
