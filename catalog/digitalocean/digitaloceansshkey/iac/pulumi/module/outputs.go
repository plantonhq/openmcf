package module

const (
	// OpSshKeyId is the numeric key id as a string (the API identity and
	// the only import id -- fingerprints do not import).
	OpSshKeyId = "ssh_key_id"
	// OpFingerprint is the MD5 fingerprint DigitalOcean computed for the key.
	OpFingerprint = "fingerprint"
)
