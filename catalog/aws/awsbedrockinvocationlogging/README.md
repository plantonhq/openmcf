<p align="center">
  <img src="logo.svg" alt="AWS Bedrock Invocation Logging" width="80"/>
</p>

# AWS Bedrock Invocation Logging

Manage the region's [Bedrock model invocation logging](https://docs.aws.amazon.com/bedrock/latest/userguide/model-invocation-logging.html)
— the audit trail of every model call: full request/response bodies
per data type, delivered to CloudWatch Logs, S3, or both.

This is a **settings singleton**: AWS keeps exactly one invocation
logging configuration per account+region (the resource identity IS
the region), so deploy at most one instance per region.
`metadata.name` never reaches AWS. It is the observability backbone
for every Bedrock workload — without it there is no record of what
prompts were sent or what the models returned.

## What Gets Managed

- **The region's invocation logging configuration**: which data types
  are captured (text, image, embedding, video — AWS defaults all to
  on) and where they are delivered.
- **CloudWatch delivery** through an IAM role trusting
  `bedrock.amazonaws.com`, with optional S3 spillover for payloads
  larger than a log event (256 KB).
- **S3 delivery** authorized by the bucket's policy — Bedrock writes
  as its own service principal, not through the role.

Destroying this component **deletes the configuration** — the region
reverts to no invocation logging. The setting itself is free; you pay
normal CloudWatch/S3 storage for what it delivers.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
