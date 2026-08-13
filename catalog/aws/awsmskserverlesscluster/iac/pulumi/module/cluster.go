package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/msk"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cluster creates the MSK Serverless cluster. A serverless cluster is a
// single, essentially immutable resource: AWS manages brokers, storage, and
// Kafka version, so the whole declaration is WHERE it lives (one or more VPC
// placements). Everything except tags is create-time (ForceNew) -- changing
// networking or auth replaces the cluster. This mirrors the Terraform module
// field-for-field.
func cluster(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*msk.ServerlessCluster, error) {
	spec := locals.AwsMskServerlessCluster.Spec

	// One entry per VPC placement: AWS provisions client-facing network
	// interfaces in EACH declared VPC, so applications in every listed VPC
	// connect privately without peering or PrivateLink setup.
	vpcConfigs := msk.ServerlessClusterVpcConfigArray{}
	for _, vpcConfig := range spec.VpcConfigs {
		subnetIds := pulumi.StringArray{}
		for _, s := range vpcConfig.SubnetIds {
			subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
		}

		// Optional: AWS attaches the VPC's default security group when
		// omitted. The ingress rule for the SASL/IAM listener port (9098)
		// lives on the referenced first-class security group nodes.
		var securityGroupIds pulumi.StringArray
		for _, sg := range vpcConfig.SecurityGroupIds {
			securityGroupIds = append(securityGroupIds, pulumi.String(sg.GetValue()))
		}

		vpcConfigs = append(vpcConfigs, &msk.ServerlessClusterVpcConfigArgs{
			SubnetIds:        subnetIds,
			SecurityGroupIds: securityGroupIds,
		})
	}

	createdCluster, err := msk.NewServerlessCluster(ctx, locals.ClusterName, &msk.ServerlessClusterArgs{
		ClusterName: pulumi.String(locals.ClusterName),
		VpcConfigs:  vpcConfigs,
		// SASL/IAM is the ONLY client-authentication scheme serverless MSK
		// supports, and it is mandatory -- so it is hardcoded rather than
		// exposed as a spec field that could only ever hold one value.
		// Clients authenticate with AWS IAM credentials on port 9098.
		ClientAuthentication: &msk.ServerlessClusterClientAuthenticationArgs{
			Sasl: &msk.ServerlessClusterClientAuthenticationSaslArgs{
				Iam: &msk.ServerlessClusterClientAuthenticationSaslIamArgs{
					Enabled: pulumi.Bool(true),
				},
			},
		},
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create msk serverless cluster")
	}

	return createdCluster, nil
}
