//go:build !codegen
// +build !codegen

// Spec census: the catalog's contract side of the parity measurement. Walks
// every registered kind of one cloud provider and enumerates the leaf fields
// of its spec from proto descriptors -- the exact configurable surface a
// manifest author can express. Descriptor-based on purpose: regex/text
// counting undercounts nested specs and is banned for parity numbers.
//
// The walk is a sibling of pkg/secretcoverage's (registry via crkreflect,
// StringValueOrRef as a leaf, map/list handling, recursion guard) -- a reader
// who knows one walk knows both.

package providerparity

import (
	"sort"

	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// stringValueOrRefFullName is a dependency-wiring leaf: the manifest author
// configures ONE value slot (literal or reference), so it counts as one
// field, never as the message's internals.
const stringValueOrRefFullName = "dev.planton.shared.foreignkey.v1.StringValueOrRef"

// wellKnownValueLeafFullNames are well-known types whose descriptor
// internals are protobuf plumbing, never configurable surface: a
// google.protobuf.Struct field is ONE structured value the author writes
// (e.g. an event pattern authored as YAML), so it counts as one leaf --
// the same treatment as StringValueOrRef. Walking into these would emit
// meaningless paths like "spec.event_pattern.fields.bool_value" that no
// manifest could honestly account.
var wellKnownValueLeafFullNames = map[string]bool{
	"google.protobuf.Struct":    true,
	"google.protobuf.Value":     true,
	"google.protobuf.ListValue": true,
}

// isOpaqueLeafMessage reports whether a message-typed field is recorded as
// one leaf instead of being walked into.
func isOpaqueLeafMessage(fullName string) bool {
	return fullName == stringValueOrRefFullName || wellKnownValueLeafFullNames[fullName]
}

// KindCensus is one kind's spec surface.
type KindCensus struct {
	// Kind is the registry name, e.g. "GcpGcsBucket".
	Kind string
	// SpecFieldPaths is the sorted proto field-name dot paths of every leaf
	// field under spec, rooted at "spec" (matching pkg/secretcoverage's path
	// format), e.g. "spec.lifecycle_rules.condition.age_days".
	SpecFieldPaths []string
}

// SpecCensus enumerates the spec surface of every implemented kind of one
// cloud provider, sorted by kind. Kinds whose API package is not implemented
// yet (registry entry without generated code) are skipped, matching the
// kind-map codegen.
func SpecCensus(provider cloudresourcekind.CloudResourceProvider) []KindCensus {
	var out []KindCensus
	for _, kind := range crkreflect.KindsList() {
		if crkreflect.GetProvider(kind) != provider {
			continue
		}
		msg, err := crkreflect.NewInstance(kind)
		if err != nil {
			continue // enum value exists but the API package is not implemented yet
		}
		specField := msg.ProtoReflect().Descriptor().Fields().ByName("spec")
		if specField == nil || specField.Kind() != protoreflect.MessageKind {
			continue
		}
		out = append(out, KindCensus{
			Kind:           kind.String(),
			SpecFieldPaths: CollectSpecPaths(specField.Message(), "spec"),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Kind < out[j].Kind })
	return out
}

// CollectSpecPaths walks one spec message descriptor and returns its sorted
// leaf field paths. Exposed so tests can drive it against the hermetic
// testcloudresourcegeneric spec in isolation.
func CollectSpecPaths(specMd protoreflect.MessageDescriptor, prefix string) []string {
	var paths []string
	walkSpec(specMd, prefix, map[protoreflect.FullName]bool{}, &paths)
	sort.Strings(paths)
	return paths
}

func walkSpec(md protoreflect.MessageDescriptor, prefix string, visited map[protoreflect.FullName]bool, out *[]string) {
	visited[md.FullName()] = true
	defer delete(visited, md.FullName())

	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		path := string(fd.Name())
		if prefix != "" {
			path = prefix + "." + path
		}

		switch {
		case fd.IsMap():
			// The author configures one map field; message values recurse
			// because each value's shape is itself configurable surface.
			if v := fd.MapValue(); v.Kind() == protoreflect.MessageKind && !isOpaqueLeafMessage(string(v.Message().FullName())) {
				if visited[v.Message().FullName()] {
					*out = append(*out, path)
				} else {
					walkSpec(v.Message(), path, visited, out)
				}
			} else {
				*out = append(*out, path)
			}
		case fd.Kind() == protoreflect.MessageKind:
			switch {
			case isOpaqueLeafMessage(string(fd.Message().FullName())):
				*out = append(*out, path)
			case visited[fd.Message().FullName()]:
				// Recursive re-entry: the field's message is an ancestor on
				// the current walk path (a recursive grammar, e.g. WAF's
				// statement tree nesting statements inside and/or/not). The
				// author configures "the same grammar again" here, so the
				// field counts as ONE leaf — the opaque-leaf treatment.
				// Dropping it silently (the previous behavior) hid the field
				// from the census entirely, which is exactly the
				// silent-omission class this package exists to eliminate.
				*out = append(*out, path)
			default:
				walkSpec(fd.Message(), path, visited, out)
			}
		default:
			*out = append(*out, path)
		}
	}
}
