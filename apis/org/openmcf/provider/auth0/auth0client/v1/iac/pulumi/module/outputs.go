package module

import (
	"github.com/pulumi/pulumi-auth0/sdk/v3/go/auth0"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// exportOutputs exports the stack outputs for the Auth0 client
func exportOutputs(ctx *pulumi.Context, client *auth0.Client, locals *Locals) error {
	// Export client ID (the public OAuth identifier)
	ctx.Export("id", client.ID())
	ctx.Export("client_id", client.ClientId)
	ctx.Export("name", client.Name)
	ctx.Export("application_type", client.AppType)

	// Export signing keys
	ctx.Export("signing_keys", client.SigningKeys)

	// Export custom outputs for easy reference
	ctx.Export("metadata_name", pulumi.String(locals.ClientName))

	return nil
}
