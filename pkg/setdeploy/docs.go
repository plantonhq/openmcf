package setdeploy

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/planton/internal/manifest"
)

// Doc is one YAML document offered to the set, with the source label findings
// and report entries carry (a file path, or "path[docN]" for one document of
// a multi-document stream).
type Doc struct {
	Bytes  []byte
	Source string
}

// CollectDocsFromBytes splits a (possibly multi-document) YAML stream into
// docs. Single-document input returns one doc labeled with the source alone;
// documents of a stream are labeled "source[docN]" so a refusal can name the
// exact document inside the file the author wrote.
func CollectDocsFromBytes(yamlBytes []byte, source string) ([]Doc, error) {
	split, err := manifest.SplitDocuments(yamlBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to split %s into YAML documents", source)
	}
	if len(split) == 1 {
		return []Doc{{Bytes: split[0], Source: source}}, nil
	}
	docs := make([]Doc, 0, len(split))
	for i, d := range split {
		docs = append(docs, Doc{Bytes: d, Source: fmt.Sprintf("%s[doc%d]", source, i+1)})
	}
	return docs, nil
}

// CollectDocsFromDir reads every .yaml/.yml file directly under dir (sorted
// by name — the kubectl mental model: a directory of manifests is a set, and
// authored file order is the deterministic tie-break downstream). Each file
// may itself hold multiple documents; nested directories are deliberately not
// walked in v1 — an overlay tree's internals are kustomize's job, not a flat
// manifest folder's.
func CollectDocsFromDir(dir string) ([]Doc, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read manifest directory %s", dir)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		names = append(names, e.Name())
	}
	sort.Strings(names)
	if len(names) == 0 {
		return nil, errors.Errorf("no .yaml or .yml manifests found in %s", dir)
	}

	var docs []Doc
	for _, name := range names {
		path := filepath.Join(dir, name)
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read %s", path)
		}
		fileDocs, err := CollectDocsFromBytes(b, path)
		if err != nil {
			return nil, err
		}
		docs = append(docs, fileDocs...)
	}
	return docs, nil
}
