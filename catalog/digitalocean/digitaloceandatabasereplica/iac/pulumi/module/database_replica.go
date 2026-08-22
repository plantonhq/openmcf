package module

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// databaseReplica provisions the read-only replica and exports its
// outputs.
//
// Update semantics mirror the provider: only size and storage_size_mib
// change in place (a resize); every other argument change REPLACES the
// replica -- including tags, which are create-only upstream. region and
// size are required by the spec (explicit values kill the upstream
// omitted-value drift class that would otherwise schedule replacements).
func databaseReplica(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.DatabaseReplica, error) {
	spec := locals.DigitalOceanDatabaseReplica.Spec

	// User tags plus the standard Planton labels rendered as "key:value"
	// tags -- the exact set the Terraform module applies. CREATE-ONLY
	// upstream: a change to the final set replaces the replica.
	tagSet := map[string]bool{}
	var tagInputs pulumi.StringArray
	for _, t := range spec.Tags {
		if !tagSet[t] {
			tagSet[t] = true
			tagInputs = append(tagInputs, pulumi.String(t))
		}
	}
	for k, v := range locals.DigitalOceanLabels {
		t := k + ":" + v
		if !tagSet[t] {
			tagSet[t] = true
			tagInputs = append(tagInputs, pulumi.String(t))
		}
	}

	replicaArgs := &digitalocean.DatabaseReplicaArgs{
		// References are resolved to the literal cluster UUID before the
		// module runs. Enum value names are exactly the DigitalOcean
		// region slugs.
		ClusterId: pulumi.String(spec.Cluster.GetValue()),
		Name:      pulumi.String(spec.ReplicaName),
		Region:    pulumi.String(spec.Region.String()),
		Size:      pulumi.String(spec.Size),
		Tags:      tagInputs,
	}

	// Optional VPC placement for the replica's region (create-only;
	// Optional+Computed upstream, so omission is drift-safe).
	if spec.Vpc != nil && spec.Vpc.GetValue() != "" {
		replicaArgs.PrivateNetworkUuid = pulumi.StringPtr(spec.Vpc.GetValue())
	}

	// The provider's storage_size_mib is a string holding a bare MiB
	// count; the spec carries the number. Must stay >= the primary's
	// storage.
	if spec.StorageSizeMib > 0 {
		replicaArgs.StorageSizeMib = pulumi.StringPtr(strconv.FormatUint(spec.StorageSizeMib, 10))
	}

	createdReplica, err := digitalocean.NewDatabaseReplica(
		ctx,
		"replica",
		replicaArgs,
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean database replica")
	}

	// replica_id is the API UUID (the SDK's Uuid attribute) -- the
	// resource's own ID() is a legacy composite string, not the UUID.
	ctx.Export(OpReplicaId, createdReplica.Uuid)
	ctx.Export(OpClusterId, createdReplica.ClusterId)
	ctx.Export(OpReplicaName, createdReplica.Name)
	ctx.Export(OpHost, createdReplica.Host)
	ctx.Export(OpPrivateHost, createdReplica.PrivateHost)
	ctx.Export(OpPort, createdReplica.Port)
	ctx.Export(OpDatabase, createdReplica.Database)
	ctx.Export(OpUser, createdReplica.User)
	ctx.Export(OpPassword, createdReplica.Password)
	ctx.Export(OpUri, createdReplica.Uri)
	ctx.Export(OpPrivateUri, createdReplica.PrivateUri)

	return createdReplica, nil
}
