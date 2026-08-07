# Computed values for the KubernetesCertificate module. The CR spec rendered
# here is the Terraform twin of the Pulumi module's spec_builder.go — keep
# field mappings and enum translation tables in lockstep.
#
# Proto enum vocabularies are lowercase; the CRD wants its exact casing
# ("RSA", "PKCS1", "Always", "DER", "LegacyRC2", ...) — the lookup maps below
# are the Terraform side of that single translation.
#
# HCL DISCIPLINE (applies to every conditional object in this file):
# conditional entries are written as `key = cond ? value : null` inside ONE
# object literal, pruned with `{ for k, v in {...} : k => v if v != null }`.
# The two tempting alternatives are both broken: `cond ? {...} : {}`
# ternaries fail plan-time type unification when branches carry different
# attributes, and `merge(concat(cond ? [{...}] : [], ...)...)` silently
# UNIFIES primitive-only sibling objects into map(string) — numbers and
# booleans arrive at the API as strings and server-side validation rejects
# the object. The null-prune form preserves every value's type.
#
# Optional nested blocks are read with try(): HCL's && does NOT
# short-circuit, so chained null checks still dereference the null.

locals {
  certificate_name = var.metadata.name
  namespace        = var.spec.namespace
  secret_name      = var.spec.secret_name

  labels = merge(
    {
      "planton.ai/resource"      = "true"
      "planton.ai/resource-name" = var.metadata.name
      "planton.ai/resource-kind" = "KubernetesCertificate"
    },
    var.metadata.id != null && var.metadata.id != "" ? { "planton.ai/resource-id" = var.metadata.id } : {},
    var.metadata.org != null && var.metadata.org != "" ? { "planton.ai/organization" = var.metadata.org } : {},
    var.metadata.env != null && var.metadata.env != "" ? { "planton.ai/environment" = var.metadata.env } : {}
  )

  # ---- enum translation tables (proto vocabulary -> CRD casing) ----------
  algorithm_map = { rsa = "RSA", ecdsa = "ECDSA", ed25519 = "Ed25519" }
  encoding_map  = { pkcs1 = "PKCS1", pkcs8 = "PKCS8" }
  rotation_map  = { always = "Always", never = "Never" }
  pkcs12_profile_map = {
    legacy_rc2 = "LegacyRC2", legacy_des = "LegacyDES", modern2023 = "Modern2023"
  }
  output_format_map = { der = "DER", combined_pem = "CombinedPEM" }

  # ---- issuerRef ----------------------------------------------------------
  issuer_ref = {
    for k, v in {
      group = try(var.spec.issuer_ref.external, null) != null ? var.spec.issuer_ref.external.group : null
      kind = (
        try(var.spec.issuer_ref.cluster_issuer, null) != null ? "ClusterIssuer" :
        try(var.spec.issuer_ref.issuer, null) != null ? "Issuer" :
        var.spec.issuer_ref.external.kind
      )
      name = (
        try(var.spec.issuer_ref.cluster_issuer, null) != null ? var.spec.issuer_ref.cluster_issuer.name :
        try(var.spec.issuer_ref.issuer, null) != null ? var.spec.issuer_ref.issuer.name :
        var.spec.issuer_ref.external.name
      )
    } : k => v if v != null
  }

  # ---- the CR spec --------------------------------------------------------
  certificate_spec = {
    for k, v in {
      secretName = local.secret_name
      issuerRef  = local.issuer_ref

      dnsNames       = length(var.spec.dns_names) > 0 ? var.spec.dns_names : null
      ipAddresses    = length(var.spec.ip_addresses) > 0 ? var.spec.ip_addresses : null
      uris           = length(var.spec.uris) > 0 ? var.spec.uris : null
      emailAddresses = length(var.spec.email_addresses) > 0 ? var.spec.email_addresses : null
      commonName     = var.spec.common_name != "" ? var.spec.common_name : null
      literalSubject = var.spec.literal_subject != "" ? var.spec.literal_subject : null

      subject = try(var.spec.subject, null) == null ? null : {
        for sk, sv in {
          organizations       = length(try(var.spec.subject.organizations, [])) > 0 ? var.spec.subject.organizations : null
          organizationalUnits = length(try(var.spec.subject.organizational_units, [])) > 0 ? var.spec.subject.organizational_units : null
          countries           = length(try(var.spec.subject.countries, [])) > 0 ? var.spec.subject.countries : null
          provinces           = length(try(var.spec.subject.provinces, [])) > 0 ? var.spec.subject.provinces : null
          localities          = length(try(var.spec.subject.localities, [])) > 0 ? var.spec.subject.localities : null
          streetAddresses     = length(try(var.spec.subject.street_addresses, [])) > 0 ? var.spec.subject.street_addresses : null
          postalCodes         = length(try(var.spec.subject.postal_codes, [])) > 0 ? var.spec.subject.postal_codes : null
          serialNumber        = try(var.spec.subject.serial_number, "") != "" ? var.spec.subject.serial_number : null
        } : sk => sv if sv != null
      }

      otherNames = length(var.spec.other_names) > 0 ? [
        for n in var.spec.other_names : { oid = n.oid, utf8Value = n.utf8_value }
      ] : null

      duration              = try(var.spec.duration, null)
      renewBefore           = var.spec.renew_before != "" ? var.spec.renew_before : null
      renewBeforePercentage = try(var.spec.renew_before_percentage, null)

      privateKey = try(var.spec.private_key, null) == null ? null : {
        for pk, pv in {
          algorithm      = try(var.spec.private_key.algorithm, null) != null ? local.algorithm_map[var.spec.private_key.algorithm] : null
          size           = try(var.spec.private_key.size, null)
          encoding       = try(var.spec.private_key.encoding, null) != null ? local.encoding_map[var.spec.private_key.encoding] : null
          rotationPolicy = try(var.spec.private_key.rotation_policy, null) != null ? local.rotation_map[var.spec.private_key.rotation_policy] : null
        } : pk => pv if pv != null
      }

      usages                = length(var.spec.usages) > 0 ? var.spec.usages : null
      encodeUsagesInRequest = var.spec.encode_usages_in_request ? true : null
      isCA                  = var.spec.is_ca ? true : null
      signatureAlgorithm    = try(var.spec.signature_algorithm, "") != "" ? var.spec.signature_algorithm : null

      keystores = try(var.spec.keystores, null) == null ? null : {
        for kk, kv in {
          jks = try(var.spec.keystores.jks, null) == null ? null : {
            for jk, jv in {
              create   = var.spec.keystores.jks.create
              alias    = try(var.spec.keystores.jks.alias, null)
              password = var.spec.keystores.jks.password
            } : jk => jv if jv != null
          }
          pkcs12 = try(var.spec.keystores.pkcs12, null) == null ? null : {
            for pk, pv in {
              create   = var.spec.keystores.pkcs12.create
              password = var.spec.keystores.pkcs12.password
              profile  = try(var.spec.keystores.pkcs12.profile, "") != "" ? local.pkcs12_profile_map[var.spec.keystores.pkcs12.profile] : null
            } : pk => pv if pv != null
          }
        } : kk => kv if kv != null
      }

      additionalOutputFormats = length(var.spec.additional_output_formats) > 0 ? [
        for f in var.spec.additional_output_formats : { type = local.output_format_map[f.type] }
      ] : null

      nameConstraints = try(var.spec.name_constraints, null) == null ? null : {
        for nk, nv in {
          critical = var.spec.name_constraints.critical
          permitted = try(var.spec.name_constraints.permitted, null) == null ? null : {
            for ck, cv in {
              dnsDomains     = length(try(var.spec.name_constraints.permitted.dns_domains, [])) > 0 ? var.spec.name_constraints.permitted.dns_domains : null
              ipRanges       = length(try(var.spec.name_constraints.permitted.ip_ranges, [])) > 0 ? var.spec.name_constraints.permitted.ip_ranges : null
              emailAddresses = length(try(var.spec.name_constraints.permitted.email_addresses, [])) > 0 ? var.spec.name_constraints.permitted.email_addresses : null
              uriDomains     = length(try(var.spec.name_constraints.permitted.uri_domains, [])) > 0 ? var.spec.name_constraints.permitted.uri_domains : null
            } : ck => cv if cv != null
          }
          excluded = try(var.spec.name_constraints.excluded, null) == null ? null : {
            for ck, cv in {
              dnsDomains     = length(try(var.spec.name_constraints.excluded.dns_domains, [])) > 0 ? var.spec.name_constraints.excluded.dns_domains : null
              ipRanges       = length(try(var.spec.name_constraints.excluded.ip_ranges, [])) > 0 ? var.spec.name_constraints.excluded.ip_ranges : null
              emailAddresses = length(try(var.spec.name_constraints.excluded.email_addresses, [])) > 0 ? var.spec.name_constraints.excluded.email_addresses : null
              uriDomains     = length(try(var.spec.name_constraints.excluded.uri_domains, [])) > 0 ? var.spec.name_constraints.excluded.uri_domains : null
            } : ck => cv if cv != null
          }
        } : nk => nv if nv != null
      }

      secretTemplate = try(var.spec.secret_template, null) == null ? null : (
        length(try(var.spec.secret_template.labels, {})) > 0 || length(try(var.spec.secret_template.annotations, {})) > 0
        ) ? {
        for tk, tv in {
          labels      = length(try(var.spec.secret_template.labels, {})) > 0 ? var.spec.secret_template.labels : null
          annotations = length(try(var.spec.secret_template.annotations, {})) > 0 ? var.spec.secret_template.annotations : null
        } : tk => tv if tv != null
      } : null

      revisionHistoryLimit = try(var.spec.revision_history_limit, null)
    } : k => v if v != null
  }
}
