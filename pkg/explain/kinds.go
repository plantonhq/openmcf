package explain

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/protodocs"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// DefaultEngine explains this repo's own schema universe: the cloud-resource
// kinds, with the shared option family and the embedded proto docs. Hosts
// composing additional universes (platform APIs, extra option families,
// kind-valued dispatchers) construct their own Engine and reuse the pieces.
func DefaultEngine() *Engine {
	return &Engine{
		Interpreters: []OptionInterpreter{SharedOptions},
		Docs:         protodocs.Lookup,
	}
}

// KindResource builds the explainable Resource for a cloud-resource kind
// from the descriptors compiled into the binary.
func KindResource(kind cloudresourcekind.CloudResourceKind) (Resource, error) {
	instance, err := crkreflect.NewInstance(kind)
	if err != nil {
		return Resource{}, err
	}
	return Resource{
		Name:       kind.String(),
		ApiVersion: crkreflect.GroupVersion(kind),
		Message:    instance.ProtoReflect().Descriptor(),
	}, nil
}

// ResolveKindName resolves a user-typed kind name (AwsVpc, aws-vpc, aws_vpc)
// to its explainable Resource.
func ResolveKindName(name string) (Resource, error) {
	kind := crkreflect.KindFromString(name)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		return Resource{}, errors.Errorf("unknown cloud resource kind %q -- run `planton explain --list` to see all kinds", name)
	}
	return KindResource(kind)
}

// KindNames returns every explainable kind name, sorted.
func KindNames() []string {
	kinds := crkreflect.KindsList()
	names := make([]string, 0, len(kinds))
	for _, k := range kinds {
		names = append(names, k.String())
	}
	sort.Strings(names)
	return names
}
