package module

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/valuefrom"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// group provisions the IAM group, its declarative membership, and its
// policy wiring.
//
// Lifecycle facts the render below depends on:
//   - the group's name comes from metadata.name and updates IN PLACE
//     at AWS (UpdateGroup renames; the ARN recomputes but members and
//     policies persist);
//   - membership renders as ONE aws_iam_group_membership-equivalent
//     resource carrying the whole users list - the AUTHORITATIVE form:
//     users added out-of-band are removed on the next apply, and
//     clearing the spec list removes the resource (and with it every
//     membership). The users must already exist at IAM;
//   - each managed-policy attachment is its own resource keyed by the
//     policy ARN itself (sanitized for the logical name), never the
//     list index - reordering managed_policy_arns must be a no-op, not
//     a transient detach/re-attach on a live group;
//   - inline policies live and die with the group;
//   - IAM groups and every satellite here are untaggable at AWS.
func group(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	groupName := locals.Target.Metadata.Name
	spec := locals.Spec

	groupArgs := &iam.GroupArgs{
		Name: pulumi.String(groupName),
	}
	if spec.Path != "" {
		groupArgs.Path = pulumi.StringPtr(spec.Path)
	}

	createdGroup, err := iam.NewGroup(ctx, groupName, groupArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create IAM group")
	}

	// One membership resource owns the whole users list (the
	// group-centric declarative form). valueFrom references were
	// resolved to user names before the module ran.
	if users := valuefrom.ToStringArray(spec.Users); len(users) > 0 {
		_, err := iam.NewGroupMembership(ctx, fmt.Sprintf("%s-membership", groupName), &iam.GroupMembershipArgs{
			// The membership resource's own name is a provider-side
			// label (never reaches AWS); the group name keeps it
			// stable and readable.
			Name:  pulumi.String(fmt.Sprintf("%s-membership", groupName)),
			Group: createdGroup.Name,
			Users: pulumi.ToStringArray(users),
		}, pulumi.Provider(provider), pulumi.Parent(createdGroup))
		if err != nil {
			return errors.Wrap(err, "failed to create group membership")
		}
	}

	// Each managed-policy attachment reconciles individually: adding or
	// removing an entry attaches or detaches just that policy, and
	// attachments made outside this resource are left alone.
	for _, policyArn := range valuefrom.ToStringArray(spec.ManagedPolicyArns) {
		attachName := fmt.Sprintf("%s-attach-%s", groupName, sanitizeForLogicalName(policyArn))
		_, err := iam.NewGroupPolicyAttachment(ctx, attachName, &iam.GroupPolicyAttachmentArgs{
			Group:     createdGroup.Name,
			PolicyArn: pulumi.String(policyArn),
		}, pulumi.Provider(provider), pulumi.Parent(createdGroup))
		if err != nil {
			return errors.Wrapf(err, "failed to attach policy ARN %s", policyArn)
		}
	}

	// Inline policies live and die with the group -- permissions unique
	// to this one group that would be noise as standalone AwsIamPolicy
	// resources.
	for policyName, inlineStruct := range spec.InlinePolicies {
		inlinePolicyString, err := structToJSONString(inlineStruct)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal inline policy for %s", policyName)
		}

		inlineName := fmt.Sprintf("%s-inline-%s", groupName, policyName)
		_, err = iam.NewGroupPolicy(ctx, inlineName, &iam.GroupPolicyArgs{
			Name:   pulumi.String(policyName),
			Group:  createdGroup.Name,
			Policy: pulumi.String(inlinePolicyString),
		}, pulumi.Provider(provider), pulumi.Parent(createdGroup))
		if err != nil {
			return errors.Wrapf(err, "failed to create inline policy %s", policyName)
		}
	}

	ctx.Export(OpGroupArn, createdGroup.Arn)
	ctx.Export(OpGroupName, createdGroup.Name)
	ctx.Export(OpGroupId, createdGroup.UniqueId)

	return nil
}

// sanitizeForLogicalName reduces an ARN to a stable Pulumi logical-name
// segment: every character outside [A-Za-z0-9] becomes '-'. The full ARN
// stays in the name so distinct policies can never collide (bare policy
// names repeat across paths and accounts).
func sanitizeForLogicalName(s string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, s)
}

// structToJSONString converts a google.protobuf.Struct to a raw JSON string.
func structToJSONString(s *structpb.Struct) (string, error) {
	if s == nil {
		return "{}", nil
	}
	bytes, err := json.Marshal(s.AsMap())
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
