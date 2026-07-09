package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/memorydb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cluster creates the MemoryDB cluster and exports outputs.
func cluster(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdSubnetGroup *memorydb.SubnetGroup,
	createdParamGroup *memorydb.ParameterGroup,
) error {
	spec := locals.Spec

	args := &memorydb.ClusterArgs{
		// The AWS cluster name is create-time immutable and doubles as the
		// Pulumi resource name -- metadata.name on both engines (never
		// provider auto-naming, which would suffix a random token and
		// diverge from Terraform).
		Name: pulumi.String(locals.ClusterName),
		// The ACL ref arrives pre-resolved (a literal like "open-access" or
		// a flattened reference to an AwsMemorydbAcl's exported name).
		AclName:  pulumi.String(spec.AclName.GetValue()),
		NodeType: pulumi.String(spec.NodeType),
		// Description is ALWAYS sent explicitly -- when left to their own
		// defaults the two providers inject differing "Managed by ..."
		// strings and the engines' state permanently diverges.
		Description: pulumi.String(spec.Description),
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}

	// Engine (redis/valkey). Version left empty lets AWS pick the engine's
	// default; AWS supports redis -> valkey switches in place.
	if spec.Engine != "" {
		args.Engine = pulumi.String(spec.Engine)
	}
	if spec.EngineVersion != "" {
		args.EngineVersion = pulumi.String(spec.EngineVersion)
	}

	// Port (ForceNew).
	if spec.Port != nil {
		args.Port = pulumi.Int(int(*spec.Port))
	}

	// Topology -- both dials scale in place.
	if spec.NumShards != nil {
		args.NumShards = pulumi.Int(int(*spec.NumShards))
	}
	if spec.NumReplicasPerShard != nil {
		args.NumReplicasPerShard = pulumi.Int(int(*spec.NumReplicasPerShard))
	}

	// Networking: the module-managed subnet group wins when the folded arm
	// was used; otherwise the bring-your-own name (CEL guarantees the arms
	// are exclusive). Neither set falls back to AWS's account "default"
	// subnet group. The output is exported in every arm (empty when the
	// default group applies) so both engines emit the same output set.
	switch {
	case createdSubnetGroup != nil:
		args.SubnetGroupName = createdSubnetGroup.Name
		ctx.Export(OpSubnetGroupName, createdSubnetGroup.Name)
	case spec.SubnetGroupName != "":
		args.SubnetGroupName = pulumi.String(spec.SubnetGroupName)
		ctx.Export(OpSubnetGroupName, pulumi.String(spec.SubnetGroupName))
	default:
		ctx.Export(OpSubnetGroupName, pulumi.String(""))
	}

	var sgIds pulumi.StringArray
	for _, sg := range spec.SecurityGroupIds {
		if sg.GetValue() != "" {
			sgIds = append(sgIds, pulumi.String(sg.GetValue()))
		}
	}
	if len(sgIds) > 0 {
		args.SecurityGroupIds = sgIds
	}

	// Dual-stack networking (both ForceNew on the network type; discovery
	// updates in place).
	if spec.NetworkType != "" {
		args.NetworkType = pulumi.String(spec.NetworkType)
	}
	if spec.IpDiscovery != "" {
		args.IpDiscovery = pulumi.String(spec.IpDiscovery)
	}

	// Encryption. TLS default true; at-rest encryption is always on, the
	// KMS ref merely substitutes a customer-managed key (ForceNew).
	if spec.TlsEnabled != nil {
		args.TlsEnabled = pulumi.Bool(*spec.TlsEnabled)
	}
	if spec.KmsKeyArn.GetValue() != "" {
		args.KmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
	}

	// Maintenance and snapshots. The retention limit is ALWAYS sent —
	// 0 explicitly disables automatic snapshots (the spec's documented
	// contract), and Terraform sends the same explicit 0, so omitting it
	// here would let AWS pick a default and diverge the engines.
	if spec.MaintenanceWindow != "" {
		args.MaintenanceWindow = pulumi.String(spec.MaintenanceWindow)
	}
	args.SnapshotRetentionLimit = pulumi.Int(int(spec.SnapshotRetentionLimit))
	if spec.SnapshotWindow != "" {
		args.SnapshotWindow = pulumi.String(spec.SnapshotWindow)
	}
	if spec.FinalSnapshotName != "" {
		args.FinalSnapshotName = pulumi.String(spec.FinalSnapshotName)
	}

	// Create-time restore sources (mutually exclusive by CEL).
	if len(spec.SnapshotArns) > 0 {
		args.SnapshotArns = pulumi.ToStringArray(spec.SnapshotArns)
	}
	if spec.SnapshotName != "" {
		args.SnapshotName = pulumi.String(spec.SnapshotName)
	}

	// Parameter group: module-managed wins; else bring-your-own name. The
	// output is exported in every arm (empty when the family default
	// applies) so both engines emit the same output set.
	switch {
	case createdParamGroup != nil:
		args.ParameterGroupName = createdParamGroup.Name
		ctx.Export(OpParameterGroupName, createdParamGroup.Name)
	case spec.ParameterGroupName != "":
		args.ParameterGroupName = pulumi.String(spec.ParameterGroupName)
		ctx.Export(OpParameterGroupName, pulumi.String(spec.ParameterGroupName))
	default:
		ctx.Export(OpParameterGroupName, pulumi.String(""))
	}

	// Multi-region active-active membership (ForceNew): the multi-region
	// cluster is created outside this resource; this regional cluster joins
	// it by name.
	if spec.MultiRegionClusterName != "" {
		args.MultiRegionClusterName = pulumi.String(spec.MultiRegionClusterName)
	}

	// Notifications.
	if spec.SnsTopicArn.GetValue() != "" {
		args.SnsTopicArn = pulumi.String(spec.SnsTopicArn.GetValue())
	}

	// Advanced (both ForceNew).
	if spec.AutoMinorVersionUpgrade != nil {
		args.AutoMinorVersionUpgrade = pulumi.Bool(*spec.AutoMinorVersionUpgrade)
	}
	if spec.DataTiering {
		args.DataTiering = pulumi.Bool(true)
	}

	c, err := memorydb.NewCluster(ctx, locals.ClusterName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create memorydb cluster")
	}

	ctx.Export(OpClusterArn, c.Arn)
	ctx.Export(OpClusterName, c.Name)
	ctx.Export(OpEnginePatchVersion, c.EnginePatchVersion)

	// The single cluster endpoint handles slot discovery and routing;
	// Index() on an empty list yields the element type's zero value, so
	// these exports are ApplyT-free by design.
	ctx.Export(OpClusterEndpointAddress, c.ClusterEndpoints.Index(pulumi.Int(0)).Address())
	ctx.Export(OpClusterEndpointPort, c.ClusterEndpoints.Index(pulumi.Int(0)).Port())

	return nil
}
