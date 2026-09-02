// Package protobufyaml loads YAML manifests into protobuf messages and is
// the home of the catalog's canonical YAML<->JSON conversion.
//
// The conversion speaks YAML 1.2 core-schema semantics (gopkg.in/yaml.v3):
// only `true`/`false` (and case variants) are booleans, so `y`, `n`, `yes`,
// `no`, `on`, and `off` stay strings in both key and value position, and
// duplicate mapping keys are loud errors instead of silent last-wins. This
// is deliberate: manifests are validated by protojson and diagnosed against
// their original bytes with yaml.v3, so the conversion must agree with the
// diagnoser about what the document means. YAML 1.1 lineage parsers
// (sigs.k8s.io/yaml) silently rewrite `y: 2` to `"true": 2` and `country: NO`
// to `country: false` before validation ever sees the document.
//
// Boundary: every user-manifest-shaped load or write (anything headed for
// protojson, and the files written back for those loaders to re-read) goes
// through THIS package. YAML destined for Kubernetes-lineage consumers
// (Helm chart values, cert-manager CRs) deliberately keeps sigs.k8s.io/yaml
// at its call sites — a loader must match the semantics of the system that
// consumes its output. Typed internal data files (bundle manifests,
// conversion specs, finops tables) parse strictly into Go structs where the
// boolean-token ambiguity cannot arise, and keep their existing parsers.
package protobufyaml

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// YAMLToJSON converts a YAML document to JSON under YAML 1.2 core-schema
// semantics. Only the first document of a multi-document stream is read
// (matching the YAML 1.1 converter this replaced). Non-string mapping keys
// (numbers, booleans) are stringified, since JSON objects allow nothing else.
//
// The conversion walks the parsed node tree instead of decoding into `any`,
// for one load-bearing reason: yaml.v3's plain decode resolves date-like
// scalars (`2026-01-01`) into timestamp objects and would re-render them
// RFC3339 (`2026-01-01T00:00:00Z`) — a silent rewrite of the author's text.
// Walking the tree keeps every scalar's interpretation explicit and keeps
// timestamps as the exact string the author wrote.
func YAMLToJSON(yamlBytes []byte) ([]byte, error) {
	var root yaml.Node
	if err := yaml.Unmarshal(yamlBytes, &root); err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}
	if root.Kind == 0 {
		// Empty document; the previous converter produced JSON null.
		return []byte("null"), nil
	}
	doc, err := nodeToJSONValue(&root)
	if err != nil {
		return nil, err
	}
	return json.Marshal(doc)
}

// JSONToYAML converts JSON to YAML through the same yaml.v3 lineage, so a
// document written here reads back identically through YAMLToJSON. Strings
// that would be ambiguous under 1.2 (e.g. "true", "123", "2026-01-01") come
// out quoted; strings like "off" or "y" need no quoting because 1.2 never
// mistakes them for booleans.
func JSONToYAML(jsonBytes []byte) ([]byte, error) {
	var doc any
	if err := json.Unmarshal(jsonBytes, &doc); err != nil {
		return nil, errors.Wrap(err, "parsing json")
	}
	return yaml.Marshal(doc)
}

// nodeToJSONValue converts one parsed YAML node into a JSON-encodable value.
func nodeToJSONValue(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			return nil, nil
		}
		return nodeToJSONValue(node.Content[0])
	case yaml.AliasNode:
		return nodeToJSONValue(node.Alias)
	case yaml.SequenceNode:
		items := make([]any, len(node.Content))
		for i, child := range node.Content {
			item, err := nodeToJSONValue(child)
			if err != nil {
				return nil, err
			}
			items[i] = item
		}
		return items, nil
	case yaml.MappingNode:
		result := make(map[string]any, len(node.Content)/2)
		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode, valueNode := node.Content[i], node.Content[i+1]
			if keyNode.Tag == "!!merge" {
				// YAML merge keys (`<<:`) are a 1.1 feature with no usage in
				// this catalog; refusing loudly beats mis-parsing quietly.
				return nil, errors.Errorf("line %d: YAML merge keys (<<:) are not supported; spell the mapping out", keyNode.Line)
			}
			key, err := scalarKeyString(keyNode)
			if err != nil {
				return nil, err
			}
			if _, exists := result[key]; exists {
				return nil, errors.Errorf("line %d: mapping key %q is already defined", keyNode.Line, key)
			}
			value, err := nodeToJSONValue(valueNode)
			if err != nil {
				return nil, err
			}
			result[key] = value
		}
		return result, nil
	case yaml.ScalarNode:
		return scalarValue(node)
	default:
		return nil, errors.Errorf("line %d: unsupported yaml node kind %d", node.Line, node.Kind)
	}
}

// scalarKeyString renders a mapping key as the string JSON requires. YAML
// allows number and boolean keys; their canonical text form becomes the key.
func scalarKeyString(node *yaml.Node) (string, error) {
	if node.Kind != yaml.ScalarNode {
		return "", errors.Errorf("line %d: mapping keys must be scalars", node.Line)
	}
	value, err := scalarValue(node)
	if err != nil {
		return "", err
	}
	if s, ok := value.(string); ok {
		return s, nil
	}
	return fmt.Sprintf("%v", value), nil
}

// scalarValue interprets one scalar by the tag yaml.v3's 1.2 resolver
// assigned it. Timestamps deliberately stay the author's exact string.
func scalarValue(node *yaml.Node) (any, error) {
	switch node.Tag {
	case "!!null":
		return nil, nil
	case "!!bool":
		var b bool
		if err := node.Decode(&b); err != nil {
			return nil, errors.Wrapf(err, "line %d: parsing boolean %q", node.Line, node.Value)
		}
		return b, nil
	case "!!int":
		// Let yaml.v3 parse its own integer grammar (hex, octal,
		// underscores); fall through widening types on overflow.
		var i int64
		if err := node.Decode(&i); err == nil {
			return i, nil
		}
		var u uint64
		if err := node.Decode(&u); err == nil {
			return u, nil
		}
		var f float64
		if err := node.Decode(&f); err != nil {
			return nil, errors.Wrapf(err, "line %d: parsing integer %q", node.Line, node.Value)
		}
		return f, nil
	case "!!float":
		var f float64
		if err := node.Decode(&f); err != nil {
			return nil, errors.Wrapf(err, "line %d: parsing float %q", node.Line, node.Value)
		}
		return f, nil
	default:
		// !!str, !!timestamp, and any custom tag: the author's exact text.
		return node.Value, nil
	}
}
