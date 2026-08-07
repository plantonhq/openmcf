package catalogbundle

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"sigs.k8s.io/yaml"
)

// Bundle is a loaded, checksum-verified catalog bundle.
type Bundle struct {
	Manifest Manifest
	// Files resolves every schema the bundle carries (kind messages,
	// validation rules in options, the kind registry enum).
	Files *protoregistry.Files
	// entries holds the raw bundle contents by entry name (conversions/**,
	// presets/**, descriptors.binpb).
	entries map[string][]byte
}

// Load opens a bundle zip, verifies every entry against the manifest's
// checksums, and parses the descriptor set. A bundle that fails ANY check is
// refused whole -- consumers never run on partially-verified data.
func Load(path string) (*Bundle, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("opening catalog bundle %s: %w", path, err)
	}
	defer reader.Close()

	entries := map[string][]byte{}
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("reading bundle entry %s: %w", file.Name, err)
		}
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("reading bundle entry %s: %w", file.Name, err)
		}
		entries[file.Name] = content
	}

	manifestRaw, ok := entries["manifest.yaml"]
	if !ok {
		return nil, fmt.Errorf("catalog bundle %s has no manifest.yaml -- not a bundle", path)
	}
	delete(entries, "manifest.yaml")
	var manifest Manifest
	if err := yaml.UnmarshalStrict(manifestRaw, &manifest); err != nil {
		return nil, fmt.Errorf("catalog bundle manifest is malformed: %w", err)
	}
	if manifest.FormatVersion != FormatVersion {
		return nil, fmt.Errorf("catalog bundle format %q is not understood by this consumer (wants %q) -- upgrade the consumer or rebuild the bundle", manifest.FormatVersion, FormatVersion)
	}

	// Self-verification: every entry must match its recorded checksum, and
	// every recorded checksum must have its entry.
	for name, content := range entries {
		want, recorded := manifest.Checksums[name]
		if !recorded {
			return nil, fmt.Errorf("bundle entry %s is not in the manifest -- the bundle was tampered with or mis-built", name)
		}
		sum := sha256.Sum256(content)
		if hex.EncodeToString(sum[:]) != want {
			return nil, fmt.Errorf("bundle entry %s fails its checksum -- the bundle is corrupt", name)
		}
	}
	for name := range manifest.Checksums {
		if _, present := entries[name]; !present {
			return nil, fmt.Errorf("the manifest records %s but the bundle does not contain it -- the bundle is incomplete", name)
		}
	}

	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(entries["descriptors.binpb"], &fds); err != nil {
		return nil, fmt.Errorf("parsing the bundle's descriptor set: %w", err)
	}
	files, err := protodesc.NewFiles(&fds)
	if err != nil {
		return nil, fmt.Errorf("building a registry from the bundle's descriptors: %w", err)
	}

	return &Bundle{Manifest: manifest, Files: files, entries: entries}, nil
}

// ConversionSpecs returns the bundle's conversion spec contents keyed by
// entry name (conversions/<provider>/<kind>/<file>.yaml), sorted.
func (b *Bundle) ConversionSpecs() map[string][]byte {
	return b.subtree("conversions/")
}

// Presets returns the bundle's preset contents keyed by entry name
// (presets/<provider>/<kind>/<file>), sorted.
func (b *Bundle) Presets() map[string][]byte {
	return b.subtree("presets/")
}

func (b *Bundle) subtree(prefix string) map[string][]byte {
	out := map[string][]byte{}
	names := make([]string, 0, len(b.entries))
	for name := range b.entries {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if strings.HasPrefix(name, prefix) {
			out[name] = b.entries[name]
		}
	}
	return out
}
