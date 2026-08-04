package crkreflect

import (
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
)

// canonicalKindName lowercases a kind name and strips separators so that
// names differing only in case or separator style compare equal. Registry
// uniqueness is enforced on this form: two kinds whose names collide
// canonically are ambiguous on at least one resolution surface even when
// their exact spellings differ.
func canonicalKindName(name string) string {
	return strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(name))
}

// kindNameIndex maps each kind's effective name (kind_meta.name, falling
// back to the enum value name) to its kind. Built once; construction fails
// if two kinds claim the same name, because a manifest's kind must resolve
// to exactly one type — first-match-wins over an unordered walk resolved
// duplicates randomly.
var kindNameIndex = sync.OnceValues(buildKindNameIndex)

func buildKindNameIndex() (map[string]cloudresourcekind.CloudResourceKind, error) {
	index := make(map[string]cloudresourcekind.CloudResourceKind)
	canonicalOwner := make(map[string]cloudresourcekind.CloudResourceKind)
	for _, kind := range KindsList() {
		effectiveName := kind.String()
		if kindMeta, err := KindMeta(kind); err == nil && kindMeta.Name != "" {
			effectiveName = kindMeta.Name
		}
		canonicalName := canonicalKindName(effectiveName)
		if owner, taken := canonicalOwner[canonicalName]; taken {
			return nil, errors.Errorf(
				"cloud-resource kind registry is ambiguous: %s and %s both answer to the kind name %q; "+
					"kind names must be unique across the registry so every manifest resolves to exactly one kind — rename one of them",
				owner, kind, effectiveName)
		}
		canonicalOwner[canonicalName] = kind
		index[effectiveName] = kind
	}
	return index, nil
}

func KindByKindName(kindName string) (cloudresourcekind.CloudResourceKind, error) {
	index, err := kindNameIndex()
	if err != nil {
		return cloudresourcekind.CloudResourceKind_unspecified, err
	}
	if kind, ok := index[kindName]; ok {
		return kind, nil
	}
	return cloudresourcekind.CloudResourceKind_unspecified,
		errors.Errorf("no matching CloudResourceKind found for kind: %s", kindName)
}
