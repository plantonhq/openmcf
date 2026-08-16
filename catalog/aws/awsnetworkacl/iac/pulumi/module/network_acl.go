package module

import (
	"github.com/pkg/errors"
	awsnetworkaclv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsnetworkacl/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// networkAcl creates the network ACL with its rules and subnet
// associations in-line and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - VpcId replaces the ACL on change; everything else updates in
//     place;
//   - in-line rules and SubnetIds are the single declarative owner
//     (the standalone rule/association resources are identical
//     payloads and fight this form);
//   - AWS stores protocols as NUMBERS and the provider normalizes
//     names in its rule hash, so "tcp" never causes a perpetual diff;
//     AWS's catch-all rules (32767/32768) are invisible and
//     unmanageable;
//   - a subnet listed here is atomically REPLACED onto this ACL;
//     removal hands it back to the VPC's default NACL; destroy tears
//     down all associations first.
func networkAcl(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &ec2.NetworkAclArgs{
		VpcId: pulumi.String(spec.VpcId.GetValue()),
		Tags:  pulumi.ToStringMap(locals.AwsTags),
	}

	if len(spec.Ingress) > 0 {
		ingress := ec2.NetworkAclIngressArray{}
		for _, rule := range spec.Ingress {
			ingress = append(ingress, buildIngressRule(rule))
		}
		args.Ingress = ingress
	}

	if len(spec.Egress) > 0 {
		egress := ec2.NetworkAclEgressArray{}
		for _, rule := range spec.Egress {
			egress = append(egress, buildEgressRule(rule))
		}
		args.Egress = egress
	}

	if len(spec.SubnetIds) > 0 {
		subnetIds := pulumi.StringArray{}
		for _, subnetId := range spec.SubnetIds {
			subnetIds = append(subnetIds, pulumi.String(subnetId.GetValue()))
		}
		args.SubnetIds = subnetIds
	}

	createdAcl, err := ec2.NewNetworkAcl(ctx, "network_acl", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create network acl")
	}

	ctx.Export(OpNetworkAclId, createdAcl.ID())
	ctx.Export(OpNetworkAclArn, createdAcl.Arn)
	ctx.Export(OpOwnerId, createdAcl.OwnerId)
	return nil
}

func buildIngressRule(rule *awsnetworkaclv1alpha1.AwsNetworkAclRule) *ec2.NetworkAclIngressArgs {
	args := &ec2.NetworkAclIngressArgs{
		RuleNo:   pulumi.Int(int(rule.RuleNo)),
		Action:   pulumi.String(rule.Action),
		Protocol: pulumi.String(rule.Protocol),
		FromPort: pulumi.Int(int(rule.FromPort)),
		ToPort:   pulumi.Int(int(rule.ToPort)),
	}
	if rule.CidrBlock != "" {
		args.CidrBlock = pulumi.String(rule.CidrBlock)
	}
	if rule.Ipv6CidrBlock != "" {
		args.Ipv6CidrBlock = pulumi.String(rule.Ipv6CidrBlock)
	}
	if rule.IcmpType != nil {
		args.IcmpType = pulumi.Int(int(*rule.IcmpType))
	}
	if rule.IcmpCode != nil {
		args.IcmpCode = pulumi.Int(int(*rule.IcmpCode))
	}
	return args
}

func buildEgressRule(rule *awsnetworkaclv1alpha1.AwsNetworkAclRule) *ec2.NetworkAclEgressArgs {
	args := &ec2.NetworkAclEgressArgs{
		RuleNo:   pulumi.Int(int(rule.RuleNo)),
		Action:   pulumi.String(rule.Action),
		Protocol: pulumi.String(rule.Protocol),
		FromPort: pulumi.Int(int(rule.FromPort)),
		ToPort:   pulumi.Int(int(rule.ToPort)),
	}
	if rule.CidrBlock != "" {
		args.CidrBlock = pulumi.String(rule.CidrBlock)
	}
	if rule.Ipv6CidrBlock != "" {
		args.Ipv6CidrBlock = pulumi.String(rule.Ipv6CidrBlock)
	}
	if rule.IcmpType != nil {
		args.IcmpType = pulumi.Int(int(*rule.IcmpType))
	}
	if rule.IcmpCode != nil {
		args.IcmpCode = pulumi.Int(int(*rule.IcmpCode))
	}
	return args
}
