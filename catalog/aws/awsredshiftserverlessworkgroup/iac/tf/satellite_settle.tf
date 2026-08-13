# The destroy-side settle between the satellite groups. Redshift
# Serverless serializes mutating control-plane operations per workgroup,
# and the lock outlives the API call: a usage-limit create/delete
# returns in under a second but flips the workgroup to MODIFYING for
# ~15-30s afterward. Endpoint-access CREATE and DELETE both answer 400
# ConflictException ("An operation is running on the serverless
# workgroup") against that window, and the provider retries the
# conflict only on the workgroup's own delete/update -- never on the
# endpoint access (upstream gap, recorded).
#
# Create order needs no delay -- endpoint accesses apply straight after
# the workgroup (idle from the provider's own wait-for-available) and
# the conflict-immune usage limits apply last. Destroy reverses that
# order, so the usage-limit deletes run FIRST and would leave the
# workgroup busy exactly when the endpoint delete arrives (live-caught)
# -- this sleep sits between them and absorbs the window. It exists
# only when both groups exist (the only shape with the crossing), and
# only its destroy_duration is set, so applies are never delayed.
resource "time_sleep" "endpoint_access_settle" {
  count = length(var.spec.endpoint_accesses) > 0 && length(var.spec.usage_limits) > 0 ? 1 : 0

  destroy_duration = "60s"

  depends_on = [aws_redshiftserverless_endpoint_access.this]
}
