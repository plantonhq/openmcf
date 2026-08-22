package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// databaseDb provisions the logical database and exports its outputs. Both
// arguments are create-only upstream: any change replaces the logical
// database and drops its data.
func databaseDb(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.DatabaseDb, error) {
	spec := locals.DigitalOceanDatabaseDb.Spec

	createdDb, err := digitalocean.NewDatabaseDb(
		ctx,
		"database",
		&digitalocean.DatabaseDbArgs{
			// References are resolved to the literal cluster UUID before
			// the module runs.
			ClusterId: pulumi.String(spec.Cluster.GetValue()),
			Name:      pulumi.String(spec.DatabaseName),
		},
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean logical database")
	}

	ctx.Export(OpClusterId, createdDb.ClusterId)
	ctx.Export(OpDatabaseName, createdDb.Name)

	return createdDb, nil
}
