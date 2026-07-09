package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	tt "github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/plantonhq/planton/e2e/framework/provider"
	"github.com/plantonhq/planton/pkg/crkreflect"
	"github.com/plantonhq/planton/pkg/iac/importmap"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// ImportRoundTripEnvVar opts the import round-trip phase in. It stays opt-in
// because it roughly doubles a component's E2E wall time (every resource is
// re-imported and re-planned) -- the scheduled matrix enables it per lane.
const ImportRoundTripEnvVar = "PLANTON_E2E_IMPORT_ROUNDTRIP"

// importRoundTripEnabled gates the phase: opted in, terraform engine (the
// pulumi arm rides the same recipes once its lane lands), and the component
// actually ships an import map.
func importRoundTripEnabled(tc *provider.ComponentTestContext) bool {
	return os.Getenv(ImportRoundTripEnvVar) == "1" &&
		tc.Engine == "terraform" &&
		importmap.HasComponentImportMap(tc.RepoRoot, tc.Provider, tc.Component)
}

// runImportRoundTrip is the machine proof that a component's import recipes
// are CORRECT, not just well-formed: with the deployed fixture's state set
// aside, every resource is re-imported "blind" -- addresses from the real
// state, IDs derived purely through the recipes (spec/metadata/outputs, no
// human) -- and a plan against the re-imported state must show zero changes.
// A recipe that imports the wrong resource, or misses one, cannot pass.
//
// The freshly-imported state replaces the deployed one for the DESTROY phase
// that follows, which is itself part of the proof: destroy tearing the
// fixture down cleanly shows the re-imported state fully owns the resources.
func runImportRoundTrip(tc *provider.ComponentTestContext) error {
	opts, ok := tc.TerraformOpts.(*tt.Options)
	if !ok || opts == nil {
		return errors.New("terraform options not initialized (runValidate must run first)")
	}
	workDir := tc.TerraformWorkDir

	// 1. Addresses come from the REAL deployed state -- never authored.
	stateList, err := tt.RunTerraformCommandE(tc.T, opts, "state", "list")
	if err != nil {
		return errors.Wrap(err, "listing deployed state")
	}
	var addresses []string
	for _, line := range strings.Split(stateList, "\n") {
		if line = strings.TrimSpace(line); line != "" && !strings.HasPrefix(line, "data.") {
			addresses = append(addresses, line)
		}
	}
	if len(addresses) == 0 {
		return errors.New("deployed state lists no resources")
	}

	// 2. Set the deployed state aside -- the re-import starts blind.
	statePath := filepath.Join(workDir, "terraform.tfstate")
	asidePath := statePath + ".pre-roundtrip"
	if err := os.Rename(statePath, asidePath); err != nil {
		return errors.Wrap(err, "setting deployed state aside")
	}
	// A stale backup would let tofu offer state recovery paths that mask a
	// broken import; the round-trip must stand on imports alone.
	_ = os.Remove(statePath + ".backup")

	// 3. Resolve every import ID through the recipes -- no human values.
	catalog, err := importmap.LoadProviderCatalog(tc.RepoRoot, tc.Provider)
	if err != nil {
		return err
	}
	componentMap, err := importmap.LoadComponentImportMap(tc.RepoRoot, tc.Provider, tc.Component)
	if err != nil {
		return err
	}
	formats := map[string]string{}
	configOnlyAttributes := map[string]map[string]bool{}
	for _, rt := range catalog.GetSpec().GetResourceTypes() {
		formats[rt.GetTerraformType()] = rt.GetIdFormat()
		if len(rt.GetConfigOnlyAttributes()) > 0 {
			allowed := make(map[string]bool, len(rt.GetConfigOnlyAttributes()))
			for _, attr := range rt.GetConfigOnlyAttributes() {
				allowed[attr] = true
			}
			configOnlyAttributes[rt.GetTerraformType()] = allowed
		}
	}

	metadataName, spec, err := loadManifestMetadataAndSpec(tc.Component, tc.ManifestPath)
	if err != nil {
		return err
	}

	for _, address := range addresses {
		resourceType, _, instanceKey, parsed := importmap.ParseTofuAddress(address)
		if !parsed {
			return errors.Errorf("address %q is not mappable (module-nested?)", address)
		}
		idFormat, mapped := formats[resourceType]
		if !mapped {
			return errors.Errorf("no id_format for %s (address %s) -- the conformance guard should have caught this", resourceType, address)
		}
		resolved, unresolved := importmap.ResolveValues(componentMap, importmap.Placeholders(idFormat), importmap.ResolveContext{
			MetadataName: metadataName,
			Spec:         spec,
			StackOutputs: tc.FlatOutputs,
			AddressKey:   instanceKey,
		})
		if len(unresolved) > 0 {
			return errors.Errorf("address %s: values %v not derivable from spec/metadata/outputs -- blind import impossible", address, unresolved)
		}
		importID, err := importmap.RenderID(idFormat, resolved)
		if err != nil {
			return errors.Wrapf(err, "address %s", address)
		}

		fmt.Printf("  [import-rt] tofu import %s %s\n", address, importID)
		args := []string{"import", "-input=false"}
		for _, varFile := range opts.VarFiles {
			args = append(args, "-var-file="+varFile)
		}
		args = append(args, address, importID)
		if _, err := tt.RunTerraformCommandE(tc.T, opts, args...); err != nil {
			return errors.Wrapf(err, "importing %s as %q", address, importID)
		}
	}

	// 4. The truth oracle: a plan over the re-imported state must propose no
	// real change. "Real" is precise, not exit-code-blunt: a bare
	// -detailed-exitcode would fail every kind whose config sets a
	// CONFIG-ONLY attribute (aws_s3_bucket.force_destroy is engine delete
	// behavior with no cloud-side existence -- no import can ever carry it,
	// and it legitimately plans as an in-place update after any import). So
	// the plan JSON is inspected structurally: creates, destroys, and
	// replaces always fail; in-place updates pass only when every changed
	// attribute is declared config-only in the provider catalog.
	// Shallow-copy the options: the show-with-struct helper requires a plan
	// file path, and the shared opts must stay pristine for the DESTROY phase.
	planOpts := *opts
	planOpts.PlanFilePath = filepath.Join(workDir, "roundtrip.tfplan")
	planStruct, err := tt.InitAndPlanAndShowWithStructE(tc.T, &planOpts)
	if err != nil {
		return errors.Wrap(err, "plan after re-import")
	}
	tolerated := 0
	for address, rc := range planStruct.ResourceChangesMap {
		if rc.Change == nil || rc.Change.Actions.NoOp() || rc.Change.Actions.Read() || rc.Mode == "data" {
			continue
		}
		if !rc.Change.Actions.Update() {
			return errors.Errorf(
				"plan after blind re-import proposes %v on %s -- the recipes imported the wrong resource or missed one; deployed state kept at %s",
				rc.Change.Actions, address, asidePath)
		}
		changed := changedTopLevelAttributes(rc.Change.Before, rc.Change.After)
		for _, attribute := range changed {
			if !configOnlyAttributes[rc.Type][attribute] {
				return errors.Errorf(
					"plan after blind re-import updates %s.%s (changed: %v) -- not a declared config-only attribute of %s; deployed state kept at %s",
					address, attribute, changed, rc.Type, asidePath)
			}
		}
		fmt.Printf("  [import-rt] tolerating config-only update on %s: %v (declared in the %s catalog)\n",
			address, changed, tc.Provider)
		tolerated++
	}

	fmt.Printf("  [import-rt] %d resources re-imported blind; plan proposes no real change (%d config-only updates tolerated)\n",
		len(addresses), tolerated)
	return nil
}

// changedTopLevelAttributes diffs a plan change's before/after objects and
// returns the top-level attribute names whose values differ. Both sides are
// JSON objects for managed resources; anything non-object degrades to a
// sentinel so the caller fails closed rather than tolerating blindly.
func changedTopLevelAttributes(before, after interface{}) []string {
	beforeMap, beforeOK := before.(map[string]interface{})
	afterMap, afterOK := after.(map[string]interface{})
	if !beforeOK || !afterOK {
		return []string{"<non-object change>"}
	}
	keys := map[string]bool{}
	for k := range beforeMap {
		keys[k] = true
	}
	for k := range afterMap {
		keys[k] = true
	}
	var changed []string
	for k := range keys {
		if !reflect.DeepEqual(beforeMap[k], afterMap[k]) {
			changed = append(changed, k)
		}
	}
	sort.Strings(changed)
	return changed
}

// loadManifestMetadataAndSpec loads the scenario manifest into the kind's api
// message and returns the metadata.name plus the spec message the recipe
// derivations read.
func loadManifestMetadataAndSpec(component, manifestPath string) (string, proto.Message, error) {
	kind := crkreflect.KindFromString(component)
	apiMessage, err := crkreflect.NewInstance(kind)
	if err != nil {
		return "", nil, errors.Wrapf(err, "no proto instance for component %s", component)
	}
	if err := protobufyaml.Load(manifestPath, apiMessage); err != nil {
		return "", nil, errors.Wrapf(err, "loading manifest %s", manifestPath)
	}

	reflected := apiMessage.ProtoReflect()
	metadataField := reflected.Descriptor().Fields().ByName("metadata")
	specField := reflected.Descriptor().Fields().ByName("spec")
	if metadataField == nil || specField == nil {
		return "", nil, errors.Errorf("%s is not a KRM envelope", reflected.Descriptor().FullName())
	}
	nameField := metadataField.Message().Fields().ByName("name")
	metadataName := reflected.Get(metadataField).Message().Get(nameField).String()
	spec := reflected.Get(specField).Message().Interface()
	return metadataName, spec, nil
}
