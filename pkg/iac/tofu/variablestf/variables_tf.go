package variablestf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantonhq/openmcf/internal/apidocs"
	"github.com/plantonhq/openmcf/pkg/strings/caseconverter"
	gendoc "github.com/pseudomuto/protoc-gen-doc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var fieldDescriptions = map[string]string{
	"metadata": "Metadata for the resource, including name and labels",
	"spec":     "Specification for Deployment Component",
}

// terraformType is an interface representing a Terraform type.
type terraformType interface {
	format(indentLevel int) string
}

type tfPrimitive string

func (p tfPrimitive) format(indentLevel int) string {
	return string(p)
}

type tfList struct {
	elem terraformType
}

func (l tfList) format(indentLevel int) string {
	return fmt.Sprintf("list(%s)", l.elem.format(indentLevel))
}

// tfField holds a field's terraform type and description for inline comments
type tfField struct {
	name        string
	description string
	t           terraformType
}

type tfObject struct {
	fields []tfField
}

func (o tfObject) format(indentLevel int) string {
	if len(o.fields) == 0 {
		return "object({})"
	}

	indent := strings.Repeat("  ", indentLevel)
	nextIndent := strings.Repeat("  ", indentLevel+1)

	var fieldLines []string
	for _, f := range o.fields {
		if f.description != "" {
			// Add a blank line before comments to improve readability
			fieldLines = append(fieldLines, "")
			commentLines := strings.Split(f.description, "\n")
			for _, cl := range commentLines {
				fieldLines = append(fieldLines, fmt.Sprintf("%s# %s", nextIndent, cl))
			}
		}

		fieldStr := f.t.format(indentLevel + 1)
		fieldLines = append(fieldLines, fmt.Sprintf("%s%s = %s", nextIndent, f.name, fieldStr))
	}

	return fmt.Sprintf("object({\n%s\n%s})", strings.Join(fieldLines, "\n"), indent)
}

// ProtoToVariablesTF uses proto reflection to determine the Terraform variable schema.
// It now also includes proto field comments as inline comments for object fields and
// will skip the 'version' field inside a 'Metadata' message.
func ProtoToVariablesTF(msg proto.Message) (string, error) {
	apiDocsJson, err := apidocs.GetApiDocsJson()
	if err != nil {
		return "", errors.Wrap(err, "failed to get api docs")
	}

	md := msg.ProtoReflect().Descriptor()

	skipFields := map[string]bool{
		"api_version": true,
		"kind":        true,
		"status":      true,
	}

	var buf bytes.Buffer
	fields := md.Fields()

	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		fieldName := string(fd.Name())

		if skipFields[fieldName] {
			continue
		}

		tfType, err := fieldDescriptorToTerraformType(fd, md, apiDocsJson)
		if err != nil {
			return "", errors.Wrapf(err, "failed to convert field %q to terraform type", fieldName)
		}

		desc := findFieldDescription(apiDocsJson, string(md.FullName()), fieldName)
		if desc == "" && fd.Kind() == protoreflect.MessageKind {
			desc = findMessageDescription(apiDocsJson, string(fd.Message().FullName()))
		}

		if desc == "" {
			desc = fieldDescriptions[fieldName]
			if desc == "" {
				desc = fmt.Sprintf("Description for %s", fieldName)
			}
		}

		typeStr := tfType.format(1)

		fmt.Fprintf(&buf, `variable "%s" {
  description = %q
  type = %s
}

`, caseconverter.ToSnakeCase(fieldName), desc, typeStr)
	}

	return strings.TrimSpace(buf.String()), nil
}

func fieldDescriptorToTerraformType(fd protoreflect.FieldDescriptor, parentMsg protoreflect.MessageDescriptor, apiDocsJson *gendoc.Template) (terraformType, error) {
	// If repeated -> list
	if fd.IsList() {
		elemType, err := scalarOrMessageToTFType(parentMsg, fd, apiDocsJson)
		if err != nil {
			return nil, err
		}
		return tfList{elem: elemType}, nil
	}

	return scalarOrMessageToTFType(parentMsg, fd, apiDocsJson)
}

func scalarOrMessageToTFType(parentMsg protoreflect.MessageDescriptor, fd protoreflect.FieldDescriptor, apiDocsJson *gendoc.Template) (terraformType, error) {
	kind := fd.Kind()

	switch kind {
	case protoreflect.StringKind:
		return tfPrimitive("string"), nil
	case protoreflect.BoolKind:
		return tfPrimitive("bool"), nil
	case protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind,
		protoreflect.FloatKind, protoreflect.DoubleKind:
		return tfPrimitive("number"), nil
	case protoreflect.BytesKind:
		return tfPrimitive("string"), nil
	case protoreflect.EnumKind:
		return tfPrimitive("string"), nil
	case protoreflect.MessageKind:
		// Treat google.protobuf JSON wrapper types as simple strings to avoid deep recursion
		fullName := string(fd.Message().FullName())
		if isWellKnownJsonType(fullName) {
			return tfPrimitive("string"), nil
		}
		return messageToTerraformObject(fd.Message(), fd, apiDocsJson)
	default:
		return nil, fmt.Errorf("unsupported field kind: %v", kind)
	}
}

func messageToTerraformObject(md protoreflect.MessageDescriptor, fd protoreflect.FieldDescriptor, apiDocsJson *gendoc.Template) (terraformType, error) {
	fields := md.Fields()
	obj := tfObject{}

	// Now uses a suffix check on the message name to detect Metadata messages
	shouldSkipVersion := strings.HasSuffix(strings.ToLower(string(md.Name())), "metadata")
	parentFullName := string(md.FullName())

	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		fieldName := string(f.Name())

		// If this is a metadata message, skip the 'version' field
		if shouldSkipVersion && fieldName == "version" {
			continue
		}

		valType, err := fieldDescriptorToTerraformType(f, md, apiDocsJson)
		if err != nil {
			return nil, err
		}
		snakeKey := caseconverter.ToSnakeCase(fieldName)

		desc := findFieldDescription(apiDocsJson, parentFullName, fieldName)
		if desc == "" && f.Kind() == protoreflect.MessageKind {
			desc = findMessageDescription(apiDocsJson, string(f.Message().FullName()))
		}
		if desc == "" {
			desc = fieldDescriptions[fieldName]
			if desc == "" {
				desc = fmt.Sprintf("Description for %s", fieldName)
			}
		}

		obj.fields = append(obj.fields, tfField{
			name:        snakeKey,
			description: desc,
			t:           valType,
		})
	}
	return obj, nil
}

func findMessageDescription(apiDocsJson *gendoc.Template, fullName string) string {
	if apiDocsJson == nil {
		return ""
	}

	for _, f := range apiDocsJson.Files {
		for _, m := range f.Messages {
			if m.FullName == fullName {
				return strings.TrimSpace(m.Description)
			}
		}
	}
	return ""
}

func findFieldDescription(apiDocsJson *gendoc.Template, messageFullName, fieldName string) string {
	if apiDocsJson == nil {
		return ""
	}

	for _, f := range apiDocsJson.Files {
		for _, m := range f.Messages {
			if m.FullName == messageFullName {
				for _, fld := range m.Fields {
					if fld.Name == fieldName {
						return strings.TrimSpace(fld.Description)
					}
				}
			}
		}
	}

	return ""
}

// isWellKnownJsonType returns true for protobuf well-known types representing JSON
// so we map them to primitive string in Terraform variable schema
func isWellKnownJsonType(fullName string) bool {
	switch fullName {
	case "google.protobuf.Struct", "google.protobuf.Value", "google.protobuf.ListValue":
		return true
	default:
		return false
	}
}
