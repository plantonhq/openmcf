---
title: "Billing"
description: "Subscription plans, prepaid AI credits, and billing management for your Planton organization"
icon: billing
order: 20
tags:
  - Billing
  - Subscriptions
  - Plans
---

# Billing

Every Planton organization has a billing account that determines which capabilities are available. Billing is managed through Stripe and can be configured entirely from the web console. Seats are the only meter: automation runs — infrastructure deployments and CI/CD pipelines — are never billed, and AI assistance runs on prepaid credits you control.

## Subscription Plans

### Free

The starting point, free forever. Up to 3 seats, with everything else — environments, cloud connections, resources, services, and automation runs — unlimited and unmetered. No credit card required, and a free organization can never be surprise-billed.

### Team

A per-seat plan for teams shipping together: $20 per seat per month or $192 per seat per year (₹999/₹9,990 in India). Every seat adds a person; nothing else is metered.

Current prices are always displayed on the **Plans** view in the Billing section of the web console — the price you see there is the price the platform charges, resolved from the same catalog the checkout uses.

<!-- SCREENSHOT: Plans page
  Page: /orgs/{org}/settings/billing
  Action: Show the plan cards (Free, Team)
  Focus: The plan cards showing per-seat pricing
  Alt: Billing page showing the Free and Team plans with per-seat pricing
-->

## AI Credits

AI assistance is funded by prepaid credits, never usage-billed after the fact:

- **Credit packs** — Buy a pack on the Stripe-hosted checkout page; the credit value shown at the click is exactly what lands in your wallet.
- **Spend protection** — Automatic top-up is opt-in and bounded: it requires a threshold, a named pack, and a monthly ceiling you set. When the balance runs out and no top-up applies, AI usage pauses — it never runs up a bill.
- **A transparent ledger** — Every credit movement appears on your wallet statement with a plain-language description.

## Automation Is Never Billed

Infrastructure deployments and service pipelines consume no billable meter. Usage is recorded for platform insight and abuse protection, but it never appears on an invoice — a busy deployment week costs the same as a quiet one.

## Managing Your Subscription

### Viewing Your Plan

Navigate to **Settings → Billing** in the web console. You can see your current plan, standing, payment method on file, and your AI-credit balance and ledger.

### Changing Plans

Purchases and plan changes happen on Stripe's hosted checkout page — your card details never touch Planton. A new subscription becomes live the moment Stripe confirms payment. You can cancel from the billing page; paid time you already bought keeps running until the end of the period.

### Payment and Invoices

Payment is processed through Stripe. From the billing page you can open the Stripe billing portal to update payment methods, view invoices, and download receipts.

## Related Documentation

- [Teams and Access](/docs/teams-and-access) — Member management, teams, and role-based access control
- [Platform Overview](/docs/platform) — Organization structure and core concepts
