package runner

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/outputs"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"sigs.k8s.io/yaml"
)

const stringValueOrRefFullName = "dev.planton.shared.foreignkey.v1.StringValueOrRef"

// ResolveManifestRefs implements, for the standalone E2E harness, the foreign-key
// resolution the Planton orchestrator performs in production: it replaces each
// value_from reference in the component manifest with the literal value read from
// a deployed prerequisite's outputs. Standalone Planton otherwise requires literal
// values -- the tofu generator errors on an unresolved ref and the pulumi modules
// drop it -- so this is the step that makes a composed topology testable end to end.
//
// Resolution walks the entire spec message tree: singular StringValueOrRef fields,
// nested messages, and repeated message elements (e.g. backends[].group). Each ref
// resolves by the ref's explicit valueFrom.kind + fieldPath when set, falling back
// to the field descriptor's default_kind + default_kind_field_path.
//
// The resolved manifest is written to a temp file whose path is returned; the
// original is left untouched. When there is nothing to resolve, the original path
// is returned unchanged.
func ResolveManifestRefs(manifestPath string, depOutputs map[cloudresourcekind.CloudResourceKind]map[string]interface{}) (string, error) {
	if len(depOutputs) == 0 {
		return manifestPath, nil
	}

	manifestObject, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load manifest for ref resolution from %s", manifestPath)
	}

	// Flatten each prerequisite's outputs to dotted keys so a fieldPath like
	// "status.outputs.self_link" resolves the way the platform flattens outputs.
	flattened := make(map[cloudresourcekind.CloudResourceKind]map[string]string, len(depOutputs))
	for kind, out := range depOutputs {
		flattened[kind] = outputs.Flatten(out)
	}

	top := manifestObject.ProtoReflect()
	specFd := top.Descriptor().Fields().ByName("spec")
	if specFd == nil || specFd.Kind() != protoreflect.MessageKind {
		return manifestPath, nil
	}

	resolvedAny, err := resolveRefsInMessage(top.Mutable(specFd).Message(), flattened)
	if err != nil {
		return "", err
	}
	if !resolvedAny {
		return manifestPath, nil
	}

	jsonBytes, err := protojson.Marshal(manifestObject)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal resolved manifest")
	}
	yamlBytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert resolved manifest to yaml")
	}

	// A temp file (not next to the scenario) so scenario discovery never picks it up.
	tmpFile, err := os.CreateTemp("", "planton-e2e-resolved-*.yaml")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp file for resolved manifest")
	}
	if _, err := tmpFile.Write(yamlBytes); err != nil {
		tmpFile.Close()
		return "", errors.Wrap(err, "failed to write resolved manifest")
	}
	if err := tmpFile.Close(); err != nil {
		return "", errors.Wrap(err, "failed to close resolved manifest")
	}
	return tmpFile.Name(), nil
}

// resolveRefsInMessage walks msg and replaces value_from arms on every
// StringValueOrRef field (including inside nested and repeated messages) with
// literals from the matching prerequisite's outputs. Returns whether any field
// was resolved.
func resolveRefsInMessage(msg protoreflect.Message, flattened map[cloudresourcekind.CloudResourceKind]map[string]string) (bool, error) {
	resolvedAny := false
	fields := msg.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)

		if fd.Kind() == protoreflect.MessageKind {
			if string(fd.Message().FullName()) == stringValueOrRefFullName {
				if !msg.Has(fd) {
					continue
				}
				ref, ok := msg.Get(fd).Message().Interface().(*foreignkeyv1.StringValueOrRef)
				if !ok || ref.GetValueFrom() == nil {
					continue
				}
				val, err := resolveRefValue(fd, ref, flattened)
				if err != nil {
					return false, err
				}
				if val == "" {
					continue
				}
				resolved := &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
				}
				msg.Set(fd, protoreflect.ValueOfMessage(resolved.ProtoReflect()))
				resolvedAny = true
				continue
			}

			if fd.IsList() {
				list := msg.Get(fd).List()
				for j := 0; j < list.Len(); j++ {
					elem := list.Get(j).Message()
					if elem.IsValid() {
						changed, err := resolveRefsInMessage(elem, flattened)
						if err != nil {
							return false, err
						}
						if changed {
							resolvedAny = true
						}
					}
				}
				continue
			}

			if !fd.IsMap() && msg.Has(fd) {
				changed, err := resolveRefsInMessage(msg.Mutable(fd).Message(), flattened)
				if err != nil {
					return false, err
				}
				if changed {
					resolvedAny = true
				}
			}
		}
	}
	return resolvedAny, nil
}

// resolveRefValue looks up the literal for a value_from reference. The explicit
// kind and fieldPath on the ref take precedence; when kind is unspecified the
// field's default_kind annotation is used.
func resolveRefValue(fd protoreflect.FieldDescriptor, ref *foreignkeyv1.StringValueOrRef, flattened map[cloudresourcekind.CloudResourceKind]map[string]string) (string, error) {
	valueFrom := ref.GetValueFrom()
	if valueFrom == nil {
		return "", nil
	}

	kind := valueFrom.GetKind()
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		opts := fd.Options()
		if opts != nil {
			kind, _ = proto.GetExtension(opts, foreignkeyv1.E_DefaultKind).(cloudresourcekind.CloudResourceKind)
		}
	}
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		return "", nil
	}

	outs, ok := flattened[kind]
	if !ok {
		return "", nil
	}

	path := valueFrom.GetFieldPath()
	if path == "" {
		opts := fd.Options()
		if opts != nil {
			path, _ = proto.GetExtension(opts, foreignkeyv1.E_DefaultKindFieldPath).(string)
		}
	}
	key := strings.TrimPrefix(path, "status.outputs.")
	val, ok := outs[key]
	if !ok {
		return "", errors.Errorf("prerequisite %s has no output %q to resolve field %q", kind, key, fd.Name())
	}
	return val, nil
}
