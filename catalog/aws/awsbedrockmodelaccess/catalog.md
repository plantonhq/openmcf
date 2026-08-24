# AWS Bedrock Model Access

Enable foundation models for your account declaratively — the marketplace
agreement per model, with the offer token resolved at deploy time and the
Anthropic use-case form handled where required.

## What Gets Created

- A foundation-model agreement accepting the model's public offer
  (typically no subscription charge; invocations bill per token).
- Optionally, the account's use-case form — the once-per-account
  prerequisite for Anthropic model agreements.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Bedrock marketplace permissions
  (`bedrock:CreateFoundationModelAgreement`,
  `bedrock:ListFoundationModelAgreementOffers`,
  `bedrock:PutUseCaseForModelAccess` and read/delete siblings, plus the
  `aws-marketplace:Subscribe` family).

### AWS Account

- The model available in the target region (agreements are regional).
- For Anthropic models: the use-case form on file (via this component or
  recorded once in the console).

## Deploy

### Console

Create the resource from the AWS catalog, name the model, and deploy.

### CLI

```bash
planton apply -f model-access.yaml
```

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockModelAccess
metadata:
  name: mistral-7b-access
spec:
  region: us-west-2
  modelId: mistral.mistral-7b-instruct-v0:2
```

## Operational Notes

- **One instance per model per region** — charts order model-consuming
  components (agents, throughput purchases) after the access grant via
  the `model_id` output.
- **The use-case form is account-global and write-once.** Keep it in
  exactly ONE instance; AWS errors loudly on a differing re-put and
  provides no delete.
- **Destroy = access revocation.** Removing the component cancels the
  agreement; anything still invoking the model in that region starts
  failing.
- **Many models need no agreement at all.** Auto-enabled models (Amazon
  first-party, and increasingly Mistral and Meta) reject the offers API
  with "Agreement not supported" — probe with
  `aws bedrock list-foundation-model-agreement-offers --model-id <id>`
  before adding an instance; this component is for the models that DO
  carry marketplace offers (Cohere, Anthropic, ...).
