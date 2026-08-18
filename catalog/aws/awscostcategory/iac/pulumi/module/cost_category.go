package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/costexplorer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// costCategory provisions the Cost Explorer cost category and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - the category's name is spec.category_name (an explicit field -
//     names legally carry spaces metadata.name cannot) and changing it
//     replaces the category;
//   - rule_version is module-pinned to "CostCategoryExpression.v1",
//     the only value the AWS API accepts - never a spec knob;
//   - rules are ORDERED (first match wins), so the module renders them
//     exactly in spec order;
//   - the expression tree is leveled (root -> node -> leaf) to the
//     exact nesting AWS accepts; the conversion in expression.go is
//     1:1 with no depth checks because the spec shape cannot express
//     an illegal tree.
func costCategory(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.Spec

	rules := costexplorer.CostCategoryRuleArray{}
	for _, r := range spec.Rules {
		rule := costexplorer.CostCategoryRuleArgs{}
		if r.Type != "" {
			rule.Type = pulumi.StringPtr(r.Type)
		}
		if r.Value != "" {
			rule.Value = pulumi.StringPtr(r.Value)
		}
		if r.Rule != nil {
			rule.Rule = buildRuleExpression(r.Rule)
		}
		if r.InheritedValue != nil {
			inherited := &costexplorer.CostCategoryRuleInheritedValueArgs{
				DimensionName: pulumi.StringPtr(r.InheritedValue.DimensionName),
			}
			if r.InheritedValue.DimensionKey != "" {
				inherited.DimensionKey = pulumi.StringPtr(r.InheritedValue.DimensionKey)
			}
			rule.InheritedValue = inherited
		}
		rules = append(rules, rule)
	}

	args := &costexplorer.CostCategoryArgs{
		Name:        pulumi.StringPtr(spec.CategoryName),
		RuleVersion: pulumi.String("CostCategoryExpression.v1"),
		Rules:       rules,
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.DefaultValue != "" {
		args.DefaultValue = pulumi.StringPtr(spec.DefaultValue)
	}
	if spec.EffectiveStart != "" {
		args.EffectiveStart = pulumi.StringPtr(spec.EffectiveStart)
	}

	if len(spec.SplitChargeRules) > 0 {
		splitRules := costexplorer.CostCategorySplitChargeRuleArray{}
		for _, s := range spec.SplitChargeRules {
			splitRule := costexplorer.CostCategorySplitChargeRuleArgs{
				Method:  pulumi.String(s.Method),
				Source:  pulumi.String(s.Source),
				Targets: pulumi.ToStringArray(s.Targets),
			}
			if len(s.Parameters) > 0 {
				parameters := costexplorer.CostCategorySplitChargeRuleParameterArray{}
				for _, p := range s.Parameters {
					parameters = append(parameters, costexplorer.CostCategorySplitChargeRuleParameterArgs{
						Type:   pulumi.StringPtr(p.Type),
						Values: pulumi.ToStringArray(p.Values),
					})
				}
				splitRule.Parameters = parameters
			}
			splitRules = append(splitRules, splitRule)
		}
		args.SplitChargeRules = splitRules
	}

	createdCategory, err := costexplorer.NewCostCategory(ctx, "cost-category", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create cost category")
	}

	ctx.Export(OpCategoryArn, createdCategory.Arn)
	ctx.Export(OpCategoryName, createdCategory.Name)
	ctx.Export(OpEffectiveStart, createdCategory.EffectiveStart)
	ctx.Export(OpEffectiveEnd, createdCategory.EffectiveEnd)

	return nil
}
