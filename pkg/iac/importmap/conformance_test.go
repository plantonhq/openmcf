//go:build !codegen
// +build !codegen

package importmap

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	componentv1 "github.com/plantonhq/planton/apis/dev/planton/iac/componentimportmap/v1"
	"github.com/plantonhq/planton/pkg/crkreflect"
)

// mappedKinds is the allowlist of components whose import recipes exist and
// are guarded, keyed by provider directory. A kind is added here only after
// its recipes have passed the live E2E import round-trip (deploy -> state
// aside -> blind re-import -> zero-diff plan), so the guard can never be red
// for an unmapped module -- mirroring the variables.tf drift guard's
// enrollment discipline.
var mappedKinds = map[string][]string{
	"aws": {
		"awss3bucket",
		"awsvpc",
		"awssecuritygroup",
		"awsecrrepo",
	},
}

var tfResourcePattern = regexp.MustCompile(`(?m)^resource\s+"([a-z0-9_]+)"\s+"`)

// TestImportMapConformance validates every authored import recipe against the
// module sources it maps, offline:
//
//  1. The provider catalog and each component map parse against their proto
//     schemas.
//  2. Every resource type the component's OpenTofu module declares has an
//     id_format in the provider catalog -- an unmapped type would surface as
//     an unimportable enumerated address.
//  3. Every {placeholder} those formats reference is declared by the
//     component map, and every declared value either derives from somewhere
//     or tells the user where to find it.
//  4. from_spec_field paths resolve to scalar leaves on the kind's spec
//     proto; from_stack_output keys are real fields on the kind's
//     StackOutputs proto. A typo'd path would silently downgrade a zero-input
//     import to a manual one.
func TestImportMapConformance(t *testing.T) {
	root := repoRoot(t)

	for provider, components := range mappedKinds {
		catalog, err := LoadProviderCatalog(root, provider)
		if err != nil {
			t.Fatalf("provider catalog for %s: %v", provider, err)
		}
		if catalog.GetKind() != "ProviderImportCatalog" {
			t.Fatalf("provider catalog for %s: kind is %q, want ProviderImportCatalog", provider, catalog.GetKind())
		}

		formatsByTerraformType := map[string]string{}
		for _, rt := range catalog.GetSpec().GetResourceTypes() {
			if rt.GetTerraformType() == "" {
				t.Errorf("%s catalog: entry with empty terraform_type", provider)
				continue
			}
			if _, dup := formatsByTerraformType[rt.GetTerraformType()]; dup {
				t.Errorf("%s catalog: duplicate entry for %s", provider, rt.GetTerraformType())
			}
			formatsByTerraformType[rt.GetTerraformType()] = rt.GetIdFormat()
			if rt.GetIdFormat() == "" {
				t.Errorf("%s catalog: %s has an empty id_format", provider, rt.GetTerraformType())
			}
			// Config-only attributes are the round-trip proof's tolerance list;
			// an empty entry would tolerate nothing meaningful and signals a
			// YAML authoring slip.
			for _, attr := range rt.GetConfigOnlyAttributes() {
				if strings.TrimSpace(attr) == "" {
					t.Errorf("%s catalog: %s declares an empty config_only_attributes entry", provider, rt.GetTerraformType())
				}
			}
		}

		for _, component := range components {
			component := component
			t.Run(provider+"/"+component, func(t *testing.T) {
				m, err := LoadComponentImportMap(root, provider, component)
				if err != nil {
					t.Fatalf("component import map: %v", err)
				}
				if m.GetMetadata().GetName() != component {
					t.Errorf("metadata.name is %q, want %q", m.GetMetadata().GetName(), component)
				}

				declaredValues := map[string]*componentv1.ImportValue{}
				for _, v := range m.GetSpec().GetValues() {
					if v.GetName() == "" {
						t.Error("import value with empty name")
						continue
					}
					declaredValues[v.GetName()] = v
					// A value that cannot derive from anywhere MUST carry guidance --
					// the empathy contract: when we must ask, show where to find it.
					if len(v.GetDerivations()) == 0 && v.GetWhereToFind() == "" {
						t.Errorf("value %q has no derivations and no where_to_find", v.GetName())
					}
				}

				// Every module-declared resource type must be importable.
				moduleTypes := terraformResourceTypes(t, filepath.Join(root, "apis/dev/planton/provider", provider, component, "v1/iac/tf"))
				if len(moduleTypes) == 0 {
					t.Fatal("module declares no resources -- wrong path?")
				}
				for _, resourceType := range moduleTypes {
					idFormat, mapped := formatsByTerraformType[resourceType]
					if !mapped {
						t.Errorf("module resource %s has no id_format in the %s catalog", resourceType, provider)
						continue
					}
					for _, placeholder := range Placeholders(idFormat) {
						if _, declared := declaredValues[placeholder]; !declared {
							t.Errorf("placeholder {%s} (id_format of %s) not declared in the component map", placeholder, resourceType)
						}
					}
				}

				// Derivation paths must resolve against the kind's protos.
				kind := crkreflect.KindFromString(component)
				apiMessage, err := crkreflect.NewInstance(kind)
				if err != nil {
					t.Fatalf("NewInstance(%s): %v", component, err)
				}
				specDescriptor := specMessageDescriptor(t, apiMessage)
				outputsDescriptor := stackOutputsDescriptor(t, apiMessage)
				for _, v := range m.GetSpec().GetValues() {
					for _, d := range v.GetDerivations() {
						switch source := d.GetSource().(type) {
						case *componentv1.ImportValueDerivation_FromSpecField:
							if err := validateScalarPath(specDescriptor, source.FromSpecField); err != nil {
								t.Errorf("value %q from_spec_field %q: %v", v.GetName(), source.FromSpecField, err)
							}
						case *componentv1.ImportValueDerivation_FromStackOutput:
							if outputsDescriptor.Fields().ByName(protoreflect.Name(source.FromStackOutput)) == nil {
								t.Errorf("value %q from_stack_output %q: no such field on %s",
									v.GetName(), source.FromStackOutput, outputsDescriptor.FullName())
							}
						case *componentv1.ImportValueDerivation_FromArnPart:
							switch source.FromArnPart {
							case "resource_id", "resource_name", "account_id", "region":
							default:
								t.Errorf("value %q from_arn_part %q: not a recognized ARN part", v.GetName(), source.FromArnPart)
							}
						}
					}
				}
			})
		}
	}
}

// terraformResourceTypes scans a module directory's .tf files for top-level
// resource declarations and returns the distinct resource types.
func terraformResourceTypes(t *testing.T, moduleDir string) []string {
	t.Helper()
	entries, err := os.ReadDir(moduleDir)
	if err != nil {
		t.Fatalf("reading module dir %s: %v", moduleDir, err)
	}
	seen := map[string]bool{}
	var types []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tf") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(moduleDir, entry.Name()))
		if err != nil {
			t.Fatalf("reading %s: %v", entry.Name(), err)
		}
		for _, match := range tfResourcePattern.FindAllStringSubmatch(string(content), -1) {
			if !seen[match[1]] {
				seen[match[1]] = true
				types = append(types, match[1])
			}
		}
	}
	return types
}

// specMessageDescriptor returns the descriptor of the kind's spec message
// (the api envelope's `spec` field).
func specMessageDescriptor(t *testing.T, apiMessage proto.Message) protoreflect.MessageDescriptor {
	t.Helper()
	specField := apiMessage.ProtoReflect().Descriptor().Fields().ByName("spec")
	if specField == nil || specField.Kind() != protoreflect.MessageKind {
		t.Fatalf("%s has no spec message field", apiMessage.ProtoReflect().Descriptor().FullName())
	}
	return specField.Message()
}

// stackOutputsDescriptor locates the kind's StackOutputs message from the
// sibling stack_outputs.proto of the kind's api.proto.
func stackOutputsDescriptor(t *testing.T, apiMessage proto.Message) protoreflect.MessageDescriptor {
	t.Helper()
	apiPath := apiMessage.ProtoReflect().Descriptor().ParentFile().Path()
	outputsPath := filepath.Join(filepath.Dir(apiPath), "stack_outputs.proto")
	file, err := protoregistry.GlobalFiles.FindFileByPath(outputsPath)
	if err != nil {
		t.Fatalf("finding %s: %v", outputsPath, err)
	}
	for i := 0; i < file.Messages().Len(); i++ {
		msg := file.Messages().Get(i)
		if strings.HasSuffix(string(msg.Name()), "StackOutputs") {
			return msg
		}
	}
	t.Fatalf("%s declares no *StackOutputs message", outputsPath)
	return nil
}

// validateScalarPath walks a dot path of field names and requires a scalar
// leaf with no repeated/map segments (a repeated segment would make the
// derivation ambiguous -- which instance?).
func validateScalarPath(desc protoreflect.MessageDescriptor, dotPath string) error {
	if dotPath == "" {
		return fmt.Errorf("empty path")
	}
	current := desc
	segments := strings.Split(dotPath, ".")
	for i, segment := range segments {
		field := current.Fields().ByName(protoreflect.Name(segment))
		if field == nil {
			return fmt.Errorf("no field %q on %s", segment, current.FullName())
		}
		if field.IsList() || field.IsMap() {
			return fmt.Errorf("segment %q is repeated/map -- ambiguous for derivation", segment)
		}
		if i == len(segments)-1 {
			if field.Kind() == protoreflect.MessageKind || field.Kind() == protoreflect.GroupKind {
				return fmt.Errorf("leaf %q is a message, not a scalar", segment)
			}
			return nil
		}
		if field.Kind() != protoreflect.MessageKind {
			return fmt.Errorf("segment %q is a scalar but the path continues", segment)
		}
		current = field.Message()
	}
	return nil
}

// repoRoot walks up from this test file to the directory containing go.mod.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found walking up from test file")
		}
		dir = parent
	}
}
