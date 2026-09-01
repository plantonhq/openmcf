package manifest

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	goyaml "gopkg.in/yaml.v3"
)

// documentIdentity names one document of a multi-document stream for the
// refusal message: enough for the author to recognize each resource.
type documentIdentity struct {
	Kind string
	Name string
}

func (d documentIdentity) String() string {
	kind := d.Kind
	if kind == "" {
		kind = "(no kind)"
	}
	name := d.Name
	if name == "" {
		name = "(unnamed)"
	}
	return kind + "/" + name
}

// refuseMultiDocument errors when the bytes hold more than one non-empty YAML
// document, naming every document found. Empty documents (separator-only
// segments, comment-only segments) do not count, so a file that merely leads
// with `---` or ends with a stray separator loads normally.
func refuseMultiDocument(manifestYamlBytes []byte, sourceName string) error {
	identities, err := collectDocumentIdentities(manifestYamlBytes)
	if err != nil {
		// Not this gate's job: let the schema-aware loader below produce its
		// richer diagnosis for malformed YAML.
		return nil
	}
	if len(identities) <= 1 {
		return nil
	}

	var list strings.Builder
	for i, identity := range identities {
		list.WriteString(fmt.Sprintf("  %d. %s\n", i+1, identity))
	}
	return errors.Errorf(
		"%s contains %d YAML documents — this command applies exactly one manifest per run:\n%s"+
			"apply each document separately (kustomize overlays that bundle several resources render as one multi-document stream)",
		sourceName, len(identities), list.String())
}

// SplitDocuments splits a (possibly multi-document) YAML stream into
// per-document byte slices, skipping empty documents. Callers that
// legitimately handle multi-document YAML — preset validation, catalog
// tooling — split with this and load each document through
// LoadManifestBytes, whose contract is exactly one document.
// The yaml.Node round-trip preserves each document's content (including
// comments) without inventing a second YAML dialect.
func SplitDocuments(manifestYamlBytes []byte) ([][]byte, error) {
	decoder := goyaml.NewDecoder(bytes.NewReader(manifestYamlBytes))
	var docs [][]byte
	for {
		var node goyaml.Node
		err := decoder.Decode(&node)
		if err == io.EOF {
			return docs, nil
		}
		if err != nil {
			return nil, err
		}
		if isNullDocument(&node) {
			continue
		}
		docBytes, err := goyaml.Marshal(&node)
		if err != nil {
			return nil, err
		}
		docs = append(docs, docBytes)
	}
}

// isNullDocument reports whether a decoded document node holds no content —
// a separator-only or comment-only segment of the stream.
func isNullDocument(node *goyaml.Node) bool {
	if node.Kind == 0 {
		return true
	}
	if node.Kind == goyaml.DocumentNode {
		if len(node.Content) == 0 {
			return true
		}
		if len(node.Content) == 1 && node.Content[0].Tag == "!!null" {
			return true
		}
	}
	return false
}

// collectDocumentIdentities decodes the stream document by document and
// returns the identity of each non-empty one.
func collectDocumentIdentities(manifestYamlBytes []byte) ([]documentIdentity, error) {
	decoder := goyaml.NewDecoder(bytes.NewReader(manifestYamlBytes))
	var identities []documentIdentity
	for {
		var doc map[string]interface{}
		err := decoder.Decode(&doc)
		if err == io.EOF {
			return identities, nil
		}
		if err != nil {
			return nil, err
		}
		if len(doc) == 0 {
			continue
		}
		identity := documentIdentity{}
		if kind, ok := doc["kind"].(string); ok {
			identity.Kind = kind
		}
		if metadata, ok := doc["metadata"].(map[string]interface{}); ok {
			if name, ok := metadata["name"].(string); ok {
				identity.Name = name
			}
		}
		identities = append(identities, identity)
	}
}
