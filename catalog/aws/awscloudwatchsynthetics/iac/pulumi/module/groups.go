package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/synthetics"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// groups creates the owned Synthetics groups keyed by name and exports
// the group maps.
func groups(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (map[string]*synthetics.Group, error) {
	createdGroups := map[string]*synthetics.Group{}

	groupArns := pulumi.StringMap{}
	groupIds := pulumi.StringMap{}

	for _, group := range locals.Spec.Groups {
		createdGroup, err := synthetics.NewGroup(ctx, "group-"+group.Name, &synthetics.GroupArgs{
			Name: pulumi.String(group.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create group %s", group.Name)
		}
		createdGroups[group.Name] = createdGroup
		groupArns[group.Name] = createdGroup.Arn
		groupIds[group.Name] = createdGroup.GroupId
	}

	ctx.Export(OpGroupArns, groupArns)
	ctx.Export(OpGroupIds, groupIds)
	return createdGroups, nil
}
