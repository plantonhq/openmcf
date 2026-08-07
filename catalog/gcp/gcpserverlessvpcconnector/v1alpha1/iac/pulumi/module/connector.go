package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/vpcaccess"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// connector provisions the Serverless VPC Access connector — the managed
// instance fleet that bridges serverless egress into a VPC network.
func connector(
	ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider,
) (*vpcaccess.Connector, error) {
	spec := locals.GcpServerlessVpcConnector.Spec

	// Enable the Serverless VPC Access API before creating the connector so
	// a fresh project works first try. disable_on_destroy=false: turning an
	// API off on teardown is a project-wide blast radius no single resource
	// should own.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("vpcaccess.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"connector-vpcaccess.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to enable vpcaccess.googleapis.com api")
	}

	args := &vpcaccess.ConnectorArgs{
		Name:   pulumi.String(locals.ConnectorName),
		Region: pulumi.String(spec.Region),
	}

	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Placement is an exactly-one contract enforced pre-deploy by the spec's
	// CEL rules: either network + ip_cidr_range (the connector carves a new
	// /28 out of the network) or an existing dedicated /28 subnet (the
	// Shared-VPC-capable mode).
	if spec.Network.GetValue() != "" {
		args.Network = pulumi.String(spec.Network.GetValue())
	}
	if spec.IpCidrRange != "" {
		args.IpCidrRange = pulumi.String(spec.IpCidrRange)
	}
	if spec.Subnet != nil {
		subnetArgs := &vpcaccess.ConnectorSubnetArgs{
			Name: pulumi.String(spec.Subnet.Name.GetValue()),
		}
		if spec.Subnet.ProjectId != "" {
			subnetArgs.ProjectId = pulumi.String(spec.Subnet.ProjectId)
		}
		args.Subnet = subnetArgs
	}

	if spec.MachineType != "" {
		args.MachineType = pulumi.String(spec.MachineType)
	}

	// Scaling: only the instance-based contract is modeled. The legacy
	// min/max_throughput fields are deliberately not set — the provider
	// discourages them in favor of instances, they conflict with the
	// instance fields, and they force replacement on change.
	//
	// Shrink asymmetry (worth knowing before an update): the provider
	// applies INCREASES to min/max_instances in place but forces the
	// connector to be REPLACED when either value is DECREASED — a brief
	// egress outage for every workload using the connector.
	if spec.MinInstances != nil {
		args.MinInstances = pulumi.Int(int(spec.GetMinInstances()))
	}
	if spec.MaxInstances != nil {
		args.MaxInstances = pulumi.Int(int(spec.GetMaxInstances()))
	}

	createdConnector, err := vpcaccess.NewConnector(ctx,
		locals.GcpServerlessVpcConnector.Metadata.Name,
		args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create vpc access connector")
	}

	return createdConnector, nil
}
