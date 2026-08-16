package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ssm"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// patchBaseline creates the baseline with its folded patch groups and
// default designation, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - metadata.name is the baseline name on both engines; AWS
//     identifies the baseline as "pb-..." (the import ID);
//   - OperatingSystem is rendered only on an explicit choice (unset =
//     WINDOWS, the provider default) and forces replacement;
//   - patch groups are fully ForceNew registrations keyed by the group
//     name; a group registers with only ONE baseline per OS
//     account-wide (AWS state is the referee);
//   - the default designation is a REVERSIBLE pointer: destroying it
//     RESTORES AWS's own predefined default baseline for the OS (the
//     provider looks it up and re-registers it) - the TGW default-table
//     class, not the App Runner one-way class;
//   - if the baseline is deleted while holding the designation, the
//     provider restores AWS's default and retries the delete.
func patchBaseline(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &ssm.PatchBaselineArgs{
		// metadata.name is the baseline name on both engines.
		Name: pulumi.String(locals.Target.Metadata.Name),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.OperatingSystem != "" {
		args.OperatingSystem = pulumi.String(spec.OperatingSystem)
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	var approvalRules ssm.PatchBaselineApprovalRuleArray
	for _, r := range spec.ApprovalRules {
		ruleArgs := &ssm.PatchBaselineApprovalRuleArgs{}
		if r.ApproveAfterDays != nil {
			ruleArgs.ApproveAfterDays = pulumi.Int(int(*r.ApproveAfterDays))
		}
		if r.ApproveUntilDate != "" {
			ruleArgs.ApproveUntilDate = pulumi.String(r.ApproveUntilDate)
		}
		if r.ComplianceLevel != "" {
			ruleArgs.ComplianceLevel = pulumi.String(r.ComplianceLevel)
		}
		if r.EnableNonSecurity {
			ruleArgs.EnableNonSecurity = pulumi.Bool(true)
		}
		var filters ssm.PatchBaselineApprovalRulePatchFilterArray
		for _, f := range r.PatchFilters {
			filters = append(filters, &ssm.PatchBaselineApprovalRulePatchFilterArgs{
				Key:    pulumi.String(f.Key),
				Values: pulumi.ToStringArray(f.Values),
			})
		}
		ruleArgs.PatchFilters = filters
		approvalRules = append(approvalRules, ruleArgs)
	}
	if len(approvalRules) > 0 {
		args.ApprovalRules = approvalRules
	}

	var globalFilters ssm.PatchBaselineGlobalFilterArray
	for _, f := range spec.GlobalFilters {
		globalFilters = append(globalFilters, &ssm.PatchBaselineGlobalFilterArgs{
			Key:    pulumi.String(f.Key),
			Values: pulumi.ToStringArray(f.Values),
		})
	}
	if len(globalFilters) > 0 {
		args.GlobalFilters = globalFilters
	}

	if len(spec.ApprovedPatches) > 0 {
		args.ApprovedPatches = pulumi.ToStringArray(spec.ApprovedPatches)
	}
	if spec.ApprovedPatchesComplianceLevel != "" {
		args.ApprovedPatchesComplianceLevel = pulumi.String(spec.ApprovedPatchesComplianceLevel)
	}
	if spec.ApprovedPatchesEnableNonSecurity {
		args.ApprovedPatchesEnableNonSecurity = pulumi.Bool(true)
	}
	if len(spec.RejectedPatches) > 0 {
		args.RejectedPatches = pulumi.ToStringArray(spec.RejectedPatches)
	}
	if spec.RejectedPatchesAction != "" {
		args.RejectedPatchesAction = pulumi.String(spec.RejectedPatchesAction)
	}
	if spec.AvailableSecurityUpdatesComplianceStatus != "" {
		args.AvailableSecurityUpdatesComplianceStatus = pulumi.String(spec.AvailableSecurityUpdatesComplianceStatus)
	}

	var sources ssm.PatchBaselineSourceArray
	for _, s := range spec.Sources {
		sources = append(sources, &ssm.PatchBaselineSourceArgs{
			Name:          pulumi.String(s.Name),
			Configuration: pulumi.String(s.Configuration),
			Products:      pulumi.ToStringArray(s.Products),
		})
	}
	if len(sources) > 0 {
		args.Sources = sources
	}

	createdBaseline, err := ssm.NewPatchBaseline(ctx, "baseline", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create baseline")
	}

	// Folded patch-group registrations, keyed by group name (fully
	// ForceNew at the provider - any change replaces the registration).
	for _, group := range spec.PatchGroups {
		_, err := ssm.NewPatchGroup(ctx, "patch-group-"+group, &ssm.PatchGroupArgs{
			PatchGroup: pulumi.String(group),
			BaselineId: createdBaseline.ID(),
		}, pulumi.Provider(provider), pulumi.Parent(createdBaseline))
		if err != nil {
			return errors.Wrapf(err, "patch group %s", group)
		}
	}

	// The account/region default-baseline designation for this
	// baseline's OS. Destroying it restores AWS's own predefined
	// default (the provider records and reverts - a true revert, unlike
	// the App Runner one-way designation).
	if spec.SetAsDefaultBaseline {
		_, err := ssm.NewDefaultPatchBaseline(ctx, "default-designation", &ssm.DefaultPatchBaselineArgs{
			BaselineId: createdBaseline.ID(),
			// The provider echoes WINDOWS when the spec leaves the OS
			// unset, so the resolved value is always present.
			OperatingSystem: createdBaseline.OperatingSystem.Elem(),
		}, pulumi.Provider(provider), pulumi.Parent(createdBaseline))
		if err != nil {
			return errors.Wrap(err, "claim default-baseline designation")
		}
	}

	ctx.Export(OpBaselineId, createdBaseline.ID())
	ctx.Export(OpBaselineArn, createdBaseline.Arn)
	ctx.Export(OpOperatingSystem, createdBaseline.OperatingSystem.Elem())
	return nil
}
