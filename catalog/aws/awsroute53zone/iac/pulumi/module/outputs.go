package module

// Stack output keys — must stay in lockstep with AwsRoute53ZoneStackOutputs.
const (
	OpZoneId            = "zone_id"
	OpZoneName          = "zone_name"
	OpNameservers       = "nameservers"
	OpPrimaryNameServer = "primary_name_server"
	OpZoneArn           = "zone_arn"
	OpDsRecord          = "ds_record"
	OpDnskeyRecord      = "dnskey_record"
	OpKeySigningKeyTag  = "key_signing_key_tag"
)
