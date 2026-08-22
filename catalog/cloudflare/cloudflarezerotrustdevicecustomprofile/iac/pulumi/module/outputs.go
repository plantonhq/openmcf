package module

const (
	// OpPolicyId is the exported stack output containing the
	// Cloudflare-assigned identifier of the profile (its policy ID).
	OpPolicyId = "policy_id"

	// OpGatewayUniqueId is the exported stack output containing the
	// Gateway-side identifier Cloudflare assigns the profile.
	OpGatewayUniqueId = "gateway_unique_id"
)
