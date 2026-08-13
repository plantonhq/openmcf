package module

import (
	"github.com/pkg/errors"
	awssecuritygroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssecuritygroup/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/vpc"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// securityGroup creates the EC2 security group with its rules INLINE. AWS
// forbids mixing inline rules with standalone rule resources on the same
// group (each apply would fight to own the rule set), so this module must
// never emit standalone rule resources -- the inline arrays below are the
// single owner of the group's rules. This mirrors the Terraform module's
// dynamic ingress/egress blocks field-for-field.
func securityGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.AwsSecurityGroup.Spec

	sg, err := ec2.NewSecurityGroup(ctx, locals.AwsSecurityGroup.Metadata.Name, &ec2.SecurityGroupArgs{
		VpcId: pulumi.String(spec.VpcId.GetValue()),
		// The group name is metadata.name -- create-only in AWS, and the
		// basis both engines share so a manifest deploys identically on
		// either.
		Name:        pulumi.String(locals.AwsSecurityGroup.Metadata.Name),
		Description: pulumi.String(spec.Description),
		// Forcibly revoke this group's rules -- and rules in OTHER groups
		// that reference it -- before delete. Without it, destroying a group
		// a sibling group still references fails with a DependencyViolation.
		RevokeRulesOnDelete: pulumi.Bool(spec.RevokeRulesOnDelete),
		Ingress:             buildIngress(spec.Ingress),
		// With inline egress, an empty list means DENY ALL outbound: the
		// provider revokes the allow-all egress rule AWS adds to every new
		// group, so the manifest is the complete statement of what the group
		// permits.
		Egress: buildEgress(spec.Egress),
		Tags:   pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "unable to create AWS Security Group")
	}

	ctx.Export(OpSecurityGroupId, sg.ID())
	ctx.Export(OpSecurityGroupArn, sg.Arn)
	ctx.Export(OpOwnerId, sg.OwnerId)

	// Share the group into additional VPCs (same account and region), keyed
	// by the resolved VPC id -- mirrors the Terraform module's for_each. AWS
	// ignores rules that reference a security group from a different VPC
	// than the one a packet traverses, so multi-VPC groups should prefer
	// CIDR or prefix-list rules (documented on the spec field).
	associationIds := pulumi.StringMap{}
	for _, v := range spec.AdditionalVpcIds {
		vpcId := v.GetValue()
		if vpcId == "" {
			continue
		}
		_, err := vpc.NewSecurityGroupVpcAssociation(ctx,
			locals.AwsSecurityGroup.Metadata.Name+"-"+vpcId,
			&vpc.SecurityGroupVpcAssociationArgs{
				SecurityGroupId: sg.ID(),
				VpcId:           pulumi.String(vpcId),
			}, pulumi.Provider(provider), pulumi.Parent(sg))
		if err != nil {
			return errors.Wrapf(err, "unable to associate security group with vpc %s", vpcId)
		}
		// Import id form: <group_id>,<vpc_id> (the provider's documented
		// comma-separated import key).
		associationIds[vpcId] = pulumi.Sprintf("%s,%s", sg.ID(), vpcId)
	}
	ctx.Export(OpAdditionalVpcAssociationIds, associationIds)

	return nil
}

// buildIngress converts proto-based SecurityGroupRules into Pulumi's SecurityGroupIngressArgs array.
func buildIngress(rules []*awssecuritygroupv1alpha1.SecurityGroupRule) ec2.SecurityGroupIngressArray {
	var ingress ec2.SecurityGroupIngressArray
	for _, r := range rules {
		ingress = append(ingress, ruleToIngress(r))
	}
	return ingress
}

// buildEgress converts proto-based SecurityGroupRules into Pulumi's SecurityGroupEgressArgs array.
func buildEgress(rules []*awssecuritygroupv1alpha1.SecurityGroupRule) ec2.SecurityGroupEgressArray {
	var egress ec2.SecurityGroupEgressArray
	for _, r := range rules {
		egress = append(egress, ruleToEgress(r))
	}
	return egress
}

// ruleToIngress maps a single spec rule to one inline ingress entry. A rule
// may carry several sources at once (CIDRs, prefix lists, other groups,
// self); AWS expands them into individual permissions server-side.
func ruleToIngress(r *awssecuritygroupv1alpha1.SecurityGroupRule) *ec2.SecurityGroupIngressArgs {
	var sourceSGs pulumi.StringArray
	for _, sg := range r.SourceSecurityGroupIds {
		sourceSGs = append(sourceSGs, pulumi.String(sg.GetValue()))
	}
	return &ec2.SecurityGroupIngressArgs{
		Protocol:       pulumi.String(r.Protocol),
		FromPort:       pulumi.Int(int(r.FromPort)),
		ToPort:         pulumi.Int(int(r.ToPort)),
		CidrBlocks:     pulumi.ToStringArray(r.Ipv4Cidrs),
		Ipv6CidrBlocks: pulumi.ToStringArray(r.Ipv6Cidrs),
		// Managed prefix lists let a rule target a named CIDR set (an AWS
		// service like S3, or a customer-maintained office/partner range) by
		// stable ID instead of hardcoding addresses.
		PrefixListIds:  pulumi.ToStringArray(r.PrefixListIds),
		SecurityGroups: sourceSGs,
		Self:           pulumi.Bool(r.SelfReference),
		Description:    pulumi.String(r.Description),
	}
}

// ruleToEgress maps a single spec rule to one inline egress entry.
func ruleToEgress(r *awssecuritygroupv1alpha1.SecurityGroupRule) *ec2.SecurityGroupEgressArgs {
	var destSGs pulumi.StringArray
	for _, sg := range r.DestinationSecurityGroupIds {
		destSGs = append(destSGs, pulumi.String(sg.GetValue()))
	}
	return &ec2.SecurityGroupEgressArgs{
		Protocol:       pulumi.String(r.Protocol),
		FromPort:       pulumi.Int(int(r.FromPort)),
		ToPort:         pulumi.Int(int(r.ToPort)),
		CidrBlocks:     pulumi.ToStringArray(r.Ipv4Cidrs),
		Ipv6CidrBlocks: pulumi.ToStringArray(r.Ipv6Cidrs),
		PrefixListIds:  pulumi.ToStringArray(r.PrefixListIds),
		SecurityGroups: destSGs,
		Self:           pulumi.Bool(r.SelfReference),
		Description:    pulumi.String(r.Description),
	}
}
