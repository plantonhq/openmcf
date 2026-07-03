---
title: "Presets"
description: "Ready-to-deploy configuration presets for Role Definition"
type: "preset-list"
componentSlug: "role-definition"
componentTitle: "Role Definition"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-vm-operator-role"
    rank: "01"
    title: "Explicit-Actions Operator Role"
    excerpt: "This preset defines the most common custom-role shape: an explicit list of control-plane actions capturing exactly what an operations team may do -- here, observe and restart existing VMs without any..."
  - slug: "02-blob-data-reader-role"
    rank: "02"
    title: "Data-Plane Role (Blob Auditor)"
    excerpt: "This preset teaches the distinction that trips up most custom-role authors: control-plane `actions` manage Azure resources, data-plane `dataActions` access the data inside them. A role with every..."
  - slug: "03-project-admin-carve-out-role"
    rank: "03"
    title: "Broad Grant with Carve-Outs and Assignable-Scope Control"
    excerpt: "This preset is the classic \"almost-Contributor\" role: a wildcard grant with `notActions` trimming away the authorization plane. Project admins manage every resource in their environments but can..."
---

# Role Definition Presets

Ready-to-deploy configuration presets for Role Definition. Each preset is a complete manifest you can copy, customize, and deploy.
