// pkg/protodocs/distiller/main.go
// Distills proto-source documentation into pkg/protodocs' embedded index.
//
// The Go protobuf runtime strips comments from generated code, so the only
// machine-readable home of the .proto documentation is a compiled descriptor
// image built with source info (buf build). Embedding that image whole would
// carry spans, paths, and structure the CLI never reads; this tool reduces it
// to the one thing runtime reflection cannot provide -- a fully-qualified
// name -> documentation map -- so the committed artifact stays a fraction of
// the image and lookup needs no descriptor arithmetic at runtime.
//
// Run via Makefile: make generate-proto-docs
//
//	go run ./pkg/protodocs/distiller \
//	    --image build/proto-docs-image.binpb \
//	    --out pkg/protodocs/index.json.gz \
//	    --include catalog --include shared --include qa --include iac
//
// Output is deterministic (sorted file list, JSON object keys sorted by
// encoding/json, gzip with an empty header), so regenerating without proto
// changes produces a byte-identical artifact.
package main

import (
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// index mirrors pkg/protodocs' embedded payload. Kept in sync by hand: the
// two packages must not import each other (protodocs go:embeds the artifact
// this tool writes, so importing it here would make first generation
// impossible on a clean tree).
type index struct {
	Files []string          `json:"files"`
	Docs  map[string]string `json:"docs"`
}

type prefixList []string

func (p *prefixList) String() string { return strings.Join(*p, ",") }
func (p *prefixList) Set(v string) error {
	*p = append(*p, v)
	return nil
}

func main() {
	var imagePath, outPath string
	var includes prefixList
	flag.StringVar(&imagePath, "image", "", "path to a FileDescriptorSet built WITH source info (buf build -o ...)")
	flag.StringVar(&outPath, "out", "", "path to write the gzipped JSON index")
	flag.Var(&includes, "include", "proto file path prefix to include (repeatable; default catalog+shared+qa+iac)")
	flag.Parse()

	if imagePath == "" || outPath == "" {
		fmt.Fprintln(os.Stderr, "usage: distiller --image <image.binpb> --out <index.json.gz> [--include <path-prefix>]...")
		os.Exit(2)
	}
	if len(includes) == 0 {
		includes = prefixList{"catalog", "shared", "qa", "iac"}
	}

	if err := run(imagePath, outPath, includes); err != nil {
		fmt.Fprintln(os.Stderr, "distiller: "+err.Error())
		os.Exit(1)
	}
}

func run(imagePath, outPath string, includes []string) error {
	raw, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}
	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(raw, &fds); err != nil {
		return fmt.Errorf("parse descriptor image: %w", err)
	}

	out := index{Docs: map[string]string{}}
	for _, file := range fds.GetFile() {
		if !included(file.GetName(), includes) {
			continue
		}
		out.Files = append(out.Files, file.GetName())
		if file.GetSourceCodeInfo() == nil {
			return fmt.Errorf("%s has no source info -- build the image without --exclude-source-info", file.GetName())
		}
		for _, loc := range file.GetSourceCodeInfo().GetLocation() {
			doc := cleanComment(loc.GetLeadingComments())
			if doc == "" {
				continue
			}
			fullName, ok := resolvePath(file, loc.GetPath())
			if !ok {
				continue
			}
			out.Docs[fullName] = doc
		}
	}
	if len(out.Files) == 0 {
		return fmt.Errorf("no proto files matched include prefixes %v", includes)
	}
	sort.Strings(out.Files)

	body, err := json.Marshal(out)
	if err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	zw, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		return err
	}
	if _, err := zw.Write(body); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	fmt.Printf("distilled %d documented elements from %d files into %s\n",
		len(out.Docs), len(out.Files), outPath)
	return nil
}

func included(fileName string, includes []string) bool {
	for _, prefix := range includes {
		if strings.HasPrefix(fileName, prefix) {
			return true
		}
	}
	return false
}

// resolvePath maps a SourceCodeInfo location path onto the fully-qualified
// name protoreflect reports for the same element at runtime. Only the
// elements explain renders are resolved (messages, fields, oneofs, enums,
// enum values); everything else (services, options, syntax, ...) is skipped.
//
// FileDescriptorProto field numbers anchoring the paths:
//
//	4 = message_type, 5 = enum_type (file level)
//	2 = field, 3 = nested_type, 4 = enum_type, 8 = oneof_decl (in a message)
//	2 = value (in an enum)
func resolvePath(file *descriptorpb.FileDescriptorProto, path []int32) (string, bool) {
	if len(path) < 2 {
		return "", false
	}
	pkg := file.GetPackage()
	switch path[0] {
	case 4:
		msgs := file.GetMessageType()
		if int(path[1]) >= len(msgs) {
			return "", false
		}
		return resolveMessage(pkg, msgs[path[1]], path[2:])
	case 5:
		enums := file.GetEnumType()
		if int(path[1]) >= len(enums) {
			return "", false
		}
		return resolveEnum(pkg, enums[path[1]], path[2:])
	}
	return "", false
}

func resolveMessage(scope string, msg *descriptorpb.DescriptorProto, rest []int32) (string, bool) {
	name := scope + "." + msg.GetName()
	if len(rest) == 0 {
		return name, true
	}
	if len(rest) < 2 {
		return "", false
	}
	switch rest[0] {
	case 2: // field: full name is <message>.<field proto name>
		fields := msg.GetField()
		if int(rest[1]) >= len(fields) || len(rest) != 2 {
			return "", false
		}
		return name + "." + fields[rest[1]].GetName(), true
	case 3: // nested message
		nested := msg.GetNestedType()
		if int(rest[1]) >= len(nested) {
			return "", false
		}
		return resolveMessage(name, nested[rest[1]], rest[2:])
	case 4: // nested enum
		enums := msg.GetEnumType()
		if int(rest[1]) >= len(enums) {
			return "", false
		}
		return resolveEnum(name, enums[rest[1]], rest[2:])
	case 8: // oneof declaration
		oneofs := msg.GetOneofDecl()
		if int(rest[1]) >= len(oneofs) || len(rest) != 2 {
			return "", false
		}
		return name + "." + oneofs[rest[1]].GetName(), true
	}
	return "", false
}

func resolveEnum(scope string, enum *descriptorpb.EnumDescriptorProto, rest []int32) (string, bool) {
	name := scope + "." + enum.GetName()
	if len(rest) == 0 {
		return name, true
	}
	// Enum VALUES scope to the enum's parent (C++ scoping rules) -- exactly
	// how protoreflect computes EnumValueDescriptor.FullName, so runtime
	// lookups by that name land here.
	if rest[0] == 2 && len(rest) == 2 {
		values := enum.GetValue()
		if int(rest[1]) >= len(values) {
			return "", false
		}
		return scope + "." + values[rest[1]].GetName(), true
	}
	return "", false
}

// cleanComment normalizes protoc's raw leading_comments into display prose:
// the conventional single space after // is dropped, block-comment asterisk
// gutters (the repo's /** ... */ house style) are stripped, and surrounding
// blank lines are trimmed. Interior blank lines survive -- they are the
// authors' paragraph breaks.
//
// Order matters: the single leading space is dropped BEFORE gutter
// detection, and a bare "*" is only treated as a gutter when followed by a
// space (or alone). A line beginning "**bold**" is markdown, not a gutter --
// stripping its first "*" was a real defect this ordering fixes.
func cleanComment(raw string) string {
	if raw == "" {
		return ""
	}
	lines := strings.Split(raw, "\n")
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimRight(line, " \t")
		line = strings.TrimPrefix(line, " ")
		switch {
		case line == "*":
			line = ""
		case strings.HasPrefix(line, "* "):
			line = line[2:]
		case strings.HasPrefix(strings.TrimLeft(line, " \t"), "* "):
			// Indented gutter (block comments nested under message bodies).
			line = strings.TrimLeft(line, " \t")[2:]
		case strings.TrimLeft(line, " \t") == "*":
			line = ""
		}
		cleaned = append(cleaned, line)
	}
	for len(cleaned) > 0 && cleaned[0] == "" {
		cleaned = cleaned[1:]
	}
	for len(cleaned) > 0 && cleaned[len(cleaned)-1] == "" {
		cleaned = cleaned[:len(cleaned)-1]
	}
	return strings.Join(cleaned, "\n")
}
