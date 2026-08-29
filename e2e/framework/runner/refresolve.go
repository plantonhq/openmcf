package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/outputs"
	"github.com/plantonhq/planton/pkg/protobufyaml"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const stringValueOrRefFullName = "dev.planton.shared.foreignkey.v1.StringValueOrRef"

// DependencyOutputs holds the captured stack outputs of every deployed
// prerequisite, keyed by kind and then by manifest name. The second level exists
// because a prerequisite install profile may deploy several instances of the
// same kind (e.g. the two different-AZ subnets a load balancer requires), and a
// reference must be able to pick the specific instance it means.
type DependencyOutputs map[cloudresourcekind.CloudResourceKind]map[string]map[string]interface{}

// ResolveManifestRefs implements, for the standalone E2E harness, the foreign-key
// resolution the Planton orchestrator performs in production: it replaces each
// value_from reference in the component manifest that matches a deployed
// prerequisite with the literal value read from that prerequisite's outputs.
// Standalone Planton otherwise requires literal values -- the tofu generator
// errors on an unresolved ref and the pulumi modules drop it -- so this is the
// step that makes a composed (e.g. subnet -> vpc) topology testable end to end.
//
// Resolution mirrors the production resolver's semantics: the reference's own
// kind and field_path win when present, falling back to the field's default_kind
// and default_kind_field_path annotations. The instance is selected by the
// reference's name; when the name matches no deployed instance but exactly one
// instance of the kind exists, that sole instance is used -- scenario manifests
// name their references after real-world topology, while the shared install
// profiles have fixed names, and forcing them to agree would couple every
// scenario to the profile file.
//
// The resolved manifest is written to a temp file whose path is returned; the
// original is left untouched. When there is nothing to resolve, the original path
// is returned unchanged.
//
// Scope: StringValueOrRef fields anywhere in the spec tree -- singular, repeated,
// and nested inside (repeated) messages, e.g. a listener action's target-group
// refs. Each element of a repeated field resolves independently, so a list can
// mix literals (say, an AWS-managed policy ARN) with references to deployed
// prerequisites. Map-typed fields are not traversed (no spec models refs inside
// maps today).
func ResolveManifestRefs(manifestPath string, depOutputs DependencyOutputs) (string, error) {
	if len(depOutputs) == 0 {
		return manifestPath, nil
	}

	manifestObject, err := manifest.LoadManifest(manifestPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load manifest for ref resolution from %s", manifestPath)
	}

	// Flatten each prerequisite's outputs to dotted keys so a field path
	// like "status.outputs.vpc_id" resolves the way the platform flattens outputs.
	flattened := make(map[cloudresourcekind.CloudResourceKind]map[string]map[string]string, len(depOutputs))
	for kind, instances := range depOutputs {
		flattened[kind] = make(map[string]map[string]string, len(instances))
		for name, out := range instances {
			flattened[kind][name] = outputs.Flatten(out)
		}
	}

	top := manifestObject.ProtoReflect()
	specFd := top.Descriptor().Fields().ByName("spec")
	if specFd == nil || specFd.Kind() != protoreflect.MessageKind {
		return manifestPath, nil
	}

	// The visitor replaces each resolvable value_from arm with the literal read
	// from the matching prerequisite's outputs; unmatched refs stay untouched.
	resolvedAny, err := forEachRefField(top.Mutable(specFd).Message(), func(fd protoreflect.FieldDescriptor, ref *foreignkeyv1.StringValueOrRef) (*foreignkeyv1.StringValueOrRef, error) {
		if ref.GetValueFrom() == nil {
			return nil, nil
		}
		val, resolved, err := lookupRefValue(fd, ref, flattened)
		if err != nil || !resolved {
			return nil, err
		}
		return &foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
		}, nil
	})
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
	yamlBytes, err := protobufyaml.JSONToYAML(jsonBytes)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert resolved manifest to yaml")
	}

	// A temp file (not next to the scenario, so discovery never picks it up)
	// that KEEPS the scenario's basename: verifier dispatch keys behavioral
	// variants off the scenario name in the manifest path, and a
	// random-only temp name would silently demote every reference-carrying
	// behavioral scenario to its plain verifier (the same identity contract
	// the token-expansion copy honors).
	tmpDir, err := os.MkdirTemp("", "planton-e2e-resolved-*")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temp dir for resolved manifest")
	}
	tmpPath := filepath.Join(tmpDir, filepath.Base(manifestPath))
	if err := os.WriteFile(tmpPath, yamlBytes, 0o600); err != nil {
		return "", errors.Wrap(err, "failed to write resolved manifest")
	}
	return tmpPath, nil
}

// refFieldVisitor is invoked for every StringValueOrRef value in a spec tree,
// literal-valued ones included (visitors filter on GetValueFrom themselves).
// Returning a non-nil replacement swaps the value in place; returning (nil,
// nil) leaves it untouched.
type refFieldVisitor func(fd protoreflect.FieldDescriptor, ref *foreignkeyv1.StringValueOrRef) (*foreignkeyv1.StringValueOrRef, error)

// forEachRefField is the single traversal over a spec tree's StringValueOrRef
// fields -- singular, repeated, and nested inside (repeated) messages. Both the
// deploy-time resolver above and the offline fixture-integrity checker walk
// through here, so "what counts as a reference" can never diverge between
// them. Each element of a repeated ref field is visited independently, so a
// list can mix literals with references. Map-typed fields are not traversed
// (no spec models refs inside maps today). Returns whether any value was
// replaced.
func forEachRefField(msg protoreflect.Message, visit refFieldVisitor) (bool, error) {
	changedAny := false
	fields := msg.Descriptor().Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if fd.Kind() != protoreflect.MessageKind || fd.IsMap() {
			continue
		}

		isRef := string(fd.Message().FullName()) == stringValueOrRefFullName

		if fd.IsList() {
			if !msg.Has(fd) {
				continue
			}
			list := msg.Mutable(fd).List()
			for j := 0; j < list.Len(); j++ {
				if !isRef {
					// A repeated nested message (e.g. a listener's default
					// actions): recurse into each element.
					changed, err := forEachRefField(list.Get(j).Message(), visit)
					if err != nil {
						return false, err
					}
					changedAny = changedAny || changed
					continue
				}
				ref, ok := list.Get(j).Message().Interface().(*foreignkeyv1.StringValueOrRef)
				if !ok {
					continue
				}
				replacement, err := visit(fd, ref)
				if err != nil {
					return false, err
				}
				if replacement == nil {
					continue
				}
				list.Set(j, protoreflect.ValueOfMessage(replacement.ProtoReflect()))
				changedAny = true
			}
			continue
		}

		if !msg.Has(fd) {
			continue
		}

		if !isRef {
			// A singular nested message (e.g. a listener's tls config): recurse.
			changed, err := forEachRefField(msg.Mutable(fd).Message(), visit)
			if err != nil {
				return false, err
			}
			changedAny = changedAny || changed
			continue
		}

		ref, ok := msg.Get(fd).Message().Interface().(*foreignkeyv1.StringValueOrRef)
		if !ok {
			continue
		}
		replacement, err := visit(fd, ref)
		if err != nil {
			return false, err
		}
		if replacement == nil {
			continue
		}
		msg.Set(fd, protoreflect.ValueOfMessage(replacement.ProtoReflect()))
		changedAny = true
	}
	return changedAny, nil
}

// lookupRefValue resolves one value_from reference against the deployed
// prerequisites' flattened outputs. The reference's explicit kind/field_path win
// over the field's default_kind/default_kind_field_path annotations, matching
// how the production resolver treats a fully-specified reference. Returns
// (value, true, nil) on success and (_, false, nil) when no prerequisite of the
// referenced kind was deployed -- in which case the ref is left untouched. Two
// failure modes ARE errors, so misdeclared scenarios fail loudly rather than
// silently skipping: a name that matches none of several deployed instances
// (ambiguous), and a matched instance missing the referenced output.
func lookupRefValue(fd protoreflect.FieldDescriptor, ref *foreignkeyv1.StringValueOrRef, flattened map[cloudresourcekind.CloudResourceKind]map[string]map[string]string) (string, bool, error) {
	valueFrom := ref.GetValueFrom()
	if valueFrom == nil {
		return "", false, nil
	}

	kind := referencedKind(fd, ref)
	if kind == cloudresourcekind.CloudResourceKind_unspecified {
		return "", false, nil
	}

	instances, ok := flattened[kind]
	if !ok || len(instances) == 0 {
		return "", false, nil
	}

	outs, ok := instances[valueFrom.GetName()]
	if !ok {
		if len(instances) > 1 {
			return "", false, errors.Errorf(
				"reference %q on field %q matches none of the %d deployed %s prerequisites by name",
				valueFrom.GetName(), fd.Name(), len(instances), kind)
		}
		for _, sole := range instances {
			outs = sole
		}
	}

	path := valueFrom.GetFieldPath()
	if path == "" && fd.Options() != nil {
		path, _ = proto.GetExtension(fd.Options(), foreignkeyv1.E_DefaultKindFieldPath).(string)
	}
	key := strings.TrimPrefix(path, "status.outputs.")
	val, ok := outs[key]
	if !ok {
		return "", false, errors.Errorf("prerequisite %s has no output %q to resolve field %q", kind, key, fd.Name())
	}
	return val, true, nil
}

// referencedKind determines which kind a value_from reference points at: the
// reference's own explicit kind wins, falling back to the field's default_kind
// annotation. Returns unspecified for a bare polymorphic reference that
// carries neither -- such a reference can never resolve (the resolver leaves
// it untouched and the module deploys without the value).
func referencedKind(fd protoreflect.FieldDescriptor, ref *foreignkeyv1.StringValueOrRef) cloudresourcekind.CloudResourceKind {
	kind := ref.GetValueFrom().GetKind()
	if kind == cloudresourcekind.CloudResourceKind_unspecified && fd.Options() != nil {
		kind, _ = proto.GetExtension(fd.Options(), foreignkeyv1.E_DefaultKind).(cloudresourcekind.CloudResourceKind)
	}
	return kind
}
