package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// prefixList creates the customer-managed prefix list with its entries
// in-line and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - AddressFamily replaces the list on change - every referencing
//     rule breaks with the old pl- id; everything else updates in
//     place;
//   - the in-line entry set is the single declarative owner (the
//     standalone entry resource is the identical payload and fights
//     this form);
//   - the provider orders MaxEntries increases before entry changes
//     and decreases after, so a resize never transiently strands
//     entries; description-only edits cost two API round trips
//     (remove + re-add) - expected, not drift.
func prefixList(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	entries := ec2.ManagedPrefixListEntryTypeArray{}
	for _, entry := range spec.Entries {
		entryArgs := &ec2.ManagedPrefixListEntryTypeArgs{
			Cidr: pulumi.String(entry.Cidr),
		}
		if entry.Description != "" {
			entryArgs.Description = pulumi.String(entry.Description)
		}
		entries = append(entries, entryArgs)
	}

	args := &ec2.ManagedPrefixListArgs{
		Name:          pulumi.String(locals.Target.Metadata.Name),
		AddressFamily: pulumi.String(spec.AddressFamily),
		MaxEntries:    pulumi.Int(int(spec.MaxEntries)),
		Tags:          pulumi.ToStringMap(locals.AwsTags),
	}
	if len(entries) > 0 {
		args.Entries = entries
	}

	createdPrefixList, err := ec2.NewManagedPrefixList(ctx, "prefix_list", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create managed prefix list")
	}

	ctx.Export(OpPrefixListId, createdPrefixList.ID())
	ctx.Export(OpPrefixListArn, createdPrefixList.Arn)
	ctx.Export(OpOwnerId, createdPrefixList.OwnerId)
	// The version is an int at the provider; exported as a string to match
	// the outputs contract (string-typed observable state).
	ctx.Export(OpVersion, pulumi.Sprintf("%d", createdPrefixList.Version))
	return nil
}
