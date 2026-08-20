package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/backup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// settings applies the Backup settings arms and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - BOTH resources' deletes are no-ops at the provider - destroying
//     this component changes NOTHING at AWS; the last-applied settings
//     stay in effect (taught on the spec arms);
//   - both are Required full maps at the provider: AWS returns every
//     supported key on read, so a key/type missing from the spec map
//     shows as a perpetual plan difference - the spec teaches listing
//     every key deliberately;
//   - the global arm is ACCOUNT-wide (its AWS identity is the account
//     ID, no region involved); the region arm's identity is the
//     region.
func settings(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	if spec.Global != nil {
		if _, err := backup.NewGlobalSettings(ctx, "global-settings", &backup.GlobalSettingsArgs{
			GlobalSettings: pulumi.ToStringMap(spec.Global.Settings),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "apply global settings")
		}
	}

	if spec.RegionSettings != nil {
		args := &backup.RegionSettingsArgs{
			ResourceTypeOptInPreference: pulumi.ToBoolMap(spec.RegionSettings.ResourceTypeOptInPreference),
		}
		// Rendered only when set - once set at AWS, the preference
		// cannot be cleared back to unset, only flipped per type.
		if len(spec.RegionSettings.ResourceTypeManagementPreference) > 0 {
			args.ResourceTypeManagementPreference = pulumi.ToBoolMap(spec.RegionSettings.ResourceTypeManagementPreference)
		}
		if _, err := backup.NewRegionSettings(ctx, "region-settings", args,
			pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "apply region settings")
		}
	}

	// The settings resources expose no ARNs - the identities are the
	// outputs.
	callerIdentity, err := aws.GetCallerIdentity(ctx, &aws.GetCallerIdentityArgs{}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "get caller identity")
	}

	ctx.Export(OpAccountId, pulumi.String(callerIdentity.AccountId))
	ctx.Export(OpRegion, pulumi.String(locals.Spec.Region))
	return nil
}
