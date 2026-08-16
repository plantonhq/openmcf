package module

import (
	"sort"

	"github.com/pkg/errors"
	awsguarddutyv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsguardduty/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/guardduty"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// detector creates the GuardDuty detector and every satellite the
// spec declares, and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - AWS allows ONE detector per account per region; a pre-existing
//     detector (enabled by hand or by Organizations auto-enable) makes
//     create fail with "detector already exists";
//   - detector features and organization features are PATCHES onto the
//     detector (Create and Update are the same UpdateDetector call;
//     Delete is a no-op) - features removed from the spec are NOT
//     reverted by AWS, and upstream serializes feature writes per
//     detector under a global mutex, which is why every feature
//     resource below parents on the detector;
//   - the organization-configuration resource's delete is a no-op too:
//     destroying this component leaves the org posture as last applied
//     (taught in the GUIDE);
//   - the publishing destination's bucket POLICY and KMS key policy
//     must grant guardduty.amazonaws.com before create (the consumer's
//     contract on AwsS3Bucket / AwsKmsKey specs);
//   - changing TAGS on the publishing destination REPLACES it upstream
//     (ForceNew tags) - the destination is deliberately untagged so
//     tag sweeps never replace findings export;
//   - members and the invite accepter are the cross-account surface: a
//     member record needs the member account's root email, and the
//     accepter runs in the MEMBER account against a pending invite.
//
// Iteration over every satellite list is name-sorted for
// deterministic previews.
func detector(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	detectorArgs := &guardduty.DetectorArgs{
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}
	// Rendered only on an explicit choice so the module never fights
	// the provider default (enabled).
	if spec.Enable != nil {
		detectorArgs.Enable = pulumi.Bool(*spec.Enable)
	}
	// Left to AWS (SIX_HOURS) unless the spec sets it: members inherit
	// the administrator's value, and an explicit send on a member
	// detector would fight the org sync forever (the idempotency gate
	// would catch exactly that).
	if spec.FindingPublishingFrequency != "" {
		detectorArgs.FindingPublishingFrequency = pulumi.String(spec.FindingPublishingFrequency)
	}

	createdDetector, err := guardduty.NewDetector(ctx, "detector", detectorArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create detector")
	}

	// Protection plans, patch-keyed by feature name.
	for _, f := range sortedByName(spec.Features, func(f *awsguarddutyv1alpha1.AwsGuardDutyFeature) string { return f.Name }) {
		featureArgs := &guardduty.DetectorFeatureArgs{
			DetectorId: createdDetector.ID(),
			Name:       pulumi.String(f.Name),
			Status:     pulumi.String(enabledStatus(f.Enabled)),
		}
		var additional guardduty.DetectorFeatureAdditionalConfigurationArray
		for _, c := range f.AdditionalConfiguration {
			additional = append(additional, &guardduty.DetectorFeatureAdditionalConfigurationArgs{
				Name:   pulumi.String(c.Name),
				Status: pulumi.String(enabledStatus(c.Enabled)),
			})
		}
		if len(additional) > 0 {
			featureArgs.AdditionalConfigurations = additional
		}
		_, err := guardduty.NewDetectorFeature(ctx, "feature-"+f.Name, featureArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "enable feature %q", f.Name)
		}
	}

	// Finding filters, keyed by name (the for_each key both engines
	// share).
	for _, f := range sortedByName(spec.Filters, func(f *awsguarddutyv1alpha1.AwsGuardDutyFilter) string { return f.Name }) {
		var criteria guardduty.FilterFindingCriteriaCriterionArray
		for _, c := range f.Criteria {
			criterionArgs := &guardduty.FilterFindingCriteriaCriterionArgs{
				Field: pulumi.String(c.Field),
			}
			if len(c.Equals) > 0 {
				criterionArgs.Equals = pulumi.ToStringArray(c.Equals)
			}
			if len(c.NotEquals) > 0 {
				criterionArgs.NotEquals = pulumi.ToStringArray(c.NotEquals)
			}
			if len(c.Matches) > 0 {
				criterionArgs.Matches = pulumi.ToStringArray(c.Matches)
			}
			if len(c.NotMatches) > 0 {
				criterionArgs.NotMatches = pulumi.ToStringArray(c.NotMatches)
			}
			if c.GreaterThan != "" {
				criterionArgs.GreaterThan = pulumi.String(c.GreaterThan)
			}
			if c.GreaterThanOrEqual != "" {
				criterionArgs.GreaterThanOrEqual = pulumi.String(c.GreaterThanOrEqual)
			}
			if c.LessThan != "" {
				criterionArgs.LessThan = pulumi.String(c.LessThan)
			}
			if c.LessThanOrEqual != "" {
				criterionArgs.LessThanOrEqual = pulumi.String(c.LessThanOrEqual)
			}
			criteria = append(criteria, criterionArgs)
		}
		filterArgs := &guardduty.FilterArgs{
			DetectorId: createdDetector.ID(),
			Name:       pulumi.String(f.Name),
			Action:     pulumi.String(f.Action),
			Rank:       pulumi.Int(int(f.Rank)),
			FindingCriteria: &guardduty.FilterFindingCriteriaArgs{
				Criterions: criteria,
			},
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if f.Description != "" {
			filterArgs.Description = pulumi.String(f.Description)
		}
		_, err := guardduty.NewFilter(ctx, "filter-"+f.Name, filterArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create filter %q", f.Name)
		}
	}

	// Trusted IP lists (AWS activates at most one per detector).
	ipSetIds := pulumi.StringMap{}
	for _, s := range sortedByName(spec.IpSets, func(s *awsguarddutyv1alpha1.AwsGuardDutyIpSet) string { return s.Name }) {
		created, err := guardduty.NewIPSet(ctx, "ipset-"+s.Name, &guardduty.IPSetArgs{
			DetectorId: createdDetector.ID(),
			Name:       pulumi.String(s.Name),
			Format:     pulumi.String(s.Format),
			Location:   pulumi.String(s.Location),
			Activate:   pulumi.Bool(s.Activate),
			Tags:       pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create ip set %q", s.Name)
		}
		ipSetIds[s.Name] = created.IpSetId
	}

	// Threat intel lists.
	threatIntelSetIds := pulumi.StringMap{}
	for _, s := range sortedByName(spec.ThreatIntelSets, func(s *awsguarddutyv1alpha1.AwsGuardDutyThreatIntelSet) string { return s.Name }) {
		created, err := guardduty.NewThreatIntelSet(ctx, "threatintelset-"+s.Name, &guardduty.ThreatIntelSetArgs{
			DetectorId: createdDetector.ID(),
			Name:       pulumi.String(s.Name),
			Format:     pulumi.String(s.Format),
			Location:   pulumi.String(s.Location),
			Activate:   pulumi.Bool(s.Activate),
			Tags:       pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create threat intel set %q", s.Name)
		}
		threatIntelSetIds[s.Name] = created.ThreatIntelSetId
	}

	// Findings export to S3 (deliberately untagged - see the header).
	if spec.PublishingDestination != nil {
		created, err := guardduty.NewPublishingDestination(ctx, "publishing-destination", &guardduty.PublishingDestinationArgs{
			DetectorId:      createdDetector.ID(),
			DestinationArn:  pulumi.String(spec.PublishingDestination.BucketArn.GetValue()),
			KmsKeyArn:       pulumi.String(spec.PublishingDestination.KmsKeyArn.GetValue()),
			DestinationType: pulumi.String("S3"),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create publishing destination")
		}
		ctx.Export(OpPublishingDestinationId, created.DestinationId)
	} else {
		ctx.Export(OpPublishingDestinationId, pulumi.String(""))
	}

	// ----- ADMIN side: organization administration -----
	if spec.Organization != nil {
		if err := organization(ctx, locals, provider, createdDetector); err != nil {
			return err
		}
	}

	// ----- ADMIN side: members -----
	if err := members(ctx, locals, provider, createdDetector); err != nil {
		return err
	}

	// ----- MEMBER side: accept a pending invitation -----
	if spec.AcceptInvitationFromAccountId != "" {
		_, err := guardduty.NewInviteAccepter(ctx, "invite-accepter", &guardduty.InviteAccepterArgs{
			DetectorId:      createdDetector.ID(),
			MasterAccountId: pulumi.String(spec.AcceptInvitationFromAccountId),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "accept invitation")
		}
	}

	ctx.Export(OpDetectorId, createdDetector.ID())
	ctx.Export(OpDetectorArn, createdDetector.Arn)
	ctx.Export(OpAccountId, createdDetector.AccountId)
	ctx.Export(OpIpSetIds, ipSetIds)
	ctx.Export(OpThreatIntelSetIds, threatIntelSetIds)
	return nil
}

// organization renders the admin-side org surface: the account-global
// delegation act, the org configuration, and org-wide feature
// auto-enablement.
func organization(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdDetector *guardduty.Detector) error {
	org := locals.Spec.Organization

	// The account-global delegation act (one per organization,
	// performed from the MANAGEMENT account).
	var orgDeps []pulumi.Resource
	if org.AdminAccountId != "" {
		createdAdmin, err := guardduty.NewOrganizationAdminAccount(ctx, "org-admin-account", &guardduty.OrganizationAdminAccountArgs{
			AdminAccountId: pulumi.String(org.AdminAccountId),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "delegate organization admin account")
		}
		orgDeps = append(orgDeps, createdAdmin)
	}

	// A same-apply delegation must land before org configuration.
	createdOrgConfig, err := guardduty.NewOrganizationConfiguration(ctx, "org-configuration", &guardduty.OrganizationConfigurationArgs{
		DetectorId:                    createdDetector.ID(),
		AutoEnableOrganizationMembers: pulumi.String(org.AutoEnableOrganizationMembers),
	}, pulumi.Provider(provider), pulumi.DependsOn(orgDeps))
	if err != nil {
		return errors.Wrap(err, "configure organization")
	}

	// Organization-wide feature auto-enablement, patch-keyed by name.
	for _, f := range sortedByName(org.Features, func(f *awsguarddutyv1alpha1.AwsGuardDutyOrganizationFeature) string { return f.Name }) {
		featureArgs := &guardduty.OrganizationConfigurationFeatureArgs{
			DetectorId: createdDetector.ID(),
			Name:       pulumi.String(f.Name),
			AutoEnable: pulumi.String(f.AutoEnable),
		}
		var additional guardduty.OrganizationConfigurationFeatureAdditionalConfigurationArray
		for _, c := range f.AdditionalConfiguration {
			additional = append(additional, &guardduty.OrganizationConfigurationFeatureAdditionalConfigurationArgs{
				Name:       pulumi.String(c.Name),
				AutoEnable: pulumi.String(c.AutoEnable),
			})
		}
		if len(additional) > 0 {
			featureArgs.AdditionalConfigurations = additional
		}
		_, err := guardduty.NewOrganizationConfigurationFeature(ctx, "org-feature-"+f.Name, featureArgs,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdOrgConfig}))
		if err != nil {
			return errors.Wrapf(err, "auto-enable organization feature %q", f.Name)
		}
	}

	return nil
}

// members renders the admin-side member records and their per-member
// feature overrides.
func members(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdDetector *guardduty.Detector) error {
	for _, m := range sortedByName(locals.Spec.Members, func(m *awsguarddutyv1alpha1.AwsGuardDutyMember) string { return m.AccountId }) {
		memberArgs := &guardduty.MemberArgs{
			DetectorId: createdDetector.ID(),
			AccountId:  pulumi.String(m.AccountId),
			Email:      pulumi.String(m.Email),
		}
		if m.Invite != nil {
			memberArgs.Invite = pulumi.Bool(*m.Invite)
		}
		if m.InvitationMessage != "" {
			memberArgs.InvitationMessage = pulumi.String(m.InvitationMessage)
		}
		if m.DisableEmailNotification {
			memberArgs.DisableEmailNotification = pulumi.Bool(true)
		}
		createdMember, err := guardduty.NewMember(ctx, "member-"+m.AccountId, memberArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create member %q", m.AccountId)
		}

		// Per-member protection-plan overrides, keyed "account/feature"
		// (the two stable keys both engines share).
		for _, f := range sortedByName(m.Features, func(f *awsguarddutyv1alpha1.AwsGuardDutyMemberFeature) string { return f.Name }) {
			featureArgs := &guardduty.MemberDetectorFeatureArgs{
				DetectorId: createdDetector.ID(),
				AccountId:  pulumi.String(m.AccountId),
				Name:       pulumi.String(f.Name),
				Status:     pulumi.String(enabledStatus(f.Enabled)),
			}
			var additional guardduty.MemberDetectorFeatureAdditionalConfigurationArray
			for _, c := range f.AdditionalConfiguration {
				additional = append(additional, &guardduty.MemberDetectorFeatureAdditionalConfigurationArgs{
					Name:   pulumi.String(c.Name),
					Status: pulumi.String(enabledStatus(c.Enabled)),
				})
			}
			if len(additional) > 0 {
				featureArgs.AdditionalConfigurations = additional
			}
			_, err := guardduty.NewMemberDetectorFeature(ctx, "member-feature-"+m.AccountId+"-"+f.Name, featureArgs,
				pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdMember}))
			if err != nil {
				return errors.Wrapf(err, "set member %q feature %q", m.AccountId, f.Name)
			}
		}
	}
	return nil
}

// enabledStatus maps the spec's tri-state enabled bool onto the
// FeatureStatus vocabulary: unset means ENABLED (listing a feature
// means you want it).
func enabledStatus(enabled *bool) string {
	if enabled != nil && !*enabled {
		return "DISABLED"
	}
	return "ENABLED"
}

// sortedByName returns a name-sorted copy for deterministic previews.
func sortedByName[T any](in []T, key func(T) string) []T {
	out := append([]T{}, in...)
	sort.Slice(out, func(i, j int) bool { return key(out[i]) < key(out[j]) })
	return out
}
