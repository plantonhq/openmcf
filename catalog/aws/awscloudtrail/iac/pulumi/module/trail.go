package module

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cloudtrail"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// trail creates the CloudTrail trail (and the optional organization
// delegated-admin registration) and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - AWS validates the delivery bucket's POLICY at create ("Incorrect
//     S3 bucket policy is detected") -- the bucket policy is the
//     consumer's contract (AwsS3Bucket spec.policy), never this
//     module's;
//   - the classic and advanced selector styles are mutually exclusive
//     on a trail (the spec CEL guarantees only one arrives here);
//   - AWS expects the CloudWatch group ARN in its ":*" suffix form;
//     the module appends the suffix when the referenced value lacks
//     it, so both engines send the identical ARN;
//   - the delegated-admin registration is an ACCOUNT-GLOBAL act with
//     its own lifecycle (deregistered on destroy) -- it has no
//     structural edge to the trail resource.
func trail(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &cloudtrail.TrailArgs{
		// metadata.name is the trail name on both engines (AWS: 3-128
		// chars, letters/digits/._-, starts and ends alphanumeric).
		Name:         pulumi.String(locals.Target.Metadata.Name),
		S3BucketName: pulumi.String(spec.S3BucketName.GetValue()),
		Tags:         pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.S3KeyPrefix != "" {
		args.S3KeyPrefix = pulumi.String(spec.S3KeyPrefix)
	}
	if spec.IsMultiRegionTrail {
		args.IsMultiRegionTrail = pulumi.Bool(true)
	}
	if spec.IsOrganizationTrail {
		args.IsOrganizationTrail = pulumi.Bool(true)
	}
	// Rendered only on an explicit choice so the module never fights
	// the provider defaults (both default true).
	if spec.IncludeGlobalServiceEvents != nil {
		args.IncludeGlobalServiceEvents = pulumi.Bool(*spec.IncludeGlobalServiceEvents)
	}
	if spec.EnableLogging != nil {
		args.EnableLogging = pulumi.Bool(*spec.EnableLogging)
	}
	if spec.EnableLogFileValidation {
		args.EnableLogFileValidation = pulumi.Bool(true)
	}
	if spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
	}
	if spec.SnsTopicName.GetValue() != "" {
		args.SnsTopicName = pulumi.String(spec.SnsTopicName.GetValue())
	}

	if spec.CloudwatchLogs != nil {
		// AWS expects "arn:...:log-group:<name>:*" -- normalize once so
		// a bare group ARN reference renders identically on both
		// engines.
		groupArn := spec.CloudwatchLogs.LogGroupArn.GetValue()
		if !strings.HasSuffix(groupArn, ":*") {
			groupArn += ":*"
		}
		args.CloudWatchLogsGroupArn = pulumi.String(groupArn)
		args.CloudWatchLogsRoleArn = pulumi.String(spec.CloudwatchLogs.RoleArn.GetValue())
	}

	// Classic selectors: management scope plus coarse data-event
	// scopes.
	var eventSelectors cloudtrail.TrailEventSelectorArray
	for _, s := range spec.EventSelectors {
		sel := &cloudtrail.TrailEventSelectorArgs{}
		if s.ReadWriteType != "" {
			sel.ReadWriteType = pulumi.String(s.ReadWriteType)
		} else {
			sel.ReadWriteType = pulumi.String("All")
		}
		if s.IncludeManagementEvents != nil {
			sel.IncludeManagementEvents = pulumi.Bool(*s.IncludeManagementEvents)
		}
		if len(s.ExcludeManagementEventSources) > 0 {
			sel.ExcludeManagementEventSources = pulumi.ToStringArray(s.ExcludeManagementEventSources)
		}
		var dataResources cloudtrail.TrailEventSelectorDataResourceArray
		for _, d := range s.DataResources {
			dataResources = append(dataResources, &cloudtrail.TrailEventSelectorDataResourceArgs{
				Type:   pulumi.String(d.Type),
				Values: pulumi.ToStringArray(d.Values),
			})
		}
		if len(dataResources) > 0 {
			sel.DataResources = dataResources
		}
		eventSelectors = append(eventSelectors, sel)
	}
	if len(eventSelectors) > 0 {
		args.EventSelectors = eventSelectors
	}

	// Advanced selectors: fine-grained field matching.
	var advancedSelectors cloudtrail.TrailAdvancedEventSelectorArray
	for _, s := range spec.AdvancedEventSelectors {
		var fieldSelectors cloudtrail.TrailAdvancedEventSelectorFieldSelectorArray
		for _, f := range s.FieldSelectors {
			fs := &cloudtrail.TrailAdvancedEventSelectorFieldSelectorArgs{
				Field: pulumi.String(f.Field),
			}
			if len(f.Equals) > 0 {
				fs.Equals = pulumi.ToStringArray(f.Equals)
			}
			if len(f.NotEquals) > 0 {
				fs.NotEquals = pulumi.ToStringArray(f.NotEquals)
			}
			if len(f.StartsWith) > 0 {
				fs.StartsWiths = pulumi.ToStringArray(f.StartsWith)
			}
			if len(f.NotStartsWith) > 0 {
				fs.NotStartsWiths = pulumi.ToStringArray(f.NotStartsWith)
			}
			if len(f.EndsWith) > 0 {
				fs.EndsWiths = pulumi.ToStringArray(f.EndsWith)
			}
			if len(f.NotEndsWith) > 0 {
				fs.NotEndsWiths = pulumi.ToStringArray(f.NotEndsWith)
			}
			fieldSelectors = append(fieldSelectors, fs)
		}
		selArgs := &cloudtrail.TrailAdvancedEventSelectorArgs{
			FieldSelectors: fieldSelectors,
		}
		if s.Name != "" {
			selArgs.Name = pulumi.String(s.Name)
		}
		advancedSelectors = append(advancedSelectors, selArgs)
	}
	if len(advancedSelectors) > 0 {
		args.AdvancedEventSelectors = advancedSelectors
	}

	// Insights engines (anomaly detection; billed separately).
	var insightSelectors cloudtrail.TrailInsightSelectorArray
	for _, t := range spec.InsightTypes {
		insightSelectors = append(insightSelectors, &cloudtrail.TrailInsightSelectorArgs{
			InsightType: pulumi.String(t),
		})
	}
	if len(insightSelectors) > 0 {
		args.InsightSelectors = insightSelectors
	}

	createdTrail, err := cloudtrail.NewTrail(ctx, "trail", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create trail")
	}

	// The organization's delegated CloudTrail administrator - an
	// account-global registration (one per organization, performed
	// from the management account). Reads resolve through the
	// Organizations API, so the caller needs organizations:Describe*
	// alongside cloudtrail:RegisterOrganizationDelegatedAdmin.
	if spec.OrganizationDelegatedAdminAccountId != "" {
		_, err := cloudtrail.NewOrganizationDelegatedAdminAccount(ctx, "delegated-admin", &cloudtrail.OrganizationDelegatedAdminAccountArgs{
			AccountId: pulumi.String(spec.OrganizationDelegatedAdminAccountId),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "register organization delegated admin")
		}
	}

	ctx.Export(OpTrailArn, createdTrail.Arn)
	ctx.Export(OpHomeRegion, createdTrail.HomeRegion)
	ctx.Export(OpSnsTopicArn, createdTrail.SnsTopicArn)
	return nil
}
