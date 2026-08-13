package catalogschema

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ExtractFile builds a ProtoFile from a protogen.File with fully-resolved
// descriptors. RawContent is left empty here -- extraction is pure over
// descriptors; the plugin main owns reading the authored source from disk.
func ExtractFile(file *protogen.File) ProtoFile {
	filePkg := file.Desc.Package()

	imports := make([]string, 0, file.Desc.Imports().Len())
	for i := 0; i < file.Desc.Imports().Len(); i++ {
		imports = append(imports, string(file.Desc.Imports().Get(i).Path()))
	}

	fileOptions := extractFileOptions(file)

	messages := make([]ProtoMessage, 0, len(file.Messages))
	for _, msg := range file.Messages {
		if isMapEntry(msg) {
			continue
		}
		messages = append(messages, extractMessage(msg, filePkg))
	}

	enums := make([]ProtoEnum, 0, len(file.Enums))
	for _, enum := range file.Enums {
		enums = append(enums, extractEnum(enum))
	}

	return ProtoFile{
		Syntax:      syntaxString(file.Desc.Syntax()),
		PackageName: string(filePkg),
		Imports:     imports,
		Messages:    messages,
		Enums:       enums,
		Options:     fileOptions,
		RawContent:  "",
	}
}

func syntaxString(s protoreflect.Syntax) string {
	switch s {
	case protoreflect.Proto2:
		return "proto2"
	case protoreflect.Proto3:
		return "proto3"
	default:
		return s.String()
	}
}

// extractFileOptions returns well-known file-level option strings.
// protoreflect doesn't expose raw option text, so options are reconstructed
// from known fields; go_package reflects buf's managed-mode resolution.
func extractFileOptions(file *protogen.File) []string {
	opts := file.Desc.Options()
	if opts == nil {
		return []string{}
	}

	var result []string
	if goPkg := file.GoImportPath; goPkg != "" {
		result = append(result, "option go_package = \""+string(goPkg)+"\";")
	}
	return result
}

// extractMessage recursively extracts a message and its nested types.
func extractMessage(msg *protogen.Message, filePkg protoreflect.FullName) ProtoMessage {
	fields := make([]ProtoField, 0, len(msg.Fields))
	for _, field := range msg.Fields {
		// Skip fields that belong to a oneof (they are handled under oneofs)
		if field.Oneof != nil && !isSyntheticOneof(field.Oneof) {
			continue
		}
		fields = append(fields, extractField(field, filePkg))
	}

	nested := make([]ProtoMessage, 0, len(msg.Messages))
	for _, sub := range msg.Messages {
		if isMapEntry(sub) {
			continue
		}
		nested = append(nested, extractMessage(sub, filePkg))
	}

	nestedEnums := make([]ProtoEnum, 0, len(msg.Enums))
	for _, enum := range msg.Enums {
		nestedEnums = append(nestedEnums, extractEnum(enum))
	}

	oneofs := make([]ProtoOneof, 0, len(msg.Oneofs))
	for _, oo := range msg.Oneofs {
		// Skip synthetic oneofs (proto3 optional desugaring)
		if isSyntheticOneof(oo) {
			continue
		}
		oneofs = append(oneofs, extractOneof(oo, filePkg))
	}

	rawOpts := extractMessageOptions(msg)

	return ProtoMessage{
		Name:           string(msg.Desc.Name()),
		Comment:        cleanComment(msg.Comments.Leading),
		Fields:         fields,
		NestedMessages: nested,
		NestedEnums:    nestedEnums,
		Oneofs:         oneofs,
		RawOptions:     rawOpts,
	}
}

// extractField builds a ProtoField with type resolution and semantic annotations.
func extractField(field *protogen.Field, filePkg protoreflect.FullName) ProtoField {
	fd := field.Desc

	typeName, isMap, mapKeyType, mapValueType := resolveFieldType(fd, filePkg)
	label := resolveFieldLabel(fd)
	sem, rawOpts := extractFieldSemantics(fd)

	return ProtoField{
		Name:         string(fd.Name()),
		Number:       int32(fd.Number()),
		Type:         typeName,
		Label:        label,
		Comment:      cleanComment(field.Comments.Leading),
		RawOptions:   rawOpts,
		IsMap:        isMap,
		MapKeyType:   mapKeyType,
		MapValueType: mapValueType,
		Semantics:    sem,
	}
}

// resolveFieldType returns the type name, map status, and map key/value types.
func resolveFieldType(fd protoreflect.FieldDescriptor, filePkg protoreflect.FullName) (typeName string, isMap bool, mapKeyType string, mapValueType string) {
	if fd.IsMap() {
		keyType := scalarKindName(fd.MapKey().Kind())
		valType := typeNameForKind(fd.MapValue(), filePkg)
		return valType, true, keyType, valType
	}

	return typeNameForKind(fd, filePkg), false, "", ""
}

// typeNameForKind returns a type name string for a field descriptor.
func typeNameForKind(fd protoreflect.FieldDescriptor, filePkg protoreflect.FullName) string {
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return messageTypeName(fd.Message(), filePkg)
	case protoreflect.EnumKind:
		return enumTypeName(fd.Enum(), filePkg)
	default:
		return scalarKindName(fd.Kind())
	}
}

// messageTypeName returns the short name for same-package messages, fully-qualified otherwise.
func messageTypeName(md protoreflect.MessageDescriptor, filePkg protoreflect.FullName) string {
	if md.ParentFile().Package() == filePkg {
		return string(md.Name())
	}
	return string(md.FullName())
}

// enumTypeName returns the short name for same-package enums, fully-qualified otherwise.
func enumTypeName(ed protoreflect.EnumDescriptor, filePkg protoreflect.FullName) string {
	if ed.ParentFile().Package() == filePkg {
		return string(ed.Name())
	}
	return string(ed.FullName())
}

func scalarKindName(k protoreflect.Kind) string {
	switch k {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.Int32Kind:
		return "int32"
	case protoreflect.Sint32Kind:
		return "sint32"
	case protoreflect.Uint32Kind:
		return "uint32"
	case protoreflect.Int64Kind:
		return "int64"
	case protoreflect.Sint64Kind:
		return "sint64"
	case protoreflect.Uint64Kind:
		return "uint64"
	case protoreflect.Sfixed32Kind:
		return "sfixed32"
	case protoreflect.Fixed32Kind:
		return "fixed32"
	case protoreflect.FloatKind:
		return "float"
	case protoreflect.Sfixed64Kind:
		return "sfixed64"
	case protoreflect.Fixed64Kind:
		return "fixed64"
	case protoreflect.DoubleKind:
		return "double"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "bytes"
	default:
		return k.String()
	}
}

// resolveFieldLabel returns "optional", "repeated", or "" for the field's label.
func resolveFieldLabel(fd protoreflect.FieldDescriptor) string {
	if fd.IsList() && !fd.IsMap() {
		return "repeated"
	}
	// Detect explicit proto3 `optional` keyword: represented as a synthetic oneof.
	if oo := fd.ContainingOneof(); oo != nil && oo.IsSynthetic() {
		return "optional"
	}
	return ""
}

// extractOneof builds a ProtoOneof.
func extractOneof(oo *protogen.Oneof, filePkg protoreflect.FullName) ProtoOneof {
	fields := make([]ProtoField, 0, len(oo.Fields))
	for _, field := range oo.Fields {
		fields = append(fields, extractField(field, filePkg))
	}
	return ProtoOneof{
		Name:   string(oo.Desc.Name()),
		Fields: fields,
	}
}

// extractEnum builds a ProtoEnum.
func extractEnum(enum *protogen.Enum) ProtoEnum {
	values := make([]ProtoEnumValue, 0, len(enum.Values))
	for _, val := range enum.Values {
		values = append(values, ProtoEnumValue{
			Name:    string(val.Desc.Name()),
			Number:  int32(val.Desc.Number()),
			Comment: cleanComment(val.Comments.Leading),
		})
	}
	return ProtoEnum{
		Name:    string(enum.Desc.Name()),
		Comment: cleanComment(enum.Comments.Leading),
		Values:  values,
	}
}

// extractMessageOptions returns reconstructed message-level option strings.
// Message-level options (like buf.validate.message CEL rules) are complex to
// reconstruct from resolved descriptors -- a declared fidelity boundary of
// the contract: consumers render them from rawContent.
func extractMessageOptions(_ *protogen.Message) []string {
	return []string{}
}

// isMapEntry returns true for synthetic map-entry messages.
func isMapEntry(msg *protogen.Message) bool {
	return msg.Desc.IsMapEntry()
}

// isSyntheticOneof returns true for compiler-generated oneofs representing proto3 optional fields.
func isSyntheticOneof(oo *protogen.Oneof) bool {
	return oo.Desc.IsSynthetic()
}

// cleanComment trims whitespace and normalises protogen's leading comment text.
func cleanComment(c protogen.Comments) string {
	s := strings.TrimSpace(string(c))
	if s == "" {
		return ""
	}
	// protogen comments preserve leading spaces after "//". Normalise to
	// a block of text with single newlines and no per-line leading space.
	var lines []string
	for _, line := range strings.Split(s, "\n") {
		lines = append(lines, strings.TrimSpace(line))
	}
	return strings.Join(lines, "\n")
}
