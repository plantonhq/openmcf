locals {
  # The AWS collection name is create-time immutable (3-32 chars,
  # ^[a-z][0-9a-z-]+$) -- metadata.name is the naming basis both engines
  # share. The collection-scoped policies the module renders are all named
  # after it too, and their rules match exactly "collection/<name>" /
  # "index/<name>/..." -- one manifest owns one collection and everything
  # that makes it usable.
  collection_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsOpenSearchServerlessCollection"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # ------------------------------------------------------------------
  # Encryption policy document (always rendered -- AWS rejects
  # CreateCollection without a matching encryption policy). AWS-owned key
  # by default; the referenced customer-managed key otherwise (the
  # document's KmsARN member, per AWS's encryption policy schema).
  # ------------------------------------------------------------------
  use_customer_kms_key = var.spec.encryption != null && try(var.spec.encryption.kms_key_arn, "") != ""

  encryption_policy_rules = [
    {
      ResourceType = "collection"
      Resource     = ["collection/${local.collection_name}"]
    }
  ]

  # Encoded per branch: HCL conditionals demand type-consistent branches,
  # and the two document shapes legitimately differ (KmsARN exists only on
  # the customer-key arm).
  encryption_policy_json = local.use_customer_kms_key ? jsonencode({
    Rules       = local.encryption_policy_rules
    AWSOwnedKey = false
    KmsARN      = var.spec.encryption.kms_key_arn
    }) : jsonencode({
    Rules       = local.encryption_policy_rules
    AWSOwnedKey = true
  })

  # ------------------------------------------------------------------
  # Network policy document (an ARRAY of statements). An omitted
  # spec.network block still renders the PUBLIC posture (the AWS console's
  # easy-create default) -- network "public" is reachability only; data
  # access still requires SigV4 auth plus a data-access rule.
  # ------------------------------------------------------------------
  network_allow_from_public = var.spec.network != null ? coalesce(try(var.spec.network.allow_from_public, true), true) : true
  network_vpc_endpoint_ids  = var.spec.network != null ? try(var.spec.network.vpc_endpoint_ids, []) : []
  network_include_dashboard = var.spec.network != null ? coalesce(try(var.spec.network.include_dashboards, true), true) : true

  network_policy_rules = concat(
    [
      {
        ResourceType = "collection"
        Resource     = ["collection/${local.collection_name}"]
      }
    ],
    local.network_include_dashboard ? [
      {
        ResourceType = "dashboard"
        Resource     = ["collection/${local.collection_name}"]
      }
    ] : []
  )

  # Encoded per branch (the SourceVPCEs member exists only on the private
  # arm -- same type-consistency constraint as the encryption document).
  network_policy_json = local.network_allow_from_public ? jsonencode([
    {
      Rules           = local.network_policy_rules
      AllowFromPublic = true
    }
    ]) : jsonencode([
    {
      Rules           = local.network_policy_rules
      AllowFromPublic = false
      SourceVPCEs     = local.network_vpc_endpoint_ids
    }
  ])

  # ------------------------------------------------------------------
  # Data access policy document (an ARRAY of statements, one per spec
  # rule). Skipped entirely when the manifest declares no rules.
  # ------------------------------------------------------------------
  data_access_policy_document = [
    for rule in var.spec.data_access : {
      Rules = concat(
        length(rule.collection_permissions) > 0 ? [
          {
            ResourceType = "collection"
            Resource     = ["collection/${local.collection_name}"]
            Permission   = rule.collection_permissions
          }
        ] : [],
        length(rule.index_permissions) > 0 ? [
          {
            ResourceType = "index"
            # Default: all indexes of this collection.
            Resource   = [for p in(length(rule.index_patterns) > 0 ? rule.index_patterns : ["*"]) : "index/${local.collection_name}/${p}"]
            Permission = rule.index_permissions
          }
        ] : []
      )
      Principal = rule.principals
    }
  ]

  # ------------------------------------------------------------------
  # Lifecycle (retention) policy document. Skipped when the manifest
  # declares no retention rules (indexes are then retained indefinitely).
  # ------------------------------------------------------------------
  # Rules split by arm and encoded together: the NoMinIndexRetention and
  # MinIndexRetention members carry different types, so a single
  # conditional per rule would trip HCL's branch type-consistency.
  lifecycle_policy_json = jsonencode({
    Rules = concat(
      [
        for rule in var.spec.retention_rules : {
          ResourceType        = "index"
          Resource            = [for p in rule.index_patterns : "index/${local.collection_name}/${p}"]
          NoMinIndexRetention = true
        } if rule.unlimited
      ],
      [
        for rule in var.spec.retention_rules : {
          ResourceType      = "index"
          Resource          = [for p in rule.index_patterns : "index/${local.collection_name}/${p}"]
          MinIndexRetention = rule.min_index_retention
        } if !rule.unlimited
      ]
    )
  })
}
