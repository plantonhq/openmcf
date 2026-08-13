package catalogbundle

import (
	"fmt"
	"sort"
	"strings"
	"testing/fstest"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/plantonhq/planton/pkg/conversion"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
)

// CheckConformance proves the bundle serves EXACTLY what the compiled-in
// registry serves, for every registered kind:
//
//  1. the kind's api message resolves in the bundle,
//  2. the kind's Spec exposes the identical top-level JSON field set,
//  3. the served version recorded in the bundle's kind registry matches,
//  4. the catalog entries name exactly the registry's user-facing kinds
//     (both directions), with unique slugs and at least one official IaC
//     module directory each,
//  5. the version deprecations agree with the compiled registry (both
//     directions), and every deprecation the bundle announces names a
//     version whose schema the bundle carries AND has an authored
//     conversion path to the served version -- a deprecation without a
//     way out is a dead end this gate refuses to release.
//
// This is the gate between "we built a zip" and "a runtime may trust this
// zip instead of its compiled classes". It runs in CI against the buf-built
// bundle on every proto change, and in the release lane before upload.
func CheckConformance(bundle *Bundle) error {
	var problems []string
	checked := 0

	for _, kind := range crkreflect.KindsList() {
		compiled := crkreflect.ToMessageMap[kind]
		if compiled == nil {
			continue
		}
		fullName := proto.MessageName(compiled)

		bundleDesc, err := bundle.Files.FindDescriptorByName(fullName)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: the bundle does not carry %s", kind, fullName))
			continue
		}
		bundleMsg, ok := bundleDesc.(protoreflect.MessageDescriptor)
		if !ok {
			problems = append(problems, fmt.Sprintf("%s: %s is not a message in the bundle", kind, fullName))
			continue
		}

		compiledSpec := specFieldSet(compiled.ProtoReflect().Descriptor())
		bundledSpec := specFieldSet(bundleMsg)
		if compiledSpec != bundledSpec {
			problems = append(problems, fmt.Sprintf(
				"%s: spec field sets differ between the compiled registry and the bundle\n    compiled: %s\n    bundle:   %s",
				kind, compiledSpec, bundledSpec))
		}

		// The served version must agree between the two registries.
		compiledVersion, err := crkreflect.KindVersion(kind)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", kind, err))
			continue
		}
		bundleVersion := versionSegmentOf(fullName)
		if bundleVersion != compiledVersion {
			problems = append(problems, fmt.Sprintf(
				"%s: the registry serves %s but the bundle's message package says %s",
				kind, compiledVersion, bundleVersion))
		}

		checked++
	}

	problems = append(problems, entryProblems(bundle)...)
	problems = append(problems, deprecationProblems(bundle)...)

	if checked == 0 {
		return fmt.Errorf("conformance checked zero kinds -- the walk is broken")
	}
	if len(problems) > 0 {
		sort.Strings(problems)
		return fmt.Errorf("the catalog bundle does not conform to the compiled registry (%d finding(s)):\n  %s",
			len(problems), strings.Join(problems, "\n  "))
	}
	return nil
}

// entryProblems checks the bundle's catalog entries against the compiled
// registry: every user-facing kind has its entry, every entry names a
// registered kind, slugs are unique, and every entry carries at least one
// official module directory (a component with no deployable module has no
// business in the catalog).
func entryProblems(bundle *Bundle) []string {
	var problems []string

	entryByKind := map[string]CatalogEntry{}
	slugOwners := map[string]string{}
	for _, entry := range bundle.CatalogEntries() {
		entryByKind[entry.Kind] = entry
		if owner, taken := slugOwners[entry.Slug]; taken {
			problems = append(problems, fmt.Sprintf(
				"catalog entries %s and %s share the slug %q -- slugs must be unique", owner, entry.Kind, entry.Slug))
			continue
		}
		slugOwners[entry.Slug] = entry.Kind
	}

	userFacing := map[string]bool{}
	for _, kind := range crkreflect.KindsList() {
		if crkreflect.GetProvider(kind).String() == testProviderName {
			continue
		}
		kindName := crkreflect.ExtractKindNameByKind(kind)
		userFacing[kindName] = true
		entry, ok := entryByKind[kindName]
		if !ok {
			problems = append(problems, fmt.Sprintf(
				"%s: the registry serves this kind but the bundle carries no catalog entry for it", kindName))
			continue
		}
		if entry.IacModules.TerraformModuleDir == "" && entry.IacModules.PulumiModuleDir == "" {
			problems = append(problems, fmt.Sprintf(
				"%s: its catalog entry names no official IaC module directory for any engine", kindName))
		}
	}
	for kindName := range entryByKind {
		if !userFacing[kindName] {
			problems = append(problems, fmt.Sprintf(
				"%s: the bundle carries a catalog entry but the registry does not serve this kind", kindName))
		}
	}
	return problems
}

// deprecationProblems checks the deprecation half of the kind registry.
// Runtimes announce deprecations from the BUNDLE's registry, so the checks
// run against what the bundle declares, with agreement against the compiled
// registry proving the bundle is not stale (a bundle built before a
// deprecation was authored would otherwise pass every schema check and
// silently announce nothing). The compile-time facts -- grammar, duplicates,
// never the served version -- are gated by the crkreflect registry tests;
// this gate owns the facts only the built artifact can prove: the deprecated
// version's schema is aboard, and an authored conversion path to the served
// version exists among the bundle's own specs.
func deprecationProblems(bundle *Bundle) []string {
	var problems []string

	bundleDeprecations, err := bundleKindDeprecations(bundle)
	if err != nil {
		return []string{err.Error()}
	}
	specFS := conversionSpecFS(bundle)

	for _, kind := range crkreflect.KindsList() {
		compiled := crkreflect.ToMessageMap[kind]
		if compiled == nil {
			continue
		}

		compiledDeps, err := crkreflect.KindDeprecations(kind)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", kind, err))
			continue
		}
		bundleDeps := bundleDeprecations[kind.String()]

		// Agreement both directions, note included -- the note is announced
		// verbatim on user surfaces, so a stale note misinforms just like a
		// missing deprecation.
		compiledSet := renderDeprecations(compiledDeps)
		bundleSet := renderDeprecations(bundleDeps)
		if compiledSet != bundleSet {
			problems = append(problems, fmt.Sprintf(
				"%s: version deprecations differ between the compiled registry and the bundle's kind registry -- the bundle is stale or tampered\n    compiled: %s\n    bundle:   %s",
				kind, orNone(compiledSet), orNone(bundleSet)))
		}

		if len(bundleDeps) == 0 {
			continue
		}
		servedVersion, err := crkreflect.KindVersion(kind)
		if err != nil {
			// Already reported by the served-version check in the main walk.
			continue
		}
		for _, dep := range bundleDeps {
			deprecatedName := versionSiblingName(proto.MessageName(compiled), dep.GetVersion())
			if _, err := bundle.Files.FindDescriptorByName(deprecatedName); err != nil {
				problems = append(problems, fmt.Sprintf(
					"%s: version %s is announced deprecated but the bundle carries no schema for it (%s) -- a deprecation must name a version this release ships",
					kind, dep.GetVersion(), deprecatedName))
				continue
			}
			specs, err := conversion.SpecsForKind(specFS, kind)
			if err != nil {
				problems = append(problems, fmt.Sprintf("%s: reading the bundle's conversion specs: %v", kind, err))
				continue
			}
			if _, err := conversion.Path(specs, dep.GetVersion(), servedVersion); err != nil {
				problems = append(problems, fmt.Sprintf(
					"%s: version %s is announced deprecated but has no conversion path to the served version %s -- a deprecation without a way out strands its writers: %v",
					kind, dep.GetVersion(), servedVersion, err))
			}
		}
	}
	return problems
}

// kindRegistryEnumName locates the kind registry inside a bundle's
// descriptor set -- the same enum the compiled registry is generated from.
const kindRegistryEnumName = "dev.planton.shared.cloudresourcekind.CloudResourceKind"

// bundleKindDeprecations reads each kind's declared deprecations off the
// BUNDLE's kind registry enum options, keyed by kind name. The extension
// type is compiled shared-registry infrastructure (never per-kind code), so
// bundle bytes populate it directly.
func bundleKindDeprecations(bundle *Bundle) (map[string][]*cloudresourcekind.CloudResourceKindVersionDeprecation, error) {
	desc, err := bundle.Files.FindDescriptorByName(kindRegistryEnumName)
	if err != nil {
		return nil, fmt.Errorf("the bundle carries no kind registry enum: %w", err)
	}
	enum, ok := desc.(protoreflect.EnumDescriptor)
	if !ok {
		return nil, fmt.Errorf("the bundle's kind registry is not an enum")
	}
	out := map[string][]*cloudresourcekind.CloudResourceKindVersionDeprecation{}
	values := enum.Values()
	for i := 0; i < values.Len(); i++ {
		value := values.Get(i)
		opts := value.Options()
		if opts == nil {
			continue
		}
		meta, ok := proto.GetExtension(opts, cloudresourcekind.E_KindMeta).(*cloudresourcekind.CloudResourceKindMeta)
		if !ok || meta == nil || len(meta.GetDeprecations()) == 0 {
			continue
		}
		out[string(value.Name())] = meta.GetDeprecations()
	}
	return out, nil
}

// conversionSpecFS re-roots the bundle's conversion specs
// (conversions/<provider>/<kind>/<file>) into the on-disk layout the
// conversion package discovers (<provider>/<kind>/conversions/<file>), so
// spec discovery and path finding run the EXACT code the engines use --
// fstest.MapFS is the stdlib's canonical in-memory fs.FS, test-named but
// production-clean.
func conversionSpecFS(bundle *Bundle) fstest.MapFS {
	fsys := fstest.MapFS{}
	for name, content := range bundle.ConversionSpecs() {
		parts := strings.Split(name, "/")
		if len(parts) != 4 {
			continue
		}
		fsys[parts[1]+"/"+parts[2]+"/conversions/"+parts[3]] = &fstest.MapFile{Data: content}
	}
	return fsys
}

// renderDeprecations canonicalizes a deprecation list for comparison and
// display: sorted "version(note)" entries.
func renderDeprecations(deps []*cloudresourcekind.CloudResourceKindVersionDeprecation) string {
	rendered := make([]string, 0, len(deps))
	for _, dep := range deps {
		rendered = append(rendered, fmt.Sprintf("%s(%q)", dep.GetVersion(), dep.GetNote()))
	}
	sort.Strings(rendered)
	return strings.Join(rendered, ", ")
}

func orNone(rendered string) string {
	if rendered == "" {
		return "(none)"
	}
	return rendered
}

// versionSiblingName swaps the version segment of a kind message's full name
// (dev.planton.<p>.<k>.<version>.<Kind>) -- how a non-served version's schema
// is located in the bundle.
func versionSiblingName(fullName protoreflect.FullName, version string) protoreflect.FullName {
	parts := strings.Split(string(fullName), ".")
	if len(parts) < 2 {
		return fullName
	}
	parts[len(parts)-2] = version
	return protoreflect.FullName(strings.Join(parts, "."))
}

// specFieldSet renders a message's `spec` field's top-level JSON field names
// as a stable comparable string; kinds without a spec field compare as "".
func specFieldSet(desc protoreflect.MessageDescriptor) string {
	specField := desc.Fields().ByName("spec")
	if specField == nil || specField.Kind() != protoreflect.MessageKind {
		return ""
	}
	fields := specField.Message().Fields()
	names := make([]string, 0, fields.Len())
	for i := 0; i < fields.Len(); i++ {
		names = append(names, fields.Get(i).JSONName())
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

// versionSegmentOf extracts the version segment from a kind message's full
// name (dev.planton.<p>.<k>.<version>.<Kind>).
func versionSegmentOf(fullName protoreflect.FullName) string {
	parts := strings.Split(string(fullName), ".")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-2]
}
