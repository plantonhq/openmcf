package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// organization applies the Zero Trust organization configuration. Cloudflare
// has no create call for an organization: both create and update are the
// same PUT (an upsert of the singleton the account or zone already carries),
// and DESTROY IS A NO-OP -- deleting this resource abandons the live
// configuration exactly as last applied. Unset spec fields are never sent,
// leaving the live value untouched.
//
// The Access service-key rotation cadence is a separate Cloudflare surface
// with the same singleton/upsert/no-op-destroy lifecycle, folded into this
// component: it deploys only when the spec declares
// key_rotation_interval_days (account scope only -- CEL enforces that).
func organization(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustOrganization.Spec

	args := &cloudflare.ZeroTrustOrganizationArgs{}

	// Exactly one scope is set (CEL-enforced). The provider prefers
	// account_id when both are present -- our spec never produces that.
	if spec.AccountId != "" {
		args.AccountId = pulumi.String(spec.AccountId)
	}
	if spec.ZoneId.GetValue() != "" {
		args.ZoneId = pulumi.String(spec.ZoneId.GetValue())
	}

	if spec.AuthDomain != "" {
		args.AuthDomain = pulumi.String(spec.AuthDomain)
	}
	if spec.Name != "" {
		args.Name = pulumi.String(spec.Name)
	}
	if spec.SessionDuration != "" {
		args.SessionDuration = pulumi.String(spec.SessionDuration)
	}
	if spec.WarpAuthSessionDuration != "" {
		args.WarpAuthSessionDuration = pulumi.String(spec.WarpAuthSessionDuration)
	}
	if spec.UserSeatExpirationInactiveTime != "" {
		args.UserSeatExpirationInactiveTime = pulumi.String(spec.UserSeatExpirationInactiveTime)
	}
	if spec.DenyUnmatchedRequests != nil {
		args.DenyUnmatchedRequests = pulumi.BoolPtr(spec.GetDenyUnmatchedRequests())
	}
	if len(spec.DenyUnmatchedRequestsExemptedZoneNames) > 0 {
		args.DenyUnmatchedRequestsExemptedZoneNames = pulumi.ToStringArray(spec.DenyUnmatchedRequestsExemptedZoneNames)
	}
	if spec.AllowAuthenticateViaWarp != nil {
		args.AllowAuthenticateViaWarp = pulumi.BoolPtr(spec.GetAllowAuthenticateViaWarp())
	}
	if spec.AutoRedirectToIdentity != nil {
		args.AutoRedirectToIdentity = pulumi.BoolPtr(spec.GetAutoRedirectToIdentity())
	}
	if spec.MfaRequiredForAllApps != nil {
		args.MfaRequiredForAllApps = pulumi.BoolPtr(spec.GetMfaRequiredForAllApps())
	}
	if spec.IsUiReadOnly != nil {
		args.IsUiReadOnly = pulumi.BoolPtr(spec.GetIsUiReadOnly())
	}
	if spec.UiReadOnlyToggleReason != "" {
		args.UiReadOnlyToggleReason = pulumi.String(spec.UiReadOnlyToggleReason)
	}

	if spec.CustomPages != nil {
		customPages := cloudflare.ZeroTrustOrganizationCustomPagesArgs{}
		if spec.CustomPages.Forbidden != "" {
			customPages.Forbidden = pulumi.String(spec.CustomPages.Forbidden)
		}
		if spec.CustomPages.IdentityDenied != "" {
			customPages.IdentityDenied = pulumi.String(spec.CustomPages.IdentityDenied)
		}
		args.CustomPages = customPages
	}

	if spec.LoginDesign != nil {
		loginDesign := cloudflare.ZeroTrustOrganizationLoginDesignArgs{}
		if spec.LoginDesign.BackgroundColor != "" {
			loginDesign.BackgroundColor = pulumi.String(spec.LoginDesign.BackgroundColor)
		}
		if spec.LoginDesign.TextColor != "" {
			loginDesign.TextColor = pulumi.String(spec.LoginDesign.TextColor)
		}
		if spec.LoginDesign.LogoPath != "" {
			loginDesign.LogoPath = pulumi.String(spec.LoginDesign.LogoPath)
		}
		if spec.LoginDesign.HeaderText != "" {
			loginDesign.HeaderText = pulumi.String(spec.LoginDesign.HeaderText)
		}
		if spec.LoginDesign.FooterText != "" {
			loginDesign.FooterText = pulumi.String(spec.LoginDesign.FooterText)
		}
		args.LoginDesign = loginDesign
	}

	if spec.MfaConfig != nil {
		mfaConfig := cloudflare.ZeroTrustOrganizationMfaConfigArgs{}
		if len(spec.MfaConfig.AllowedAuthenticators) > 0 {
			mfaConfig.AllowedAuthenticators = pulumi.ToStringArray(spec.MfaConfig.AllowedAuthenticators)
		}
		if spec.MfaConfig.SessionDuration != "" {
			mfaConfig.SessionDuration = pulumi.String(spec.MfaConfig.SessionDuration)
		}
		if spec.MfaConfig.AmrMatchingSessionDuration != "" {
			mfaConfig.AmrMatchingSessionDuration = pulumi.String(spec.MfaConfig.AmrMatchingSessionDuration)
		}
		if spec.MfaConfig.RequiredAaguids != "" {
			mfaConfig.RequiredAaguids = pulumi.String(spec.MfaConfig.RequiredAaguids)
		}
		args.MfaConfig = mfaConfig
	}

	if spec.MfaSshPivKeyRequirements != nil {
		pivReqs := cloudflare.ZeroTrustOrganizationMfaSshPivKeyRequirementsArgs{}
		if spec.MfaSshPivKeyRequirements.PinPolicy != "" {
			pivReqs.PinPolicy = pulumi.String(spec.MfaSshPivKeyRequirements.PinPolicy)
		}
		if spec.MfaSshPivKeyRequirements.TouchPolicy != "" {
			pivReqs.TouchPolicy = pulumi.String(spec.MfaSshPivKeyRequirements.TouchPolicy)
		}
		if len(spec.MfaSshPivKeyRequirements.SshKeyType) > 0 {
			pivReqs.SshKeyTypes = pulumi.ToStringArray(spec.MfaSshPivKeyRequirements.SshKeyType)
		}
		if len(spec.MfaSshPivKeyRequirements.SshKeySize) > 0 {
			sizes := pulumi.IntArray{}
			for _, size := range spec.MfaSshPivKeyRequirements.SshKeySize {
				sizes = append(sizes, pulumi.Int(int(size)))
			}
			pivReqs.SshKeySizes = sizes
		}
		if spec.MfaSshPivKeyRequirements.RequireFipsDevice != nil {
			pivReqs.RequireFipsDevice = pulumi.BoolPtr(spec.MfaSshPivKeyRequirements.GetRequireFipsDevice())
		}
		args.MfaSshPivKeyRequirements = pivReqs
	}

	createdOrganization, err := cloudflare.NewZeroTrustOrganization(
		ctx,
		"organization",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to apply zero trust organization")
	}

	// The folded service-key rotation cadence: its own provider resource on
	// the same account, deployed only when declared. Ordered after the
	// organization so a fresh Zero Trust account applies the org first.
	if spec.KeyRotationIntervalDays != nil {
		_, err := cloudflare.NewZeroTrustAccessKeyConfiguration(
			ctx,
			"key_configuration",
			&cloudflare.ZeroTrustAccessKeyConfigurationArgs{
				AccountId: pulumi.String(spec.AccountId),
				// The provider models the day count as a float (a Stainless
				// artifact of the OpenAPI number type); the spec keeps the
				// honest int32.
				KeyRotationIntervalDays: pulumi.Float64(float64(spec.GetKeyRotationIntervalDays())),
			},
			pulumi.Provider(cloudflareProvider),
			pulumi.DependsOn([]pulumi.Resource{createdOrganization}),
		)
		if err != nil {
			return errors.Wrap(err, "failed to apply access key configuration")
		}
	}

	ctx.Export(OpAuthDomain, createdOrganization.AuthDomain)
	ctx.Export(OpAccountId, pulumi.String(spec.AccountId))

	return nil
}
