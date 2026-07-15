//go:build !codegen
// +build !codegen

package refcheck

import (
	"regexp"
	"strings"

	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ResolveValueFromPath validates a manifest-level valueFrom reference (a ValueFromRef's
// kind + field_path) against the referenced kind's descriptor and returns a human-readable
// reason when it does not resolve, or "" when it does.
//
// This is the manifest-side sibling of the annotation check in Analyze: the annotation
// check guards `default_kind_field_path` values authored in protos, while this guards
// `fieldPath` values authored in YAML (infra-chart templates, presets). The resolution
// semantics deliberately mirror the control plane's reference validator so a reference
// accepted here is accepted there:
//   - the path is rooted at the kind's top-level API message (so it starts with
//     "status.", "spec.", or "metadata.");
//   - segments resolve by exact proto field name first, then by a camelCase->snake_case
//     fallback (chart authors write either form);
//   - bracketed or bare integer segments index into repeated message fields;
//   - every non-terminal segment must be a message, and the terminal field must be a
//     string (a valueFrom always feeds a string-valued reference).
func ResolveValueFromPath(kind cloudresourcekind.CloudResourceKind, fieldPath string) string {
	inst, err := crkreflect.NewInstance(kind)
	if err != nil {
		return "kind " + kind.String() + " is not a registered/implemented kind"
	}
	if fieldPath == "" {
		return "fieldPath is empty"
	}

	current := inst.ProtoReflect().Descriptor()
	var fd protoreflect.FieldDescriptor
	segments := strings.Split(fieldPath, ".")

	for i, seg := range segments {
		if isIndexSegment(seg) {
			if fd == nil || !fd.IsList() || fd.Kind() != protoreflect.MessageKind {
				return "cannot index into '" + seg + "' -- preceding field is not a repeated message"
			}
			current = fd.Message()
			continue
		}

		fd = current.Fields().ByName(protoreflect.Name(seg))
		if fd == nil {
			fd = current.Fields().ByName(protoreflect.Name(camelToSnake(seg)))
		}
		if fd == nil {
			return "no field named '" + seg + "' on " + string(current.FullName())
		}
		if i < len(segments)-1 {
			if fd.Kind() != protoreflect.MessageKind {
				return "cannot descend into scalar field '" + seg + "'"
			}
			current = fd.Message()
		}
	}

	if fd.Kind() != protoreflect.StringKind {
		return "terminal field '" + fieldPath + "' is " + fd.Kind().String() + ", expected string"
	}
	return ""
}

var camelBoundary = regexp.MustCompile(`([a-z0-9])([A-Z])`)

// camelToSnake converts camelCase to snake_case, matching the control plane's fallback
// so chart authors may write either the protojson (camelCase) or proto (snake_case) form.
func camelToSnake(s string) string {
	return strings.ToLower(camelBoundary.ReplaceAllString(s, "${1}_${2}"))
}
