# Organization Admin

This preset runs GuardDuty for the WHOLE organization from the
delegated administrator account: new member accounts enroll
automatically and protection plans roll out org-wide.

## When to Use

- The organization's delegated GuardDuty administrator account (the
  management account delegates once — set
  `organization.admin_account_id` there, or in the console)
- Rolling threat detection out to every current and future account

## What You Get

- Every NEW organization account enrolled automatically, with its own
  detector
- S3 protection enabled for ALL members immediately; runtime
  monitoring for new members as they join

## Customize

- `autoEnableOrganizationMembers: ALL` also enrolls existing accounts
  that predate the delegation
- Add explicit `members` entries only for exceptions or
  non-Organizations accounts (they need the account's root email and
  an invitation)
- The org posture is a PATCH: destroying this component leaves it as
  last applied — flip values to NONE first when standing the
  administration down

## Composing

Run this in the delegated-admin account; member accounts need nothing
(their detectors are created by the enrollment).
