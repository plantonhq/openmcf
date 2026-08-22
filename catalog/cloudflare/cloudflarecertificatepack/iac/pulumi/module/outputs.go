package module

const (
	// OpCertificatePackId is the exported stack output containing the pack id.
	OpCertificatePackId = "certificate_pack_id"
	// OpStatus is the exported stack output containing the order/issuance status.
	OpStatus = "status"
	// OpPrimaryCertificate is the exported stack output containing the primary
	// certificate id.
	OpPrimaryCertificate = "primary_certificate"
	// OpZoneId is the exported stack output containing the Cloudflare zone ID
	// the pack was ordered in (a pack's API identity is zone_id + certificate_pack_id).
	OpZoneId = "zone_id"
)
