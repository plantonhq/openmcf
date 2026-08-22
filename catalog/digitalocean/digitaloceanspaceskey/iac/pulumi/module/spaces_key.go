package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// spacesKey provisions the access-key pair and exports its outputs. Name
// and grants update in place (the grant list is replaced wholesale); the
// key material itself never changes -- rotation is destroy and recreate.
func spacesKey(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.SpacesKey, error) {
	spec := locals.DigitalOceanSpacesKey.Spec

	// The provider's grant grammar: read/readwrite grants name their
	// bucket; a fullaccess grant carries an EMPTY bucket string (spec
	// validation guarantees the pairing). Bucket references resolve to
	// literal bucket names before the module runs.
	var grants digitalocean.SpacesKeyGrantArray
	for _, grant := range spec.Grants {
		grants = append(grants, &digitalocean.SpacesKeyGrantArgs{
			Bucket:     pulumi.String(grant.GetBucket().GetValue()),
			Permission: pulumi.String(grant.Permission),
		})
	}

	createdKey, err := digitalocean.NewSpacesKey(
		ctx,
		"spaces-key",
		&digitalocean.SpacesKeyArgs{
			Name:   pulumi.String(spec.KeyName),
			Grants: grants,
		},
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean spaces key")
	}

	ctx.Export(OpAccessKey, createdKey.AccessKey)
	// The SDK does not secret-flag the secret key; wrap it so it is
	// encrypted in state and masked in every output surface.
	ctx.Export(OpSecretKey, pulumi.ToSecret(createdKey.SecretKey))

	return createdKey, nil
}
