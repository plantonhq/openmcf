# Anthropic Model with Use-Case Form

This preset enables an Anthropic model, supplying the account use-case
form Anthropic agreements require. The form is ACCOUNT-GLOBAL and
write-once — keep it in exactly one instance (this one); every other
model-access instance in the account omits `useCaseForm`.

## When to Use

- The FIRST Anthropic model enabled in an account
- Accounts managing their Bedrock access fully as code

## Key Configuration Choices

- **Fill the form truthfully** — the fields mirror the Bedrock console's
  first-time Anthropic questionnaire; AWS reviews them.
- **One owner instance.** AWS keeps a single form per account, errors on
  differing re-puts, and provides no delete — a second instance with
  different content fails loudly.

## After Deployment

Further Anthropic models need only the agreement (the marketplace-model
preset shape) — the form is already on file.
