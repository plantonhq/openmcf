package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// manifestDocHeader is the subset of a Kubernetes document needed to build
// the applied-resource inventory.
type manifestDocHeader struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name string `json:"name"`
	} `json:"metadata"`
}

// parseAppliedResources derives the applied-resource inventory
// ("apiVersion/Kind/name" per document, manifest order) by parsing the input
// YAML. The Terraform module derives the identical list the same way (its
// locals split and decode the same documents), so both engines export the
// same inventory by construction — engine-side child-resource reflection is
// never consulted.
//
// Documents are split on the YAML document separator (a line starting with
// `---`); blank documents (a leading/trailing ---, comment-only fragments)
// are skipped. The leading newline prepended before splitting makes a
// manifest that STARTS with `---` yield an empty first chunk instead of a
// missed document — the Terraform module's locals use the identical
// prepend-then-split rule.
func parseAppliedResources(manifestYaml string) ([]string, error) {
	var inventory []string
	for _, doc := range strings.Split("\n"+manifestYaml, "\n---") {
		if strings.TrimSpace(doc) == "" {
			continue
		}
		var header manifestDocHeader
		if err := yaml.Unmarshal([]byte(doc), &header); err != nil {
			return nil, errors.Wrap(err, "manifest_yaml contains a document that is not valid YAML")
		}
		if header.Kind == "" && header.APIVersion == "" && header.Metadata.Name == "" {
			// Comment-only / empty document — nothing is applied for it.
			continue
		}
		inventory = append(inventory,
			fmt.Sprintf("%s/%s/%s", header.APIVersion, header.Kind, header.Metadata.Name))
	}
	return inventory, nil
}
