package module

import (
	"strings"

	azuremonitormetricalertv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuremonitormetricalert/v1alpha1"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMonitorMetricAlert *azuremonitormetricalertv1alpha1.AzureMonitorMetricAlert
	ResourceGroupName       string
	AzureTags               map[string]string
}

// aggregationStrings maps the aggregation enum to ARM's wire values.
var aggregationStrings = map[azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertAggregation]string{
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertAggregation_AVERAGE: "Average",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertAggregation_COUNT:   "Count",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertAggregation_MINIMUM: "Minimum",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertAggregation_MAXIMUM: "Maximum",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertAggregation_TOTAL:   "Total",
}

// operatorStrings is one shared vocabulary for both criteria families (the
// spec CELs keep each family to its legal subset).
var operatorStrings = map[azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertOperator]string{
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertOperator_EQUALS:                "Equals",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertOperator_GREATER_THAN:          "GreaterThan",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertOperator_GREATER_THAN_OR_EQUAL: "GreaterThanOrEqual",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertOperator_LESS_THAN:             "LessThan",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertOperator_LESS_THAN_OR_EQUAL:    "LessThanOrEqual",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertOperator_GREATER_OR_LESS_THAN:  "GreaterOrLessThan",
}

// dimensionOperatorStrings maps the dimension-filter enum to ARM's values.
var dimensionOperatorStrings = map[azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertDimensionOperator]string{
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertDimensionOperator_INCLUDE:     "Include",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertDimensionOperator_EXCLUDE:     "Exclude",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertDimensionOperator_STARTS_WITH: "StartsWith",
}

// sensitivityStrings maps the dynamic-threshold sensitivity enum to ARM's
// values.
var sensitivityStrings = map[azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertSensitivity]string{
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertSensitivity_LOW:    "Low",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertSensitivity_MEDIUM: "Medium",
	azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertSensitivity_HIGH:   "High",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremonitormetricalertv1alpha1.AzureMonitorMetricAlertStackInput) *Locals {
	locals := &Locals{}

	locals.AzureMonitorMetricAlert = stackInput.Target
	target := stackInput.Target

	// The resource_group field is a StringValueOrRef. The platform middleware
	// resolves valueFrom references before IaC modules run, so .GetValue()
	// always returns the resolved literal string.
	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// Identity tags derived from metadata; user tags merge OVER these (the
	// governance surface belongs to the user).
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureMonitorMetricAlert.String()),
	}

	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}

	// PARITY-EXCEPTION: the Terraform module's base tags use the snake_case
	// literal "azure_monitor_metric_alert" for resource_kind and fall back to
	// metadata.name for resource_id, while this module emits the lowered enum
	// string and omits resource_id when metadata.id is unset -- the
	// family-wide tag-shape divergence documented across the Azure catalog.
	// Output-neutral: stack outputs never carry tags.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
