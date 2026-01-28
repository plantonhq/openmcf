package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func natGateway(ctx *pulumi.Context, locals *Locals, createdVpc *ec2.Vpc,
	subnetName string, createdPrivateSubnet *ec2.Subnet) (*ec2.NatGateway, error) {

	//create elastic ip for nat gateway
	createdElasticIp, err := ec2.NewEip(ctx,
		fmt.Sprintf("nat-%s", subnetName),
		&ec2.EipArgs{
			Tags: AddEntryToPulumiStringMap(pulumi.ToStringMap(locals.AwsTags), "Name",
				pulumi.Sprintf("%s-nat", createdPrivateSubnet.ID())),
		}, pulumi.Parent(createdPrivateSubnet))
	if err != nil {
		return nil, errors.Wrap(err, "error creating eip for nat gateway")
	}

	//create nat gateway
	createdNatGateway, err := ec2.NewNatGateway(ctx,
		subnetName,
		&ec2.NatGatewayArgs{
			SubnetId:     createdPrivateSubnet.ID(),
			AllocationId: createdElasticIp.ID(),
			Tags: AddIdValueEntryToPulumiStringMap(pulumi.ToStringMap(locals.AwsTags),
				"Name", createdPrivateSubnet.ID()),
		}, pulumi.Parent(createdPrivateSubnet))
	if err != nil {
		return nil, errors.Wrap(err, "error creating nat gateway")
	}

	// private route table to route traffic through nat gateway
	createdPrivateRouteTable, err := ec2.NewRouteTable(ctx,
		fmt.Sprintf("private-route-table-%s", subnetName),
		&ec2.RouteTableArgs{
			VpcId: createdVpc.ID(),
			Routes: ec2.RouteTableRouteArray{
				&ec2.RouteTableRouteArgs{
					CidrBlock:    pulumi.String("0.0.0.0/0"),
					NatGatewayId: createdNatGateway.ID(),
				},
			},
			Tags: AddEntryToPulumiStringMap(pulumi.ToStringMap(locals.AwsTags), "Name",
				pulumi.Sprintf("%s-private", createdPrivateSubnet.ID())),
		}, pulumi.Parent(createdNatGateway))
	if err != nil {
		return nil, errors.Wrap(err, "error creating private route table")
	}

	// associate private route table with private subnets
	_, err = ec2.NewRouteTableAssociation(ctx,
		fmt.Sprintf("private-route-assoc-%s", subnetName),
		&ec2.RouteTableAssociationArgs{
			RouteTableId: createdPrivateRouteTable.ID(),
			SubnetId:     createdPrivateSubnet.ID(),
		}, pulumi.Parent(createdPrivateRouteTable))
	if err != nil {
		return nil, errors.Wrap(err, "error associating private route table")
	}
	return createdNatGateway, nil
}

func AddIdValueEntryToPulumiStringMap(m pulumi.StringMap, key string, id pulumi.IDOutput) pulumi.StringMap {
	m[key] = id
	return m
}

func AddEntryToPulumiStringMap(m pulumi.StringMap, key string, id pulumi.StringOutput) pulumi.StringMap {
	m[key] = id
	return m
}
