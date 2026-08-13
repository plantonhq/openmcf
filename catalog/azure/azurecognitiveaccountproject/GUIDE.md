# Azure Cognitive Account Project -- Operational Guide

Judgment that saves real time when running AI Foundry projects. The field reference lives in the API Explorer; this is the operational layer above it.

## The first project is special -- plan which one that is

ARM marks the first project created on an account as the account's DEFAULT project (the `is_default` output), and portal and SDK flows that do not name a project land there. Create the intended default first; the marker is ARM's, not yours to set.

## The parent's shape is a prerequisite, not a suggestion

A project create fails unless the account is kind `AIServices` with `projectManagementEnabled: true` -- and enabling project management requires the account to carry a managed identity. When a project create is rejected, fix the ACCOUNT (the AI Foundry Account preset is the right shape); nothing in the project manifest is wrong.

## Never clear description or display name in place

ARM cannot update either field to an empty value -- the module honestly replaces the project when one is cleared, and a replaced project takes its agents, evaluations and files with it. Set both at creation or leave them empty forever; change them freely to other non-empty values.

## The project identity is the grant surface

Agents and evaluations act as the PROJECT's identity, not the account's. Grant storage containers, search indexes and Key Vault secrets to `system_assigned_identity_principal_id` (or bring a pre-granted user-assigned identity when the grants must exist before the project). Granting the account's identity instead is the classic silent-403 mistake.

## Deleting a project deletes its contents

A project is the container for its agents, evaluation runs and files -- teardown removes them with it. Deleting the DEFAULT project while sibling projects exist also leaves the account without a default; recreate deliberately, not as cleanup.
