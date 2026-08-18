package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// organization creates THE organization with its folded service
// access, delegated administrators, resource policy, and root-access
// management, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - feature_set upgrades apply in place (EnableAllFeatures); the
//     downgrade forces replacement, which is delete-and-recreate of
//     the ENTIRE organization;
//   - service-access principals and policy types are applied as
//     enable/disable calls diffed on update (disables first); trusted
//     access, policy types, delegated administrators, and the
//     resource policy all require ALL features (spec CELs front-load
//     this);
//   - destroy calls DeleteOrganization - the whole organization ends;
//   - the resource policy is a per-organization SINGLETON
//     (PutResourcePolicy upserts it), so the arm renders at most one;
//   - root-access management is per-organization too (its provider ID
//     is the organization ID) and requires iam.amazonaws.com trusted
//     access (a spec CEL front-loads it); destroying the arm DISABLES
//     every enabled feature;
//   - the organization resource is untaggable; the resource policy is
//     the kind's one taggable surface.
func organization(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &organizations.OrganizationArgs{}

	// Empty keeps the provider default (ALL - the level every advanced
	// arm requires).
	if spec.FeatureSet != "" {
		args.FeatureSet = pulumi.StringPtr(spec.FeatureSet)
	}
	if len(spec.AwsServiceAccessPrincipals) > 0 {
		args.AwsServiceAccessPrincipals = pulumi.ToStringArray(spec.AwsServiceAccessPrincipals)
	}
	if len(spec.EnabledPolicyTypes) > 0 {
		args.EnabledPolicyTypes = pulumi.ToStringArray(spec.EnabledPolicyTypes)
	}

	createdOrganization, err := organizations.NewOrganization(ctx, "organization", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create organization")
	}

	// Each registration names one member account as the administrator
	// for one AWS service. Both leaves are immutable - a change
	// re-registers. Keyed by the same composite the provider imports
	// them under ("{account_id}/{service_principal}").
	for _, delegatedAdministrator := range spec.DelegatedAdministrators {
		_, err := organizations.NewDelegatedAdministrator(ctx,
			"delegated-administrator-"+delegatedAdministrator.AccountId+"-"+delegatedAdministrator.ServicePrincipal,
			&organizations.DelegatedAdministratorArgs{
				AccountId:        pulumi.String(delegatedAdministrator.AccountId),
				ServicePrincipal: pulumi.String(delegatedAdministrator.ServicePrincipal),
			}, pulumi.Provider(provider), pulumi.Parent(createdOrganization))
		if err != nil {
			return errors.Wrapf(err, "delegated administrator %s/%s",
				delegatedAdministrator.AccountId, delegatedAdministrator.ServicePrincipal)
		}
	}

	// The organization's single resource-based delegation policy. The
	// spec carries the document structured; it is serialized to JSON
	// here.
	resourcePolicyId := pulumi.String("").ToStringOutput()
	if spec.ResourcePolicy != nil {
		resourcePolicyJson, err := json.Marshal(spec.ResourcePolicy.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal resource policy")
		}
		createdResourcePolicy, err := organizations.NewResourcePolicy(ctx, "resource-policy",
			&organizations.ResourcePolicyArgs{
				Content: pulumi.String(string(resourcePolicyJson)),
				Tags:    pulumi.ToStringMap(locals.AwsTags),
			}, pulumi.Provider(provider), pulumi.Parent(createdOrganization))
		if err != nil {
			return errors.Wrap(err, "create resource policy")
		}
		resourcePolicyId = createdResourcePolicy.ID().ToStringOutput()
	}

	// Centralized root-access management: the org-wide IAM features.
	// The resource depends on the organization's trusted-access list
	// carrying iam.amazonaws.com, so it parents the organization and
	// applies after it.
	if spec.RootAccessManagement != nil {
		if _, err := iam.NewOrganizationsFeatures(ctx, "root-access-management",
			&iam.OrganizationsFeaturesArgs{
				EnabledFeatures: pulumi.ToStringArray(spec.RootAccessManagement.EnabledFeatures),
			}, pulumi.Provider(provider), pulumi.Parent(createdOrganization)); err != nil {
			return errors.Wrap(err, "enable root-access management features")
		}
	}

	ctx.Export(OpOrganizationId, createdOrganization.ID())
	ctx.Export(OpArn, createdOrganization.Arn)
	ctx.Export(OpManagementAccountId, createdOrganization.MasterAccountId)
	ctx.Export(OpManagementAccountArn, createdOrganization.MasterAccountArn)
	ctx.Export(OpManagementAccountEmail, createdOrganization.MasterAccountEmail)
	ctx.Export(OpRootId, createdOrganization.Roots.Index(pulumi.Int(0)).Id())
	ctx.Export(OpResourcePolicyId, resourcePolicyId)
	return nil
}
