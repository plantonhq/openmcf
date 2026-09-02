package manifestgraph

import (
	"fmt"

	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const stringValueOrRefFullName = "dev.planton.shared.foreignkey.v1.StringValueOrRef"

// RefUse is one populated valueFrom reference found in a manifest: the spec
// field it sits on (whose descriptor carries the foreign-key annotations),
// the reference itself, and the capability to replace the whole
// StringValueOrRef in place — so collection, validation, edge derivation,
// and literal rewriting all ride ONE traversal and can never disagree about
// what counts as a reference.
type RefUse struct {
	// FieldPath is the proto field-name dot path to the reference,
	// e.g. "spec.kms_key_name", "spec.iam_members[0].member", or
	// "spec.ref_map.primary" for a map entry.
	FieldPath string

	// Field is the field declaring the StringValueOrRef — the carrier of the
	// default_kind / default_kind_field_path / containment_exempt annotations.
	Field protoreflect.FieldDescriptor

	// Ref is the populated reference.
	Ref *foreignkeyv1.ValueFromRef

	// replace swaps the whole StringValueOrRef value at this position.
	replace func(*foreignkeyv1.StringValueOrRef)
}

// Replace substitutes the StringValueOrRef at this reference's position —
// the resolution seam. Passing a literal-armed value is how a resolved
// reference becomes a plain value.
func (u RefUse) Replace(v *foreignkeyv1.StringValueOrRef) {
	u.replace(v)
}

// CollectRefUses walks every populated field of a loaded manifest — singular,
// repeated, nested, and MAP-typed containers included — and returns each
// StringValueOrRef carrying a valueFrom reference. Literal values are not
// references and are skipped.
func CollectRefUses(msg proto.Message) []RefUse {
	var out []RefUse
	walkPopulated(msg.ProtoReflect(), "", &out)
	return out
}

func walkPopulated(m protoreflect.Message, prefix string, out *[]RefUse) {
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		path := string(fd.Name())
		if prefix != "" {
			path = prefix + "." + path
		}
		switch {
		case fd.IsMap():
			if fd.MapValue().Kind() != protoreflect.MessageKind {
				return true
			}
			mp := v.Map()
			mp.Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
				visitMessage(fd, mv.Message(), path+"."+k.String(), out,
					func(replacement *foreignkeyv1.StringValueOrRef) {
						mp.Set(k, protoreflect.ValueOfMessage(replacement.ProtoReflect()))
					})
				return true
			})
		case fd.IsList():
			if fd.Kind() != protoreflect.MessageKind {
				return true
			}
			list := v.List()
			for i := 0; i < list.Len(); i++ {
				idx := i
				visitMessage(fd, list.Get(i).Message(), fmt.Sprintf("%s[%d]", path, i), out,
					func(replacement *foreignkeyv1.StringValueOrRef) {
						list.Set(idx, protoreflect.ValueOfMessage(replacement.ProtoReflect()))
					})
			}
		case fd.Kind() == protoreflect.MessageKind:
			visitMessage(fd, v.Message(), path, out,
				func(replacement *foreignkeyv1.StringValueOrRef) {
					m.Set(fd, protoreflect.ValueOfMessage(replacement.ProtoReflect()))
				})
		}
		return true
	})
}

// visitMessage either records a populated StringValueOrRef reference (the
// annotations are read off the declaring field fd) or recurses.
func visitMessage(fd protoreflect.FieldDescriptor, m protoreflect.Message, path string, out *[]RefUse,
	replace func(*foreignkeyv1.StringValueOrRef)) {
	md := m.Descriptor()
	if string(md.FullName()) == stringValueOrRefFullName {
		svor, ok := m.Interface().(*foreignkeyv1.StringValueOrRef)
		if !ok || svor.GetValueFrom() == nil {
			return
		}
		*out = append(*out, RefUse{FieldPath: path, Field: fd, Ref: svor.GetValueFrom(), replace: replace})
		return
	}
	walkPopulated(m, path, out)
}
