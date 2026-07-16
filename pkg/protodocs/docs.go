// Package protodocs serves the proto-source documentation for the schemas
// compiled into this binary, fully offline.
//
// The Go protobuf runtime deliberately strips comments from generated code,
// so the rich field/message documentation written in the .proto sources is
// invisible to runtime reflection. This package closes that gap: a build
// step (make generate-proto-docs) compiles the proto module with source
// info, distills it to a full-name -> documentation index, and commits the
// gzipped result here for go:embed. Lookup is the read side.
//
// The index is a published surface: `planton explain` renders these strings
// to humans and AI agents as the API reference. Keep the distiller's
// normalization stable and extend the index format additively.
package protodocs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"

	_ "embed"
)

//go:embed index.json.gz
var indexGz []byte

// index is the distilled documentation payload. Files carries every proto
// file path the index was built from so tests can prove the committed
// artifact still covers everything compiled into the binary (a new kind
// without a regenerated index fails loudly instead of silently rendering
// undocumented).
type index struct {
	Files []string          `json:"files"`
	Docs  map[string]string `json:"docs"`
}

var (
	loadOnce sync.Once
	loaded   index
)

// load decompresses and parses the embedded index exactly once. A corrupt
// artifact is a build defect, not a runtime condition, so failures panic.
func load() {
	loadOnce.Do(func() {
		zr, err := gzip.NewReader(bytes.NewReader(indexGz))
		if err != nil {
			panic("protodocs: embedded index is not valid gzip: " + err.Error())
		}
		body, err := io.ReadAll(zr)
		if err != nil {
			panic("protodocs: failed to decompress embedded index: " + err.Error())
		}
		if err := json.Unmarshal(body, &loaded); err != nil {
			panic("protodocs: embedded index is not valid JSON: " + err.Error())
		}
	})
}

// Lookup returns the proto-source documentation for a message, field, enum,
// or enum value, keyed by its fully-qualified proto name. Empty string means
// the element carries no documentation.
func Lookup(fullName protoreflect.FullName) string {
	load()
	return loaded.Docs[string(fullName)]
}

// Files returns the proto file paths the embedded index was distilled from,
// for coverage tests.
func Files() []string {
	load()
	return loaded.Files
}
