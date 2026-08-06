package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/memorydb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// acl provisions the aws_memorydb_acl. Membership add/remove applies in
// place -- AWS diffs the user set on update, so granting an application
// access never disturbs the others.
func acl(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.AwsMemorydbAcl.Spec

	// Membership refs arrive pre-resolved to plain user names (the platform
	// flattens valueFrom references before the module runs), so GetValue()
	// is the whole extraction -- literals and resolved refs look identical.
	// Empty values are skipped so an unresolved ref never becomes a phantom
	// "" member AWS would reject.
	userNames := make([]string, 0, len(spec.UserNames))
	for _, ref := range spec.UserNames {
		if ref.GetValue() != "" {
			userNames = append(userNames, ref.GetValue())
		}
	}

	args := &memorydb.AclArgs{
		// The AWS ACL name is create-time immutable and doubles as the
		// Pulumi resource name -- metadata.name on both engines.
		Name: pulumi.String(locals.AclName),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}
	// An empty ACL is valid (MemoryDB has no mandatory member); send the
	// set only when it has entries so the create call mirrors Terraform's.
	if len(userNames) > 0 {
		args.UserNames = pulumi.ToStringArray(userNames)
	}

	created, err := memorydb.NewAcl(ctx, locals.AclName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create ACL")
	}

	ctx.Export(OpAclName, created.Name)
	ctx.Export(OpAclArn, created.Arn)
	ctx.Export(OpMinimumEngineVersion, created.MinimumEngineVersion)

	return nil
}
