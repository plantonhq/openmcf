package module

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

// deriveEndpointName maps a resource identity to a stable numeric Vertex AI
// endpoint ID in [1000000000, 9999999999] — always 10 digits, never a
// leading zero, satisfying the API's numeric-name contract.
//
// The algorithm is implemented IDENTICALLY in the Terraform module
// (locals.tf): sha256 of "org/env/name", first 12 hex characters parsed as
// a 48-bit integer, mapped into the 9-billion-wide band. Keep the two in
// lockstep — the whole point is that the same manifest yields the same
// endpoint ID on either engine.
func deriveEndpointName(org, env, name string) string {
	identity := fmt.Sprintf("%s/%s/%s", org, env, name)
	digest := sha256.Sum256([]byte(identity))
	prefix := hex.EncodeToString(digest[:])[:12]
	// 12 hex chars = 48 bits, always within uint64 range; the hex prefix of
	// a sha256 digest cannot fail to parse.
	value, _ := strconv.ParseUint(prefix, 16, 64)
	return strconv.FormatUint(1000000000+value%9000000000, 10)
}
