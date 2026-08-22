package verify

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
	resolvertypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"
	pkgerrors "github.com/pkg/errors"
)

// Verifiers for the Route 53 Resolver family: the endpoint (with its
// folded forwarding rules and their VPC associations), the DNS
// Firewall rule group (with its domain lists and VPC associations),
// and the query logging configuration (with its VPC associations).
// Satellites are walked through the outputs path over the keyed maps
// the modules export.

func resolverClientForRegion(cfg aws.Config, region string) *route53resolver.Client {
	return route53resolver.NewFromConfig(cfg, func(o *route53resolver.Options) {
		if region != "" {
			o.Region = region
		}
	})
}

func isResolverNotFound(err error) bool {
	var notFound *resolvertypes.ResourceNotFoundException
	return errors.As(err, &notFound)
}

// --- AwsRoute53ResolverEndpoint -------------------------------------------

// resolverEndpointVerifier asserts the endpoint is OPERATIONAL, every
// folded rule exists, and every rule association is COMPLETE (the
// association can fail asynchronously - mere existence is not enough).
type resolverEndpointVerifier struct{}

func (*resolverEndpointVerifier) IDOutputKey() string { return "endpoint_id" }

func (v *resolverEndpointVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := resolverClientForRegion(cfg, region).GetResolverEndpoint(ctx, &route53resolver.GetResolverEndpointInput{
		ResolverEndpointId: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsroute53resolverendpoint verify-exists failed for %q", id)
	}
	if out.ResolverEndpoint.Status != resolvertypes.ResolverEndpointStatusOperational {
		return pkgerrors.Errorf("awsroute53resolverendpoint %q status %q, want OPERATIONAL", id, out.ResolverEndpoint.Status)
	}
	return nil
}

func (v *resolverEndpointVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := resolverClientForRegion(cfg, region).GetResolverEndpoint(ctx, &route53resolver.GetResolverEndpointInput{
		ResolverEndpointId: aws.String(id),
	})
	if err == nil {
		return pkgerrors.Errorf("awsroute53resolverendpoint %q still exists after destroy", id)
	}
	if isResolverNotFound(err) {
		return nil
	}
	return pkgerrors.Wrapf(err, "awsroute53resolverendpoint verify-absent failed for %q", id)
}

func (v *resolverEndpointVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	endpointId := stringOutput(outputs, "endpoint_id")
	if endpointId == "" {
		return pkgerrors.New("awsroute53resolverendpoint verify-exists: no endpoint_id in outputs")
	}
	if err := v.VerifyExists(ctx, cfg, endpointId, region); err != nil {
		return err
	}
	client := resolverClientForRegion(cfg, region)
	for ruleName, ruleId := range stringMapOutput(outputs["rule_ids"]) {
		if _, err := client.GetResolverRule(ctx, &route53resolver.GetResolverRuleInput{
			ResolverRuleId: aws.String(ruleId),
		}); err != nil {
			return pkgerrors.Wrapf(err, "awsroute53resolverendpoint rule %q (%s) verify-exists failed", ruleName, ruleId)
		}
	}
	for pairKey, associationId := range stringMapOutput(outputs["rule_association_ids"]) {
		out, err := client.GetResolverRuleAssociation(ctx, &route53resolver.GetResolverRuleAssociationInput{
			ResolverRuleAssociationId: aws.String(associationId),
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awsroute53resolverendpoint association %q (%s) verify-exists failed", pairKey, associationId)
		}
		if out.ResolverRuleAssociation.Status != resolvertypes.ResolverRuleAssociationStatusComplete {
			return pkgerrors.Errorf("awsroute53resolverendpoint association %q status %q, want COMPLETE", pairKey, out.ResolverRuleAssociation.Status)
		}
	}
	return nil
}

func (v *resolverEndpointVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	client := resolverClientForRegion(cfg, region)
	for pairKey, associationId := range stringMapOutput(outputs["rule_association_ids"]) {
		_, err := client.GetResolverRuleAssociation(ctx, &route53resolver.GetResolverRuleAssociationInput{
			ResolverRuleAssociationId: aws.String(associationId),
		})
		if err == nil {
			return pkgerrors.Errorf("awsroute53resolverendpoint association %q still exists after destroy", pairKey)
		}
		if !isResolverNotFound(err) {
			return pkgerrors.Wrapf(err, "awsroute53resolverendpoint association %q verify-absent failed", pairKey)
		}
	}
	for ruleName, ruleId := range stringMapOutput(outputs["rule_ids"]) {
		_, err := client.GetResolverRule(ctx, &route53resolver.GetResolverRuleInput{
			ResolverRuleId: aws.String(ruleId),
		})
		if err == nil {
			return pkgerrors.Errorf("awsroute53resolverendpoint rule %q still exists after destroy", ruleName)
		}
		if !isResolverNotFound(err) {
			return pkgerrors.Wrapf(err, "awsroute53resolverendpoint rule %q verify-absent failed", ruleName)
		}
	}
	if endpointId := stringOutput(outputs, "endpoint_id"); endpointId != "" {
		return v.VerifyAbsent(ctx, cfg, endpointId, region)
	}
	return nil
}

// --- AwsRoute53ResolverFirewall -------------------------------------------

// resolverFirewallVerifier asserts the rule group, every owned domain
// list, every rule (counted against the rule_match_ids map), and
// every VPC association (COMPLETE - the association carries an
// asynchronous status) exist.
type resolverFirewallVerifier struct{}

func (*resolverFirewallVerifier) IDOutputKey() string { return "rule_group_id" }

func (v *resolverFirewallVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := resolverClientForRegion(cfg, region).GetFirewallRuleGroup(ctx, &route53resolver.GetFirewallRuleGroupInput{
		FirewallRuleGroupId: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsroute53resolverfirewall verify-exists failed for %q", id)
	}
	return nil
}

func (v *resolverFirewallVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := resolverClientForRegion(cfg, region).GetFirewallRuleGroup(ctx, &route53resolver.GetFirewallRuleGroupInput{
		FirewallRuleGroupId: aws.String(id),
	})
	if err == nil {
		return pkgerrors.Errorf("awsroute53resolverfirewall %q still exists after destroy", id)
	}
	if isResolverNotFound(err) {
		return nil
	}
	return pkgerrors.Wrapf(err, "awsroute53resolverfirewall verify-absent failed for %q", id)
}

func (v *resolverFirewallVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	ruleGroupId := stringOutput(outputs, "rule_group_id")
	if ruleGroupId == "" {
		return pkgerrors.New("awsroute53resolverfirewall verify-exists: no rule_group_id in outputs")
	}
	if err := v.VerifyExists(ctx, cfg, ruleGroupId, region); err != nil {
		return err
	}
	client := resolverClientForRegion(cfg, region)
	for listName, listId := range stringMapOutput(outputs["domain_list_ids"]) {
		if _, err := client.GetFirewallDomainList(ctx, &route53resolver.GetFirewallDomainListInput{
			FirewallDomainListId: aws.String(listId),
		}); err != nil {
			return pkgerrors.Wrapf(err, "awsroute53resolverfirewall domain list %q (%s) verify-exists failed", listName, listId)
		}
	}
	// Rules have no Get API keyed by our exports alone; the
	// rule_match_ids map's size is the declared rule count - assert
	// AWS agrees.
	ruleMatchIds := stringMapOutput(outputs["rule_match_ids"])
	if len(ruleMatchIds) > 0 {
		rules, err := client.ListFirewallRules(ctx, &route53resolver.ListFirewallRulesInput{
			FirewallRuleGroupId: aws.String(ruleGroupId),
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awsroute53resolverfirewall list-rules failed for %q", ruleGroupId)
		}
		if len(rules.FirewallRules) != len(ruleMatchIds) {
			return pkgerrors.Errorf("awsroute53resolverfirewall %q carries %d rules, want %d", ruleGroupId, len(rules.FirewallRules), len(ruleMatchIds))
		}
	}
	for associationName, associationId := range stringMapOutput(outputs["association_ids"]) {
		out, err := client.GetFirewallRuleGroupAssociation(ctx, &route53resolver.GetFirewallRuleGroupAssociationInput{
			FirewallRuleGroupAssociationId: aws.String(associationId),
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awsroute53resolverfirewall association %q (%s) verify-exists failed", associationName, associationId)
		}
		if out.FirewallRuleGroupAssociation.Status != resolvertypes.FirewallRuleGroupAssociationStatusComplete {
			return pkgerrors.Errorf("awsroute53resolverfirewall association %q status %q, want COMPLETE", associationName, out.FirewallRuleGroupAssociation.Status)
		}
	}
	return nil
}

func (v *resolverFirewallVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	client := resolverClientForRegion(cfg, region)
	for associationName, associationId := range stringMapOutput(outputs["association_ids"]) {
		_, err := client.GetFirewallRuleGroupAssociation(ctx, &route53resolver.GetFirewallRuleGroupAssociationInput{
			FirewallRuleGroupAssociationId: aws.String(associationId),
		})
		if err == nil {
			return pkgerrors.Errorf("awsroute53resolverfirewall association %q still exists after destroy", associationName)
		}
		if !isResolverNotFound(err) {
			return pkgerrors.Wrapf(err, "awsroute53resolverfirewall association %q verify-absent failed", associationName)
		}
	}
	for listName, listId := range stringMapOutput(outputs["domain_list_ids"]) {
		_, err := client.GetFirewallDomainList(ctx, &route53resolver.GetFirewallDomainListInput{
			FirewallDomainListId: aws.String(listId),
		})
		if err == nil {
			return pkgerrors.Errorf("awsroute53resolverfirewall domain list %q still exists after destroy", listName)
		}
		if !isResolverNotFound(err) {
			return pkgerrors.Wrapf(err, "awsroute53resolverfirewall domain list %q verify-absent failed", listName)
		}
	}
	if ruleGroupId := stringOutput(outputs, "rule_group_id"); ruleGroupId != "" {
		return v.VerifyAbsent(ctx, cfg, ruleGroupId, region)
	}
	return nil
}

// --- AwsRoute53ResolverQueryLog -------------------------------------------

// resolverQueryLogVerifier asserts the configuration is CREATED and
// every VPC association is ACTIVE - the association FAILS
// asynchronously when the resolver cannot write to the destination,
// so mere existence is never enough.
type resolverQueryLogVerifier struct{}

func (*resolverQueryLogVerifier) IDOutputKey() string { return "query_log_config_id" }

func (v *resolverQueryLogVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := resolverClientForRegion(cfg, region).GetResolverQueryLogConfig(ctx, &route53resolver.GetResolverQueryLogConfigInput{
		ResolverQueryLogConfigId: aws.String(id),
	})
	if err != nil {
		return pkgerrors.Wrapf(err, "awsroute53resolverquerylog verify-exists failed for %q", id)
	}
	if out.ResolverQueryLogConfig.Status != resolvertypes.ResolverQueryLogConfigStatusCreated {
		return pkgerrors.Errorf("awsroute53resolverquerylog %q status %q, want CREATED", id, out.ResolverQueryLogConfig.Status)
	}
	return nil
}

func (v *resolverQueryLogVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := resolverClientForRegion(cfg, region).GetResolverQueryLogConfig(ctx, &route53resolver.GetResolverQueryLogConfigInput{
		ResolverQueryLogConfigId: aws.String(id),
	})
	if err == nil {
		return pkgerrors.Errorf("awsroute53resolverquerylog %q still exists after destroy", id)
	}
	if isResolverNotFound(err) {
		return nil
	}
	return pkgerrors.Wrapf(err, "awsroute53resolverquerylog verify-absent failed for %q", id)
}

func (v *resolverQueryLogVerifier) VerifyExistsFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	configId := stringOutput(outputs, "query_log_config_id")
	if configId == "" {
		return pkgerrors.New("awsroute53resolverquerylog verify-exists: no query_log_config_id in outputs")
	}
	if err := v.VerifyExists(ctx, cfg, configId, region); err != nil {
		return err
	}
	client := resolverClientForRegion(cfg, region)
	for vpcId, associationId := range stringMapOutput(outputs["association_ids"]) {
		out, err := client.GetResolverQueryLogConfigAssociation(ctx, &route53resolver.GetResolverQueryLogConfigAssociationInput{
			ResolverQueryLogConfigAssociationId: aws.String(associationId),
		})
		if err != nil {
			return pkgerrors.Wrapf(err, "awsroute53resolverquerylog association for vpc %q verify-exists failed", vpcId)
		}
		if out.ResolverQueryLogConfigAssociation.Status != resolvertypes.ResolverQueryLogConfigAssociationStatusActive {
			return pkgerrors.Errorf("awsroute53resolverquerylog association for vpc %q status %q (error: %v %v), want ACTIVE",
				vpcId, out.ResolverQueryLogConfigAssociation.Status,
				out.ResolverQueryLogConfigAssociation.Error, aws.ToString(out.ResolverQueryLogConfigAssociation.ErrorMessage))
		}
	}
	return nil
}

func (v *resolverQueryLogVerifier) VerifyAbsentFromOutputs(ctx context.Context, cfg aws.Config, outputs map[string]interface{}, region string) error {
	client := resolverClientForRegion(cfg, region)
	for vpcId, associationId := range stringMapOutput(outputs["association_ids"]) {
		_, err := client.GetResolverQueryLogConfigAssociation(ctx, &route53resolver.GetResolverQueryLogConfigAssociationInput{
			ResolverQueryLogConfigAssociationId: aws.String(associationId),
		})
		if err == nil {
			return pkgerrors.Errorf("awsroute53resolverquerylog association for vpc %q still exists after destroy", vpcId)
		}
		if !isResolverNotFound(err) {
			return pkgerrors.Wrapf(err, "awsroute53resolverquerylog association for vpc %q verify-absent failed", vpcId)
		}
	}
	if configId := stringOutput(outputs, "query_log_config_id"); configId != "" {
		return v.VerifyAbsent(ctx, cfg, configId, region)
	}
	return nil
}
