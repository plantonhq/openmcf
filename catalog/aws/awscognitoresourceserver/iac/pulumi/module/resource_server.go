package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cognito"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func resourceServer(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &cognito.ResourceServerArgs{
		// The identifier is the resource server's identity within the pool
		// (ForceNew) and the prefix of every scope it mints.
		Identifier: pulumi.String(spec.Identifier),
		Name:       pulumi.String(spec.Name),
		UserPoolId: pulumi.String(spec.UserPoolId.GetValue()),
	}

	if len(spec.Scopes) > 0 {
		var scopes cognito.ResourceServerScopeArray
		for _, s := range spec.Scopes {
			scopes = append(scopes, &cognito.ResourceServerScopeArgs{
				ScopeName:        pulumi.String(s.ScopeName),
				ScopeDescription: pulumi.String(s.ScopeDescription),
			})
		}
		args.Scopes = scopes
	}

	created, err := cognito.NewResourceServer(ctx, locals.Target.Metadata.Name, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create Cognito resource server")
	}

	// ---------------------------------------------------------------------------
	// Exports. Scope identifiers are computed from the spec ("{identifier}/
	// {scope_name}") rather than read back from the provider, so the export
	// order matches the spec order deterministically on both engines.
	// ---------------------------------------------------------------------------

	ctx.Export(OpResourceServerIdentifier, created.Identifier)

	scopeIdentifiers := make([]string, 0, len(spec.Scopes))
	for _, s := range spec.Scopes {
		scopeIdentifiers = append(scopeIdentifiers, fmt.Sprintf("%s/%s", spec.Identifier, s.ScopeName))
	}
	ctx.Export(OpScopeIdentifiers, pulumi.ToStringArray(scopeIdentifiers))

	// Echo the resolved pool id: AWS keys resource servers by the
	// (pool id, identifier) pair.
	ctx.Export(OpUserPoolId, created.UserPoolId)

	return nil
}
