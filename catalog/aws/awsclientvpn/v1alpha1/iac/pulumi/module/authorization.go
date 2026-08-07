package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2clientvpn"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// authorizationRules creates the endpoint's access grants. Rules are keyed
// by (destination CIDR, grantee) -- the same identity AWS uses -- so the
// common pattern of authorizing one CIDR for several IdP groups is
// representable, and editing one grant never disturbs the others.
func authorizationRules(
	ctx *pulumi.Context,
	locals *Locals,
	provider pulumi.ProviderResource,
	createdEndpoint *ec2clientvpn.Endpoint,
) error {
	spec := locals.AwsClientVpn.Spec

	for _, rule := range spec.AuthorizationRules {
		args := &ec2clientvpn.AuthorizationRuleArgs{
			ClientVpnEndpointId: createdEndpoint.ID(),
			TargetNetworkCidr:   pulumi.String(rule.TargetNetworkCidr),
		}
		// Exactly one grantee arm (CEL-enforced): a specific IdP group or
		// every authenticated client.
		grantee := "all-groups"
		if rule.AccessGroupId != "" {
			args.AccessGroupId = pulumi.String(rule.AccessGroupId)
			grantee = rule.AccessGroupId
		} else {
			args.AuthorizeAllGroups = pulumi.Bool(true)
		}
		if rule.Description != "" {
			args.Description = pulumi.String(rule.Description)
		}

		if _, err := ec2clientvpn.NewAuthorizationRule(ctx,
			fmt.Sprintf("authorization-%s-%s", rule.TargetNetworkCidr, grantee),
			args, pulumi.Provider(provider), pulumi.Parent(createdEndpoint)); err != nil {
			return errors.Wrapf(err, "authorization rule %s for %s", rule.TargetNetworkCidr, grantee)
		}
	}

	return nil
}
