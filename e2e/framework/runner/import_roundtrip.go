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
// actually ships an import map. File presence is the recipes' single
// enrollment signal everywhere -- this gate, the offline conformance guard,
// and the platform's catalog bundler all key off the same import-map.yaml,
// so a map cannot ship while dodging its checks.
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
	// Both declared tolerance classes merge into one allowlist for the plan
	// oracle: config-only attributes never exist cloud-side, write-normalized
	// attributes read back in a provider-normalized form -- either way an
	// in-place update touching only them is the documented post-import shape,
	// not a wrong import.
	// A declared tolerance is either a top-level attribute name ("force_destroy")
	// or a dotted sub-path ("spec.update_strategy") for the narrower class where
	// a provider's importer fails to read back ONE nested block while the rest
	// of the attribute must stay under the oracle. Dotted entries are collected
	// separately: the changed attribute is then re-compared with only those
	// sub-paths pruned, so an undeclared sibling drift still fails.
	formats := map[string]string{}
	toleratedAttributes := map[string]map[string]bool{}
	toleratedSubPaths := map[string][][]string{}
	for _, rt := range catalog.GetSpec().GetResourceTypes() {
		formats[rt.GetTerraformType()] = rt.GetIdFormat()
		declared := append(append([]string{}, rt.GetConfigOnlyAttributes()...), rt.GetWriteNormalizedAttributes()...)
		for _, attr := range declared {
			if strings.Contains(attr, ".") {
				toleratedSubPaths[rt.GetTerraformType()] = append(
					toleratedSubPaths[rt.GetTerraformType()], strings.Split(attr, "."))
				continue
			}
			if toleratedAttributes[rt.GetTerraformType()] == nil {
				toleratedAttributes[rt.GetTerraformType()] = map[string]bool{}
			}
			toleratedAttributes[rt.GetTerraformType()][attr] = true
		}
	}

	metadataName, spec, err := loadManifestMetadataAndSpec(tc.Component, tc.ManifestPath)
	if err != nil {
		return err
	}

	// The round-trip has no user-pasted ARN, but the ACCOUNT-LEVEL ARN parts
	// (account_id, region) are properties of the deployed account itself, not
	// of any one resource -- the platform always knows them from the
	// connection, and here every ARN-shaped stack output carries the same
	// pair. Recipes for IDs that embed the account id (e.g. DynamoDB
	// contributor insights) stay blind-derivable. Per-resource parts
	// (resource_id/resource_name/arn) are deliberately NOT filled: those
	// genuinely require the user's ARN, and faking them from an unrelated
	// output would let a wrong recipe pass.
	accountArnParts := accountLevelArnParts(tc.FlatOutputs)

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
			ArnParts:     accountArnParts,
		})
		// Only required placeholders abort the blind import -- optional
		// ("{name?}") segments are provider-documented as legitimately empty
		// for some variants and render as "".
		if missing := requiredUnresolved(idFormat, unresolved); len(missing) > 0 {
			return errors.Errorf("address %s: values %v not derivable from spec/metadata/outputs -- blind import impossible", address, missing)
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
	// and it legitimately plans as an in-place update after any import) or a
	// WRITE-NORMALIZED one (policy documents read back provider-normalized).
	// So the plan JSON is inspected structurally: creates, destroys, and
	// replaces always fail; in-place updates pass only when every changed
	// attribute is declared in the provider catalog's tolerance lists.
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
			if toleratedAttributes[rc.Type][attribute] {
				continue
			}
			if changeCoveredBySubPaths(rc.Change.Before, rc.Change.After, attribute, toleratedSubPaths[rc.Type]) {
				continue
			}
			return errors.Errorf(
				"plan after blind re-import updates %s.%s (changed: %v) -- not a declared config-only/write-normalized attribute of %s; deployed state kept at %s",
				address, attribute, changed, rc.Type, asidePath)
		}
		fmt.Printf("  [import-rt] tolerating declared update on %s: %v (config-only/write-normalized in the %s catalog)\n",
			address, changed, tc.Provider)
		tolerated++
	}

	fmt.Printf("  [import-rt] %d resources re-imported blind; plan proposes no real change (%d config-only updates tolerated)\n",
		len(addresses), tolerated)
	return nil
}

// accountLevelArnParts extracts the account_id and region from the first
// ARN-shaped stack output. Within one deployment every ARN carries the same
// account (and, for regional services, the same region), so these two parts
// are deployment-level facts -- the same facts the platform derives from the
// provider connection in the real import flow.
func accountLevelArnParts(flatOutputs map[string]string) map[string]string {
	parts := map[string]string{}
	for _, value := range flatOutputs {
		if !strings.HasPrefix(value, "arn:") {
			continue
		}
		// arn:partition:service:region:account-id:resource
		segments := strings.SplitN(value, ":", 6)
		if len(segments) < 6 {
			continue
		}
		// Global-service ARNs (S3, IAM) leave region -- and S3 even the
		// account -- empty, so each part fills from the first ARN that
		// carries it.
		if parts["account_id"] == "" && segments[4] != "" {
			parts["account_id"] = segments[4]
		}
		if parts["region"] == "" && segments[3] != "" {
			parts["region"] = segments[3]
		}
		if parts["account_id"] != "" && parts["region"] != "" {
			break
		}
	}
	if len(parts) == 0 {
		return nil
	}
	return parts
}

// requiredUnresolved filters the unresolved placeholder names down to the
// ones the id_format marks required -- the only ones that block a blind
// import.
func requiredUnresolved(idFormat string, unresolved []string) []string {
	required := map[string]bool{}
	for _, name := range importmap.RequiredPlaceholders(idFormat) {
		required[name] = true
	}
	var missing []string
	for _, name := range unresolved {
		if required[name] {
			missing = append(missing, name)
		}
	}
	return missing
}

// changeCoveredBySubPaths reports whether the drift on ONE changed top-level
// attribute is fully explained by the declared dotted sub-paths: both sides
// are re-compared with those sub-paths pruned, and only if the remainder is
// identical is the change tolerated. Sub-path segments walk JSON objects and
// apply element-wise through arrays (Terraform plan JSON renders blocks as
// arrays), so "spec.update_strategy" prunes spec[i].update_strategy on both
// sides. Any structural surprise fails closed (returns false).
func changeCoveredBySubPaths(before, after interface{}, attribute string, subPaths [][]string) bool {
	if len(subPaths) == 0 {
		return false
	}
	var relevant [][]string
	for _, p := range subPaths {
		if p[0] == attribute {
			relevant = append(relevant, p[1:])
		}
	}
	if len(relevant) == 0 {
		return false
	}
	beforeMap, beforeOK := before.(map[string]interface{})
	afterMap, afterOK := after.(map[string]interface{})
	if !beforeOK || !afterOK {
		return false
	}
	prunedBefore := pruneSubPaths(beforeMap[attribute], relevant)
	prunedAfter := pruneSubPaths(afterMap[attribute], relevant)
	return reflect.DeepEqual(prunedBefore, prunedAfter)
}

// pruneSubPaths returns a deep copy of value with every declared sub-path
// removed. Arrays are traversed element-wise; scalar leaves along a declared
// path are dropped wherever the path bottoms out.
func pruneSubPaths(value interface{}, paths [][]string) interface{} {
	if len(paths) == 0 {
		return value
	}
	switch typed := value.(type) {
	case map[string]interface{}:
		copied := make(map[string]interface{}, len(typed))
		for key, child := range typed {
			var childPaths [][]string
			drop := false
			for _, p := range paths {
				if len(p) > 0 && p[0] == key {
					if len(p) == 1 {
						drop = true
						break
					}
					childPaths = append(childPaths, p[1:])
				}
			}
			if drop {
				continue
			}
			copied[key] = pruneSubPaths(child, childPaths)
		}
		return copied
	case []interface{}:
		copied := make([]interface{}, len(typed))
		for i, element := range typed {
			copied[i] = pruneSubPaths(element, paths)
		}
		return copied
	default:
		return value
	}
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
