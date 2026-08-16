package module

import (
	"github.com/pkg/errors"
	awsorganizationaccountv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsorganizationaccount/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/account"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// memberAccount creates the member account with its folded contact
// and region settings, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - email, role_name, iam_user_access_to_billing, and
//     create_govcloud are creation-time facts (the first three force
//     replacement; a create_govcloud change is silently ignored by
//     AWS);
//   - the display name renames in place through the Account Management
//     API; a parent_id change MOVES the account between OUs in place;
//   - role_name has NO read API - both engines deliberately ignore
//     later changes (without that, importing an existing account would
//     plan a destructive replacement) - declared config-only in the
//     import catalog;
//   - destroy honors close_on_deletion: false REMOVES the account from
//     the organization (it survives standalone), true CLOSES it
//     (~90-day PENDING_CLOSURE, quota-limited per rolling 30 days);
//   - the contact satellites use idempotent Put APIs (the provider
//     polls until the write is visible); primary-contact delete is a
//     NO-OP (the last-written contact stays on file);
//   - region enable/disable are long operations (up to ~60 minutes
//     each way) and region delete is a NO-OP (the region keeps its
//     last state);
//   - the account-settings satellites require trusted access for AWS
//     Account Management ("account.amazonaws.com") on the
//     organization.
func memberAccount(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &organizations.AccountArgs{
		Name:  pulumi.StringPtr(spec.AccountName),
		Email: pulumi.String(spec.Email),
		// Always sent explicitly to match the Terraform module's
		// render (both engines pin the engine-behavior booleans).
		CloseOnDeletion: pulumi.Bool(spec.CloseOnDeletion),
		CreateGovcloud:  pulumi.Bool(spec.CreateGovcloud),
		Tags:            pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.ParentId.GetValue() != "" {
		args.ParentId = pulumi.StringPtr(spec.ParentId.GetValue())
	}
	if spec.RoleName != "" {
		args.RoleName = pulumi.StringPtr(spec.RoleName)
	}
	// Empty keeps AWS's default (ALLOW).
	if spec.IamUserAccessToBilling != "" {
		args.IamUserAccessToBilling = pulumi.StringPtr(spec.IamUserAccessToBilling)
	}

	createdAccount, err := organizations.NewAccount(ctx, "account", args,
		pulumi.Provider(provider),
		// role_name is write-only at AWS (no read API): ignored on
		// later changes, matching the Terraform module's lifecycle
		// ignore.
		pulumi.IgnoreChanges([]string{"roleName"}))
	if err != nil {
		return errors.Wrap(err, "create account")
	}

	// At most one alternate contact per category, keyed by the same
	// type token the provider imports them under
	// ("{account_id}/{BILLING|OPERATIONS|SECURITY}").
	alternateContacts := []struct {
		contactType string
		contact     *awsorganizationaccountv1alpha1.AwsOrganizationAccountAlternateContact
	}{
		{"BILLING", spec.GetAlternateContacts().GetBilling()},
		{"OPERATIONS", spec.GetAlternateContacts().GetOperations()},
		{"SECURITY", spec.GetAlternateContacts().GetSecurity()},
	}
	for _, entry := range alternateContacts {
		if entry.contact == nil {
			continue
		}
		_, err := account.NewAlternativeContact(ctx, "alternate-contact-"+entry.contactType,
			&account.AlternativeContactArgs{
				AccountId:            createdAccount.ID().ToStringOutput(),
				AlternateContactType: pulumi.String(entry.contactType),
				Name:                 pulumi.StringPtr(entry.contact.Name),
				Title:                pulumi.String(entry.contact.Title),
				EmailAddress:         pulumi.String(entry.contact.EmailAddress),
				PhoneNumber:          pulumi.String(entry.contact.PhoneNumber),
			}, pulumi.Provider(provider), pulumi.Parent(createdAccount))
		if err != nil {
			return errors.Wrapf(err, "alternate contact %s", entry.contactType)
		}
	}

	// The account's primary contact information. Optional leaves are
	// sent only when set - clearing one in the spec leaves the last
	// value on file at AWS (the API has no unset semantics).
	if spec.PrimaryContact != nil {
		primaryContactArgs := &account.PrimaryContactArgs{
			AccountId:    createdAccount.ID().ToStringOutput(),
			FullName:     pulumi.String(spec.PrimaryContact.FullName),
			AddressLine1: pulumi.String(spec.PrimaryContact.AddressLine_1),
			City:         pulumi.String(spec.PrimaryContact.City),
			PostalCode:   pulumi.String(spec.PrimaryContact.PostalCode),
			CountryCode:  pulumi.String(spec.PrimaryContact.CountryCode),
			PhoneNumber:  pulumi.String(spec.PrimaryContact.PhoneNumber),
		}
		if spec.PrimaryContact.CompanyName != "" {
			primaryContactArgs.CompanyName = pulumi.StringPtr(spec.PrimaryContact.CompanyName)
		}
		if spec.PrimaryContact.AddressLine_2 != "" {
			primaryContactArgs.AddressLine2 = pulumi.StringPtr(spec.PrimaryContact.AddressLine_2)
		}
		if spec.PrimaryContact.AddressLine_3 != "" {
			primaryContactArgs.AddressLine3 = pulumi.StringPtr(spec.PrimaryContact.AddressLine_3)
		}
		if spec.PrimaryContact.DistrictOrCounty != "" {
			primaryContactArgs.DistrictOrCounty = pulumi.StringPtr(spec.PrimaryContact.DistrictOrCounty)
		}
		if spec.PrimaryContact.StateOrRegion != "" {
			primaryContactArgs.StateOrRegion = pulumi.StringPtr(spec.PrimaryContact.StateOrRegion)
		}
		if spec.PrimaryContact.WebsiteUrl != "" {
			primaryContactArgs.WebsiteUrl = pulumi.StringPtr(spec.PrimaryContact.WebsiteUrl)
		}
		_, err := account.NewPrimaryContact(ctx, "primary-contact", primaryContactArgs,
			pulumi.Provider(provider), pulumi.Parent(createdAccount))
		if err != nil {
			return errors.Wrap(err, "primary contact")
		}
	}

	// Opt-in region enablement, keyed by region name (each imports as
	// "{account_id},{region_name}").
	for _, region := range spec.Regions {
		_, err := account.NewRegion(ctx, "region-"+region.RegionName,
			&account.RegionArgs{
				AccountId:  createdAccount.ID().ToStringOutput(),
				RegionName: pulumi.String(region.RegionName),
				Enabled:    pulumi.Bool(region.Enabled),
			}, pulumi.Provider(provider), pulumi.Parent(createdAccount))
		if err != nil {
			return errors.Wrapf(err, "region %s", region.RegionName)
		}
	}

	ctx.Export(OpAccountId, createdAccount.ID())
	ctx.Export(OpArn, createdAccount.Arn)
	ctx.Export(OpState, createdAccount.State)
	ctx.Export(OpGovcloudId, createdAccount.GovcloudId)
	return nil
}
