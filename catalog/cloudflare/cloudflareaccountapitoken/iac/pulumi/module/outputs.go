package module

const (
	// OpTokenId is the exported stack output containing the
	// Cloudflare-assigned token ID (the token's identity for management
	// calls, not the credential).
	OpTokenId = "token_id"

	// OpValue is the exported stack output containing the token's secret
	// value -- returned by Cloudflare exactly once, on create (secret-marked).
	OpValue = "value"
)
