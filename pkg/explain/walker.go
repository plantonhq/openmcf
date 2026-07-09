package explain

import (
	"fmt"
	"strings"

	validatepb "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const stringValueOrRefFullName = "dev.planton.shared.foreignkey.v1.StringValueOrRef"

// walkFields renders every field of a message. visited holds the message
// full names on the current descent path so self/mutually-recursive schemas
// terminate (the recursion point is marked instead of expanded).
func (e *Engine) walkFields(md protoreflect.MessageDescriptor, visited map[protoreflect.FullName]bool) []Field {
	fields := md.Fields()
	out := make([]Field, 0, fields.Len())
	for i := 0; i < fields.Len(); i++ {
		out = append(out, e.renderField(fields.Get(i), visited))
	}
	return out
}

func (e *Engine) renderField(fd protoreflect.FieldDescriptor, visited map[protoreflect.FullName]bool) Field {
	f := Field{
		Name:     fd.JSONName(),
		Doc:      e.doc(fd.FullName()),
		Optional: fd.HasOptionalKeyword(),
	}
	for _, interpret := range e.Interpreters {
		interpret(fd, &f)
	}
	applyValidateRules(fd, &f)

	switch {
	case fd.IsMap():
		f.Type = fmt.Sprintf("map<%s, %s>", fd.MapKey().Kind().String(), TypeLabel(fd.MapValue()))
		if valueMsg := fd.MapValue().Message(); valueMsg != nil && expandable(valueMsg) {
			f.Fields, f.Constraints = e.expandMessage(valueMsg, visited, f.Constraints)
		}
	case fd.IsList():
		f.Type = "[]" + TypeLabel(fd)
		f.Fields, f.Constraints = e.expandIfMessage(fd, visited, f.Constraints)
	default:
		f.Type = TypeLabel(fd)
		f.Fields, f.Constraints = e.expandIfMessage(fd, visited, f.Constraints)
	}

	// The foreign-key wrapper's docs describe the concept ("a literal or a
	// reference") but not the YAML serialization, and the wrapper is
	// deliberately not expanded -- so without this line authors guess the
	// shape and guess wrong (a bare string does NOT parse). Spell out both
	// authorable forms, concretized with the field's default reference
	// target when the schema declares one.
	if fd.Message() != nil && fd.Message().FullName() == stringValueOrRefFullName {
		f.Constraints = append(f.Constraints, ValueFromContract(f.RefKind, f.RefFieldPath))
	}

	if fd.Kind() == protoreflect.EnumKind {
		values := fd.Enum().Values()
		for i := 0; i < values.Len(); i++ {
			value := values.Get(i)
			f.Enum = append(f.Enum, EnumValue{
				Name: string(value.Name()),
				Doc:  e.doc(value.FullName()),
			})
		}
	}
	return f
}

// TypeLabel names a field's element type the way a manifest author thinks
// about it. Well-known types collapse to friendly names; the foreign-key
// wrapper is surfaced as its two authorable shapes.
//
// Exported as shared vocabulary: every surface that talks to manifest
// authors about schemas (this package's reference reports, the YAML
// diagnoser's error messages) must name types identically, or the two
// surfaces drift into teaching different languages.
func TypeLabel(fd protoreflect.FieldDescriptor) string {
	if fd.Kind() != protoreflect.MessageKind && fd.Kind() != protoreflect.GroupKind {
		if fd.Kind() == protoreflect.EnumKind {
			return "enum"
		}
		return fd.Kind().String()
	}
	md := fd.Message()
	switch md.FullName() {
	case stringValueOrRefFullName:
		return "string | valueFrom"
	case "google.protobuf.Struct":
		return "object"
	case "google.protobuf.Value":
		return "any"
	case "google.protobuf.ListValue":
		return "[]any"
	case "google.protobuf.Timestamp":
		return "timestamp"
	case "google.protobuf.Duration":
		return "duration"
	case "google.protobuf.StringValue":
		return "string"
	case "google.protobuf.BoolValue":
		return "bool"
	case "google.protobuf.Int32Value", "google.protobuf.Int64Value":
		return "int"
	}
	return string(md.Name())
}

// expandable reports whether a message type should be recursed into. The
// well-known and foreign-key wrappers are terminal: their internal shape is
// framework plumbing, not manifest surface.
func expandable(md protoreflect.MessageDescriptor) bool {
	name := md.FullName()
	return name != stringValueOrRefFullName && !strings.HasPrefix(string(name), "google.protobuf.")
}

func (e *Engine) expandIfMessage(fd protoreflect.FieldDescriptor, visited map[protoreflect.FullName]bool, constraints []string) ([]Field, []string) {
	md := fd.Message()
	if md == nil || !expandable(md) {
		return nil, constraints
	}
	return e.expandMessage(md, visited, constraints)
}

// expandMessage recurses into a nested message, folding the nested message's
// own cross-field CEL messages into the field's constraints so an agent sees
// every rule at the place it authors the value.
func (e *Engine) expandMessage(md protoreflect.MessageDescriptor, visited map[protoreflect.FullName]bool, constraints []string) ([]Field, []string) {
	if visited[md.FullName()] {
		return nil, append(constraints, "recursive: same shape as enclosing "+string(md.Name()))
	}
	visited[md.FullName()] = true
	defer delete(visited, md.FullName())

	for _, rule := range messageCelRules(md) {
		constraints = append(constraints, rule.Message)
	}
	return e.walkFields(md, visited), constraints
}

// applyValidateRules folds the field's buf.validate rules into the report:
// the required flag is derived (explicit `required`, or a min-length/min-items
// floor), CEL rules contribute their human-worded messages, and the remaining
// structured rules are kept verbatim as one compact JSON constraint so no
// constraint class is silently dropped.
func applyValidateRules(fd protoreflect.FieldDescriptor, f *Field) {
	opts, ok := fd.Options().(proto.Message)
	if !ok || opts == nil || !proto.HasExtension(opts, validatepb.E_Field) {
		return
	}
	rules, ok := proto.GetExtension(opts, validatepb.E_Field).(*validatepb.FieldRules)
	if !ok || rules == nil {
		return
	}

	if rules.GetRequired() ||
		rules.GetString().GetMinLen() > 0 ||
		rules.GetRepeated().GetMinItems() > 0 ||
		rules.GetMap().GetMinPairs() > 0 {
		f.Required = true
	}

	for _, cel := range rules.GetCel() {
		if msg := cel.GetMessage(); msg != "" {
			f.Constraints = append(f.Constraints, msg)
		} else {
			f.Constraints = append(f.Constraints, cel.GetExpression())
		}
	}

	structural := proto.Clone(rules).(*validatepb.FieldRules)
	structural.Cel = nil
	if body, err := protojson.Marshal(structural); err == nil && string(body) != "{}" {
		f.Constraints = append(f.Constraints, string(body))
	}
}

// ValueFromContract renders the exact YAML shapes a StringValueOrRef field
// accepts, using the declared default reference target when present.
// Exported for the same shared-vocabulary reason as TypeLabel: reference
// reports and load-error diagnoses must teach the identical authoring shape.
func ValueFromContract(refKind, refFieldPath string) string {
	kind := refKind
	if kind == "" {
		kind = "<Kind>"
	}
	fieldPath := refFieldPath
	if fieldPath == "" {
		fieldPath = "status.outputs.<output>"
	}
	return fmt.Sprintf(
		"write as {value: <literal>} or {valueFrom: {kind: %s, name: <that resource's name>, fieldPath: %s}} -- a bare string does not parse",
		kind, fieldPath)
}

// messageCelRules extracts a message's cross-field CEL rules with their
// authored ids and messages.
func messageCelRules(md protoreflect.MessageDescriptor) []Rule {
	opts, ok := md.Options().(proto.Message)
	if !ok || opts == nil || !proto.HasExtension(opts, validatepb.E_Message) {
		return nil
	}
	rules, ok := proto.GetExtension(opts, validatepb.E_Message).(*validatepb.MessageRules)
	if !ok || rules == nil {
		return nil
	}
	out := make([]Rule, 0, len(rules.GetCel()))
	for _, cel := range rules.GetCel() {
		out = append(out, Rule{Id: cel.GetId(), Message: cel.GetMessage()})
	}
	return out
}
