// protoc-gen-catalog-schema emits one JSON schema document per catalog .proto
// file (the pkg/catalogschema published contract), driven by buf generate with
// the repository itself as input (see buf.gen.catalog-schema.yaml).
//
// Two responsibilities live here rather than in the pure extraction:
//
//   - rawContent: the authored .proto source is read from disk and carried
//     verbatim -- the contract's fidelity escape hatch. buf generate runs from
//     the repository root, so file.Desc.Path() is directly readable.
//   - the _test boundary: the synthetic proving kinds never reach a serving
//     surface, and the schema artifact is one -- their files are skipped, the
//     same posture as the bundle's catalog entries.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/plantonhq/planton/pkg/catalogschema"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const testProviderPathPrefix = "catalog/_test/"

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		for _, file := range plugin.Files {
			if !file.Generate {
				continue
			}
			if strings.HasPrefix(file.Desc.Path(), testProviderPathPrefix) {
				continue
			}

			schema := catalogschema.ExtractFile(file)

			rawContent, err := os.ReadFile(file.Desc.Path())
			if err != nil {
				return fmt.Errorf("read authored source of %s: %w", file.Desc.Path(), err)
			}
			schema.RawContent = string(rawContent)

			data, err := json.MarshalIndent(schema, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal %s: %w", file.Desc.Path(), err)
			}
			outputName := file.Desc.Path() + ".json"
			g := plugin.NewGeneratedFile(outputName, "")
			if _, err := g.Write(data); err != nil {
				return fmt.Errorf("write %s: %w", outputName, err)
			}
		}
		return nil
	})
}
