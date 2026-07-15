---
title: "Preset: API Scopes"
description: "**Rank**: 1 (most common starting point)"
type: "preset"
rank: "01"
presetSlug: "01-api-scopes"
componentSlug: "cognito-resource-server"
componentTitle: "Cognito Resource Server"
provider: "aws"
icon: "package"
order: 1
---

# Preset: API Scopes

**Rank**: 1 (most common starting point)

## When to Use

- Defining custom OAuth scopes for an API a user pool protects
- The prerequisite for any machine-to-machine (`client_credentials`) client -- M2M clients can only request custom scopes

## What It Provides

- A resource server identified by the API's audience URI (`identifier` is ForceNew -- choose deliberately)
- A `read`/`write` scope vocabulary, requestable by clients as `{identifier}/read` and `{identifier}/write`

## What You Might Add

- Finer-grained scopes (`orders:read`, `orders:write`, `admin`) as the API's permission model grows
- An `AwsCognitoUserPoolClient` with `allowedOauthFlows: [client_credentials]` listing these scopes
