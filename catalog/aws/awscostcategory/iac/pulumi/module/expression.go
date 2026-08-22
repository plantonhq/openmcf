package module

import (
	awscostcategoryv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscostcategory/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/costexplorer"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The spec's expression tree is LEVELED (root -> node -> leaf),
// mirroring the exact nesting AWS accepts - and the provider SDK
// mirrors the same cap by generating a DISTINCT Go type per tree
// path (root, And, Or, Not, AndAnd, AndOr, ...). The builders below
// are therefore a mechanical 1:1 walk: one function per path, no
// depth checks anywhere, because neither side can express an illegal
// tree. The Terraform module renders the same leveled variables
// natively.

func buildRuleExpression(root *awscostcategoryv1alpha1.AwsCostCategoryExpression) *costexplorer.CostCategoryRuleRuleArgs {
	out := &costexplorer.CostCategoryRuleRuleArgs{}
	if root.Dimension != nil {
		out.Dimension = &costexplorer.CostCategoryRuleRuleDimensionArgs{
			Key:          pulumi.StringPtr(root.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(root.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(root.Dimension.Values),
		}
	}
	if root.Tag != nil {
		out.Tags = &costexplorer.CostCategoryRuleRuleTagsArgs{
			Key:          pulumi.StringPtr(root.Tag.Key),
			MatchOptions: pulumi.ToStringArray(root.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(root.Tag.Values),
		}
	}
	if root.CostCategory != nil {
		out.CostCategory = &costexplorer.CostCategoryRuleRuleCostCategoryArgs{
			Key:          pulumi.StringPtr(root.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(root.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(root.CostCategory.Values),
		}
	}
	if len(root.And) > 0 {
		ands := costexplorer.CostCategoryRuleRuleAndArray{}
		for _, node := range root.And {
			ands = append(ands, buildAndNode(node))
		}
		out.Ands = ands
	}
	if len(root.Or) > 0 {
		ors := costexplorer.CostCategoryRuleRuleOrArray{}
		for _, node := range root.Or {
			ors = append(ors, buildOrNode(node))
		}
		out.Ors = ors
	}
	if root.Not != nil {
		out.Not = buildNotNode(root.Not)
	}
	return out
}

func buildAndNode(node *awscostcategoryv1alpha1.AwsCostCategoryExpressionNode) costexplorer.CostCategoryRuleRuleAndArgs {
	out := costexplorer.CostCategoryRuleRuleAndArgs{}
	if node.Dimension != nil {
		out.Dimension = &costexplorer.CostCategoryRuleRuleAndDimensionArgs{
			Key:          pulumi.StringPtr(node.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(node.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(node.Dimension.Values),
		}
	}
	if node.Tag != nil {
		out.Tags = &costexplorer.CostCategoryRuleRuleAndTagsArgs{
			Key:          pulumi.StringPtr(node.Tag.Key),
			MatchOptions: pulumi.ToStringArray(node.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(node.Tag.Values),
		}
	}
	if node.CostCategory != nil {
		out.CostCategory = &costexplorer.CostCategoryRuleRuleAndCostCategoryArgs{
			Key:          pulumi.StringPtr(node.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(node.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(node.CostCategory.Values),
		}
	}
	if len(node.And) > 0 {
		children := costexplorer.CostCategoryRuleRuleAndAndArray{}
		for _, leaf := range node.And {
			child := costexplorer.CostCategoryRuleRuleAndAndArgs{}
			if leaf.Dimension != nil {
				child.Dimension = &costexplorer.CostCategoryRuleRuleAndAndDimensionArgs{
					Key:          pulumi.StringPtr(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &costexplorer.CostCategoryRuleRuleAndAndTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategory = &costexplorer.CostCategoryRuleRuleAndAndCostCategoryArgs{
					Key:          pulumi.StringPtr(leaf.CostCategory.Key),
					MatchOptions: pulumi.ToStringArray(leaf.CostCategory.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.CostCategory.Values),
				}
			}
			children = append(children, child)
		}
		out.Ands = children
	}
	if len(node.Or) > 0 {
		children := costexplorer.CostCategoryRuleRuleAndOrArray{}
		for _, leaf := range node.Or {
			child := costexplorer.CostCategoryRuleRuleAndOrArgs{}
			if leaf.Dimension != nil {
				child.Dimension = &costexplorer.CostCategoryRuleRuleAndOrDimensionArgs{
					Key:          pulumi.StringPtr(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &costexplorer.CostCategoryRuleRuleAndOrTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategory = &costexplorer.CostCategoryRuleRuleAndOrCostCategoryArgs{
					Key:          pulumi.StringPtr(leaf.CostCategory.Key),
					MatchOptions: pulumi.ToStringArray(leaf.CostCategory.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.CostCategory.Values),
				}
			}
			children = append(children, child)
		}
		out.Ors = children
	}
	if node.Not != nil {
		child := &costexplorer.CostCategoryRuleRuleAndNotArgs{}
		if node.Not.Dimension != nil {
			child.Dimension = &costexplorer.CostCategoryRuleRuleAndNotDimensionArgs{
				Key:          pulumi.StringPtr(node.Not.Dimension.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Dimension.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Dimension.Values),
			}
		}
		if node.Not.Tag != nil {
			child.Tags = &costexplorer.CostCategoryRuleRuleAndNotTagsArgs{
				Key:          pulumi.StringPtr(node.Not.Tag.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Tag.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Tag.Values),
			}
		}
		if node.Not.CostCategory != nil {
			child.CostCategory = &costexplorer.CostCategoryRuleRuleAndNotCostCategoryArgs{
				Key:          pulumi.StringPtr(node.Not.CostCategory.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.CostCategory.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.CostCategory.Values),
			}
		}
		out.Not = child
	}
	return out
}

func buildOrNode(node *awscostcategoryv1alpha1.AwsCostCategoryExpressionNode) costexplorer.CostCategoryRuleRuleOrArgs {
	out := costexplorer.CostCategoryRuleRuleOrArgs{}
	if node.Dimension != nil {
		out.Dimension = &costexplorer.CostCategoryRuleRuleOrDimensionArgs{
			Key:          pulumi.StringPtr(node.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(node.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(node.Dimension.Values),
		}
	}
	if node.Tag != nil {
		out.Tags = &costexplorer.CostCategoryRuleRuleOrTagsArgs{
			Key:          pulumi.StringPtr(node.Tag.Key),
			MatchOptions: pulumi.ToStringArray(node.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(node.Tag.Values),
		}
	}
	if node.CostCategory != nil {
		out.CostCategory = &costexplorer.CostCategoryRuleRuleOrCostCategoryArgs{
			Key:          pulumi.StringPtr(node.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(node.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(node.CostCategory.Values),
		}
	}
	if len(node.And) > 0 {
		children := costexplorer.CostCategoryRuleRuleOrAndArray{}
		for _, leaf := range node.And {
			child := costexplorer.CostCategoryRuleRuleOrAndArgs{}
			if leaf.Dimension != nil {
				child.Dimension = &costexplorer.CostCategoryRuleRuleOrAndDimensionArgs{
					Key:          pulumi.StringPtr(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &costexplorer.CostCategoryRuleRuleOrAndTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategory = &costexplorer.CostCategoryRuleRuleOrAndCostCategoryArgs{
					Key:          pulumi.StringPtr(leaf.CostCategory.Key),
					MatchOptions: pulumi.ToStringArray(leaf.CostCategory.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.CostCategory.Values),
				}
			}
			children = append(children, child)
		}
		out.Ands = children
	}
	if len(node.Or) > 0 {
		children := costexplorer.CostCategoryRuleRuleOrOrArray{}
		for _, leaf := range node.Or {
			child := costexplorer.CostCategoryRuleRuleOrOrArgs{}
			if leaf.Dimension != nil {
				child.Dimension = &costexplorer.CostCategoryRuleRuleOrOrDimensionArgs{
					Key:          pulumi.StringPtr(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &costexplorer.CostCategoryRuleRuleOrOrTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategory = &costexplorer.CostCategoryRuleRuleOrOrCostCategoryArgs{
					Key:          pulumi.StringPtr(leaf.CostCategory.Key),
					MatchOptions: pulumi.ToStringArray(leaf.CostCategory.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.CostCategory.Values),
				}
			}
			children = append(children, child)
		}
		out.Ors = children
	}
	if node.Not != nil {
		child := &costexplorer.CostCategoryRuleRuleOrNotArgs{}
		if node.Not.Dimension != nil {
			child.Dimension = &costexplorer.CostCategoryRuleRuleOrNotDimensionArgs{
				Key:          pulumi.StringPtr(node.Not.Dimension.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Dimension.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Dimension.Values),
			}
		}
		if node.Not.Tag != nil {
			child.Tags = &costexplorer.CostCategoryRuleRuleOrNotTagsArgs{
				Key:          pulumi.StringPtr(node.Not.Tag.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Tag.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Tag.Values),
			}
		}
		if node.Not.CostCategory != nil {
			child.CostCategory = &costexplorer.CostCategoryRuleRuleOrNotCostCategoryArgs{
				Key:          pulumi.StringPtr(node.Not.CostCategory.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.CostCategory.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.CostCategory.Values),
			}
		}
		out.Not = child
	}
	return out
}

func buildNotNode(node *awscostcategoryv1alpha1.AwsCostCategoryExpressionNode) *costexplorer.CostCategoryRuleRuleNotArgs {
	out := &costexplorer.CostCategoryRuleRuleNotArgs{}
	if node.Dimension != nil {
		out.Dimension = &costexplorer.CostCategoryRuleRuleNotDimensionArgs{
			Key:          pulumi.StringPtr(node.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(node.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(node.Dimension.Values),
		}
	}
	if node.Tag != nil {
		out.Tags = &costexplorer.CostCategoryRuleRuleNotTagsArgs{
			Key:          pulumi.StringPtr(node.Tag.Key),
			MatchOptions: pulumi.ToStringArray(node.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(node.Tag.Values),
		}
	}
	if node.CostCategory != nil {
		out.CostCategory = &costexplorer.CostCategoryRuleRuleNotCostCategoryArgs{
			Key:          pulumi.StringPtr(node.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(node.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(node.CostCategory.Values),
		}
	}
	if len(node.And) > 0 {
		children := costexplorer.CostCategoryRuleRuleNotAndArray{}
		for _, leaf := range node.And {
			child := costexplorer.CostCategoryRuleRuleNotAndArgs{}
			if leaf.Dimension != nil {
				child.Dimension = &costexplorer.CostCategoryRuleRuleNotAndDimensionArgs{
					Key:          pulumi.StringPtr(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &costexplorer.CostCategoryRuleRuleNotAndTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategory = &costexplorer.CostCategoryRuleRuleNotAndCostCategoryArgs{
					Key:          pulumi.StringPtr(leaf.CostCategory.Key),
					MatchOptions: pulumi.ToStringArray(leaf.CostCategory.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.CostCategory.Values),
				}
			}
			children = append(children, child)
		}
		out.Ands = children
	}
	if len(node.Or) > 0 {
		children := costexplorer.CostCategoryRuleRuleNotOrArray{}
		for _, leaf := range node.Or {
			child := costexplorer.CostCategoryRuleRuleNotOrArgs{}
			if leaf.Dimension != nil {
				child.Dimension = &costexplorer.CostCategoryRuleRuleNotOrDimensionArgs{
					Key:          pulumi.StringPtr(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &costexplorer.CostCategoryRuleRuleNotOrTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategory = &costexplorer.CostCategoryRuleRuleNotOrCostCategoryArgs{
					Key:          pulumi.StringPtr(leaf.CostCategory.Key),
					MatchOptions: pulumi.ToStringArray(leaf.CostCategory.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.CostCategory.Values),
				}
			}
			children = append(children, child)
		}
		out.Ors = children
	}
	if node.Not != nil {
		child := &costexplorer.CostCategoryRuleRuleNotNotArgs{}
		if node.Not.Dimension != nil {
			child.Dimension = &costexplorer.CostCategoryRuleRuleNotNotDimensionArgs{
				Key:          pulumi.StringPtr(node.Not.Dimension.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Dimension.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Dimension.Values),
			}
		}
		if node.Not.Tag != nil {
			child.Tags = &costexplorer.CostCategoryRuleRuleNotNotTagsArgs{
				Key:          pulumi.StringPtr(node.Not.Tag.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Tag.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Tag.Values),
			}
		}
		if node.Not.CostCategory != nil {
			child.CostCategory = &costexplorer.CostCategoryRuleRuleNotNotCostCategoryArgs{
				Key:          pulumi.StringPtr(node.Not.CostCategory.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.CostCategory.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.CostCategory.Values),
			}
		}
		out.Not = child
	}
	return out
}
