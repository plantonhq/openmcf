package verify

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/pkg/errors"
)

// Organizations is a GLOBAL service: every verifier below ignores the
// lane's region (the client resolves the service's global endpoint).
// All four verifiers require the caller to be the organization's
// management account - the same wall the kinds' recorded live
// deferrals document.

// organizationVerifier verifies THE organization by its "o-..." ID.
type organizationVerifier struct{}

func (v *organizationVerifier) IDOutputKey() string { return "organization_id" }

func (v *organizationVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := organizations.NewFromConfig(cfg)
	out, err := client.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return errors.Wrapf(err, "describe organization %s", id)
	}
	if out.Organization == nil || aws.ToString(out.Organization.Id) != id {
		return errors.Errorf("organization %s not found (the account belongs to %q)", id, aws.ToString(out.Organization.Id))
	}
	return nil
}

func (v *organizationVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := organizations.NewFromConfig(cfg)
	out, err := client.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		// The account no longer belongs to ANY organization - the
		// deleted-organization terminal state.
		var notInUse *orgtypes.AWSOrganizationsNotInUseException
		if errors.As(err, &notInUse) {
			return nil
		}
		return errors.Wrapf(err, "describe organization %s during absence check", id)
	}
	if out.Organization != nil && aws.ToString(out.Organization.Id) == id {
		return errors.Errorf("organization %s still exists", id)
	}
	return nil
}

// organizationalUnitVerifier verifies an OU by its "ou-..." ID.
type organizationalUnitVerifier struct{}

func (v *organizationalUnitVerifier) IDOutputKey() string { return "ou_id" }

func (v *organizationalUnitVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := organizations.NewFromConfig(cfg)
	out, err := client.DescribeOrganizationalUnit(ctx, &organizations.DescribeOrganizationalUnitInput{
		OrganizationalUnitId: aws.String(id),
	})
	if err != nil {
		return errors.Wrapf(err, "describe organizational unit %s", id)
	}
	if out.OrganizationalUnit == nil || aws.ToString(out.OrganizationalUnit.Id) != id {
		return errors.Errorf("organizational unit %s not found", id)
	}
	return nil
}

func (v *organizationalUnitVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := organizations.NewFromConfig(cfg)
	_, err := client.DescribeOrganizationalUnit(ctx, &organizations.DescribeOrganizationalUnitInput{
		OrganizationalUnitId: aws.String(id),
	})
	if err == nil {
		return errors.Errorf("organizational unit %s still exists", id)
	}
	var notFound *orgtypes.OrganizationalUnitNotFoundException
	if errors.As(err, &notFound) {
		return nil
	}
	// The whole organization being gone also proves the OU is gone
	// (composed destroys tear down the fixture chain).
	var notInUse *orgtypes.AWSOrganizationsNotInUseException
	if errors.As(err, &notInUse) {
		return nil
	}
	return errors.Wrapf(err, "describe organizational unit %s during absence check", id)
}

// organizationAccountVerifier verifies a member account by its
// 12-digit ID.
type organizationAccountVerifier struct{}

func (v *organizationAccountVerifier) IDOutputKey() string { return "account_id" }

func (v *organizationAccountVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := organizations.NewFromConfig(cfg)
	out, err := client.DescribeAccount(ctx, &organizations.DescribeAccountInput{AccountId: aws.String(id)})
	if err != nil {
		return errors.Wrapf(err, "describe account %s", id)
	}
	if out.Account == nil || aws.ToString(out.Account.Id) != id {
		return errors.Errorf("account %s not found", id)
	}
	if out.Account.State != orgtypes.AccountStateActive {
		return errors.Errorf("account %s is %s, want ACTIVE", id, out.Account.State)
	}
	return nil
}

func (v *organizationAccountVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := organizations.NewFromConfig(cfg)
	out, err := client.DescribeAccount(ctx, &organizations.DescribeAccountInput{AccountId: aws.String(id)})
	if err == nil {
		// A closing account lingers in PENDING_CLOSURE/SUSPENDED for
		// ~90 days - that IS the destroy path's terminal state (the
		// provider treats CLOSED as gone; removal-only makes the
		// account invisible to the org immediately).
		if out.Account != nil && out.Account.State != orgtypes.AccountStateActive {
			return nil
		}
		return errors.Errorf("account %s still exists and is ACTIVE", id)
	}
	var notFound *orgtypes.AccountNotFoundException
	if errors.As(err, &notFound) {
		return nil
	}
	var notInUse *orgtypes.AWSOrganizationsNotInUseException
	if errors.As(err, &notInUse) {
		return nil
	}
	return errors.Wrapf(err, "describe account %s during absence check", id)
}

// organizationPolicyVerifier verifies a policy by its "p-..." ID.
type organizationPolicyVerifier struct{}

func (v *organizationPolicyVerifier) IDOutputKey() string { return "policy_id" }

func (v *organizationPolicyVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := organizations.NewFromConfig(cfg)
	out, err := client.DescribePolicy(ctx, &organizations.DescribePolicyInput{PolicyId: aws.String(id)})
	if err != nil {
		return errors.Wrapf(err, "describe policy %s", id)
	}
	if out.Policy == nil || out.Policy.PolicySummary == nil || aws.ToString(out.Policy.PolicySummary.Id) != id {
		return errors.Errorf("policy %s not found", id)
	}
	// AWS-managed policies can never be this kind's surface - a
	// managed hit means the ID belongs to the wrong object.
	if out.Policy.PolicySummary.AwsManaged {
		return errors.Errorf("policy %s is AWS-managed - not a customer policy", id)
	}
	return nil
}

func (v *organizationPolicyVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, _ string) error {
	client := organizations.NewFromConfig(cfg)
	_, err := client.DescribePolicy(ctx, &organizations.DescribePolicyInput{PolicyId: aws.String(id)})
	if err == nil {
		return errors.Errorf("policy %s still exists", id)
	}
	var notFound *orgtypes.PolicyNotFoundException
	if errors.As(err, &notFound) {
		return nil
	}
	var notInUse *orgtypes.AWSOrganizationsNotInUseException
	if errors.As(err, &notInUse) {
		return nil
	}
	// Some deletes surface as a generic 400 naming the policy - treat
	// an explicit not-found message as absence, everything else as a
	// real failure.
	if strings.Contains(err.Error(), "PolicyNotFound") {
		return nil
	}
	return errors.Wrapf(err, "describe policy %s during absence check", id)
}
