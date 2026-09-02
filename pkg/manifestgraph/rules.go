package manifestgraph

import (
	"fmt"
	"strings"

	"github.com/plantonhq/planton/pkg/refcheck"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

// Target is a checked reference's resolved destination: which kind and name
// it points at, the env it names (empty means "the consumer's own env"), and
// the EFFECTIVE field path after annotation defaults applied.
type Target struct {
	Kind      cloudresourcekind.CloudResourceKind
	Name      string
	Env       string
	FieldPath string
}

// Identity derives the target's graph identity, falling back to the
// consumer's env when the reference names none — the same fallback the
// platform's edge selectors apply. The slug derives through the ONE slug
// function (see the phantom-node warning in the package doc).
func (t Target) Identity(consumerEnv string) Identity {
	env := t.Env
	if env == "" {
		env = consumerEnv
	}
	return Identity{Kind: t.Kind, Slug: GenerateSlug(t.Name), Env: env}
}

// EffectiveKind resolves which kind a reference points at: the reference's
// explicit kind wins, else the field's default_kind annotation, else
// unspecified. This is the kind-half of CheckRef, exported separately because
// resolution lookups need it without the full rule evaluation.
func EffectiveKind(use RefUse) cloudresourcekind.CloudResourceKind {
	kind := use.Ref.GetKind()
	if kind == cloudresourcekind.CloudResourceKind_unspecified && use.Field.Options() != nil {
		kind, _ = proto.GetExtension(use.Field.Options(), foreignkeyv1.E_DefaultKind).(cloudresourcekind.CloudResourceKind)
	}
	return kind
}

// CheckRef validates one valueFrom reference against the foreign-key
// annotations on its declaring field and the referenced kind's proto surface.
// It returns the resolved target (for dependency-graph construction and
// resolution) and any problems found.
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
//     otherwise only surfaces at deploy time. A path that EXTENDS the
//     annotated path is not an override: map-typed composition keys are
//     addressed by entry key (`status.outputs.backend_pool_ids.web`), and
//     the entry key is data the annotation cannot name.
//  4. The effective field path must resolve against the target kind's actual
//     proto surface (stack outputs, spec, or metadata).
func CheckRef(use RefUse) (Target, []string) {
	var problems []string

	var annotatedKind cloudresourcekind.CloudResourceKind
	var annotatedPath string
	if opts := use.Field.Options(); opts != nil {
		annotatedKind, _ = proto.GetExtension(opts, foreignkeyv1.E_DefaultKind).(cloudresourcekind.CloudResourceKind)
		annotatedPath, _ = proto.GetExtension(opts, foreignkeyv1.E_DefaultKindFieldPath).(string)
	}

	targetKind := use.Ref.GetKind()
	if targetKind == cloudresourcekind.CloudResourceKind_unspecified {
		targetKind = annotatedKind
	}
	if targetKind == cloudresourcekind.CloudResourceKind_unspecified {
		problems = append(problems,
			fmt.Sprintf("%s: valueFrom does not name a kind and the field declares no default kind — add an explicit `kind:`", use.FieldPath))
		return Target{}, problems
	}

	effectivePath := use.Ref.GetFieldPath()
	if effectivePath == "" {
		if targetKind == annotatedKind && annotatedPath != "" {
			effectivePath = annotatedPath
		} else {
			problems = append(problems,
				fmt.Sprintf("%s: valueFrom targets %s but has no fieldPath, and no annotated default applies — add an explicit `fieldPath:`", use.FieldPath, targetKind))
			return Target{Kind: targetKind, Name: use.Ref.GetName(), Env: use.Ref.GetEnv()}, problems
		}
	} else if targetKind == annotatedKind && annotatedPath != "" && effectivePath != annotatedPath &&
		!strings.HasPrefix(effectivePath, annotatedPath+".") {
		problems = append(problems,
			fmt.Sprintf("%s: valueFrom overrides the annotated composition key for %s — the field's contract is %q but the reference names %q (id/name/self-link format mismatches only surface at deploy time; use the annotated path)",
				use.FieldPath, targetKind, annotatedPath, effectivePath))
	}

	if reason := refcheck.ResolveValueFromPath(targetKind, effectivePath); reason != "" {
		problems = append(problems,
			fmt.Sprintf("%s: valueFrom fieldPath %q does not resolve on %s: %s", use.FieldPath, effectivePath, targetKind, reason))
	}

	return Target{Kind: targetKind, Name: use.Ref.GetName(), Env: use.Ref.GetEnv(), FieldPath: effectivePath}, problems
}
