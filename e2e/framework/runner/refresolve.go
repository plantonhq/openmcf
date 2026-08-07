package runner

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/manifest"
	"github.com/plantonhq/planton/pkg/outputs"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"sigs.k8s.io/yaml"
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

// resolveRefsInMessage replaces value_from arms on the message's StringValueOrRef
// fields with literals from the matching prerequisite's outputs, and recurses
// into nested (repeated) messages so refs buried inside structured spec fields
// resolve too. Returns whether any field was resolved.
func resolveRefsInMessage(msg protoreflect.Message, flattened map[cloudresourcekind.CloudResourceKind]map[string]map[string]string) (bool, error) {
	resolvedAny := false
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
					resolved, err := resolveRefsInMessage(list.Get(j).Message(), flattened)
					if err != nil {
						return false, err
					}
					resolvedAny = resolvedAny || resolved
					continue
				}
				// Repeated refs (e.g. a role's managed_policy_arns): each element
				// resolves independently, so literals and references can mix in
				// one list.
				ref, ok := list.Get(j).Message().Interface().(*foreignkeyv1.StringValueOrRef)
				if !ok || ref.GetValueFrom() == nil {
					continue
				}
				val, resolved, err := lookupRefValue(fd, ref, flattened)
				if err != nil {
					return false, err
				}
				if !resolved {
					continue // no prerequisite matches this element; leave it untouched
				}
				list.Set(j, protoreflect.ValueOfMessage((&foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
				}).ProtoReflect()))
				resolvedAny = true
			}
			continue
		}

		if !msg.Has(fd) {
			continue
		}

		if !isRef {
			// A singular nested message (e.g. a listener's tls config): recurse.
			resolved, err := resolveRefsInMessage(msg.Mutable(fd).Message(), flattened)
			if err != nil {
				return false, err
			}
			resolvedAny = resolvedAny || resolved
			continue
		}

		ref, ok := msg.Get(fd).Message().Interface().(*foreignkeyv1.StringValueOrRef)
		if !ok || ref.GetValueFrom() == nil {
			continue
		}
		val, resolved, err := lookupRefValue(fd, ref, flattened)
		if err != nil {
			return false, err
		}
		if !resolved {
			continue
		}
		msg.Set(fd, protoreflect.ValueOfMessage((&foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
		}).ProtoReflect()))
		resolvedAny = true
	}
	return resolvedAny, nil
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

	kind := valueFrom.GetKind()
	if kind == cloudresourcekind.CloudResourceKind_unspecified && fd.Options() != nil {
		kind, _ = proto.GetExtension(fd.Options(), foreignkeyv1.E_DefaultKind).(cloudresourcekind.CloudResourceKind)
	}
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
