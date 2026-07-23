---
title: "Presets"
description: "Ready-to-deploy configuration presets for Manifest"
type: "preset-list"
componentSlug: "manifest"
componentTitle: "Manifest"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-config-bundle"
    rank: "01"
    title: "Config Bundle"
    excerpt: "This preset applies a multi-document manifest — a ConfigMap and a Secret — anchored in a single namespace. Neither document declares its own `metadata.namespace`, so both land in the anchor namespace..."
  - slug: "02-crd-and-custom-resource"
    rank: "02"
    title: "CRD and Custom Resource"
    excerpt: "This preset applies a CustomResourceDefinition and a custom resource of that new type in one manifest — the ordering problem that breaks a naive `kubectl apply -f` (the custom resource is rejected..."
  - slug: "03-vendor-install-manifest"
    rank: "03"
    title: "Vendor Install Manifest"
    excerpt: "This preset is the \"paste the vendor's install YAML\" pattern: take the manifest a project publishes for `kubectl apply -f` — often hundreds of documents spanning CRDs, RBAC, Services, Deployments,..."
---

# Manifest Presets

Ready-to-deploy configuration presets for Manifest. Each preset is a complete manifest you can copy, customize, and deploy.
