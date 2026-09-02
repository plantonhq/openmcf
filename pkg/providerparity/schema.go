//go:build !codegen
// +build !codegen

package providerparity

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// DefaultSchemaDir is repo-root-relative (where the CLI runs). Tests read the
// directory by its bare name because `go test` runs with the package directory
// as cwd. The artifacts are loaded from the repo tree, NOT go:embed-ded: they
// are repo-audit data that grows with every provider program, and embedding
// them would carry megabytes into the customer-facing binary for data only
// repo gates read.
const DefaultSchemaDir = "pkg/providerparity/schemas"

// Schema is the committed, distilled description of ONE Terraform provider at
// ONE pinned version -- the provider describing itself via
// `providers schema -json`, reduced to what parity accounting reads: resource
// names and their configurable argument surface. Descriptions are dropped
// (the design floor for modeling a field remains the provider source at the
// pin); data sources are out of scope (parity accounting is over resources).
type Schema struct {
	// Provider is the registry short name and artifact identity, e.g.
	// "google", "google-beta". Loaders key on this field, never on filenames.
	Provider string `json:"provider"`
	// Source is the full registry address, e.g. "hashicorp/google".
	Source string `json:"source"`
	// Version is the exact resolved release the artifact was distilled from,
	// e.g. "6.55.0" -- the named pin parity is declared against.
	Version string `json:"version"`
	// ProviderConfig is the provider's own configuration block: the arguments
	// a `provider "<name>" {}` block accepts (assume_role, default_tags,
	// endpoints, ...). The provider-config accounting reads it the same way
	// depth accounting reads resource blocks. Named ProviderConfig because
	// Provider (above) is the artifact's identity field.
	ProviderConfig *Block `json:"providerConfig,omitempty"`
	// Resources maps resource type (e.g. "google_storage_bucket") to its
	// root configuration block.
	Resources map[string]*Block `json:"resources"`
}

// Block is one configuration block: its direct attributes plus nested blocks.
type Block struct {
	Attributes map[string]*Attribute   `json:"attributes,omitempty"`
	Blocks     map[string]*NestedBlock `json:"blocks,omitempty"`
	Deprecated bool                    `json:"deprecated,omitempty"`
}

// Attribute is one provider argument. Computed-only attributes (Computed set,
// neither Required nor Optional) are provider outputs, not configuration --
// they are kept in the artifact for output-parity work but excluded from the
// configurable surface.
type Attribute struct {
	// Type is the provider's own type expression, passed through verbatim
	// (a JSON string like "string" or an array like ["list","string"]).
	Type       json.RawMessage `json:"type,omitempty"`
	Required   bool            `json:"required,omitempty"`
	Optional   bool            `json:"optional,omitempty"`
	Computed   bool            `json:"computed,omitempty"`
	Sensitive  bool            `json:"sensitive,omitempty"`
	Deprecated bool            `json:"deprecated,omitempty"`
}

// NestedBlock is a named nested configuration block and its nesting mode
// ("single", "list", "set", "map").
type NestedBlock struct {
	NestingMode string `json:"nestingMode,omitempty"`
	Block       *Block `json:"block,omitempty"`
}

// Arg is one configurable argument flattened to a dotted path, e.g.
// "lifecycle_rule.condition.age". The unit the total-accounting check
// matches, maps, or requires an exclusion for.
type Arg struct {
	Path       string
	Required   bool
	Sensitive  bool
	Deprecated bool // set on the attribute itself or inherited from a deprecated enclosing block
}

// ConfigurableArgs flattens a block's configurable surface (Required or
// Optional attributes; computed-only excluded) into sorted dotted paths,
// descending nested blocks. prefix is empty for a resource's root block.
func (b *Block) ConfigurableArgs(prefix string) []Arg {
	var args []Arg
	b.collectArgs(prefix, b.Deprecated, &args)
	sort.Slice(args, func(i, j int) bool { return args[i].Path < args[j].Path })
	return args
}

func (b *Block) collectArgs(prefix string, deprecated bool, out *[]Arg) {
	join := func(name string) string {
		if prefix == "" {
			return name
		}
		return prefix + "." + name
	}
	for name, a := range b.Attributes {
		if !a.Required && !a.Optional {
			continue // computed-only: an output, not configuration
		}
		*out = append(*out, Arg{
			Path:       join(name),
			Required:   a.Required,
			Sensitive:  a.Sensitive,
			Deprecated: deprecated || a.Deprecated,
		})
	}
	for name, nb := range b.Blocks {
		if nb.Block == nil {
			continue
		}
		nb.Block.collectArgs(join(name), deprecated || nb.Block.Deprecated, out)
	}
}

// ConfigurableArgCount is the size of the non-deprecated configurable
// surface -- the denominator-side unit of depth parity.
func (b *Block) ConfigurableArgCount() int {
	n := 0
	for _, a := range b.ConfigurableArgs("") {
		if !a.Deprecated {
			n++
		}
	}
	return n
}

// LoadSchemas reads every artifact in dir and returns them keyed by provider
// name. Artifact identity comes from the content's `provider` field, never
// from the filename (filenames carry the version for humans and diffs).
// Exactly one artifact per provider: a second one means a stale artifact
// survived a pin bump, which would make "the pin" ambiguous -- fail loudly.
func LoadSchemas(dir string) (map[string]*Schema, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "reading schema dir %s", dir)
	}
	schemas := map[string]*Schema{}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json.gz") {
			continue
		}
		s, err := ReadSchemaFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		if prev, ok := schemas[s.Provider]; ok {
			return nil, errors.Errorf(
				"two schema artifacts for provider %q in %s (versions %s and %s) -- remove the stale one; the pin must be unambiguous",
				s.Provider, dir, prev.Version, s.Version)
		}
		schemas[s.Provider] = s
	}
	return schemas, nil
}

// LoadSchema returns dir's artifact for one provider, failing plainly when it
// is absent (the artifact is committed; absence means the distiller was never
// run for this provider or the tree is incomplete).
func LoadSchema(dir, provider string) (*Schema, error) {
	schemas, err := LoadSchemas(dir)
	if err != nil {
		return nil, err
	}
	s, ok := schemas[provider]
	if !ok {
		return nil, errors.Errorf(
			"no schema artifact for provider %q in %s -- run `make generate-provider-schemas`",
			provider, dir)
	}
	return s, nil
}

// ReadSchemaFile decodes one gzipped JSON artifact.
func ReadSchemaFile(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "opening schema artifact %s", path)
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return nil, errors.Wrapf(err, "schema artifact %s is not valid gzip", path)
	}
	defer zr.Close()
	raw, err := io.ReadAll(zr)
	if err != nil {
		return nil, errors.Wrapf(err, "decompressing schema artifact %s", path)
	}
	var s Schema
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, errors.Wrapf(err, "schema artifact %s is not valid JSON", path)
	}
	if s.Provider == "" || s.Version == "" {
		return nil, errors.Errorf("schema artifact %s carries no provider/version identity", path)
	}
	return &s, nil
}

// WriteSchemaFile writes s to dir as <provider>-<version>.json.gz, first
// removing any existing artifact with the same provider identity so a pin
// bump replaces the artifact instead of accumulating versions. Output is
// deterministic (encoding/json sorts map keys; gzip header left empty), so
// re-distilling an unchanged pin produces a byte-identical file.
func WriteSchemaFile(dir string, s *Schema) (string, error) {
	existing, err := LoadSchemas(dir)
	if err != nil && !os.IsNotExist(errors.Cause(err)) {
		return "", err
	}
	if prev, ok := existing[s.Provider]; ok {
		stale := filepath.Join(dir, prev.Provider+"-"+prev.Version+".json.gz")
		if err := os.Remove(stale); err != nil && !os.IsNotExist(err) {
			return "", errors.Wrapf(err, "removing stale artifact %s", stale)
		}
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", errors.Wrapf(err, "creating schema dir %s", dir)
	}
	raw, err := json.Marshal(s)
	if err != nil {
		return "", errors.Wrap(err, "encoding schema")
	}
	path := filepath.Join(dir, s.Provider+"-"+s.Version+".json.gz")
	f, err := os.Create(path)
	if err != nil {
		return "", errors.Wrapf(err, "creating %s", path)
	}
	zw := gzip.NewWriter(f)
	if _, err := zw.Write(raw); err != nil {
		f.Close()
		return "", errors.Wrapf(err, "writing %s", path)
	}
	if err := zw.Close(); err != nil {
		f.Close()
		return "", errors.Wrapf(err, "finalizing %s", path)
	}
	return path, f.Close()
}
