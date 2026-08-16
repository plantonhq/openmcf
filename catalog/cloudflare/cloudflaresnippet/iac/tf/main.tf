# Cloudflare Snippet: a small JavaScript module at the zone's edge, invoked by
# snippet rules (managed by the CloudflareSnippetRules kind). The snippet NAME
# is the identity, and Cloudflare's create is an upsert -- deploying a name that
# already exists in the zone silently adopts and overwrites it. Renames replace
# the resource.
#
# The provider refetches stored content on refresh (rebuilt from the API's
# multipart response), so server-side normalization of the source can read back
# as drift -- keep file content byte-stable.
resource "cloudflare_snippet" "main" {
  zone_id      = var.spec.zone_id
  snippet_name = var.spec.snippet_name

  files = [
    for file in var.spec.files : {
      name    = file.name
      content = file.content
    }
  ]

  metadata = {
    main_module = var.spec.main_module
  }
}
