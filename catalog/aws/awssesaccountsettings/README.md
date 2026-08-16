<p align="center">
  <img src="logo.svg" alt="AWS SES Account Settings" width="80"/>
</p>

# AWS SES Account Settings

Manage the region's [SES account-level settings](https://docs.aws.amazon.com/ses/latest/dg/sending-email-suppression-list.html)
— the account suppression list and the
[Virtual Deliverability Manager](https://docs.aws.amazon.com/ses/latest/dg/vdm.html)
(VDM) posture.

This is a **settings singleton**: AWS keeps exactly one SES account
object per account+region, so deploy at most one instance per region.
`metadata.name` never reaches AWS. These settings are deliberately not
fields on [AwsSesConfigurationSet](../awssesconfigurationset) or
[AwsSesEmailIdentity](../awssesemailidentity) — multiple sets or
identities would fight over the one account object.

## What Gets Managed

- **The account suppression list**: which events (bounces,
  complaints) automatically suppress recipient addresses across every
  send from the account. An empty list explicitly turns
  auto-suppression off.
- **The VDM posture**: deliverability analytics (engagement
  dashboards, Guardian delivery optimization). VDM carries its own
  AWS pricing — enabling it is a billing decision.

Destroy semantics **differ per arm**: suppression settings persist
after destroy (SES has no delete for them); VDM is reset to disabled.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
