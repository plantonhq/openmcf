package module

const (
	OpKeyId      = "key_id"
	OpKeyArn     = "key_arn"
	OpAliasNames = "alias_names"
	// Keyed by the grant's position in spec.grants -- the same key the
	// Terraform for_each uses, so import recipes resolve identically on
	// both engines.
	OpGrantIds = "grant_ids"
)
