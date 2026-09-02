<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Model Access" width="80"/>
</p>

# AWS Bedrock Model Access

Manage [Amazon Bedrock model access](https://docs.aws.amazon.com/bedrock/latest/userguide/model-access.html)
declaratively — the marketplace agreement (and, where required, the
account use-case form) that entitles an AWS account to invoke a foundation
model in a region. One component instance per model: the infrastructure
form of the console's "Model access" page.

## What Gets Created

- **A foundation-model agreement** accepting the model's PUBLIC offer.
  The offer token is looked up at deploy time — the manifest names only
  the model. Most offers carry no subscription charge; invocations
  bill per-token regardless.
- **The account use-case form** (optional, `use_case_form`) — required
  once per account before Anthropic model agreements activate.

## Destroy Cancels Access

Destroying the component cancels the agreement — the account loses the
model in that region. The use-case form, however, is account-global and
write-once: AWS provides no delete for it, and it survives the component.

## Quick Start

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBedrockModelAccess
metadata:
  name: mistral-7b-access
spec:
  region: us-west-2
  modelId: mistral.mistral-7b-instruct-v0:2
```

For Anthropic models, add the account's use-case form (once, in exactly
one instance):

```yaml
spec:
  region: us-west-2
  modelId: anthropic.claude-3-5-haiku-20241022-v1:0
  useCaseForm:
    companyName: Example Corp
    companyWebsite: https://example.com
    intendedUsers: Internal employees
    industryOption: Technology
    useCases: Customer support assistant
```

## Spec Reference

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
