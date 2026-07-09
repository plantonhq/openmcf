package explain

import (
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Dispatcher continues path resolution across a kind-valued boundary: a
// field whose VALUE selects another schema (e.g. a platform envelope whose
// `kind` enum types an otherwise-unstructured payload field). When path
// resolution stands on a field a dispatcher claims and segments remain, the
// next segment is handed to the dispatcher, which resolves it to a new
// resource root; the walk then continues inside that resource.
type Dispatcher interface {
	// Claims reports whether this dispatcher owns continuation below fd.
	Claims(fd protoreflect.FieldDescriptor) bool
	// Resolve maps one path segment (e.g. a kind name) to the resource it
	// selects. The error should name the segment and how to list valid ones.
	Resolve(segment string) (Resource, error)
	// Hint is appended to the field's constraints when the field is rendered
	// WITHOUT a further segment, so the drill-down path is discoverable
	// (e.g. "shape determined by spec.kind -- drill with ...").
	Hint() string
}

// resolvePath walks dotted protojson segments from the resource envelope and
// renders the resolved node. Segment names are the exact keys a manifest
// author writes (protojson names); proto names are accepted as a fallback so
// paths copied from proto sources also resolve.
func (e *Engine) resolvePath(res Resource, path []string) (*Report, error) {
	md := res.Message
	var fd protoreflect.FieldDescriptor

	for i, segment := range path {
		// Dispatch takes precedence over structural descent: the claimed
		// field is typically a scalar (an enum discriminator), so this must
		// run before the no-fields-to-drill-into check below.
		if fd != nil {
			if dispatcher := e.dispatcherFor(fd); dispatcher != nil {
				next, err := dispatcher.Resolve(segment)
				if err != nil {
					return nil, err
				}
				report, err := e.Explain(next, path[i+1:])
				if err != nil {
					return nil, err
				}
				// Re-anchor the report on the path the user typed so the
				// answer names what they asked for.
				report.Path = joinPath(res.Name, path)
				return report, nil
			}
		}
		if md == nil {
			return nil, errors.Errorf(
				"%s.%s does not resolve: %q is a %s, which has no fields to drill into",
				res.Name, strings.Join(path[:i], "."), path[i-1], fd.Kind())
		}

		next := fieldBySegment(md, segment)
		if next == nil {
			return nil, errors.Errorf(
				"%s has no field %q under %q -- run `explain %s` to see its fields",
				res.Name, segment, joinPath(res.Name, path[:i]), joinPath(res.Name, path[:i]))
		}
		fd = next
		md = messageBehind(fd)
	}

	// The recursion guard starts EMPTY here: renderField routes the field's
	// own message through expandMessage, which marks it on entry. Pre-seeding
	// it (as the root view does for walking a message's direct fields) would
	// make every resolved message field report itself as recursive.
	field := e.renderField(fd, map[protoreflect.FullName]bool{})
	if dispatcher := e.dispatcherFor(fd); dispatcher != nil && dispatcher.Hint() != "" {
		field.Constraints = append(field.Constraints, dispatcher.Hint())
	}
	return &Report{
		Kind:       res.Name,
		ApiVersion: res.ApiVersion,
		Path:       joinPath(res.Name, path),
		Field:      &field,
	}, nil
}

func (e *Engine) dispatcherFor(fd protoreflect.FieldDescriptor) Dispatcher {
	for _, d := range e.Dispatchers {
		if d.Claims(fd) {
			return d
		}
	}
	return nil
}

// fieldBySegment matches a path segment to a field: protojson name first
// (the manifest surface), proto name as fallback.
func fieldBySegment(md protoreflect.MessageDescriptor, segment string) protoreflect.FieldDescriptor {
	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		if fields.Get(i).JSONName() == segment {
			return fields.Get(i)
		}
	}
	return fields.ByName(protoreflect.Name(segment))
}

// messageBehind returns the message to continue walking into: the value
// type for maps, the element type for lists, the message itself otherwise;
// nil for scalars and for wrappers explain treats as terminal.
func messageBehind(fd protoreflect.FieldDescriptor) protoreflect.MessageDescriptor {
	var md protoreflect.MessageDescriptor
	if fd.IsMap() {
		md = fd.MapValue().Message()
	} else {
		md = fd.Message()
	}
	if md == nil || !expandable(md) {
		return nil
	}
	return md
}

func joinPath(name string, segments []string) string {
	if len(segments) == 0 {
		return name
	}
	return name + "." + strings.Join(segments, ".")
}
