// Package specpath validates spec-relative proto field paths -- the
// vocabulary catalog sidecars (cost profiles, control profiles) use to name
// the spec fields they describe. A path here is proto field names,
// snake_case, dot-separated, relative to the kind's spec message
// ("default_node_pool.node_count"). Unlike import-value derivation paths
// (pkg/iac/importmap), sidecar paths may traverse and terminate on repeated,
// map, and message fields: "replicas" is a legitimate cost driver precisely
// BECAUSE it is repeated (each entry bills), and a message field can be the
// toggle a control profile points at. The only requirement is that every
// segment exists on the descriptor -- that is what keeps sidecars honest
// across schema revisions: a renamed field orphans the path and fails the
// conformance gate loudly.
package specpath

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// Validate walks dotPath from desc and returns an error naming the first
// segment that does not exist. Traversal continues through message fields
// (including repeated message fields and map values that are messages); a
// path that tries to continue past a scalar fails.
func Validate(desc protoreflect.MessageDescriptor, dotPath string) error {
	_, err := Terminal(desc, dotPath)
	return err
}

// Terminal walks dotPath from desc and returns the terminal field's
// descriptor -- the descriptor-level twin of Resolve, for gates that need
// to reason about a path's kind and cardinality (is this repeated? is it
// numeric?) without a live message.
func Terminal(desc protoreflect.MessageDescriptor, dotPath string) (protoreflect.FieldDescriptor, error) {
	if dotPath == "" {
		return nil, fmt.Errorf("empty path")
	}
	current := desc
	segments := strings.Split(dotPath, ".")
	for i, segment := range segments {
		field := current.Fields().ByName(protoreflect.Name(segment))
		if field == nil {
			return nil, fmt.Errorf("no field %q on %s", segment, current.FullName())
		}
		if i == len(segments)-1 {
			return field, nil
		}
		next := fieldMessage(field)
		if next == nil {
			return nil, fmt.Errorf("segment %q on %s is a scalar but the path continues", segment, current.FullName())
		}
		current = next
	}
	return nil, fmt.Errorf("unreachable")
}

// ResolvableTerminal walks dotPath from desc under Resolve's OWN
// traversal contract -- refusing to continue THROUGH repeated or map
// fields while allowing a repeated terminal -- and returns the terminal
// field's descriptor. Gates over paths that a live evaluator will later
// hand to Resolve must validate with THIS function, not Terminal: the
// permissive Terminal accepts a path Resolve refuses (which element
// would the value come from?), and that disagreement surfaces only at
// replay, as an engine error instead of a CI failure naming the file.
func ResolvableTerminal(desc protoreflect.MessageDescriptor, dotPath string) (protoreflect.FieldDescriptor, error) {
	if dotPath == "" {
		return nil, fmt.Errorf("empty path")
	}
	current := desc
	segments := strings.Split(dotPath, ".")
	for i, segment := range segments {
		field := current.Fields().ByName(protoreflect.Name(segment))
		if field == nil {
			return nil, fmt.Errorf("no field %q on %s", segment, current.FullName())
		}
		if i == len(segments)-1 {
			return field, nil
		}
		if field.IsList() || field.IsMap() {
			return nil, fmt.Errorf("segment %q on %s is %s -- resolution cannot pick an element; only the terminal segment may be repeated",
				segment, current.FullName(), cardinality(field))
		}
		next := fieldMessage(field)
		if next == nil {
			return nil, fmt.Errorf("segment %q on %s is a scalar but the path continues", segment, current.FullName())
		}
		current = next
	}
	return nil, fmt.Errorf("unreachable")
}

// fieldMessage returns the message descriptor a path can continue into, or
// nil for scalar leaves. Map fields continue into their value type when the
// value is a message.
func fieldMessage(field protoreflect.FieldDescriptor) protoreflect.MessageDescriptor {
	if field.IsMap() {
		value := field.MapValue()
		if value.Kind() == protoreflect.MessageKind {
			return value.Message()
		}
		return nil
	}
	if field.Kind() == protoreflect.MessageKind || field.Kind() == protoreflect.GroupKind {
		return field.Message()
	}
	return nil
}

// Resolved is the outcome of resolving a path against a live message: the
// terminal field, its value on this message, and whether the message
// actually carries it. Value is meaningful only when Present is true;
// Field is always the terminal descriptor, so callers can reason about
// kind and cardinality even for absent values.
type Resolved struct {
	Field   protoreflect.FieldDescriptor
	Value   protoreflect.Value
	Present bool
}

// Resolve walks dotPath from msg and returns the terminal field's value.
// It is the evaluation half of Validate: the same path vocabulary, read
// against a live message instead of a descriptor, so validation and
// resolution can never disagree about what a path means.
//
// Presence is explicit-presence semantics: a message field is present when
// set, a repeated field when non-empty, and a scalar when the message
// reports Has (proto3 scalars without optional therefore read as absent at
// their zero value -- callers that treat zero as meaningful read
// Field.Default via Value regardless of Present). An unset INTERMEDIATE
// message resolves as absent with the terminal descriptor intact.
//
// Resolution refuses to traverse THROUGH repeated or map fields -- which
// element would the value come from? -- while terminal repeated fields
// resolve normally (callers count them or demand a single element).
func Resolve(msg protoreflect.Message, dotPath string) (Resolved, error) {
	if dotPath == "" {
		return Resolved{}, fmt.Errorf("empty path")
	}
	current := msg
	currentDesc := msg.Descriptor()
	segments := strings.Split(dotPath, ".")
	for i, segment := range segments {
		field := currentDesc.Fields().ByName(protoreflect.Name(segment))
		if field == nil {
			return Resolved{}, fmt.Errorf("no field %q on %s", segment, currentDesc.FullName())
		}
		terminal := i == len(segments)-1
		if terminal {
			if current == nil {
				return Resolved{Field: field}, nil
			}
			return Resolved{Field: field, Value: current.Get(field), Present: current.Has(field)}, nil
		}
		if field.IsList() || field.IsMap() {
			return Resolved{}, fmt.Errorf("segment %q on %s is %s -- resolution cannot pick an element; only the terminal segment may be repeated",
				segment, currentDesc.FullName(), cardinality(field))
		}
		next := fieldMessage(field)
		if next == nil {
			return Resolved{}, fmt.Errorf("segment %q on %s is a scalar but the path continues", segment, currentDesc.FullName())
		}
		currentDesc = next
		if current != nil && current.Has(field) {
			current = current.Get(field).Message()
		} else {
			current = nil
		}
	}
	return Resolved{}, fmt.Errorf("unreachable")
}

// cardinality names a field's shape for error messages.
func cardinality(field protoreflect.FieldDescriptor) string {
	if field.IsMap() {
		return "a map"
	}
	return "repeated"
}
