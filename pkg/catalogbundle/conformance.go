package catalogbundle

import (
	"fmt"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/plantonhq/planton/pkg/crkreflect"
)

// CheckConformance proves the bundle serves EXACTLY what the compiled-in
// registry serves, for every registered kind:
//
//  1. the kind's api message resolves in the bundle,
//  2. the kind's Spec exposes the identical top-level JSON field set,
//  3. the served version recorded in the bundle's kind registry matches,
//  4. the catalog entries name exactly the registry's user-facing kinds
//     (both directions), with unique slugs and at least one official IaC
//     module directory each.
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
