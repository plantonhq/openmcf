package module

import (
	"strings"

	azuremonitoractivitylogalertv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azuremonitoractivitylogalert/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureMonitorActivityLogAlert *azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlert

	ResourceGroupName string

	// ScopeIds are the resolved ARM IDs the alert watches.
	ScopeIds []string

	// Location is the ARM region string for the alert resource ("global"
	// when unspecified).
	Location string

	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertStackInput) *Locals {
	locals := &Locals{}

	locals.AzureMonitorActivityLogAlert = stackInput.Target
	target := stackInput.Target
	spec := target.Spec

	locals.ResourceGroupName = spec.ResourceGroup.GetValue()

	for _, s := range spec.Scopes {
		locals.ScopeIds = append(locals.ScopeIds, s.GetValue())
	}

	locals.Location = activityLogAlertLocation(spec.Location)

	// PARITY-EXCEPTION: resource_kind here is the lowered CloudResourceKind
	// enum string and resource_id is omitted when metadata.id is empty,
	// while the Terraform module emits the family-wide snake-case literal
	// and falls back to metadata.name. Output-neutral (tags never feed stack
	// outputs); aligning the two shapes is a family-wide convention change.
	locals.AzureTags = map[string]string{
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureMonitorActivityLogAlert.String()),
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
	for k, v := range spec.Tags {
		locals.AzureTags[k] = v
	}

	return locals
}

// activityLogAlertLocation maps the location enum to its ARM string,
// defaulting to "global" when unspecified.
func activityLogAlertLocation(l azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLocation) string {
	switch l {
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLocation_WEST_EUROPE:
		return "westeurope"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLocation_NORTH_EUROPE:
		return "northeurope"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLocation_EAST_US_2_EUAP:
		return "eastus2euap"
	default:
		return "global"
	}
}

func categoryString(c azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertCategory) string {
	switch c {
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertCategory_ADMINISTRATIVE:
		return "Administrative"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertCategory_AUTOSCALE:
		return "Autoscale"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertCategory_POLICY:
		return "Policy"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertCategory_RECOMMENDATION:
		return "Recommendation"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertCategory_RESOURCE_HEALTH:
		return "ResourceHealth"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertCategory_SECURITY:
		return "Security"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertCategory_SERVICE_HEALTH:
		return "ServiceHealth"
	}
	return ""
}

func levelString(l azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLevel) string {
	switch l {
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLevel_VERBOSE:
		return "Verbose"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLevel_INFORMATIONAL:
		return "Informational"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLevel_WARNING:
		return "Warning"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLevel_ERROR:
		return "Error"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLevel_CRITICAL:
		return "Critical"
	}
	return ""
}

func levelStrings(ls []azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertLevel) []string {
	out := make([]string, 0, len(ls))
	for _, l := range ls {
		if s := levelString(l); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func recommendationCategoryString(c azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationCategory) string {
	switch c {
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationCategory_COST:
		return "Cost"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationCategory_RELIABILITY:
		return "Reliability"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationCategory_OPERATIONAL_EXCELLENCE:
		return "OperationalExcellence"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationCategory_PERFORMANCE:
		return "Performance"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationCategory_HIGH_AVAILABILITY:
		return "HighAvailability"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationCategory_SECURITY_RECOMMENDATION:
		return "Security"
	}
	return ""
}

func recommendationImpactString(i azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationImpact) string {
	switch i {
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationImpact_HIGH:
		return "High"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationImpact_MEDIUM:
		return "Medium"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertRecommendationImpact_LOW:
		return "Low"
	}
	return ""
}

func healthStatusString(h azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthStatus) string {
	switch h {
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthStatus_AVAILABLE:
		return "Available"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthStatus_DEGRADED:
		return "Degraded"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthStatus_UNAVAILABLE:
		return "Unavailable"
	case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthStatus_UNKNOWN:
		return "Unknown"
	}
	return ""
}

func healthStatusStrings(hs []azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthStatus) []string {
	out := make([]string, 0, len(hs))
	for _, h := range hs {
		if s := healthStatusString(h); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func healthReasonStrings(rs []azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthReason) []string {
	out := make([]string, 0, len(rs))
	for _, r := range rs {
		switch r {
		case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthReason_PLATFORM_INITIATED:
			out = append(out, "PlatformInitiated")
		case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthReason_USER_INITIATED:
			out = append(out, "UserInitiated")
		case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertHealthReason_REASON_UNKNOWN:
			out = append(out, "Unknown")
		}
	}
	return out
}

func serviceHealthEventStrings(es []azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertServiceHealthEvent) []string {
	out := make([]string, 0, len(es))
	for _, e := range es {
		switch e {
		case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertServiceHealthEvent_INCIDENT:
			out = append(out, "Incident")
		case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertServiceHealthEvent_MAINTENANCE:
			out = append(out, "Maintenance")
		case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertServiceHealthEvent_EVENT_INFORMATIONAL:
			out = append(out, "Informational")
		case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertServiceHealthEvent_ACTION_REQUIRED:
			out = append(out, "ActionRequired")
		case azuremonitoractivitylogalertv1.AzureMonitorActivityLogAlertServiceHealthEvent_EVENT_SECURITY:
			out = append(out, "Security")
		}
	}
	return out
}
