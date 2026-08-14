package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudtrail"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// eventDataStore creates the CloudTrail Lake event data store and
// exports outputs.
//
// Lifecycle facts the render below depends on:
//   - deletion is REFUSED while termination protection is on (AWS
//     behavior, not a module choice) - the teardown is two steps:
//     apply with termination_protection_enabled = false, then destroy;
//   - a destroyed store lingers in PENDING_DELETION for 7 days and
//     its name stays reserved until the purge completes;
//   - "suspend" is write-only at AWS (never reported back), so it is
//     asserted on every apply and invisible to imports;
//   - an omitted selector list makes AWS materialize a default
//     all-management-events selector - the first import after an
//     omitted-selector create shows that server-side default.
func eventDataStore(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &cloudtrail.EventDataStoreArgs{
		// metadata.name is the store name on both engines (AWS: 3-128
		// characters).
		Name: pulumi.String(locals.Target.Metadata.Name),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// Rendered only on an explicit choice so the module never fights
	// the provider defaults (billing_mode EXTENDABLE_RETENTION_PRICING,
	// multi_region true, retention 2555, termination protection true).
	if spec.BillingMode != "" {
		args.BillingMode = pulumi.String(spec.BillingMode)
	}
	if spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
	}
	if spec.MultiRegionEnabled != nil {
		args.MultiRegionEnabled = pulumi.Bool(*spec.MultiRegionEnabled)
	}
	if spec.OrganizationEnabled {
		args.OrganizationEnabled = pulumi.Bool(true)
	}
	if spec.RetentionPeriodDays != 0 {
		args.RetentionPeriod = pulumi.Int(int(spec.RetentionPeriodDays))
	}
	if spec.TerminationProtectionEnabled != nil {
		args.TerminationProtectionEnabled = pulumi.Bool(*spec.TerminationProtectionEnabled)
	}
	// The provider models suspend as a nullable string ("true"/"false").
	if spec.Suspend != nil {
		args.Suspend = pulumi.String(strconv.FormatBool(*spec.Suspend))
	}

	// Ingestion scope: fine-grained field matching. AWS requires every
	// selector to carry an eventCategory condition (server-side rule).
	var selectors cloudtrail.EventDataStoreAdvancedEventSelectorArray
	for _, s := range spec.AdvancedEventSelectors {
		var fieldSelectors cloudtrail.EventDataStoreAdvancedEventSelectorFieldSelectorArray
		for _, f := range s.FieldSelectors {
			fs := &cloudtrail.EventDataStoreAdvancedEventSelectorFieldSelectorArgs{
				Field: pulumi.String(f.Field),
			}
			if len(f.Equals) > 0 {
				fs.Equals = pulumi.ToStringArray(f.Equals)
			}
			if len(f.NotEquals) > 0 {
				fs.NotEquals = pulumi.ToStringArray(f.NotEquals)
			}
			if len(f.StartsWith) > 0 {
				fs.StartsWiths = pulumi.ToStringArray(f.StartsWith)
			}
			if len(f.NotStartsWith) > 0 {
				fs.NotStartsWiths = pulumi.ToStringArray(f.NotStartsWith)
			}
			if len(f.EndsWith) > 0 {
				fs.EndsWiths = pulumi.ToStringArray(f.EndsWith)
			}
			if len(f.NotEndsWith) > 0 {
				fs.NotEndsWiths = pulumi.ToStringArray(f.NotEndsWith)
			}
			fieldSelectors = append(fieldSelectors, fs)
		}
		selArgs := &cloudtrail.EventDataStoreAdvancedEventSelectorArgs{
			FieldSelectors: fieldSelectors,
		}
		if s.Name != "" {
			selArgs.Name = pulumi.String(s.Name)
		}
		selectors = append(selectors, selArgs)
	}
	if len(selectors) > 0 {
		args.AdvancedEventSelectors = selectors
	}

	createdStore, err := cloudtrail.NewEventDataStore(ctx, "event-data-store", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create event data store")
	}

	ctx.Export(OpEventDataStoreArn, createdStore.Arn)
	return nil
}
