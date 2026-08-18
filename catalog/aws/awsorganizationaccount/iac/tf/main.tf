# A MEMBER account of the organization, with its account-level
# settings satellites (alternate contacts, primary contact, opt-in
# region enablement) folded onto the created account's ID.
#
# Lifecycle facts the render below depends on:
#   - email, role_name, iam_user_access_to_billing, and create_govcloud
#     are creation-time facts (the first three force replacement; a
#     create_govcloud change is silently ignored by AWS);
#   - the display name renames in place through the Account Management
#     API; a parent_id change MOVES the account between OUs in place;
#   - role_name has NO read API - imports and drift detection never see
#     it (declared config-only in the import catalog);
#   - destroy honors close_on_deletion: false REMOVES the account from
#     the organization (it survives standalone), true CLOSES it
#     (~90-day PENDING_CLOSURE, quota-limited per rolling 30 days);
#   - the contact satellites use idempotent Put APIs - the provider
#     polls until the write is visible (contacts are eventually
#     consistent); primary-contact delete is a NO-OP (the last-written
#     contact stays on file);
#   - region enable/disable are long operations (up to ~60 minutes
#     each way, 60m timeouts in the provider) and region delete is a
#     NO-OP (the region keeps its last state);
#   - the account-settings satellites require trusted access for AWS
#     Account Management ("account.amazonaws.com") on the organization.

resource "aws_organizations_account" "this" {
  name  = var.spec.account_name
  email = var.spec.email

  parent_id = var.spec.parent_id != "" ? var.spec.parent_id : null
  role_name = var.spec.role_name != "" ? var.spec.role_name : null

  # Empty keeps AWS's default (ALLOW).
  iam_user_access_to_billing = var.spec.iam_user_access_to_billing != "" ? var.spec.iam_user_access_to_billing : null

  close_on_deletion = var.spec.close_on_deletion
  create_govcloud   = var.spec.create_govcloud

  tags = local.aws_tags

  # role_name is write-only at AWS (no read API): without this, any
  # imported account would plan a destructive replacement to "set" a
  # value AWS can never echo back.
  lifecycle {
    ignore_changes = [role_name]
  }
}

# At most one alternate contact per category (BILLING / OPERATIONS /
# SECURITY).
resource "aws_account_alternate_contact" "this" {
  for_each = local.alternate_contacts

  account_id             = aws_organizations_account.this.id
  alternate_contact_type = each.key

  name          = each.value.name
  title         = each.value.title
  email_address = each.value.email_address
  phone_number  = each.value.phone_number
}

# The account's primary contact information. Optional leaves are sent
# only when set - clearing one in the spec leaves the last value on
# file at AWS (the API has no unset semantics).
resource "aws_account_primary_contact" "this" {
  count = var.spec.primary_contact != null ? 1 : 0

  account_id = aws_organizations_account.this.id

  full_name      = var.spec.primary_contact.full_name
  address_line_1 = var.spec.primary_contact.address_line_1
  city           = var.spec.primary_contact.city
  postal_code    = var.spec.primary_contact.postal_code
  country_code   = var.spec.primary_contact.country_code
  phone_number   = var.spec.primary_contact.phone_number

  company_name       = var.spec.primary_contact.company_name != "" ? var.spec.primary_contact.company_name : null
  address_line_2     = var.spec.primary_contact.address_line_2 != "" ? var.spec.primary_contact.address_line_2 : null
  address_line_3     = var.spec.primary_contact.address_line_3 != "" ? var.spec.primary_contact.address_line_3 : null
  district_or_county = var.spec.primary_contact.district_or_county != "" ? var.spec.primary_contact.district_or_county : null
  state_or_region    = var.spec.primary_contact.state_or_region != "" ? var.spec.primary_contact.state_or_region : null
  website_url        = var.spec.primary_contact.website_url != "" ? var.spec.primary_contact.website_url : null
}

# Opt-in region enablement for the member account.
resource "aws_account_region" "this" {
  for_each = local.regions

  account_id  = aws_organizations_account.this.id
  region_name = each.value.region_name
  enabled     = each.value.enabled
}
