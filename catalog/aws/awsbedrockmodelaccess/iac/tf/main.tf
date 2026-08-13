# Amazon Bedrock model access: accepts the marketplace agreement that
# entitles this account (in this region) to invoke a foundation model --
# the declarative form of the console's "Model access" page.
#
# The offer token is looked up at deploy time from the model's PUBLIC
# offer, so the manifest only names the model. Destroying the agreement
# cancels the account's access to the model in this region; the use-case
# form (when rendered) is account-global and write-once -- AWS provides no
# delete for it.

# The model's current public offer -- offer tokens are short-lived, so the
# lookup happens on every deploy rather than ever living in a manifest.
data "aws_bedrock_foundation_model_agreement_offers" "this" {
  model_id   = var.spec.model_id
  offer_type = "PUBLIC"
}

# The account-global use-case form. PutUseCaseForModelAccess is
# upsert-convergent (re-putting identical form data is a no-op) and the
# provider imports an existing identical form instead of failing; a
# DIFFERENT existing form fails the deploy loudly -- the form cannot be
# updated through this API.
resource "aws_bedrock_use_case_for_model_access" "this" {
  count = local.has_use_case_form ? 1 : 0

  form_data = jsonencode(var.spec.use_case_form)
}

# The agreement itself. Create waits for the agreement to reach AVAILABLE
# (the provider polls GetFoundationModelAvailability); delete cancels the
# agreement and retries through AWS's transient ConflictException window.
resource "aws_bedrock_foundation_model_agreement" "this" {
  model_id    = var.spec.model_id
  offer_token = data.aws_bedrock_foundation_model_agreement_offers.this.offers[0].offer_token

  # Anthropic agreements activate only once the account's use-case form is
  # on file -- order the agreement after the form whenever one is managed.
  depends_on = [aws_bedrock_use_case_for_model_access.this]
}
