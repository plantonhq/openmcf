package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type EipResult struct {
	AllocationId  pulumi.StringOutput
	PublicIp      pulumi.StringOutput
	Arn           pulumi.StringOutput
	PublicDns     pulumi.StringOutput
	AssociationId pulumi.StringOutput
	PtrRecord     pulumi.StringOutput
}

func eip(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*EipResult, error) {
	spec := locals.AwsElasticIp.Spec

	args := &ec2.EipArgs{
		// Every managed EIP is a VPC EIP: AWS retired EC2-Classic ("standard"
		// domain) addresses and the provider itself refuses to read or delete
		// them, so the legacy domain is deliberately not configurable.
		Domain: pulumi.String("vpc"),
		Tags:   pulumi.ToStringMap(locals.AwsTags),
	}

	// BYOIP: allocate from a specific IPv4 address pool.
	if spec.PublicIpv4Pool != "" {
		args.PublicIpv4Pool = pulumi.StringPtr(spec.PublicIpv4Pool)
	}

	// BYOIP / IPAM: request a specific IP address the pool holds.
	if spec.Address != "" {
		args.Address = pulumi.StringPtr(spec.Address)
	}

	// IPAM: allocate from an Amazon VPC IP Address Manager public pool.
	if spec.IpamPoolId != "" {
		args.IpamPoolId = pulumi.StringPtr(spec.IpamPoolId)
	}

	// Location scope for Local Zones and Wavelength zones.
	if spec.NetworkBorderGroup != "" {
		args.NetworkBorderGroup = pulumi.StringPtr(spec.NetworkBorderGroup)
	}

	// Association: attach the address to an instance XOR a network interface
	// (spec CEL enforces at-most-one; AWS associates with exactly one target).
	if spec.Instance.GetValue() != "" {
		args.Instance = pulumi.StringPtr(spec.Instance.GetValue())
	}
	if spec.NetworkInterface.GetValue() != "" {
		args.NetworkInterface = pulumi.StringPtr(spec.NetworkInterface.GetValue())
	}
	if spec.AssociateWithPrivateIp != "" {
		args.AssociateWithPrivateIp = pulumi.StringPtr(spec.AssociateWithPrivateIp)
	}

	createdEip, err := ec2.NewEip(ctx, locals.AwsElasticIp.Metadata.Name, args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create elastic ip")
	}

	result := &EipResult{
		AllocationId:  createdEip.AllocationId,
		PublicIp:      createdEip.PublicIp,
		Arn:           createdEip.Arn,
		PublicDns:     createdEip.PublicDns,
		AssociationId: createdEip.AssociationId,
		// Empty unless the reverse-DNS resource below overrides it.
		PtrRecord: pulumi.String("").ToStringOutput(),
	}

	// Reverse DNS (PTR) record for the address. AWS validates SERVER-SIDE that
	// a forward A record for the domain already resolves to this EIP before
	// granting the PTR — a fresh EIP therefore typically sets this on a
	// follow-up apply, after DNS points at the address.
	if spec.ReverseDnsDomainName != "" {
		createdDomainName, err := ec2.NewEipDomainName(ctx,
			locals.AwsElasticIp.Metadata.Name,
			&ec2.EipDomainNameArgs{
				AllocationId: createdEip.AllocationId,
				DomainName:   pulumi.String(spec.ReverseDnsDomainName),
			}, pulumi.Provider(provider), pulumi.Parent(createdEip))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create elastic ip reverse dns domain name")
		}
		result.PtrRecord = createdDomainName.PtrRecord
	}

	return result, nil
}
