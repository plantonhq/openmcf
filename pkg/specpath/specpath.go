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
	if dotPath == "" {
		return fmt.Errorf("empty path")
	}
	current := desc
	segments := strings.Split(dotPath, ".")
	for i, segment := range segments {
		field := current.Fields().ByName(protoreflect.Name(segment))
		if field == nil {
			return fmt.Errorf("no field %q on %s", segment, current.FullName())
		}
		if i == len(segments)-1 {
			return nil
		}
		next := fieldMessage(field)
		if next == nil {
			return fmt.Errorf("segment %q on %s is a scalar but the path continues", segment, current.FullName())
		}
		current = next
	}
	return nil
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
