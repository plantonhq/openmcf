package crkreflect

import "github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"

func ProvidersList() []cloudresourcekind.CloudResourceProvider {
	resp := make([]cloudresourcekind.CloudResourceProvider, 0)
	// Iterate over all the enum values in ApiResourceKind
	for _, enumValue := range cloudresourcekind.CloudResourceProvider_value {
		resp = append(resp, cloudresourcekind.CloudResourceProvider(enumValue))
	}
	return resp
}
