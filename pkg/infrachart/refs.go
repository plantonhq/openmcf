package infrachart

import (
	"fmt"

	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"github.com/plantonhq/planton/pkg/refcheck"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// RefError is one valueFrom reference in a rendered manifest that cannot resolve against
// the referenced kind -- a composition that would fail at deploy time.
type RefError struct {
	// FieldPath is the proto field-name dot path of the referencing field in the manifest.
	FieldPath string
	// TargetKind is the referenced kind as written in the manifest.
	TargetKind string
	// RefPath is the fieldPath written in the manifest's valueFrom block.
	RefPath string
	// Reason is why the reference does not resolve.
	Reason string
}

func (e RefError) Error() string {
	return fmt.Sprintf("%s -> %s %q: %s", e.FieldPath, e.TargetKind, e.RefPath, e.Reason)
}

// CheckValueFromRefs walks a loaded manifest and validates every populated
// StringValueOrRef.valueFrom: the referenced kind must be registered and the fieldPath
// must resolve to a string field on it. The walk mirrors the control plane's reference
// validator (recurse into every message, repeated message, and map value; treat
// StringValueOrRef as a leaf) so a chart accepted offline is accepted at publish time.
func CheckValueFromRefs(manifest proto.Message) []RefError {
	var errs []RefError
	walkRefs(manifest.ProtoReflect(), "", func(path string, vf *foreignkeyv1.ValueFromRef) {
		checkRef(vf, path, &errs)
	})
	return errs
}

// collectValueFromRefs returns every populated valueFrom reference in a loaded
// manifest with its field path -- the raw sites, unvalidated, for checks that need
// to see the reference targets themselves (e.g. intra-chart target resolution).
func collectValueFromRefs(manifest proto.Message) []refSite {
	var sites []refSite
	walkRefs(manifest.ProtoReflect(), "", func(path string, vf *foreignkeyv1.ValueFromRef) {
		sites = append(sites, refSite{FieldPath: path, Ref: vf})
	})
	return sites
}

// refSite is one populated valueFrom reference at a field path in a manifest.
type refSite struct {
	FieldPath string
	Ref       *foreignkeyv1.ValueFromRef
}

// walkRefs visits every populated StringValueOrRef.valueFrom in the message tree,
// recursing into nested messages, repeated messages, and map values.
func walkRefs(m protoreflect.Message, prefix string, visit func(path string, vf *foreignkeyv1.ValueFromRef)) {
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
			v.Map().Range(func(mk protoreflect.MapKey, mv protoreflect.Value) bool {
				visitMessageValue(mv.Message(), path+"."+mk.String(), visit)
				return true
			})
		case fd.IsList():
			if fd.Kind() != protoreflect.MessageKind {
				return true
			}
			list := v.List()
			for i := 0; i < list.Len(); i++ {
				visitMessageValue(list.Get(i).Message(), fmt.Sprintf("%s[%d]", path, i), visit)
			}
		case fd.Kind() == protoreflect.MessageKind:
			visitMessageValue(v.Message(), path, visit)
		}
		return true
	})
}

// visitMessageValue visits a message value: a StringValueOrRef leaf's populated
// valueFrom is handed to the visitor, any other message is recursed into.
func visitMessageValue(m protoreflect.Message, path string, visit func(path string, vf *foreignkeyv1.ValueFromRef)) {
	if svor, ok := m.Interface().(*foreignkeyv1.StringValueOrRef); ok {
		if vf := svor.GetValueFrom(); vf != nil {
			visit(path, vf)
		}
		return
	}
	walkRefs(m, path, visit)
}

func checkRef(vf *foreignkeyv1.ValueFromRef, path string, out *[]RefError) {
	if vf.GetKind() == cloudresourcekind.CloudResourceKind_unspecified {
		*out = append(*out, RefError{
			FieldPath:  path,
			TargetKind: vf.GetKind().String(),
			RefPath:    vf.GetFieldPath(),
			Reason:     "valueFrom.kind is missing or not a known kind",
		})
		return
	}
	if reason := refcheck.ResolveValueFromPath(vf.GetKind(), vf.GetFieldPath()); reason != "" {
		*out = append(*out, RefError{
			FieldPath:  path,
			TargetKind: vf.GetKind().String(),
			RefPath:    vf.GetFieldPath(),
			Reason:     reason,
		})
	}
}
