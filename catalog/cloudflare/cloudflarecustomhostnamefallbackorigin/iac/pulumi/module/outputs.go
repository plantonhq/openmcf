package module

const (
	// OpStatus is the exported stack output containing the fallback origin status.
	OpStatus = "status"
	// OpCreatedAt is the exported stack output containing the creation timestamp.
	OpCreatedAt = "created_at"
	// OpUpdatedAt is the exported stack output containing the last-update timestamp.
	OpUpdatedAt = "updated_at"
	// OpErrors is the exported stack output containing any deployment errors.
	OpErrors = "errors"
	// OpZoneId is the exported stack output containing the Cloudflare zone ID
	// this singleton belongs to (the fallback origin has no resource id; its
	// API identity IS the zone).
	OpZoneId = "zone_id"
)
