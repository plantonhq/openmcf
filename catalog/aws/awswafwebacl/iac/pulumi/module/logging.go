package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/wafv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// logging creates a WAFv2 Web ACL Logging Configuration resource. This is a
// separate AWS resource that links logging to the Web ACL by ARN.
func logging(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, webAclResource *wafv2.WebAcl) error {
	loggingSpec := locals.WebAcl.Spec.Logging

	args := &wafv2.WebAclLoggingConfigurationArgs{
		ResourceArn:           webAclResource.Arn,
		LogDestinationConfigs: pulumi.StringArray{pulumi.String(loggingSpec.DestinationArn.GetValue())},
	}

	// Build redacted fields from the simplified spec fields.
	var redactedFields wafv2.WebAclLoggingConfigurationRedactedFieldArray

	for _, headerName := range loggingSpec.RedactedHeaderNames {
		redactedFields = append(redactedFields, &wafv2.WebAclLoggingConfigurationRedactedFieldArgs{
			SingleHeader: &wafv2.WebAclLoggingConfigurationRedactedFieldSingleHeaderArgs{
				Name: pulumi.String(headerName),
			},
		})
	}

	if loggingSpec.RedactUriPath {
		redactedFields = append(redactedFields, &wafv2.WebAclLoggingConfigurationRedactedFieldArgs{
			UriPath: &wafv2.WebAclLoggingConfigurationRedactedFieldUriPathArgs{},
		})
	}

	if loggingSpec.RedactQueryString {
		redactedFields = append(redactedFields, &wafv2.WebAclLoggingConfigurationRedactedFieldArgs{
			QueryString: &wafv2.WebAclLoggingConfigurationRedactedFieldQueryStringArgs{},
		})
	}

	if loggingSpec.RedactMethod {
		redactedFields = append(redactedFields, &wafv2.WebAclLoggingConfigurationRedactedFieldArgs{
			Method: &wafv2.WebAclLoggingConfigurationRedactedFieldMethodArgs{},
		})
	}

	if len(redactedFields) > 0 {
		args.RedactedFields = redactedFields
	}

	// Log filtering: keep or drop records by the action WAF applied or by
	// labels on the request. Each spec condition carries exactly one of
	// action / label_name (spec CEL), so each renders as a one-armed
	// condition block.
	if filterSpec := loggingSpec.Filter; filterSpec != nil {
		var filters wafv2.WebAclLoggingConfigurationLoggingFilterFilterArray
		for _, filterRule := range filterSpec.Filters {
			var conditions wafv2.WebAclLoggingConfigurationLoggingFilterFilterConditionArray
			for _, condition := range filterRule.Conditions {
				conditionArgs := &wafv2.WebAclLoggingConfigurationLoggingFilterFilterConditionArgs{}
				if condition.Action != "" {
					conditionArgs.ActionCondition = &wafv2.WebAclLoggingConfigurationLoggingFilterFilterConditionActionConditionArgs{
						Action: pulumi.String(condition.Action),
					}
				}
				if condition.LabelName != "" {
					conditionArgs.LabelNameCondition = &wafv2.WebAclLoggingConfigurationLoggingFilterFilterConditionLabelNameConditionArgs{
						LabelName: pulumi.String(condition.LabelName),
					}
				}
				conditions = append(conditions, conditionArgs)
			}
			filters = append(filters, &wafv2.WebAclLoggingConfigurationLoggingFilterFilterArgs{
				Behavior:    pulumi.String(filterRule.Behavior),
				Requirement: pulumi.String(filterRule.Requirement),
				Conditions:  conditions,
			})
		}
		args.LoggingFilter = &wafv2.WebAclLoggingConfigurationLoggingFilterArgs{
			DefaultBehavior: pulumi.String(filterSpec.DefaultBehavior),
			Filters:         filters,
		}
	}

	_, err := wafv2.NewWebAclLoggingConfiguration(
		ctx,
		locals.WebAcl.Metadata.Name+"-logging",
		args,
		pulumi.Provider(provider),
		pulumi.DependsOn([]pulumi.Resource{webAclResource}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create WAF logging configuration")
	}

	return nil
}
