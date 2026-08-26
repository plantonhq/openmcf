package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// canonicalFqdn appends the trailing dot AWS stores on every firewall
// domain, so the state matches the read-back echo whichever form the
// manifest authored.
func canonicalFqdn(domain string) string {
	if strings.HasSuffix(domain, ".") {
		return domain
	}
	return domain + "."
}

func canonicalFqdns(domains []string) []string {
	out := make([]string, len(domains))
	for i, d := range domains {
		out[i] = canonicalFqdn(d)
	}
	return out
}

// firewall creates the rule group, its domain lists, rules, and VPC
// associations, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the group itself is a name-and-tags container (name ForceNew,
//     update path is tags-only);
//   - domain list contents push through a separate update call after
//     create - a partially failed import surfaces as a retry error,
//     not silent success;
//   - a rule's match source (domain list vs threat protection) is
//     ForceNew either way; action, priority, and the block-response
//     shape update in place;
//   - BlockOverrideDnsType has exactly one legal value (CNAME), so the
//     module pins it whenever the response is OVERRIDE - a dead knob
//     is never surfaced;
//   - associations update priority/mutation_protection in place.
func firewall(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	createdRuleGroup, err := route53.NewResolverFirewallRuleGroup(ctx, "rule-group",
		&route53.ResolverFirewallRuleGroupArgs{
			Name: pulumi.String(locals.Target.Metadata.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create rule group")
	}

	// Owned domain lists, keyed by list name. AWS stores every firewall
	// domain as a trailing-dot FQDN and echoes that form on read, with
	// no provider-side diff suppression - the modules compose the
	// canonical form so a bare-authored domain never re-plans forever
	// (live-caught 2026-08-26 on block_override_domain; same storage
	// rule here).
	createdDomainLists := map[string]*route53.ResolverFirewallDomainList{}
	domainListIds := pulumi.StringMap{}
	for _, domainList := range spec.DomainLists {
		args := &route53.ResolverFirewallDomainListArgs{
			Name: pulumi.String(domainList.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if len(domainList.Domains) > 0 {
			args.Domains = pulumi.ToStringArray(canonicalFqdns(domainList.Domains))
		}

		createdDomainList, err := route53.NewResolverFirewallDomainList(ctx,
			fmt.Sprintf("domain-list-%s", domainList.Name), args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create domain list %s", domainList.Name)
		}
		createdDomainLists[domainList.Name] = createdDomainList
		domainListIds[domainList.Name] = createdDomainList.ID().ToStringOutput()
	}

	// Filtering rules, keyed by rule name. The match source resolves
	// to an owned list's generated ID, a literal external/managed
	// list ID, or the threat-protection arm. Each rule's match
	// identity (domain list ID or threat-protection ID) is exported
	// keyed by rule name - the second half of its composite import ID.
	ruleMatchIds := pulumi.StringMap{}
	for _, rule := range spec.Rules {
		args := &route53.ResolverFirewallRuleArgs{
			Name:                pulumi.String(rule.Name),
			FirewallRuleGroupId: createdRuleGroup.ID(),
			Priority:            pulumi.Int(int(rule.Priority)),
			Action:              pulumi.String(rule.Action),
		}
		if rule.DomainListName != "" {
			args.FirewallDomainListId = createdDomainLists[rule.DomainListName].ID()
		} else if rule.DomainListId.GetValue() != "" {
			args.FirewallDomainListId = pulumi.String(rule.DomainListId.GetValue())
		}
		if rule.DnsThreatProtection != "" {
			args.DnsThreatProtection = pulumi.String(rule.DnsThreatProtection)
		}
		if rule.ConfidenceThreshold != "" {
			args.ConfidenceThreshold = pulumi.String(rule.ConfidenceThreshold)
		}
		if rule.BlockResponse != "" {
			args.BlockResponse = pulumi.String(rule.BlockResponse)
		}
		// AWS stores the override domain as a trailing-dot FQDN and
		// echoes it back, so a bare-authored value re-plans forever -
		// compose the canonical form (live-caught 2026-08-26; upstream
		// has no diff suppression). PARITY: the Terraform module
		// normalizes identically.
		if rule.BlockOverrideDomain != "" {
			args.BlockOverrideDomain = pulumi.String(canonicalFqdn(rule.BlockOverrideDomain))
		}
		if rule.BlockOverrideTtl != nil {
			args.BlockOverrideTtl = pulumi.Int(int(*rule.BlockOverrideTtl))
		}
		// The one legal override record type - module-owned constant.
		if rule.BlockResponse == "OVERRIDE" {
			args.BlockOverrideDnsType = pulumi.String("CNAME")
		}
		if rule.FirewallDomainRedirectionAction != "" {
			args.FirewallDomainRedirectionAction = pulumi.String(rule.FirewallDomainRedirectionAction)
		}
		if rule.QType != "" {
			args.QType = pulumi.String(rule.QType)
		}

		createdRule, err := route53.NewResolverFirewallRule(ctx,
			fmt.Sprintf("rule-%s", rule.Name), args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create rule %s", rule.Name)
		}
		if rule.DnsThreatProtection != "" {
			ruleMatchIds[rule.Name] = createdRule.FirewallThreatProtectionId
		} else {
			ruleMatchIds[rule.Name] = createdRule.FirewallDomainListId.Elem()
		}
	}

	// VPC associations, keyed by association name.
	associationIds := pulumi.StringMap{}
	for _, association := range spec.VpcAssociations {
		args := &route53.ResolverFirewallRuleGroupAssociationArgs{
			Name:                pulumi.String(association.Name),
			FirewallRuleGroupId: createdRuleGroup.ID(),
			VpcId:               pulumi.String(association.VpcId.GetValue()),
			Priority:            pulumi.Int(int(association.Priority)),
			Tags:                pulumi.ToStringMap(locals.AwsTags),
		}
		if association.MutationProtection != "" {
			args.MutationProtection = pulumi.String(association.MutationProtection)
		}

		createdAssociation, err := route53.NewResolverFirewallRuleGroupAssociation(ctx,
			fmt.Sprintf("association-%s", association.Name), args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "associate to vpc %s", association.Name)
		}
		associationIds[association.Name] = createdAssociation.ID().ToStringOutput()
	}

	ctx.Export(OpRuleGroupId, createdRuleGroup.ID())
	ctx.Export(OpRuleGroupArn, createdRuleGroup.Arn)
	ctx.Export(OpShareStatus, createdRuleGroup.ShareStatus)
	ctx.Export(OpDomainListIds, domainListIds)
	ctx.Export(OpAssociationIds, associationIds)
	ctx.Export(OpRuleMatchIds, ruleMatchIds)
	return nil
}
