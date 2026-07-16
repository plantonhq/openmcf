package module

import (
	"strings"

	azureapplicationinsightsv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureapplicationinsights/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureApplicationInsights *azureapplicationinsightsv1.AzureApplicationInsights
	ResourceGroupName        string
	AzureTags                map[string]string
}

// applicationTypeStrings maps the spec enum to Azure's CASE-SENSITIVE and
// irregular API strings ("Node.JS", "MobileCenter") -- an unmatched value
// would be silently treated as ASP.NET by Azure, which is why the spec
// closes the vocabulary and this map carries the exact wire strings. The
// unspecified row deploys "web" (an unmapped enum would send the empty
// string, which the provider rejects).
var applicationTypeStrings = map[azureapplicationinsightsv1.AzureApplicationInsightsApplicationType]string{
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_azure_application_insights_application_type_unspecified: "web",
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_WEB:                                                     "web",
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_JAVA:                                                    "java",
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_NODE_JS:                                                 "Node.JS",
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_OTHER:                                                   "other",
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_IOS:                                                     "ios",
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_PHONE:                                                   "phone",
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_STORE:                                                   "store",
	azureapplicationinsightsv1.AzureApplicationInsightsApplicationType_MOBILE_CENTER:                                           "MobileCenter",
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureapplicationinsightsv1.AzureApplicationInsightsStackInput) *Locals {
	locals := &Locals{}

	locals.AzureApplicationInsights = stackInput.Target
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
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureApplicationInsights.String()),
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
	// literal "azure_application_insights" for resource_kind and fall back to
	// metadata.name for resource_id, while this module emits the lowered enum
	// string and omits resource_id when metadata.id is unset -- the
	// family-wide tag-shape divergence documented across the Azure catalog.
	// Output-neutral: stack outputs never carry tags.
	for key, value := range target.Spec.Tags {
		locals.AzureTags[key] = value
	}

	return locals
}
