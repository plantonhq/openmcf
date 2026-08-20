package module

import (
	awsbudgetv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbudget/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/budgets"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The spec's filter tree is LEVELED (root -> node -> leaf), mirroring
// the exact nesting AWS accepts on budgets - and the provider SDK
// mirrors the same cap by generating a DISTINCT Go type per tree path
// (root, And, Or, Not, AndAnd, AndOr, ...). The builders below are
// therefore a mechanical 1:1 walk: one function per path, no depth
// checks anywhere, because neither side can express an illegal tree.
// The Terraform module renders the same leveled variables natively.

func buildFilterExpression(root *awsbudgetv1alpha1.AwsBudgetFilterExpression) *budgets.BudgetFilterExpressionArgs {
	out := &budgets.BudgetFilterExpressionArgs{}
	if root.Dimension != nil {
		out.Dimensions = &budgets.BudgetFilterExpressionDimensionsArgs{
			Key:          pulumi.String(root.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(root.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(root.Dimension.Values),
		}
	}
	if root.Tag != nil {
		out.Tags = &budgets.BudgetFilterExpressionTagsArgs{
			Key:          pulumi.StringPtr(root.Tag.Key),
			MatchOptions: pulumi.ToStringArray(root.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(root.Tag.Values),
		}
	}
	if root.CostCategory != nil {
		out.CostCategories = &budgets.BudgetFilterExpressionCostCategoriesArgs{
			Key:          pulumi.StringPtr(root.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(root.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(root.CostCategory.Values),
		}
	}
	if len(root.And) > 0 {
		ands := budgets.BudgetFilterExpressionAndArray{}
		for _, node := range root.And {
			ands = append(ands, buildFilterAndNode(node))
		}
		out.Ands = ands
	}
	if len(root.Or) > 0 {
		ors := budgets.BudgetFilterExpressionOrArray{}
		for _, node := range root.Or {
			ors = append(ors, buildFilterOrNode(node))
		}
		out.Ors = ors
	}
	if root.Not != nil {
		out.Not = buildFilterNotNode(root.Not)
	}
	return out
}

func buildFilterAndNode(node *awsbudgetv1alpha1.AwsBudgetFilterExpressionNode) budgets.BudgetFilterExpressionAndArgs {
	out := budgets.BudgetFilterExpressionAndArgs{}
	if node.Dimension != nil {
		out.Dimensions = &budgets.BudgetFilterExpressionAndDimensionsArgs{
			Key:          pulumi.String(node.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(node.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(node.Dimension.Values),
		}
	}
	if node.Tag != nil {
		out.Tags = &budgets.BudgetFilterExpressionAndTagsArgs{
			Key:          pulumi.StringPtr(node.Tag.Key),
			MatchOptions: pulumi.ToStringArray(node.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(node.Tag.Values),
		}
	}
	if node.CostCategory != nil {
		out.CostCategories = &budgets.BudgetFilterExpressionAndCostCategoriesArgs{
			Key:          pulumi.StringPtr(node.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(node.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(node.CostCategory.Values),
		}
	}
	if len(node.And) > 0 {
		children := budgets.BudgetFilterExpressionAndAndArray{}
		for _, leaf := range node.And {
			child := budgets.BudgetFilterExpressionAndAndArgs{}
			if leaf.Dimension != nil {
				child.Dimensions = &budgets.BudgetFilterExpressionAndAndDimensionsArgs{
					Key:          pulumi.String(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &budgets.BudgetFilterExpressionAndAndTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategories = &budgets.BudgetFilterExpressionAndAndCostCategoriesArgs{
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
		children := budgets.BudgetFilterExpressionAndOrArray{}
		for _, leaf := range node.Or {
			child := budgets.BudgetFilterExpressionAndOrArgs{}
			if leaf.Dimension != nil {
				child.Dimensions = &budgets.BudgetFilterExpressionAndOrDimensionsArgs{
					Key:          pulumi.String(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &budgets.BudgetFilterExpressionAndOrTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategories = &budgets.BudgetFilterExpressionAndOrCostCategoriesArgs{
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
		child := &budgets.BudgetFilterExpressionAndNotArgs{}
		if node.Not.Dimension != nil {
			child.Dimensions = &budgets.BudgetFilterExpressionAndNotDimensionsArgs{
				Key:          pulumi.String(node.Not.Dimension.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Dimension.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Dimension.Values),
			}
		}
		if node.Not.Tag != nil {
			child.Tags = &budgets.BudgetFilterExpressionAndNotTagsArgs{
				Key:          pulumi.StringPtr(node.Not.Tag.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Tag.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Tag.Values),
			}
		}
		if node.Not.CostCategory != nil {
			child.CostCategories = &budgets.BudgetFilterExpressionAndNotCostCategoriesArgs{
				Key:          pulumi.StringPtr(node.Not.CostCategory.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.CostCategory.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.CostCategory.Values),
			}
		}
		out.Not = child
	}
	return out
}

func buildFilterOrNode(node *awsbudgetv1alpha1.AwsBudgetFilterExpressionNode) budgets.BudgetFilterExpressionOrArgs {
	out := budgets.BudgetFilterExpressionOrArgs{}
	if node.Dimension != nil {
		out.Dimensions = &budgets.BudgetFilterExpressionOrDimensionsArgs{
			Key:          pulumi.String(node.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(node.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(node.Dimension.Values),
		}
	}
	if node.Tag != nil {
		out.Tags = &budgets.BudgetFilterExpressionOrTagsArgs{
			Key:          pulumi.StringPtr(node.Tag.Key),
			MatchOptions: pulumi.ToStringArray(node.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(node.Tag.Values),
		}
	}
	if node.CostCategory != nil {
		out.CostCategories = &budgets.BudgetFilterExpressionOrCostCategoriesArgs{
			Key:          pulumi.StringPtr(node.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(node.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(node.CostCategory.Values),
		}
	}
	if len(node.And) > 0 {
		children := budgets.BudgetFilterExpressionOrAndArray{}
		for _, leaf := range node.And {
			child := budgets.BudgetFilterExpressionOrAndArgs{}
			if leaf.Dimension != nil {
				child.Dimensions = &budgets.BudgetFilterExpressionOrAndDimensionsArgs{
					Key:          pulumi.String(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &budgets.BudgetFilterExpressionOrAndTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategories = &budgets.BudgetFilterExpressionOrAndCostCategoriesArgs{
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
		children := budgets.BudgetFilterExpressionOrOrArray{}
		for _, leaf := range node.Or {
			child := budgets.BudgetFilterExpressionOrOrArgs{}
			if leaf.Dimension != nil {
				child.Dimensions = &budgets.BudgetFilterExpressionOrOrDimensionsArgs{
					Key:          pulumi.String(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &budgets.BudgetFilterExpressionOrOrTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategories = &budgets.BudgetFilterExpressionOrOrCostCategoriesArgs{
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
		child := &budgets.BudgetFilterExpressionOrNotArgs{}
		if node.Not.Dimension != nil {
			child.Dimensions = &budgets.BudgetFilterExpressionOrNotDimensionsArgs{
				Key:          pulumi.String(node.Not.Dimension.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Dimension.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Dimension.Values),
			}
		}
		if node.Not.Tag != nil {
			child.Tags = &budgets.BudgetFilterExpressionOrNotTagsArgs{
				Key:          pulumi.StringPtr(node.Not.Tag.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Tag.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Tag.Values),
			}
		}
		if node.Not.CostCategory != nil {
			child.CostCategories = &budgets.BudgetFilterExpressionOrNotCostCategoriesArgs{
				Key:          pulumi.StringPtr(node.Not.CostCategory.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.CostCategory.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.CostCategory.Values),
			}
		}
		out.Not = child
	}
	return out
}

func buildFilterNotNode(node *awsbudgetv1alpha1.AwsBudgetFilterExpressionNode) *budgets.BudgetFilterExpressionNotArgs {
	out := &budgets.BudgetFilterExpressionNotArgs{}
	if node.Dimension != nil {
		out.Dimensions = &budgets.BudgetFilterExpressionNotDimensionsArgs{
			Key:          pulumi.String(node.Dimension.Key),
			MatchOptions: pulumi.ToStringArray(node.Dimension.MatchOptions),
			Values:       pulumi.ToStringArray(node.Dimension.Values),
		}
	}
	if node.Tag != nil {
		out.Tags = &budgets.BudgetFilterExpressionNotTagsArgs{
			Key:          pulumi.StringPtr(node.Tag.Key),
			MatchOptions: pulumi.ToStringArray(node.Tag.MatchOptions),
			Values:       pulumi.ToStringArray(node.Tag.Values),
		}
	}
	if node.CostCategory != nil {
		out.CostCategories = &budgets.BudgetFilterExpressionNotCostCategoriesArgs{
			Key:          pulumi.StringPtr(node.CostCategory.Key),
			MatchOptions: pulumi.ToStringArray(node.CostCategory.MatchOptions),
			Values:       pulumi.ToStringArray(node.CostCategory.Values),
		}
	}
	if len(node.And) > 0 {
		children := budgets.BudgetFilterExpressionNotAndArray{}
		for _, leaf := range node.And {
			child := budgets.BudgetFilterExpressionNotAndArgs{}
			if leaf.Dimension != nil {
				child.Dimensions = &budgets.BudgetFilterExpressionNotAndDimensionsArgs{
					Key:          pulumi.String(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &budgets.BudgetFilterExpressionNotAndTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategories = &budgets.BudgetFilterExpressionNotAndCostCategoriesArgs{
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
		children := budgets.BudgetFilterExpressionNotOrArray{}
		for _, leaf := range node.Or {
			child := budgets.BudgetFilterExpressionNotOrArgs{}
			if leaf.Dimension != nil {
				child.Dimensions = &budgets.BudgetFilterExpressionNotOrDimensionsArgs{
					Key:          pulumi.String(leaf.Dimension.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Dimension.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Dimension.Values),
				}
			}
			if leaf.Tag != nil {
				child.Tags = &budgets.BudgetFilterExpressionNotOrTagsArgs{
					Key:          pulumi.StringPtr(leaf.Tag.Key),
					MatchOptions: pulumi.ToStringArray(leaf.Tag.MatchOptions),
					Values:       pulumi.ToStringArray(leaf.Tag.Values),
				}
			}
			if leaf.CostCategory != nil {
				child.CostCategories = &budgets.BudgetFilterExpressionNotOrCostCategoriesArgs{
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
		child := &budgets.BudgetFilterExpressionNotNotArgs{}
		if node.Not.Dimension != nil {
			child.Dimensions = &budgets.BudgetFilterExpressionNotNotDimensionsArgs{
				Key:          pulumi.String(node.Not.Dimension.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Dimension.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Dimension.Values),
			}
		}
		if node.Not.Tag != nil {
			child.Tags = &budgets.BudgetFilterExpressionNotNotTagsArgs{
				Key:          pulumi.StringPtr(node.Not.Tag.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.Tag.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.Tag.Values),
			}
		}
		if node.Not.CostCategory != nil {
			child.CostCategories = &budgets.BudgetFilterExpressionNotNotCostCategoriesArgs{
				Key:          pulumi.StringPtr(node.Not.CostCategory.Key),
				MatchOptions: pulumi.ToStringArray(node.Not.CostCategory.MatchOptions),
				Values:       pulumi.ToStringArray(node.Not.CostCategory.Values),
			}
		}
		out.Not = child
	}
	return out
}
