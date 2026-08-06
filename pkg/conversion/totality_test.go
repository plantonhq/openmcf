package conversion

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	validatepb "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/crkreflect"
)

// The no-new-version-without-a-bridge gate.
//
// A kind directory holding MORE than one api-version directory has committed
// to serving a version pair -- from that moment a conversion path between
// them must exist, and the upgrade must be TOTAL for the top-level spec
// fields: every field the old version had is either present unchanged in the
// new version or explicitly handled by an op (rename source, map source, or
// a drop that declares its loss), and every REQUIRED field the new version
// added is produced by an op. A field that simply vanishes is exactly the
// silent-loss defect this machinery exists to kill.

var versionDirRe = regexp.MustCompile(`^v\d+((alpha|beta)\d+)?$`)

func TestEveryVersionPairHasATotalBridge(t *testing.T) {
	base := providerBaseDir(t)
	fsys := os.DirFS(base)

	providers, err := os.ReadDir(base)
	if err != nil {
		t.Fatal(err)
	}
	checkedPairs := 0
	for _, provider := range providers {
		if !provider.IsDir() {
			continue
		}
		kinds, err := os.ReadDir(filepath.Join(base, provider.Name()))
		if err != nil {
			t.Fatal(err)
		}
		for _, kindDir := range kinds {
			if !kindDir.IsDir() {
				continue
			}
			var versions []string
			entries, err := os.ReadDir(filepath.Join(base, provider.Name(), kindDir.Name()))
			if err != nil {
				t.Fatal(err)
			}
			for _, e := range entries {
				if e.IsDir() && versionDirRe.MatchString(e.Name()) {
					versions = append(versions, e.Name())
				}
			}
			if len(versions) < 2 {
				continue
			}

			kind := crkreflect.KindFromString(kindDir.Name())
			if kind == cloudresourcekind.CloudResourceKind_unspecified {
				t.Errorf("%s/%s has version directories but is not a registered kind", provider.Name(), kindDir.Name())
				continue
			}
			specs, err := SpecsForKind(fsys, kind)
			if err != nil {
				t.Errorf("%s/%s: %v", provider.Name(), kindDir.Name(), err)
				continue
			}

			// Every non-served version must have a conversion path to the
			// served (registry) version.
			served, err := crkreflect.KindVersion(kind)
			if err != nil {
				t.Errorf("%s/%s: %v", provider.Name(), kindDir.Name(), err)
				continue
			}
			for _, version := range versions {
				if version == served {
					continue
				}
				steps, err := Path(specs, version, served)
				if err != nil {
					t.Errorf("%s/%s serves %s but still carries %s: %v", provider.Name(), kindDir.Name(), served, version, err)
					continue
				}
				for _, step := range steps {
					checkTotality(t, provider.Name(), kindDir.Name(), step.Spec)
					checkedPairs++
				}
			}
		}
	}
	if checkedPairs == 0 {
		t.Fatal("no version pairs were checked -- the torture kind's pair must exist; the gate walk is broken")
	}
}

// checkTotality verifies the spec accounts for every top-level spec-field
// difference between the two versions' Spec descriptors.
func checkTotality(t *testing.T, provider, kindDir string, spec *Spec) {
	t.Helper()
	fromFields, err := specFieldJSONNames(spec.Kind, provider, kindDir, spec.From)
	if err != nil {
		t.Errorf("%s/%s: %v", provider, kindDir, err)
		return
	}
	toFields, err := specFieldJSONNames(spec.Kind, provider, kindDir, spec.To)
	if err != nil {
		t.Errorf("%s/%s: %v", provider, kindDir, err)
		return
	}

	handledSources := map[string]bool{}
	producedTargets := map[string]bool{}
	for _, op := range spec.Ops {
		switch {
		case op.Rename != nil:
			handledSources[topSpecField(op.Rename.FromPath)] = true
			producedTargets[topSpecField(op.Rename.ToPath)] = true
		case op.Map != nil:
			handledSources[topSpecField(op.Map.Path)] = true
			to := op.Map.To
			if to == "" {
				to = op.Map.Path
			}
			producedTargets[topSpecField(to)] = true
		case op.Default != nil:
			producedTargets[topSpecField(op.Default.Path)] = true
		case op.Drop != nil:
			handledSources[topSpecField(op.Drop.Path)] = true
		}
	}

	for name := range fromFields {
		if _, stillThere := toFields[name]; stillThere {
			continue
		}
		if !handledSources[name] {
			t.Errorf("%s/%s %s->%s: spec field %q exists in %s but not in %s and NO op handles it -- rename it, map it, or drop it with a declared loss; a field must never just vanish",
				provider, kindDir, spec.From, spec.To, name, spec.From, spec.To)
		}
	}
	for name, required := range toFields {
		if _, existed := fromFields[name]; existed {
			continue
		}
		if required && !producedTargets[name] {
			t.Errorf("%s/%s %s->%s: %s added REQUIRED spec field %q and no op produces it -- upgraded documents would all be invalid; add a default or map op",
				provider, kindDir, spec.From, spec.To, name, spec.To)
		}
	}
}

// topSpecField extracts the top-level spec field a path addresses
// ("spec.displayName" -> "displayName"); paths outside spec return "".
func topSpecField(path string) string {
	segments := strings.SplitN(path, ".", 3)
	if len(segments) < 2 || segments[0] != "spec" {
		return ""
	}
	return segments[1]
}

// specFieldJSONNames returns the JSON names of a version's top-level Spec
// fields mapped to whether each is required (buf.validate required).
func specFieldJSONNames(kindName, provider, kindDir, version string) (map[string]bool, error) {
	fullName := protoreflect.FullName(fmt.Sprintf("dev.planton.provider.%s.%s.%s.%sSpec", provider, kindDir, version, kindName))
	msgType, err := protoregistry.GlobalTypes.FindMessageByName(fullName)
	if err != nil {
		return nil, fmt.Errorf("the %s Spec descriptor (%s) is not linked into this test -- import the version's package in corpus_test.go so totality can be checked: %w", version, fullName, err)
	}
	fields := map[string]bool{}
	descriptor := msgType.Descriptor()
	for i := 0; i < descriptor.Fields().Len(); i++ {
		field := descriptor.Fields().Get(i)
		fields[field.JSONName()] = isRequiredField(field)
	}
	return fields, nil
}

// isRequiredField reports whether the field carries buf.validate required.
func isRequiredField(field protoreflect.FieldDescriptor) bool {
	opts := field.Options()
	if opts == nil {
		return false
	}
	rules, ok := proto.GetExtension(opts, validatepb.E_Field).(*validatepb.FieldRules)
	if !ok || rules == nil {
		return false
	}
	return rules.GetRequired()
}
