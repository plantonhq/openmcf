locals {
  # Flatten the spec's records for iteration: a record with multiple values
  # becomes one DigitalOcean record per value (same name and type). Record
  # values arrive flattened as plain strings: the Planton orchestrator
  # resolves valueFrom references before Terraform runs. The spec's shared
  # enum value names ARE the DigitalOcean record types (A, AAAA, CNAME, ...),
  # so each record's type wires through directly.
  dns_records = flatten([
    for idx, record in coalesce(var.spec.records, []) : [
      for val_idx, value in record.values : {
        key = "${record.name}-${idx}-${val_idx}"

        name  = record.name
        type  = record.type
        value = value

        # 0 means unset: the ttl attribute is then Computed and DigitalOcean
        # applies its default (1800 seconds).
        ttl = coalesce(record.ttl_seconds, 0) > 0 ? record.ttl_seconds : null

        # Per-type fields carry the spec's presence semantics: unset arrives
        # as null and is simply not sent. Spec CEL rules already guarantee
        # the fields each record type requires are present.
        priority = record.priority
        weight   = record.weight
        port     = record.port
        flags    = record.flags
        tag      = record.tag != null && record.tag != "" ? record.tag : null
      }
    ]
  ])
}
