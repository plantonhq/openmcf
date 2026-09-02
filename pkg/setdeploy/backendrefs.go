package setdeploy

import (
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Backend-resolved value references. `$var/...` and `$secret/...` string
// values resolve ONLY through a Planton backend — the grammar itself has one
// home, server-side. This detector is deliberately prefix detection and
// nothing more: parsing the slug shape here would create a second grammar
// that could drift from the real one. A backendless deploy refuses these
// values by naming every field that carries one.
const (
	varRefPrefix    = "$var/"
	secretRefPrefix = "$secret/"
)

// BackendRefUse is one string value carrying a backend-resolved prefix: the
// dotted field path where it sits and the prefix that identifies its class.
type BackendRefUse struct {
	FieldPath string
	// Prefix is "$var/" or "$secret/".
	Prefix string
}

// CollectBackendRefs walks every populated string value of a loaded manifest
// — singular fields, repeated elements, map values, and strings nested in any
// message container (which includes StringValueOrRef literal arms) — and
// returns each value that begins with a backend-resolved prefix.
func CollectBackendRefs(msg proto.Message) []BackendRefUse {
	var out []BackendRefUse
	walkStrings(msg.ProtoReflect(), "", &out)
	return out
}

func walkStrings(m protoreflect.Message, prefix string, out *[]BackendRefUse) {
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		path := string(fd.Name())
		if prefix != "" {
			path = prefix + "." + path
		}
		switch {
		case fd.IsMap():
			mapValue := fd.MapValue()
			v.Map().Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
				entryPath := path + "." + k.String()
				switch mapValue.Kind() {
				case protoreflect.StringKind:
					recordBackendRef(mv.String(), entryPath, out)
				case protoreflect.MessageKind:
					walkStrings(mv.Message(), entryPath, out)
				}
				return true
			})
		case fd.IsList():
			list := v.List()
			for i := 0; i < list.Len(); i++ {
				elemPath := path + "[" + strconv.Itoa(i) + "]"
				switch fd.Kind() {
				case protoreflect.StringKind:
					recordBackendRef(list.Get(i).String(), elemPath, out)
				case protoreflect.MessageKind:
					walkStrings(list.Get(i).Message(), elemPath, out)
				}
			}
		case fd.Kind() == protoreflect.StringKind:
			recordBackendRef(v.String(), path, out)
		case fd.Kind() == protoreflect.MessageKind:
			walkStrings(v.Message(), path, out)
		}
		return true
	})
}

func recordBackendRef(value, path string, out *[]BackendRefUse) {
	switch {
	case strings.HasPrefix(value, varRefPrefix):
		*out = append(*out, BackendRefUse{FieldPath: path, Prefix: varRefPrefix})
	case strings.HasPrefix(value, secretRefPrefix):
		*out = append(*out, BackendRefUse{FieldPath: path, Prefix: secretRefPrefix})
	}
}
