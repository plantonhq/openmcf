package module

const (
	// OpCertificatePackId is the exported stack output containing the pack id.
	// No status output: pack issuance is asynchronous (initializing ->
	// pending_validation -> active), and a point-in-time phase is never a
	// stable stack output -- it flips on the first refresh after the
	// transition and re-plans forever.
	OpCertificatePackId = "certificate_pack_id"
	// OpPrimaryCertificate is the exported stack output containing the primary
	// certificate id.
	OpPrimaryCertificate = "primary_certificate"
	// OpZoneId is the exported stack output containing the Cloudflare zone ID
	// the pack was ordered in (a pack's API identity is zone_id + certificate_pack_id).
	OpZoneId = "zone_id"
)
