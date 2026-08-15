package conversion

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Canonicalize returns the document's protobuf-faithful form — the ONLY form
// the conversion engine's op paths and envelope restamp speak: JSON field
// names (camelCase), the protojson value representation (64-bit integers as
// strings, enums as names), and only the fields the message actually carries.
// A zero-valued field without explicit presence reads as absent — exactly
// what the document durably holds once stored — while `optional` fields
// survive even at zero, so the engine's presence-based default op operates
// on protobuf presence, by design.
//
// Documents legally arrive in two spellings: stored documents carry
// proto-name keys (api_version — the storage print's contract) while
// authored manifests usually use camelCase, and protobuf JSON accepts both,
// mixed included. The engine is deliberately descriptor-free and cannot
// bridge the spellings itself, so every lane canonicalizes at its entry into
// conversion — feeding a non-canonical document to the engine silently
// no-ops its ops, which is a defect class, never a shortcut.
//
// The parse is strict: a field the schema does not declare is an error
// naming the field. Canonicalization is idempotent — an already-canonical
// document round-trips unchanged.
func Canonicalize(md protoreflect.MessageDescriptor, doc map[string]any) (map[string]any, error) {
	raw, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("the document is not representable as JSON: %w", err)
	}
	msg := dynamicpb.NewMessage(md)
	if err := (protojson.UnmarshalOptions{}).Unmarshal(raw, msg); err != nil {
		return nil, fmt.Errorf("the document does not fit the %s schema: %w", md.FullName(), err)
	}
	canonical, err := protojson.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("canonicalizing the document: %w", err)
	}
	var out map[string]any
	if err := json.Unmarshal(canonical, &out); err != nil {
		return nil, fmt.Errorf("canonicalizing the document: %w", err)
	}
	return out, nil
}
