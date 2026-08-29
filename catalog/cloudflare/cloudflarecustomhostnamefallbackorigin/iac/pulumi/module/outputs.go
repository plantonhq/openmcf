package module

const (
	// OpCreatedAt is the exported stack output containing the creation timestamp.
	// No status output: fallback-origin deployment is asynchronous
	// (pending_deployment -> active), and a point-in-time phase is never a
	// stable stack output -- it flips on the first refresh after the
	// transition and re-plans forever.
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
