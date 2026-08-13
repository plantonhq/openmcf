<p align="center">
  <img src="logo.svg" alt="AWS Bedrock AgentCore Tools" width="80"/>
</p>

# AWS Bedrock AgentCore Tools

Create and manage [Amazon Bedrock AgentCore built-in tools](https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/builtin-tools.html)
— managed, sandboxed execution environments agents drive at runtime, as
ONE bundle:

## What Gets Created

- **Browsers** — cloud web browsers (navigation, form filling,
  scraping) with optional session recording to S3, cryptographic
  traffic signing, Chrome enterprise policies, and mTLS client
  certificates from Secrets Manager.
- **Browser profiles** — reusable saved browser state (cookies, logins)
  sessions can start from.
- **Code interpreters** — sandboxes that run model-written code, in
  SANDBOX (no network — the safest for untrusted code), PUBLIC, or VPC
  posture.

Every arm is optional — author the ones this bundle owns. Tools are
free to create; AWS bills per session at runtime.

## No Update — By Design

AWS exposes NO update for any of these resources: every field change
recreates the tool. That is cheap — the tools are configuration shells;
sessions, not tools, carry state.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
