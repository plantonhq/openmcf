package crkreflect

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"google.golang.org/protobuf/proto"
)

// ProviderMeta returns the provider_meta extension (group, display name) of
// the provider that owns one kind.
func ProviderMeta(kind cloudresourcekind.CloudResourceKind) (*cloudresourcekind.CloudResourceProviderMeta, error) {
	kindMeta, err := KindMeta(kind)
	if err != nil {
		return nil, errors.Wrap(err, "while getting cloud resource kind meta")
	}
	return ProviderMetaOf(kindMeta.Provider)
}

// ProviderMetaOf returns the provider_meta extension (group, display name)
// of one provider enum value.
func ProviderMetaOf(provider cloudresourcekind.CloudResourceProvider) (*cloudresourcekind.CloudResourceProviderMeta, error) {
	// Get the descriptor for the enum value (CloudResourceProvider)
	enumValueDescriptor := provider.Descriptor().Values().ByNumber(provider.Number())
	if enumValueDescriptor == nil {
		return nil, errors.Errorf("no descriptor found for provider: %v", provider)
	}

	// Get the options from the enum value descriptor
	options := enumValueDescriptor.Options()
	if options == nil {
		return nil, errors.Errorf("no options found for provider: %v", provider)
	}

	// Extract the meta field from the options
	providerMeta, ok := proto.GetExtension(options, cloudresourcekind.E_ProviderMeta).(*cloudresourcekind.CloudResourceProviderMeta)
	if !ok || providerMeta == nil {
		return nil, errors.Errorf("no meta information found for provider: %v", provider)
	}

	// Return the meta information
	return providerMeta, nil
}
