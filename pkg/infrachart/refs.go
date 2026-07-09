package infrachart

import (
	"fmt"

	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"github.com/plantonhq/planton/pkg/refcheck"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const stringValueOrRefFullName = "dev.planton.shared.foreignkey.v1.StringValueOrRef"

// refUse is one populated valueFrom reference found in a rendered manifest:
// the spec field it sits on (whose descriptor carries the FK annotations) and
// the reference itself.
type refUse struct {
	// fieldPath is the proto field-name dot path to the reference,
	// e.g. "spec.kms_key_name" or "spec.iam_members[0].member".
	fieldPath string

	// fd is the field declaring the StringValueOrRef — the carrier of the
	// default_kind / default_kind_field_path annotations.
	fd protoreflect.FieldDescriptor

	// ref is the populated reference.
	ref *foreignkeyv1.ValueFromRef
}

// refTarget is a fully resolved reference: which kind and which resource name
// this refUse points at. Used to build the chart's dependency graph.
type refTarget struct {
	kind cloudresourcekind.CloudResourceKind
	name string
}

// collectRefUses walks every populated field of a loaded manifest and returns
// each StringValueOrRef that carries a valueFrom reference. Literal values
// are not references and are skipped.
func collectRefUses(msg proto.Message) []refUse {
	var out []refUse
	walkPopulated(msg.ProtoReflect(), "", &out)
	return out
}

func walkPopulated(m protoreflect.Message, prefix string, out *[]refUse) {
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
			v.Map().Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
				visitMessage(fd, mv.Message(), path+"."+k.String(), out)
				return true
			})
		case fd.IsList():
			if fd.Kind() != protoreflect.MessageKind {
				return true
			}
			list := v.List()
			for i := 0; i < list.Len(); i++ {
				visitMessage(fd, list.Get(i).Message(), fmt.Sprintf("%s[%d]", path, i), out)
			}
		case fd.Kind() == protoreflect.MessageKind:
			visitMessage(fd, v.Message(), path, out)
		}
		return true
	})
}

// visitMessage either records a populated StringValueOrRef reference (the
// annotation is read off the declaring field fd) or recurses.
func visitMessage(fd protoreflect.FieldDescriptor, m protoreflect.Message, path string, out *[]refUse) {
	md := m.Descriptor()
	if string(md.FullName()) == stringValueOrRefFullName {
		svor, ok := m.Interface().(*foreignkeyv1.StringValueOrRef)
		if !ok || svor.GetValueFrom() == nil {
			return
		}
		*out = append(*out, refUse{fieldPath: path, fd: fd, ref: svor.GetValueFrom()})
		return
	}
	walkPopulated(m, path, out)
}

// checkRef validates one valueFrom reference against the FK annotations on
// its declaring field and the referenced kind's proto surface. It returns the
// resolved target (for dependency-graph construction) and any problems found.
//
// The rules, in order:
//
//  1. A reference must have a target kind: an explicit valueFrom.kind, or the
//     field's default_kind annotation.
//  2. A reference must have a field path: an explicit valueFrom.fieldPath, or
//     — only when the target IS the field's default kind — the annotated
//     default_kind_field_path.
//  3. When the target is the field's default kind and the reference spells
//     out a DIFFERENT field path than the annotation, that is an error: the
//     annotated path is the composition key the modules are proven to accept,
//     and overriding it is the id/name/self-link mismatch class that
//     otherwise only surfaces at deploy time.
//  4. The effective field path must resolve against the target kind's actual
//     proto surface (stack outputs, spec, or metadata).
func checkRef(use refUse) (refTarget, []string) {
	var problems []string

	var annotatedKind cloudresourcekind.CloudResourceKind
	var annotatedPath string
	if opts := use.fd.Options(); opts != nil {
		annotatedKind, _ = proto.GetExtension(opts, foreignkeyv1.E_DefaultKind).(cloudresourcekind.CloudResourceKind)
		annotatedPath, _ = proto.GetExtension(opts, foreignkeyv1.E_DefaultKindFieldPath).(string)
	}

	targetKind := use.ref.GetKind()
	if targetKind == cloudresourcekind.CloudResourceKind_unspecified {
		targetKind = annotatedKind
	}
	if targetKind == cloudresourcekind.CloudResourceKind_unspecified {
		problems = append(problems,
			fmt.Sprintf("%s: valueFrom does not name a kind and the field declares no default kind — add an explicit `kind:`", use.fieldPath))
		return refTarget{}, problems
	}

	effectivePath := use.ref.GetFieldPath()
	if effectivePath == "" {
		if targetKind == annotatedKind && annotatedPath != "" {
			effectivePath = annotatedPath
		} else {
			problems = append(problems,
				fmt.Sprintf("%s: valueFrom targets %s but has no fieldPath, and no annotated default applies — add an explicit `fieldPath:`", use.fieldPath, targetKind))
			return refTarget{kind: targetKind, name: use.ref.GetName()}, problems
		}
	} else if targetKind == annotatedKind && annotatedPath != "" && effectivePath != annotatedPath {
		problems = append(problems,
			fmt.Sprintf("%s: valueFrom overrides the annotated composition key for %s — the field's contract is %q but the chart references %q (id/name/self-link format mismatches only surface at deploy time; use the annotated path)",
				use.fieldPath, targetKind, annotatedPath, effectivePath))
	}

	if reason := refcheck.ResolveRefPath(targetKind, effectivePath); reason != "" {
		problems = append(problems,
			fmt.Sprintf("%s: valueFrom fieldPath %q does not resolve on %s: %s", use.fieldPath, effectivePath, targetKind, reason))
	}

	return refTarget{kind: targetKind, name: use.ref.GetName()}, problems
}
