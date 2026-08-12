package module

import (
	azuredatafactorylinkedservicev1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactorylinkedservice/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The five shared optional fields every linked service resource
// carries. Each SDK resource declares its own Args struct, so these
// helpers return the VALUES (nil when unset, leaving the argument
// unsent) and every builder assigns them field by field.

func descriptionPtr(spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec) pulumi.StringPtrInput {
	if spec.Description == "" {
		return nil
	}
	return pulumi.StringPtr(spec.Description)
}

func annotationsArray(spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec) pulumi.StringArrayInput {
	if len(spec.Annotations) == 0 {
		return nil
	}
	return pulumi.ToStringArray(spec.Annotations)
}

func parametersMap(spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec) pulumi.StringMapInput {
	if len(spec.Parameters) == 0 {
		return nil
	}
	return pulumi.ToStringMap(spec.Parameters)
}

func additionalPropertiesMap(spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec) pulumi.StringMapInput {
	if len(spec.AdditionalProperties) == 0 {
		return nil
	}
	return pulumi.ToStringMap(spec.AdditionalProperties)
}

func integrationRuntimeNamePtr(spec *azuredatafactorylinkedservicev1alpha1.AzureDataFactoryLinkedServiceSpec) pulumi.StringPtrInput {
	if spec.IntegrationRuntimeName == nil || spec.IntegrationRuntimeName.GetValue() == "" {
		return nil
	}
	return pulumi.StringPtr(spec.IntegrationRuntimeName.GetValue())
}

// useManagedIdentityOrDefault applies the platform default: managed
// identity stays off unless the spec turns it on.
func useManagedIdentityOrDefault(value *bool) bool {
	if value != nil {
		return *value
	}
	return false
}
