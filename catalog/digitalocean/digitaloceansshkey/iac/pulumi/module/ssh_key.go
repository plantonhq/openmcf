package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// sshKey registers the SSH public key and exports its outputs.
func sshKey(
	ctx *pulumi.Context,
	locals *Locals,
	digitalOceanProvider *digitalocean.Provider,
) (*digitalocean.SshKey, error) {
	spec := locals.DigitalOceanSshKey.Spec

	// The key material is create-only upstream: DigitalOcean trims only
	// leading/trailing whitespace before comparing, so any in-line change
	// REPLACES the key and rotates the id and fingerprint.
	createdSshKey, err := digitalocean.NewSshKey(
		ctx,
		"ssh_key",
		&digitalocean.SshKeyArgs{
			Name:      pulumi.StringPtr(spec.KeyName),
			PublicKey: pulumi.String(spec.PublicKey),
		},
		pulumi.Provider(digitalOceanProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create digitalocean ssh key")
	}

	ctx.Export(OpSshKeyId, createdSshKey.ID())
	ctx.Export(OpFingerprint, createdSshKey.Fingerprint)

	return createdSshKey, nil
}
