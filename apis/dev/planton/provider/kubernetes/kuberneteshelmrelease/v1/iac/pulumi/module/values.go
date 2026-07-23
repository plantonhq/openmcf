package module

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
	kuberneteshelmreleasev1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteshelmrelease/v1"
	"helm.sh/helm/v3/pkg/strvals"
	"sigs.k8s.io/yaml"
)

// buildHelmValues merges the spec's values layers into the final values map
// handed to Helm, with the documented precedence: values_yaml first, then
// set, then set_string, then set_sensitive (highest).
//
// PARITY: the Terraform module reaches the same result through the
// helm_release resource's native mechanisms — values = [values_yaml] plus
// set blocks (type auto for `set`, type string for `set_string`) and
// set_sensitive blocks — which the provider merges in exactly this order.
// Both paths run Helm's own strvals parser on the override entries, so
// dotted-path syntax and type coercion behave identically on both engines.
// Map entries are applied in sorted-key order on both engines (Terraform's
// for_each iterates maps lexically), keeping even same-path collisions
// deterministic.
func buildHelmValues(spec *kuberneteshelmreleasev1.KubernetesHelmReleaseSpec) (map[string]interface{}, error) {
	values := map[string]interface{}{}

	if spec.GetValuesYaml() != "" {
		if err := yaml.Unmarshal([]byte(spec.GetValuesYaml()), &values); err != nil {
			return nil, errors.Wrap(err, "failed to parse values_yaml as a YAML document")
		}
	}

	// `set` entries use Helm's --set coercion ("true" -> bool, digits ->
	// number, "null" removes the key).
	for _, key := range sortedKeys(spec.GetSet()) {
		entry := fmt.Sprintf("%s=%s", key, spec.GetSet()[key])
		if err := strvals.ParseInto(entry, values); err != nil {
			return nil, errors.Wrapf(err, "failed to apply set entry %q", key)
		}
	}

	// `set_string` and `set_sensitive` entries keep values as literal
	// strings (Helm's --set-string).
	for _, key := range sortedKeys(spec.GetSetString()) {
		entry := fmt.Sprintf("%s=%s", key, spec.GetSetString()[key])
		if err := strvals.ParseIntoString(entry, values); err != nil {
			return nil, errors.Wrapf(err, "failed to apply set_string entry %q", key)
		}
	}

	for _, key := range sortedKeys(spec.GetSetSensitive()) {
		entry := fmt.Sprintf("%s=%s", key, spec.GetSetSensitive()[key])
		if err := strvals.ParseIntoString(entry, values); err != nil {
			return nil, errors.Wrapf(err, "failed to apply set_sensitive entry %q", key)
		}
	}

	return values, nil
}

// sortedKeys returns the map's keys in lexical order, mirroring Terraform's
// deterministic map iteration so both engines apply overrides in the same
// order.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
