package module

import (
	azuredatafactorydatasetv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactorydataset/v1alpha1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The shared optional fields every dataset resource carries. Each SDK
// resource declares its own Args struct, so these helpers return the
// VALUES (nil when unset, leaving the argument unsent) and every
// builder assigns them field by field.

func descriptionPtr(spec *azuredatafactorydatasetv1alpha1.AzureDataFactoryDatasetSpec) pulumi.StringPtrInput {
	if spec.Description == "" {
		return nil
	}
	return pulumi.StringPtr(spec.Description)
}

func annotationsArray(spec *azuredatafactorydatasetv1alpha1.AzureDataFactoryDatasetSpec) pulumi.StringArrayInput {
	if len(spec.Annotations) == 0 {
		return nil
	}
	return pulumi.ToStringArray(spec.Annotations)
}

func parametersMap(spec *azuredatafactorydatasetv1alpha1.AzureDataFactoryDatasetSpec) pulumi.StringMapInput {
	if len(spec.Parameters) == 0 {
		return nil
	}
	return pulumi.ToStringMap(spec.Parameters)
}

func additionalPropertiesMap(spec *azuredatafactorydatasetv1alpha1.AzureDataFactoryDatasetSpec) pulumi.StringMapInput {
	if len(spec.AdditionalProperties) == 0 {
		return nil
	}
	return pulumi.ToStringMap(spec.AdditionalProperties)
}

func folderPtr(spec *azuredatafactorydatasetv1alpha1.AzureDataFactoryDatasetSpec) pulumi.StringPtrInput {
	if spec.Folder == "" {
		return nil
	}
	return pulumi.StringPtr(spec.Folder)
}

// linkedServiceName resolves the shared name-form linked service
// reference (the wire form 11 of the 13 variants use; the spec
// guarantees it is set for exactly those variants).
func linkedServiceName(spec *azuredatafactorydatasetv1alpha1.AzureDataFactoryDatasetSpec) pulumi.StringInput {
	return pulumi.String(spec.LinkedServiceName.GetValue())
}

// stringPtrWhenSet leaves an optional string argument unsent when the
// spec field is empty (several dataset arguments reject an explicit
// empty string).
func stringPtrWhenSet(value string) pulumi.StringPtrInput {
	if value == "" {
		return nil
	}
	return pulumi.StringPtr(value)
}

// boolPtrWhenTrue sends a flag only when it is on -- unset means
// false on every dynamic_*_enabled flag (the provider's default).
func boolPtrWhenTrue(value bool) pulumi.BoolPtrInput {
	if !value {
		return nil
	}
	return pulumi.BoolPtr(true)
}
